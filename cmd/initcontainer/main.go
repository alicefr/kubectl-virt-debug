package main

import (
	"flag"
	"fmt"
	"github.com/alicefr/kubectl-virt-guestfs/utils"
	corev1 "k8s.io/api/core/v1"
	log "k8s.io/klog/v2"
	"os"
	"strings"
)

func ExitWithError(msg string) {
	log.Errorf(msg)
	os.Exit(1)
}

func ExitWithoutError(msg string) {
	log.Infof(msg)
	os.Exit(0)
}

type DuplicateFlags struct {
	labels map[string]string
}

func (i *DuplicateFlags) String() string {
	return fmt.Sprintf("%v", i.labels)
}

func (i *DuplicateFlags) Set(value string) error {
	s := strings.Split(value, "=")
	if len(s) != 2 {
		return fmt.Errorf("Wrong format for the label %s. It has to be <label>=<value>", value)
	}
	i.labels[s[0]] = s[1]
	return nil
}

func (i *DuplicateFlags) Contains(labels map[string]string) bool {
	for k, v := range labels {
		if val, ok := i.labels[k]; ok && val == v {
			return true
		}

	}
	return false
}

func contains(s []string, v string) bool {
	for _, t := range s {
		if t == v {
			return true
		}
	}
	return false
}

func main() {
	var pvc, ns, pods string
	l := DuplicateFlags{
		labels: make(map[string]string),
	}
	flag.StringVar(&pvc, "pvc", "", "pvc to check")
	flag.StringVar(&ns, "ns", "", "namespace where the pvc is located")
	flag.StringVar(&pods, "skip-pods", "", "list of pods to skip the check <pod1,pod2,..>")
	flag.Var(&l, "label", "label for the pods to be skipped <label>=<value>. This option can be specified multiple times")
	flag.Parse()

	log.Infof("label %s", l)
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
	var exist bool
	exist, err = client.ExistsPVC(pvc, ns)
	if err != nil {
		ExitWithError(fmt.Sprintf("Failed to get the pvc %s in the namespace %s: %v)", pvc, ns, err))
	}
	if !exist {
		ExitWithError(fmt.Sprintf("The PVC %s doesn't exist", pvc))
	}
	var podsUsingPvc []corev1.Pod
	podsUsingPvc, err = client.GetPodsForPVC(pvc, ns)
	if err != nil {
		ExitWithError(fmt.Sprintf("Failed getting the pvc: %v", err))
	}

	if len(podsUsingPvc) > 0 {
		var podNames []string
		// Inside the cluster the pvc is in use by the pod its self. Hence, we need a way to skip certain pods
		for _, p := range podsUsingPvc {
			if l.Contains(p.ObjectMeta.Labels) {
				continue
			}
			if contains(strings.Split(pods, ","), p.Name) {
				continue
			}
			log.Errorf("Pod %s is using the PVC %s", p.Name, pvc)
			podNames = append(podNames, p.Name)
		}
		if len(podNames) > 0 {
			ExitWithError(fmt.Sprintf("PVC %s is in use by the pods %v", pvc, podNames))
		}
	}
	ExitWithoutError("PVC can be used")
}
