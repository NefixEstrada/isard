package desktopbuilder_test

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/isard-vdi/isard/hyper/hyper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"libvirt.org/libvirt-go"
)

func TestViewerGet(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	cases := map[string]struct {
		PrepareDesktop func(h *hyper.Hyper) *libvirt.Domain
		ExpectedErr    string
	}{
		"save the desktop correctly": {
			PrepareDesktop: func(h *hyper.Hyper) *libvirt.Domain {
				desktop, err := h.Start(hyper.TestMinDesktopXML(t), &hyper.StartOptions{})
				require.NoError(err)

				return desktop
			},
		},
		"should return an error if there's an error suspending the desktop": {
			PrepareDesktop: func(h *hyper.Hyper) *libvirt.Domain {
				return &libvirt.Domain{}
			},
			ExpectedErr: "virError(Code=7, Domain=6, Message='invalid domain pointer in virDomainSave')",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			h, err := hyper.New(nil, hyper.TestLibvirtDriver(t))
			require.NoError(err)

			defer h.Close()

			desktop := tc.PrepareDesktop(h)
			if desktop != nil {
				defer desktop.Free()
			}

			dir, err := ioutil.TempDir("", "dumps")
			if err != nil {
				log.Fatal(err)
			}
			defer os.RemoveAll(dir)

			err = h.Save(desktop, filepath.Join(dir, "test.dump"))

			if tc.ExpectedErr == "" {
				assert.NoError(err)
			} else {
				assert.EqualError(err, tc.ExpectedErr)
			}
		})
	}
}
