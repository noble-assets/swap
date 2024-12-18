package stableswap

import "github.com/cosmos/gogoproto/proto"

// PoolWrapper is a necessary interface to allow correct `Any` marshalling.
type PoolWrapper interface {
	proto.Message
}
