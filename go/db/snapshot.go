package db

import (
	"encoding/json"
	"fmt"

	"github.com/rickyfitts/distributed-evolution/go/api"
)

const SNAPSHOT = "snapshot"

func (db *DB) SaveSnapshot(snapshot api.MasterSnapshot) error {
	encoded, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}

	err = db.Set(SNAPSHOT, string(encoded))
	if err != nil {
		return fmt.Errorf("saving snapshot: %v", err)
	}

	return nil
}

func (db *DB) GetSnapshot() (api.MasterSnapshot, error) {
	data, err := db.Get(SNAPSHOT)
	if err != nil {
		return api.MasterSnapshot{}, fmt.Errorf("fetching snapshot: %v", err)
	}

	var snapshot api.MasterSnapshot
	err = json.Unmarshal([]byte(data), &snapshot)

	return snapshot, err
}
