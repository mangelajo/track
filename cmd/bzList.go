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
	"github.com/mangelajo/track/pkg/bugzilla"
	"fmt"
)

// bzListCmd represents the bzList command
var bzListCmd = &cobra.Command{
	Use:   "bz-list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: bzList,
}

func init() {
	rootCmd.AddCommand(bzListCmd)
}

func bzList(cmd *cobra.Command, args []string) {

	fmt.Println("Login in...")
	client, _ := bugzilla.NewClient(bzURL, bzEmail,  bzPassword)

	fmt.Println("List query ..")
	buglist, _:= client.BugList(50,0)

	for _, bz := range buglist {
		fmt.Printf("%v\n", bz)

	}
	fmt.Println("bug info")

	bug,_ := client.BugInfo(1546996)
	fmt.Printf("%v\n", bug)

	fmt.Println("history:")
	hist,_ := client.BugHistory(1546996)

	fmt.Printf("%v\n", hist)

	fmt.Println("show_bug:")
	bi, _ := client.ShowBug(1546996)
	fmt.Printf("%v\n", bi)

	fmt.Println(bi.Cassigned_to.Content)
}


