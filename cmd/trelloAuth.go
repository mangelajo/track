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
	"github.com/adlio/trello"
	"github.com/mangelajo/track/pkg/storecache"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Receives, tests, and stores the trello auth token",
	Long: `Using the trello API requires an auth token which
can be obtained from the website here:

` + GetTrelloAuthURL() + `
Once obtained you can pass it to this command and it will be stored
in the track database.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			trelloAuthHelp()
		}

		client := trello.NewClient(trelloAppKey, args[0])
		member, err := client.GetMember("me", trello.Defaults())
		if err == nil && member != nil {
			fmt.Printf("hello %s, you're logged into Trello\n", member.FullName)
			storecache.StoreTrelloToken(args[0])
		} else {
			fmt.Printf("\nWe had an issue logging, in.\n\n")
			trelloAuthHelp()
		}
	},
}

func trelloAuthHelp() {
	fmt.Println("You need to pass the auth token:")
	fmt.Println()
	fmt.Println("\ttrack trello auth << TOKEN >>")
	fmt.Println("")
	fmt.Printf("The token can be retrieved from:\n%s\n", GetTrelloAuthURL())
	fmt.Println()
	os.Exit(1)
}

func init() {
	trelloCmd.AddCommand(authCmd)
}
