/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/miekg/dns"
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
		if nsFlag, err := cmd.Flags().GetString("nameserver"); err != nil {
			return err
		} else if nsFlag == "" {
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
	Use:     "soa",
	Aliases: []string{"sao"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			args = append(args, ns)
		}

		c := new(dns.Client)
		msg := new(dns.Msg)
		msg.SetEdns0(4096, true)
		fmt.Println("sending question using name server", ns)

		for _, e := range args {
			msg.SetQuestion(dns.Fqdn(e), dns.TypeSOA)
			resp, _, err := c.Exchange(msg, ns+":53")
			if err != nil {
				log.Printf("Error: %s\n", err)
				return nil
			}

			fmt.Printf("SOA response for %s:\n%v", e, resp)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(nsCmd)
	nsCmd.AddCommand(lookupCmd)
	nsCmd.AddCommand(soaCmd)

	nsCmd.PersistentFlags().String("nameserver", "", "name server to use for DNS resolution")
}
