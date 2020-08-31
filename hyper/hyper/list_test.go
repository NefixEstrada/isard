package hyper_test

import (
	"errors"
	"testing"

	"github.com/isard-vdi/isard/hyper/hyper"
	"libvirt.org/libvirt-go"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	cases := map[string]struct {
		PrepareTest          func(h *hyper.Hyper)
		ExpectedErr          string
		ExpectedDesktopNames []string
	}{
		"should list the desktops correctly": {
			ExpectedDesktopNames: []string{"test"},
		},
		"should return a libvirt error message if there's an error listing the desktops": {
			PrepareTest: func(h *hyper.Hyper) {
				h.Close()
			},
			ExpectedDesktopNames: []string{},
			ExpectedErr: libvirt.Error{
				Code:    libvirt.ERR_INVALID_CONN,
				Domain:  libvirt.ErrorDomain(20),
				Message: "invalid connection pointer in virConnectListAllDomains",
			}.Error(),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			h, err := hyper.New(nil, hyper.TestLibvirtDriver(t))
			require.NoError(err)

			defer h.Close()

			if tc.PrepareTest != nil {
				tc.PrepareTest(h)
			}

			desktops, err := h.List()

			if tc.ExpectedErr == "" {
				assert.NoError(err)
			} else {
				var e libvirt.Error
				if errors.As(err, &e) {
					assert.Equal(tc.ExpectedErr, e.Error())
				} else {
					assert.EqualError(err, tc.ExpectedErr)
				}
			}

			names := []string{}
			for _, desktop := range desktops {
				name, err := desktop.GetName()
				assert.NoError(err)
				names = append(names, name)
			}

			assert.Equal(tc.ExpectedDesktopNames, names)
		})
	}
}
