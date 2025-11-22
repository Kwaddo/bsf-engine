package entities

import "context"

var (
	CTX = context.Background()

	DefaultSpeeds = map[string]string{
		"motorway":      "120",
		"trunk":         "100",
		"primary":       "90",
		"secondary":     "80",
		"tertiary":      "70",
		"unclassified":  "60",
		"residential":   "50",
		"service":       "30",
		"living_street": "20",
	}
)
