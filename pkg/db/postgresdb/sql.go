/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package postgresdb

/**
 * Code
 **/

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

// const deleteCodeByGroupSQL = `
// DELETE
// FROM "edgecraft"."tbl_code" A
// WHERE
// 	A.group_id = $1
// `
