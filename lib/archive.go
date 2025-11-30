package lib

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ExtractArchive extracts supported archive formats (zip, tar, tar.gz, tgz) to destDir.
// It respects the limits provided (e.g. max archive size is handled by caller generally, 
// but individual file extraction limits could be added here if needed).
func ExtractArchive(file *os.File, filename string, destDir string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	
	// Handle .tar.gz or .tgz
	if ext == ".gz" || ext == ".tgz" {
		return extractTarGz(file, destDir)
	}
	if ext == ".tar" {
		return extractTar(file, destDir)
	}
	if ext == ".zip" {
		stat, err := file.Stat()
		if err != nil {
			return err
		}
		return extractZip(file, stat.Size(), destDir)
	}
	
	return fmt.Errorf("unsupported archive format: %s", ext)
}

func extractZip(reader io.ReaderAt, size int64, destDir string) error {
	r, err := zip.NewReader(reader, size)
	if err != nil {
		return err
	}

	for _, f := range r.File {
		// Skip root directory entry
		if f.Name == "." || f.Name == "" {
			continue
		}
		fpath := filepath.Join(destDir, f.Name)
		cleanDest := filepath.Clean(destDir)
		// Allow exact match or paths that start with destDir + separator
		if fpath != cleanDest && !strings.HasPrefix(fpath, cleanDest+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func extractTar(reader io.Reader, destDir string) error {
	tr := tar.NewReader(reader)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Skip root directory entry
		if header.Name == "." || header.Name == "" {
			continue
		}
		fpath := filepath.Join(destDir, header.Name)
		cleanDest := filepath.Clean(destDir)
		// Allow exact match or paths that start with destDir + separator
		if fpath != cleanDest && !strings.HasPrefix(fpath, cleanDest+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}
	return nil
}

func extractTarGz(reader io.Reader, destDir string) error {
	gzr, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer gzr.Close()
	return extractTar(gzr, destDir)
}

// CreateArchive creates an archive of the sourceDir in the specified format at outputPath.
func CreateArchive(sourceDir string, format string, outputPath string) error {
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	switch format {
	case "zip":
		return createZip(sourceDir, outFile)
	case "tar":
		return createTar(sourceDir, outFile)
	case "tar.gz", "tgz":
		return createTarGz(sourceDir, outFile)
	default:
		return fmt.Errorf("unsupported archive format for creation: %s", format)
	}
}

func createZip(sourceDir string, w io.Writer) error {
	zw := zip.NewWriter(w)
	defer zw.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Set the relative path for the header name
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		// Skip root directory entry
		if relPath == "." {
			return nil
		}

		// Create a header based on the file info
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})
}

func createTar(sourceDir string, w io.Writer) error {
	tw := tar.NewWriter(w)
	defer tw.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		header.Name = relPath

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(tw, file)
		return err
	})
}

func createTarGz(sourceDir string, w io.Writer) error {
	gw := gzip.NewWriter(w)
	defer gw.Close()
	
	// We pass the gzip writer to createTar, which wraps it in a tar writer
	// createTar closes the tar writer, and defer gw.Close() closes gzip
	return createTar(sourceDir, gw)
}

// DetectArchiveFormat returns "zip", "tar", "tar.gz", or "" if unknown
func DetectArchiveFormat(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == ".zip" {
		return "zip"
	}
	if ext == ".tar" {
		return "tar"
	}
	if ext == ".tgz" {
		return "tar.gz"
	}
	if ext == ".gz" && strings.HasSuffix(strings.ToLower(filename), ".tar.gz") {
		return "tar.gz"
	}
	return ""
}

