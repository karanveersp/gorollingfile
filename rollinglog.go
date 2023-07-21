package gorollinglog

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// RollingLogFile implements io.Writer to allow logging to rolling files.
// The active file will always have the given name, but backup files
// will have bkp_1, bkp_2, ... appended to the name.
type RollingLogFile struct {
	LogDir   string
	FileName string
	Path     string

	maxSizeBytes int64
	backupCount  int
	fileHandle   *os.File
}

// NewRollingFile creates a new rolling file struct.
func NewRollingFile(dir, name string, mode rune, maxSizeMb float64, backups int) (*RollingLogFile, error) {
	path := filepath.Join(dir, name)

	var fhandle *os.File
	var err error
	switch mode {
	case 'w':
		fhandle, err = os.Create(path)
	case 'a':
		fhandle, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	default:
		return nil, fmt.Errorf("unrecognized mode '%c'. Pass 'w' for overwrite and 'a' for append", mode)
	}
	if err != nil {
		return nil, err
	}

	return &RollingLogFile{
		LogDir:       dir,
		FileName:     name,
		Path:         path,
		maxSizeBytes: megabytesToBytes(maxSizeMb),
		backupCount:  backups,
		fileHandle:   fhandle,
	}, nil
}

func (f *RollingLogFile) rotate() error {
	f.fileHandle.Close()
	parts := strings.Split(f.FileName, ".")
	fname, ext := parts[0], parts[1]

	fstBkpName := fmt.Sprintf("%s_bkp_%d.%s", fname, 1, ext)
	fstBkpPath := filepath.Join(f.LogDir, fstBkpName)
	var err error
	if !fileExists(fstBkpPath) {

		// create first backup
		err := copyFile(f.Path, fstBkpPath)
		if err != nil {
			return err
		}
	} else {
		// cascade copy, starting from oldest file
		for i := f.backupCount - 1; i > 0; i-- {
			bkpName := fmt.Sprintf("%s_bkp_%d.%s", fname, i, ext)
			backupPath := filepath.Join(f.LogDir, bkpName)
			var newerBkpName string
			if i-1 == 0 {
				newerBkpName = f.FileName
			} else {
				newerBkpName = fmt.Sprintf("%s_bkp_%d.%s", fname, i-1, ext)
			}
			newerBackupPath := filepath.Join(f.LogDir, newerBkpName)
			if fileExists(newerBackupPath) {
				// overwrite nth file with n-1th file, effectively discarding nth file
				err := copyFile(newerBackupPath, backupPath)
				if err != nil {
					return err
				}
			}
		}
	}
	// reopen file for fresh logging
	newHandle, err := os.Create(f.Path)
	if err != nil {
		return err
	}
	f.fileHandle = newHandle
	return nil
}

// Close will close the open file handle.
func (f *RollingLogFile) Close() error {
	return f.fileHandle.Close()
}

func (f *RollingLogFile) Write(data []byte) (int, error) {
	i, err := f.fileHandle.Write(data)
	if err != nil {
		return i, err
	}
	size, err := getFileSizeBytes(f.Path)
	if err != nil {
		return 0, err
	}
	if size >= f.maxSizeBytes {
		err := f.rotate()
		if err != nil {
			return 0, err
		}
	}
	return i, nil
}

func getFileSizeBytes(path string) (int64, error) {
	finfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return finfo.Size(), nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return !info.IsDir()
}

func copyFile(src, dest string) error {
	// create first backup
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()
	destf, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destf.Close()
	_, err = io.Copy(destf, sf)
	if err != nil {
		return err
	}
	return nil
}

func megabytesToBytes(megabytes float64) int64 {
	bytes := int64(megabytes * 1024 * 1024)
	return bytes
}
