// +build !windows

package termcolor

// windowsLevel returns the level of the windows terminal. If the OS is indeed windows,
// the level is returned the boolean is true. Otherwise, LevelNone is returned with a false boolean.
func windowsLevel() (Level, bool) {
	return LevelNone, false
}
