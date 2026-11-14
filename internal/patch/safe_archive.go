package patch

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreatePatchArchiveSafe(
	manifest PatchManifest,
	repoPath string,
	outputPath string,
) error {

	//================== Create a temporary file in the output directory =================
	outputDir := filepath.Dir(outputPath)

	tempFile, err := os.CreateTemp(outputDir, ".patch-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary file: %w", err)
	}

	tempPath := tempFile.Name()

	//========== Always remove the temporary file when finished =================
	defer os.Remove(tempPath)

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("close temporary file: %w", err)
	}

	//================== Build the archive in the temporary file ==================
	if err := CreatePatchArchive(
		manifest,
		repoPath,
		tempPath,
	); err != nil {
		return fmt.Errorf("build patch archive: %w", err)
	}

	//=============== Verify the completed archive before publishing ========================
	if err := VerifyPatchArchive(tempPath); err != nil {
		return fmt.Errorf("verify generated patch: %w", err)
	}

	//================ Flush completed file data before publishing =====================
	completedFile, err := os.OpenFile(tempPath, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("open completed patch: %w", err)
	}

	if err := completedFile.Sync(); err != nil {
		completedFile.Close()
		return fmt.Errorf("sync patch archive: %w", err)
	}

	if err := completedFile.Close(); err != nil {
		return fmt.Errorf("close completed patch: %w", err)
	}

	//============= Publish without replacing an existing output file =========
	if err := os.Link(tempPath, outputPath); err != nil {
		return fmt.Errorf("publish patch archive: %w", err)
	}

	return nil
}
