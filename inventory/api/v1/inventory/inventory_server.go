package inventory

import (
	"context"
	"time"

	"nms/api/v1/common"
	api "nms/api/v1/inventory"
	inv "nms/inventory/pkg/inventory"
	"nms/lib/repo"

	"google.golang.org/grpc"
)

type InventoryServer struct {
	api.UnimplementedInventoriesServer
}

func RegisterServices(s *grpc.Server) error {
	api.RegisterInventoriesServer(s, &InventoryServer{})
	api.RegisterSnapshotsServer(s, &SnpashotServer{})
	return nil
}

var TIMEOUT = time.Second * 10

// GetInventory API handler
func (s *InventoryServer) Get(ctx context.Context, req *api.GetInventoryRequest) (*api.GetInventoryResponse, error) {
	if req == nil {
		return nil, inv.ErrNullPoint
	}
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	r, err := repo.GetRepo(ctx)
	if err != nil {
		return nil, err
	}

	invs, err := inv.GetInventory(r.DB(), req.Id)
	if err != nil {
		return nil, err
	}

	if len(invs) != 1 {
		return nil, inv.ErrNotUnique(len(invs))
	}
	inv := invs[0].Unmarshal()

	return &api.GetInventoryResponse{
		Success:   true,
		Message:   "",
		Inventory: inv,
	}, nil
}

func (s *InventoryServer) List(ctx context.Context, req *api.ListInventoriesRequest) (*api.ListInventoriesResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	r, err := repo.GetRepo(ctx)
	if err != nil {
		return nil, err
	}
	invs, err := inv.ListInventories(ctx, r.DB(), req.Pagination)
	return &api.ListInventoriesResponse{
		Success: true,
		Message: "",
		Pagination: &common.PaginationResponse{
			Page:  req.Pagination.Page,
			Size:  req.Pagination.Size,
			Total: int32(len(invs)),
		},
		Inventories: invs,
	}, nil
}
