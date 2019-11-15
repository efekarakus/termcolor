// +build !windows

package termcolor

// windowsLevel returns the level of the windows terminal. If the OS is windows, the terminal level is returned and
// the boolean is set to true. If an error occurs then LevelBasic and true is returned.
// If the OS is not windows, then LevelNone and false is returned.
func windowsLevel() (Level, bool) {
	return LevelNone, false
}
