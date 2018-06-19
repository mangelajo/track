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
	"github.com/mangelajo/trello"
	"fmt"
	"os"
)

// cardListCmd represents the cardList command
var boardListCmd = &cobra.Command{
	Use:   "list-boards",
	Short: "List boards available to user",
	Long: ``,
	Run: trelloListBoards,
}

func init() {
	trelloCmd.AddCommand(boardListCmd)
}


func trelloListBoards(cmd *cobra.Command, args []string) {

	trelloClient := GetTrelloClient()

	member, _ := trelloClient.GetMember("me", trello.Defaults())
	boards, err := member.GetBoards(trello.Defaults())

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Boards:")

	for _, b := range boards {
		if b.Closed {
			continue
		}
		fmt.Printf(" - %s\t%s\t%s\n", b.ID, b.ShortUrl, b.Name)
	}
}
