package utils

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func OKResp(ctx context.Context, c *app.RequestContext, data any) {
	c.JSON(consts.StatusOK, utils.H{
		"data": data,
	})
}

func ErrResp(ctx context.Context, c *app.RequestContext, err error) {
	c.JSON(consts.StatusInternalServerError, utils.H{
		"error": err.Error(),
	})
}
