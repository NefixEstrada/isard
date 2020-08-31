package hyper_test

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/isard-vdi/isard/hyper/hyper"
	"libvirt.org/libvirt-go"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRestore(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	cases := map[string]struct {
		PrepareTest     func(h *hyper.Hyper) (string, string)
		ExpectedErr     string
		ExpectedDesktop func(desktop *libvirt.Domain)
	}{
		"restore the desktop correctly": {
			PrepareTest: func(h *hyper.Hyper) (string, string) {
				desktop, err := h.Start(hyper.TestMinDesktopXML(t), &hyper.StartOptions{})
				require.NoError(err)

				desktop_name, err := desktop.GetName()
				require.NoError(err)

				dir, err := ioutil.TempDir("", "isard-test-restore")
				if err != nil {
					log.Fatal(err)
				}

				path := filepath.Join(dir, "test.dump")
				err = h.Save(desktop, path)
				require.NoError(err)

				return desktop_name, path
			},
			ExpectedDesktop: func(desktop *libvirt.Domain) {
				state, _, err := desktop.GetState()

				assert.NoError(err)
				assert.Equal(libvirt.DOMAIN_RUNNING, state)
			},
		},
		"should return an error as the file is invalid": {
			PrepareTest: func(h *hyper.Hyper) (string, string) {
				dir, err := ioutil.TempDir("", "isard-test-restore")
				if err != nil {
					log.Fatal(err)
				}

				path := filepath.Join(dir, "test.dump")

				file, err := os.Create(path)
				if err != nil {
					log.Fatalf("failed creating file: %s", err)
				}
				defer file.Close()

				_, err = file.WriteString("This will be invalid file content.")

				if err != nil {
					log.Fatalf("failed writing to file: %s", err)
				}

				return "", path
			},
			ExpectedErr: libvirt.Error{
				Code:    libvirt.ERR_INTERNAL_ERROR,
				Domain:  libvirt.ErrorDomain(12),
				Message: "internal error: mismatched header magic",
			}.Error(),
			ExpectedDesktop: func(desktop *libvirt.Domain) {
				assert.Nil(desktop)
			},
		},
		"should return an error if the path is incorrect or file missing": {
			PrepareTest: func(h *hyper.Hyper) (string, string) {
				return "", ""
			},
			ExpectedErr: os.ErrNotExist.Error(),
			ExpectedDesktop: func(desktop *libvirt.Domain) {
				assert.Nil(desktop)
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			h, err := hyper.New(nil, hyper.TestLibvirtDriver(t))
			require.NoError(err)

			defer h.Close()

			desktop_name, path := tc.PrepareTest(h)
			defer os.RemoveAll(path)

			err = h.Restore(path)

			if tc.ExpectedErr == "" {
				assert.NoError(err)
				desktop, _ := h.Get(desktop_name)
				tc.ExpectedDesktop(desktop)
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
