// Copyright 2023 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package cli

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/spf13/cobra"

	"github.com/sourcenetwork/defradb/client"
)

func MakeP2PReplicatorSetCommand() *cobra.Command {
	var collections []string
	var cmd = &cobra.Command{
		Use:   "set [-c, --collection] <peer>",
		Short: "Set a P2P replicator",
		Long: `Add a new target replicator.
A replicator replicates one or all collection(s) from this node to another.
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store := mustGetStoreContext(cmd)

			addr, err := peer.AddrInfoFromString(args[0])
			if err != nil {
				return err
			}
			rep := client.Replicator{
				Info:    *addr,
				Schemas: collections,
			}
			return store.SetReplicator(cmd.Context(), rep)
		},
	}

	cmd.Flags().StringSliceVarP(&collections, "collection", "c",
		[]string{}, "Define the collection for the replicator")
	return cmd
}
