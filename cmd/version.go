/*
Copyright © 2022 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/louislef299/lnet/pkg/version"
	"github.com/spf13/cobra"
)

var short bool

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"ver", "vers"},
	Short:   "print the version for lnet",
	Long:    `print the version for lnet`,
	Hidden:  true,
	Run: func(cmd *cobra.Command, args []string) {
		if short {
			fmt.Println(version.String())
		} else {
			err := version.PrintVersion(os.Stdout, rootCmd)
			if err != nil {
				log.Fatal("couldn't print version:", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&short, "short", "s", false, "print out just the lnet version number")
}
