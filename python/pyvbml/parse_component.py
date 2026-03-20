"""Parse component.

Port of Vestaboard/vbml/src/parseComponent.ts
"""

from __future__ import annotations

from .character_codes import COLOR_CODES, convert_characters_to_character_codes
from .create_empty_board import create_empty_board
from .horizontal_align import horizontal_align
from .random_colors import random_colors
from .render_component import render_component
from .template_parts import resolve_template_lines
from .types import Align, IVBMLComponent, Justify, VBMLProps
from .vertical_align import vertical_align


def parse_component(
    default_height: int,
    default_width: int,
    props: VBMLProps,
    component: IVBMLComponent,
) -> list[list[int]]:
    """Parse component."""
    if "rawCharacters" in component:
        return component["rawCharacters"]

    style = component.get("style") or {}
    width = style.get("width") or default_width
    height = style.get("height") or default_height

    if "randomColors" in component:
        colors = component["randomColors"].get("colors") or COLOR_CODES
        return random_colors(height, width, colors)

    empty_component = create_empty_board(height, width)
    template = component.get("template", "") if "template" in component else ""
    justify = style.get("justify") or Justify.LEFT
    align = style.get("align") or Align.TOP

    lines = resolve_template_lines(width, props, template)
    codes = [convert_characters_to_character_codes(line) for line in lines]
    codes = vertical_align(height, align, codes)
    codes = horizontal_align(width, justify, codes)
    return render_component(empty_component, codes)


def parse_absolute_component(
    default_height: int,
    default_width: int,
    props: VBMLProps,
    component: IVBMLComponent,
) -> dict:
    """Parse absolute component."""
    style = component.get("style") or {}
    abs_pos = style.get("absolutePosition")
    return {
        "characters": parse_component(default_height, default_width, props, component),
        "x": abs_pos.get("x", 0) if abs_pos else 0,
        "y": abs_pos.get("y", 0) if abs_pos else 0,
    }
