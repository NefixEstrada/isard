package grpc

import (
	"context"
	"errors"

	"github.com/isard-vdi/isard/hyper/hyper"
	"github.com/isard-vdi/isard/hyper/pkg/proto"

	"github.com/isard-vdi/isard/common/pkg/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"libvirt.org/libvirt-go"
)

// DesktopValidate validates the XML definition
// TODO: Capabilities?
func (h *HyperServer) DesktopValidate(ctx context.Context, req *proto.DesktopValidateRequest) (*proto.DesktopValidateResponse, error) {
	if err := grpc.Required(grpc.RequiredParams{
		"xml": req.Xml,
	}); err != nil {
		return nil, err
	}

	desktop, err := h.hyper.Start(req.Xml, &hyper.StartOptions{Paused: true})
	if err != nil {
		var e libvirt.Error
		if errors.As(err, &e) {
			switch e.Code {
			case libvirt.ERR_XML_ERROR, libvirt.ERR_XML_DETAIL:
				return nil, status.Errorf(codes.InvalidArgument, "invalid XML: %s", e.Message)
			}
		}

		return nil, status.Error(codes.Unknown, err.Error())
	}
	defer desktop.Free()

	err = h.hyper.Stop(desktop)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "stop desktop: %v", err)
	}

	return &proto.DesktopValidateResponse{}, nil
}
