package node

import (
	"context"

	"github.com/cxd-scu/bbot/biz/model"
)

type CommonNode interface {
	Name() string
	Run(ctx context.Context, req *model.IOParam) (*model.IOParam, error)
}
