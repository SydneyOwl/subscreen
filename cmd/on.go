/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	ipc "github.com/james-barrow/golang-ipc"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"sscreen/disp_func"
	"sscreen/driver"
	"strconv"
	"strings"
)

var background bool = false

// onCmd represents the on command
var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Start Application",
	Long:  `Start Application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if background {
			// check if i am the only process
			s, err := ipc.StartServer("sscreenPipe", nil)
			if err != nil {
				fmt.Printf("server error: %v", err)
				return
			}
			// start screen update
			if err := driver.InitSerial(); err != nil {
				fmt.Printf("Cannot connect serial: %v", err)
				return
			}
			disp_func.InitScreen()
			disp_func.StartAll()
			// start listening to requests: blocked!
			fmt.Println("Server started")
			for {
				message, err := s.Read()
				if err == nil {
					if message.MsgType == -1 {
						fmt.Println("server status", s.Status())
					} else {
						fmt.Println("Server received: " + string(message.Data))
						order := strings.Split(string(message.Data), ",")[0]
						value := strings.Split(string(message.Data), ",")[1]
						switch order {
						case "off":
							s.Write(5, []byte("R"))
							disp_func.SetBrightness(255)
							disp_func.StopAll()
							return
						case "light":
							valueI, _ := strconv.Atoi(value)
							disp_func.SetBrightness(valueI)
						default:
							fmt.Println("No such order")
						}
						s.Write(5, []byte("R"))
						//s.Close()
						//return
					}
				} else {
					fmt.Printf("Failed to read message: %v", err)
				}
			}
		} else {
			cmd := exec.Command(os.Args[0], "on", "-b")
			fmt.Println(os.Args[0])
			cmd.Start()
			fmt.Printf("Starting new application which has pid %d...", cmd.Process.Pid)
		}
	},
}

func init() {
	rootCmd.AddCommand(onCmd)
	onCmd.PersistentFlags().BoolVarP(&background, "run-background", "b", false, "Run at Background")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// onCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// onCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
