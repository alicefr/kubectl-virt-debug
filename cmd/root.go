package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	log "k8s.io/klog/v2"
	"os"
)

var (
	PvcClaimName string
	defaultImage = "libguestfs-tools"
	Image        string
	Namespace    string
	Config       string
	Running      bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "guestfs",
	Short: "kubectl plugin to create libguestfs pod",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("Failed: %v", err)
		os.Exit(1)
	}
}

func init() {
	config := os.Getenv("KUBECONFIG")
	if config == "" {
		config = "~/.kube/config"
	}
	rootCmd.PersistentFlags().StringVarP(&PvcClaimName, "pvc", "p", "", "pvc claim name")
	rootCmd.MarkPersistentFlagRequired("pvc")
	rootCmd.PersistentFlags().StringVarP(&Namespace, "ns", "n", "default", "namspace of the pvc")
	rootCmd.PersistentFlags().StringVarP(&Config, "config", "c", config, "path to kubernetes config file")
	rootCmd.PersistentFlags().StringVarP(&Image, "image", "i", defaultImage, fmt.Sprintf("overwrite default container image"))
	rootCmd.PersistentFlags().BoolVar(&Running, "running", false, "let the libguestfs-tool pod running")
}
