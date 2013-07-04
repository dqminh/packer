package lxc

type artifact struct {
}

func (*artifact) BuilderId() string {
	return ""
}

func (*artifact) Files() []string {
	// We have no files
	return nil
}

func (a *artifact) Id() string {
	return ""
}

func (a *artifact) String() string {
	return ""
}

func (a *artifact) Destroy() error {
	return nil
}
