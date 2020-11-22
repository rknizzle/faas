package runner

// ContainerRunner contains all the methods required to handle a function invocation by pulling down
// a function container image from a remote registry and running the function code in the container
type ContainerRunner interface {
	PullImage(string) error
	RunContainer(string) error
}
