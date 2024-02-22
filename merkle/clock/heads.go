// Copyright 2022 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package clock

import (
	"bytes"
	"context"
	"encoding/binary"
	"sort"

	cid "github.com/ipfs/go-cid"

	"github.com/sourcenetwork/corekv"
	"github.com/sourcenetwork/defradb/core"
	"github.com/sourcenetwork/defradb/datastore"
	"github.com/sourcenetwork/defradb/logging"
)

// heads manages the current Merkle-CRDT heads.
type heads struct {
	store     datastore.DSReaderWriter
	namespace core.HeadStoreKey
}

func NewHeadSet(store datastore.DSReaderWriter, namespace core.HeadStoreKey) *heads {
	return &heads{
		store:     store,
		namespace: namespace,
	}
}

func (hh *heads) key(c cid.Cid) core.HeadStoreKey {
	return hh.namespace.WithCid(c)
}

func (hh *heads) Write(ctx context.Context, c cid.Cid, height uint64) error {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, height)

	return hh.store.Set(ctx, hh.key(c).ToDS().Bytes(), buf[0:n])
}

// IsHead returns if a given cid is among the current heads.
func (hh *heads) IsHead(ctx context.Context, c cid.Cid) (bool, error) {
	return hh.store.Has(ctx, hh.key(c).ToDS().Bytes())
}

// Replace replaces a head with a new CID.
func (hh *heads) Replace(ctx context.Context, old cid.Cid, new cid.Cid, height uint64) error {
	log.Info(
		ctx,
		"Replacing DAG head",
		logging.NewKV("Old", old),
		logging.NewKV("CID", new),
		logging.NewKV("Height", height))

	err := hh.store.Delete(ctx, hh.key(old).ToDS().Bytes())
	if err != nil {
		return err
	}

	err = hh.Write(ctx, new, height)
	if err != nil {
		return err
	}

	return nil
}

// List returns the list of current heads plus the max height.
// @todo Document Heads.List function
func (hh *heads) List(ctx context.Context) ([]cid.Cid, uint64, error) {
	iter := hh.store.Iterator(ctx, corekv.IterOptions{
		Prefix: hh.namespace.Bytes(),
	})

	defer func() {
		err := iter.Close(ctx)
		if err != nil {
			log.ErrorE(ctx, "Error closing results", err)
		}
	}()

	heads := make([]cid.Cid, 0)
	var maxHeight uint64
	for ; iter.Valid(); iter.Next() {

		headKey, err := core.NewHeadStoreKey(string(iter.Key()))
		if err != nil {
			return nil, 0, err
		}

		height, n := binary.Uvarint(iter.Value())
		if n <= 0 {
			return nil, 0, ErrDecodingHeight
		}
		heads = append(heads, headKey.Cid)
		if height > maxHeight {
			maxHeight = height
		}
	}
	sort.Slice(heads, func(i, j int) bool {
		ci := heads[i].Bytes()
		cj := heads[j].Bytes()
		return bytes.Compare(ci, cj) < 0
	})

	return heads, maxHeight, nil
}
