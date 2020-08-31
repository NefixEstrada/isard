package hyper_test

import (
	"errors"
	"testing"

	"libvirt.org/libvirt-go"

	"github.com/isard-vdi/isard/hyper/hyper"

	"github.com/stretchr/testify/assert"
)

func TestHyperNew(t *testing.T) {
	assert := assert.New(t)

	cases := map[string]struct {
		URI         string
		ExpectedErr string
	}{
		"should create the hypervisor correctly": {},
		"should return an error if there's an error connecting to the libvirt daemon": {
			URI: ":::://///",
			ExpectedErr: libvirt.Error{
				Code:    libvirt.ERR_INTERNAL_ERROR,
				Domain:  libvirt.ErrorDomain(45),
				Message: "internal error: Unable to parse URI :::://///",
			}.Error(),
			//ExpectedError: "connect to libvirt: virError(Code=1, Domain=45, Message='internal error: Unable to parse URI :::://///')",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			h, err := hyper.New(nil, tc.URI)
			if h != nil {
				defer h.Close()
			}

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
		})
	}
}
