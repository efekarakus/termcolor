# Termcolor [![Actions Status](https://github.com/efekarakus/termcolor/workflows/Go/badge.svg)](https://github.com/efekarakus/termcolor/actions)
Detects what level of color support your terminal has.
This package is heavily inspired by [chalk's support-color](https://github.com/chalk/supports-color) module.

## Install
```sh
go get github.com/efekarakus/termcolor
```

## Examples
Colorize output by finding out which level of color your terminal support:
```go
switch l := termcolor.SupportLevel(os.Stderr); l {
case termcolor.Level16M:
    // wrap text with 24 bit color https://en.wikipedia.org/wiki/ANSI_escape_code#24-bit
case termcolor.Level256:
    // wrap text with 8 bit color https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit
case termcolor.LevelBasic:
    // wrap text with 3/4 bit color https://en.wikipedia.org/wiki/ANSI_escape_code#3/4_bit
case default:
    // no color, return text as is.
}
```

Alternatively, you can use:
```go
if termcolor.Supports16M(os.Stderr) {}
if termcolor.Supports256(os.Stderr) {}
if termcolor.SupportsBasic(os.Stderr) {}
if termcolor.SupportsNone(os.Stderr) {}
```

## Credits
* [Efe Karakus](https://www.efekarakus.com/)
* [chalk/supports-color](https://github.com/chalk/supports-color/)

## License
The MIT License (MIT) - see [LICENSE](https://github.com/efekarakus/termcolor/blob/master/LICENSE) for more details.
