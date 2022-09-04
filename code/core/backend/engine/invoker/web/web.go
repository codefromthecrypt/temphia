package web

import (
	"github.com/gin-gonic/gin"
	"github.com/temphia/temphia/code/core/backend/engine/invoker"
	"github.com/temphia/temphia/code/core/backend/libx/lazydata"
	"github.com/temphia/temphia/code/core/backend/xtypes"
	"github.com/temphia/temphia/code/core/backend/xtypes/etypes/job"
	"github.com/temphia/temphia/code/core/backend/xtypes/service"
)

var _ job.Invoker = (*WebRequest)(nil)

type WebRequest struct {
	rctx   *gin.Context
	signer service.Signer
}

func NewWeb(ctx *gin.Context, signer service.Signer) *WebRequest {
	return &WebRequest{
		rctx:   ctx,
		signer: signer,
	}
}

type request struct {
	Name string `json:"name,omitempty"`
}

func (r *WebRequest) Handle(method string, data xtypes.LazyData) (xtypes.LazyData, error) {
	switch method {
	case "get_header":
		name, err := getTargetName(data)
		if err != nil {
			return nil, err
		}
		return lazydata.NewAnyData(r.rctx.GetHeader(name)), nil
	case "get_headers":
		return lazydata.NewAnyData(r.rctx.Request.Header), nil

	case "get_cookie":
		name, err := getTargetName(data)
		if err != nil {
			return nil, err
		}

		cookie, err := r.rctx.Cookie(name)
		if err != nil {
			return nil, err
		}

		return lazydata.NewAnyData(cookie), nil
	case "get_cookies":
		return lazydata.NewAnyData(r.rctx.Request.Cookies()), nil
	case "get_client_ip":
		return lazydata.NewAnyData(r.rctx.ClientIP()), nil
	default:
		return nil, invoker.ErrInvokerActionNotImplemented
	}

}

func (r *WebRequest) Name() string {
	return invoker.TypeWebAdmin
}

func (r *WebRequest) CurrentUser() *job.InvokeUser {
	sclaim, err := invoker.ParseClaim(
		r.signer,
		r.rctx,
	)
	if err != nil {
		return nil
	}

	return &job.InvokeUser{
		UserId:    sclaim.UserID,
		UserGroup: sclaim.UserGroup,
		SessionId: sclaim.SessionID,
	}

}

func getTargetName(data xtypes.LazyData) (string, error) {
	req := &request{}
	err := data.AsObject(req)
	if err != nil {
		return "", err
	}
	return req.Name, nil
}