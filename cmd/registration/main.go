package main

import (
	goflag "flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	utilflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"

	"github.com/open-cluster-management/registration/pkg/cmd/hub"
	"github.com/open-cluster-management/registration/pkg/cmd/spoke"
	"github.com/open-cluster-management/registration/pkg/cmd/webhook"
	"github.com/open-cluster-management/registration/pkg/version"
)

// The registration binary contains both the hub-side controllers for the
// registration API and the spoke agent.

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	pflag.CommandLine.SetNormalizeFunc(utilflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	logs.InitLogs()
	defer logs.FlushLogs()

	command := newRegistrationCommand()
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func newRegistrationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registration",
		Short: "Spoke Cluster Registration",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
			os.Exit(1)
		},
	}

	if v := version.Get().String(); len(v) == 0 {
		cmd.Version = "<unknown>"
	} else {
		cmd.Version = v
	}

	cmd.AddCommand(hub.NewController())
	cmd.AddCommand(spoke.NewAgent())
	cmd.AddCommand(webhook.NewAdmissionHook())

	return cmd
}
