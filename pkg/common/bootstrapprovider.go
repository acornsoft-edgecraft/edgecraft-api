/*
Copyright 2023 Acronsoft Authors. All right reserved.
*/
package common

// ===== [ Constants and Variables ] =====

// ===== [ BootstrapProviders ] =====

type BootstrapProvider int

const (
	Kubeadm  BootstrapProvider = 1
	MicroK8s BootstrapProvider = 2
	K3s      BootstrapProvider = 3
)
