// Copyright 2017 GRAIL, Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package runner

import (
	"context"

	"github.com/grailbio/reflow"
	"github.com/grailbio/reflow/pool"
)

// Cluster is a kind of pool.Pool that also allows the user to
// directly reserve an alloc. This way, the cluster can be responsive
// to demand.
type Cluster interface {
	pool.Pool

	// Allocate reserves an alloc of at least min, and at most max resources.
	// The cluster may scale elastically in order to meet this demand.
	// Labels are passed down to the underlying pool.
	Allocate(ctx context.Context, min, max reflow.Resources, labels pool.Labels) (pool.Alloc, error)
}

// A StaticCluster implements a pass-through Cluster on top of a pool.Pool.
type StaticCluster struct {
	pool.Pool
}

// Allocate reserves an alloc from the underlying static pool.
func (s *StaticCluster) Allocate(ctx context.Context, min, max reflow.Resources, labels pool.Labels) (pool.Alloc, error) {
	return pool.Allocate(ctx, s, min, max, labels)
}
