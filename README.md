Kubectl plugin for running libguestfs-tool with kubevirt
---
This repo contains a kubectl plugin to use libguestfs tools inside a kubernetes cluster.
```bash
kubectl plugin to create libguestfs pod

Usage:
  guestfs [command]

Available Commands:
  check       Check if the pvc is in use or not
  customize   Customize disk on the pvc
  help        Help about any command
  rescue      Run virt-rescue
  shell       Start a shell to the libguestfs pod
  sparsify    Sparsify the disk on the pvc

Flags:
  -c, --config string   path to kubernetes config file (default "/home/afrosi/go/src/github.com/kubevirt/kubevirt/_ci-configs/k8s-1.18/.kubeconfig")
  -h, --help            help for guestfs
  -i, --image string    overwrite default container image (default "libguestfs-tools")
  -n, --ns string       namspace of the pvc (default "default")
  -p, --pvc string      pvc claim name
      --running         let the libguestfs-tool pod running

Use "guestfs [command] --help" for more information about a command.

```
