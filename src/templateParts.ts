import {
  IVBMLTemplatePart,
  TemplateWrap,
  VBMLProps,
} from "./types";

import { emojisToCharacterCodes } from "./emojisToCharacterCodes";
import { parseProps } from "./parseProps";
import { sanitizeSpecialCharacters } from "./sanitizeSpecialCharacters";

interface ResolvedTemplatePart {
  template: string;
  wrap: TemplateWrap;
}

interface InlineSegment {
  atomic: boolean;
  cells: string[];
}

interface TemplateWordToken {
  kind: "word";
  segments: InlineSegment[];
}

interface TemplateSpaceToken {
  kind: "space";
}

interface TemplateNewlineToken {
  kind: "newline";
}

type TemplateToken =
  | TemplateWordToken
  | TemplateSpaceToken
  | TemplateNewlineToken;

interface LineState {
  line: string;
  length: number;
}

const SPACE_CELL = "{0}";

const normalizeTemplatePart = (
  part: IVBMLTemplatePart
): ResolvedTemplatePart => ({
  template: part.template,
  wrap: part.wrap === "never" ? "never" : "normal",
});

const normalizeTemplate = (
  template: string | IVBMLTemplatePart[] | undefined
): ResolvedTemplatePart[] => {
  if (Array.isArray(template)) {
    return template.map(normalizeTemplatePart);
  }

  return [
    {
      template: template || "",
      wrap: "normal",
    },
  ];
};

const preprocessTemplatePart = (
  props: VBMLProps,
  part: ResolvedTemplatePart
): ResolvedTemplatePart => ({
  ...part,
  template: sanitizeSpecialCharacters(
    parseProps(props, emojisToCharacterCodes(part.template))
  ),
});

const tokenizeDisplayCells = (template: string): string[] => {
  const cells: string[] = [];

  for (let index = 0; index < template.length; index += 1) {
    const current = template[index];

    if (current === "{") {
      const closeIndex = template.indexOf("}", index + 1);

      if (closeIndex !== -1) {
        cells.push(template.slice(index, closeIndex + 1));
        index = closeIndex;
        continue;
      }
    }

    if (current === " ") {
      cells.push(SPACE_CELL);
      continue;
    }

    cells.push(current);
  }

  return cells;
};

const cloneSegment = (segment: InlineSegment): InlineSegment => ({
  atomic: segment.atomic,
  cells: [...segment.cells],
});

const cloneSegments = (segments: InlineSegment[]): InlineSegment[] =>
  segments.map(cloneSegment);

const createWordToken = (segments: InlineSegment[]): TemplateWordToken => ({
  kind: "word",
  segments,
});

const buildTemplateTokens = (parts: ResolvedTemplatePart[]): TemplateToken[] => {
  const tokens: TemplateToken[] = [];
  let currentSegments: InlineSegment[] = [];

  const pushWordToken = () => {
    if (currentSegments.length === 0) {
      return;
    }

    tokens.push(createWordToken(currentSegments));
    currentSegments = [];
  };

  const pushSegment = (atomic: boolean, cells: string[]) => {
    if (cells.length === 0) {
      return;
    }

    currentSegments = [
      ...currentSegments,
      {
        atomic,
        cells,
      },
    ];
  };

  const appendNormalCells = (cells: string[]) => {
    let currentCells: string[] = [];

    const flushCurrentCells = () => {
      pushSegment(false, currentCells);
      currentCells = [];
    };

    cells.forEach((cell) => {
      if (cell === "\n") {
        flushCurrentCells();
        pushWordToken();
        tokens.push({ kind: "newline" });
        return;
      }

      if (cell === SPACE_CELL) {
        flushCurrentCells();
        pushWordToken();
        tokens.push({ kind: "space" });
        return;
      }

      currentCells = [...currentCells, cell];
    });

    flushCurrentCells();
  };

  const appendAtomicCells = (cells: string[]) => {
    let currentCells: string[] = [];

    const flushCurrentCells = () => {
      pushSegment(true, currentCells);
      currentCells = [];
    };

    cells.forEach((cell) => {
      if (cell === "\n") {
        flushCurrentCells();
        pushWordToken();
        tokens.push({ kind: "newline" });
        return;
      }

      currentCells = [...currentCells, cell];
    });

    flushCurrentCells();
  };

  parts.forEach((part) => {
    const cells = tokenizeDisplayCells(part.template);

    if (part.wrap === "never") {
      appendAtomicCells(cells);
      return;
    }

    appendNormalCells(cells);
  });

  pushWordToken();

  return tokens;
};

const measureWordWidth = (token: TemplateWordToken): number =>
  token.segments.reduce((width, segment) => width + segment.cells.length, 0);

const wordTokenToString = (token: TemplateWordToken): string =>
  token.segments
    .map((segment) => segment.cells.join(""))
    .join("");

