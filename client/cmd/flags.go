package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type flag struct {
	name   string
	desc   string
	envVar string

	defaultValue interface{}
	required     bool

	kind string
}

func addFlag(flagset *pflag.FlagSet, f *flag) {
	switch f.kind {
	case "bool":
		if f.defaultValue != nil {
			flagset.Bool(f.name, f.defaultValue.(bool), f.desc)
		} else {
			flagset.Bool(f.name, false, f.desc)
		}
	case "int":
		if f.defaultValue != nil {
			flagset.Int(f.name, f.defaultValue.(int), f.desc)
		} else {
			flagset.Int(f.name, 0, f.desc)
		}
	default:
		if f.defaultValue != nil {
			flagset.String(f.name, f.defaultValue.(string), f.desc)
		} else {
			flagset.String(f.name, "", f.desc)
		}
	}

	if f.envVar != "" {
		check(viper.BindEnv(f.name, f.envVar))
	}

}

// check prints the error & exits the program with code 1 if err is non-nil
func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err) // nolint: errcheck
		os.Exit(1)
	}
}
