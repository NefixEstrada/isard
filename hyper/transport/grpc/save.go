package grpc

import (
	"context"
	"errors"

	"github.com/isard-vdi/isard/common/pkg/grpc"
	"github.com/isard-vdi/isard/hyper/pkg/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"libvirt.org/libvirt-go"
)

// DesktopSave saves a desktop in the hypervisor
func (h *HyperServer) DesktopSave(ctx context.Context, req *proto.DesktopSaveRequest) (*proto.DesktopSaveResponse, error) {
	if err := grpc.Required(grpc.RequiredParams{
		"id":   req.Id,
		"path": req.Path,
	}); err != nil {
		return nil, err
	}

	desktop, err := h.hyper.Get(req.Id)
	if err != nil {
		var e libvirt.Error
		if errors.As(err, &e) {
			switch e.Code {
			case libvirt.ERR_NO_DOMAIN:
				return nil, status.Error(codes.NotFound, "desktop not found")
			}
		}

		return nil, status.Errorf(codes.Unknown, "get desktop: %v", err)
	}
	defer desktop.Free()

	if err := h.hyper.Save(desktop, req.Path); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &proto.DesktopSaveResponse{}, nil
}
