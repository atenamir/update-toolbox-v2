package patch

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ========== ApplyPatchToCopy applies a verified patch to a temporary copy of an installation ================
func ApplyPatchToCopy(archivePath, installDir string) (string, string, error) {
	workDir, err := os.MkdirTemp("", "toolbox-apply-*")
	if err != nil {
		return "", "", err
	}

	complete := false
	defer func() {
		if !complete {
			os.RemoveAll(workDir)
		}
	}()

	// Verify and apply the same local copy of the patch.
	archiveData, err := os.ReadFile(archivePath)
	if err != nil {
		return "", "", fmt.Errorf("read patch: %w", err)
	}
	snapshotPath := filepath.Join(workDir, "patch.tar.gz")
	if err := os.WriteFile(snapshotPath, archiveData, 0o600); err != nil {
		return "", "", fmt.Errorf("save patch snapshot: %w", err)
	}
	defer os.Remove(snapshotPath)

	manifest, err := ReadVerifiedPatchManifest(snapshotPath)
	if err != nil {
		return "", "", err
	}

	for _, file := range manifest.Files {
		if strings.EqualFold(file.Path, "version.txt") {
			return "", "", fmt.Errorf("patch cannot modify the installation version marker")
		}
	}

	updatedDir := filepath.Join(workDir, "installation")
	if err := os.CopyFS(updatedDir, os.DirFS(installDir)); err != nil {
		return "", "", fmt.Errorf("copy installation: %w", err)
	}

	versionFile := filepath.Join(updatedDir, "version.txt")
	if err := CheckInstalledVersion(manifest, versionFile); err != nil {
		return "", "", err
	}

	backupDir, err := CreateBackup(manifest, updatedDir)
	if err != nil {
		return "", "", err
	}

	root, err := os.OpenRoot(updatedDir)
	if err != nil {
		return "", backupDir, err
	}
	defer root.Close()

	if err := applyArchiveFiles(root, snapshotPath, manifest); err != nil {
		return "", backupDir, err
	}

	if err := root.WriteFile("version.txt", []byte(manifest.ToVersion+"\n"), 0o644); err != nil {
		return "", backupDir, fmt.Errorf("write installed version: %w", err)
	}

	complete = true
	return updatedDir, backupDir, nil
}

func applyArchiveFiles(root *os.Root, archivePath string, manifest PatchManifest) error {
	archive, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	gz, err := gzip.NewReader(archive)
	if err != nil {
		return err
	}
	defer gz.Close()

	actions := make(map[string]string, len(manifest.Files))
	for _, file := range manifest.Files {
		actions[file.Path] = file.Status
	}

	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if header.Name == "manifest.json" {
			continue
		}

		patchPath := strings.TrimPrefix(header.Name, "files/")
		relativePath, err := filepath.Localize(patchPath)
		if err != nil {
			return err
		}

		flags := os.O_WRONLY | os.O_TRUNC
		if actions[patchPath] == "A" {
			if err := root.MkdirAll(filepath.Dir(relativePath), 0o755); err != nil {
				return err
			}
			flags = os.O_WRONLY | os.O_CREATE | os.O_EXCL
		}

		target, err := root.OpenFile(relativePath, flags, 0o644)
		if err != nil {
			return fmt.Errorf("open destination %q: %w", patchPath, err)
		}
		_, copyErr := io.Copy(target, tr)
		closeErr := target.Close()
		if copyErr != nil {
			return fmt.Errorf("write %q: %w", patchPath, copyErr)
		}
		if closeErr != nil {
			return fmt.Errorf("close %q: %w", patchPath, closeErr)
		}
	}

	for _, file := range manifest.Files {
		if file.Status == "D" {
			relativePath, err := filepath.Localize(file.Path)
			if err != nil {
				return err
			}
			if err := root.Remove(relativePath); err != nil {
				return fmt.Errorf("delete %q: %w", file.Path, err)
			}
		}
	}

	return nil
}
