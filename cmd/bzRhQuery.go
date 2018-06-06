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
	"net/http"
	"strings"
	"os"
)


var bzRhQueryCmd = &cobra.Command{
	Use:   "bz-rh-query",
	Short: "Grab query parameters from https://url.corp.redhat.com/< name >",
	Long: `This command will grab the bugzilla query parameters from the
url redirector, and provide a list, it also allows caching of bugs for
bz-show <id> with the -x flag`,
	Run: bzRhQuery,
}

func init() {
	rootCmd.AddCommand(bzRhQueryCmd)
}

func findRHQuery(name string) string {
	url := "https://url.corp.redhat.com/" + name

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
	} }

	resp, err := client.Get(url)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 301 {
		location := resp.Header.Get("Location")
		if !strings.Contains(location, "https://bugzilla.redhat.com/buglist.cgi?") {
			return ""
		} else {
			return strings.Replace(location,
				"https://bugzilla.redhat.com/buglist.cgi?", "", 1)
		}
	}

	return ""
}

func bzRhQuery(cmd *cobra.Command, args []string) {

	if len(args)<1 {
		fmt.Println("We need at least one target URL, for example:")
		fmt.Println("")
		fmt.Println(" - network-dfg-untriaged")
		fmt.Println(" - network-dfg-untriaged-customers")
		fmt.Println(" - network-dfg-untriaged-ci")
		fmt.Println(" - network-dfg-triaged-error")
		fmt.Println("")
		os.Exit(1)
	}

	urlQuery := findRHQuery(args[0])

	if urlQuery == "" {
		fmt.Printf("No bugzilla query found for %s\n", args[0])
		os.Exit(1)
	}

	client := getClient()

	buglist, _ := client.BugList(&bugzilla.BugListQuery{CustomQuery: urlQuery})

	for _, bz := range buglist {
		fmt.Printf("%s\n", bz.String())
	}

	bzChan := grabBugzillasConcurrently(client, buglist)

	for bi := range bzChan {
		if !changedBugs || (changedBugs && !bi.Cached) {
			bi.Bug.ShortSummary(bugzilla.USE_COLOR)
		}
	}

	if preCacheHTML {
		fmt.Println("Pre caching HTML")
		grabBugzillasHTMLConcurrently(client, buglist)

	}

}



