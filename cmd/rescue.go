/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/alicefr/kubectl-virt-debug/utils"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
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
			klog.Infof("The PVC %s doesn't exist", PvcClaimName)
			os.Exit(1)
		}
		inUse, err = client.IsPVCinUse(PvcClaimName, Namespace)
		if err != nil {
			return err
		}
		if inUse {
			klog.Infof("PVC %s is in use, and virt-rescue cannot be run on the pvc until is in used", PvcClaimName)
			os.Exit(0)
		}
		defer client.RemovePod(Namespace)
		klog.Infof("Attach to libguestfs pod")
		argsRescue := []string{"-a", "disk.img"}
		err = client.CreateInteractivePodWithPVC(Config, PvcClaimName, Image, Namespace, "virt-rescue", argsRescue)
		return err
	},
}

func init() {
	rootCmd.AddCommand(rescueCmd)
}
