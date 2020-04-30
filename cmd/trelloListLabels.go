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
	"fmt"
	"os"

	"github.com/mangelajo/trello"
	"github.com/spf13/cobra"
)

// cardListCmd represents the cardList command
var labelListCmd = &cobra.Command{
	Use:   "labels",
	Short: "List labels available in board",
	Long:  ``,
	Run:   trelloListLabels,
}

func init() {
	trelloCmd.AddCommand(labelListCmd)
}

func trelloListLabels(cmd *cobra.Command, args []string) {

	if len(args) != 1 {
		fmt.Println("This command needs at least a board ID")
		os.Exit(1)
	}

	trelloClient := GetTrelloClient()

	board := FindBoard(trelloClient, args[0])

	labels, err := board.GetLabels(trello.Defaults())
	checkError(err)

	fmt.Println("Labels:")

	for _, label := range labels {
		fmt.Printf(" - %s\t%s\t%s\n", label.ID, label.Color, label.Name)
	}
}
