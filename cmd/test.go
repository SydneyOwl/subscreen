/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"sscreen/disp_func"
	"sscreen/driver"
	"time"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test",
	Long:  `Test if it works.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting Application...")
		if err := driver.InitSerial(); err != nil {
			fmt.Printf("Cannot initialize serial port: %v\n", err)
			return
		}
		fmt.Println("Init Application...")
		disp_func.InitScreen()
		fmt.Println("Start Update screen...")
		disp_func.StartAll()
		time.Sleep(time.Second * 10)
		disp_func.StopAll()
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
