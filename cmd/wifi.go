/*
Copyright © 2024 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/mdlayher/wifi"
	"github.com/spf13/cobra"
)

// wifiCmd represents the wifi command
var wifiCmd = &cobra.Command{
	Use:   "wifi",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("wifi called")

		c, err := wifi.New()
		if err != nil {
			log.Fatal(err)
		}

		ifs, err := net.Interfaces()
		if err != nil {
			log.Fatal(err)
		}

		for _, i := range ifs {
			if i.Flags&net.FlagUp == 0 {
				continue
			}

			if i.Flags&net.FlagLoopback != 0 {
				continue
			}

			if i.Flags&net.FlagPointToPoint != 0 {
				continue
			}

			if i.Flags&net.FlagMulticast != 0 {
				continue
			}

			if i.HardwareAddr.String() == "" {
				continue
			}

			fmt.Printf("Interface: %s (%s)\n", i.Name, i.HardwareAddr)

			sta, err := c.Station(i.Name)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("  Station: %s\n", sta)

			aps, err := c.Scan(i.Name)
			if err != nil {
				log.Fatal(err)
			}

			for _, ap := range aps {
				fmt.Printf("  AP: %s\n", ap)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(wifiCmd)
}
