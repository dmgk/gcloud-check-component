package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

const snapshotURL = "https://dl.google.com/dl/cloudsdk/channels/rapid/components-2.json"

type Snapshot struct {
	Components    []*Component   `json:"components"`
	Revision      int64          `json:"revision"`
	SchemaVersion *SchemaVersion `json:"schema_version"`
	Version       string         `json:"version"`
}

type Component struct {
	Dependencies    []string  `json:"dependencies"`
	Details         Details   `json:"details"`
	ID              string    `json:"id"`
	IsConfiguration bool      `json:"is_configuration"`
	IsHidden        bool      `json:"is_hidden"`
	IsRequired      bool      `json:"is_required"`
	Platform        *Platform `json:"platform"`
	Version         *Version  `json:"version"`
}

type Details struct {
	Description string `json:"description"`
	DisplayName string `json:"display_name"`
}

type Platform struct {
	Architectures    []string `json:"architectures,omitempty"`
	OperatingSystems []string `json:"operating_systems,omitempty"`
}

type SchemaVersion struct {
	NoUpdate bool   `json:"no_update"`
	URL      string `json:"url"`
	Version  int64  `json:"version"`
}

func (s *Snapshot) getComponent(id string) *Component {
	for _, c := range s.Components {
		if c.ID == id {
			return c
		}
	}
	return nil
}

type Version struct {
	BuildNumber   int64  `json:"build_number"`
	VersionString string `json:"version_string"`
}

func getSnapshot() (*Snapshot, error) {
	resp, err := http.Get(snapshotURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s Snapshot

	d := json.NewDecoder(resp.Body)
	if err := d.Decode(&s); err != nil {
		return nil, err
	}

	return &s, nil
}

func main() {
	var id string
	var buildNumber int64

	flag.Parse()
	args := flag.Args()

	if flagManifest {
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "Missing component_id parameter")
			flag.Usage()
			os.Exit(255)
		}
		id = args[0]
	} else {
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Missing component_id and build_number parameters")
			flag.Usage()
			os.Exit(255)
		}
		var err error
		id = args[0]
		buildNumber, err = strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(255)
		}
	}

	s, err := getSnapshot()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(255)
	}

	c := s.getComponent(id)
	if c == nil {
		fmt.Fprintf(os.Stderr, "No such component: %s\n", id)
		os.Exit(255)
	}

	// Print manifest and exit
	if flagManifest {
		m := *s
		m.Components = []*Component{c}
		// Reset platform restrictions so the component is visible in gcloud on FreeBSD
		m.Components[0].Platform = &Platform{}
		b, err := json.MarshalIndent(&m, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(255)
		}
		fmt.Println(string(b))
		os.Exit(0)
	}

	// Exit with zero status if new version is found
	if c.Version.BuildNumber > buildNumber {
		if !flagQuiet {
			fmt.Printf("Found new build for %s: %d (have %d)\n", id, c.Version.BuildNumber, buildNumber)
		}
		os.Exit(0)
	}

	// Exit with status 1 if there's no new version
	if !flagQuiet {
		fmt.Printf("No new builds for %s\n", id)
	}
	os.Exit(1)
}

var flagManifest bool
var flagQuiet bool

func init() {
	flag.BoolVar(&flagManifest, "manifest", false, "do not check for new version, write component manifest to standard output instead")
	flag.BoolVar(&flagQuiet, "quiet", false, "be quiet, exit with zero status if new version is found")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] component_id [build_number]\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
}
