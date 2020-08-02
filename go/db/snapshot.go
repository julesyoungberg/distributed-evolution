package db

import (
	"github.com/rickyfitts/distributed-evolution/go/api"
)

const SnapshotKey = "snapshot"

func (db *DB) SaveSnapshot(snapshot api.MasterSnapshot) error {
	return db.SetData(SnapshotKey, snapshot)
}

func (db *DB) GetSnapshot() (api.MasterSnapshot, error) {
	var snapshot api.MasterSnapshot
	err := db.GetData(SnapshotKey, &snapshot)
	return snapshot, err
}
