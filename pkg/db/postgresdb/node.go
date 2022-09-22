/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package postgresdb

import "github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"

// InsertNode - Insert a new Baremetal Node
func (db *DB) InsertNode(node *model.NodeTable) error {
	return db.GetClient().Insert(node)
}
