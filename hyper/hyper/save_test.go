package hyper_test

import (
	"errors"
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

func TestSave(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	cases := map[string]struct {
		PrepareDesktop func(h *hyper.Hyper) *libvirt.Domain
		ExpectedErr    string
		AfterTest      func(h *hyper.Hyper, desktop_name string, path string)
	}{
		"save the desktop correctly": {
			PrepareDesktop: func(h *hyper.Hyper) *libvirt.Domain {
				desktop, err := h.Start(hyper.TestMinDesktopXML(t), &hyper.StartOptions{})
				require.NoError(err)

				return desktop
			},
			AfterTest: func(h *hyper.Hyper, desktop_name string, path string) {
				err := h.Restore(path)
				assert.NoError(err)

				desktop, _ := h.Get(desktop_name)
				state, _, err := desktop.GetState()
				assert.NoError(err)

				assert.Equal(libvirt.DOMAIN_RUNNING, state)
			},
		},
		"should return an error if there's an error saving the desktop": {
			PrepareDesktop: func(h *hyper.Hyper) *libvirt.Domain {
				return &libvirt.Domain{}
			},
			ExpectedErr: libvirt.Error{
				Code:    libvirt.ERR_INVALID_DOMAIN,
				Domain:  libvirt.ErrorDomain(6),
				Message: "invalid domain pointer in virDomainSave",
			}.Error(),
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

			dir, err := ioutil.TempDir("", "isard-test-restore")
			if err != nil {
				log.Fatal(err)
			}
			defer os.RemoveAll(dir)

			err = h.Save(desktop, filepath.Join(dir, "test.dump"))

			if tc.ExpectedErr == "" {
				assert.NoError(err)

				desktop_name, err := desktop.GetName()
				assert.NoError(err)

				tc.AfterTest(h, desktop_name, filepath.Join(dir, "test.dump"))
			} else {
				var e libvirt.Error
				if errors.As(err, &e) {
					assert.Equal(tc.ExpectedErr, e.Error())
				} else {
					assert.EqualError(err, tc.ExpectedErr)
				}
			}

		})
	}
}
