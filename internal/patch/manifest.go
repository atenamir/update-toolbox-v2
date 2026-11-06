package patch

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Manifest struct {
	Version        string         `json:"version"`
	PackageVersion PackageVersion `json:"package_version"`
}
type PackageVersion struct {
	Baadbaan  string `json:"baadbaan"`
	Technical string `json:"technical"`
	Docker    string `json:"docker"`
	Toolbox   string `json:"toolbox"`
	Discovery string `json:"discovery"`
}

// ===================== LoadManifest reads and parses a JSON manifest file ===============
func LoadManifest(path string) ([]Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read manifest %q: %w", path, err)
	}

	var manifests []Manifest

	if err := json.Unmarshal(data, &manifests); err != nil {
		return nil, fmt.Errorf("parse manifest %q: %w", path, err)
	}

	return manifests, nil
}

// ================= ValidateManifest checks the release manifest ==================

func ValidateManifest(manifests []Manifest) error {
	if len(manifests) != 2 {
		return fmt.Errorf("expected exactly 2 releases, got %d", len(manifests))
	}

	var errorsFound []string

	for i, manifest := range manifests {
		if strings.TrimSpace(manifest.Version) == "" {
			errorsFound = append(
				errorsFound,
				fmt.Sprintf("release %d: version is required", i+1),
			)
		}

		components := []struct {
			name    string
			version string
		}{
			{"baadbaan", manifest.PackageVersion.Baadbaan},
			{"technical", manifest.PackageVersion.Technical},
			{"docker", manifest.PackageVersion.Docker},
			{"toolbox", manifest.PackageVersion.Toolbox},
			{"discovery", manifest.PackageVersion.Discovery},
		}

		for _, component := range components {
			if strings.TrimSpace(component.version) == "" {
				errorsFound = append(
					errorsFound,
					fmt.Sprintf(
						"release %d: package_version.%s is required",
						i+1,
						component.name,
					),
				)
			}
		}
	}

	if len(errorsFound) > 0 {
		return fmt.Errorf("invalid manifest:\n%s",
			strings.Join(errorsFound, "\n"))
	}

	if manifests[0].Version == manifests[1].Version {
		return fmt.Errorf("source and target versions must be different")
	}

	return nil
}
