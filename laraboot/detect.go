package laraboot

import (
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry/packit"
	"os"
	"path/filepath"
)

func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {

		// The DetectContext includes a WorkingDir field that specifies the
		// location of the application source code. This field can be combined with
		// other paths to find and inspect files included in the application source
		// code that is provided to the buildpack.
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
			fmt.Printf("	--> An error ocurred while parsing laraboot file: '%s'", err)
			return packit.DetectResult{}, fmt.Errorf("invalid laraboot file")
		}

		// Once the laraboot.json file has been parsed, the detect phase can return
		// a result that indicates the provision of xxxxx and the requirement of
		// xxxxx. As can be seen below, the BuildPlanRequirement may also include
		// optional metadata information to such as the source of the version
		// information for a given requirement.

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{Name: "laravel-commander"},
				},
				Requires: []packit.BuildPlanRequirement{
					{
						Name:    "laravel-commander",
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
