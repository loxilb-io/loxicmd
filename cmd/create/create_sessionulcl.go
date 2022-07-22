package create

import (
	"context"
	"loxicmd/pkg/api"
	"net/http"
	"time"
)

func SessionUlClAPICall(restOptions *api.RESTOptions, ulclModel api.SessionUlClMod) (*http.Response, error) {
	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}

	return client.SessionUlCL().Create(ctx, ulclModel)
}
