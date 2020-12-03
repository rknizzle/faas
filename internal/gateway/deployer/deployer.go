package deployer

import (
	"github.com/rknizzle/faas/internal/models"
	"os"
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
