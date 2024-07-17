package indexer

import (
	"context"

	"github.com/pkg/errors"
)

// Requires `indexerApi` implements the `Interface` interface.
var _ Interface = (*IndexerApi)(nil)

// IndexerApi indexer service configuration
type IndexerApi struct {
	Namespace string
	manager   *NodeManager
}

// NewIndexerApi creates indexer service configuration
func NewIndexerApi(manager *NodeManager) *IndexerApi {
	return &IndexerApi{"indexer", manager}
}

// GetShardedNodes return storage node list
func (api *IndexerApi) GetShardedNodes(ctx context.Context) (ShardedNodes, error) {
	trusted, err := api.manager.Trusted()
	if err != nil {
		return ShardedNodes{}, errors.WithMessage(err, "Failed to retrieve trusted nodes")
	}

	return ShardedNodes{
		Trusted:    trusted,
		Discovered: api.manager.Discovered(),
	}, nil
}

// GetNodes return storage nodes with IP location information.
func (api *IndexerApi) GetNodes(ctx context.Context) ([]*NodeInfo, error) {
	var nodes []*NodeInfo

	trusted, err := api.manager.Trusted()
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to retrieve trusted nodes")
	}

	for _, v := range trusted {
		node := &NodeInfo{
			ShardedNode: v,
		}

		if loc, ok := api.manager.Location(v.URL); ok {
			node.Location = loc
		}

		nodes = append(nodes, node)
	}

	for _, v := range api.manager.Discovered() {
		node := &NodeInfo{
			ShardedNode: v,
		}

		if loc, ok := api.manager.Location(v.URL); ok {
			node.Location = loc
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}
