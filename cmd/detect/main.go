package main

import (
	"github.com/cloudfoundry/packit"
	"laraboot-buildpacks/laraboot-commander/laraboot"
)

func main() {
	packit.Detect(laraboot.Detect())
}
