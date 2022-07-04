/*
 Inventory package has basic inventory type.
*/

package inventory

import (
	"context"
	"errors"
	"nms/api/v1/common"
	api "nms/api/v1/inventory"

	pg "nms/lib/pgutils"
	"nms/lib/repo"

	"github.com/google/uuid"
	lop "github.com/samber/lo/parallel"
	"google.golang.org/protobuf/types/known/structpb"
)

var ErrRoutineCancel = errors.New("")
var ErrNullPoint = errors.New("null point")

const TIME_LAYOUT = pg.TIME_LAYOUT

func parallel(funs ...func(context.Context, repo.IRepo) error) func(context.Context, repo.IRepo) error {
	return func(ctx context.Context, r repo.IRepo) error {
		out := make(chan error)
		for _, fun := range funs {
			go func(f func(context.Context, repo.IRepo) error) {
				out <- f(ctx, r)
			}(fun)
		}
		for i := 0; i < len(funs); i++ {
			select {
			case <-ctx.Done():
				return ErrRoutineCancel
			case err := <-out:
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
}

// Init initialize inventory init routines
// var Init = parallel(InitPostgresTable, ProcessScanResult)

type FirmwareInfo struct {
	Kernel string
	Ap     string
}

type Location struct {
	Path string
}

type Inventory struct {
	Name string
	// e.g.  devices.atop.switch.egh5910, be careful the category should
	// be the first element
	DeviceType          string
	Id                  string `pg:",pk"`
	Owner               string
	Model               string
	Location            Location
	IpAddress           string
	MacAddress          string
	HostName            string
	FirmwareInformation FirmwareInfo
	CreatedAt           string
	LastSeen            string
	LastMissing         string
	LastRecovered       string
	SupportProtocols    []string
	More                map[string]any
}

func (inv *Inventory) Unmarshal() *api.Inventory {
	more, err := structpb.NewStruct(inv.More)
	if err != nil {
		more, _ = structpb.NewStruct(map[string]any{})
	}
	return &api.Inventory{
		DeviceType: inv.DeviceType,
		Id:         inv.Id,
		Owner:      inv.Owner,
		Model:      inv.Model,
		Location: &api.Location{
			Path: inv.Location.Path,
		},
		IpAddress:  inv.IpAddress,
		MacAddress: inv.MacAddress,
		HostName:   inv.HostName,
		FirmwareInformation: &api.FirmwareInfo{
			Kernel: inv.FirmwareInformation.Kernel,
			Ap:     inv.FirmwareInformation.Ap,
		},
		CreatedAt:        inv.CreatedAt,
		LastSeen:         inv.LastSeen,
		LastMissing:      inv.LastMissing,
		LastRecovered:    inv.LastRecovered,
		SupportProtocols: inv.SupportProtocols,
		More:             more,
	}
}

func GetInventory(db pg.IClient, id string, opts ...pg.QueryExpr) ([]Inventory, error) {
	var invs []Inventory
	var opt = pg.QueryExpr{
		Expr:  "id = ?",
		Value: id,
	}
	if len(opts) > 0 {
		opt = opts[0]
	}
	err := db.Query(&invs, opt)
	return invs, err
}

func ListInventories(ctx context.Context, db pg.IClient, pg *common.Pagination) ([]*api.Inventory, error) {
	var invs []Inventory
	c, err := db.GetDB()
	if err != nil {
		return nil, err
	}
	err = c.Model(&invs).Limit(int(pg.Size)).Offset(int(pg.Page - 1)).Select()
	if err != nil {
		return nil, err
	}
	apiInvs := lop.Map(invs, func(i Inventory, _ int) *api.Inventory {
		return i.Unmarshal()
	})
	return apiInvs, nil
}

func InitPostgresTable(ctx context.Context, r repo.IRepo) error {
	db := r.DB()
	err := db.CreateTable(&Inventory{}, pg.CreateTableOpt{IfNotExists: true})
	if err != nil {
		return err
	}
	return nil
}

func newInventoryKey() string {
	return uuid.New().String()
}

func initInv(inv Inventory) Inventory {
	ret := inv
	ret.Id = newInventoryKey()
	return ret
}

var updateInv = pipe(updateInvTS, checkInvRecoveringTS)
var missingInv = pipe(missingInvTS)
var newInv = pipe(newInvTS, initInv)
