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
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/reaandrew/surge/cli"
	surgeHTTP "github.com/reaandrew/surge/infrastructure/http"
	"github.com/reaandrew/surge/utils"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func RightPad2Len(s string, padStr string, overallLen int) string {
	//https://github.com/git-time-metric/gtm/blob/master/util/string.go#L53-L88
	var padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}

var (
	cfgFile         string
	urlFile         string
	random          bool
	workerCount     int
	iterations      int
	processes       int
	output          string
	server          bool
	serverHost      string
	serverPort      int
	workerEndpoints []string
	Timer           utils.Timer          = &utils.DefaultTimer{}
	HttpClient      surgeHTTP.HttpClient = surgeHTTP.NewDefaultHttpClient()
)

const (
	TransactionsKey           = "Transactions"
	AvailabilityKey           = "Availability (%)"
	ElapsedTimeKey            = "Elapsed Time (ms)"
	TotalBytesSentKey         = "Total Bytes Sent"
	TotalBytesReceivedKey     = "Total Bytes Received"
	AverageResponseTimeKey    = "Average Response Time (ms)"
	AverageTransactionRateKey = "Average Transaction Rate (requests/sec)"
	ConcurrencyKey            = "Concurrency"
	DataSendRateKey           = "Data Send Rate (bytes/sec)"
	DataReceiveRateKey        = "Data Receive Rate (bytes/sec)"
	SuccessfulTransactionsKey = "Successfull Transactions"
	FailedTransactionsKey     = "Failed Transactions"
	LongestTransactionKey     = "Longest Transaction"
	ShortestTransactionKey    = "Shortest Transaction"
	WorkerCountKey            = "Worker Count"
	RandomKey                 = "Random"
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
		surgeClient := cli.NewSurgeCLIBuilder().
			SetURLFilePath(urlFile).
			SetRandom(random).
			SetWorkers(workerCount).
			SetIterations(iterations).
			SetServer(server).
			SetServerHost(serverHost).
			SetServerPort(serverPort).
			SetProcesses(processes).
			Build()

		result, err := surgeClient.Run()

		transactions := fmt.Sprintf("%v", result.Transactions)
		availability := fmt.Sprintf("%v", result.Availability*100)
		elapsedTime := fmt.Sprintf("%v", result.ElapsedTime.String())
		totalBytesSent := fmt.Sprintf("%v", humanize.Bytes(uint64(result.TotalBytesSent)))
		totalBytesReceived := fmt.Sprintf("%v", humanize.Bytes(uint64(result.TotalBytesReceived)))
		averageResponseTime := fmt.Sprintf("%.2f", result.AverageResponseTime/(float64(time.Millisecond)))
		transactionRate := fmt.Sprintf("%.2f", result.TransactionRate)
		concurency := fmt.Sprintf("%.2f", result.ConcurrencyRate)
		dataSendRate := fmt.Sprintf("%v", humanize.Bytes(uint64(result.DataSendRate)))
		dataReceiveRate := fmt.Sprintf("%v", humanize.Bytes(uint64(result.DataReceiveRate)))
		successfulTransactions := fmt.Sprintf("%v", result.SuccessfulTransactions)
		failedTransactions := fmt.Sprintf("%v", result.FailedTransactions)
		longestTransaction := time.Duration(result.LongestTransaction).String()
		shortestTransaction := time.Duration(result.ShortestTransaction).String()
		workerCount := fmt.Sprintf("%v", workerCount)
		randomEnabled := fmt.Sprintf("%v", random)

		cmd.Println(`
 ____  _   _ ____   ____ _____ 
/ ___|| | | |  _ \ / ___| ____|
\___ \| | | | |_) | |  _|  _|  
 ___) | |_| |  _ <| |_| | |___ 
|____/ \___/|_| \_\\____|_____|
		`)

		if err == nil {
			records := [][]string{
				{
					TransactionsKey,
					AvailabilityKey,
					ElapsedTimeKey,
					TotalBytesSentKey,
					TotalBytesReceivedKey,
					AverageResponseTimeKey,
					AverageTransactionRateKey,
					ConcurrencyKey,
					DataSendRateKey,
					DataReceiveRateKey,
					SuccessfulTransactionsKey,
					FailedTransactionsKey,
					LongestTransactionKey,
					ShortestTransactionKey,
					WorkerCountKey,
					RandomKey,
				},
				{
					transactions,
					availability,
					elapsedTime,
					totalBytesSent,
					totalBytesReceived,
					averageResponseTime,
					transactionRate,
					concurency,
					dataSendRate,
					dataReceiveRate,
					successfulTransactions,
					failedTransactions,
					longestTransaction,
					shortestTransaction,
					workerCount,
					randomEnabled,
				},
			}

			fileLines := records[:]

			usr, err := user.Current()
			if err != nil {
				log.Fatal(err)

			}
			var resultsFile *os.File
			defer resultsFile.Close()
			name := fmt.Sprintf("%v/.surge.results", usr.HomeDir)
			if _, err := os.Stat(name); os.IsNotExist(err) {
				// path/to/whatever does not exist
				resultsFile, err = os.Create(name)
				if err != nil {
					panic(err)
				}
			} else {
				fileLines = [][]string{records[1]}
				resultsFile, err = os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					panic(err)
				}
			}

			w := csv.NewWriter(resultsFile)

			for _, record := range fileLines {
				if err := w.Write(record); err != nil {
					log.Fatalln("error writing record to csv:", err)
				}
			}
			w.Flush()
			switch output {
			case "csv":
				w := csv.NewWriter(os.Stdout)

				for _, record := range records {
					if err := w.Write(record); err != nil {
						log.Fatalln("error writing record to csv:", err)
					}
				}
				w.Flush()
			default:
				cmd.Println(fmt.Sprintf("%v: %v", RightPad2Len(TransactionsKey, ".", 45), transactions))
				cmd.Println(fmt.Sprintf("%v: %v", RightPad2Len(AvailabilityKey, ".", 45), availability))
				cmd.Println(fmt.Sprintf("%v: %v", RightPad2Len(ElapsedTimeKey, ".", 45), elapsedTime))
				cmd.Println(fmt.Sprintf("%v: %v", RightPad2Len(TotalBytesSentKey, ".", 45), totalBytesSent))
				cmd.Println(fmt.Sprintf("%v: %v", RightPad2Len(TotalBytesReceivedKey, ".", 45), totalBytesReceived))
				cmd.Println(fmt.Sprintf("%v: %v", RightPad2Len(AverageResponseTimeKey, ".", 45), averageResponseTime))
				cmd.Println(fmt.Sprintf("%v: %v", RightPad2Len(AverageTransactionRateKey, ".", 45), transactionRate))
				cmd.Println(fmt.Sprintf("%v: %v", RightPad2Len(ConcurrencyKey, ".", 45), concurency))
				cmd.Println(fmt.Sprintf("%v: %v", RightPad2Len(DataSendRateKey, ".", 45), dataSendRate))
				cmd.Println(fmt.Sprintf("%v: %v", RightPad2Len(DataReceiveRateKey, ".", 45), dataReceiveRate))
				cmd.Println(fmt.Sprintf("%v: %v", RightPad2Len(SuccessfulTransactionsKey, ".", 45), successfulTransactions))
				cmd.Println(fmt.Sprintf("%v: %v", RightPad2Len(FailedTransactionsKey, ".", 45), failedTransactions))
				cmd.Println(fmt.Sprintf("%v: %v", RightPad2Len(LongestTransactionKey, ".", 45), longestTransaction))
				cmd.Println(fmt.Sprintf("%v: %v", RightPad2Len(ShortestTransactionKey, ".", 45), shortestTransaction))
				cmd.Println(fmt.Sprintf("%v: %v", RightPad2Len(WorkerCountKey, ".", 45), workerCount))
				cmd.Println(fmt.Sprintf("%v: %v", RightPad2Len(RandomKey, ".", 45), randomEnabled))
			}
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
	RootCmd.PersistentFlags().StringVarP(&output, "output", "o", "default", "The output format to use for the results")
	RootCmd.PersistentFlags().BoolVarP(&random, "random", "r", false, "Read the urls in random order")
	RootCmd.PersistentFlags().IntVarP(&workerCount, "worker-count", "c", 1, "The number of concurrent virtual users")
	RootCmd.PersistentFlags().IntVarP(&iterations, "number-iterations", "n", 1, "The number of iterations per virtual user")
	RootCmd.PersistentFlags().IntVarP(&processes, "processes", "p", 1, "The number of processes to run virtual users")

	RootCmd.PersistentFlags().BoolVar(&server, "server", false, "Set in server mode")
	RootCmd.PersistentFlags().IntVar(&serverPort, "server-port", 51234, "The port thew server should bind to")
	RootCmd.PersistentFlags().StringVar(&serverHost, "server-host", "localhost", "The hostname the server should bind to")
	RootCmd.PersistentFlags().StringArrayVar(&workerEndpoints, "worker-endpoints", []string{}, "The number of processes to run virtual users")

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
