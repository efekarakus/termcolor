// Package termcolor detects what level of color support your terminal has.
package termcolor

import (
	"os"
	"strconv"

	"github.com/mattn/go-isatty"
)

// FileDescriptor is the interface that wraps the file descriptor method.
type FileDescriptor interface {
	Fd() uintptr
}

// Level represents the number of colors the terminal supports.
type Level int

// Color levels that can be supported by a terminal.
// See https://en.wikipedia.org/wiki/ANSI_escape_code#Colors
const (
	// LevelNone represents a terminal that does not support colors.
	LevelNone Level = iota
	// LevelBasic represents a terminal that can support the basic 16 colors.
	LevelBasic
	// Level256 represents a terminal that can support 256 colors.
	Level256
	// Level16M represents a terminal that can support "true colors".
	Level16M
)

// Supports16M returns true if the file descriptor can support true colors.
func Supports16M(f FileDescriptor) bool {
	return SupportLevel(f) == Level16M
}

// Supports256 returns true if the file descriptor can support 256 colors.
func Supports256(f FileDescriptor) bool {
	return SupportLevel(f) == Level256
}

// SupportsBasic returns true if the file descriptor can support the basic 16 colors.
func SupportsBasic(f FileDescriptor) bool {
	return SupportLevel(f) == LevelBasic
}

// SupportNone returns true if the file descriptor cannot support colors.
func SupportsNone(f FileDescriptor) bool {
	return SupportLevel(f) == LevelNone
}

// SupportLevel returns the color level that's supported by the file descriptor.
// If the environment variables set no color, then returns LevelNone.
func SupportLevel(f FileDescriptor) Level {
	if hasDisabledEnv() {
		return LevelNone
	}
	if has16MEnv() {
		return Level16M
	}
	if has256Env() {
		return Level256
	}
	if !isTerminal(f.Fd()) {
		return LevelNone
	}
	min := minLevel()
	if isDumbTerminal() {
		return min
	}
	return LevelNone
}

// Point to dependency for testing.
var isTerminal = isatty.IsTerminal

func hasDisabledEnv() bool {
	if os.Getenv("no-color") != "" {
		return true
	}
	if os.Getenv("no-colors") != "" {
		return true
	}
	c := os.Getenv("color")
	return c == "false" || c == "never"
}

func has16MEnv() bool {
	c := os.Getenv("color")
	if c == "16m" {
		return true
	}
	if c == "full" {
		return true
	}
	return c == "truecolor"
}

func has256Env() bool {
	return os.Getenv("color") == "256"
}

func minLevel() Level {
	if len(os.Getenv("FORCE_COLOR")) > 0 {
		return forceColorValue()
	}
	if len(os.Getenv("color")) > 0 {
		return LevelBasic
	}
	if len(os.Getenv("colors")) > 0 {
		return LevelBasic
	}
	return LevelNone
}

func isDumbTerminal() bool {
	return os.Getenv("TERM") == "dumb"
}

func forceColorValue() Level {
	fc := os.Getenv("FORCE_COLOR")
	if fc == "true" {
		return LevelBasic
	}
	if fc == "false" {
		return LevelNone
	}
	num, err := strconv.Atoi(fc)
	if err != nil {
		// If not a number then return basic colors.
		return LevelBasic
	}
	switch l := Level(num); l {
	case LevelNone:
		return LevelNone
	case Level256:
		return Level256
	case Level16M:
		return Level16M
	default:
		// If the number is out of bounds default to basic.
		return LevelBasic
	}
}