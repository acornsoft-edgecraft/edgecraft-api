# Installation with Helm

## requirement
- helm 3

## Postgres Installation

- 1. Bitnami 차트 repo 추가:
```sh
helm repo add bitnami https://charts.bitnami.com/bitnami

helm repo upgrade
```

- 2. 명령을 사용하여 Postgres를 배포합니다.
```sh
## // -- values.yaml 수정 내용
## global.postgresql.username="edgecraft"
## global.global.postgresql.password="edgecraft"
## global.postgresql.database="edgecraft"
## global.storageClass=nfs-storageclass
## service.type=NodePort
## service.nodePorts.postgresql="31000"

helm -n edgecraft upgrade postgresql bitnami/postgresql -i --create-namespace --cleanup-on-fail --values values.yaml
```