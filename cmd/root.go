package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/svrana/heimdal"
	"github.com/svrana/heimdal/pkg/config"
	"github.com/svrana/heimdal/pkg/exec"
	"github.com/svrana/heimdal/pkg/zapped"
)

var cfgFile string
var cfg config.Blob
var logLevel string

var rootCmd = &cobra.Command{
	Use:   "heimdal",
	Short: "Run commands as your go code or its dependencies change",
	Run: func(cmd *cobra.Command, _ []string) {
		if len(cfg.Targets) == 0 {
			fmt.Println("Must specify a path or package to watch")
			os.Exit(1)
		}

		l, _ := zap.NewDevelopment()
		s := l.Sugar()
		defer l.Sync()
		ctx := zapped.NewContext(context.Background(), s)

		runners := make(map[string]*exec.Runner)
		for _, t := range cfg.Targets {
			runners[t.Pkg] = exec.NewRunner(t.Name, t.Command, t.Args...).
				WithDelayDuration(cfg.DelayMS)
		}

		heimdal.Start(ctx, runners)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $CWD/.heimdal.toml)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.PersistentFlags().StringVar(&logLevel, "logs", "info", "log level (debug, info, fatal)")
	viper.BindPFlag("logs", rootCmd.PersistentFlags().Lookup("logs"))

	rootCmd.PersistentFlags().IntVar(&cfg.DelayMS, "delay", 500, "milliseconds to delay before triggering command")
	viper.BindPFlag("delay", rootCmd.PersistentFlags().Lookup("delay"))
}

func initConfig() {
	viper.SetEnvPrefix("HEIMDAL")
	viper.AutomaticEnv()

	configFile := viper.GetString("config")
	if configFile != "" {
		fmt.Println("loading config file:", configFile)
		viper.SetConfigFile(configFile)
	} else {
		// Search config in current working directory with name ".heimdal"
		viper.AddConfigPath(".")
		viper.SetConfigName(".heimdal")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// no config file, ok.. command line maybe
		} else {
			fmt.Println("Error parsing config:", err)
			os.Exit(1)
		}
	} else {
		if err = viper.Unmarshal(&cfg); err != nil {
			fmt.Println("unable to decode into struct:", err)
			os.Exit(1)
		}
	}
}
