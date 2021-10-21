// Package laraboot .
package Larabootcommander

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/packit"
)

// Detect fcn.
func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		file, err := os.Open(filepath.Join(context.WorkingDir, "laraboot.json"))

		if err != nil {
			fmt.Printf("Spec file '%s' was not found", filepath.Join(context.WorkingDir, "laraboot.json"))
			return packit.DetectResult{}, fmt.Errorf("laraboot file not found")
		}

		var config struct {
			PhpConfig struct {
				Version string `json:"version"`
			} `json:"php"`
		}

		err = json.NewDecoder(file).Decode(&config)
		if err != nil {
			fmt.Printf("	--> An error occurred while parsing laraboot file: '%s'", err)
			return packit.DetectResult{}, fmt.Errorf("invalid laraboot file")
		}

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{Name: "laraboot-commander"},
				},
				Requires: []packit.BuildPlanRequirement{
					{
						Name:    "laraboot-commander",
						Version: config.PhpConfig.Version,
						Metadata: map[string]string{
							"version-source": "laraboot.json",
						},
					},
				},
			},
		}, nil
	}
}
