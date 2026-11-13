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

	// ================== Create a temporary file in the output directory ===========================
	outputDir := filepath.Dir(outputPath)

	tempFile, err := os.CreateTemp(outputDir, ".patch-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary file: %w", err)
	}

	tempPath := tempFile.Name()

	if err := tempFile.Close(); err != nil {
		os.Remove(tempPath)
		return err
	}

	//====================== Remove temporary file on failure ===================
	defer os.Remove(tempPath)

	// ================= Build the archive into the temporary path =====================
	if err := CreatePatchArchive(manifest, repoPath, tempPath); err != nil {
		return err
	}

	// ================ Check whether the final output already exists ===================
	if _, err := os.Stat(outputPath); err == nil {
		return fmt.Errorf("output already exists: %s", outputPath)
	} else if !os.IsNotExist(err) {
		return err
	}

	// ============== Rename the completed archive=====================
	if err := os.Rename(tempPath, outputPath); err != nil {
		return fmt.Errorf("finalize patch archive: %w", err)
	}

	return nil
}
