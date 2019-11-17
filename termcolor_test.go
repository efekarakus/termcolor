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
			args: []string{"cli"},
			envs: map[string]string {
				"CI": "",
				"TRAVIS": "",
			},
			isTerminal: true,
			wantedLevel: LevelBasic,
		},
		"with circle CI": {
			args: []string{"cli"},
			envs: map[string]string {
				"CI": "",
				"CIRCLECI": "",
			},
			isTerminal: true,
			wantedLevel: LevelBasic,
		},
		"with appveyor CI": {
			args: []string{"cli"},
			envs: map[string]string {
				"CI": "",
				"APPVEYOR": "",
			},
			isTerminal: true,
			wantedLevel: LevelBasic,
		},
		"with gitlab CI": {
			args: []string{"cli"},
			envs: map[string]string {
				"CI": "",
				"GITLAB_CI": "",
			},
			isTerminal: true,
			wantedLevel: LevelBasic,
		},
		"with codeship CI": {
			args: []string{"cli"},
			envs: map[string]string {
				"CI": "",
				"CI_NAME": "codeship",
			},
			isTerminal: true,
			wantedLevel: LevelBasic,
		},
		"with unknown CI": {
			args: []string{"cli"},
			envs: map[string]string {
				"CI": "",
				"FORCE_COLOR": "3",
			},
			isTerminal: true,
			wantedLevel: Level16M,
		},
	}
	
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Given
			for k, v := range tc.envs {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
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
