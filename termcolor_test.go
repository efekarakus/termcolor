package termcolor

import (
	"os"
	"testing"
)

func TestSupportLevel(t *testing.T) {
	testCases := map[string] struct {
		envs map[string]string
		isTerminal bool

		wantedLevel Level
	} {
		"with disabled colors: no-color": {
			envs: map[string] string {
				"no-color": "true",
			},
			wantedLevel: LevelNone,
		},
		"with disabled colors: no-colors": {
			envs: map[string] string {
				"no-colors": "true",
			},
			wantedLevel: LevelNone,
		},
		"with disabled colors: color=false": {
			envs: map[string] string {
				"color": "false",
			},
			wantedLevel: LevelNone,
		},
		"with disabled colors: color=never": {
			envs: map[string]string{
				"color": "never",
			},
			wantedLevel: LevelNone,
		},
		"with true colors: color=16m": {
			envs: map[string]string {
				"color": "16m",
			},
			wantedLevel: Level16M,
		},
		"with true colors: color=full": {
			envs: map[string]string {
				"color": "full",
			},
			wantedLevel: Level16M,
		},
		"with true colors: color=truecolor": {
			envs: map[string]string {
				"color": "truecolor",
			},
			wantedLevel: Level16M,
		},
		"with 256 colors: color=256": {
			envs: map[string]string {
				"color": "256",
			},
			wantedLevel: Level256,
		},
		"with a fd that's not a terminal": {
			wantedLevel: LevelNone,
		},
		"with a dumb terminal: FORCE_COLOR=true": {
			envs: map[string]string {
				"FORCE_COLOR": "true",
				"TERM": "dumb",
			},
			isTerminal: true,
			wantedLevel: Level8,
		},
		"with a dumb terminal: FORCE_COLOR=false": {
			envs: map[string]string {
				"FORCE_COLOR": "false",
				"TERM": "dumb",
			},
			isTerminal: true,
			wantedLevel: LevelNone,
		},
		"with a dumb terminal: FORCE_COLOR is out of bounds": {
			envs: map[string]string {
				"FORCE_COLOR": "123",
				"TERM": "dumb",
			},
			isTerminal: true,
			wantedLevel: Level8,
		},
		"with a dumb terminal: FORCE_COLOR is within bounds": {
			envs: map[string]string {
				"FORCE_COLOR": "3",
				"TERM": "dumb",
			},
			isTerminal: true,
			wantedLevel: Level16M,
		},
		"with a dumb terminal: color is set": {
			envs: map[string]string {
				"color": "true",
				"TERM": "dumb",
			},
			isTerminal: true,
			wantedLevel: Level8,
		},
		"with a dumb terminal: colors is set": {
			envs: map[string]string {
				"colors": "true",
				"TERM": "dumb",
			},
			isTerminal: true,
			wantedLevel: Level8,
		},
	}
	
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Given
			for k, v := range tc.envs {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}
			old := isTerminal
			isTerminal = mockFalseTty()
			if tc.isTerminal {
				isTerminal = mockTrueTty()
			}
			defer func() { isTerminal = old }()

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
