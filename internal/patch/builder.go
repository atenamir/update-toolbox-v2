package patch

import (
	"encoding/json"
	"fmt"
	"os"
)

// ============= PatchManifest describes the changes included in a patch============
type PatchManifest struct {
	FromVersion string      `json:"from_version"`
	ToVersion   string      `json:"to_version"`
	Files       []PatchFile `json:"files"`
}

// =============== PatchFile describes one changed file================
type PatchFile struct {
	Path   string `json:"path"`
	Status string `json:"status"`
}

func BuildPatchManifest(
	fromVersion string,
	toVersion string,
	changedFiles []ChangedFile,
) PatchManifest {

	manifest := PatchManifest{
		FromVersion: fromVersion,
		ToVersion:   toVersion,
		Files:       make([]PatchFile, 0, len(changedFiles)),
	}

	for _, file := range changedFiles {
		manifest.Files = append(manifest.Files, PatchFile{
			Path:   file.Path,
			Status: file.Status,
		})
	}

	return manifest
}

func WritePatchManifest(
	manifest PatchManifest,
	outputPath string,
) error {

	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("encode patch manifest: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("write patch manifest: %w", err)
	}

	return nil
}
