<?php

namespace Vestaboard\Vbml;

/**
 * Helpers for normalized template-part wrapping.
 */
class TemplateParts
{
    private const SPACE_CELL = '{0}';

    public static function resolveTemplateLines(
        int $width,
        array $props,
        string $template
    ): array {
        $parts = array_map(
            fn(array $part) => self::preprocessTemplatePart($props, $part),
            self::normalizeLegacyTemplate($template)
        );

        return self::wrapTemplateTokens($width, self::buildTemplateTokens($parts));
    }

    private static function normalizeLegacyTemplate(string $template): array
    {
        return [[
            'template' => $template,
            'wrap' => 'normal',
        ]];
    }

    private static function preprocessTemplatePart(array $props, array $part): array
    {
        return [
            'template' => SanitizeSpecialCharacters::sanitize(
                ParseProps::parse($props, EmojisToCharacterCodes::convert($part['template']))
            ),
            'wrap' => $part['wrap'],
        ];
    }

    private static function tokenizeDisplayCells(string $template): array
    {
        $cells = [];
        $characters = mb_str_split($template);
        $count = count($characters);

        for ($index = 0; $index < $count; $index++) {
            $current = $characters[$index];

            if ($current === '{') {
                $closeIndex = $index + 1;
                while ($closeIndex < $count && $characters[$closeIndex] !== '}') {
                    $closeIndex++;
                }

                if ($closeIndex < $count) {
                    $cells[] = implode('', array_slice(
                        $characters,
                        $index,
                        $closeIndex - $index + 1
                    ));
                    $index = $closeIndex;
                    continue;
                }
            }

            if ($current === ' ') {
                $cells[] = self::SPACE_CELL;
                continue;
            }

            $cells[] = $current;
        }

        return $cells;
    }

    private static function cloneSegment(array $segment): array
    {
        return [
            'atomic' => $segment['atomic'],
            'cells' => array_values($segment['cells']),
        ];
    }

    private static function cloneSegments(array $segments): array
    {
        return array_map(fn(array $segment) => self::cloneSegment($segment), $segments);
    }

    private static function createWordToken(array $segments): array
    {
        return [
            'kind' => 'word',
            'segments' => $segments,
        ];
    }

    private static function buildTemplateTokens(array $parts): array
    {
        $tokens = [];
        $currentSegments = [];

        $pushWordToken = function () use (&$tokens, &$currentSegments): void {
            if (empty($currentSegments)) {
                return;
            }

            $tokens[] = self::createWordToken($currentSegments);
            $currentSegments = [];
        };

        $pushSegment = function (bool $atomic, array $cells) use (&$currentSegments): void {
            if (empty($cells)) {
                return;
            }

            $currentSegments[] = [
                'atomic' => $atomic,
                'cells' => array_values($cells),
            ];
        };

        $appendNormalCells = function (array $cells) use (&$tokens, $pushSegment, $pushWordToken): void {
            $currentCells = [];

            $flushCurrentCells = function () use (&$currentCells, $pushSegment): void {
                $pushSegment(false, $currentCells);
                $currentCells = [];
            };

            foreach ($cells as $cell) {
                if ($cell === "\n") {
                    $flushCurrentCells();
                    $pushWordToken();
                    $tokens[] = ['kind' => 'newline'];
                    continue;
                }

                if ($cell === self::SPACE_CELL) {
                    $flushCurrentCells();
                    $pushWordToken();
                    $tokens[] = ['kind' => 'space'];
                    continue;
                }

                $currentCells[] = $cell;
            }

            $flushCurrentCells();
        };

        $appendAtomicCells = function (array $cells) use (&$tokens, $pushSegment, $pushWordToken): void {
            $currentCells = [];

            $flushCurrentCells = function () use (&$currentCells, $pushSegment): void {
                $pushSegment(true, $currentCells);
                $currentCells = [];
            };

            foreach ($cells as $cell) {
                if ($cell === "\n") {
                    $flushCurrentCells();
                    $pushWordToken();
                    $tokens[] = ['kind' => 'newline'];
                    continue;
                }

                $currentCells[] = $cell;
            }

            $flushCurrentCells();
        };

        foreach ($parts as $part) {
            $cells = self::tokenizeDisplayCells($part['template']);
            if ($part['wrap'] === 'never') {
                $appendAtomicCells($cells);
                continue;
            }

            $appendNormalCells($cells);
        }

        $pushWordToken();

        return $tokens;
    }

    private static function measureWordWidth(array $token): int
    {
        return array_reduce(
            $token['segments'],
            fn(int $width, array $segment) => $width + count($segment['cells']),
            0
        );
    }

    private static function wordTokenToString(array $token): string
    {
        return implode('', array_map(
            fn(array $segment) => implode('', $segment['cells']),
            $token['segments']
        ));
    }

