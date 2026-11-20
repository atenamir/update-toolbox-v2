package patch

// CreatePatch builds and publishes a verified patch archive.
func CreatePatch(manifestPath, repoPath, outputPath string) (PatchManifest, error) {
	manifests, err := LoadManifest(manifestPath)
	if err != nil {
		return PatchManifest{}, err
	}

	if err := ValidateManifest(manifests); err != nil {
		return PatchManifest{}, err
	}

	fromTag := manifests[0].Version
	toTag := manifests[1].Version

	changedFiles, err := GetChangedFiles(repoPath, fromTag, toTag)
	if err != nil {
		return PatchManifest{}, err
	}

	patchManifest, err := BuildPatchManifestWithHashes(
		repoPath,
		fromTag,
		toTag,
		changedFiles,
	)
	if err != nil {
		return PatchManifest{}, err
	}

	if err := CreatePatchArchiveSafe(patchManifest, repoPath, outputPath); err != nil {
		return PatchManifest{}, err
	}

	return patchManifest, nil
}
