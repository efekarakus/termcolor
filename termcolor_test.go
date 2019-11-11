package termcolor

import (
	"os"
	"testing"
)

func TestSupportLevel(t *testing.T) {
	testCases := map[string] struct {
		f FileDescriptor
		envs map[string]string

		wantedLevel Level
	} {
		"with disabled colors: no-color": {
			f: os.Stdout,
			envs: map[string] string {
				"no-color": "true",
			},
			wantedLevel: LevelNone,
		},
		"with disabled colors: no-colors": {
			f: os.Stdout,
			envs: map[string] string {
				"no-colors": "true",
			},
			wantedLevel: LevelNone,
		},
		"with disabled colors: color=false": {
			f: os.Stdout,
			envs: map[string] string {
				"color": "false",
			},
			wantedLevel: LevelNone,
		},
		"with disabled colors: color=never": {
			f: os.Stdout,
			envs: map[string]string{
				"color": "never",
			},
			wantedLevel: LevelNone,
		},
	}
	
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Given
			for k, v := range tc.envs {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}
			// When
			l := SupportLevel(tc.f)

			// Then
			if l != tc.wantedLevel {
				t.Errorf("expected %v, got %v", tc.wantedLevel, l)
			}
		})
	}
}
