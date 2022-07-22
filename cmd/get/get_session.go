// Author: Inho gog <inhogog2@netlox.io>
package get

import (
	"context"
	"fmt"
	"io/ioutil"
	"loxicmd/pkg/api"
	"net/http"
	"os"
	"strings"
	"time"
)

func SessionAPICall(restOptions *api.RESTOptions) (*http.Response, error) {
	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}
	resp, err := client.Session().SetUrl("/config/session/all").Get(ctx)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return nil, err
	}
	return resp, nil
}

func Sessiondump(restOptions *api.RESTOptions) (string, error) {
	// File Open
	fileP := []string{"sessionconfig_", ".txt"}
	t := time.Now()
	file := strings.Join(fileP, t.Local().Format("2006-01-02_15:04:05"))
	f, err := os.Create(file)
	if err != nil {
		fmt.Printf("Can't create dump file\n")
		os.Exit(1)
	}
	defer f.Close()

	// API Call
	resp, err := SessionAPICall(restOptions)
	if err != nil {
		return "", err
	}
	resultByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
	}
	// Write
	f.Write(resultByte)

	return file, nil
}
