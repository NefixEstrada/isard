package grpc

import (
	"context"
	"os"

	"github.com/isard-vdi/isard/hyper/pkg/proto"

	"github.com/isard-vdi/isard/common/pkg/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DesktopResume resumes a suspended desktop in the hypervisor
func (h *HyperServer) DesktopRestore(ctx context.Context, req *proto.DesktopRestoreRequest) (*proto.DesktopRestoreResponse, error) {
	if err := grpc.Required(grpc.RequiredParams{
		"path": req.Path,
	}); err != nil {
		return nil, err
	}

	if err := h.hyper.Restore(req.Path); err != nil {
		if os.IsNotExist(err) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}

	if err := os.Remove(req.Path); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &proto.DesktopRestoreResponse{}, nil
}
