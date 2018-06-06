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
)

// bzListCmd represents the bzList command
var bzListCmd = &cobra.Command{
	Use:   "bz-list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: bzList,
}

var myBugs bool

func init() {

	rootCmd.AddCommand(bzListCmd)
	bzListCmd.Flags().StringP("dfg", "d", "", "Openstack DFG")
	bzListCmd.Flags().StringP("squad", "s", "", "Openstack DFG Squad")
	bzListCmd.Flags().BoolVarP(&myBugs,"me", "m", false,"List only my bugs")

}

const RH_OPENSTACK_PRODUCT = "Red Hat OpenStack"
const CMP_NETWORKING_OVN = "python-networking-ovn"
var BZ_NEW_ASSIGNED = []string {"NEW", "ASSIGNED"}
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

	client := getClient()

	query := bugzilla.BugListQuery{
		Limit:          50,
		Offset:         0,
		Classification: CLS_REDHAT,
		Product:        RH_OPENSTACK_PRODUCT,
		BugStatus:      BZ_NEW_ASSIGNED,
		WhiteBoard:     getWhiteBoardQuery(),
	}

	if myBugs {
		query.AssignedTo = bzEmail
	}

	buglist, _:= client.BugList(&query)

	for _, bz := range buglist {
		fmt.Printf("%v\n", bz)
	}

	bugs := make(chan bugzilla.Bug, 64)
	bzChan := make(chan bugzilla.Cbug, 64)

	var wg sync.WaitGroup

	for i := 0 ; i < workers; i++ {
		wg.Add(1)
		go func() {
			for bz := range bugs {
				bi, err := client.ShowBug(bz.ID, bz.Changed.String())
				if bi == nil || err != nil {
					fmt.Printf("Error grabbing bug %d : %s", bz.ID, err)
				} else {
					bzChan <- *bi
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

	for _, bug := range buglist {
		bugs <- bug
	}
	close(bugs)


	for bi := range bzChan {
		bi.ShortSummary(bugzilla.USE_COLOR)
	}
}



