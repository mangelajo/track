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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"github.com/mangelajo/track/pkg/bugzilla"
	"github.com/mangelajo/track/pkg/storecache"
)

var BzEmail string
var BzPassword string
var BzURL string
var workers int
var preCacheHTML bool
var dropInteractiveShell bool
var summary bool
// bzCmd represents the bz command
var bzCmd = &cobra.Command{
	Use:   "bz",
	Short: "Bugzilla related commands",
	Long: ``,
}

func init() {
	rootCmd.AddCommand(bzCmd)

	bzCmd.PersistentFlags().StringP("bzurl", "b", "https://bugzilla.redhat.com", "Bugzilla URL")
	bzCmd.PersistentFlags().String("bzemail", "", "Bugzilla login email")
	bzCmd.PersistentFlags().StringP("bzpass", "k", "", "Bugzilla login password")
	bzCmd.PersistentFlags().BoolVarP(&preCacheHTML, "html", "x", false, "Pre-cache html for bz show command")
	bzCmd.PersistentFlags().BoolVar(&dropInteractiveShell, "shell", false, "Start an interactive shell once the command is done")
	bzCmd.PersistentFlags().BoolVarP(&summary, "summary", "u", false, "Show a summary of the bugs we retrieve")
}

func initBzConfig() {
	for _, k := range []string {"bzurl", "bzemail", "bzpass"} {
		viper.BindPFlag(k, bzCmd.PersistentFlags().Lookup(k))
	}

	for _, k := range []string {"dfg", "squad"} {
		viper.BindPFlag(k, bzListCmd.Flags().Lookup(k))
	}


	BzURL = viper.GetString("bzurl")
	BzPassword = viper.GetString("bzpass")
	BzEmail = viper.GetString("bzemail")

	if BzEmail == "" {
		fmt.Println("No bz url provided either in parameters or ~/.track.yaml file")
		exampleTrackYaml()
		os.Exit(1)
	}

}


func GetBzClient() *bugzilla.Client {
	client, err := bugzilla.NewClient(BzURL, BzEmail, BzPassword, storecache.GetBzAuth,
		storecache.StoreBzAuth)
	if err != nil || client == nil {
		fmt.Printf("Problem during login to bugzilla: %s\n", err)
		os.Exit(1)
	}
	return client
}