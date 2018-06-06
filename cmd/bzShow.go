// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"strconv"
	"github.com/mangelajo/track/pkg/show"
)


var bzShowCmd = &cobra.Command{
	Use:   "bz-show",
	Short: "Open cached HTML for bugzilla",
	Long: `This command will open a cached HTML, or grab it from bugzilla and
open it. It will use the command specified in the -openHtmlCommand parameter,
or in the ~/.track.yaml file`,
	Run: bzShow,
}

func init() {
	rootCmd.AddCommand(bzShowCmd)
}

func bzShow(cmd *cobra.Command, args []string) {

	if len(args) < 1 {
		fmt.Println("We need at least one bz")
		os.Exit(1)
	}

	bzid, err := strconv.Atoi(args[0])

	if err != nil {
		fmt.Println("Cannot parse bz id %s", args[0])
	}

	res := show.OpenBz(bzid, GetBzClient)
	os.Exit(res)
}
