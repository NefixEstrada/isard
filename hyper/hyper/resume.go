package hyper

import (
	"libvirt.org/libvirt-go"
)

func (h *Hyper) Resume(desktop *libvirt.Domain) error {
	return desktop.Resume()
}
