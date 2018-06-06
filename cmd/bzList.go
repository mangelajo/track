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

func init() {
	rootCmd.AddCommand(bzListCmd)
}

const RH_OPENSTACK_PRODUCT = "Red Hat OpenStack"
const CMP_NETWORKING_OVN = "python-networking-ovn"
var BZ_NEW_ASSIGNED = []string {"NEW", "ASSIGNED"}
const CLS_REDHAT = "Red Hat"

func bzList(cmd *cobra.Command, args []string) {

	if bzEmail == "" {
		panic(fmt.Errorf("No email address provided either in parameters or ~/.track.yaml file"))
	}
	if bzPassword == "" {
		panic(fmt.Errorf("No bz password provided either in parameters or ~/.track.yaml file"))
	}

	client, err := bugzilla.NewClient(bzURL, bzEmail,  bzPassword)

	if err != nil || client == nil {
		panic(fmt.Errorf("Problem during login to bugzilla: %s", err))
	}

	query := bugzilla.BugListQuery{
		Limit:          50,
		Offset:         0,
		Classification: CLS_REDHAT,
		Product:        RH_OPENSTACK_PRODUCT,
		BugStatus:      BZ_NEW_ASSIGNED,
		WhiteBoard:     whiteBoardQuery,
	}

	fmt.Println(whiteBoardQuery)

	if myBugs {
		query.AssignedTo = bzEmail
	}

	buglist, _:= client.BugList(&query)

	for _, bz := range buglist {
		fmt.Printf("%v\n", bz)

	}
	//fmt.Println("bug info")

	//bug,_ := client.BugInfo(1546996)
	//hist,_ := client.BugHistory(1546996)

	//fmt.Printf("%v\n", hist)

	fmt.Println("show_bug:")

	for _, bug := range buglist {
		bi, err := client.ShowBug(bug.ID, bug.Changed.String())
		if bi == nil || err != nil {
			fmt.Printf("Error grabbing bug %d : %s", bug.ID, err)
		} else {

			fmt.Printf("\nBZ %d (%8s) %s\n", bi.Cbug_id.Number, bi.Cbug_status.Content, bi.Cshort_desc.Content)

			if bi.Cassigned_to != nil {
				fmt.Printf("  Assigned to: %s\n",bi.Cassigned_to.Content)
			}

			for _, x:= range bi.Cexternal_bugs {

				fmt.Printf("  ext bug on %s : %s\n", x.Attrname, x.Content)
			}
		}
	}



}


