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
	"strings"

	"github.com/mangelajo/trello"
	"github.com/spf13/cobra"
)

// cardListCmd represents the cardList command
var cardListCmd = &cobra.Command{
	Use:   "cards",
	Short: "List cards available in board",
	Long:  ``,
	Run:   trelloListCards,
}

func init() {
	trelloCmd.AddCommand(cardListCmd)
	cardListCmd.Flags().BoolVarP(&myCards, "me", "m", false, "List only cards assigned to me")
	cardListCmd.Flags().StringVarP(&cardDFG, "dfg", "d", "", "List only cards with DFG custom field")
	cardListCmd.Flags().BoolVar(&cardNoDFG, "no-dfg", false, "List only cards with no DFG custom field")
	cardListCmd.Flags().StringVar(&cardList, "list", "", "List only cards on a specific List")
}

var myCards bool
var cardDFG string
var cardNoDFG bool
var cardList string

func trelloListCards(cmd *cobra.Command, args []string) {

	var listID *string = nil

	if len(args) != 1 {
		fmt.Println("This command needs at least a board ID")
		os.Exit(1)
	}

	trelloClient := GetTrelloClient()

	me, err := trelloClient.GetMember("me", trello.Defaults())
	checkError(err)

	board := FindBoard(trelloClient, args[0])

	cfields, err := board.GetCustomFields(trello.Defaults())
	checkError(err)

	cards, err := board.GetCards(map[string]string{
		"customFieldItems": "true",
		"fields":           "all",
	})
	checkError(err)

	lists, err := board.GetLists(trello.Defaults())
	checkError(err)

	listMap := map[string]string{}
	for _, list := range lists {
		listMap[list.ID] = list.Name
		if cardList != "" && strings.Contains(strings.ToUpper(list.Name),
			strings.ToUpper(cardList)) {
			listID = &list.ID
		}
	}

	if cardList != "" && listID == nil {
		fmt.Printf("Unknown list %s for board %s", cardList, args[0])
		os.Exit(1)
	}

	currentList := ""

	for _, card := range cards {

		if myCards && !meOnCard(card, me) {
			continue
		}

		fields := card.CustomFields(cfields)

		dfg, ok := fields["DFG"]

		if ok && cardNoDFG {
			continue
		}

		if cardDFG != "" {
			if !ok {
				continue
			}
			if !strings.Contains(strings.ToUpper(dfg.(string)), strings.ToUpper(cardDFG)) {
				continue
			}
		}

		if currentList != card.IDList {
			if listID == nil || card.IDList == *listID {
				fmt.Println("")
				fmt.Println(listMap[card.IDList])
			}
			currentList = card.IDList
		}

		if listID == nil || currentList == *listID {
			fmt.Printf(" - %s\t%v\n", card.ShortUrl, card.Name)
		}

	}
}

func meOnCard(card *trello.Card, me *trello.Member) bool {
	meOnCard := false
	for _, id := range card.IDMembers {
		if id == me.ID {
			meOnCard = true
			break
		}
	}
	return meOnCard
}
