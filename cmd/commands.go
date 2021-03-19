package cmd

import (
	"fmt"
	"github.com/alicefr/kubectl-virt-guestfs/utils"
	"github.com/spf13/cobra"
	log "k8s.io/klog/v2"
	"os"
)

var (
	interactive bool
	runPod      = runInteractivePod
)

func createPod(command string, args []string) error {
	return fmt.Errorf("Method not implemented")
}

func runInteractivePod(command string, args []string) error {
	var inUse bool
	client, err := utils.CreateClient(Config)
	if err != nil {
		return err
	}
	exist, _ := client.ExistsPVC(PvcClaimName, Namespace)
	if !exist {
		log.Infof("The PVC %s doesn't exist", PvcClaimName)
		os.Exit(1)
	}
	inUse, err = client.IsPVCinUse(PvcClaimName, Namespace)
	if err != nil {
		return err
	}
	if inUse {
		log.Infof("PVC %s is in use by another pod", PvcClaimName)
		os.Exit(0)
	}
	if !Running {
		defer client.RemovePod(Namespace)
	}
	return client.CreateInteractivePodWithPVC(Config, PvcClaimName, Image, Namespace, command, args)
}

// rescueCmd represents the rescue command
var rescueCmd = &cobra.Command{
	Use:   "rescue",
	Short: "Run virt-rescue",
	RunE: func(cmd *cobra.Command, args []string) error {
		argsRescue := []string{"-a", "disk.img"}
		return runInteractivePod("virt-rescue", argsRescue)
	},
}

// sparsifyCmd represents the sparsify command
var sparsifyCmd = &cobra.Command{
	Use:   "sparsify",
	Short: "Sparsify the disk on the pvc",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !interactive {
			runPod = createPod
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		argsSparsify := []string{"--in-place", "disk.img"}
		return runPod("virt-sparsify", argsSparsify)
	},
}

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Start a shell to the libguestfs pod",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInteractivePod("/bin/bash", []string{})
	},
}

func init() {
	rootCmd.AddCommand(rescueCmd)
	rootCmd.AddCommand(sparsifyCmd)
	rootCmd.AddCommand(shellCmd)
	// TODO transform it into job
	//	sparsifyCmd.PersistentFlags().BoolVar(&interactive, "interactive", true, "Pod run in detached mode")
}
