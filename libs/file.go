// 文件操作，linux下一切皆文件，目录和文件的很多操作方法类似

package libs

import (
	"archive/tar"
	"compress/bzip2"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
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
			err = os.MkdirAll(fileName, f.Mode())
			if err != nil {
				return err
			}

		} else {
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

	fs, err := f.Stat()
	if err != nil {
		return err
	}

	defer f.Close()

	destf, err := os.Create(destFile)
	if err != nil {
		return err
	}

	defer destf.Close()

	err = destf.Chmod(fs.Mode())
	if err != nil {
		return err
	}

	_, err = io.Copy(destf, f)
	if err != nil {
		return err
	}

	return nil
}

// EmptyDir 清空目录内容
func EmptyDir(dir string) (err error) {
	fm, err := GetFileMode(dir)
	if err != nil {
		return err
	}

	if err = os.RemoveAll(dir); err != nil {
		return err
	}

	if err = os.MkdirAll(dir, fm); err != nil {
		return err
	}

	return nil
}

// GetFileMode 获取文件或者目录的权限
func GetFileMode(name string) (fm os.FileMode, err error) {
	f, err := os.Open(name)
	if err != nil {
		return fm, err
	}

	defer f.Close()

	fs, err := f.Stat()
	if err != nil {
		return fm, err
	}

	fm = fs.Mode()
	return fm, nil
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

// MD5File 计算文件MD5值
func MD5File(filename string) (md5v string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return md5v, err
	}

	defer f.Close()

	body, err := ioutil.ReadAll(f)
	if err != nil {
		return md5v, err
	}

	md5v = fmt.Sprintf("%x", md5.Sum(body))
	return md5v, nil
}

// GetFileList 获取某个目录下的目录列表和文件列表
func GetFileList(dirPath string) (dirs, files []string, err error) {
	l, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return dirs, files, err
	}

	for _, f := range l {
		if f.IsDir() {
			dirs = append(dirs, f.Name())
		} else {
			files = append(files, f.Name())
		}
	}

	return dirs, files, nil
}

// UnGzip 解压文件：tar.gz
func UnGzip(compressFile, destDir string) (err error) {
	f, err := os.Open(compressFile)
	if err != nil {
		return err
	}

	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}

	defer gr.Close()

	err = UnTar(gr, destDir)
	if err != nil {
		return err
	}

	return nil
}

// UnBzip2 解压文件：tar.bzip2
func UnBzip2(compressFile, destDir string) (err error) {
	f, err := os.Open(compressFile)
	if err != nil {
		return err
	}

	defer f.Close()

	br := bzip2.NewReader(f)
	err = UnTar(br, destDir)
	if err != nil {
		return err
	}

	return nil
}

// Untar 读取tar io.Reader，并写入目录
func UnTar(rd io.Reader, destDir string) (err error) {
	tr := tar.NewReader(rd)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		filename := path.Join(destDir, hdr.Name)
		hdrF := hdr.FileInfo()
		if hdrF.IsDir() {
			err = os.MkdirAll(filename, hdrF.Mode())
			if err != nil {
				return err
			}

			continue
		}

		file, err := os.Create(filename)
		if err != nil {
			return err
		}

		err = file.Chmod(hdrF.Mode())
		if err != nil {
			return err
		}

		_, err = io.Copy(file, tr)
		if err != nil {
			return err
		}
	}

	return nil
}
