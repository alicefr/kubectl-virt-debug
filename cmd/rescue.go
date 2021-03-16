package cmd

import (
	"github.com/alicefr/kubectl-virt-debug/utils"
	"github.com/spf13/cobra"
	log "k8s.io/klog/v2"
	"os"
)

// rescueCmd represents the rescue command
var rescueCmd = &cobra.Command{
	Use:   "rescue",
	Short: "Run virt-rescue",
	RunE: func(cmd *cobra.Command, args []string) error {
		var inUse bool
		client, err := utils.CreateClient(Config)
		if err != nil {
			return err
		}
		if !client.ExistsPVC(PvcClaimName, Namespace) {
			log.Infof("The PVC %s doesn't exist", PvcClaimName)
			os.Exit(1)
		}
		inUse, err = client.IsPVCinUse(PvcClaimName, Namespace)
		if err != nil {
			return err
		}
		if inUse {
			log.Infof("PVC %s is in use, and virt-rescue cannot be run on the pvc until is in used", PvcClaimName)
			os.Exit(0)
		}
		defer client.RemovePod(Namespace)
		log.Infof("Attach to libguestfs pod")
		argsRescue := []string{"-a", "disk.img"}
		err = client.CreateInteractivePodWithPVC(Config, PvcClaimName, Image, Namespace, "virt-rescue", argsRescue)
		return err
	},
}

func init() {
	rootCmd.AddCommand(rescueCmd)
}
