package rgx

import (
	"fmt"
	"os"
	"rgx/candidates"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install a package",
	Run:   install,
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolP("lts", "", false, "only consider LTS releases")
}

func install(cmd *cobra.Command, args []string) {
	var usage = `Usage: rgx install openjdk <version>
e.g.
	rgx install openjdk latest
	rgx install openjdk latest --lts
	rgx install openjdk 20`

	setDebug(cmd)
	lts, _ := cmd.Flags().GetBool("lts")
	if len(args) != 2 {
		fmt.Println(usage)
		os.Exit(1)
	}
	candidates.Install(args[0], args[1], lts)
}
