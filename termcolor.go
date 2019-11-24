// Package termcolor detects what level of color support your terminal has.
package termcolor

import (
	"os"
	"regexp"
	"strconv"
	"strings"

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
	return SupportLevel(f) >= Level256
}

// SupportsBasic returns true if the file descriptor can support the basic 16 colors.
func SupportsBasic(f FileDescriptor) bool {
	return SupportLevel(f) >= LevelBasic
}

// SupportNone returns true if the file descriptor cannot support colors.
func SupportsNone(f FileDescriptor) bool {
	return SupportLevel(f) >= LevelNone
}

// SupportLevel returns the color level that's supported by the file descriptor.
// If the environment variables set no color, then returns LevelNone.
func SupportLevel(f FileDescriptor) Level {
	// Flags take priority over anything else.
	if hasDisabledFlag() {
		return LevelNone
	}
	if has16MFlag() {
		return Level16M
	}
	if has256Flag() {
		return Level256
	}

	if !isTerminal(f.Fd())  {
		// If the user forces colors proceed even though it's not a terminal.
		if _, ok := os.LookupEnv("FORCE_COLOR"); !ok {
			return LevelNone
		}
	}

	min := minLevel()
	// Retrieve color from environment variables.
	if isDumbTerminal() {
		return min
	}
	if l, isWindows := lookupWindows(); isWindows {
		return l
	}
	if l, isCI := lookupCI(min); isCI {
		return l
	}
	if isTrueColorTerminal() {
		return Level16M
	}
	if l, isMacOS := lookupMacOS(); isMacOS {
		return l
	}
	if is256Terminal() {
		return Level256
	}
	if isBasicTerminal() {
		return LevelBasic
	}
	return min
}


// Point to dependencies for testing.
var isTerminal = isatty.IsTerminal

func hasDisabledFlag() bool {
	if hasFlag("no-color") {
		return true
	}
	if hasFlag("no-colors") {
		return true
	}
	if hasFlag("color=false") {
		return true
	}
	return hasFlag("color=never")
}

func has16MFlag() bool {
	if hasFlag("color=16m") {
		return true
	}
	if hasFlag("color=full") {
		return true
	}
	return hasFlag("color=truecolor")
}

func has256Flag() bool {
	return hasFlag("color=256")
}

func minLevel() Level {
	if _, ok := os.LookupEnv("FORCE_COLOR"); ok {
		return forceColorValue()
	}
	if hasFlag("color") {
		return LevelBasic
	}
	if hasFlag("colors") {
		return LevelBasic
	}
	if hasFlag("color=true") {
		return LevelBasic
	}
	if hasFlag("color=always") {
		return LevelBasic
	}
	return LevelNone
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

func isDumbTerminal() bool {
	return os.Getenv("TERM") == "dumb"
}

func isTrueColorTerminal() bool {
	return os.Getenv("COLORTERM") == "truecolor"
}

// colored256Screen matches terminals containing "-256" or "-256color".
var colored256Screen = regexp.MustCompile(`-256(color)`)

func is256Terminal() bool {
	return colored256Screen.MatchString(os.Getenv("TERM"))
}

// coloredScreen matches other well known basic colored terminals.
var coloredScreen = regexp.MustCompile(`^screen|^xterm|^vt100|^vt220|^rxvt|color|ansi|cygwin|linux`)

func isBasicTerminal() bool {
	if coloredScreen.MatchString(os.Getenv("TERM")) {
		return true
	}
	_, isColored := os.LookupEnv("COLORTERM")
	return isColored
}

// teamCityVersion matches if the version is greater than 9.1.0.
var teamCityVersion = regexp.MustCompile(`^(9\.(0*[1-9]\d*)\.|\d{2,}\.)`)

func lookupCI(min Level) (Level, bool) {
	if _, ok := os.LookupEnv("TEAMCITY_VERSION"); ok {
		if teamCityVersion.MatchString(os.Getenv("TEAMCITY_VERSION")) {
			return LevelBasic, true
		}
		return LevelNone, true
	}
	if _, ok := os.LookupEnv("GITHUB_ACTIONS"); ok {
		return LevelBasic, true
	}

	// Other CI products set the env CI=true.
	if _, ok := os.LookupEnv("CI"); !ok {
		return LevelNone, false
	}
	if _, ok := os.LookupEnv("TRAVIS"); ok {
		return LevelBasic, true
	}
	if _, ok := os.LookupEnv("CIRCLECI"); ok {
		return LevelBasic, true
	}
	if _, ok := os.LookupEnv("APPVEYOR"); ok {
		return LevelBasic, true
	}
	if _, ok := os.LookupEnv("GITLAB_CI"); ok {
		return LevelBasic, true
	}
	if os.Getenv("CI_NAME") == "codeship" {
		return LevelBasic, true
	}
	return min, true
}

func lookupMacOS() (Level, bool) {
	prog, isMacOS := os.LookupEnv("TERM_PROGRAM")
	if !isMacOS {
		return LevelNone, false
	}
	switch prog {
	case "iTerm.app":
		// Default is 0 if can't convert to integer.
		v, _ := strconv.Atoi(strings.Split(os.Getenv("TERM_PROGRAM_VERSION"), ".")[0])
		if v >= 3 {
			return Level16M, true
		}
		return Level256, true
	case "Apple_Terminal":
		return Level256, true
	default:
		return LevelNone, false
	}
}