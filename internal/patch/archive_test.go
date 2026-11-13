package patch

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreatePatchArchive(t *testing.T) {
	manifest := PatchManifest{
		FromVersion: "v1.0.0",
		ToVersion:   "v1.1.0",
	}

	outputPath := filepath.Join(t.TempDir(), "patch.tar.gz")

	if err := CreatePatchArchive(manifest, "", outputPath); err != nil {
		t.Fatalf("create archive failed: %v", err)
	}

	file, err := os.Open(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		t.Fatal(err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)

	header, err := tr.Next()
	if err != nil {
		t.Fatal(err)
	}

	if header.Name != "manifest.json" {
		t.Fatalf("unexpected archive entry: %s", header.Name)
	}
}

func TestCreatePatchArchiveWithFiles(t *testing.T) {
	repoPath := filepath.Join("..", "..", "..", "patch-lab")

	changedFiles, err := GetChangedFiles(
		repoPath,
		"v1.0.0",
		"v1.1.0",
	)
	if err != nil {
		t.Fatal(err)
	}

	manifest := BuildPatchManifest(
		"v1.0.0",
		"v1.1.0",
		changedFiles,
	)

	outputPath := filepath.Join(t.TempDir(), "patch.tar.gz")

	err = CreatePatchArchive(manifest, repoPath, outputPath)
	if err != nil {
		t.Fatalf("create archive failed: %v", err)
	}

	file, err := os.Open(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		t.Fatal(err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)

	expected := map[string]bool{
		"manifest.json":            false,
		"files/app/NewService.php": false,
		"files/app/User.php":       false,
	}

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}

		if _, exists := expected[header.Name]; !exists {
			t.Fatalf("unexpected archive entry: %s", header.Name)
		}

		expected[header.Name] = true

		data, err := io.ReadAll(tr)
		if err != nil {
			t.Fatal(err)
		}

		if header.Name == "files/app/User.php" &&
			!strings.Contains(string(data), "Version 2") {
			t.Fatal("User.php has incorrect content")
		}
	}

	for name, found := range expected {
		if !found {
			t.Errorf("missing archive entry: %s", name)
		}
	}
}