    private static function takeWordPrefix(array $token, int $maxWidth): array
    {
        if ($maxWidth <= 0) {
            return [
                'prefix' => null,
                'remainder' => self::createWordToken(
                    self::cloneSegments($token['segments'])
                ),
            ];
        }

        $prefixSegments = [];
        $remainderSegments = [];
        $consumedWidth = 0;
        $segments = $token['segments'];

        $appendRemainder = function (int $startIndex) use (&$remainderSegments, $segments): void {
            $remainderSegments = array_merge(
                $remainderSegments,
                self::cloneSegments(array_slice($segments, $startIndex))
            );
        };

        foreach ($segments as $index => $segment) {
            $segmentWidth = count($segment['cells']);

            if ($consumedWidth + $segmentWidth <= $maxWidth) {
                $prefixSegments[] = self::cloneSegment($segment);
                $consumedWidth += $segmentWidth;
                continue;
            }

            if ($segment['atomic']) {
                if ($consumedWidth === 0) {
                    $prefixSegments[] = [
                        'atomic' => true,
                        'cells' => array_slice($segment['cells'], 0, $maxWidth),
                    ];

                    if ($segmentWidth > $maxWidth) {
                        $remainderSegments[] = [
                            'atomic' => true,
                            'cells' => array_slice($segment['cells'], $maxWidth),
                        ];
                    } else {
                        $remainderSegments[] = self::cloneSegment($segment);
                    }

                    $appendRemainder($index + 1);
                } else {
                    $remainderSegments[] = self::cloneSegment($segment);
                    $appendRemainder($index + 1);
                }

                return [
                    'prefix' => !empty($prefixSegments)
                        ? self::createWordToken($prefixSegments)
                        : null,
                    'remainder' => !empty($remainderSegments)
                        ? self::createWordToken($remainderSegments)
                        : null,
                ];
            }

            $remainingWidth = $maxWidth - $consumedWidth;
            if ($remainingWidth > 0) {
                $prefixSegments[] = [
                    'atomic' => false,
                    'cells' => array_slice($segment['cells'], 0, $remainingWidth),
                ];
                $remainderSegments[] = [
                    'atomic' => false,
                    'cells' => array_slice($segment['cells'], $remainingWidth),
                ];
            } else {
                $remainderSegments[] = self::cloneSegment($segment);
            }

            $appendRemainder($index + 1);

            return [
                'prefix' => !empty($prefixSegments)
                    ? self::createWordToken($prefixSegments)
                    : null,
                'remainder' => !empty($remainderSegments)
                    ? self::createWordToken($remainderSegments)
                    : null,
            ];
        }

        return [
            'prefix' => self::createWordToken($prefixSegments),
            'remainder' => null,
        ];
    }

    private static function chunkWordTokens(int $width, array $tokens): array
    {
        $chunkedTokens = [];

        foreach ($tokens as $token) {
            if ($token['kind'] !== 'word' || self::measureWordWidth($token) <= $width) {
                $chunkedTokens[] = $token;
                continue;
            }

            $remaining = $token;
            while ($remaining !== null && self::measureWordWidth($remaining) > $width) {
                $split = self::takeWordPrefix($remaining, $width);
                if ($split['prefix'] === null) {
                    break;
                }

                $chunkedTokens[] = $split['prefix'];
                $remaining = $split['remainder'];
            }

            if ($remaining !== null) {
                $chunkedTokens[] = $remaining;
            }
        }

        return $chunkedTokens;
    }

    private static function wrapTemplateTokens(int $width, array $tokens): array
    {
        $lines = [[
            'line' => '',
            'length' => 0,
        ]];
        $chunkedTokens = self::chunkWordTokens($width, $tokens);

        foreach ($chunkedTokens as $index => $token) {
            $lastIndex = count($lines) - 1;
            $lineLength = $lines[$lastIndex]['length'];
            $emptyLine = ($lines[$lastIndex]['line'] === '');

            if (
                $token['kind'] === 'newline'
                && $index > 0
                && $chunkedTokens[$index - 1]['kind'] === 'newline'
            ) {
                $lines[] = [
                    'line' => '',
                    'length' => 0,
                ];
                continue;
            }

            if ($token['kind'] === 'newline') {
                $lines[$lastIndex] = [
                    'line' => $lines[$lastIndex]['line'],
                    'length' => $width,
                ];
                continue;
            }

            if ($token['kind'] === 'space') {
                if ($lineLength + 1 > $width) {
                    continue;
                }

                $lines[$lastIndex] = [
                    'line' => $lines[$lastIndex]['line'] . self::SPACE_CELL,
                    'length' => $lineLength + 1,
                ];
                continue;
            }

            $wordWidth = self::measureWordWidth($token);
            if ($width >= $wordWidth + $lineLength && !$emptyLine) {
                $lines[$lastIndex] = [
                    'line' => $lines[$lastIndex]['line'] . self::wordTokenToString($token),
                    'length' => $lineLength + $wordWidth,
                ];
                continue;
            }

            $lines[] = [
                'line' => self::wordTokenToString($token),
                'length' => $wordWidth,
            ];
        }

        if ($lines[0]['line'] === '') {
            $lines = array_slice($lines, 1);
        }

        return array_map(fn(array $line) => $line['line'], $lines);
    }
}
