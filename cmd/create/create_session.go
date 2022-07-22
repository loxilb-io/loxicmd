package create

import (
	"context"
	"loxicmd/pkg/api"
	"net/http"
	"time"
)

func SessionAPICall(restOptions *api.RESTOptions, sessionModel api.SessionMod) (*http.Response, error) {
	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}

	return client.Session().Create(ctx, sessionModel)
}
