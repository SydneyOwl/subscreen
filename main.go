/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/spf13/cobra"
	"sscreen/cmd"
)

func main() {
	cobra.MousetrapHelpText = ""
	cmd.Execute()
}
