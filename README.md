# VBML - Vestaboard Markup Language

## Installation

### JavaScript / TypeScript

```bash
yarn install @vestaboard/vbml
```

or

```bash
npm i @vestaboard/vbml
```

### Python

```bash
pip install pyvbml
```

### PHP

```bash
composer require vestaboard/vbml
```

### Go

```bash
go get github.com/brenank/vbml/go
```

## Usage

### JavaScript / TypeScript

```typescript
import { vbml } from "@vestaboard/vbml";

// Generate an array of 6 rows of 22 character codes representing the template
const characters = vbml.parse({
  components: [
    {
      style: {
        justify: "center",
        align: "center",
      },
      template: "Hello World!",
    },
  ],
});
```

### Python

```python
from pyvbml import vbml

# Generate an array of 6 rows of 22 character codes representing the template
characters = vbml.parse({
    "components": [
        {
            "style": {
                "justify": "center",
                "align": "center",
            },
            "template": "Hello World!",
        }
    ]
})
```

### PHP

```php
use Vestaboard\Vbml\Vbml;

// Generate an array of 6 rows of 22 character codes representing the template
$characters = Vbml::parse([
    'components' => [
        [
            'style' => [
                'justify' => 'center',
                'align' => 'center',
            ],
            'template' => 'Hello World!',
        ],
    ],
]);
```

### Go

```go
package main

import (
    "log"

    vbml "github.com/brenank/vbml/go"
)

func main() {
    characters, err := vbml.Parse(vbml.Input{
        Components: []vbml.Component{
            {
                Style: &vbml.ComponentStyle{
                    Justify: vbml.JustifyCenter,
                    Align:   vbml.AlignCenter,
                },
                Template: "Hello World!",
            },
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    _ = characters
}
```

## Template Parts

The shared JSON contract accepts either a legacy string template or an array of
template parts:

```json
{
  "components": [
    {
      "style": {
        "width": 22,
        "height": 4,
        "justify": "center",
        "align": "top"
      },
      "template": [
        { "template": "rainfall amounts " },
        { "template": "{{amount}}", "wrap": "never" },
        { "template": " possible." }
      ]
    }
  ]
}
```

TypeScript, Python, and PHP accept that `template` shape directly.

Template-part boundaries do not insert spaces or line breaks. Each part is
processed independently for emoji conversion, Mustache props, and special
character sanitization before wrapping.

Use `wrap: "never"` to keep a resolved part atomic during wrapping. If an
atomic part is wider than the available line width, VBML falls back to
splitting it at display-cell boundaries instead of overflowing or failing.

### Go

Go keeps the canonical JSON contract at the serialization boundary, but exposes
an idiomatic runtime API:

```go
component := vbml.Component{
    Style: &vbml.ComponentStyle{
        Width:  22,
        Height: 4,
    },
    TemplateParts: []vbml.TemplatePart{
        {Template: "rainfall amounts "},
        {Template: "{{amount}}", Wrap: vbml.TemplateWrapNever},
        {Template: " possible."},
    },
}
```

`Template` and `TemplateParts` are mutually exclusive in the Go runtime API.

## Docs

Full documentation is available at [https://docs.vestaboard.com/docs/vbml](https://docs.vestaboard.com/docs/vbml)
