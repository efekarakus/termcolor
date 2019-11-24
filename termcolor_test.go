package termcolor

import (
	"os"
	"testing"
)

func TestSupportLevel(t *testing.T) {
	testCases := map[string] struct {
		args []string
		envs map[string]string
		isTerminal bool

		wantedLevel Level
	} {
		"with disabled colors: no-color": {
			args: []string{"cli", "--no-color"},
			wantedLevel: LevelNone,
		},
		"with disabled colors: no-colors": {
			args: []string{"cli", "--no-colors"},
			wantedLevel: LevelNone,
		},
		"with disabled colors: color=false": {
			args: []string{"cli", "--color=false"},
			wantedLevel: LevelNone,
		},
		"with disabled colors: color=never": {
			args: []string{"cli", "--color=never"},
			wantedLevel: LevelNone,
		},
		"with true colors: color=16m": {
			args: []string{"cli", "--color=16m"},
			wantedLevel: Level16M,
		},
		"with true colors: color=full": {
			args: []string{"cli", "--color=full"},
			wantedLevel: Level16M,
		},
		"with true colors: color=truecolor": {
			args: []string{"cli", "--color=truecolor"},
			wantedLevel: Level16M,
		},
		"with 256 colors: color=256": {
			args: []string{"cli", "--color=256"},
			wantedLevel: Level256,
		},
		"with a fd that's not a terminal": {
			args: []string{"cli"},
			wantedLevel: LevelNone,
		},
		"with a dumb terminal: FORCE_COLOR=true": {
			args: []string{"cli"},
			envs: map[string]string {
				"FORCE_COLOR": "true",
				"TERM": "dumb",
			},
			isTerminal:  true,
			wantedLevel: LevelBasic,
		},
		"with a dumb terminal: FORCE_COLOR=false": {
			args: []string{"cli"},
			envs: map[string]string {
				"FORCE_COLOR": "false",
				"TERM": "dumb",
			},
			isTerminal: true,
			wantedLevel: LevelNone,
		},
		"with a dumb terminal: FORCE_COLOR is out of bounds": {
			args: []string{"cli"},
			envs: map[string]string {
				"FORCE_COLOR": "123",
				"TERM": "dumb",
			},
			isTerminal:  true,
			wantedLevel: LevelBasic,
		},
		"with a dumb terminal: FORCE_COLOR is within bounds": {
			args: []string{"cli"},
			envs: map[string]string {
				"FORCE_COLOR": "3",
				"TERM": "dumb",
			},
			isTerminal: true,
			wantedLevel: Level16M,
		},
		"with a dumb terminal: color is set": {
			args: []string{"cli", "--color"},
			envs: map[string]string {
				"TERM": "dumb",
			},
			isTerminal:  true,
			wantedLevel: LevelBasic,
		},
		"with travis CI": {
			envs: map[string]string {
				"CI": "true",
				"TRAVIS": "1",
			},
			isTerminal: true,
			wantedLevel: LevelBasic,
		},
		"with circle CI": {
			envs: map[string]string {
				"CI": "true",
				"CIRCLECI": "true",
			},
			isTerminal: true,
			wantedLevel: LevelBasic,
		},
		"with appveyor CI": {
			envs: map[string]string {
				"CI": "true",
				"APPVEYOR": "true",
			},
			isTerminal: true,
			wantedLevel: LevelBasic,
		},
		"with gitlab CI": {
			envs: map[string]string {
				"CI": "true",
				"GITLAB_CI": "true",
			},
			isTerminal: true,
			wantedLevel: LevelBasic,
		},
		"with codeship CI": {
			envs: map[string]string {
				"CI": "true",
				"CI_NAME": "codeship",
			},
			isTerminal: true,
			wantedLevel: LevelBasic,
		},
		"with unknown CI": {
			envs: map[string]string {
				"CI": "true",
				"FORCE_COLOR": "3",
			},
			isTerminal: true,
			wantedLevel: Level16M,
		},
		"with teamcity version < 9.1": {
			envs: map[string]string {
				"TEAMCITY_VERSION": "9.0.5 (build 32523)",
			},
			isTerminal: true,
			wantedLevel: LevelNone,
		},
		"with teamcity version >= 9.1": {
			envs: map[string]string {
				"TEAMCITY_VERSION": "9.1.0 (build 32523)",
			},
			isTerminal: true,
			wantedLevel: LevelBasic,
		},
		"with github actions": {
			envs: map[string]string {
				"GITHUB_ACTIONS": "true",
			},
			isTerminal: true,
			wantedLevel: LevelBasic,
		},
		"with COLORTERM set to truecolor": {
			envs: map[string]string {
				"COLORTERM": "truecolor",
			},
			isTerminal: true,
			wantedLevel: Level16M,
		},
		"level should be 16M when using iTerm 3.0": {
			envs: map[string]string {
				"TERM_PROGRAM": "iTerm.app",
				"TERM_PROGRAM_VERSION": "3.0.10",
			},
			isTerminal: true,
			wantedLevel: Level16M,
		},
		"level should be 256 when using iTerm 2.9": {
			envs: map[string]string {
				"TERM_PROGRAM": "iTerm.app",
				"TERM_PROGRAM_VERSION": "2.9.3",
			},
			isTerminal: true,
			wantedLevel: Level256,
		},
		"level should be 256 on default apple terminal": {
			envs: map[string]string {
				"TERM_PROGRAM": "Apple_Terminal",
			},
			isTerminal: true,
			wantedLevel: Level256,
		},
		"level should be 256 on xterm": {
			envs: map[string]string {
				"TERM": "xterm-256color",
			},
			isTerminal: true,
			wantedLevel: Level256,
		},
		"level should be 256 for screen-256color": {
			envs: map[string]string {
				"TERM": "screen-256color",
			},
			isTerminal: true,
			wantedLevel: Level256,
		},
		"level should be 256 for putty-256color": {
			envs: map[string]string {
				"TERM": "putty-256color",
			},
			isTerminal: true,
			wantedLevel: Level256,
		},
		"level should be basic for colored terminals": {
			envs: map[string]string {
				"COLORTERM": "",
			},
			isTerminal: true,
			wantedLevel: LevelBasic,
		},
	}
	
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Given
			os.Clearenv() // Start the tests from a clean state.
			for k, v := range tc.envs {
				os.Setenv(k, v)
			}
			oldIsTerminal := isTerminal
			oldArgs := args

			isTerminal = mockFalseTty()
			if tc.isTerminal {
				isTerminal = mockTrueTty()
			}
			args = tc.args
			defer func() {
				isTerminal = oldIsTerminal
				args = oldArgs
			}()

			// When
			l := SupportLevel(os.Stdout)

			// Then
			if l != tc.wantedLevel {
				t.Errorf("expected %v, got %v", tc.wantedLevel, l)
			}
		})
	}
}

