package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/chadgrant/dynamodb-go-sample/store/server"
	"github.com/chadgrant/go-http-infra/infra/cmds"
	"github.com/chadgrant/go-http-infra/infra/health"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dynamodb-go-sample",
	Short: "A sample app using dynamodb and golang",
	RunE: func(cmd *cobra.Command, args []string) error {

		cfg, err := server.Load(cfgFile)
		if err != nil {
			return err
		}

		srv, err := server.New(cfg)
		if err != nil {
			return err
		}

		done := make(chan interface{}, 1)
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		go func(quit <-chan os.Signal, srv *server.Server) {
			<-quit
			srv.Shutdown(done)
		}(quit, srv)

		if err := srv.Serve(done, cfg.Service.Address); err != nil {
			return err
		}

		<-done
		fmt.Println("\nServer Stopped")
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	getConfig := func(file string) (interface{}, error) {
		return server.Load(file)
	}
	getHealth := func(cfg interface{}) health.HealthChecker {
		return server.RegisterHealthChecks(cfg.(*server.Configuration))
	}

	cmds.Register(rootCmd, getConfig, getHealth)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	fs := rootCmd.PersistentFlags()
	fs.StringVar(&cfgFile, "config", "config.yaml", "path to config file")
}
