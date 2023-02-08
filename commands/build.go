package commands

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var buildOutput string

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Add log info before each return statement in target package",
	Long: `
Build command will modify the source files in place, so you can continue to work in the origin place.
`,
	Example: `
# Generate stubbed source file using return command.
return_trace_log build .
`, Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Fail to build: %v", err)
		}

		err = runBuild(args, wd)
		if err != nil {
			panic(err)
		}
	},
}

func runBuild(args []string, wd string) error {
	pkgs, err := listPackage(wd, "-json "+strings.Join(args, " "))
	if err != nil {
		return err
	}

	// FIXME: iterate each pkg for generated
	_ = pkgs

	// FIXME: implement it
	return nil
}

func listPackage(dir string, args string) (pkgs map[string]*Package, err error) {
	goListCmd := exec.Command("/bin/bash", "-c", "go list "+args)
	goListCmd.Dir = dir
	goListCmdOutput, err := goListCmd.Output()
	if err != nil {
		return nil, GoListDependenciesError
	}

	pkgs = make(map[string]*Package, 0)
	decoder := json.NewDecoder(bytes.NewReader(goListCmdOutput))
	for {
		var p *Package
		if err = decoder.Decode(&p); err != nil {
			if err == io.EOF {
				break
			}

			return nil, GoListDependenciesError
		}

		pkgs[p.ImportPath] = p
	}

	return pkgs, nil
}

// Package map a package output by go list
// this is subset of package struct in: https://github.com/golang/go/blob/master/src/cmd/go/internal/load/pkg.go#L58
type Package struct {
	Dir        string `json:"Dir"`        // directory containing package sources
	ImportPath string `json:"ImportPath"` // import path of package in dir
	Name       string `json:"Name"`       // package name
	Target     string `json:",omitempty"` // installed target for this package (may be executable)
	Root       string `json:",omitempty"` // Go root, Go path dir, or module root dir containing this package

	Module   *ModulePublic `json:",omitempty"`         // info about package's module, if any
	Goroot   bool          `json:"Goroot,omitempty"`   // is this package in the Go root?
	Standard bool          `json:"Standard,omitempty"` // is this package part of the standard Go library?
	DepOnly  bool          `json:"DepOnly,omitempty"`  // package is only a dependency, not explicitly listed

	// Source files
	GoFiles  []string `json:"GoFiles,omitempty"`  // .go source files (excluding CgoFiles, TestGoFiles, XTestGoFiles)
	CgoFiles []string `json:"CgoFiles,omitempty"` // .go source files that import "C"

	// Dependency information
	Deps      []string          `json:"Deps,omitempty"` // all (recursively) imported dependencies
	Imports   []string          `json:",omitempty"`     // import paths used by this package
	ImportMap map[string]string `json:",omitempty"`     // map from source import to ImportPath (identity entries omitted)

	// Error information
	Incomplete bool            `json:"Incomplete,omitempty"` // this package or a dependency has an error
	Error      *PackageError   `json:"Error,omitempty"`      // error loading package
	DepsErrors []*PackageError `json:"DepsErrors,omitempty"` // errors loading dependencies
}

// ModulePublic represents the package info of a module
type ModulePublic struct {
	Path      string        `json:",omitempty"` // module path
	Version   string        `json:",omitempty"` // module version
	Versions  []string      `json:",omitempty"` // available module versions
	Replace   *ModulePublic `json:",omitempty"` // replaced by this module
	Time      *time.Time    `json:",omitempty"` // time version was created
	Update    *ModulePublic `json:",omitempty"` // available update (with -u)
	Main      bool          `json:",omitempty"` // is this the main module?
	Indirect  bool          `json:",omitempty"` // module is only indirectly needed by main module
	Dir       string        `json:",omitempty"` // directory holding local copy of files, if any
	GoMod     string        `json:",omitempty"` // path to go.mod file describing module, if any
	GoVersion string        `json:",omitempty"` // go version used in module
	Error     *ModuleError  `json:",omitempty"` // error loading module
}

// ModuleError represents the error loading module
type ModuleError struct {
	Err string // error text
}

// PackageError is the error info for a package when list failed
type PackageError struct {
	ImportStack []string // shortest path from package named on command line to this one
	Pos         string   // position of error (if present, file:line:col)
	Err         string   // the error itself
}

func init() {
	addBuildFlags(buildCmd.Flags())
	buildCmd.Flags().StringVarP(&buildOutput, "output", "o", "", "it forces build to write the resulting executable to the named output file")
	rootCmd.AddCommand(buildCmd)
}
