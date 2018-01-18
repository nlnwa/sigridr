package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sigridrctl",
	Short: "Twitter API client",
	Long:  `Twitter API client`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatal()
	}
}
