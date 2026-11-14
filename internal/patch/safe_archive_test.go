package patch

import (
	"os"
	"path/filepath"
	"testing"
)

// =========== An existing archive must never be overwritten ==============
func TestCreatePatchArchiveSafe_NoOverwrite(t *testing.T) {
	outputDir := t.TempDir()
	outputPath := filepath.Join(outputDir, "patch.tar.gz")

	//================ Create an existing file with known content ===============
	original := []byte("DO NOT OVERWRITE")

	if err := os.WriteFile(outputPath, original, 0644); err != nil {
		t.Fatal(err)
	}

	manifest := PatchManifest{
		FromVersion: "v1.0.0",
		ToVersion:   "v1.1.0",
		Files:       []PatchFile{},
	}

	//================== This should fail because the destination already exists ==================
	err := CreatePatchArchiveSafe(manifest, "", outputPath)

	if err == nil {
		t.Fatal("expected an error because output already exists")
	}

	//=============== Ensure the original file has not changed ===========
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != string(original) {
		t.Fatal("existing output file was modified")
	}
}

// ========== A valid patch should be created successfully =================
func TestCreatePatchArchiveSafe_Success(t *testing.T) {
	outputDir := t.TempDir()
	outputPath := filepath.Join(outputDir, "patch.tar.gz")

	manifest := PatchManifest{
		FromVersion: "v1.0.0",
		ToVersion:   "v1.1.0",
		Files:       []PatchFile{},
	}

	err := CreatePatchArchiveSafe(manifest, "", outputPath)
	if err != nil {
		t.Fatalf("expected success, got: %v", err)
	}

	//================ Check that the final archive exists ============
	if _, err := os.Stat(outputPath); err != nil {
		t.Fatalf("output file does not exist: %v", err)
	}

	//========== Verify the archive content ===========
	if err := VerifyPatchArchive(outputPath); err != nil {
		t.Fatalf("invalid archive: %v", err)
	}

	//============== Ensure temporary files were removed =============
	tempFiles, err := filepath.Glob(
		filepath.Join(outputDir, ".patch-*.tmp"),
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(tempFiles) != 0 {
		t.Fatalf("temporary files remain: %v", tempFiles)
	}
}

// ============ A failed build must not publish an archive ==================
func TestCreatePatchArchiveSafe_BuildFailure(t *testing.T) {
	outputDir := t.TempDir()
	outputPath := filepath.Join(outputDir, "patch.tar.gz")

	manifest := PatchManifest{
		FromVersion: "v1.0.0",
		ToVersion:   "v1.1.0",
		Files: []PatchFile{
			{
				Path:   "../secret.txt",
				Status: "A",
			},
		},
	}

	//============== The invalid file path must cause an error ===================
	err := CreatePatchArchiveSafe(manifest, "", outputPath)

	if err == nil {
		t.Fatal("expected build failure")
	}

	//============= The final archive must not exist ==================
	if _, err := os.Stat(outputPath); !os.IsNotExist(err) {
		t.Fatalf("expected no output file, got: %v", err)
	}

	//============= Ensure temporary files were removed ===============
	tempFiles, err := filepath.Glob(
		filepath.Join(outputDir, ".patch-*.tmp"),
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(tempFiles) != 0 {
		t.Fatalf("temporary files remain: %v", tempFiles)
	}
}
