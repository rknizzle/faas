package builder

type ContainerDeployer interface {
	PushImage(string) error
	BuildImage(string, string) error
}
