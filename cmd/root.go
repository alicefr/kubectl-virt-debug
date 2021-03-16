package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	PvcClaimName string
	defaultImage = "libguestfs-tools"
	Image        string
	Namespace    string
	Config       string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "virt-debug",
	Short: "kubectl plugin to create libguestfs pod",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	config := os.Getenv("KUBECONFIG")
	if config == "" {
		config = "~/.kube/config"
	}
	rootCmd.PersistentFlags().StringVarP(&PvcClaimName, "pvc", "p", "", "pvc claim name")
	rootCmd.PersistentFlags().StringVarP(&Namespace, "ns", "n", "default", "namspace of the pvc")
	rootCmd.PersistentFlags().StringVarP(&Config, "config", "c", config, "path to kubernetes config file")
	rootCmd.PersistentFlags().StringVarP(&Image, "image", "i", defaultImage, fmt.Sprintf("overwrite default container image"))
	rootCmd.MarkFlagRequired("pvc")
}
