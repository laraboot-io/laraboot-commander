// Package larabootcommander .
package larabootcommander

import (
	"github.com/paketo-buildpacks/packit"
)

// Detect fcn.
func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{Name: "laraboot-commander"},
				},
				Requires: []packit.BuildPlanRequirement{
					{
						Name:     "laraboot-commander",
						Metadata: map[string]string{},
					},
				},
			},
		}, nil
	}
}
