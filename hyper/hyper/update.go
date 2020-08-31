package hyper

import (
	"libvirt.org/libvirt-go"
)

// UpdateOptions are a set of parameters that modify how the desktop is updated
type UpdateOptions struct {
	// Live specifies that the device shall be changed on the active domain instance only and is not added to the persisted domain configuration
	Live bool
	// Current specifies that the device change is made based on current domain state
	Current bool
	// Config specifies that the device shall be changed on the persisted domain configuration only.
	Config bool
}

// UpdateLive live migrates a running desktop to another hypervisor
func (h *Hyper) Update(desktop *libvirt.Domain, xml string, options *UpdateOptions) error {
	flags := libvirt.DOMAIN_DEVICE_MODIFY_FORCE
	if options.Live {
		flags = libvirt.DOMAIN_DEVICE_MODIFY_LIVE
	} else if options.Current {
		flags = libvirt.DOMAIN_DEVICE_MODIFY_CURRENT
	} else if options.Config {
		flags = libvirt.DOMAIN_DEVICE_MODIFY_CONFIG
	}

	if err := desktop.UpdateDeviceFlags(xml, flags); err != nil {
		return err
	}

	return nil
}
