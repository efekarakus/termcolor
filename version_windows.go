// +build windows

package termcolor

import (
	"strconv"

	"golang.org/x/sys/windows/registry"
)

// lookupWindows returns the level of the windows terminal. If the OS is windows, the terminal level is returned and
// the boolean is set to true. If an error occurs then LevelBasic and true is returned.
// If the OS is not windows, then LevelNone and false is returned.
func windowsLevel() (Level, bool) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return LevelBasic, true
	}
	defer key.Close()
	maj, _, err := key.GetIntegerValue("CurrentMajorVersionNumber")
	if err != nil {
		return LevelBasic, true
	}
	if maj < 10 {
		return LevelBasic, true
	}
	// Windows 10 build 10586 is the first Windows release that supports 256 colors.
	// Windows 10 build 14931 is the first release that supports truecolor.
	cb, _, err := key.GetStringValue("CurrentBuild")
	cbv, err := strconv.Atoi(cb)
	if err != nil {
		return LevelBasic, true
	}
	if cbv < 10586 {
		return LevelBasic, true
	}
	if cbv < 14931 {
		return Level256, true
	}
	return Level16M, true
}
