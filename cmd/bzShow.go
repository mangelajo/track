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
	"os"
	"strconv"
	"github.com/mangelajo/track/pkg/storecache"
	"strings"
	"github.com/spf13/viper"
	"os/exec"
)


var bzShowCmd = &cobra.Command{
	Use:   "bz-show",
	Short: "Grab cached HTML for bugzilla",
	Long: ``,
	Run: bzShow,
}

func init() {
	rootCmd.AddCommand(bzShowCmd)
}

func bzShow(cmd *cobra.Command, args []string) {

	if len(args)<1 {
		fmt.Println("We need at least one bz")
		os.Exit(1)
	}

	bzid, err := strconv.Atoi(args[0])

	if err != nil {
		fmt.Println("Cannot parse bz id %s", args[0])
	}

	html, err := storecache.RetrieveCache(bzid, "", false)

	if err == nil {
		openHTML(bzid, html)
		os.Exit(0)
	}

	client := getClient()

	html, _, err = client.ShowBugHTML(bzid, "")
	if err == nil {
		openHTML(bzid, html)
		os.Exit(0)

	} else {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
}

func openHTML(bzid int, html *[]byte) {
	filename := fmt.Sprintf("/tmp/bz%d.html", bzid)
	fmt.Printf("Wrote %s\n", filename)
	writeHTML(html, filename)
	err := exec.Command(viper.Get("htmlOpenCommand").(string), filename).Run()
	if err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}
}

func writeHTML(html *[]byte, outputFile string) {
	htmlStr := string(*html)

	htmlStr = strings.Replace(htmlStr,"src=\"", "src=\"" + viper.Get("bzurl").(string) + "/" , -1)
	htmlStr = strings.Replace(htmlStr,"href=\"", "href=\"" + viper.Get("bzurl").(string) + "/" , -1)
	htmlStr = strings.Replace(htmlStr,"action=\"", "action=\"" + viper.Get("bzurl").(string) + "/" , -1)



	f, _ := os.Create(outputFile)
	defer f.Close()

	data := []byte(htmlStr)
	f.Write(data)


}


