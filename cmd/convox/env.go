package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/convox/rack/cmd/convox/stdcli"
	"gopkg.in/urfave/cli.v1"
)

func init() {
	stdcli.RegisterCommand(cli.Command{
		Name:        "env",
		Description: "manage an app's environment variables",
		Usage:       "",
		Action:      cmdEnvList,
		Flags:       []cli.Flag{appFlag, rackFlag},
		Subcommands: []cli.Command{
			{
				Name:        "get",
				Description: "get all environment variables",
				Usage:       "VARIABLE",
				Action:      cmdEnvGet,
				Flags:       []cli.Flag{appFlag, rackFlag},
			},
			{
				Name:        "set",
				Description: "set an environment variable",
				Usage:       "VARIABLE=VALUE",
				Action:      cmdEnvSet,
				Flags: []cli.Flag{
					appFlag,
					rackFlag,
					cli.BoolFlag{
						Name:  "promote",
						Usage: "promote the release after env change",
					},
					cli.BoolFlag{
						Name:  "id",
						Usage: "env set logs on stderr, release id on stdout",
					},
					cli.BoolFlag{
						Name:   "wait",
						EnvVar: "CONVOX_WAIT",
						Usage:  "wait for release to finish promoting before returning",
					},
				},
			},
			{
				Name:        "unset",
				Description: "delete an environment varible",
				Usage:       "VARIABLE",
				Action:      cmdEnvUnset,
				Flags: []cli.Flag{
					appFlag,
					rackFlag,
					cli.BoolFlag{
						Name:  "promote",
						Usage: "promote the release after env change",
					},
					cli.BoolFlag{
						Name:  "id",
						Usage: "env set logs on stderr, release id on stdout",
					},
					cli.BoolFlag{
						Name:   "wait",
						EnvVar: "CONVOX_WAIT",
						Usage:  "wait for release to finish promoting before returning",
					},
				},
			},
		},
	})
}

func cmdEnvList(c *cli.Context) error {
	_, app, err := stdcli.DirApp(c, ".")
	if err != nil {
		return stdcli.Error(err)
	}

	if len(c.Args()) > 0 {
		return stdcli.Error(fmt.Errorf("`convox env` does not take arguments. Perhaps you meant `convox env set`?"))
	}

	if c.Bool("help") {
		stdcli.Usage(c, "")
		return nil
	}

	env, err := rackClient(c).GetEnvironment(app)
	if err != nil {
		return stdcli.Error(err)
	}

	keys := []string{}

	for key := range env {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		fmt.Printf("%s=%s\n", key, env[key])
	}

	return nil
}

func cmdEnvGet(c *cli.Context) error {
	_, app, err := stdcli.DirApp(c, ".")
	if err != nil {
		return stdcli.Error(err)
	}

	if len(c.Args()) == 0 {
		return stdcli.Error(errors.New("No variable specified"))
	}

	if len(c.Args()) > 1 {
		return stdcli.Error(errors.New("Only 1 variable can be retrieved at a time"))
	}

	variable := c.Args()[0]

	env, err := rackClient(c).GetEnvironment(app)
	if err != nil {
		return stdcli.Error(err)
	}

	fmt.Println(env[variable])
	return nil
}

func cmdEnvSet(c *cli.Context) error {
	_, app, err := stdcli.DirApp(c, ".")
	if err != nil {
		return stdcli.Error(err)
	}

	// get current app environment
	env, err := rackClient(c).GetEnvironment(app)
	if err != nil {
		return stdcli.Error(err)
	}

	data := ""

	for key, value := range env {
		data += fmt.Sprintf("%s=%s\n", key, value)
	}

	// handle data from stdin
	if !stdcli.IsTerminal(os.Stdin) {
		in, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return stdcli.Error(err)
		}

		scanner := bufio.NewScanner(bytes.NewReader(in))
		for scanner.Scan() {
			data += fmt.Sprintf("%s\n", scanner.Text())
		}
	}

	// handle args
	for _, value := range c.Args() {
		data += fmt.Sprintf("%s\n", value)
	}

	output := os.Stdout
	if c.Bool("id") {
		output = os.Stderr
	}

	output.Write([]byte("Updating environment... "))

	_, release, err := rackClient(c).SetEnvironment(app, strings.NewReader(data))
	if err != nil {
		return stdcli.Error(err)
	}

	output.Write([]byte("OK\n"))

	if release != "" {
		if c.Bool("id") {
			os.Stdout.Write([]byte(fmt.Sprintf("%s\n", release)))
		}

		if c.Bool("promote") {
			output.Write([]byte(fmt.Sprintf("Promoting %s... ", release)))

			_, err = rackClient(c).PromoteRelease(app, release)
			if err != nil {
				return stdcli.Error(err)
			}

			output.Write([]byte("OK\n"))

			if c.Bool("wait") {
				output.Write([]byte(fmt.Sprintf("Waiting for stabilization... ")))

				if err := waitForReleasePromotion(c, app, release); err != nil {
					return stdcli.Error(err)
				}

				output.Write([]byte("OK\n"))
			}
		} else {
			output.Write([]byte(fmt.Sprintf("To deploy these changes run `convox releases promote %s`\n", release)))
		}
	}

	return nil
}

func cmdEnvUnset(c *cli.Context) error {
	_, app, err := stdcli.DirApp(c, ".")
	if err != nil {
		return stdcli.Error(err)
	}

	if len(c.Args()) == 0 {
		return stdcli.Error(errors.New("No variable specified"))
	}

	if len(c.Args()) > 1 {
		return stdcli.Error(errors.New("Only 1 variable can be unset at a time"))
	}

	key := c.Args()[0]

	output := os.Stdout
	if c.Bool("id") {
		output = os.Stderr
	}

	output.Write([]byte("Updating environment... "))

	_, release, err := rackClient(c).DeleteEnvironment(app, key)
	if err != nil {
		return stdcli.Error(err)
	}

	output.Write([]byte("OK\n"))

	if release != "" {
		if c.Bool("id") {
			os.Stdout.Write([]byte(fmt.Sprintf("%s\n", release)))
		}

		if c.Bool("promote") {
			output.Write([]byte(fmt.Sprintf("Promoting %s... ", release)))

			_, err = rackClient(c).PromoteRelease(app, release)
			if err != nil {
				return stdcli.Error(err)
			}

			output.Write([]byte("OK\n"))

			if c.Bool("wait") {
				output.Write([]byte(fmt.Sprintf("Waiting for stabilization... ")))

				if err := waitForReleasePromotion(c, app, release); err != nil {
					return stdcli.Error(err)
				}

				output.Write([]byte("OK\n"))
			}
		} else {
			output.Write([]byte(fmt.Sprintf("To deploy these changes run `convox releases promote %s`\n", release)))
		}
	}

	return nil
}
