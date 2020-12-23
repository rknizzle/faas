package deployer

import (
	"bytes"
	"encoding/base64"
	"io"
	"strings"

	"github.com/rknizzle/faas/internal/models"
	"github.com/spf13/afero"
)

// containerDeployer contains all the methods required to turn function code into a container image
// and then sent to a container registry where it can be pulled for invocation
type containerDeployer interface {
	BuildImage(io.Reader, string) error
	PushImage(string) error
}

// Deployer handles deploying new functions that a user submits by unpacking the function data from
// a users request and then using a ContainerDeployer to build and push a container image for later
// invocation
type Deployer struct {
	c  containerDeployer
	fs afero.Fs
}

// NewDeployer initializes a Deployer with a ContainerDeploy for building and pushing images
func NewDeployer(c containerDeployer, fs afero.Fs) Deployer {
	return Deployer{c, fs}
}

// Deploy unpacks the function code and then builds and pushes a container image
func (d Deployer) Deploy(data models.FnData) error {
	tarBuf, err := d.tarFromBase64Data(data.File)
	if err != nil {
		return err
	}

	tag := data.Name
	err = d.c.BuildImage(tarBuf, tag)
	if err != nil {
		return err
	}
	// d.c.PushImage()

	return nil
}

// tarFromBase64Data decodes a base64 string into a tar file as a buffer
func (d Deployer) tarFromBase64Data(base64Data string) (*bytes.Buffer, error) {
	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64Data))
	buf := &bytes.Buffer{}

	// write decoded data to a tar file as a buffer
	_, err := io.Copy(buf, decoder)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
