package patch

import (
	"fmt"
	"os/exec"
	"strings"
)

type ChangedFile struct {
	Path   string
	Status string
}

func GetChangedFiles(repoPath, fromTag, toTag string) ([]ChangedFile, error) {
	//================= Run Git inside the specified repository ==================
	cmd := exec.Command(
		"git",
		"diff",
		"--name-status",
		"-z",
		"--no-renames",
		fromTag,
		toTag,
		"--",
	)

	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git diff failed: %w", err)
	}

	// ================= With -z, Git separates fields using NUL bytes ======================
	if len(output) == 0 {
		return []ChangedFile{}, nil
	}

	fields := strings.Split(
		strings.TrimSuffix(string(output), "\x00"),
		"\x00",
	)

	if len(fields)%2 != 0 {
		return nil, fmt.Errorf("invalid git diff output")
	}

	changedFiles := make([]ChangedFile, 0, len(fields)/2)

	for i := 0; i < len(fields); i += 2 {
		status := fields[i]
		path := fields[i+1]

		switch status {
		case "A", "M", "D", "T":
			changedFiles = append(changedFiles, ChangedFile{
				Path:   path,
				Status: status,
			})
		default:
			return nil, fmt.Errorf("unsupported git status %q", status)
		}
	}

	return changedFiles, nil
}

func GetFileAtTag(repoPath string, tag string, filePath string) ([]byte, error) {
	cmd := exec.Command(
		"git",
		"show",
		tag+":"+filePath,
	)

	cmd.Dir = repoPath

	data, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf(
			"read file %q at tag %q: %w",
			filePath,
			tag,
			err,
		)
	}

	return data, nil
}
