package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/svrana/heimdal"
	"github.com/svrana/heimdal/pkg/exec"
	"github.com/svrana/heimdal/pkg/zapped"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "heimdal",
	Short: "Run commands as your go code or its dependencies change",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Must specify a path or package to watch")
			os.Exit(1)
		}

		l, _ := zap.NewDevelopment()
		s := l.Sugar()
		defer l.Sync()
		ctx := zapped.NewContext(context.Background(), s)

		runners := make(map[string]*exec.Runner)

		for _, filename := range args {
			runners[filename] = exec.NewRunner("make")
		}

		heimdal.Start(ctx, runners)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $CWD/.heimdal.yaml)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
}

func initConfig() {
	viper.SetEnvPrefix("HEIMDAL")
	viper.AutomaticEnv()

	configFile := viper.GetString("config")
	if configFile != "" {
		fmt.Println("loading config file:", configFile)
		viper.SetConfigFile(configFile)
	} else {
		// Search config in current working directory with name ".heimdal" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigName(".heimdal.yaml")
	}

	viper.ReadInConfig()
}
