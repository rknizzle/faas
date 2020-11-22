package manager

type ContainerRunner interface {
	PullImage(string) error
	RunContainer(string) error
}
