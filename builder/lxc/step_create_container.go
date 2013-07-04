package lxc

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

type stepCreateContainer struct {
	container *Container
}

func (s *stepCreateContainer) Run(state map[string]interface{}) multistep.StepAction {
	config := state["config"].(config)
	ui := state["ui"].(packer.Ui)
	s.container = &Container{Name: config.Name, Template: config.Template}

	ui.Say("Launching a container with template: " + config.Template)
	if err := s.container.Create(); err != nil {
		err := fmt.Errorf("Error creating a container named %s with template %s: %s",
			config.Name, config.Template, err)
		state["error"] = err
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	state["container"] = s.container
	return multistep.ActionContinue
}

func (s *stepCreateContainer) Cleanup(state map[string]interface{}) {
	ui := state["ui"].(packer.Ui)
	if s.container == nil {
		return
	}

	ui.Say("Terminating the LXC container...")
	if err := s.container.Destroy(); err != nil {
		ui.Error(fmt.Sprintf("Error terminating container, may still be around: %s", err))
		return
	}

	s.container.WaitForState("STOPPED")
}
