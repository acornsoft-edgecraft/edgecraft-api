# Install CSI driver with Helm 3

> [참고] https://github.com/kubernetes-csi/csi-driver-nfs/blob/master/charts/README.md


## install a specific version

```sh
## step-1. helm repo 추가
$ helm repo add csi-driver-nfs https://raw.githubusercontent.com/kubernetes-csi/csi-driver-nfs/master/charts

## 노드에서 kubelet root_dir 확인
$ ps -ef | grep kubelet | grep 'root-dir' | grep -Po '\-\-root\-dir=\K[^\s]+'
/data/kubelet

## chart value 구성 및 배포
$ helm install csi-driver-nfs csi-driver-nfs/csi-driver-nfs --namespace kube-system --version v4.3.0 \
--set kubeletDir="/data/kubelet" \
--set controller.runOnControlPlane=true

## StorageClass 구성및 배포
kubectl apply -f - <<EOF
---
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: nfs-csi
provisioner: nfs.csi.k8s.io
parameters:
  server: 192.168.88.11
  share: /data/hdd/nfs-storage
  # csi.storage.k8s.io/provisioner-secret is only needed for providing mountOptions in DeleteVolume
  # csi.storage.k8s.io/provisioner-secret-name: "mount-options"
  # csi.storage.k8s.io/provisioner-secret-namespace: "default"
reclaimPolicy: Delete
volumeBindingMode: Immediate
mountOptions:
  - nfsvers=4.1
EOF
```