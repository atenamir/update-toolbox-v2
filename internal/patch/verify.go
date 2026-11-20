package patch

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// VerifyPatchArchive checks the archive and returns any verification error.
func VerifyPatchArchive(archivePath string) error {
	_, err := ReadVerifiedPatchManifest(archivePath)
	return err
}

// ReadVerifiedPatchManifest checks the archive and returns its manifest.
func ReadVerifiedPatchManifest(archivePath string) (PatchManifest, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return PatchManifest{}, fmt.Errorf("open patch: %w", err)
	}
	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return PatchManifest{}, fmt.Errorf("open gzip: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)

	var manifest PatchManifest
	manifestFound := false
	archiveFiles := make(map[string]bool)
	archiveHashes := make(map[string]string)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return PatchManifest{}, fmt.Errorf("read archive: %w", err)
		}

		if header.Typeflag != tar.TypeReg && header.Typeflag != tar.TypeRegA {
			return PatchManifest{}, fmt.Errorf("unsupported archive entry: %s", header.Name)
		}

		if header.Name == "manifest.json" {
			if manifestFound {
				return PatchManifest{}, fmt.Errorf("duplicate manifest.json")
			}
			manifestFound = true

			const maxManifestSize = 1 << 20 // 1 MiB
			if header.Size > maxManifestSize {
				return PatchManifest{}, fmt.Errorf("manifest is too large")
			}

			data, err := io.ReadAll(io.LimitReader(tr, maxManifestSize+1))
			if err != nil {
				return PatchManifest{}, fmt.Errorf("read manifest: %w", err)
			}
			if len(data) > maxManifestSize {
				return PatchManifest{}, fmt.Errorf("manifest is too large")
			}

			if err := json.Unmarshal(data, &manifest); err != nil {
				return PatchManifest{}, fmt.Errorf("decode manifest: %w", err)
			}
			continue
		}

		if !strings.HasPrefix(header.Name, "files/") {
			return PatchManifest{}, fmt.Errorf("unexpected archive entry: %s", header.Name)
		}

		relativePath := strings.TrimPrefix(header.Name, "files/")
		if !validPatchPath(relativePath) {
			return PatchManifest{}, fmt.Errorf("unsafe path: %s", header.Name)
		}

		if archiveFiles[relativePath] {
			return PatchManifest{}, fmt.Errorf("duplicate file: %s", relativePath)
		}
		archiveFiles[relativePath] = true

		hasher := sha256.New()
		if _, err := io.Copy(hasher, tr); err != nil {
			return PatchManifest{}, fmt.Errorf("read file %s: %w", relativePath, err)
		}

		archiveHashes[relativePath] = hex.EncodeToString(hasher.Sum(nil))
	}

	if !manifestFound {
		return PatchManifest{}, fmt.Errorf("manifest.json not found")
	}

	if manifest.FromVersion == "" || manifest.ToVersion == "" {
		return PatchManifest{}, fmt.Errorf("patch versions are required")
	}

	seenPaths := make(map[string]bool)

	for _, changed := range manifest.Files {
		if !validPatchPath(changed.Path) {
			return PatchManifest{}, fmt.Errorf("unsafe manifest path: %s", changed.Path)
		}

		if seenPaths[changed.Path] {
			return PatchManifest{}, fmt.Errorf("duplicate manifest path: %s", changed.Path)
		}
		seenPaths[changed.Path] = true

		switch changed.Status {
		case "A", "M":
			if !archiveFiles[changed.Path] {
				return PatchManifest{}, fmt.Errorf("missing file: %s", changed.Path)
			}

			if len(changed.SHA256) != hex.EncodedLen(sha256.Size) {
				return PatchManifest{}, fmt.Errorf("invalid SHA256 for %s", changed.Path)
			}

			expectedHash, err := hex.DecodeString(changed.SHA256)
			if err != nil {
				return PatchManifest{}, fmt.Errorf("invalid SHA256 for %s", changed.Path)
			}

			actualHash, err := hex.DecodeString(archiveHashes[changed.Path])
			if err != nil {
				return PatchManifest{}, fmt.Errorf("invalid calculated SHA256 for %s", changed.Path)
			}

			if !equalSHA256(expectedHash, actualHash) {
				return PatchManifest{}, fmt.Errorf("SHA256 mismatch for %s", changed.Path)
			}

		case "D":
			if archiveFiles[changed.Path] {
				return PatchManifest{}, fmt.Errorf("deleted file must not be archived: %s", changed.Path)
			}
			if changed.SHA256 != "" {
				return PatchManifest{}, fmt.Errorf("deleted file must not have SHA256: %s", changed.Path)
			}

		default:
			return PatchManifest{}, fmt.Errorf("unsupported status %q", changed.Status)
		}
	}

	for filePath := range archiveFiles {
		if !seenPaths[filePath] {
			return PatchManifest{}, fmt.Errorf("unexpected file: %s", filePath)
		}
	}

	return manifest, nil
}

func equalSHA256(a, b []byte) bool {
	if len(a) != sha256.Size || len(b) != sha256.Size {
		return false
	}

	var difference byte
	for i := range a {
		difference |= a[i] ^ b[i]
	}
	return difference == 0
}