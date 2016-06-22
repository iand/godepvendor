package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Godeps struct {
	ImportPath   string
	GoVersion    string   // Abridged output of 'go version'.
	GodepVersion string   // Abridged output of 'godep version'
	Packages     []string // Arguments to godep save, if any.
	Deps         []struct {
		ImportPath string
		Comment    string // Description of commit, if present.
		Rev        string // VCS-specific commit ID.
	}
}

type Vendor struct {
	RootPath string          `json:"rootPath,omitempty"`
	Comment  string          `json:"comment,omitempty"`
	Package  []VendorPackage `json:"package,omitempty"`
}

type VendorPackage struct {
	Path         string `json:"path"`
	Origin       string `json:"origin,omitempty"`
	Revision     string `json:"revision"`
	RevisionTime string `json:"revisionTime,omitempty"`
	Comment      string `json:"comment,omitempty"`
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "missing required argument (expected path to Godeps.json)\n")
		os.Exit(1)
	}

	fin, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer fin.Close()

	dec := json.NewDecoder(fin)
	var godeps Godeps
	err = dec.Decode(&godeps)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	vendor := Vendor{
		RootPath: godeps.ImportPath,
		Comment:  fmt.Sprintf("Converted from %s", flag.Arg(0)),
	}

	for _, dep := range godeps.Deps {
		vendor.Package = append(vendor.Package, VendorPackage{
			Path:     dep.ImportPath,
			Revision: dep.Rev,
			Comment:  dep.Comment,
		})
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(vendor)

}