const takeWordPrefix = (
  token: TemplateWordToken,
  maxWidth: number
): {
  prefix: TemplateWordToken | null;
  remainder: TemplateWordToken | null;
} => {
  if (maxWidth <= 0) {
    return {
      prefix: null,
      remainder: createWordToken(cloneSegments(token.segments)),
    };
  }

  const prefixSegments: InlineSegment[] = [];
  const remainderSegments: InlineSegment[] = [];
  let consumedWidth = 0;

  const appendRemainder = (startIndex: number) => {
    remainderSegments.push(...cloneSegments(token.segments.slice(startIndex)));
  };

  for (let index = 0; index < token.segments.length; index += 1) {
    const segment = token.segments[index];
    const segmentWidth = segment.cells.length;

    if (consumedWidth + segmentWidth <= maxWidth) {
      prefixSegments.push(cloneSegment(segment));
      consumedWidth += segmentWidth;
      continue;
    }

    if (segment.atomic) {
      if (consumedWidth === 0) {
        prefixSegments.push({
          atomic: true,
          cells: segment.cells.slice(0, maxWidth),
        });

        if (segmentWidth > maxWidth) {
          remainderSegments.push({
            atomic: true,
            cells: segment.cells.slice(maxWidth),
          });
        } else {
          remainderSegments.push(cloneSegment(segment));
        }

        appendRemainder(index + 1);
      } else {
        remainderSegments.push(cloneSegment(segment));
        appendRemainder(index + 1);
      }

      return {
        prefix: prefixSegments.length
          ? createWordToken(prefixSegments)
          : null,
        remainder: remainderSegments.length
          ? createWordToken(remainderSegments)
          : null,
      };
    }

    const remainingWidth = maxWidth - consumedWidth;

    if (remainingWidth > 0) {
      prefixSegments.push({
        atomic: false,
        cells: segment.cells.slice(0, remainingWidth),
      });

      remainderSegments.push({
        atomic: false,
        cells: segment.cells.slice(remainingWidth),
      });
    } else {
      remainderSegments.push(cloneSegment(segment));
    }

    appendRemainder(index + 1);

    return {
      prefix: prefixSegments.length ? createWordToken(prefixSegments) : null,
      remainder: remainderSegments.length
        ? createWordToken(remainderSegments)
        : null,
    };
  }

  return {
    prefix: createWordToken(prefixSegments),
    remainder: null,
  };
};

const chunkWordTokens = (
  width: number,
  tokens: TemplateToken[]
): TemplateToken[] =>
  tokens.flatMap((token) => {
    if (token.kind !== "word" || measureWordWidth(token) <= width) {
      return [token];
    }

    const chunkedTokens: TemplateWordToken[] = [];
    let remaining: TemplateWordToken | null = token;

    while (remaining && measureWordWidth(remaining) > width) {
      const { prefix, remainder } = takeWordPrefix(remaining, width);

      if (!prefix) {
        break;
      }

      chunkedTokens.push(prefix);
      remaining = remainder;
    }

    if (remaining) {
      chunkedTokens.push(remaining);
    }

    return chunkedTokens;
  });

const wrapTemplateTokens = (width: number, tokens: TemplateToken[]): string[] => {
  const lines = chunkWordTokens(width, tokens).reduce(
    (acc: LineState[], token, index, chunkedTokens) => {
      const lastIndex = acc.length - 1;
      const lineLength = acc[lastIndex].length;
      const emptyLine = !acc[lastIndex].line;

      if (
        token.kind === "newline" &&
        chunkedTokens[index - 1]?.kind === "newline"
      ) {
        return [
          ...acc,
          {
            line: "",
            length: 0,
          },
        ];
      }

      if (token.kind === "newline") {
        const line = {
          line: acc[lastIndex].line,
          length: width,
        };
        const previousLines = acc.slice(0, lastIndex);
        return [...previousLines, line];
      }

      if (token.kind === "space") {
        if (1 + lineLength > width) {
          return acc;
        }

        const line = {
          line: `${acc[lastIndex].line}${SPACE_CELL}`,
          length: lineLength + 1,
        };
        const previousLines = acc.slice(0, lastIndex);
        return [...previousLines, line];
      }

      if (width >= measureWordWidth(token) + lineLength && !emptyLine) {
        const line = {
          line: `${acc[lastIndex].line}${wordTokenToString(token)}`,
          length: lineLength + measureWordWidth(token),
        };
        const previousLines = acc.slice(0, lastIndex);
        return [...previousLines, line];
      }

      return [
        ...acc,
        {
          line: wordTokenToString(token),
          length: measureWordWidth(token),
        },
      ];
    },
    [
      {
        line: "",
        length: 0,
      },
    ]
  );

  const [firstLine, ...rest] = lines;

  if (!firstLine.line) {
    return rest.map(({ line }) => line);
  }

  return lines.map(({ line }) => line);
};

export const resolveTemplateLines = (
  width: number,
  props: VBMLProps,
  template: string | IVBMLTemplatePart[] | undefined
): string[] => {
  const parts = normalizeTemplate(template).map((part) =>
    preprocessTemplatePart(props, part)
  );

  return wrapTemplateTokens(width, buildTemplateTokens(parts));
};
