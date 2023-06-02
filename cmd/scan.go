/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"log"
	"time"

	"github.com/louislef299/lnet/pkg/port"
	"github.com/spf13/cobra"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "port scans on a specified target",
	Long: `A port scanner is an application designed to probe a server or 
host for open ports. Such an application may be used by 
administrators to verify security policies of their networks 
and by attackers to identify network services running on a 
host and exploit vulnerabilities.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		t, err := cmd.Flags().GetString("target")
		if err != nil {
			return err
		}
		r, err := cmd.Flags().GetInt("range")
		if err != nil {
			return err
		}

		log.Printf("Port Scanning %d range on %s", r, t)

		start := time.Now()
		results, done := port.PortScan(ctx, t, r)
		for {
			select {
			case r := <-results:
				if port.IsOpen(r.State) {
					log.Printf("%s:%d is Open\n", r.Protocol, r.Port)
				}
			case <-done:
				log.Printf("scan took %v\n", time.Since(start).Truncate(time.Second))
				return nil
			case <-ctx.Done():
				log.Println("recieved SIGINT, exiting")
				return nil
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringP("target", "t", "localhost", "the target to scan")
	scanCmd.Flags().IntP("range", "r", 49152, "port range to scan")
}
