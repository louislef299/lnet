/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"log"
	"os"
	"strconv"

	"github.com/louislef299/lnet/pkg/icmp"
	"github.com/spf13/cobra"
)

// icmpCmd represents the icmp command
var icmpCmd = &cobra.Command{
	Use:   "icmp",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if usrID := os.Getuid(); usrID != 0 {
			log.Println("current user id:", usrID)
			log.Fatal("icmp command requires 'root' privileges(rerun with 'sudo')")
		}
		if i, _ := strconv.Atoi(os.Args[1]); i == 1 {
			icmp.BottomOfIt()
		} else {
			icmp.IcmpScan()
		}
	},
}

func init() {
	rootCmd.AddCommand(icmpCmd)
}
