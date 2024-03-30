package cmd

import (
	"app/build"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s v%s\n", build.AppName, build.Version)
		fmt.Println("BuildCommit:\t", build.CommitID)
		fmt.Println("BuildDate:\t", build.Date)
		fmt.Println("BuildUser:\t", build.User)
	},
}
