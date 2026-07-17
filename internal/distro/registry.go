package distro

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Distro struct {
	VerJSON          string `json:"verjson"`
	Name             string `json:"name"`
	URL              string `json:"url"`
	InstallationType string `json:"installationtype"`
	Info             string `json:"info,omitempty"`
	Checksum         string `json:"checksum,omitempty"`
	ChecksumType     string `json:"checksumtype,omitempty"`
	Sig              string `json:"sig,omitempty"`
	SigType          string `json:"sigtype,omitempty"`
}

// distrosDir returns the path to the distros/ folder next to the executable.
func distrosDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "distros"
	}
	return filepath.Join(filepath.Dir(exe), "distros")
}

// LoadAll reads every *.json file in the distros directory.
func LoadAll() ([]Distro, error) {
	dir := distrosDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var list []Distro
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		var d Distro
		if err := json.Unmarshal(data, &d); err != nil {
			continue
		}
		list = append(list, d)
	}
	return list, nil
}

// Find returns the distro with the given name (case-insensitive).
func Find(name string) (Distro, bool) {
	list, err := LoadAll()
	if err != nil {
		return Distro{}, false
	}
	for _, d := range list {
		if strings.EqualFold(d.Name, name) {
			return d, true
		}
	}
	return Distro{}, false
}
