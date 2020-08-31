package hyper

import (
	"os"
)

func (h *Hyper) Restore(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return os.ErrNotExist
		}
		return err
	}

	return h.conn.DomainRestore(path)
}
