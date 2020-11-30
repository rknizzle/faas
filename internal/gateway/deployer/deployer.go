package deployer

import (
	"github.com/rknizzle/faas/internal/models"
	"os"
)

// Deployer handles deploying new functions that a user submits by unpacking the function data from
// a users request and then using a ContainerDeployer to build and push a container image for later
// invocation
type Deployer interface {
	Deploy(models.FnData) error
}

type deployer struct {
	c ContainerDeployer
	e Extractor
}

// NewDeployer initializes a Deployer with a ContainerDeploy for building and pushing images
func NewDeployer(c ContainerDeployer, e Extractor) deployer {
	return deployer{c, e}
}

// Deploy unpacks the function code and then builds and pushes a container image
func (d deployer) Deploy(data models.FnData) error {
	dir, err := d.e.FnDirFromBase64Data(data.Name, data.File)
	if err != nil {
		return err
	}

	tag := data.Name
	d.c.BuildImage(dir, tag)
	//b.cBuilder.PushImage()

	// remove temporary directory used to build the image
	os.RemoveAll(dir)

	return nil
}
