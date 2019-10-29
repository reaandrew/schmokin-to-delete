/*
Copyright © 2019 Andy Rea <email@andrewrea.co.uk>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/reaandrew/surge/client"
	"github.com/reaandrew/surge/utils"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	urlFile     string
	random      bool
	workerCount int
	iterations  int
	Timer       utils.Timer       = &utils.DefaultTimer{}
	HttpClient  client.HttpClient = client.NewDefaultHttpClient()
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "surge",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		surgeClient := client.NewSurgeClientBuilder().
			SetURLFilePath(urlFile).
			SetRandom(random).
			SetWorkers(workerCount).
			SetIterations(iterations).
			SetHTTPClient(HttpClient).
			SetTimer(Timer).
			Build()

		result, err := surgeClient.Run()

		if err == nil {
			cmd.Println(fmt.Sprintf("Transactions: %v", result.Transactions))
			cmd.Println(fmt.Sprintf("Availability: %v%%", result.Availability*100))
			cmd.Println(fmt.Sprintf("Elapsed Time: %v", result.ElapsedTime.String()))
			cmd.Println(fmt.Sprintf("Total Bytes Sent: %v", humanize.Bytes(uint64(result.TotalBytesSent))))
			cmd.Println(fmt.Sprintf("Total Bytes Received: %v", humanize.Bytes(uint64(result.TotalBytesReceived))))
		}
		return err
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.surge.yaml)")
	RootCmd.PersistentFlags().StringVarP(&urlFile, "urls", "u", "", "The urls file to use")
	RootCmd.MarkPersistentFlagRequired("urls")
	RootCmd.PersistentFlags().BoolVarP(&random, "random", "r", false, "Read the urls in random order")
	RootCmd.PersistentFlags().IntVarP(&workerCount, "worker-count", "c", 1, "The number of concurrent virtual users")
	RootCmd.PersistentFlags().IntVarP(&iterations, "number-iterations", "n", 1, "The number of iterations per virtual user")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
			os.Exit(1)
		}

		// Search config in home directory with name ".surge" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".surge")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
