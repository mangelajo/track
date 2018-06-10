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
	"strings"
	"os"
	"github.com/spf13/viper"
)


var bzQueryCmd = &cobra.Command{
	Use:   "bz-query",
	Short: "Grab query parameters from your config",
	Long: `This command will grab the bugzilla query parameters from your
configuration under the queries key, for example:
queries:
    query1: https://bugzilla.redhat.com/buglist.cgi?bug_status=NEW&bug_status=ASSIGNED&classification=Red%20Hat&component=python-networking-ovn&component=openstack-neutron&list_id=8959823&product=Red%20Hat%20OpenStack&query_format=advanced
    query2: https://bugzilla.redhat.com/buglist.cgi?bug_status=NEW&classification=Red%20Hat&keywords=Triaged%2C%20&keywords_type=nowords&list_id=8959824&product=Red%20Hat%20OpenStack&query_format=advanced`,
	Run: bzQuery,
}

func init() {
	rootCmd.AddCommand(bzQueryCmd)
}

func findQuery(name string) string {

	if !viper.InConfig("queries") {
		fmt.Println("No queries config in your .track.yaml file")
		os.Exit(1)
	}

	queries := viper.GetStringMapString("queries")
	location, ok := queries[name]

	if !ok {
		fmt.Printf("Query %s unknown\n", name)
		os.Exit(1)
	}

	toReplace := viper.Get("bzurl").(string) + "/buglist.cgi?"
	if strings.Index(location, toReplace) == -1 {
		fmt.Printf("Your query %s\nis not under the %s url\n", location, toReplace)
		os.Exit(1)
	}
	return strings.Replace(location,
				       toReplace, "", 1)


}

func bzQuery(cmd *cobra.Command, args []string) {

	if len(args)<1 {
		fmt.Println("We need at least one target query.")
		os.Exit(1)
	}

	urlQuery := findQuery(args[0])

	client := GetBzClient()

	bugList, _ := client.BugList(&bugzilla.BugListQuery{CustomQuery: urlQuery})

	listBugs(bugList, client)

}