// TestParityWithChalk tests parity with chalk's supports-color module by
// copying their tests and validating it with our source code.
func TestParityWithChalk(t *testing.T) {
	testCases := map[string] struct {
		args []string
		envs map[string]string
		isTerminal bool

		wantedLevel Level
	} {
		// https://github.com/chalk/supports-color/blob/8a40054cdbcd3f42b4f68eaefb41c3064835b991/test.js#L296
		"return level 2 when FORCE_COLOR is set when not TTY in xterm256": {
			envs: map[string]string{
				"FORCE_COLOR": "true",
				"TERM": "xterm-256color",
			},
			isTerminal: false,
			wantedLevel: Level256,
		},
		// https://github.com/chalk/supports-color/blob/8a40054cdbcd3f42b4f68eaefb41c3064835b991/test.js#L305
		"supports setting a color level using FORCE_COLOR=1": {
			envs: map[string]string{
				"FORCE_COLOR": "1",
			},
			wantedLevel: LevelBasic,
		},
		"supports setting a color level using FORCE_COLOR=2": {
			envs: map[string]string{
				"FORCE_COLOR": "2",
			},
			wantedLevel: Level256,
		},
		"supports setting a color level using FORCE_COLOR=3": {
			envs: map[string]string{
				"FORCE_COLOR": "3",
			},
			wantedLevel: Level16M,
		},
		"supports setting a color level using FORCE_COLOR=0": {
			envs: map[string]string{
				"FORCE_COLOR": "0",
			},
			wantedLevel: LevelNone,
		},
		// https://github.com/chalk/supports-color/blob/8a40054cdbcd3f42b4f68eaefb41c3064835b991/test.js#L359
		"return false when `TERM` is set to dumb when `TERM_PROGRAM` is set": {
			envs: map[string]string{
				"TERM": "dumb",
				"TERM_PROGRAM": "Apple_Terminal",
			},
			isTerminal: true,
			wantedLevel: LevelNone,
		},
		// https://github.com/chalk/supports-color/blob/8a40054cdbcd3f42b4f68eaefb41c3064835b991/test.js#L379
		"return level 1 when `TERM` is set to dumb when `FORCE_COLOR` is set": {
			envs: map[string]string{
				"FORCE_COLOR": "1",
				"TERM": "dumb",
			},
			isTerminal: true,
			wantedLevel: LevelBasic,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Given
			os.Clearenv() // Start the tests from a clean state.
			for k, v := range tc.envs {
				os.Setenv(k, v)
			}
			oldIsTerminal := isTerminal
			oldArgs := args

			isTerminal = mockFalseTty()
			if tc.isTerminal {
				isTerminal = mockTrueTty()
			}
			args = tc.args
			defer func() {
				isTerminal = oldIsTerminal
				args = oldArgs
			}()

			// When
			l := SupportLevel(os.Stdout)

			// Then
			if l != tc.wantedLevel {
				t.Errorf("expected %v, got %v", tc.wantedLevel, l)
			}
		})
	}
}

