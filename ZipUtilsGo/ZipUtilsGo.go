package ZipUtilsGo

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// ZipFolder compresses the contents of a folder into a ZIP file.
func ZipFolder(source, target string) error {
	zipFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath := filepath.ToSlash(path[len(source)+1:])
		if relPath == "." {
			return nil
		}

		if info.IsDir() {
			_, err := zipWriter.Create(relPath + "/")
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		writer, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		return err
	})
}

// ZipFile compresses a single file into a ZIP file.
func ZipFile(source, target string) error {
	zipFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()

	writer, err := zipWriter.Create(filepath.Base(source))
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}

// Unzip extracts the contents of a ZIP file into a specified folder.
func Unzip(source, target string) error {
	zipReader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	os.MkdirAll(target, os.ModePerm)

	for _, file := range zipReader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		os.MkdirAll(filepath.Dir(path), os.ModePerm)
		destFile, err := os.Create(path)
		if err != nil {
			return err
		}
		defer destFile.Close()

		srcFile, err := file.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return err
		}
	}

	return nil
}

// UnzipFile extracts a specific file from a ZIP archive.
func UnzipFile(source, targetFile, output string) error {
	zipReader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		if file.Name == filepath.Base(targetFile) {
			srcFile, err := file.Open()
			if err != nil {
				return err
			}
			defer srcFile.Close()

			destFile, err := os.Create(output)
			if err != nil {
				return err
			}
			defer destFile.Close()

			_, err = io.Copy(destFile, srcFile)
			return err
		}
	}

	return os.ErrNotExist
}
