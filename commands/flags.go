package commands

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func addBuildFlags(cmdSet *pflag.FlagSet) {
	// bind to viper
	err := viper.BindPFlags(cmdSet)
	if err != nil {
		panic(err)
	}
}
