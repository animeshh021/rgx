package rgx

import (
	"os"
	"rgx/common/log"
	"rgx/common/utils"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     utils.ApplicationName,
	Version: utils.Version,
	Short:   utils.ApplicationName + ":" + utils.ApplicationShortDescription,
	Long:    utils.ApplicationDescription,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error("Error running command: ", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().Bool("debug", false, "Display debug messages (false, by default)")
	rootCmd.PersistentFlags().Bool("trace", false, "Display trace messages (false, by default)")
}

// TODO refactor this
func setDebug(cmd *cobra.Command) {
	trace, _ := cmd.Flags().GetBool("trace")
	if trace {
		log.EnableTrace()
		return
	}
	debug, _ := cmd.Flags().GetBool("debug")
	if debug {
		log.EnableDebug()
	}
}
