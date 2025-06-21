package main

import (
	cc "github.com/rafa-mori/selfrestart/cmd/cli"
	gl "github.com/rafa-mori/selfrestart/logger"
	vs "github.com/rafa-mori/selfrestart/version"
	"github.com/spf13/cobra"

	"os"
	"strings"
)

type SelfRestart struct {
	parentCmdName string
	printBanner   bool
}

func (m *SelfRestart) Alias() string {
	return ""
}
func (m *SelfRestart) ShortDescription() string {
	return "SelfRestart is a Go library for automatic process restart functionality."
}
func (m *SelfRestart) LongDescription() string {
	return `SelfRestart: A Go library that allows applications to restart themselves automatically in a safe and elegant way.`
}
func (m *SelfRestart) Usage() string {
	return "selfrestart [command] [args]"
}
func (m *SelfRestart) Examples() []string {
	return []string{
		"selfrestart start --debug --daemon",
		"selfrestart restart --wait",
		"selfrestart restart --pid 12345",
		"selfrestart status --pid 12345",
		"selfrestart check",
	}
}
func (m *SelfRestart) Active() bool {
	return true
}
func (m *SelfRestart) Module() string {
	return "selfrestart"
}
func (m *SelfRestart) Execute() error {
	return m.Command().Execute()
}
func (m *SelfRestart) Command() *cobra.Command {
	gl.Log("debug", "Starting SelfRestart CLI...")

	var rtCmd = &cobra.Command{
		Use:     m.Module(),
		Aliases: []string{m.Alias()},
		Example: m.concatenateExamples(),
		Version: vs.GetVersion(),
		Annotations: cc.GetDescriptions([]string{
			m.LongDescription(),
			m.ShortDescription(),
		}, m.printBanner),
	}

	rtCmd.AddCommand(cc.ServiceCmdList()...)
	rtCmd.AddCommand(vs.CliCommand())

	// Set usage definitions for the command and its subcommands
	setUsageDefinition(rtCmd)
	for _, c := range rtCmd.Commands() {
		setUsageDefinition(c)
		if !strings.Contains(strings.Join(os.Args, " "), c.Use) {
			if c.Short == "" {
				c.Short = c.Annotations["description"]
			}
		}
	}

	return rtCmd
}
func (m *SelfRestart) SetParentCmdName(rtCmd string) {
	m.parentCmdName = rtCmd
}
func (m *SelfRestart) concatenateExamples() string {
	examples := ""
	rtCmd := m.parentCmdName
	if rtCmd != "" {
		rtCmd = rtCmd + " "
	}
	for _, example := range m.Examples() {
		examples += rtCmd + example + "\n  "
	}
	return examples
}
func RegX() *SelfRestart {
	var printBannerV = os.Getenv("SELFRESTART_PRINT_BANNER")
	if printBannerV == "" {
		printBannerV = "true"
	}

	return &SelfRestart{
		printBanner: strings.ToLower(printBannerV) == "true",
	}
}
