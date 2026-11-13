package patch

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWritePatchManifest(t *testing.T) {
	changedFiles := []ChangedFile{
		{Path: "app/NewService.php", Status: "A"},
		{Path: "app/OldService.php", Status: "D"},
		{Path: "app/User.php", Status: "M"},
	}

	manifest := BuildPatchManifest(
		"v1.0.0",
		"v1.1.0",
		changedFiles,
	)

	outputPath := filepath.Join(t.TempDir(), "manifest.json")

	if err := WritePatchManifest(manifest, outputPath); err != nil {
		t.Fatalf("failed to write manifest: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read manifest: %v", err)
	}

	var result PatchManifest
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if result.FromVersion != "v1.0.0" {
		t.Errorf("unexpected source version: %s", result.FromVersion)
	}

	if result.ToVersion != "v1.1.0" {
		t.Errorf("unexpected target version: %s", result.ToVersion)
	}

	if len(result.Files) != 3 {
		t.Fatalf("expected 3 files, got %d", len(result.Files))
	}

	expected := map[string]string{
		"app/NewService.php": "A",
		"app/OldService.php": "D",
		"app/User.php":       "M",
	}

	for _, file := range result.Files {
		status, exists := expected[file.Path]
		if !exists {
			t.Errorf("unexpected file: %s", file.Path)
			continue
		}
		if file.Status != status {
			t.Errorf("wrong status for %s: got %s, want %s",
				file.Path, file.Status, status)
		}
	}
}
