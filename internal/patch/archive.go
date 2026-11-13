package patch

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

// writeTarEntry adds one regular file to a TAR archive.
func writeTarEntry(tw *tar.Writer, name string, data []byte) error {
	header := &tar.Header{
		Name: name,
		Mode: 0644,
		Size: int64(len(data)),
	}

	if err := tw.WriteHeader(header); err != nil {
		return fmt.Errorf("write header %q: %w", name, err)
	}

	if _, err := tw.Write(data); err != nil {
		return fmt.Errorf("write content %q: %w", name, err)
	}

	return nil
}

// validPatchPath checks the path before adding it to the archive.
func validPatchPath(name string) bool {
	return name != "" &&
		name != "." &&
		!strings.HasPrefix(name, "/") &&
		!strings.Contains(name, "\\") &&
		!strings.Contains(name, ":") &&
		!strings.ContainsRune(name, '\x00') &&
		path.Clean(name) == name &&
		name != ".." &&
		!strings.HasPrefix(name, "../")
}

func CreatePatchArchive(
	manifest PatchManifest,
	repoPath string,
	outputPath string,
) (err error) {

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create archive: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	gzipWriter := gzip.NewWriter(file)
	tarWriter := tar.NewWriter(gzipWriter)

	defer func() {
		if closeErr := gzipWriter.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	defer func() {
		if closeErr := tarWriter.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	// Add the patch manifest.
	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("encode manifest: %w", err)
	}

	if err := writeTarEntry(tarWriter, "manifest.json", manifestData); err != nil {
		return err
	}

	// Add changed files from the target Git tag.
	for _, changed := range manifest.Files {
		if !validPatchPath(changed.Path) {
			return fmt.Errorf("unsafe file path: %q", changed.Path)
		}

		switch changed.Status {
		case "D":
			// Deleted files are recorded in the manifest only.
			continue

		case "A", "M":
			data, err := GetFileAtTag(
				repoPath,
				manifest.ToVersion,
				changed.Path,
			)
			if err != nil {
				return err
			}

			archivePath := "files/" + changed.Path

			if err := writeTarEntry(tarWriter, archivePath, data); err != nil {
				return err
			}

		default:
			return fmt.Errorf(
				"unsupported patch status %q for %q",
				changed.Status,
				changed.Path,
			)
		}
	}

	return nil
}
