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
