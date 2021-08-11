package laraboot

import (
	"fmt"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/pexec"
	"os"
	"os/exec"
	"strings"
)

func GitSign() error {

	commands := []string{
		"config --global user.email \"larabot@laraboot.io\"",
		"config --global user.name \"LarabootProject\"",
		"init",
		"add .",
		"status",
		"commit -m add-models",
	}

	for _, command := range commands {
		split := strings.Split(command, " ")
		err := gitCmd(split...)
		go func(error) {
			if err != nil {
				return
			}
		}(err)
	}

	return nil
}

func gitCmd(args ...string) error {

	git := pexec.NewExecutable("php")

	fmt.Printf("Running git command : `git %s` \n", strings.Join(args, " "))

	phpErr := git.Execute(pexec.Execution{
		Args:   args,
		Stdout: os.Stdout,
	})

	if phpErr != nil {
		return phpErr
	}

	return nil
}

func ComposerCommand(context packit.BuildContext, customIni string, arg ...string) ([]byte, error) {

	php_path, err := exec.LookPath("php")

	composerDir := "/layers/paketo-buildpacks_php-composer/composer"
	composerPhar := fmt.Sprintf("%s/composer.phar", composerDir)

	if err != nil {
		panic(err)
	}

	args := append([]string{
		fmt.Sprintf("-dextension_dir=%s", os.Getenv("PHP_EXTENSION_DIR")),
		fmt.Sprintf("-derror_reporting=%s", "E_ALL"),
		"-c",
		customIni,
		composerPhar,
	}, arg...)

	cmd := exec.Command(php_path, args...)
	cmd.Dir = context.WorkingDir

	return cmd.CombinedOutput()
}

func LarasedCommand(customIni string, arg ...string) error {

	php := pexec.NewExecutable("php")

	args := append([]string{
		fmt.Sprintf("-dextension_dir=%s", os.Getenv("PHP_EXTENSION_DIR")),
		fmt.Sprintf("-derror_reporting=%s", "E_ALL"),
		"-c",
		customIni,
		"/home/cnb/.config/composer/vendor/bin/larased",
	}, arg...)

	phpErr := php.Execute(pexec.Execution{
		Args:   args,
		Stdout: os.Stdout,
	})

	if phpErr != nil {
		return phpErr
	}

	return nil

}

func ArtisanCommand(customIni string, arg ...string) error {

	php := pexec.NewExecutable("php")

	args := append([]string{
		fmt.Sprintf("-dextension_dir=%s", os.Getenv("PHP_EXTENSION_DIR")),
		fmt.Sprintf("-derror_reporting=%s", "E_ALL"),
		"-c",
		customIni,
		"artisan",
	}, arg...)

	phpErr := php.Execute(pexec.Execution{
		Args:   args,
		Stdout: os.Stdout,
	})

	if phpErr != nil {
		return phpErr
	}

	return nil

}

func LsCommand(context packit.BuildContext, arg ...string) ([]byte, error) {

	ls, err := exec.LookPath("ls")

	if err != nil {
		panic(err)
	}

	args := append([]string{}, arg...)

	cmd := exec.Command(ls, args...)
	cmd.Dir = context.WorkingDir

	return cmd.CombinedOutput()
}
