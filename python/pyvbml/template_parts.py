"""Template-part helpers for text components."""

from __future__ import annotations

from dataclasses import dataclass

from .emojis_to_character_codes import emojis_to_character_codes
from .parse_props import parse_props
from .sanitize_special_characters import sanitize_special_characters
from .types import IVBMLTemplatePart, TemplateValue, TemplateWrap, VBMLProps

SPACE_CELL = "{0}"


@dataclass
class TemplatePart:
    """A normalized template part."""

    template: str
    wrap: TemplateWrap


@dataclass
class InlineSegment:
    """A contiguous run of display cells within a word."""

    atomic: bool
    cells: list[str]


@dataclass
class TemplateToken:
    """A wrapping token."""

    kind: str
    segments: list[InlineSegment] | None = None


@dataclass
class LineState:
    """Current rendered line state."""

    line: str
    length: int


def _normalize_template_part(part: IVBMLTemplatePart) -> TemplatePart:
    return TemplatePart(
        template=part["template"],
        wrap="never" if part.get("wrap") == "never" else "normal",
    )


def _normalize_template_value(template: TemplateValue | None) -> list[TemplatePart]:
    if isinstance(template, list):
        return [_normalize_template_part(part) for part in template]

    return [TemplatePart(template=template or "", wrap="normal")]


def _preprocess_template_part(props: VBMLProps, part: TemplatePart) -> TemplatePart:
    return TemplatePart(
        template=sanitize_special_characters(
            parse_props(props, emojis_to_character_codes(part.template))
        ),
        wrap=part.wrap,
    )


def _tokenize_display_cells(template: str) -> list[str]:
    cells: list[str] = []
    index = 0

    while index < len(template):
        current = template[index]
        if current == "{":
            close_index = template.find("}", index + 1)
            if close_index != -1:
                cells.append(template[index : close_index + 1])
                index = close_index + 1
                continue

        if current == " ":
            cells.append(SPACE_CELL)
        else:
            cells.append(current)

        index += 1

    return cells


def _clone_segments(segments: list[InlineSegment]) -> list[InlineSegment]:
    return [InlineSegment(atomic=segment.atomic, cells=segment.cells[:]) for segment in segments]


def _build_template_tokens(parts: list[TemplatePart]) -> list[TemplateToken]:
    tokens: list[TemplateToken] = []
    current_segments: list[InlineSegment] = []

    def push_word_token() -> None:
        nonlocal current_segments
        if not current_segments:
            return

        tokens.append(TemplateToken(kind="word", segments=current_segments))
        current_segments = []

    def push_segment(atomic: bool, cells: list[str]) -> None:
        if not cells:
            return

        current_segments.append(InlineSegment(atomic=atomic, cells=cells[:]))

    def append_normal_cells(cells: list[str]) -> None:
        current_cells: list[str] = []

        def flush_current_cells() -> None:
            nonlocal current_cells
            push_segment(False, current_cells)
            current_cells = []

        for cell in cells:
            if cell == "\n":
                flush_current_cells()
                push_word_token()
                tokens.append(TemplateToken(kind="newline"))
                continue

            if cell == SPACE_CELL:
                flush_current_cells()
                push_word_token()
                tokens.append(TemplateToken(kind="space"))
                continue

            current_cells.append(cell)

        flush_current_cells()

    def append_atomic_cells(cells: list[str]) -> None:
        current_cells: list[str] = []

        def flush_current_cells() -> None:
            nonlocal current_cells
            push_segment(True, current_cells)
            current_cells = []

        for cell in cells:
            if cell == "\n":
                flush_current_cells()
                push_word_token()
                tokens.append(TemplateToken(kind="newline"))
                continue

            current_cells.append(cell)

        flush_current_cells()

    for part in parts:
        cells = _tokenize_display_cells(part.template)
        if part.wrap == "never":
            append_atomic_cells(cells)
        else:
            append_normal_cells(cells)

    push_word_token()

    return tokens


def _measure_word_width(token: TemplateToken) -> int:
    return sum(len(segment.cells) for segment in token.segments or [])


def _word_token_to_string(token: TemplateToken) -> str:
    return "".join("".join(segment.cells) for segment in token.segments or [])


