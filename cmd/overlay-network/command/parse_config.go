package command

import (
	"flag"
	"fmt"
	"strings"

	flags "github.com/da-moon/coe865-final/cmd/overlay-network/flags"
	config "github.com/da-moon/coe865-final/pkg/config"
	cli "github.com/mitchellh/cli"
)

// ParseConfigCommand is a Command implementation that generates an encryption
// key.
type ParseConfigCommand struct {
	Ui cli.Ui
}

var _ cli.Command = &ParseConfigCommand{}

const entrypoint = "parse-config"

// Run ...
func (c *ParseConfigCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet(entrypoint, flag.ContinueOnError)
	cmdConfigFactory := config.DefaultConfigFactory()
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	var configFiles []string
	cmdFlags.Var((*flags.AppendSliceValue)(&configFiles), "config-file",
		"raw file to read config from")
	cmdFlags.Var((*flags.AppendSliceValue)(&configFiles), "config-dir",
		"directory of raw config files to read")
	dev := flags.DevFlag(cmdFlags)
	logLevel := flags.LogLevelFlag(cmdFlags)
	port := flags.RPCPortFlag(cmdFlags)
	cron := flags.CronFlag(cmdFlags)
	err := cmdFlags.Parse(args)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("[ERROR] could not parse arguments : %v", err))
		return -1
	}
	cmdConfigFactory.DevelopmentMode = *dev
	cmdConfigFactory.LogLevel = *logLevel
	cmdConfigFactory.Port = *port
	cmdConfigFactory.Cron = *cron
	factory := config.DefaultConfigFactory()
	factory = config.MergeFactory(factory, cmdConfigFactory)
	if len(configFiles) > 0 {
		mapping, err := factory.ReadConfigPaths(configFiles, config.CONF)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("[ERROR]: %s", err.Error()))
			return -1
		}
		for k, v := range mapping {
			err := v.SaveAsJSON(k)
			if err != nil {
				c.Ui.Error(fmt.Sprintf("[ERROR]: %s", err.Error()))
				return -1
			}
		}
	}
	return 0
}

// Synopsis ...
func (c *ParseConfigCommand) Synopsis() string {

	return "transform a given config file to sane format"
}

// Help ...
func (c *ParseConfigCommand) Help() string {

	helpText := `
Usage: overlay-network parse-config
  reads a config file as defined in project specification and
  converts is into a normal mashalling format such as JSON.
  it stores the converted file with the same name.
    -config-file=foo            Path to a config file you wish to convert.
                                This can be specified multiple times.
	-config-dir=foo             Path to a directory to read and convert configurations from.
	                            This can be specified multiple times.
	-rpc-port=8080 	            Override converted settings files rpc port value.
    -dev                        Override converted settings files development mode to true
    -log-level=info             Override converted settings files log level.
    -cost-estimator-path=foo    Override converted settings files cost estimator plugin Path.
    -cron='@every 20s'          Override converted settings files message sending interval.
  sample config file (before transform) :
  1 100 10.2.2.1	; RCID ASN IP Address (local rc info)
  2	                ; No. of RC connected
  2 200 10.1.1.2	; RCID ASN IP Address
  3 300 11.1.1.2	; RCID ASN IP Address
  4	                ; No. of ASN connected
  10 2 5 	        ; ASN Mbps(link capacity) cost
  20 5 5	        ; ASN Mbps(link capacity) cost
  200 10 5          ; ASN Mbps(link capacity) cost
  300 10 5          ; ASN Mbps(link capacity) cost
`
	return strings.TrimSpace(helpText)
}
