package utils

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	*kubernetes.Clientset
}

func CreateClient(kubeconfig string) (*K8sClient, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &K8sClient{client}, nil

}

func (client *K8sClient) ExistsPVC(pvc, ns string) bool {
	p, err := client.CoreV1().PersistentVolumeClaims(ns).Get(context.TODO(), pvc, metav1.GetOptions{})
	if err != nil {
		return false
	}
	fmt.Printf("XXX %v \n", p)
	if p.Name == "" {
		return false
	}
	return true
}

// IsPVCinUse returns if the pvc is currently used by a pod
func (client *K8sClient) IsPVCinUse(pvc, ns string) (bool, error) {
	pods, err := client.getPodsForPVC(pvc, ns)
	if err != nil {
		return false, err
	}
	if len(pods) > 0 {
		return true, nil
	}
	return false, nil
}

func (client *K8sClient) getPodsForPVC(pvcName, ns string) ([]corev1.Pod, error) {
	nsPods, err := client.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return []corev1.Pod{}, err
	}

	var pods []corev1.Pod

	for _, pod := range nsPods.Items {
		for _, volume := range pod.Spec.Volumes {
			if volume.VolumeSource.PersistentVolumeClaim != nil && volume.VolumeSource.PersistentVolumeClaim.ClaimName == pvcName {
				pods = append(pods, pod)
			}
		}
	}

	return pods, nil
}
