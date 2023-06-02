/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"context"
	"errors"
	"log"
	"os"
	"path"

	"github.com/louislef299/lnet/pkg/dns"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lnet",
	Short: "A network tool for the modern system administrator",
	Long: `A network tool for the modern system administrator
                    _______ ______
                    |     / |    /
         O          |    /  |   /
                    |   /   |  /
      o  O 0         \  \   \  \
      o               \  \   \  \
         o            /  /   /  /
          o     /\_  /\\\   /  /
           O  /    /    /     /
   ..       /    /    /\=    /
  .  ))))))) = /====/    \
  . (((((((( /    /\=  _ }
  . |-----_|_+( /   \}
  . \_<\_//|  \  \ }
   ...=Q=  |==)\  \
     \----/     ) )
               / /
              /=/ 
            \|/
            o}`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx context.Context) {
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lnet.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		cfgName := ".lnet"
		// Search config in home directory with name ".lnet" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(cfgName)
		if _, err := os.Stat(path.Join(home, cfgName)); errors.Is(err, os.ErrNotExist) {
			f, err := os.OpenFile(path.Join(home, cfgName), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0744)
			if err != nil {
				log.Fatal(err)
			}
			f.Close()
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Couldn't use config file:", viper.ConfigFileUsed())
	}

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
