// Copyright 2022 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package db

import (
	"context"
	"testing"

	badger "github.com/dgraph-io/badger/v4"
	badgerkv "github.com/sourcenetwork/corekv/badger"
	"github.com/sourcenetwork/defradb/datastore"
)

func newMemoryDB(ctx context.Context) (*implicitTxnDB, error) {
	opts := badger.DefaultOptions("").WithInMemory(true)
	rootstore, err := badgerkv.NewDatastore("", opts)
	if err != nil {
		return nil, err
	}
	return newDB(ctx, rootstore.(datastore.RootStore))
}

func TestNewDB(t *testing.T) {
	ctx := context.Background()
	opts := badger.DefaultOptions("").WithInMemory(true)
	rootstore, err := badgerkv.NewDatastore("", opts)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = NewDB(ctx, rootstore.(datastore.RootStore))
	if err != nil {
		t.Error(err)
	}
}
