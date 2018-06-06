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
	"sync"
	"github.com/spf13/viper"
	"strings"
	"github.com/mangelajo/track/pkg/shell"
)

// bzListCmd represents the bzList command
var bzListCmd = &cobra.Command{
	Use:   "bz-list",
	Short: "List bugzillas based on parameters and configuration",
	Long: `This command will list and retrieve details for bugzillas
based on configuration and query.`,
	Run: bzList,
}

var myBugs bool
var changedBugs bool
var componentStr string
var statusStr string


func init() {

	rootCmd.AddCommand(bzListCmd)
	bzListCmd.Flags().StringP("dfg", "d", "", "Openstack DFG")
	bzListCmd.Flags().StringP("squad", "", "", "Openstack DFG Squad")
	bzListCmd.Flags().StringVarP(&statusStr, "status", "s", "NEW,ASSIGNED", "Status list separated by commas")
	bzListCmd.Flags().StringVarP(&componentStr, "component", "c", "", "Component")
	bzListCmd.Flags().BoolVarP(&myBugs,"me", "m", false,"List only my bugs")
	bzListCmd.Flags().BoolVarP(&changedBugs,"changed", "", false,"Show bugs changed since last run")

}

const RH_OPENSTACK_PRODUCT = "Red Hat OpenStack"
const CMP_NETWORKING_OVN = "python-networking-ovn"
const CLS_REDHAT = "Red Hat"

func getWhiteBoardQuery() string {
	whiteBoardQuery := ""

	if viper.GetString("dfg") != "" {
		whiteBoardQuery += fmt.Sprintf("DFG:%s", viper.GetString("dfg"))
	}

	if viper.GetString("squad") != "" {
		if viper.GetString("dfg") != "" {
			whiteBoardQuery += " "
		}
		whiteBoardQuery += fmt.Sprintf("Squad:%s", viper.GetString("squad"))
	}
	return whiteBoardQuery
}

func bzList(cmd *cobra.Command, args []string) {

	client := GetBzClient()

	statusSelectors := strings.Split(statusStr, ",")

	query := bugzilla.BugListQuery{
		Limit:          50,
		Offset:         0,
		Classification: CLS_REDHAT,
		Product:        RH_OPENSTACK_PRODUCT,
		Component: 		componentStr,
		BugStatus:      statusSelectors,
		WhiteBoard:     getWhiteBoardQuery(),
	}

	if myBugs {
		query.AssignedTo = BzEmail
	}

	buglist, _:= client.BugList(&query)

	for _, bz := range buglist {
		fmt.Printf("%s\n", bz.String())
	}

	bzChan := grabBugzillasConcurrently(client, buglist)

	var bugs []bugzilla.Cbug

	for bi := range bzChan {
		if !changedBugs || (changedBugs && !bi.Cached) {
			bi.Bug.ShortSummary(bugzilla.USE_COLOR)
		}
		bugs = append(bugs, bi.Bug)
	}

	if preCacheHTML {
		fmt.Println("Pre caching HTML")
		grabBugzillasHTMLConcurrently(client, buglist)

	}

	if dropInteractiveShell {
		shell.Shell(&bugs, GetBzClient)
	}



}

type BugzillaResponse struct {
	Cached bool
	Bug bugzilla.Cbug
}

func grabBugzillasConcurrently(client *bugzilla.Client, buglist []bugzilla.Bug) chan BugzillaResponse {

	bugs := make(chan bugzilla.Bug, 64)
	bzChan := make(chan BugzillaResponse, 64)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			for bz := range bugs {
				bi, cached, err := client.ShowBug(bz.ID, bz.Changed.String())
				if bi == nil || err != nil {
					fmt.Printf("Error grabbing bug %d : %s", bz.ID, err)
				} else {
					bzChan <- BugzillaResponse{
						Cached: cached,
						Bug:    *bi,
					}
				}
			}
			wg.Done()
		}()
	}
	// when all workers have finished we close the output channel
	go func() {
		wg.Wait()
		close(bzChan)
	}()

	// Grab all bugs from the list and push them into the channel
	for _, bug := range buglist {
		bugs <- bug
	}
	close(bugs)
	return bzChan
}


func grabBugzillasHTMLConcurrently(client *bugzilla.Client, buglist []bugzilla.Bug) {

	bugs := make(chan bugzilla.Bug, 64)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			for bz := range bugs {
				_, cached, err := client.ShowBugHTML(bz.ID, bz.Changed.String())

				if err != nil {
					fmt.Printf("Error grabbing %d : %s\n", bz.ID, err)
				} else if !cached {
					fmt.Printf(" - bz#%d cached\n", bz.ID)
				}

			}
			wg.Done()
		}()
	}

	// Grab all bugs from the list and push them into the channel
	for _, bug := range buglist {
		bugs <- bug
	}
	close(bugs)
	wg.Wait()
}