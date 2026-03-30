package util

import (
	"io"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/sys/windows"
)

func CreateHiddenFolder(folderPath string) error {
    err := os.MkdirAll(folderPath, 0o755)
    if err != nil {
        return err
    }

    // Windows only
    if runtime.GOOS == "windows" {
		ptr, err := windows.UTF16PtrFromString(folderPath)
		if err != nil {
			return err
		}

		attrs, err := windows.GetFileAttributes(ptr)
		if err != nil {
			return err
		}

		// Add BOTH hidden + system
		attrs |= windows.FILE_ATTRIBUTE_HIDDEN | windows.FILE_ATTRIBUTE_SYSTEM
        windows.SetFileAttributes(ptr, attrs)
    }

    return nil
}

func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

func RemoveAndCreateDir(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	return os.MkdirAll(path, 0o755)
}

func AtomicWriteFile(path string, data []byte, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, perm); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func CopyFile(src string, dst string, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
