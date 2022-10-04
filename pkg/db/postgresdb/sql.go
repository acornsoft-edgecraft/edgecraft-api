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
	(SELECT COUNT(node_uid) FROM tbl_cloud_node WHERE cloud_uid = A.cloud_uid) AS nodeCount
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
