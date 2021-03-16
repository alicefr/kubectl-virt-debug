package cmd

import (
	"github.com/spf13/cobra"
)

// sparsifyCmd represents the sparsify command
var sparsifyCmd = &cobra.Command{
	Use:   "sparsify",
	Short: "Sparsify the disk on the pvc",
}

func init() {
	rootCmd.AddCommand(sparsifyCmd)
}
