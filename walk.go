package replace

// TODO
// SubvertPathname renames a file or directory by replacing the old name with the new name.
// func SubvertPathname(filePath string, perm int, o string, n string) error {
// }

import (
	"os"
	"path/filepath"
)

// Checks whether a path is excluded based on the exclusion list
func isExcluded(path string, excl []string) bool {
	for _, ex := range excl {
		if matched, _ := filepath.Match(ex, path); matched {
			return true
		}
	}
	return false
}

// Rlist returns a list of files and directories in the given path.
// The returned list of files and directories are relative to the given path.
// excl is a list of glob patterns to exclude.
func Rlist(path string, excl []string) (files []string, dirs []string, err error) {
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Ignore the root path
		if path == "." {
			return nil
		}
		// Ignore specified paths
		for _, exclusion := range excl {
			if matched, _ := filepath.Match(exclusion, filepath.Base(path)); matched {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}
		// Ignore device and pipe files
		if (info.Mode()&os.ModeDevice) != 0 || (info.Mode()&os.ModeNamedPipe) != 0 {
			return nil
		}

		if info.IsDir() {
			dirs = append(dirs, path)
		} else if info.Mode()&os.ModeSymlink != 0 {
			files = append(files, path)
		} else {
			files = append(files, path)
		}

		return nil
	})

	return files, dirs, err
}
