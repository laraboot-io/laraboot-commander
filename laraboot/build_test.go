package larabootcommander

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	gobuild "github.com/paketo-buildpacks/go-build"
	"github.com/paketo-buildpacks/go-build/fakes"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/chronos"
	"github.com/paketo-buildpacks/packit/scribe"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

//nolint:funlen //meh
func testBuild(t *testing.T, _ spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect //nolint:govet //meh
		layersDir  string
		workingDir string
		cnbDir     string
		timestamp  time.Time
		logs       *bytes.Buffer

		buildProcess  *fakes.BuildProcess
		pathManager   *fakes.PathManager
		sourceRemover *fakes.SourceRemover
		parser        *fakes.ConfigurationParser

		build packit.BuildFunc
	)

	it.Before(func() {
		var err error
		layersDir, err = os.MkdirTemp("", "layers")
		Expect(err).NotTo(HaveOccurred())

		cnbDir, err = os.MkdirTemp("", "cnb")
		Expect(err).NotTo(HaveOccurred())

		workingDir, err = os.MkdirTemp("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		buildProcess = &fakes.BuildProcess{}
		buildProcess.ExecuteCall.Returns.Binaries = []string{"path/some-start-command", "path/another-start-command"}

		pathManager = &fakes.PathManager{}
		pathManager.SetupCall.Returns.GoPath = "some-go-path"
		pathManager.SetupCall.Returns.Path = "some-app-path"

		timestamp = time.Now()
		clock := chronos.NewClock(func() time.Time {
			return timestamp
		})

		logs = bytes.NewBuffer(nil)

		sourceRemover = &fakes.SourceRemover{}

		parser = &fakes.ConfigurationParser{}
		parser.ParseCall.Returns.BuildConfiguration = gobuild.BuildConfiguration{
			Targets:    []string{"some-target", "other-target"},
			Flags:      []string{"some-flag", "other-flag"},
			ImportPath: "some-import-path",
		}

		build = gobuild.Build(
			parser,
			buildProcess,
			pathManager,
			clock,
			scribe.NewEmitter(logs),
			sourceRemover,
		)
	})

	it.After(func() {
		Expect(os.RemoveAll(layersDir)).To(Succeed())
		Expect(os.RemoveAll(cnbDir)).To(Succeed())
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	it("returns a result that builds correctly", func() {
		_, err := build(packit.BuildContext{
			WorkingDir: workingDir,
			CNBPath:    cnbDir,
			Stack:      "some-stack",
			BuildpackInfo: packit.BuildpackInfo{
				Name:    "Some Buildpack",
				Version: "some-version",
			},
			Layers: packit.Layers{Path: layersDir},
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(parser.ParseCall.Receives.BuildpackVersion).To(Equal("some-version"))
		Expect(parser.ParseCall.Receives.WorkingDir).To(Equal(workingDir))

		Expect(pathManager.SetupCall.Receives.Workspace).To(Equal(workingDir))
		Expect(pathManager.SetupCall.Receives.ImportPath).To(Equal("some-import-path"))

		Expect(buildProcess.ExecuteCall.Receives.Config).To(Equal(gobuild.GoBuildConfiguration{
			Workspace: "some-app-path",
			Output:    filepath.Join(layersDir, "targets", "bin"),
			GoPath:    "some-go-path",
			GoCache:   filepath.Join(layersDir, "gocache"),
			Flags:     []string{"some-flag", "other-flag"},
			Targets:   []string{"some-target", "other-target"},
		}))

		Expect(pathManager.TeardownCall.Receives.GoPath).To(Equal("some-go-path"))

		Expect(sourceRemover.ClearCall.Receives.Path).To(Equal(workingDir))

		Expect(logs.String()).To(ContainSubstring("Some Buildpack some-version"))
		Expect(logs.String()).To(ContainSubstring("Assigning launch processes"))
		Expect(logs.String()).To(ContainSubstring("web: path/some-start-command"))
		Expect(logs.String()).To(ContainSubstring("some-start-command: path/some-start-command"))
		Expect(logs.String()).To(ContainSubstring("another-start-command: path/another-start-command"))
	})
}

func TestUnitGoBuild(t *testing.T) {
	suite := spec.New("go-build", spec.Report(report.Terminal{}))
	suite("Build", testBuild, spec.Sequential())
	suite.Run(t)
}
