package main

import (
	"errors"
	"fmt"
	"math"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	_flagAddress     = "addr"
	_flagBrightness  = "brightness"
	_flagTemperature = "temperature"
)

func main() {
	cmdRoot := &cobra.Command{
		Use:                   "elgato-cli",
		Short:                 "Control Elgato light",
		DisableFlagsInUseLine: true,
		CompletionOptions:     cobra.CompletionOptions{HiddenDefaultCmd: true},
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := cmd.Flags().GetIP(_flagAddress)
			if err != nil {
				return err
			}

			l, err := NewClient(addr).GetLight()
			if err != nil {
				return err
			}

			fmt.Println(l)

			return nil
		},
	}
	cmdOn := &cobra.Command{
		Use:   "on",
		Short: "Turn Elgato light ON",
		RunE: func(cmd *cobra.Command, args []string) error {
			return handleState(1, cmd.Flags())
		},
	}
	cmdOff := &cobra.Command{
		Use:   "off",
		Short: "Turn Elgato light OFF",
		RunE: func(cmd *cobra.Command, args []string) error {
			return handleState(0, cmd.Flags())
		},
	}

	cmdRoot.AddCommand(cmdOn, cmdOff)
	cmdRoot.PersistentFlags().IPP(_flagAddress, "a", nil, "IP address of the light")
	cmdRoot.PersistentFlags().IntP(_flagBrightness, "b", 0, "brightness in percent; a value between 0 and 100")
	cmdRoot.PersistentFlags().IntP(_flagTemperature, "t", 0, "temperature in Kelvins; a value between 2900 and 7000")

	if err := cmdRoot.MarkPersistentFlagRequired(_flagAddress); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := cmdRoot.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleState(on int, flags *pflag.FlagSet) error {
	addr, err := flags.GetIP(_flagAddress)
	if err != nil {
		return err
	}

	l := &Light{On: &on}

	if flags.Lookup(_flagBrightness).Changed {
		b, err := flags.GetInt(_flagBrightness)
		if err != nil {
			return err
		}

		if b < 0 || b > 100 {
			return errors.New("incorrect brightness")
		}

		l.Brightness = &b
	}

	if flags.Lookup(_flagTemperature).Changed {
		t, err := flags.GetInt(_flagTemperature)
		if err != nil {
			return err
		}

		if t < 2900 || t > 7000 {
			return errors.New("incorrect temperature")
		}

		t = int(math.Round(1_000_000 / float64(t)))
		l.Temperature = &t
	}

	return NewClient(addr).UpdateLight(l)
}
