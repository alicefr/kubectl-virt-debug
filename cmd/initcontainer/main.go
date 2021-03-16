package main

import (
	"flag"
	"fmt"
	"github.com/alicefr/kubectl-virt-guestfs/utils"
	log "k8s.io/klog/v2"
	"os"
)

func ExitWithError(msg string) {
	log.Errorf(msg)
	os.Exit(1)
}

func ExitWithoutError(msg string) {
	log.Infof(msg)
	os.Exit(0)
}

func main() {
	var pvc, ns string
	flag.StringVar(&pvc, "pvc", "", "pvc to check")
	flag.StringVar(&ns, "ns", "", "namespace where the pvc is located")
	flag.Parse()
	if pvc == "" {
		ExitWithError("pvc cannot be empty")
	}
	if pvc == "" {
		ExitWithError("ns cannot be empty")
	}

	client, err := utils.CreateClientInCluster()
	if err != nil {
		ExitWithError(fmt.Sprintf("Failed create k8s client: %v", err))
	}
	if !client.ExistsPVC(pvc, ns) {
		ExitWithError(fmt.Sprintf("The PVC %s doesn't exist", pvc))
	}
	var inUse bool
	inUse, err = client.IsPVCinUse(pvc, ns)
	if err != nil {
		ExitWithError(fmt.Sprintf("Failed getting the pvc: %v", err))
	}
	if inUse {
		ExitWithError(fmt.Sprintf("PVC %s is in use, and libguestfs-tool cannot be run on the pvc until is in used", pvc))
	}
	ExitWithoutError("PVC can be used")
}
