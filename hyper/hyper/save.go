package hyper

import (
	"libvirt.org/libvirt-go"
)

func (h *Hyper) Save(desktop *libvirt.Domain, path string) error {
	return desktop.Save(path)
}
