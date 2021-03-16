package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootPassword string
	packages     string
)

// customizeCmd represents the customize command
var customizeCmd = &cobra.Command{
	Use:   "customize",
	Short: "Customize disk on the pvc",
	RunE: func(cmd *cobra.Command, args []string) error {
		argsCustomize := []string{
			"-a", "disk.img"}
		if rootPassword != "" {
			argsCustomize = append(argsCustomize, "--root-password", rootPassword)
		}

		if packages != "" {
			argsCustomize = append(argsCustomize, "--install", packages)
		}

		return runInteractivePod("virt-customize", argsCustomize)
	},
}

func init() {
	rootCmd.AddCommand(customizeCmd)
	customizeCmd.PersistentFlags().StringVar(&rootPassword, "root-password", "", "Set root password")
	customizeCmd.PersistentFlags().StringVar(&packages, "install", "", "Packages to install <PKG,PKG..>")
}
