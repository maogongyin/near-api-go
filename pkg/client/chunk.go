package client

import (
	"context"

	"github.com/maogongyin/near-api-go/pkg/types/hash"
)

// https://docs.near.org/docs/api/rpc#chunk-details
func (c *Client) ChunkDetails(ctx context.Context, chunkHash hash.CryptoHash) (res ChunkView, err error) {
	params := map[string]string{
		"chunk_id": chunkHash.String(),
	}
	_, err = c.doRPC(ctx, &res, "chunk", nil, params)

	return
}
