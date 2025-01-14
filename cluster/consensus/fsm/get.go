package fsm

import (
	"encoding/json"
	"errors"
	"github.com/dgraph-io/badger/v3"
	"github.com/narvikd/errorskit"
)

// Get is a DatabaseFSM's method which gets a value from a key from the LOCAL NODE.
//
// This method isn't committed since there's no need for it.
func (dbFSM DatabaseFSM) Get(k string) (any, error) {
	var result any
	dbResultValue := make([]byte, 0)

	txn := dbFSM.db.NewTransaction(false)
	defer txn.Discard()
	dbResult, errGet := txn.Get([]byte(k))
	if errGet != nil {
		return nil, errGet
	}

	errDBResultValue := dbResult.Value(func(val []byte) error {
		dbResultValue = append(dbResultValue, val...)
		return nil
	})
	if errDBResultValue != nil {
		return nil, errDBResultValue
	}

	if dbResultValue == nil || len(dbResultValue) <= 0 {
		return nil, errors.New("no result for key")
	}

	errUnmarshal := json.Unmarshal(dbResultValue, &result)
	if errUnmarshal != nil {
		return nil, errorskit.Wrap(errUnmarshal, "couldn't unmarshal get results from DB")
	}

	errCommit := txn.Commit()
	if errCommit != nil {
		return nil, errorskit.Wrap(errCommit, "couldn't commit transaction")
	}

	return result, nil
}

func (dbFSM DatabaseFSM) GetKeys() []string {
	var keys []string
	txn := dbFSM.db.NewTransaction(false)
	defer txn.Discard()

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		key := it.Item().KeyCopy(nil)
		keys = append(keys, string(key))
	}
	return keys
}
