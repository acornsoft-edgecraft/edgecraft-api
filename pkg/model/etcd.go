/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
)

type Endpoints []UrlInfo

// Value Marshal
func (a Endpoints) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *Endpoints) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

// EtcdInfo - Data for ETCD
type EtcdInfo struct {
	UseExternalEtcd bool       `json:"use_external_etcd" example:"false"`
	Endpoints       *Endpoints `json:"endpoints"`
	CAFile          string     `json:"ca_file" example:""`
	CertFile        string     `json:"cert_file" example:""`
	KeyFile         string     `json:"key_file" example:""`
}

// ToTable - ETCD 정보를 테이블로 설정
func (ei *EtcdInfo) ToTable(clusterTable *ClusterTable) {
	clusterTable.ExternalEtcdUse = utils.BoolPtr(ei.UseExternalEtcd)
	clusterTable.ExternalEtcdCertificateCa = utils.StringPtr(ei.CAFile)
	clusterTable.ExternalEtcdCertificateCert = utils.StringPtr(ei.CertFile)
	clusterTable.ExternalEtcdCertificateKey = utils.StringPtr(ei.KeyFile)
	clusterTable.ExternalEtcdEndPoints = ei.Endpoints
}

// FromTable - 테이블 정보를 ETCD 정보로 설정
func (ei *EtcdInfo) FromTable(clusterTable *ClusterTable) {
	ei.UseExternalEtcd = *clusterTable.ExternalEtcdUse
	ei.CAFile = *clusterTable.ExternalEtcdCertificateCa
	ei.CertFile = *clusterTable.ExternalEtcdCertificateCert
	ei.KeyFile = *clusterTable.ExternalEtcdCertificateKey
	ei.Endpoints = clusterTable.ExternalEtcdEndPoints
}

// ExternalEtcdEndPoints - Data from ETCD Endpoints
type ExternalEtcdEndPoints []interface{}

// Value Marshal
func (a ExternalEtcdEndPoints) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *ExternalEtcdEndPoints) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

// TODO: 화면에서는 URLInfo 배열을 DB에서는 ExternalEtcdEndPoints를 사용
