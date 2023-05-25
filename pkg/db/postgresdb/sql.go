/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package postgresdb

/***************************
 * Code
 ***************************/

const getCodeGroupListSQL = `
SELECT
	A.*
FROM
	"edgecraft"."tbl_code_group" A
`

const getCodeListSQL = `
SELECT
	A.*
FROM
	"edgecraft"."tbl_code" A
	LEFT JOIN "edgecraft"."tbl_code_group" B ON (B.group_id = A.group_id)
`

const getCodeByGroupSQL = `
SELECT
	A.*
FROM
	"edgecraft"."tbl_code" A
	LEFT JOIN "edgecraft"."tbl_code_group" B ON (A.group_id = B.group_id)
WHERE
	A.group_id = $1
`

const deleteCodeByGroupSQL = `
DELETE
FROM
	"edgecraft"."tbl_code" A
WHERE
	1 = 1
{{- if ne .GroupID "" }}
AND group_id = :group_id
{{- end }}
{{- if (gt .Code 0) }}
AND code = :code
{{- end }}
`

/***************************
 * Cloud
 ***************************/

const getCloudsSQL = `
SELECT
	A.cloud_uid,
	A.name,
	A.type,
	A.state,
	A.description,
	A.created_at,
	B.k8s_version as version,
	(SELECT COUNT(node_uid) FROM "edgecraft"."tbl_cloud_node" WHERE cloud_uid = A.cloud_uid) AS nodeCount
FROM
	"edgecraft"."tbl_cloud" A
	LEFT JOIN "edgecraft"."tbl_cloud_cluster" B ON (B.cloud_uid = A.cloud_uid)
`

/***************************
 * Cloud - Cluster
 ***************************/

const getClustersSQL = `
SELECT 
	A.*
FROM 
	"edgecraft"."tbl_cloud_cluster" A
WHERE
	A.cloud_uid = $1
`

const deleteClusters = `
DELETE
FROM "edgecraft"."tbl_cloud_cluster" A
WHERE
	A.cloud_uid = $1
`

/***************************
 * Cloud - Node
 ***************************/

const getNodeSQL = `
SELECT 
	*
FROM 
	"edgecraft"."tbl_cloud_node" A
WHERE
	1 = 1
{{- if ne (isNull .CloudUid) false }}
AND cloud_uid = :cloud_uid
{{- end }}
{{- if ne (isNull .ClusterUid) false }}
AND cluster_uid = :cluster_uid
{{- end }}
{{- if ne (isNull .NodeUid) false }}
AND node_uid =:node_uid
{{- end }}
`

const deleteNodeSQL = `
DELETE
FROM 
	"edgecraft"."tbl_cloud_node" A
WHERE
	1 = 1
{{- if ne (isNull .CloudUid) false }}
AND cloud_uid = :cloud_uid
{{- end }}
{{- if ne (isNull .ClusterUid) false }}
AND cluster_uid = :cluster_uid
{{- end }}
{{- if ne (isNull .NodeUid) false }}
AND node_uid =:node_uid
{{- end }}
`

/******************************
 * Cloud - Cluster (Openstack)
 ******************************/

const getOpenstackClustersSQL = `
SELECT 
	A.cloud_uid,
	A.cluster_uid,
	A.name,
	A.state,
	(SELECT SUM(node_count) FROM "edgecraft"."tbl_nodeset" B WHERE B.cluster_uid = A.cluster_uid GROUP BY B.cluster_uid) AS node_count,
	A.version,
	A.created_at
FROM 
	"edgecraft"."tbl_cluster" A
WHERE
	A.cloud_uid = $1
`

const deleteOpenstackClustersSQL = `
DELETE
FROM "edgecraft"."tbl_cluster" A
WHERE
	A.cloud_uid = $1
`

const updateProvisionStatusSQL = `
UPDATE	"edgecraft"."tbl_cluster" A
   SET	state = $3
WHERE
		cloud_uid = $1
AND		cluster_uid = $2
`

/***************************************
 * Cloud - Cluster - NodeSet (Openstack)
 ****************************************/

const getNodeSetsSQL = `
SELECT
	A.*
FROM
	"edgecraft"."tbl_nodeset" A
WHERE
	A.cluster_uid = $1
`

const deleteNodeSetsSQL = `
DELETE
FROM "edgecraft"."tbl_nodeset" A
WHERE
	A.cluster_uid = $1
`

/***************************************
 * Cloud - Cluster - Benchmarks (Openstack)
 ****************************************/

const getOpenstackBenchmarksListSQL = `
SELECT 
	A.cluster_uid,
	A.benchmarks_uid,
	A.totals,
	A.state,
	A.reason,
	A.created_at
FROM 
	"edgecraft"."tbl_cluster_benchmarks" A
WHERE
		A.cluster_uid = $1
ORDER BY A.created_at DESC
`

const deleteOpenstackBenchmarksSQL = `
DELETE 
FROM "edgecraft"."tbl_cluster_benchmarks" A
WHERE
	A.cluster_uid = $1
`

/***************************************
 * Cloud - Cluster - Backup
 ****************************************/

const getBackupListSQL = `
 SELECT 
	 A.cloud_uid,
	 A.cluster_uid,
	 A.backres_uid,
	 A.name,
	 A.type,
	 A.status,
	 A.reason,
	 A.backup_name,
	 A.creator,
	 A.created_at
 FROM 
	 "edgecraft"."tbl_cluster_backres" A
 WHERE
		A.cloud_uid = $1
 AND	A.type = 'B'
 ORDER BY A.created_at DESC
`

const getBackupSQL = `
 SELECT 
	 A.cloud_uid,
	 A.cluster_uid,
	 A.backres_uid,
	 A.name,
	 A.type,
	 A.status,
	 A.reason,
	 A.backup_name,
	 A.creator,
	 A.created_at
 FROM 
	 "edgecraft"."tbl_cluster_backres" A
 WHERE
		A.backres_uid = $1
 AND	A.type = 'B'
`

const getRestoreListSQL = `
 SELECT 
	 A.cloud_uid,
	 A.cluster_uid,
	 A.backres_uid,
	 A.name,
	 A.type,
	 A.status,
	 A.reason,
	 A.backup_name,
	 A.creator,
	 A.created_at
 FROM 
	 "edgecraft"."tbl_cluster_backres" A
 WHERE
		A.cloud_uid = $1
 AND	A.type = 'R'
 ORDER BY A.created_at DESC
`

const deleteBackResSQL = `
DELETE
FROM 
	"edgecraft"."tbl_cluster_backres" A
WHERE
	A.cloud_uid = $1
AND A.cluster_uid = $2
AND A.backres_uid = $3
`

const updateBackRresStatusSQL = `
UPDATE	"edgecraft"."tbl_cluster_backres"
   SET	status = $4
WHERE
		cloud_uid = $1
AND		cluster_uid = $2
AND		backres_uid = $3
`

const getBackResDuplicate = `
SELECT EXISTS(SELECT 1 FROM "edgecraft"."tbl_cluster_backres" WHERE name = $1)
`
