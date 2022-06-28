/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

// whatCmd represents the what command
var whatCmd = &cobra.Command{
	Use:   "what",
	Short: "print 'what the hell'",
	Long:  `print "what the hell". there is no mean.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("what called")
		fmt.Printf("args: %v\n", args)

		hellFlag := cmd.Flag("hell")
		fmt.Printf("hellFlag: %v\n", hellFlag)
		fmt.Printf("flags: hell: Type = %s. value = %s\n", hellFlag.Value.Type(), hellFlag.Value.String())

		argsStr := strings.Join(args, "-")
		url := "http://192.168.20.253:8090/what/the/" + argsStr

		//resp, err := http.Get(url)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			fmt.Println("failed to create http request")
			return
		}

		if hellFlag.Value.String() != "" {
			fmt.Println("누가 hell 옵션을 넣었어?")
			q := req.URL.Query()
			q.Add("o", "wide")
			req.URL.RawQuery = q.Encode()
		}

		fmt.Printf("call to URL: %s\n", req.URL.RequestURI())

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("failed to call url(%s)\n", url)
			return
		}
		defer resp.Body.Close()

		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("faield to read response body")
			return
		}

		fmt.Printf("response: %s\n", result)
	},
}
