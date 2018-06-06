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
)

var cfgFile string
var bzEmail string
var bzPassword string
var bzURL string
var whiteBoardQuery string
var myBugs bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "track",
	Short: "Track helps you handle tasks, bugs, rfes across platforms",
	Long: `Track helps hou track tasks, bugs, RFEs linked across platforms
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
	rootCmd.PersistentFlags().StringP("bzurl", "a", "https://bugzilla.redhat.com", "Bugzilla URL")
	rootCmd.PersistentFlags().StringP("bzemail", "u", "", "Bugzilla login email")
	rootCmd.PersistentFlags().StringP("bzpass", "p", "", "Bugzilla login password")
	rootCmd.PersistentFlags().StringP("dfg", "d", "", "Openstack DFG")
	rootCmd.PersistentFlags().StringP("squad", "s", "", "Openstack DFG Squad")
	rootCmd.PersistentFlags().BoolVarP(&myBugs,"my-bugs", "m", false,"List only my bugs")
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

	for _, k := range []string {"bzurl", "bzemail", "bzpass", "dfg", "squad"} {
		viper.BindPFlag(k, rootCmd.PersistentFlags().Lookup(k))
	}

	bzURL = viper.GetString("bzurl")
	bzPassword = viper.GetString("bzpass")
	bzEmail = viper.GetString("bzemail")
	whiteBoardQuery = ""
	if viper.GetString("dfg") != "" {
		whiteBoardQuery += fmt.Sprintf("DFG:%s", viper.GetString("dfg"))
	}

	if viper.GetString("squad") != "" {
		if viper.GetString("dfg") == "" {
			panic(fmt.Errorf("When specifying a squad, please provide DFG too"))
		}
		whiteBoardQuery += fmt.Sprintf(" Squad:%s", viper.GetString("squad"))
	}


}
