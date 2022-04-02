<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# hexes

```go
import "github.com/hhhhhhhhhn/hexes"
```

## Index

- [Variables](<#variables>)
- [type Attribute](<#type-attribute>)
  - [func Join(attributes ...Attribute) Attribute](<#func-join>)
  - [func TrueColor(red, green, blue int) Attribute](<#func-truecolor>)
  - [func TrueColorBg(red, green, blue int) Attribute](<#func-truecolorbg>)
- [type Renderer](<#type-renderer>)
  - [func New(in io.Reader, out io.Writer) *Renderer](<#func-new>)
  - [func (r *Renderer) End()](<#func-renderer-end>)
  - [func (r *Renderer) MoveCursor(row, col int)](<#func-renderer-movecursor>)
  - [func (r *Renderer) NewAttribute(attributes ...Attribute) Attribute](<#func-renderer-newattribute>)
  - [func (r *Renderer) Refresh()](<#func-renderer-refresh>)
  - [func (r *Renderer) Set(row, col int, value rune)](<#func-renderer-set>)
  - [func (r *Renderer) SetAttribute(attribute Attribute)](<#func-renderer-setattribute>)
  - [func (r *Renderer) SetDefaultAttribute(attribute Attribute)](<#func-renderer-setdefaultattribute>)
  - [func (r *Renderer) SetString(row, col int, value string)](<#func-renderer-setstring>)
  - [func (r *Renderer) Start()](<#func-renderer-start>)


## Variables

```go
var (
    NORMAL     = Attribute("\033[0m")
    BOLD       = Attribute("\033[1m")
    FAINT      = Attribute("\033[2m")
    ITALIC     = Attribute("\033[3m")
    UNDERLINE  = Attribute("\033[4m")
    SLOW_BLINK = Attribute("\033[5m")
    FAST_BLINK = Attribute("\033[6m")
    REVERSE    = Attribute("\033[7m")
    STRIKE     = Attribute("\033[8m")

    BLACK   = Attribute("\033[30m")
    RED     = Attribute("\033[31m")
    GREEN   = Attribute("\033[32m")
    YELLOW  = Attribute("\033[33m")
    BLUE    = Attribute("\033[34m")
    MAGENTA = Attribute("\033[35m")
    CYAN    = Attribute("\033[36m")
    WHITE   = Attribute("\033[37m")

    BG_BLACK   = Attribute("\033[40m")
    BG_RED     = Attribute("\033[41m")
    BG_GREEN   = Attribute("\033[42m")
    BG_YELLOW  = Attribute("\033[43m")
    BG_BLUE    = Attribute("\033[44m")
    BG_MAGENTA = Attribute("\033[45m")
    BG_CYAN    = Attribute("\033[46m")
    BG_WHITE   = Attribute("\033[47m")
)
```

## type Attribute

Type representing the style of the text being written \(e\.g\. color\, font weight\)\. Internally\, just command sequences \(a\.k\.a\. ANSI escape codes\)\.

```go
type Attribute []byte
```

### func Join

```go
func Join(attributes ...Attribute) Attribute
```

Create an attribute joining the given\.

### func TrueColor

```go
func TrueColor(red, green, blue int) Attribute
```

Creates an attribute that sets the foreground to the RGB value\.

### func TrueColorBg

```go
func TrueColorBg(red, green, blue int) Attribute
```

Creates an attribute that sets the background to the RGB value\.

## type Renderer

Gives an abstraction to render text in any position in the terminal\.

```go
type Renderer struct {
    Lines            [][]rune      // (RO) The virtual output characters.
    Attributes       [][]Attribute // (RO) The virtual output Attributes.
    Rows             int           // (RO) The amount of rows in virtual output.
    Cols             int           // (RO) The amount of columns in virtual output.
    CursorRow        int           // (RO) The row the cursor is in the terminal.
    CursorCol        int           // (RO) The column the cursor is in the terminal.
    CurrentAttribute Attribute     // (RO) The terminal's current attribute.
    DefaultAttribute Attribute     // (RO) The preferred default attribute of the renderer.
    Out              io.Writer     // (RO) The default writer to send terminal data to.
    In               io.Reader     // (RO) The default reader to get data from.
}
```

### func New

```go
func New(in io.Reader, out io.Writer) *Renderer
```

Creates a new Renderer using in \(e\.g\. os\.Stdin\) as input and out \(e\.g\. os\.Stdout\) as output\.

### func \(\*Renderer\) End

```go
func (r *Renderer) End()
```

Restores terminal to default state\.

### func \(\*Renderer\) MoveCursor

```go
func (r *Renderer) MoveCursor(row, col int)
```

Moves terminal cursor to a position \(0 indexed\)\.

### func \(\*Renderer\) NewAttribute

```go
func (r *Renderer) NewAttribute(attributes ...Attribute) Attribute
```

Creates a new attribute based on the default attribute\.

### func \(\*Renderer\) Refresh

```go
func (r *Renderer) Refresh()
```

Redraws virtual output to the terminal\, handling resizes\.

### func \(\*Renderer\) Set

```go
func (r *Renderer) Set(row, col int, value rune)
```

Sets the cell at row\, col \(0 indexed\) to the character given\.

### func \(\*Renderer\) SetAttribute

```go
func (r *Renderer) SetAttribute(attribute Attribute)
```

Sets the Attribute of the text being written\.

### func \(\*Renderer\) SetDefaultAttribute

```go
func (r *Renderer) SetDefaultAttribute(attribute Attribute)
```

Sets the preferred Attribute with which to prepend new ones\.

### func \(\*Renderer\) SetString

```go
func (r *Renderer) SetString(row, col int, value string)
```

Sets the cells starting at row\, col \(0 indexed\) to value\, accounting for wide characters\.

### func \(\*Renderer\) Start

```go
func (r *Renderer) Start()
```

Initializes and configures the Renderer and user's terminal\.



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)