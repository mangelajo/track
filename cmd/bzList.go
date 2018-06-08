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
var productStr string
var assignedTo string
var statusStr string
var flagOn string


func init() {

	rootCmd.AddCommand(bzListCmd)
	bzListCmd.Flags().StringVarP(&flagOn, "flags-on", "f", "", "List bugs with flags on somebody (you can use 'me')")
	bzListCmd.Flags().StringP("dfg", "d", "", "Openstack DFG")
	bzListCmd.Flags().StringP("squad", "", "", "Openstack DFG Squad")
	bzListCmd.Flags().StringVarP(&statusStr, "status", "s", "NEW,ASSIGNED,POST,MODIFIED,ON_DEV,ON_QA,VERIFIED,RELEASE_PENDING", "Status list separated by commas")
	bzListCmd.Flags().StringVarP(&componentStr, "component", "c", "", "Component")
	bzListCmd.Flags().StringVarP(&productStr, "product", "p", "", "Product")
	bzListCmd.Flags().StringVarP(&assignedTo, "assignee", "a", "", "Filter by assignee (you can use 'me'")
	bzListCmd.Flags().BoolVarP(&myBugs,"me", "m", false,"List only bugs assigned to me")
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

	statusSelectors := strings.Split(statusStr, ",")

	query := bugzilla.BugListQuery{
		Limit:          50,
		Offset:         0,
		Classification: CLS_REDHAT,
		Product:        productStr,
		Component: 		componentStr,
		BugStatus:      statusSelectors,
		WhiteBoard:     getWhiteBoardQuery(),
		AssignedTo:		assignedTo,
		FlagRequestee:  flagOn,
	}


	if myBugs || query.AssignedTo == "me" {
		query.AssignedTo = BzEmail
	}

	if query.FlagRequestee == "me" {
		query.FlagRequestee = BzEmail
	}

	client := GetBzClient()

	buglist, _:= client.BugList(&query)

	listBugs(buglist, client)
}

func listBugs(buglist []bugzilla.Bug, client *bugzilla.Client) {
	for _, bz := range buglist {
		fmt.Printf("%s\n", bz.String())
	}
	bzChan := grabBugzillasConcurrently(client, buglist)
	var bugs []bugzilla.Cbug
	fmt.Print("Grabbing bug details:")
	for bi := range bzChan {
		if !changedBugs || (changedBugs && !bi.Cached) {
			if !dropInteractiveShell {
				bi.Bug.ShortSummary(bugzilla.USE_COLOR)
			} else if !bi.Cached {
				fmt.Printf(" bz#%d", bi.Bug.Cbug_id.Number)
			}
			bugs = append(bugs, bi.Bug)
		}
	}
	if dropInteractiveShell {
		fmt.Println(" done.")
	}
	if preCacheHTML {
		fmt.Print("Pre caching HTML:")
		grabBugzillasHTMLConcurrently(client, buglist)
		fmt.Println(" done.")
	}
	fmt.Println("")
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
					fmt.Printf(" bz#%d", bz.ID)
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