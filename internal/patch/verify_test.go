package patch

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVerifyPatchArchive_RejectsSHA256Mismatch(t *testing.T) {
	manifest := PatchManifest{
		FromVersion: "v1.0.0",
		ToVersion:   "v1.1.0",
		Files: []PatchFile{
			{
				Path:   "app/User.php",
				Status: "M",
				SHA256: CalculateSHA256([]byte("Version 2")),
			},
		},
	}

	manifestData, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}

	var archive bytes.Buffer
	gz := gzip.NewWriter(&archive)
	tw := tar.NewWriter(gz)

	if err := writeTarEntry(tw, "manifest.json", manifestData); err != nil {
		t.Fatal(err)
	}
	if err := writeTarEntry(tw, "files/app/User.php", []byte("Version 3")); err != nil {
		t.Fatal(err)
	}

	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}

	archivePath := filepath.Join(t.TempDir(), "patch.tar.gz")
	if err := os.WriteFile(archivePath, archive.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}

	err = VerifyPatchArchive(archivePath)
	if err == nil {
		t.Fatal("expected verification to reject the changed file")
	}

	want := "SHA256 mismatch for app/User.php"
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("expected error containing %q, got %v", want, err)
	}
}
