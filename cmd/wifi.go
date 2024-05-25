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
		c, err := wifi.New()
		if err != nil {
			log.Fatal(err)
		}

		ifs, err := net.Interfaces()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("scanning network interfaces")
		for _, i := range ifs {
			if i.Flags&net.FlagUp == 0 {
				log.Printf("%s FlagUp\n", i.Name)
				continue
			}

			if i.Flags&net.FlagLoopback != 0 {
				log.Printf("%s FlagLoopback\n", i.Name)
				continue
			}

			if i.Flags&net.FlagPointToPoint != 0 {
				log.Printf("%s FlagPointToPoint\n", i.Name)
				continue
			}

			// if i.Flags&net.FlagMulticast != 0 {
			// 	log.Printf("%s FlagMulticast\n", i.Name)
			// 	continue
			// }

			if i.HardwareAddr.String() == "" {
				log.Printf("%s EmptyHardwareAddr\n", i.Name)
				continue
			}
			fmt.Printf("Interface: %s (%s)\n", i.Name, i.HardwareAddr)

			wifs, err := c.Interfaces()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("found %d wifi interfaces\n", len(wifs))
			var wi *wifi.Interface
			for _, w := range wifs {
				if w.Name == i.Name {
					wi = w
				}
			}
			if wi == nil {
				fmt.Println("no wifi interfaces matching")
				continue
			}

			fmt.Printf("using wifi interface %v\n", *wi)
			sta, err := c.StationInfo(wi)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("found %d stations\n", len(sta))
			for _, s := range sta {
				fmt.Printf("  Station Info: %v\n", s)
			}
		}
		fmt.Println("wifi exiting...")
	},
}

func init() {
	rootCmd.AddCommand(wifiCmd)
}
