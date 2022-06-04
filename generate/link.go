package generate

import (
	"fmt"
	"os"
	"path/filepath"
)

func fileExists(name string) bool {
	info, err := os.Stat(name)
	return err == nil && info != nil
}

func getBackupName(name string, i int) (string, error) {
	dir, file := filepath.Split(name)

	backupName := filepath.Join(
		dir,
		fmt.Sprintf("%s.%d.backup", file, i),
	)

	exists := fileExists(backupName)
	if !exists {
		return backupName, nil
	}

	return getBackupName(name, i+1)
}

func prepareLink(name string) error {
	exists := fileExists(name)
	lInfo, err1 := os.Lstat(name)
	isLink := err1 == nil && lInfo != nil && !lInfo.Mode().IsRegular()

	if isLink {
		return os.Remove(name)
	}

	if !exists {
		if err := makeDirAll(name); err != nil {
			return err
		}
		return nil
	}

	// getting here means the file exists and is not a link
	backupName, err2 := getBackupName(name, 0)
	if err2 != nil {
		return err2
	}

	return os.Rename(name, backupName)
}

func makeLink(relativeName, destinationName string) error {
	home, err1 := os.UserHomeDir()
	if err1 != nil {
		return err1
	}

	absName, err2 := filepath.Abs(destinationName)
	if err2 != nil {
		return err2
	}

	linkName := filepath.Join(home, relativeName)
	if err := prepareLink(linkName); err != nil {
		return err
	}

	return os.Symlink(absName, linkName)
}
