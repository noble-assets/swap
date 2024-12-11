package stableswap

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreatePool{}, "swap/stableswap/CreatePool", nil)
	cdc.RegisterConcrete(&MsgUpdatePool{}, "swap/stableswap/UpdatePool", nil)
	cdc.RegisterConcrete(&MsgAddLiquidity{}, "swap/stableswap/AddLiquidity", nil)
	cdc.RegisterConcrete(&MsgRemoveLiquidity{}, "swap/stableswap/RemoveLiquidity", nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgCreatePool{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUpdatePool{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgRemoveLiquidity{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgAddLiquidity{})

	registry.RegisterInterface(
		"swap.v1.Pool",
		(*PoolWrapper)(nil),
		&Pool{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var amino = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}
