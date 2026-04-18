package unit

import (
	"reflect"
	"testing"

	"github.com/KirillRg/cli-tool/internal/profile"
)

func TestBuild_ProfileModes(t *testing.T) {
	testCases := []struct {
		name               string
		input              profile.BuildInput
		expectedMode       profile.LoadProfileMode
		expectedThresholds map[string][]string
	}{
		{
			name: "constant vus profile",
			input: profile.BuildInput{
				InputPath:  "input.yaml",
				OutputPath: "output.js",
				VUs:        10,
				Duration:   "30s",
			},
			expectedMode: profile.ModeConstantVUs,
		},
		{
			name: "shared iterations profile",
			input: profile.BuildInput{
				InputPath:  "input.yaml",
				OutputPath: "output.js",
				VUs:        5,
				Iterations: 100,
			},
			expectedMode: profile.ModeSharedIterations,
		},
		{
			name: "stages profile",
			input: profile.BuildInput{
				InputPath:  "input.yaml",
				OutputPath: "output.js",
				StagesRaw:  []string{"30s:5", "1m:20", "30s:12"},
			},
			expectedMode: profile.ModeStages,
		},
		{
			name: "constant vus profile with thresholds",
			input: profile.BuildInput{
				InputPath:  "input.yaml",
				OutputPath: "output.js",
				VUs:        10,
				Duration:   "30s",
				ThresholdsRaw: []string{
					"http_req_duration=p(95)<500",
					"http_req_failed=rate<0.01",
				},
			},
			expectedMode: profile.ModeConstantVUs,
			expectedThresholds: map[string][]string{
				"http_req_duration": {"p(95)<500"},
				"http_req_failed":   {"rate<0.01"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			loadProfile, err := profile.Build(testCase.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if loadProfile.Mode != testCase.expectedMode {
				t.Fatalf("expected mode %q, got %q", testCase.expectedMode, loadProfile.Mode)
			}

			if testCase.expectedThresholds != nil &&
				!reflect.DeepEqual(loadProfile.Thresholds, testCase.expectedThresholds) {
				t.Fatalf("expected thresholds %#v, got %#v", testCase.expectedThresholds, loadProfile.Thresholds)
			}
		})
	}
}
