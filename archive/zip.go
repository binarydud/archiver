package archive

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Archiver struct {
	filepath string
	// filewriter *os.File
	writer  *zip.Writer
	zipFile *os.File
}

func assertValidFile(infilename string) (os.FileInfo, error) {
	fi, err := os.Stat(infilename)
	if err != nil && os.IsNotExist(err) {
		return fi, fmt.Errorf("could not archive missing file: %s", infilename)
	}
	return fi, err
}

func GetZipArchiver(path string) *Archiver {

	return &Archiver{
		filepath: path,
	}
}

func (a *Archiver) Start() error {
	zipFile, err := os.Create(a.filepath)
	if err != nil {
		return err
	}
	a.zipFile = zipFile
	a.writer = zip.NewWriter(a.zipFile)
	return nil
}

func (a *Archiver) AddFile(filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := a.writer.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}
func (a *Archiver) AddItem(source string) error {
	info, err := os.Stat(source)
	if err != nil {
		return nil
	}
	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}
	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		header.SetModTime(time.Time{})
		/*
			// check if file is executable
			inf, err := os.Stat(path)
			if err != nil {
				return err
			}
			mode := inf.Mode().Perm()
			if mode&0100 == 0100 {
				header.SetMode(0555)
			} else {
				header.SetMode(0444)
			}
		*/

		writer, err := a.writer.CreateHeader(header)
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

	return err
}

func (a *Archiver) GetHash() (string, error) {
	data, err := ioutil.ReadFile(a.filepath)
	if err != nil {
		return "", fmt.Errorf("could not compute file '%s' checksum: %s", a.filepath, err)
	}
	return genSha(data)
}

/*
func (a *Archiver) open() error {
	f, err := os.Create(a.filepath)
	if err != nil {
		return err
	}
	//a.filewriter = f
	a.writer = zip.NewWriter(f)
	return nil
}
*/

func (a *Archiver) CloseWriter() {
	a.close()
	a.zipFile.Close()
}
func (a *Archiver) close() {
	if a.writer != nil {
		a.writer.Close()
		a.writer = nil
	}
}
