package patch

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// ========== runGit executes Git commands in the test repository=================
func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, output)
	}
}

func TestGetChangedFiles(t *testing.T) {
	//================= Create a temporary directory====================
	repoPath := t.TempDir()

	// Initialize a Git repository.
	runGit(t, repoPath, "init")
	runGit(t, repoPath, "config", "user.name", "Test User")
	runGit(t, repoPath, "config", "user.email", "test@example.com")

	appDir := filepath.Join(repoPath, "app")

	if err := os.MkdirAll(appDir, 0755); err != nil {
		t.Fatal(err)
	}

	//================== Version 1: create two files================
	userFile := filepath.Join(appDir, "User.php")
	oldFile := filepath.Join(appDir, "OldService.php")
	newFile := filepath.Join(appDir, "NewService.php")

	if err := os.WriteFile(userFile, []byte("Version 1"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(oldFile, []byte("Old Service"), 0644); err != nil {
		t.Fatal(err)
	}

	runGit(t, repoPath, "add", ".")
	runGit(t, repoPath, "commit", "-m", "Version 1")
	runGit(t, repoPath, "tag", "v1.0.0")

	//================ Version 2: modify, add and delete files =================
	if err := os.WriteFile(userFile, []byte("Version 2"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(newFile, []byte("New Service"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.Remove(oldFile); err != nil {
		t.Fatal(err)
	}

	runGit(t, repoPath, "add", "-A")
	runGit(t, repoPath, "commit", "-m", "Version 2")
	runGit(t, repoPath, "tag", "v1.1.0")

	//================ Analyze changes between the two tags=======================
	files, err := GetChangedFiles(repoPath, "v1.0.0", "v1.1.0")
	if err != nil {
		t.Fatalf("GetChangedFiles failed: %v", err)
	}

	expected := map[string]string{
		"app/NewService.php": "A",
		"app/OldService.php": "D",
		"app/User.php":       "M",
	}

	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d", len(expected), len(files))
	}

	for _, file := range files {
		status, exists := expected[file.Path]

		if !exists {
			t.Errorf("unexpected file: %s", file.Path)
			continue
		}

		if file.Status != status {
			t.Errorf(
				"file %s: expected %s, got %s",
				file.Path,
				status,
				file.Status,
			)
		}
	}
}

func TestGetFileAtTag(t *testing.T) {
	repoPath := filepath.Join("..", "..", "..", "patch-lab")

	data, err := GetFileAtTag(
		repoPath,
		"v1.1.0",
		"app/User.php",
	)

	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	content := string(data)

	if !strings.Contains(content, "Version 2") {
		t.Fatalf("unexpected file content: %s", content)
	}
}
