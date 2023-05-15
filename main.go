package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

			timeout, err := cmd.Flags().GetDuration("timeout")
			if err != nil {
				return err
			}

			l, err := NewClient(addr, timeout).GetLight()
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

	cmdSettings := &cobra.Command{
		Use:   "settings",
		Short: "Update light's settings",
		Long:  "Running settings command without any flags will fetch and print the current settings.",
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := cmd.Flags().GetIP(_flagAddress)
			if err != nil {
				return err
			}

			timeout, err := cmd.Flags().GetDuration("timeout")
			if err != nil {
				return err
			}

			c := NewClient(addr, timeout)
			if !cmd.Flags().Lookup("bypass").Changed {
				s, err := c.GetSettings()
				if err != nil {
					return err
				}
				fmt.Println(s)
			}

			bypass, err := cmd.Flags().GetInt("bypass")
			if err != nil {
				return err
			}
			if bypass != 0 && bypass != 1 {
				return errors.New("invalid value: bypass must be 0 or 1")
			}
			return c.SetBypass(bypass)
		},
	}

	cmdSettings.Flags().Int("bypass", 0, "configure the light's battery bypass (0 or 1)")

	cmdServe := &cobra.Command{
		Use:   "serve",
		Short: "Serve a webpage for controlling Elgato Key Light",
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := cmd.Flags().GetIP(_flagAddress)
			if err != nil {
				return err
			}

			timeout, err := cmd.Flags().GetDuration("timeout")
			if err != nil {
				return err
			}

			port, err := cmd.Flags().GetInt(_flagPort)
			if err != nil {
				return err
			}

			log.Printf("INFO listening on port :%d\n", port)

			return http.ListenAndServe(fmt.Sprint(":", port), NewClient(addr, timeout))
		},
	}

	cmdServe.Flags().IntP(_flagPort, "p", 9123, "port to listen on")

	cmdRoot.AddCommand(cmdOn, cmdOff, cmdServe, cmdSettings)
	cmdRoot.PersistentFlags().IPP(_flagAddress, "a", nil, "IP address of the light")
	cmdRoot.PersistentFlags().IntP(_flagBrightness, "b", 0, "brightness in percent; a value between 0 and 100")
	cmdRoot.PersistentFlags().IntP(_flagTemperature, "t", 0, "temperature in Kelvins; a value between 2900 and 7000")
	cmdRoot.PersistentFlags().Duration("timeout", 5*time.Second, "light API timeout duration")

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

	timeout, err := flags.GetDuration("timeout")
	if err != nil {
		return err
	}

	result, err := NewClient(addr, timeout).UpdateLight(l)
	if err != nil {
		return err
	}

	fmt.Println(result)

	return nil
}
