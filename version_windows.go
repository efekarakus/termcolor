// +build windows

package termcolor

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

// windowsLevel returns the level of the windows terminal. If the OS is indeed windows,
// the level is returned the boolean is true. If an error occurs then LevelBasic and true is returned.
// If the OS is not windows, then LevelNone and false is returned.
func windowsLevel() (Level, bool) {
	key, err :=registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return LevelBasic, true
	}
	defer key.Close()
	maj, _, err := key.GetIntegerValue("CurrentMajorVersionNumber")
	fmt.Printf("CurrentMajorVersionNumber: %d\n", maj)
	if maj < 10 {
		return LevelBasic, true
	}
	// Windows 10 build 10586 is the first Windows release that supports 256 colors.
	// Windows 10 build 14931 is the first release that supports truecolor.
	cb, _, err := key.GetIntegerValue("CurrentBuild")
	fmt.Printf("CurrentBuild: %d\n", cb)
	if cb < 10586 {
		return LevelBasic, true
	}
	if cb < 14931 {
		return Level256, true
	}
	return Level16M, true
}