def _take_word_prefix(
    token: TemplateToken,
    max_width: int,
) -> tuple[TemplateToken | None, TemplateToken | None]:
    if max_width <= 0:
        return None, TemplateToken(kind="word", segments=_clone_segments(token.segments or []))

    prefix_segments: list[InlineSegment] = []
    remainder_segments: list[InlineSegment] = []
    consumed_width = 0
    segments = token.segments or []

    def append_remainder(start_index: int) -> None:
        remainder_segments.extend(_clone_segments(segments[start_index:]))

    for index, segment in enumerate(segments):
        segment_width = len(segment.cells)

        if consumed_width + segment_width <= max_width:
            prefix_segments.append(InlineSegment(atomic=segment.atomic, cells=segment.cells[:]))
            consumed_width += segment_width
            continue

        if segment.atomic:
            if consumed_width == 0:
                prefix_segments.append(
                    InlineSegment(atomic=True, cells=segment.cells[:max_width])
                )

                if segment_width > max_width:
                    remainder_segments.append(
                        InlineSegment(atomic=True, cells=segment.cells[max_width:])
                    )
                else:
                    remainder_segments.append(
                        InlineSegment(atomic=True, cells=segment.cells[:])
                    )

                append_remainder(index + 1)
            else:
                remainder_segments.append(
                    InlineSegment(atomic=True, cells=segment.cells[:])
                )
                append_remainder(index + 1)

            return (
                TemplateToken(kind="word", segments=prefix_segments)
                if prefix_segments
                else None,
                TemplateToken(kind="word", segments=remainder_segments)
                if remainder_segments
                else None,
            )

        remaining_width = max_width - consumed_width
        if remaining_width > 0:
            prefix_segments.append(
                InlineSegment(atomic=False, cells=segment.cells[:remaining_width])
            )
            remainder_segments.append(
                InlineSegment(atomic=False, cells=segment.cells[remaining_width:])
            )
        else:
            remainder_segments.append(
                InlineSegment(atomic=False, cells=segment.cells[:])
            )

        append_remainder(index + 1)

        return (
            TemplateToken(kind="word", segments=prefix_segments)
            if prefix_segments
            else None,
            TemplateToken(kind="word", segments=remainder_segments)
            if remainder_segments
            else None,
        )

    return TemplateToken(kind="word", segments=prefix_segments), None


def _chunk_word_tokens(width: int, tokens: list[TemplateToken]) -> list[TemplateToken]:
    chunked_tokens: list[TemplateToken] = []

    for token in tokens:
        if token.kind != "word" or _measure_word_width(token) <= width:
            chunked_tokens.append(token)
            continue

        remaining: TemplateToken | None = token
        while remaining is not None and _measure_word_width(remaining) > width:
            prefix, remainder = _take_word_prefix(remaining, width)
            if prefix is None:
                break

            chunked_tokens.append(prefix)
            remaining = remainder

        if remaining is not None:
            chunked_tokens.append(remaining)

    return chunked_tokens


def _wrap_template_tokens(width: int, tokens: list[TemplateToken]) -> list[str]:
    lines = [LineState(line="", length=0)]
    chunked_tokens = _chunk_word_tokens(width, tokens)

    for index, token in enumerate(chunked_tokens):
        last = lines[-1]
        line_length = last.length
        empty_line = not last.line

        if (
            token.kind == "newline"
            and index > 0
            and chunked_tokens[index - 1].kind == "newline"
        ):
            lines.append(LineState(line="", length=0))
            continue

        if token.kind == "newline":
            lines[-1] = LineState(line=last.line, length=width)
            continue

        if token.kind == "space":
            if line_length + 1 > width:
                continue

            lines[-1] = LineState(line=last.line + SPACE_CELL, length=line_length + 1)
            continue

        word_width = _measure_word_width(token)
        if width >= word_width + line_length and not empty_line:
            lines[-1] = LineState(
                line=last.line + _word_token_to_string(token),
                length=line_length + word_width,
            )
            continue

        lines.append(
            LineState(line=_word_token_to_string(token), length=word_width)
        )

    if not lines[0].line:
        return [line.line for line in lines[1:]]

    return [line.line for line in lines]


def resolve_template_lines(
    width: int,
    props: VBMLProps,
    template: TemplateValue | None,
) -> list[str]:
    """Resolve string or template-part input into wrapped lines."""

    parts = [
        _preprocess_template_part(props, part)
        for part in _normalize_template_value(template)
    ]
    return _wrap_template_tokens(width, _build_template_tokens(parts))
