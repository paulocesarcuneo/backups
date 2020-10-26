package archiver

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"path/filepath"
	"io"
	"os"
	"strings"
	"log"
)


func tarFolder(src string, writers ...io.Writer) error {
	info, err := os.Stat(src);
	if err != nil {
		return err
	}

	mw := io.MultiWriter(writers...)

	gzipWriter := gzip.NewWriter(mw)
	defer gzipWriter.Close()

	tarWriter:= tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(src)
	}

	return filepath.Walk(src, func(fileName string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !fileInfo.Mode().IsRegular() {
			return nil
		}

		header, err := tar.FileInfoHeader(fileInfo, fileInfo.Name())

		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(fileName, src))
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		file, err := os.Open(fileName)

		if err != nil {
			return err
		}

		defer file.Close()
		_, err = io.Copy(tarWriter, file)
		log.Println("Adding ",fileName, header)
		return err
	})
}

func untarFolder(dst string, r io.Reader) error {

	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}
		target := filepath.Join(dst, header.Name)

		switch header.Typeflag {

		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		case tar.TypeReg:
			file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(file, tarReader); err != nil {
				return err
			}

			file.Close()
		}
	}
}

func Tar(path string, dst string) (string, error) {
	file, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer file.Close()
	log.Println("making tar gz for", path)
	md5writer := md5.New()
	err = tarFolder(path, file, md5writer)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(md5writer.Sum(nil)), nil
}
