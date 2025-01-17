package rgx

import (
	"fmt"
	"os"
	"rgx/candidates"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "display details about packages available on the server",
}

var serverListCmd = &cobra.Command{
	Use:   "list",
	Short: "list packages available on the server",
	Run:   list,
}

var serverShowCmd = &cobra.Command{
	Use:   "show",
	Short: "show versions of the package available on the server",
	Run:   show,
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(serverListCmd)
	serverCmd.AddCommand(serverShowCmd)
	serverShowCmd.Flags().BoolP("lts", "", false, "only consider LTS releases")
}

func list(cmd *cobra.Command, _ []string) {
	setDebug(cmd)
	candidates.PrintServerPackages()
}

func show(cmd *cobra.Command, args []string) {
	var usage = `Usage: rgx server show <package> [options]
e.g.
	rgx server show openjdk
	rgx server show openjdk --lts`

	setDebug(cmd)
	lts, _ := cmd.Flags().GetBool("lts")
	if len(args) != 1 {
		fmt.Println(usage)
		os.Exit(1)
	}

	candidates.PrintMajorVersions(args[0], lts)
}
