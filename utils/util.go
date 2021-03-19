package utils

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	apiextentionclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	log "k8s.io/klog/v2"
	//	"kubevirt.io/kubevirt/pkg/virt-controller/services"
	"os"
	"os/exec"
	"time"
)

var (
	volume   = "volume"
	contName = "virt"
	diskDir  = "/disks"
	podName  = "libguestfs-tools"
)

var (
	timeout = 200 * time.Second
)

const KvmDevice = "devices.kubevirt.io/kvm"

type K8sClient struct {
	k8sClient    *kubernetes.Clientset
	apiExtclient *apiextentionclient.Clientset
}

func newK8sClient(config *rest.Config) (*K8sClient, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return &K8sClient{}, err
	}
	apiextclient, err := apiextentionclient.NewForConfig(config)
	if err != nil {
		return &K8sClient{}, err
	}

	return &K8sClient{
		k8sClient:    clientset,
		apiExtclient: apiextclient,
	}, nil
}

func CreateClient(kubeconfig string) (*K8sClient, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return newK8sClient(config)

}

func CreateClientInCluster() (*K8sClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return &K8sClient{}, err
	}

	return newK8sClient(config)
}

func (client *K8sClient) ExistsPVC(pvc, ns string) (bool, error) {
	p, err := client.k8sClient.CoreV1().PersistentVolumeClaims(ns).Get(context.TODO(), pvc, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	if p.Name == "" {
		return false, nil
	}
	return true, nil
}

func (client *K8sClient) ExistsPod(pod, ns string) bool {
	p, err := client.k8sClient.CoreV1().Pods(ns).Get(context.TODO(), pod, metav1.GetOptions{})
	if err != nil {
		return false
	}
	if p.Name == "" {
		return false
	}
	return true
}

// IsPVCinUse returns if the pvc is currently used by a pod
func (client *K8sClient) IsPVCinUse(pvc, ns string) (bool, error) {
	pods, err := client.GetPodsForPVC(pvc, ns)
	if err != nil {
		return false, err
	}
	if len(pods) > 0 {
		return true, nil
	}
	return false, nil
}

func (client *K8sClient) IsKubevirtInstalled() bool {
	c, _ := client.apiExtclient.ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), "kubevirts.kubevirt.io", metav1.GetOptions{})

	if c.ObjectMeta.Name == "kubevirts.kubevirt.io" {
		return true
	}
	return false
}

func (client *K8sClient) waitForContainerRunning(pod, cont, ns string, timeout time.Duration) error {
	log.Infof("Wait for the pod to be started")
	c := make(chan string, 1)
	go func() {
		for {
			pod, err := client.k8sClient.CoreV1().Pods(ns).Get(context.TODO(), pod, metav1.GetOptions{})
			if err != nil {
				c <- err.Error()
			}
			if pod.Status.Phase != corev1.PodPending {
				c <- string(pod.Status.Phase)

			}
			time.Sleep(1 * time.Millisecond)
		}
	}()
	select {
	case res := <-c:
		if res == string(corev1.PodRunning) {
			log.Infof("Pod started")
			return nil
		}
		return fmt.Errorf("Pod is not in running state but got %s", res)
	case <-time.After(timeout):
		return fmt.Errorf("timeout in waiting for the containers to be started in pod %s", pod)
	}

}

func (client *K8sClient) GetPodsForPVC(pvcName, ns string) ([]corev1.Pod, error) {
	nsPods, err := client.k8sClient.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
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

func createPod(pvc, image, cmd string, args []string, kvm bool) *corev1.Pod {
	var resources corev1.ResourceRequirements
	if kvm {
		log.Infof("Run %s pod with KVM enabled", podName)
		resources = corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				KvmDevice: resource.MustParse("1"),
			},
		}
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: podName,
		},
		Spec: corev1.PodSpec{
			Volumes: []corev1.Volume{
				{
					Name: volume,
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: pvc,
							ReadOnly:  false,
						},
					},
				},
			},
			Containers: []corev1.Container{
				{
					Name:       contName,
					Image:      image,
					Command:    []string{cmd},
					Args:       args,
					WorkingDir: diskDir,
					Env: []corev1.EnvVar{
						{
							Name:  "LIBGUESTFS_BACKEND",
							Value: "direct",
						},
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      volume,
							ReadOnly:  false,
							MountPath: diskDir,
						},
					},
					ImagePullPolicy: corev1.PullIfNotPresent,
					Stdin:           true,
					TTY:             true,
					Resources:       resources,
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
}

func (client *K8sClient) CreateInteractivePodWithPVC(config, pvc, image, ns, command string, args []string) error {
	var err error
	kvm := client.IsKubevirtInstalled()
	if !client.ExistsPod(podName, ns) {
		log.Infof("Pod %s doesn't exist. create", podName)
		pod := createPod(pvc, image, command, args, kvm)
		_, err = client.k8sClient.CoreV1().Pods(ns).Create(context.TODO(), pod, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}
	err = client.waitForContainerRunning(podName, contName, ns, timeout)
	if err != nil {
		return err
	}
	os.Setenv("KUBECONFIG", config)
	argsKubectl := []string{
		"attach",
		podName,
		"-ti",
		"-c", contName,
	}
	cmd := exec.Command("kubectl", argsKubectl...)
	log.Infof("Execute: %s", cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (client *K8sClient) RemovePod(ns string) error {
	log.Infof("Remove pod %s", podName)
	return client.k8sClient.CoreV1().Pods(ns).Delete(context.TODO(), podName, metav1.DeleteOptions{})
}
