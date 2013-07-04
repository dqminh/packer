package lxc

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/builder/common"
	"github.com/mitchellh/packer/packer"
	"log"
	"time"
)

const BuilderId = "mitchellh.lxc"

type Builder struct {
	config config
	runner multistep.Runner
}

type config struct {
	Name              string `mapstruture:"name"`
	Template          string `mapstructure:"template"`
	SSHPrivateKeyPath string `mapstructure:"ssh_private_key_path"`
	SSHPublicKeyPath  string `mapstructure:"ssh_public_key_path"`
	SSHUsername       string `mapstructure:"ssh_username"`
	SSHPort           int    `mapstructure:"ssh_port"`
	SSHTimeout        time.Duration

	PackerDebug   bool   `mapstructure:"packer_debug"`
	RawSSHTimeout string `mapstructure:"ssh_timeout"`
}

func (b *Builder) Prepare(raws ...interface{}) error {
	var err error

	for _, raw := range raws {
		err := mapstructure.Decode(raw, &b.config)
		if err != nil {
			return err
		}
	}

	if b.config.SSHPort == 0 {
		b.config.SSHPort = 22
	}

	// Accumulate any errors
	errs := make([]error, 0)

	if b.config.Name == "" {
		errs = append(errs, errors.New("A name must be specified"))
	}

	if b.config.Template == "" {
		errs = append(errs, errors.New("A template must be specified"))
	}

	if b.config.SSHPrivateKeyPath == "" {
		errs = append(errs, errors.New("A ssh_private_key_path must be specified"))
	}

	if b.config.SSHPublicKeyPath == "" {
		errs = append(errs, errors.New("A ssh_public_key_path key path must be specified"))
	}

	if b.config.SSHUsername == "" {
		errs = append(errs, errors.New("A ssh_username must be specified"))
	}

	if b.config.RawSSHTimeout == "" {
		b.config.RawSSHTimeout = "1m"
	}

	b.config.SSHTimeout, err = time.ParseDuration(b.config.RawSSHTimeout)
	if err != nil {
		errs = append(errs, fmt.Errorf("Failed parsing ssh_timeout: %s", err))
	}

	if len(errs) > 0 {
		return &packer.MultiError{errs}
	}

	log.Printf("Config: %+v", b.config)
	return nil
}

func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	// Setup the state bag and initial state for the steps
	state := make(map[string]interface{})
	state["config"] = b.config
	state["hook"] = hook
	state["ui"] = ui

	// Build the steps
	steps := []multistep.Step{
		&stepCreateContainer{},
	}

	// Run!
	if b.config.PackerDebug {
		b.runner = &multistep.DebugRunner{
			Steps:   steps,
			PauseFn: common.MultistepDebugFn(ui),
		}
	} else {
		b.runner = &multistep.BasicRunner{Steps: steps}
	}

	b.runner.Run(state)

	// If there was an error, return that
	if rawErr, ok := state["error"]; ok {
		return nil, rawErr.(error)
	}

	// If there are no AMIs, then just return
	if _, ok := state["container"]; !ok {
		return nil, nil
	}

	// Build the artifact and return it
	artifact := &artifact{}
	return artifact, nil
}

func (b *Builder) Cancel() {
	if b.runner != nil {
		log.Println("Cancelling the step runner...")
		b.runner.Cancel()
	}
}
