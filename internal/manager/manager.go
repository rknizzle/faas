package manager

type Manager interface {
	PushImage(string) error
	PullImage(string) error
	BuildImage(string, string) error
	RunContainer(string) error
}
