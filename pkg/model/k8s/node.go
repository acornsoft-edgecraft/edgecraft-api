/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package k8s

import (
	"strconv"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
)

// Node - Kubernetes Node 정보
type Node struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Role    string `json:"role"`
	Age     string `json:"age"`
	Version string `json:"version"`
}

// isReady - 노드 상태 반환
func isReady(node *v1.Node) string {
	var cond v1.NodeCondition
	for _, c := range node.Status.Conditions {
		if c.Type == v1.NodeReady {
			cond = c
			break
		}
	}

	if cond.Status == v1.ConditionTrue {
		return "Ready"
	} else {
		return "Not ready"
	}
}

// getRoles - 노드의 역할들 반환
func getRoles(node *v1.Node) string {
	var roles []string

	if _, ok := node.Labels["node-role.kubernetes.io/control-plane"]; ok {
		roles = append(roles, "control-plane")
	}
	if _, ok := node.Labels["node-role.kubernetes.io/master"]; ok {
		roles = append(roles, "master")
	}

	if len(roles) == 0 {
		roles = append(roles, "<none>")
	}

	return strings.Join(roles, ",")
}

// getAge - 노드의 생성이후 시간 반환
func getAge(node *v1.Node) string {
	diff := uint64(time.Since(node.CreationTimestamp.Time).Hours())
	days := diff / 24
	if days > 1 {
		return strconv.FormatUint(days, 10) + "d"
	}
	if diff > 1 {
		return strconv.FormatUint(diff, 10) + "h"
	}
	minutes := diff * 60
	if minutes > 1 {
		return strconv.FormatUint(minutes, 10) + "m"
	}

	return strconv.FormatUint(minutes*60, 10) + "s"
}

// getKubeletVersion - 노드의 쿠버네티스 버전 반환
func getKubeletVersion(node *v1.Node) string {
	return node.Status.NodeInfo.KubeletVersion
}

// ConvertToNodeList - Kubernetes NodeList를 화면에서 사용할 수 있는 정보로 전환
func ConvertToNodeList(nodeList *v1.NodeList) ([]Node, error) {
	var nodes []Node

	for _, item := range nodeList.Items {
		node := Node{
			Name:    item.Name,
			Status:  isReady(&item),
			Role:    getRoles(&item),
			Age:     getAge(&item),
			Version: getKubeletVersion(&item),
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}
