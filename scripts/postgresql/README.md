# Installation with Helm

## requirement
- helm 3

## Postgres Installation

- 1. Bitnami 차트 repo 추가:
```sh
helm repo add bitnami https://charts.bitnami.com/bitnami

helm repo update
```

- 2. 명령을 사용하여 Postgres를 배포합니다.
```sh
## // -- values.yaml 수정 내용
# global.postgresql.auth.username="edgecraft"
# global.global.postgresql.auth.password="edgecraft"
# global.postgresql.auth.database="edgecraft"
# global.storageClass=nfs-csi
# primary.service.type="NodePort"
# primary.service.nodePorts.postgresql="31000"
# volumePermissions.enabled=true

helm -n edgecraft upgrade postgresql bitnami/postgresql -i --create-namespace --cleanup-on-fail --values values.yaml

## 또는 --set 사용배포
helm -n edgecraft upgrade postgresql bitnami/postgresql -i --create-namespace --cleanup-on-fail \
--set global.postgresql.auth.username="edgecraft" \
--set global.global.postgresql.auth.password="edgecraft" \
--set global.postgresql.auth.database="edgecraft" \
--set global.storageClass=nfs-csi \
--set primary.service.type="NodePort" \
--set primary.service.nodePorts.postgresql="31000" \
--set volumePermissions.enabled=true
```

- 3. 접속 정보
```sh
NAME: postgresql
LAST DEPLOYED: Fri Jul  7 17:11:33 2023
NAMESPACE: edgecraft
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
CHART NAME: postgresql
CHART VERSION: 12.6.3
APP VERSION: 15.3.0

** Please be patient while the chart is being deployed **

PostgreSQL can be accessed via port 5432 on the following DNS names from within your cluster:

    postgresql.edgecraft.svc.cluster.local - Read/Write connection

To get the password for "postgres" run:

    export POSTGRES_ADMIN_PASSWORD=$(kubectl get secret --namespace edgecraft postgresql -o jsonpath="{.data.postgres-password}" | base64 -d)

To get the password for "edgecraft" run:

    export POSTGRES_PASSWORD=$(kubectl get secret --namespace edgecraft postgresql -o jsonpath="{.data.password}" | base64 -d)

To connect to your database run the following command:

    kubectl run postgresql-client --rm --tty -i --restart='Never' --namespace edgecraft --image docker.io/bitnami/postgresql:15.3.0-debian-11-r7 --env="PGPASSWORD=$POSTGRES_PASSWORD" \
      --command -- psql --host postgresql -U edgecraft -d edgecraft -p 5432

    > NOTE: If you access the container using bash, make sure that you execute "/opt/bitnami/scripts/postgresql/entrypoint.sh /bin/bash" in order to avoid the error "psql: local user with ID 1001} does not exist"

To connect to your database from outside the cluster execute the following commands:

    export NODE_IP=$(kubectl get nodes --namespace edgecraft -o jsonpath="{.items[0].status.addresses[0].address}")
    export NODE_PORT=$(kubectl get --namespace edgecraft -o jsonpath="{.spec.ports[0].nodePort}" services postgresql)
    PGPASSWORD="$POSTGRES_PASSWORD" psql --host $NODE_IP --port $NODE_PORT -U edgecraft -d edgecraft
```

- 유저 권한 부여
```sh
## 그러나 슈퍼 사용자가 아닌 사용자에게 언어 C에 대한 권한을 부여하려는 경우 아래와 같은 오류가 발생합니다.
## 오류 내용
SQL Error [42501]: ERROR: permission denied for language c

## 해결방법:
# SUPERUSER 권한 부여
$ PGPASSWORD="$POSTGRES_ADMIN_PASSWORD" psql --host 192.168.88.201 --port 31000 -U postgres -d postgres
postgres=# \du
postgres=# ALTER USER edgecraft WITH SUPERUSER;
```