func TestSupports16M(t *testing.T) {
	testCases := map[string] struct {
		args []string
		envs map[string]string
		isTerminal bool

		wanted bool
	} {
		"level truecolor": {
			args:   []string{"cli", "--color=truecolor"},
			wanted: true,
		},
		"level 256": {
			args: []string{"cli", "--color=256"},
		},
		"level 16": {
			args: []string{"cli"},
			envs: map[string]string {
				"FORCE_COLOR": "true",
				"TERM": "dumb",
			},
			isTerminal:  true,
		},
		"level none": {
			args: []string{"cli"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Given
			os.Clearenv() // Start the tests from a clean state.
			for k, v := range tc.envs {
				os.Setenv(k, v)
			}
			oldIsTerminal := isTerminal
			oldArgs := args

			isTerminal = mockFalseTty()
			if tc.isTerminal {
				isTerminal = mockTrueTty()
			}
			args = tc.args
			defer func() {
				isTerminal = oldIsTerminal
				args = oldArgs
			}()

			// WHEN
			supports := Supports16M(os.Stdout)

			// Then
			if tc.wanted != supports {
				t.Errorf("expected %v, got %v", tc.wanted, supports)
			}
		})
	}
}

func TestSupports256(t *testing.T) {
	testCases := map[string] struct {
		args []string
		envs map[string]string
		isTerminal bool

		wanted bool
	} {
		"level truecolor": {
			args:   []string{"cli", "--color=truecolor"},
			wanted: true,
		},
		"level 256": {
			args: []string{"cli", "--color=256"},
			wanted: true,
		},
		"level 16": {
			args: []string{"cli"},
			envs: map[string]string {
				"FORCE_COLOR": "true",
				"TERM": "dumb",
			},
			isTerminal:  true,
		},
		"level none": {
			args: []string{"cli"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Given
			os.Clearenv() // Start the tests from a clean state.
			for k, v := range tc.envs {
				os.Setenv(k, v)
			}
			oldIsTerminal := isTerminal
			oldArgs := args

			isTerminal = mockFalseTty()
			if tc.isTerminal {
				isTerminal = mockTrueTty()
			}
			args = tc.args
			defer func() {
				isTerminal = oldIsTerminal
				args = oldArgs
			}()

			// WHEN
			supports := Supports256(os.Stdout)

			// Then
			if tc.wanted != supports {
				t.Errorf("expected %v, got %v", tc.wanted, supports)
			}
		})
	}
}

func TestSupportsBasic(t *testing.T) {
	testCases := map[string] struct {
		args []string
		envs map[string]string
		isTerminal bool

		wanted bool
	} {
		"level truecolor": {
			args:   []string{"cli", "--color=truecolor"},
			wanted: true,
		},
		"level 256": {
			args: []string{"cli", "--color=256"},
			wanted: true,
		},
		"level 16": {
			args: []string{"cli"},
			envs: map[string]string {
				"FORCE_COLOR": "true",
				"TERM": "dumb",
			},
			isTerminal:  true,
			wanted: true,
		},
		"level none": {
			args: []string{"cli"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Given
			os.Clearenv() // Start the tests from a clean state.
			for k, v := range tc.envs {
				os.Setenv(k, v)
			}
			oldIsTerminal := isTerminal
			oldArgs := args

			isTerminal = mockFalseTty()
			if tc.isTerminal {
				isTerminal = mockTrueTty()
			}
			args = tc.args
			defer func() {
				isTerminal = oldIsTerminal
				args = oldArgs
			}()

			// WHEN
			supports := SupportsBasic(os.Stdout)

			// Then
			if tc.wanted != supports {
				t.Errorf("expected %v, got %v", tc.wanted, supports)
			}
		})
	}
}

func TestSupportsNone(t *testing.T) {
	testCases := map[string] struct {
		args []string
		envs map[string]string
		isTerminal bool

		wanted bool
	} {
		"level truecolor": {
			args:   []string{"cli", "--color=truecolor"},
			wanted: true,
		},
		"level 256": {
			args: []string{"cli", "--color=256"},
			wanted: true,
		},
		"level 16": {
			args: []string{"cli"},
			envs: map[string]string {
				"FORCE_COLOR": "true",
				"TERM": "dumb",
			},
			isTerminal:  true,
			wanted: true,
		},
		"level none": {
			args: []string{"cli"},
			wanted: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Given
			os.Clearenv() // Start the tests from a clean state.
			for k, v := range tc.envs {
				os.Setenv(k, v)
			}
			oldIsTerminal := isTerminal
			oldArgs := args

			isTerminal = mockFalseTty()
			if tc.isTerminal {
				isTerminal = mockTrueTty()
			}
			args = tc.args
			defer func() {
				isTerminal = oldIsTerminal
				args = oldArgs
			}()

			// WHEN
			supports := SupportsNone(os.Stdout)

			// Then
			if tc.wanted != supports {
				t.Errorf("expected %v, got %v", tc.wanted, supports)
			}
		})
	}
}

func mockFalseTty() func(fd uintptr) bool {
	return func(fd uintptr) bool {
		return false
	}
}

func mockTrueTty() func(fd uintptr) bool {
	return func(fd uintptr) bool {
		return true
	}
}
