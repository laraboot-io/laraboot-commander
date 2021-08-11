package laraboot

import (
	"encoding/json"
	"fmt"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/cargo"
	"github.com/paketo-buildpacks/packit/chronos"
	"github.com/paketo-buildpacks/packit/pexec"
	"github.com/paketo-buildpacks/packit/postal"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var larabootStruct struct {
	Version   string `json:"version"`
	ProjectID string `json:"project_id"`
	Php       struct {
		Version string `json:"version"`
	} `json:"php"`
	Framework struct {
		Config struct {
			Overrides []struct {
				Key     string `json:"key"`
				Envs    string `json:"envs"`
				Default string `json:"default"`
			} `json:"overrides"`
		} `json:"config"`
		Auth struct {
			Stack string `json:"stack"`
		} `json:"auth"`
		Models []struct {
			Name    string `json:"name"`
			Columns []struct {
				Name string `json:"name"`
				Type string `json:"type"`
			} `json:"columns"`
		} `json:"models"`
	} `json:"Framework"`
}

func Build(logger LogEmitter, clock chronos.Clock) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {

		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)
		logger.Process("Reading laraboot.json")

		thisLayer, blueprintGenErr := context.Layers.Get("laravel-model")
		if blueprintGenErr != nil {
			return packit.BuildResult{}, blueprintGenErr
		}

		lfile, errReadingJsonFile := os.Open(filepath.Join(context.WorkingDir, "laraboot.json"))
		if errReadingJsonFile != nil {
			return packit.BuildResult{}, blueprintGenErr
		}

		errDecoding := json.NewDecoder(lfile).Decode(&larabootStruct)

		if errDecoding != nil {
			fmt.Printf("	--> An error ocurred while parsing laraboot file: '%s'", blueprintGenErr)

		}

		// ---- Buildpack read process
		yamlFile, yamlError := ioutil.ReadFile(filepath.Join(context.WorkingDir, "buildpack.yml"))

		var m struct {
			LaravelModel struct {
				WorkingDir       string `yaml:"directory"`
				RectorArgs       string `yaml:"args"`
				BlueprintVersion string `yaml:"blueprint-version"`
				Git              struct {
					Enabled bool `yaml:"enabled"`
					Commit  bool `yaml:"commit"`
				} `yaml:"git"`
				Clean bool `yaml:"cleanup"`
			} `yaml:"laravel-model"`
		}
		//_, blueprintGenErr = toml.DecodeReader(file, &m)
		if yamlError != nil {
			return packit.BuildResult{}, yamlError
		}

		unmarshallErr := yaml.Unmarshal(yamlFile, &m)

		if unmarshallErr != nil {
			fmt.Printf("unmarshallErr: %v", unmarshallErr)
			return packit.BuildResult{}, yamlError
		}

		// ---- Ends Buildpack read process

		// --- Dependencies
		// Using pack tools to install dependency
		transport := cargo.NewTransport()
		dependencyService := postal.NewService(transport)
		//entryResolver := draft.NewPlanner()

		dependency, dependencyErr := dependencyService.Resolve(filepath.Join(context.CNBPath, "buildpack.yml"), "blueprintgen", "default", context.Stack)
		if dependencyErr != nil {
			return packit.BuildResult{}, dependencyErr
		}
		//logger.SelectedDependency(entry, dependency.Version)
		binPath := fmt.Sprintf("%s/bin", thisLayer.Path)
		logger.Subprocess("Installing BlueprintGen %s %s into %s", dependency.Version, dependency.SHA256, binPath)

		duration, blueprintGenErr := clock.Measure(func() error {
			return dependencyService.Deliver(dependency, context.CNBPath, binPath, "/platform")
		})

		if blueprintGenErr != nil {
			return packit.BuildResult{}, blueprintGenErr
		}
		logger.Action("Completed in %s", duration.Round(time.Millisecond))
		logger.Break()

		logger.Process("Configuring environment")
		thisLayer.SharedEnv.Append("PATH", binPath, ":")
		blueprintBin := fmt.Sprintf("%s/blueprint-gen", binPath)
		thisLayer.SharedEnv.Default("BLUEPRINTGEN_BIN", blueprintBin)
		logger.Environment(thisLayer.SharedEnv)
		// expanding path and setting BLUEPRINTGEN_BIN for runtime use
		envErr := os.Setenv("PATH", fmt.Sprintf("%s:%s", os.Getenv("PATH"), binPath))
		if envErr != nil {
			return packit.BuildResult{}, blueprintGenErr
		}
		envErr = os.Setenv("BLUEPRINTGEN_BIN", blueprintBin)
		if envErr != nil {
			return packit.BuildResult{}, blueprintGenErr
		}

		lsOutput, lserr := LsCommand(context, binPath)
		fmt.Printf("Listing bin folder: %s", lsOutput)
		if lserr != nil {
			return packit.BuildResult{}, lserr
		}

		// --- End Dependencies

		customIni := fmt.Sprintf("%s/php.ini", thisLayer.Path)

		_ = ProcessTemplateToFile("extension=openssl\nextension=mbstring\nextension=fileinfo\nextension=curl",
			customIni,
			"")

		blueprintVersion := m.LaravelModel.BlueprintVersion

		// Use blueprint to generate code for models and controllers
		blueprintFQN := "laravel-shift/blueprint"

		// use specific blueprint version
		if blueprintVersion != "" {
			blueprintFQN = fmt.Sprintf("laravel-shift/blueprint:%s", blueprintVersion)
		}

		// Use larased to tweak config files
		larasedFQN := "oscarnevarezleal/laravel-sed"
		CleanUp := m.LaravelModel.Clean

		composerDevDependencies := append([]string{}, blueprintFQN)
		globalDependencies := append([]string{}, larasedFQN)

		for k, v := range composerDevDependencies {
			logger.Process("Installing dev dependency [%d]: %s", k, v)
			composerOutput, composerErr := ComposerCommand(context, customIni, "require", v, "--dev")
			if composerErr != nil || strings.Contains(string(composerOutput), "An error occurred") {
				fmt.Printf("	--> An error ocurred while running composer install: '%s'", composerOutput)
				return packit.BuildResult{}, blueprintGenErr
			}
			logger.Break()
		}

		for kk, vv := range globalDependencies {
			logger.Process("Installing global dependency [%d]: %s", kk, vv)
			gOutput, gErr := ComposerCommand(context, customIni, "global", "require", vv, "--dev")
			logger.Detail("%s", gOutput)
			if gErr != nil || strings.Contains(string(gOutput), "An error occurred") {
				fmt.Printf("	--> An error ocurred while running composer (global) install: '%s'", gOutput)
				return packit.BuildResult{}, blueprintGenErr
			}
			logger.Break()
		}

		logger.Process("Publishing vendor config")

		vendorConfigErr := ArtisanCommand(customIni, "vendor:publish", "--tag", "blueprint-config")

		if vendorConfigErr != nil {
			fmt.Printf("	--> An error ocurred while plublishing vendor config: '%s'", "")
			return packit.BuildResult{}, blueprintGenErr
		}

		logger.Break()

		// Codegen

		logger.Process("Larased version")
		err := LarasedCommand(customIni, "--version")
		if err != nil {
			return packit.BuildResult{}, err
		}

		logger.Process("Set blueprint config models_namespace = Models")

		//# https://github.com/laravel-shift/blueprint/issues/366
		lerr := LarasedCommand(customIni, "larased:config-edit", "-d", context.WorkingDir, "config.blueprint/models_namespace", "Models", "-vvv")

		if lerr != nil {
			fmt.Printf("	--> An error ocurred while running larased command: '%s'", "")
			return packit.BuildResult{}, blueprintGenErr
		}

		logger.Break()

		logger.Process("Set blueprint config generate_fqcn_route = true")

		//# https://github.com/laravel-shift/blueprint/issues/384
		//# https://github.com/laravel-shift/blueprint/issues/377
		lerr = LarasedCommand(customIni, "larased:config-edit", "-d", context.WorkingDir, "config.blueprint/generate_fqcn_route", "true", "-vvv")

		if lerr != nil {
			fmt.Printf("	--> An error ocurred while running larased command: '%s'", "")
			return packit.BuildResult{}, blueprintGenErr
		}

		logger.Process("Initializing draft.yml file")

		// initialize

		blueprint := pexec.NewExecutable("blueprint-gen")

		blueprintGenErr = blueprint.Execute(pexec.Execution{
			Args:   []string{"init"},
			Stdout: os.Stdout,
		})
		if blueprintGenErr != nil {
			fmt.Printf("	--> An error ocurred creating a blueprint file: '%s'", blueprintGenErr)
			panic(blueprintGenErr)
		}
		logger.Break()

		logger.Process("Creating a data model")
		logger.Detail("Found %d models", len(larabootStruct.Framework.Models))

		for _, o := range larabootStruct.Framework.Models {
			logger.Detail("	Working on model: %s", o.Name)

			// model $Name --name $NewPost --with title=string:400,content=longtext,author_id=id:user,published="nullable timestamp"
			command_args := []string{"model"}

			command_args = append(command_args,
				o.Name,
				"--name",
				o.Name,
				"--with",
				NewColumnsFormatter(o.Columns))

			blueprintGenErr = blueprint.Execute(pexec.Execution{
				Args:   command_args,
				Stdout: os.Stdout,
			})

			if blueprintGenErr != nil {
				panic(blueprintGenErr)
			}

		}

		// After draft.yml has been created we're ready to run codegen by calling blueprint:build command
		artisanErr := ArtisanCommand(customIni, "blueprint:build", "draft.yml")
		fmt.Printf("%s", "")

		if artisanErr != nil {
			fmt.Printf("	--> An error ocurred while running blueprint codegen: '%s'", "")
			return packit.BuildResult{}, blueprintGenErr
		}

		// end Codegen

		thisLayer.Metadata = map[string]interface{}{
			"built_at":  clock.Now().Format(time.RFC3339Nano),
			"cache_sha": dependency.SHA256,
		}

		// Clean before leave if was instructed.
		if CleanUp {
			logger.Process("Cleaning up")
			for k, v := range composerDevDependencies {
				logger.Process("Uninstalling dev dependency [%d]: %s", k, v)
				uninstallOutput, uninstallErr := ComposerCommand(context, customIni, "remove", "--dev", v)
				if uninstallErr != nil || strings.Contains(string(uninstallOutput), "An error occurred") {
					fmt.Printf("	--> An error ocurred while running composer install: '%s'", uninstallOutput)
					return packit.BuildResult{}, blueprintGenErr
				}
			}
		}

		if m.LaravelModel.Git.Enabled && m.LaravelModel.Git.Commit {
			err2 := GitSign()
			if err2 != nil {
				return packit.BuildResult{}, err2
			}
		}

		return packit.BuildResult{
			Layers: []packit.Layer{thisLayer},
		}, nil
	}
}

func NewColumnsFormatter(columns []struct {
	Name string `json:"name"`
	Type string `json:"type"`
}) string {
	var res []string
	for _, col := range columns {
		res = append(res, fmt.Sprintf("%s=\"%s\"", col.Name, col.Type))
	}
	return strings.Join(res, ",")
}
