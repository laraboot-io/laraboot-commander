package main

import (
	"github.com/cloudfoundry/packit"
	Larabootcommander "laraboot-buildpacks/laraboot-commander/laraboot"
)

func main() {
	packit.Detect(Larabootcommander.Detect())
}
