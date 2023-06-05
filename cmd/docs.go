/*
Copyright © 2022 Louis Lefebvre <lefebl4@medtronic.com>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var thisRoot *cobra.Command

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:    "docs",
	Short:  "Generate clctl docs",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			return err
		}
		if dir == "" {
			dir = os.TempDir()
		}

		return docsAction(os.Stdout, dir)
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)

	docsCmd.Flags().StringP("dir", "d", "", "Destination directory for docs")
	thisRoot = rootCmd
}

func docsAction(out io.Writer, dir string) error {
	if err := doc.GenMarkdownTree(thisRoot, dir); err != nil {
		return err
	}

	_, err := fmt.Fprintf(out, "Documentation successfully created in %s\n", dir)
	return err
}
