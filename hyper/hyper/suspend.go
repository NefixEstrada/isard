package hyper

import (
	"libvirt.org/libvirt-go"
)

func (h *Hyper) Suspend(desktop *libvirt.Domain) error {
	return desktop.Suspend()
}
