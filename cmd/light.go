/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"sscreen/util"
	"strconv"
)

// lightCmd represents the light command
var lightCmd = &cobra.Command{
	Use:   "light",
	Short: "Set brightness",
	Long:  `Set brightness.`,
	Run: func(cmd *cobra.Command, args []string) {
		bg, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Input should be an integer")
			return
		}
		if bg < 0 || bg > 255 {
			fmt.Println("Invalid input: should between 0 and 255!")
			return
		}
		if err := util.Send(fmt.Sprintf("light,%d", bg)); err != nil {
			fmt.Printf("Error: %v", err)
		}
	},
}

func init() {
	setCmd.AddCommand(lightCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lightCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lightCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
