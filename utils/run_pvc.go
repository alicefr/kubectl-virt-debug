package utils

import (
	"encoding/json"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"os"
	"os/exec"
)

const (
	volume   = "volume"
	contName = "virt"
	diskDir  = "/disks"
	podName  = "libguestfs-tools"
)

func createPodSpec(pvc, image, cmd string) (string, error) {
	spec := corev1.PodSpec{
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
				WorkingDir: diskDir,
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
			},
		},
		RestartPolicy: corev1.RestartPolicyNever,
	}

	j, err := json.Marshal(spec)
	if err != nil {
		return "", err
	}

	return string(j), nil
}

func CreateInteractivePodWithPVC(config, pvc, image string) error {
	command := "bash"
	os.Setenv("KUBECONFIG", config)
	args := []string{
		"run",
		podName,
		"-ti",
		"--restart=Never",
		"--rm",
		fmt.Sprintf("--image=%s", image),
	}
	o, err := createPodSpec(pvc, image, command)
	if err != nil {
		return err
	}
	args = append(args, fmt.Sprintf("--overrides='%s'", o))
	cmd := exec.Command("kubectl", args...)
	klog.Infof("Execute: %s", cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
