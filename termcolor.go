// Package termcolor detects what level of color support your terminal has.
package termcolor

import "os"

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
	// Level8 represents a terminal that can support the basic colors.
	Level8
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

// Supports8 returns true if the file descriptor can support basic colors.
func Supports8(f FileDescriptor) bool {
	return SupportLevel(f) == Level8
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
	return LevelNone
}

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
	return false
}