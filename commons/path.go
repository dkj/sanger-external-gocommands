package commons

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	irodsclient_fs "github.com/cyverse/go-irodsclient/fs"
)

func MakeIRODSPath(cwd string, homedir string, zone string, path string) string {
	if strings.HasPrefix(path, fmt.Sprintf("/%s/~", zone)) {
		// compat to icommands
		// relative path from user's home
		partLen := 3 + len(zone)
		newPath := filepath.Join(homedir, path[partLen:])
		return filepath.Clean(newPath)
	}

	if strings.HasPrefix(path, "/") {
		// absolute path
		return filepath.Clean(path)
	}

	if strings.HasPrefix(path, "~") {
		// relative path from user's home
		newPath := filepath.Join(homedir, path[1:])
		return filepath.Clean(newPath)
	}

	// relative path from current woring dir
	newPath := filepath.Join(cwd, path)
	return filepath.Clean(newPath)
}

func MakeLocalPath(path string) string {
	if strings.HasPrefix(path, "/") {
		return filepath.Clean(path)
	}

	wd, _ := os.Getwd()

	newPath := filepath.Join(wd, path)
	return filepath.Clean(newPath)
}

func EnsureTargetIRODSFilePath(filesystem *irodsclient_fs.FileSystem, source string, target string) string {
	if filesystem.ExistsDir(target) {
		// make full file name for target
		filename := filepath.Base(source)
		return filepath.Join(target, filename)
	}
	return target
}

func EnsureTargetLocalFilePath(source string, target string) string {
	st, err := os.Stat(target)
	if err == nil {
		if st.IsDir() {
			// make full file name for target
			filename := filepath.Base(source)
			return filepath.Join(target, filename)
		}
	}
	return target
}

func GetFileExtension(path string) string {
	base := filepath.Base(path)

	idx := strings.Index(base, ".")
	if idx >= 0 {
		return path[idx:]
	}
	return path
}
