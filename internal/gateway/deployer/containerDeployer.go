package deployer

// ContainerDeployer contains all the methods required to turn function code into a container image
// and then sent to a container registry where it can be pulled for invocation
type ContainerDeployer interface {
	PushImage(string) error
	BuildImage(string, string) error
}
