package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the application",
	Run: func(cmd *cobra.Command, args []string) {

		if configFile == "" {
			configFile = cofigPath

		}

		config, err := loadConfig(configFile)
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}

		// // Ensure the PID file is deleted on exit
		// defer func() {
		// 	if err := os.Remove(pidPath); err != nil {
		// 		fmt.Println("Error removing PID file:", err)
		// 	}
		// }()

		// // Handle interrupt signal to remove PID file on Ctrl+C
		// c := make(chan os.Signal, 1)
		// signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		// go func() {
		// 	<-c
		// 	if err := os.Remove(pidPath); err != nil {
		// 		fmt.Println("Error removing PID file:", err)
		// 	}
		// 	os.Exit(1)
		// }()

		// check if the application is already running

		pid, err := readPidFile(config.SearchPid)
		if err != nil {
			fmt.Println("Error reading PID file:", err)
			os.Exit(1)
		}
		if pid != 0 {
			fmt.Println("The application is already running")
			os.Exit(1)
		}
		runCommand(getCallerFuncName(), config.Program)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
