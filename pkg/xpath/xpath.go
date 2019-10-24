package xpath

import (
	"os"
	"os/user"
)

// IsExists if $path is a dir or file, return the true
func IsExists(path string) (_ bool, _ error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// IsDir check $dir is a directory or not
func IsDir(dir string) bool {
	fi, err := os.Stat(dir)
	return err == nil && fi.IsDir()
}

// IsFile check $file is a regular file or not
func IsFile(file string) bool {
	fi, err := os.Stat(file)
	return err == nil && fi.Mode().IsRegular()
}

// HomeDir return the user home directory
func HomeDir() (home string) {
	u, err := user.Current()
	if err != nil {
		return "/tmp"
	}
	return u.HomeDir
}
