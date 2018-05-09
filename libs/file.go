package libs

import (
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// CopyDir 复制文件夹及其子目录和文件到目标文件下
func CopyDir(srcDir, destDir string) (err error) {
	err = filepath.Walk(srcDir, func(src string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		rel, err := filepath.Rel(srcDir, src)
		if err != nil {
			return err
		}

		fileName := path.Join(destDir, path.Base(srcDir), rel)

		if f.IsDir() {
			err = os.MkdirAll(fileName, 0755)
			if err != nil {
				return err
			}

		} else {
			println(fileName)
			err = CopyFile(src, fileName)
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

// HttpDownFile 下载http文件
func HttpDownFile(url, dest string) (err error) {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	fileName := path.Base(url)
	destFile := path.Join(dest, fileName)
	f, err := os.Create(destFile)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return err
	}

	return nil
}
