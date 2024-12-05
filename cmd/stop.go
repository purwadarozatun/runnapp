package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the application",
	Run: func(cmd *cobra.Command, args []string) {

		if configFile == "" {
			configFile = cofigPath

		}

		config, err := loadConfig(configFile)
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}
		// // Get the PID of the running application
		pid, err := readPidFile(config.SearchPid)
		if err != nil {
			fmt.Println("Error reading PID file:", err)
			os.Exit(1)
		}
		// kill the application
		if err := killProcess(pid); err != nil {
			fmt.Println("Error killing process:", err)
		}
		if err := os.Remove(pidPath); err != nil {
			fmt.Println("Error removing PID file:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
