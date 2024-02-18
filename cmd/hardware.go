/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/jaypipes/ghw"
	"github.com/spf13/cobra"
)

var disableWarnings bool

// hardwareCmd represents the hardware command
var hardwareCmd = &cobra.Command{
	Use:     "hardware",
	Aliases: []string{"hw"},
	Short:   "Gather general system hardware information",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if disableWarnings {
			os.Setenv("GHW_DISABLE_WARNINGS", "1")
		}
		return nil
	},
}

// hardwareBIOSCmd represents the bios command
var hardwareBIOSCmd = &cobra.Command{
	Use:   "bios",
	Short: "Gather information about the host computer's basis input/output system (BIOS)",
	Run: func(cmd *cobra.Command, args []string) {
		bios, err := ghw.BIOS()
		if err != nil {
			fmt.Printf("Error getting BIOS info: %v", err)
		}

		fmt.Printf("%v\n", bios)
	},
}

// hardwareCPUCmd represents the cpu command
var hardwareCPUCmd = &cobra.Command{
	Use:   "cpu",
	Short: "Get CPU system information",
	Run: func(cmd *cobra.Command, args []string) {
		cpu, err := ghw.CPU()
		if err != nil {
			fmt.Printf("Error getting CPU info: %v", err)
		}

		fmt.Printf("%v\n", cpu)

		for _, proc := range cpu.Processors {
			fmt.Printf(" %v\n", proc)
			fmt.Printf(" %s (%s)\n", proc.Vendor, proc.Model)
			for _, core := range proc.Cores {
				fmt.Printf("  %v\n", core)
			}
			if len(proc.Capabilities) > 0 {
				// pretty-print the (large) block of capability strings into rows
				// of 6 capability strings
				rows := int(math.Ceil(float64(len(proc.Capabilities)) / float64(6)))
				for row := 1; row < rows; row = row + 1 {
					rowStart := (row * 6) - 1
					rowEnd := int(math.Min(float64(rowStart+6), float64(len(proc.Capabilities))))
					rowElems := proc.Capabilities[rowStart:rowEnd]
					capStr := strings.Join(rowElems, " ")
					if row == 1 {
						fmt.Printf("  capabilities: [%s\n", capStr)
					} else if rowEnd < len(proc.Capabilities) {
						fmt.Printf("                 %s\n", capStr)
					} else {
						fmt.Printf("                 %s]\n", capStr)
					}
				}
			}
		}
	},
}

// hardwareNICCmd represents the nic command
var hardwareNICCmd = &cobra.Command{
	Use:     "nics",
	Aliases: []string{"nic"},
	Short:   "List the network interfaces on the system",
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
	hardwareCmd.AddCommand(hardwareNICCmd)
	hardwareCmd.AddCommand(hardwareCPUCmd)
	hardwareCmd.AddCommand(hardwareBIOSCmd)

	hardwareCmd.Flags().BoolVar(&disableWarnings, "disableWarnings", false, "disable verbose warning output when looking at hardware information")
}
