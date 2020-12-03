package deployer

import (
	"archive/zip"
	"encoding/base64"
	"fmt"
	"github.com/rknizzle/faas/internal/models"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// containerDeployer contains all the methods required to turn function code into a container image
// and then sent to a container registry where it can be pulled for invocation
type containerDeployer interface {
	BuildImage(string, string) error
	PushImage(string) error
}

// Deployer handles deploying new functions that a user submits by unpacking the function data from
// a users request and then using a ContainerDeployer to build and push a container image for later
// invocation
type Deployer struct {
	c containerDeployer
}

// NewDeployer initializes a Deployer with a ContainerDeploy for building and pushing images
func NewDeployer(c containerDeployer) Deployer {
	return Deployer{c}
}

// Deploy unpacks the function code and then builds and pushes a container image
func (d Deployer) Deploy(data models.FnData) error {
	dir, err := d.e.FnDirFromBase64Data(data.Name, data.File)
	if err != nil {
		return err
	}

	tag := data.Name
	d.c.BuildImage(dir, tag)
	// d.c.PushImage()

	// remove temporary directory used to build the image
	os.RemoveAll(dir)

	return nil
}

// fnDirFromBase64Data decodes a base64 string into a directory containing the function code
func (d Deployer) FnDirFromBase64Data(fnName string, base64Data string) (string, error) {
	filename, err := writeDataToZip(fnName, base64Data)
	if err != nil {
		return "", err
	}

	// unzip the file into a directory called the function name
	dir := fnName
	_, err = Unzip(filename, dir)
	if err != nil {
		return "", err
	}

	// remove the zip file now that the directory has been extracted
	err = os.Remove(filename)
	if err != nil {
		return "", err
	}

	return dir, nil
}

// writeDataToZip converts a base64 encoding string back into the original zip file
func writeDataToZip(fnName string, fnData string) (string, error) {
	filename := fnName + ".zip"

	zipFile, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(fnData))

	// write decoded data to zip file
	io.Copy(zipFile, decoder)
	return filename, nil
}

// Unzip a zip file into its file contents
func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
