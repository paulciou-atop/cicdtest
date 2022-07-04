package inventory

import (
	"context"
	"errors"
	api "nms/api/v1/inventory"
)

type SnpashotServer struct {
	api.UnimplementedSnapshotsServer
}

// Create(context.Context, *CreateSnapshotRequest) (*CreateSnapshotResponse, error)
// // GET /v1/snapshots?pagination.page=1&pagination.size=20
// List(context.Context, *ListSnapshotRequest) (*ListSnapshotResponse, error)
// GetInventoris(context.Context, *GetSnapshotInventorisRequest) (*ListInventoriesResponse, error)
type newSnapshotRes = api.CreateSnapshotResponse
type newSnapshotReq = api.CreateSnapshotRequest

func (s *SnpashotServer) Create(ctx context.Context, req *newSnapshotReq) (*newSnapshotRes, error) {
	return &newSnapshotRes{}, errors.New("not implemnt")
}

type listSnapshotRes = api.ListSnapshotResponse
type listSnapshotReq = api.ListSnapshotRequest

func (s *SnpashotServer) List(ctx context.Context, req *listSnapshotReq) (*listSnapshotRes, error) {
	return &listSnapshotRes{}, errors.New("not implemnt")
}

type getSnapshotInvReq = api.GetSnapshotInventorisRequest

func (s *SnpashotServer) GetInventoris(ctx context.Context, req *getSnapshotInvReq) (*api.ListInventoriesResponse, error) {
	return &api.ListInventoriesResponse{}, errors.New("not implemnt")
}
