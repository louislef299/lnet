/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/louislef299/lnet/pkg/dns"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ns string

// nsCmd represents the ns command
var nsCmd = &cobra.Command{
	Use:   "ns",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if ns == "" {
			nameservers := viper.GetStringSlice("nameservers")
			ns = nameservers[0]
		}

		if ip := net.ParseIP(ns); ip != nil {
			addr, err := net.LookupAddr(ip.String())
			if err != nil {
				return err
			}
			ns = addr[0]
		}
		return nil
	},
}

// lookupCmd represents the lookup command
var lookupCmd = &cobra.Command{
	Use:     "lookup",
	Aliases: []string{"lookup", "lkup", "lk", "lup"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		for i, ns := range args {
			if i != 0 {
				fmt.Println()
			}
			if ip := net.ParseIP(ns); ip != nil {
				cname, err := net.LookupAddr(ip.String())
				if err != nil {
					log.Println(err)
					continue
				}
				fmt.Printf("CNAME for %s:\n", ns)
				for _, c := range cname {
					fmt.Println(c)
				}
			} else {
				ips, err := net.LookupHost(ns)
				if err != nil {
					log.Println(err)
					continue
				}
				fmt.Printf("IPs for %s:\n", ns)
				for _, addr := range ips {
					fmt.Printf("  - %s\n", addr)
				}
			}
		}
		return nil
	},
}

// soaCmd represents the soa command
var soaCmd = &cobra.Command{
	Use:   "soa",
	Short: "Retrieve a start of authority record from requested servers",
	Long: `Returns a SOA record containing administrative information about 
the zone, especially regarding zone transfers. This can include 
the email address of the administrator, when the domain was last
updated, and how long the server should wait between refreshes.

Response will follow this format:

[Owner name] [SOA record TTL] SOA [MNAME] [RNAME] (
	[Serial number]
	[Refresh time in seconds]
	[Retry time in seconds]
	[Expire time in seconds]
	[Minimum TTL (negative response cache TTL)]
)
  - The Primary Name Server (MNAME)
  - The Responsible Person (RNAME)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If endpoint isn't given, just send msg to currrent NS
		if len(args) < 1 {
			args = append(args, ns)
		}

		resp, err := dns.GetSoa(ns, args)
		if err != nil {
			log.Fatal("could not get soa response: ", err)
		}

		for _, r := range resp {
			fmt.Printf("SOA response for %s: %v\n", r.Server, r.Msg)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(nsCmd)
	nsCmd.AddCommand(lookupCmd)
	nsCmd.AddCommand(soaCmd)

	nsCmd.PersistentFlags().StringVar(&ns, "nameserver", "", "name server to use for DNS resolution")
}
