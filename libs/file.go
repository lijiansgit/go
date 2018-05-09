package libs

import (
	"io"
	"os"
	"path/filepath"
)

// CopyDir 复制文件夹及其子目录和文件
func CopyDir(srcDir, destDir string) (err error) {
	err = filepath.Walk(srcDir, func(src string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		if f.IsDir() {
			dir := destDir + f.Name()
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				return err
			}
		} else {
			file := destDir + f.Name()
			err = CopyFile(src, file)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return nil
}

// CopyFile 复制单个文件
func CopyFile(srcFile, destFile string) (err error) {
	f, err := os.Open(srcFile)
	if err != nil {
		return err
	}

	defer f.Close()

	destf, err := os.Create(destFile)
	if err != nil {
		return err
	}

	defer destf.Close()

	_, err = io.Copy(destf, f)
	if err != nil {
		return err
	}

	return nil
}
