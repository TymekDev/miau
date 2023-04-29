package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var version string

const (
	_flagAddress     = "addr"
	_flagBrightness  = "brightness"
	_flagPort        = "port"
	_flagTemperature = "temperature"
)

func main() {
	cmdRoot := &cobra.Command{
		Use:               "miau",
		Short:             "Control Elgato Key Light",
		Version:           version,
		CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().Lookup(_flagBrightness).Changed || cmd.Flags().Lookup(_flagTemperature).Changed {
				return handleState(-1, cmd.Flags())
			}

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
		Short: "Turn the light ON",
		RunE: func(cmd *cobra.Command, args []string) error {
			return handleState(1, cmd.Flags())
		},
	}
	cmdOff := &cobra.Command{
		Use:   "off",
		Short: "Turn the light OFF",
		RunE: func(cmd *cobra.Command, args []string) error {
			return handleState(0, cmd.Flags())
		},
	}
	cmdBypass := &cobra.Command{
		Use:   "bypass",
		Short: "Control battery bypass",
	}
	cmdBypass.AddCommand(
		&cobra.Command{
			Use:   "toggle",
			Short: "Toggle battery bypass",
			RunE: func(cmd *cobra.Command, args []string) error {
				addr, err := cmd.Flags().GetIP(_flagAddress)
				if err != nil {
					return err
				}

				return NewClient(addr).ToggleBypass()
			},
		})

	cmdServe := &cobra.Command{
		Use:   "serve",
		Short: "Serve a webpage for controlling Elgato Key Light",
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := cmd.Flags().GetIP(_flagAddress)
			if err != nil {
				return err
			}

			port, err := cmd.Flags().GetInt(_flagPort)
			if err != nil {
				return err
			}

			log.Printf("INFO listening on port :%d\n", port)

			return http.ListenAndServe(fmt.Sprint(":", port), NewClient(addr))
		},
	}

	cmdServe.Flags().IntP(_flagPort, "p", 9123, "port to listen on")

	cmdRoot.AddCommand(cmdOn, cmdOff, cmdBypass, cmdServe)
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

	l := &Light{}
	if on == 0 || on == 1 {
		l.On = &on
	}

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

		t, err = KelvinToAPI(t)
		if err != nil {
			return err
		}

		l.Temperature = &t
	}

	result, err := NewClient(addr).UpdateLight(l)
	if err != nil {
		return err
	}

	fmt.Println(result)

	return nil
}
