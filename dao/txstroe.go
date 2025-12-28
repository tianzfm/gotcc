package dao

import (
	"context"
)

type TXRecordDAO interface {
	GetTXRecords(ctx context.Context, opts ...expdao.QueryOption) ([]*expdao.TXRecordPO, error)
	CreateTXRecord(ctx context.Context, record *expdao.TXRecordPO) (uint, error)
	UpdateComponentStatus(ctx context.Context, id uint, componentID string, status string) error
	UpdateTXRecord(ctx context.Context, record *expdao.TXRecordPO) error
	LockAndDo(ctx context.Context, id uint, do func(ctx context.Context, dao *expdao.TXRecordDAO, record *expdao.TXRecordPO) error) error
}
