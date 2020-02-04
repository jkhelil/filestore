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
	short  string

	defaultValue interface{}
	kind string
}

func addFlag(flagset *pflag.FlagSet, f *flag) {
	switch f.kind {
	case "bool":
		if f.defaultValue != nil {
			flagset.BoolP(f.name, f.short, f.defaultValue.(bool), f.desc)
		} else {
			flagset.BoolP(f.name, f.short, false, f.desc)
		}
	case "int":
		if f.defaultValue != nil {
			flagset.IntP(f.name, f.short, f.defaultValue.(int), f.desc)
		} else {
			flagset.IntP(f.name, f.short, 0, f.desc)
		}
	default:
		if f.defaultValue != nil {
			flagset.StringP(f.name, f.short, f.defaultValue.(string), f.desc)
		} else {
			flagset.StringP(f.name, f.short, "", f.desc)
		}
	}
	flagset.VisitAll(func(flag *pflag.Flag) {
		check(viper.BindPFlag(flag.Name, flag))
	})
}

// check prints the error & exits the program with code 1 if err is non-nil
func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err) // nolint: errcheck
		os.Exit(1)
	}
}
