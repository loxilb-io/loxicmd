/*
 * Copyright (c) 2022 NetLOX Inc
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
package dump

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

type ConfigFiles struct {
	IpConfigFile string
	LBConfigFile string
}

// applyCmd represents the save command
func ApplyCmd(cfgFiles *ConfigFiles) *cobra.Command {
	applyCmd := &cobra.Command{
	Use:   "apply",
	Short: "Apply configuration",
	Long:  `Reads and apply configuration from the text file`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd
		_ = args
		if len(cfgFiles.IpConfigFile) == 0 && len(cfgFiles.LBConfigFile) == 0 {
			fmt.Println("Provide valid filename")
			return
		}
		if len(cfgFiles.IpConfigFile) > 0 {
			ApplyIpConfig(cfgFiles.IpConfigFile)
			fmt.Printf("Configuration applied - %s\n", cfgFiles.IpConfigFile)
		}
		if len(cfgFiles.LBConfigFile) > 0 {
			fmt.Println("LB Configuration apply not implemented yet!")
		}
	},
	}
	return applyCmd
}

func ApplyIpConfig(file string) {
	// open file
	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	// remember to close the file at the end of the program
	defer f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		// do something with a line
		fmt.Printf("%s\n", scanner.Text())
		cmd := exec.Command("bash", "-c", scanner.Text())
		output, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("%v\n", string(output))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}
