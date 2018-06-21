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
	"strings"
	"log"
)

// cardListCmd represents the cardList command
var moveCardsCmd = &cobra.Command{
	Use:   "move",
	Short: "Move cards from board A to board B",
	Long: ``,
	Run: trelloMoveCards,
}

func init() {
	trelloCmd.AddCommand(moveCardsCmd)
	moveCardsCmd.Flags().StringVar(&labelFilter, "label-contains", "", "Case insensitive filter for card labels")
	moveCardsCmd.Flags().BoolVar(&forceMove, "force", false,"Move each card without asking one by one")
}

var labelFilter string
var forceMove bool

func trelloMoveCards(cmd *cobra.Command, args []string) {

	if len(args)<2 {
		fmt.Println("This command requires two arguments, origin board, destination board")
		os.Exit(1)
	}

	if labelFilter == "" {
		fmt.Println("We require at least one keyword to look for in --label-contains")
		os.Exit(1)
	}

	trelloClient := GetTrelloClient()

	// Grab necessary data

	srcBoard:= FindBoard(trelloClient, args[0])

	srcLists, err := srcBoard.GetLists(trello.Defaults())
	checkError(err)

	srcCards, err := srcBoard.GetCards(trello.Defaults())
	checkError(err)

	dstBoard:= FindBoard(trelloClient, args[1])

	dstCards, err := dstBoard.GetCards(trello.Defaults())
	checkError(err)

	dstLists, err := dstBoard.GetLists(trello.Defaults())
	checkError(err)


	// Build maps for quick lookups

	dstCardNameToCardMap := make(map[string]*trello.Card)
	for _, card := range dstCards {
		dstCardNameToCardMap[card.Name] = card
	}

	srcListIDtoNameMap := make(map[string]*trello.List)
	for _, list := range srcLists {
		srcListIDtoNameMap[list.ID] = list
	}

	dstListMap := make(map[string]string)

	for _, list := range dstLists {

		dstListMap[list.Name] = list.ID
	}


	// Start moving


	fmt.Println("")
	fmt.Printf("Moving srcCards: %s -> %s, filtering labels by: '%s'\n", srcBoard.Name, dstBoard.Name, labelFilter)

	for _, card := range srcCards {

		if labelFilter != "" {
			labelFound := false
			for _, label := range card.Labels {
				if strings.Contains(strings.ToUpper(label.Name), strings.ToUpper(labelFilter)) {
					labelFound = true
				}
			}
			if !labelFound {
				continue
			}

		}

		dstCard, ok := dstCardNameToCardMap[card.Name]
		if ok {
			fmt.Printf("Skipping: dst board already has card [%s]\n", card.Name)
			fmt.Printf("- src url: %s\n", card.ShortUrl)
			fmt.Printf("- dst url: %s\n", dstCard.ShortUrl)
			continue
		}

		listName := srcListIDtoNameMap[card.IDList].Name

		dstListID, ok := dstListMap[listName]

		if !ok {
			fmt.Printf("Skipping: dst board does not have list [%s], card [%s]\n", listName,
				card.Name)
			continue
		}

		fmt.Printf("\n%s : [%s]\n", card.ShortUrl, card.Name)
		confirmed := askForConfirmation("Do you want to move this card?")
		if !confirmed {
			continue
		}
		newCard, err := card.CopyToList(dstListID, map[string]string{"keepFromSource": "all"})
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}

		// Link the new card, and the old card to each others via comments
		newCard.AddComment(fmt.Sprintf("Moved from original card: %s", card.ShortUrl),
			trello.Defaults())
		card.AddComment(fmt.Sprintf("Card moved to %s", newCard.ShortUrl), trello.Defaults())

		// Then close the old card
		card.Closed = true
		upd_err := card.Update(map[string]string{"closed": "true"})
		if upd_err != nil {
			fmt.Print(err)
			os.Exit(1)
		}

		fmt.Printf("moved %s -> %s [%s]\n", card.ShortUrl, newCard.ShortUrl, card.Name)
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: %s\n", err)
		os.Exit(1)
	}
}

func askForConfirmation(question string) bool {
	var response string

	if question != "" {
		fmt.Print(question)
	}
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"Y", "YES"}
	nokayResponses := []string{"N", "NO", "S", "SKIP"}
	if isOneOf(okayResponses, response) {
		return true
	} else if isOneOf(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes, no or skip and then press enter:")
		return askForConfirmation("")
	}
}

func isOneOf(responses []string, response string) bool {
	for _, resp := range responses {
		if strings.ToUpper(resp) == strings.ToUpper(response) {
			return true
		}
	}
	return false
}