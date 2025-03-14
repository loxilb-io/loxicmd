/*
 * Copyright (c) 2025 LoxiLB Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at:
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package set

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"loxicmd/pkg/api"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// NewLogLevelCmd represents the save command
func NewSetLogInCmd(restOptions *api.RESTOptions) *cobra.Command {
	SetOptions := SetOptions{}

	var loginCmd = &cobra.Command{

		Use:   "login",
		Short: "login and set token",

		Run: func(cmd *cobra.Command, args []string) {
			o := api.LoginModel{}
			if SetOptions.Provider == "" {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Enter ID: ")
				userID, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("ID error:", err)
					return
				}
				userID = strings.TrimSpace(userID)
				fmt.Print("Enter Password: ")
				bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
				if err != nil {
					fmt.Println("\nPassword error:", err)
					return
				}

				// Make loginModel
				if err := ReadSetLogInOptions(&o, userID, bytePassword); err != nil {
					fmt.Printf("Error: %s\n", err.Error())
					return
				} // API Call

				resp, err := LoginAPICall(restOptions, o)
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
					return
				}
				// Save token
				// save the token in the file tmp/token.json
				if resp.StatusCode == http.StatusOK {
					PrintAndSaveTokenResult(resp, *restOptions)
					return
				}
			} else if SetOptions.Provider == "google" {
				// Google Login at first
				fmt.Printf("This process is how to log in with Google oauth. When you log in to the GUI you can get a token (access token) and a refresh token. After that, it is a registration process.\n")
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Enter Access Token: ")
				AccessToken, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Access Token error:", err)
					return
				}
				AccessToken = strings.TrimSpace(AccessToken)
				fmt.Print("Enter Refresh Token: ")
				RefreshToken, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Refresh token error", err)
					return
				}
				RefreshToken = strings.TrimSpace(RefreshToken)
				fmt.Printf("AccessToken: %v\n", AccessToken)
				fmt.Printf("RefreshToken: %v\n", RefreshToken)
				tokenFilePath := "/tmp/loxilbtoken"
				err = os.WriteFile(tokenFilePath, []byte(AccessToken), 0644)
				if err != nil {
					fmt.Printf("Error: Failed to write token to file: (%s)\n", err.Error())
					return
				}
				refreshTokenFilePath := "/tmp/loxilbrefreshtoken"
				err = os.WriteFile(refreshTokenFilePath, []byte(RefreshToken), 0644)
				if err != nil {
					fmt.Printf("Error: Failed to write token to file: (%s)\n", err.Error())
					return
				}
				fmt.Println("Login Success")

			} else if SetOptions.Provider == "manual" {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Enter Access Token: ")
				AccessToken, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Access Token error:", err)
					return
				}
				AccessToken = strings.TrimSpace(AccessToken)
				tokenFilePath := "/tmp/loxilbtoken"
				err = os.WriteFile(tokenFilePath, []byte(AccessToken), 0644)
				if err != nil {
					fmt.Printf("Error: Failed to write token to file: (%s)\n", err.Error())
					return
				}
				fmt.Println("Login Success")

			} else {
				fmt.Println("Error: Invalid provider name")
				return
			}

		},
	}
	loginCmd.Flags().StringVarP(&SetOptions.Provider, "provider", "", "", "Define the provider name ex) google, manual")

	return loginCmd
}

func ReadSetLogInOptions(o *api.LoginModel, userID string, password []byte) error {
	o.Username = userID
	o.Password = string(password)
	return nil
}

func LoginAPICall(restOptions *api.RESTOptions, loginModel api.LoginModel) (*http.Response, error) {
	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}

	return client.Login().Create(ctx, loginModel)
}

func PrintAndSaveTokenResult(resp *http.Response, o api.RESTOptions) {
	Tokenresp := api.TokenModel{}
	resultByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to read HTTP response: (%s)\n", err.Error())
		return
	}

	if err := json.Unmarshal(resultByte, &Tokenresp); err != nil {
		fmt.Printf("Error: Failed to unmarshal HTTP response: (%s)\n", err.Error())
		return
	}

	if Tokenresp.Token == "" {
		fmt.Println("Error: Failed to get token Please check your ID or Password")
		return
	}

	// Save the token to /tmp/loxilbtoken
	tokenFilePath := "/tmp/loxilbtoken"
	err = os.WriteFile(tokenFilePath, []byte(Tokenresp.Token), 0644)
	if err != nil {
		fmt.Printf("Error: Failed to write token to file: (%s)\n", err.Error())
		return
	}

	// if json options enable, it print as a json format.
	if o.PrintOption == "json" {
		resultIndent, _ := json.MarshalIndent(Tokenresp, "", "    ")
		fmt.Println(string(resultIndent))
		return
	}
	fmt.Println("Login Success")
}

// NewLogLevelCmd represents the save command
func NewSetLogOutCmd(restOptions *api.RESTOptions) *cobra.Command {
	SetOptions := SetOptions{}
	var logoutCmd = &cobra.Command{

		Use:   "logout",
		Short: "logout and remove token",

		Run: func(cmd *cobra.Command, args []string) {
			if SetOptions.Provider == "" {
				resp, err := LogOutAPICall(restOptions)
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
					return
				}
				// Remove token
				// Remove the token in the file tmp/token.json
				if resp.StatusCode == http.StatusOK {
					PrintAndRemoveTokenResult()
					return
				}
			} else if SetOptions.Provider == "google" {
				PrintAndRemoveRefreshTokenResult()
				PrintAndRemoveTokenResult()
			} else if SetOptions.Provider == "manual" {
				PrintAndRemoveTokenResult()
			} else {
				fmt.Println("Error: Invalid provider name")
				return
			}

		},
	}

	logoutCmd.Flags().StringVarP(&SetOptions.Provider, "provider", "", "", "Define the provider name ex) google, manual")

	return logoutCmd
}

func LogOutAPICall(restOptions *api.RESTOptions) (*http.Response, error) {
	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}
	// Call the Logout api. It will delete the token file in the DB.
	// It is create function, but it is actually delete function.
	return client.Login().SetUrl("/auth/logout").Create(ctx, nil)
}

func PrintAndRemoveTokenResult() {
	// Delete the token to /tmp/loxilbtoken
	tokenFilePath := "/tmp/loxilbtoken"
	err := os.Remove(tokenFilePath)
	if err != nil {
		fmt.Printf("Error: Failed to remove token file: (%s)\n", err.Error())
		return
	}
	fmt.Println("Logout Success")
}

func PrintAndRemoveRefreshTokenResult() {
	// Delete the token to /tmp/loxilbtoken
	tokenFilePath := "/tmp/loxilbrefreshtoken"
	err := os.Remove(tokenFilePath)
	if err != nil {
		fmt.Printf("Error: Failed to remove refreshtoken file: (%s)\n", err.Error())
		return
	}
}

// NewLogLevelCmd represents the save command
func NewSetRefreshTokenCmd(restOptions *api.RESTOptions) *cobra.Command {
	SetOptions := SetOptions{}

	var loginCmd = &cobra.Command{

		Use:   "refresh",
		Short: "refresh token",

		Run: func(cmd *cobra.Command, args []string) {
			if SetOptions.Provider == "google" {
				// Get the token from the file /tmp/loxilbtoken
				tokenFilePath := "/tmp/loxilbtoken"
				tokenByte, err := os.ReadFile(tokenFilePath)
				if err != nil {
					fmt.Printf("Error: Failed to read token file: (%s)\n", err.Error())
					return
				}
				// Get the refresh token from the file /tmp/loxilbrefreshtoken
				refreshTokenFilePath := "/tmp/loxilbrefreshtoken"
				refreshTokenByte, err := os.ReadFile(refreshTokenFilePath)
				if err != nil {
					fmt.Printf("Error: Failed to read refreshtoken file: (%s)\n", err.Error())
					return
				}
				o := api.TokenModel{}
				o.Token = string(tokenByte)
				o.RefreshToken = string(refreshTokenByte)

				resp, err := RefreshTokenAPICall(restOptions, o)
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
					return
				}
				// Save token
				// save the token in the file tmp/token.json
				if resp.StatusCode == http.StatusOK {
					PrintAndSaveTokenResult(resp, *restOptions)
					return
				} else if resp.StatusCode == http.StatusForbidden {
					fmt.Println("Error: Failed to refresh token. Please login again")
					return
				} else {
					fmt.Printf("resp.StatusCode: %v\n", resp.StatusCode)
					fmt.Println("Error: Failed to refresh token")
					return
				}
			} else {
				fmt.Println("Error: Invalid provider name")
				return
			}

		},
	}
	loginCmd.Flags().StringVarP(&SetOptions.Provider, "provider", "", "", "Define the provider name ex)google, github, manual")

	return loginCmd
}

func RefreshTokenAPICall(restOptions *api.RESTOptions, TokenModel api.TokenModel) (*http.Response, error) {
	client := api.NewLoxiClient(restOptions)
	ctx := context.TODO()
	var cancel context.CancelFunc
	if restOptions.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.TODO(), time.Duration(restOptions.Timeout)*time.Second)
		defer cancel()
	}
	queryArgs := make(map[string]string)
	queryArgs["token"] = TokenModel.Token
	queryArgs["refreshtoken"] = TokenModel.RefreshToken
	return client.Login().SetUrl("oauth/google/token").Query(queryArgs).Get(ctx)
}
