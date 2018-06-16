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
	"github.com/adlio/trello"
	"github.com/mangelajo/track/pkg/storecache"
	"os"
)

// trelloCmd represents the trello command
var trelloCmd = &cobra.Command{
	Use:   "trello",
	Short: "Trello related commands",
	Long: ``,
}

func init() {
	rootCmd.AddCommand(trelloCmd)
}

func GetTrelloAuthURL() string {
	return fmt.Sprintf("https://trello.com/1/authorize?expiration=never&scope=read,write,account&" +
		"response_type=token&name=Track&key=%s", trelloAppKey)
}

func GetTrelloClient() *trello.Client {
	token := storecache.GetTrelloToken()
	if token == nil || *token == "" {
		fmt.Println("You need a token from trello, please visit: ")
		fmt.Println(GetTrelloAuthURL())
		fmt.Println("")
		fmt.Println("Then run: track trello auth <<TOKEN>>")
		os.Exit(1)
	}
	return nil
}
