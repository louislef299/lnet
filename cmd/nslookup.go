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
	Short: "Gathers DNS records for a domain name",
	Long: `A flexible tool for interrogating name servers. Also
can gather/return/configure local DNS services.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		initNameServer()
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
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := dns.GetPath()
		if err != nil {
			return err
		}
		fmt.Printf("Gathering local information from %s\n", path)

		ns, err := dns.GetLocalNS()
		if err != nil {
			return err
		}
		printInfo("local name servers:", ns)

		sd, err := dns.GetLocalSearchDomains()
		if err != nil {
			return err
		}
		printInfo("local search domains:", sd)

		opt, err := dns.GetLocalOptions()
		if err != nil {
			return err
		}
		printInfo("local options:", opt)

		return nil
	},
}

// lookupCmd represents the lookup command
var lookupCmd = &cobra.Command{
	Use:     "lookup",
	Aliases: []string{"lookup", "lkup", "lk", "lup"},
	Short:   "Lookup IP address for specified domain",
	Long: `A DNS lookup, or DNS record lookup, is the process 
through which human-readable domain names are 
translated into a computer-readable IP address`,
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
				printInfo(fmt.Sprintf("CNAME for %s:", ns), cname)
			} else {
				ips, err := net.LookupHost(ns)
				if err != nil {
					log.Println(err)
					continue
				}
				printInfo(fmt.Sprintf("IPs for %s:", ns), ips)
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
		// If endpoint isn't given, just send msg to current NS
		if len(args) < 1 {
			args = append(args, ns)
		}

		resp, err := dns.GetSoa(ns, args)
		if err != nil {
			log.Fatal("could not get soa response: ", err)
		}

		if len(resp) != len(args) {
			log.Fatal("mismatched entry lengths")
		}

		raw, err := cmd.Flags().GetBool("raw")
		if err != nil {
			return err
		}
		for i, r := range resp {
			fmt.Printf("SOA response for %s:\n", args[i])
			if raw {
				fmt.Println(r.Msg)
			} else {
				fmt.Printf("%s\n", r.String())
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(nsCmd)
	nsCmd.AddCommand(lookupCmd)
	nsCmd.AddCommand(soaCmd)

	nsCmd.PersistentFlags().StringVar(&ns, "nameserver", "", "name server to use for DNS resolution")

	soaCmd.Flags().Bool("raw", false, "prints out raw dns value")
}

// Run to initialize local name servers if commands are concerned with DNS
func initNameServer() {
	if n := viper.GetString("nameserver"); n == "" {
		ns, err := dns.GetLocalNS()
		if err != nil {
			log.Println("could not gather local name servers:", err)
		}
		viper.Set("nameservers", ns)
	}

	if err := viper.WriteConfig(); err != nil {
		log.Println("couldn't write to config:", err)
	}
}

func printInfo(header string, retval []string) {
	if len(retval) > 0 {
		fmt.Println(header)
		for _, r := range retval {
			fmt.Printf("  -%s\n", r)
		}
	}
}
