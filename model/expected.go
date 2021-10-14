package model

import (
	"context"
	"strings"

	"google.golang.org/appengine/v2/datastore"
	"google.golang.org/appengine/v2/log"
)

const (
	DatastoreKind = "gas#esun"
	DatastoreKey  = "expected"
)

// Expected is default expected from Firestore
type Expected struct {
	Expected float64 `json:"expected"`
}

// PutDatastore is store to Datastore
func (ex *Expected) PutDatastore(ctx context.Context, value interface{}) {
	datastoreKey := datastore.NewKey(ctx, DatastoreKind, DatastoreKey, 0, nil)

	log.Infof(ctx, "PUT into Datastore with Key: %v", DatastoreKey)
	if _, err := datastore.Put(ctx, datastoreKey, value); err != nil {
		log.Errorf(ctx, "PUT into Datastore failed %v", err)
	}
}

// GetDatastore is get from Datastore
func (ex *Expected) GetDatastore(ctx context.Context, value interface{}) {
	datastoreKey := datastore.NewKey(ctx, DatastoreKind, DatastoreKey, 0, nil)

	log.Infof(ctx, "GET Datastore with Key: %v", DatastoreKey)
	if err := datastore.Get(ctx, datastoreKey, value); err != nil {
		if !strings.HasPrefix(err.Error(), `datastore: cannot load field`) {
			log.Errorf(ctx, "GET Datastore failed %v", err)
		}
	}
}
