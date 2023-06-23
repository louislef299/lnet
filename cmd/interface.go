/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/spf13/cobra"
)

// interfaceCmd represents the interface command
var interfaceCmd = &cobra.Command{
	Use:     "interface",
	Aliases: []string{"if", "inter"},
	Short:   "configure and find system network interfaces",
	Long:    `Used to configure and find system network interfaces.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var interfaces []net.Interface
		if len(args) == 0 {
			interfaces, err = net.Interfaces()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			for _, n := range args {
				iname, err := net.InterfaceByName(n)
				if err != nil {
					log.Fatal(err)
				}
				interfaces = append(interfaces, *iname)
			}
		}
		for _, i := range interfaces {
			fmt.Println(printInterface(i))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(interfaceCmd)
}

func printInterface(i net.Interface) string {
	return fmt.Sprintf("(%d)%s:\n  flags: <%v>\n  mtu: %d\n  hardware address: %s",
		i.Index, i.Name, i.Flags, i.MTU, i.HardwareAddr.String())
}
