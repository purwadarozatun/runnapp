package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the application",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the PID of the running application
		if configFile == "" {
			configFile = cofigPath

		}

		config, err := loadConfig(configFile)
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}
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
		time.Sleep(2 * time.Second)
		startCmd.Run(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
