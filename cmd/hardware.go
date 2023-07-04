/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/jaypipes/ghw"
	"github.com/spf13/cobra"
)

var disableWarnings bool

// hardwareCmd represents the hardware command
var hardwareCmd = &cobra.Command{
	Use:     "hardware",
	Aliases: []string{"hw"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if disableWarnings {
			os.Setenv("GHW_DISABLE_WARNINGS", "1")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		net, err := ghw.Network()
		if err != nil {
			fmt.Printf("Error getting network info: %v", err)
		}

		fmt.Printf("%v\n", net)

		for _, nic := range net.NICs {
			fmt.Printf(" %v\n", nic)

			enabledCaps := make([]int, 0)
			for x, cap := range nic.Capabilities {
				if cap.IsEnabled {
					enabledCaps = append(enabledCaps, x)
				}
			}
			if len(enabledCaps) > 0 {
				fmt.Printf("  enabled capabilities:\n")
				for _, x := range enabledCaps {
					fmt.Printf("   - %s\n", nic.Capabilities[x].Name)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(hardwareCmd)

	hardwareCmd.Flags().BoolVar(&disableWarnings, "disableWarnings", false, "disable verbose warning output when looking at hardware information")
}
