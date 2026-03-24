package ast

import (
	loadprofile "github.com/KirillRg/cli-tool/internal/profile"
)

// Сборка дерева для Профиля Нагрузки
func BuildOptionsAST(profile *loadprofile.LoadProfile) *ObjectExpression {
	if profile == nil {
		return &ObjectExpression{
			Type:       "ObjectExpression",
			Properties: []*Property{},
		}
	}

	var props []*Property

	switch profile.Mode {
	case loadprofile.ModeConstantVUs:
		props = append(props, newInitProperty("vus", &Literal{
			Type:  "Literal",
			Value: profile.VUs,
		}))
		props = append(props, newInitProperty("duration", &Literal{
			Type:  "Literal",
			Value: profile.Duration,
		}))

	case loadprofile.ModeSharedIterations:
		props = append(props, newInitProperty("vus", &Literal{
			Type:  "Literal",
			Value: profile.VUs,
		}))
		props = append(props, newInitProperty("iterations", &Literal{
			Type:  "Literal",
			Value: profile.Iterations,
		}))

	case loadprofile.ModeStages:
		props = append(props, newInitProperty("stages", buildStagesArray(profile.Stages)))
	}

	if len(profile.Thresholds) > 0 {
		props = append(props, newInitProperty("thresholds", buildThresholdsObject(profile.Thresholds)))
	}

	return &ObjectExpression{
		Type:       "ObjectExpression",
		Properties: props,
	}
}

func buildStagesArray(stages []loadprofile.StageConfig) *ArrayExpression {
	elements := make([]Expression, 0, len(stages))

	for _, stage := range stages {
		stageObj := &ObjectExpression{
			Type: "ObjectExpression",
			Properties: []*Property{
				newInitProperty("duration", &Literal{
					Type:  "Literal",
					Value: stage.Duration,
				}),
				newInitProperty("target", &Literal{
					Type:  "Literal",
					Value: stage.Target,
				}),
			},
		}

		elements = append(elements, stageObj)
	}

	return &ArrayExpression{
		Type:     "ArrayExpression",
		Elements: elements,
	}
}

func buildThresholdsObject(thresholds map[string][]string) *ObjectExpression {
	var props []*Property

	for metric, conditions := range thresholds {
		arrayElements := make([]Expression, 0, len(conditions))
		for _, condition := range conditions {
			arrayElements = append(arrayElements, &Literal{
				Type:  "Literal",
				Value: condition,
			})
		}

		props = append(props, &Property{
			Type: "Property",
			Key: &Literal{
				Type:  "Literal",
				Value: metric,
			},
			Value: &ArrayExpression{
				Type:     "ArrayExpression",
				Elements: arrayElements,
			},
			Kind:      "init",
			Method:    false,
			Shorthand: false,
			Computed:  false,
		})
	}

	return &ObjectExpression{
		Type:       "ObjectExpression",
		Properties: props,
	}
}

func newInitProperty(name string, value Expression) *Property {
	return &Property{
		Type:      "Property",
		Key:       &Identifier{Type: "Identifier", Name: name},
		Value:     value,
		Kind:      "init",
		Method:    false,
		Shorthand: false,
		Computed:  false,
	}
}
