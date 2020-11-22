package builder

type ContainerBuilder interface {
	PushImage(string) error
	BuildImage(string, string) error
}
