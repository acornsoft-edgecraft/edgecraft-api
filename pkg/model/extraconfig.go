/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// ExtraConfig - Extra configuration for kubeadm
type ExtraConfig struct {
	PreKubeadmCommands  string `json:"pre_kubeadm_commands" example:""`
	PostKubeadmCommands string `json:"post_kubeadm_commands" example:"kubectl --kubeconfig=/etc/kubernetes/admin.conf apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.24.1/manifests/calico.yaml"`
	Files               string `json:"files" example:""`
	Users               string `json:"users" example:""`
	Ntp                 string `json:"ntp" example:""`
	Format              string `json:"format" example:""`
}

// ToTable - ExtraConfig 정보를 테이블 정보로 설정
func (ec *ExtraConfig) ToTable(config *ExtraConfig) {
	config.PreKubeadmCommands = ec.PreKubeadmCommands
	config.PostKubeadmCommands = ec.PostKubeadmCommands
	config.Files = ec.Files
	config.Users = ec.Users
	config.Ntp = ec.Ntp
	config.Format = ec.Format
}

// FromTable - 테이블의 정보를 ExtraConfig 정보로 설정
func (ec *ExtraConfig) FromTable(config *ExtraConfig) {
	ec.PreKubeadmCommands = config.PreKubeadmCommands
	ec.PostKubeadmCommands = config.PostKubeadmCommands
	ec.Files = config.Files
	ec.Users = config.Users
	ec.Ntp = config.Ntp
	ec.Format = config.Format
}

// Value Marshal
func (a ExtraConfig) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *ExtraConfig) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
