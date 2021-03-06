// Copyright © 2018 Miguel Angel Ajo
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
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var cfgFile string

var listOffset int
var listLimit int
var ignoreSSLCerts bool

const trelloAppKey string = "95b7346a1c21c8dbe392f3b2d267fed3"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "track",
	Short: "Track helps you handle tasks, bugs, rfes across platforms",
	Long: `Track helps you track tasks, bugs, RFEs linked across platforms
like bugzilla, trello, launchpad, etc.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.track.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringP("htmlOpenCommand", "", "xdg-open", "Command to open an html file")
	rootCmd.PersistentFlags().IntVarP(&listOffset, "offset", "o", 0, "Offset on the bug listing")
	rootCmd.PersistentFlags().IntVarP(&listLimit, "limit", "l", 50, "Max entries to list")
	rootCmd.PersistentFlags().IntVarP(&workers, "workers", "w", 4, "Workers for http retrieval")
	rootCmd.PersistentFlags().BoolVarP(&ignoreSSLCerts, "ignorecerts", "i", false, "Ignore SSL certificates")
}

func exampleTrackYaml() {
	fmt.Print(`
An example ~/.track.yaml:

bzurl: https://bugzilla.redhat.com
bzemail: xxxxx@redhat.com
bzpass: xxxxxxxx
dfg: Networking
htmlOpenCommand: xdg-open  # note: for OSX use open instead
queries:
    ovn-new: https://bugzilla.redhat.com/buglist.cgi?bug_status=NEW&classification=Red%20Hat&component=python-networking-ovn&list_id=8959835&product=Red%20Hat%20OpenStack&query_format=advanced
    ovn-rfes: https://bugzilla.redhat.com/buglist.cgi?bug_status=NEW&bug_status=ASSIGNED&bug_status=MODIFIED&bug_status=ON_DEV&bug_status=POST&bug_status=ON_QA&classification=Red%20Hat&component=python-networking-ovn&f1=keywords&f2=short_desc&j_top=OR&list_id=8959855&o1=substring&o2=substring&product=Red%20Hat%20OpenStack&query_format=advanced&v1=RFE&v2=RFE
users:
    colleague1@email.com
    colleague2@email.com

`)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".track" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".track")
	}

	viper.SetEnvPrefix("TRACK")
	viper.AutomaticEnv() // read in environment variables that match
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		fmt.Printf("Could not read config file: %s \n", err)
	}

	for _, k := range []string {"htmlOpenCommand"} {
		viper.BindPFlag(k, rootCmd.PersistentFlags().Lookup(k))
	}

	initBzConfig()
}

func findEmail(user string) (email string){

	// if the user is "me", it's our email
	if user == "me" {
		return BzEmail
	}

	// If the user already has a @ we assume it's a valid email
	if strings.Contains(user, "@") ||  !viper.InConfig("users") {
		return user
	}

	users := viper.GetStringMapString("users")
	email, ok := users[user]
	if ok {
		return email
	} else {
		return user
	}
}