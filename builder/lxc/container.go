package lxc

import (
	"os/exec"
	"path"
)

const LXC_ROOT = "/var/lib/lxc"

type Container struct {
	Name     string
	Template string
}

func (c *Container) Root() string {
	return path.Join(LXC_ROOT, c.Name, "rootfs")
}

// lxc-create requires sudo permission to create a container at /var/lib/lxc
func (c *Container) Create() (err error) {
	return exec.Command("sudo", "lxc-create", "-n", c.Name, "-t",
		c.Template).Run()
}

func (c *Container) Start() error {
	return exec.Command("sudo", "lxc-start", "-n", c.Name, "-d").Run()
}

// 10.0.3.1 is hardcoded lxcbr0 ip
func (c *Container) Ip() (ip string, err error) {
	out, err := exec.Command("host", c.Name, "10.0.3.1").CombinedOutput()
	return string(out), err
}

func (c *Container) Destroy() error {
	return exec.Command("sudo", "lxc-destroy", "-n", c.Name, "-f").Run()
}

func (c *Container) WaitForState(state string) error {
	return exec.Command("sudo", "lxc-wait", "-n", c.Name, "-s", state).Run()
}
