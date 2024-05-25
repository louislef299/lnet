/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"tinygo.org/x/bluetooth"
)

// bleCmd represents the icmp command
var bleCmd = &cobra.Command{
	Use:   "bluetooth",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		adapter := bluetooth.DefaultAdapter
		// Enable BLE interface.
		err := adapter.Enable()
		if err != nil {
			log.Fatal(err)
		}

		// Start scanning.
		println("scanning...")
		err = adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
			println("found device:", device.Address.String(), device.RSSI, device.LocalName())
		})
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(bleCmd)
}
