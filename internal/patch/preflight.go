package patch

import (
	"fmt"
	"os"
	"strings"
)

func CheckInstalledVersion(manifest PatchManifest, versionFile string) error {
	data, err := os.ReadFile(versionFile)
	if err != nil {
		return fmt.Errorf("read installed version: %w", err)
	}

	installed := strings.TrimSpace(string(data))
	if installed == "" {
		return fmt.Errorf("installed version file %q is empty", versionFile)
	}

	if installed != manifest.FromVersion {
		return fmt.Errorf(
			"patch requires version %q, installed version is %q",
			manifest.FromVersion,
			installed,
		)
	}

	return nil
}
