package breakdown

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GenerateFilesystemStructure recursively creates folders and files for the task plan.
func (n *Node) GenerateFilesystemStructure(baseDir string) error {
	return n.generateWithIndex(baseDir, -1)
}

func (n *Node) generateWithIndex(baseDir string, index int) error {
	// Sanitize task title for use as a filename/directory name
	safeName := sanitizeFilename(n.Task)

	if index >= 0 {
		safeName = fmt.Sprintf("%02d-%s", index+1, safeName)
	}

	var currentPath string
	if n.Type == TaskTypeComposite {
		// Create a directory for composite tasks
		currentPath = filepath.Join(baseDir, safeName)
		if err := os.MkdirAll(currentPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", currentPath, err)
		}

		// Create a README file for the composite task's details
		readmePath := filepath.Join(currentPath, "README.md")
		content := fmt.Sprintf("# %s\n\n%s\n", n.Task, n.Details)
		if err := os.WriteFile(readmePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create README %s: %w", readmePath, err)
		}

		// Recurse for children
		for i, child := range n.Children {
			if err := child.generateWithIndex(currentPath, i); err != nil {
				return err
			}
		}
	} else {
		// Create a file for atomic tasks
		currentPath = filepath.Join(baseDir, safeName+".md")
		content := fmt.Sprintf("# %s\n\n%s\n", n.Task, n.Details)
		if err := os.WriteFile(currentPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %w", currentPath, err)
		}
	}

	return nil
}

func sanitizeFilename(name string) string {
	sanitized := strings.ReplaceAll(name, "/", "-")
	sanitized = strings.ReplaceAll(sanitized, " ", "_")
	if len(sanitized) > 50 {
		sanitized = sanitized[:50]
	}
	return sanitized
}
