package replace

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// SubvertFileContent replaces all occurrences of the old string with the new string in the file content.
func SubvertFileContent(filePath string, o string, n string) error {
	stat, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	perm := stat.Mode().Perm()
	file, err := os.OpenFile(filePath, os.O_RDWR, perm)
	if err != nil {
		return err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	newContent := Subvert(string(content), o, n)
	if string(content) == newContent {
		return nil
	}
	err = file.Truncate(0) // maybe a  bit dangerous
	if err != nil {
		return err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = file.WriteString(newContent)
	if err != nil {
		return err
	}

	return nil
}

// TODO
// SubvertPathname renames a file or directory by replacing the old name with the new name.
// func SubvertPathname(filePath string, perm int, o string, n string) error {
// }

// WalkAndSubvert walks the path and subverts all files tha have the old string in them.
func WalkAndSubvert(path string, excludes []string, o string, n string) error {
	return filepath.Walk(path, func(fp string, info os.FileInfo, err error) error {
		for _, exclude := range excludes {
			if strings.Contains(fp, exclude) {
				return nil
			}
		}

		// If it's a regular file, add it to the list
		if info.Mode().IsRegular() {
			err = SubvertFileContent(fp, o, n)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
