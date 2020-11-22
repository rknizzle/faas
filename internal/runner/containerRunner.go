package runner

type ContainerRunner interface {
	PullImage(string) error
	RunContainer(string) error
}
