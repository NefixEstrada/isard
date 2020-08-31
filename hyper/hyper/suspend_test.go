package hyper_test

import (
	"errors"
	"testing"

	"github.com/isard-vdi/isard/hyper/hyper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"libvirt.org/libvirt-go"
)

func TestSuspend(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	cases := map[string]struct {
		PrepareDesktop  func(h *hyper.Hyper) *libvirt.Domain
		ExpectedErr     string
		ExpectedDesktop func(desktop *libvirt.Domain)
	}{
		"suspend the desktop correctly": {
			PrepareDesktop: func(h *hyper.Hyper) *libvirt.Domain {
				desktop, err := h.Start(hyper.TestMinDesktopXML(t), &hyper.StartOptions{})
				require.NoError(err)

				return desktop
			},
			ExpectedDesktop: func(desktop *libvirt.Domain) {
				state, _, err := desktop.GetState()

				assert.NoError(err)
				assert.Equal(libvirt.DOMAIN_PAUSED, state)
			},
		},
		"should return an error if there's an error suspending the desktop": {
			PrepareDesktop: func(h *hyper.Hyper) *libvirt.Domain {
				return &libvirt.Domain{}
			},
			ExpectedErr: libvirt.Error{
				Code:    libvirt.ERR_INVALID_DOMAIN,
				Domain:  libvirt.ErrorDomain(6),
				Message: "invalid domain pointer in virDomainSuspend",
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

			err = h.Suspend(desktop)

			if tc.ExpectedErr == "" {
				assert.NoError(err)
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
