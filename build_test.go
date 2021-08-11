package laraboot_test

import (
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/go-build/fakes"
	"github.com/sclevine/spec"
	"io/ioutil"
	"os"
	"testing"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		layersDir    string
		workingDir   string
		cnbDir       string
		buildProcess *fakes.BuildProcess
		pathManager  *fakes.PathManager
	)

	it.Before(func() {
		var err error
		layersDir, err = ioutil.TempDir("", "layers")
		Expect(err).NotTo(HaveOccurred())

		cnbDir, err = ioutil.TempDir("", "cnb")
		Expect(err).NotTo(HaveOccurred())

		workingDir, err = ioutil.TempDir("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		buildProcess = &fakes.BuildProcess{}
		buildProcess.ExecuteCall.Returns.Binaries = []string{"path/some-start-command", "path/another-start-command"}

		pathManager = &fakes.PathManager{}
		pathManager.SetupCall.Returns.GoPath = "some-go-path"
		pathManager.SetupCall.Returns.Path = "some-app-path"
	})

	it.After(func() {
		Expect(os.RemoveAll(layersDir)).To(Succeed())
		Expect(os.RemoveAll(cnbDir)).To(Succeed())
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})
}
