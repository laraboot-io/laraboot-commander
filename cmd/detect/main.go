package main

import (
	"github.com/paketo-buildpacks/packit"
	Larabootcommander "laraboot-buildpacks/laraboot-commander/laraboot"
)

func main() {
	packit.Detect(Larabootcommander.Detect())
}
