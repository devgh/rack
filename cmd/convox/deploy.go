package main

import (
	"fmt"
	"os"

	"gopkg.in/urfave/cli.v1"

	"github.com/convox/rack/cmd/convox/stdcli"
)

func init() {
	stdcli.RegisterCommand(cli.Command{
		Name:        "deploy",
		Description: "deploy an app to AWS",
		Usage:       "<directory>",
		Action:      cmdDeploy,
		Flags: append(
			buildCreateFlags,
			cli.BoolFlag{
				Name:   "wait",
				EnvVar: "CONVOX_WAIT",
				Usage:  "wait for release to finish promoting before returning",
			},
		),
	})
}

func cmdDeploy(c *cli.Context) error {
	wd := "."

	if len(c.Args()) > 0 {
		wd = c.Args()[0]
	}

	dir, app, err := stdcli.DirApp(c, wd)
	if err != nil {
		return stdcli.Error(err)
	}

	output := os.Stdout

	if c.Bool("id") {
		output = os.Stderr
	}

	output.Write([]byte(fmt.Sprintf("Deploying %s\n", app)))

	a, err := rackClient(c).GetApp(app)
	if err != nil {
		return stdcli.Error(err)
	}

	switch a.Status {
	case "creating":
		return stdcli.Error(fmt.Errorf("app is still creating: %s", app))
	case "running", "updating":
	default:
		return stdcli.Error(fmt.Errorf("unable to build app: %s", app))
	}

	// build
	_, release, err := executeBuild(c, dir, app, c.String("file"), c.String("description"), output)
	if err != nil {
		return stdcli.Error(err)
	}

	if release == "" {
		return nil
	}

	output.Write([]byte(fmt.Sprintf("Release: %s\n", release)))

	if c.Bool("id") {
		os.Stdout.Write([]byte(release))
		output.Write([]byte("\n"))
	}

	output.Write([]byte(fmt.Sprintf("Promoting %s... ", release)))

	_, err = rackClient(c).PromoteRelease(app, release)
	if err != nil {
		return stdcli.Error(err)
	}

	output.Write([]byte("UPDATING\n"))

	if c.Bool("wait") {
		output.Write([]byte(fmt.Sprintf("Waiting for %s... ", release)))

		if err := waitForReleasePromotion(c, app, release); err != nil {
			return stdcli.Error(err)
		}

		output.Write([]byte("OK\n"))
	}

	return nil
}
