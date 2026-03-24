package profile

import (
	"fmt"
	"strconv"
	"strings"
)

func Build(in BuildInput) (*LoadProfile, error) {
	profile := &LoadProfile{
		InputPath:  in.InputPath,
		OutputPath: in.OutputPath,
		Thresholds: make(map[string][]string),
	}

	switch {
	case in.Duration != "":
		profile.Mode = ModeConstantVUs
		profile.VUs = in.VUs
		profile.Duration = in.Duration

	case in.Iterations > 0:
		profile.Mode = ModeSharedIterations
		profile.VUs = in.VUs
		profile.Iterations = in.Iterations

	case len(in.StagesRaw) > 0:
		stages, err := parseStages(in.StagesRaw)
		if err != nil {
			return nil, err
		}

		profile.Mode = ModeStages
		profile.Stages = stages

	default:
		return nil, fmt.Errorf("load profile is not specified")
	}

	thresholds, err := parseThresholds(in.ThresholdsRaw)
	if err != nil {
		return nil, err
	}
	profile.Thresholds = thresholds

	return profile, nil
}

func parseStages(rawStages []string) ([]StageConfig, error) {
	result := make([]StageConfig, 0, len(rawStages))

	for _, rawStage := range rawStages {
		parts := strings.SplitN(rawStage, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("failed to parse stage %q", rawStage)
		}

		duration := strings.TrimSpace(parts[0])
		targetRaw := strings.TrimSpace(parts[1])

		target, err := strconv.Atoi(targetRaw)
		if err != nil {
			return nil, fmt.Errorf("failed to parse stage target in %q: %w", rawStage, err)
		}

		result = append(result, StageConfig{
			Duration: duration,
			Target:   target,
		})
	}

	return result, nil
}

func parseThresholds(rawThresholds []string) (map[string][]string, error) {
	result := make(map[string][]string)

	for _, rawThreshold := range rawThresholds {
		parts := strings.SplitN(rawThreshold, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("failed to parse threshold %q", rawThreshold)
		}

		metric := strings.TrimSpace(parts[0])
		condition := strings.TrimSpace(parts[1])

		result[metric] = append(result[metric], condition)
	}

	return result, nil
}
