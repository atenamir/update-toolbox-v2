package patch

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

//=========== CreateBackup copies files that the verified patch will replace or delete ==================
func CreateBackup(manifest PatchManifest, installDir string) (string, error) {
	installed, err := os.OpenRoot(installDir)
	if err != nil {
		return "", fmt.Errorf("open installation: %w", err)
	}
	defer installed.Close()

	backupDir, err := os.MkdirTemp("", "toolbox-backup-*")
	if err != nil {
		return "", fmt.Errorf("create backup directory: %w", err)
	}

	complete := false
	defer func() {
		if !complete {
			os.RemoveAll(backupDir)
		}
	}()

	for _, file := range manifest.Files {
		switch file.Status {
		case "M", "D":
			if err := copyBackupFile(installed, backupDir, file.Path); err != nil {
				return "", fmt.Errorf("back up %q: %w", file.Path, err)
			}
		}
	}

	if err := WritePatchManifest(manifest, filepath.Join(backupDir, "manifest.json")); err != nil {
		return "", err
	}

	complete = true
	return backupDir, nil
}

func copyBackupFile(installed *os.Root, backupDir, patchPath string) error {
	relativePath, err := filepath.Localize(patchPath)
	if err != nil {
		return err
	}

	info, err := installed.Lstat(relativePath)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("expected a regular file")
	}

	source, err := installed.Open(relativePath)
	if err != nil {
		return err
	}
	defer source.Close()

	destination := filepath.Join(backupDir, "files", relativePath)
	if err := os.MkdirAll(filepath.Dir(destination), 0o700); err != nil {
		return err
	}

	target, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}

	_, copyErr := io.Copy(target, source)
	closeErr := target.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}
