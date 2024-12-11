package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"swap.noble.xyz/types/stableswap"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	stableswap.RegisterLegacyAminoCodec(cdc)

	cdc.RegisterConcrete(&MsgSwap{}, "swap/Swap", nil)
	cdc.RegisterConcrete(&MsgPauseByAlgorithm{}, "swap/PauseByAlgorithm", nil)
	cdc.RegisterConcrete(&MsgPauseByPoolIds{}, "swap/PauseByPoolIds", nil)
	cdc.RegisterConcrete(&MsgUnpauseByAlgorithm{}, "swap/UnpauseByAlgorithm", nil)
	cdc.RegisterConcrete(&MsgUnpauseByPoolIds{}, "swap/UnpauseByPoolIds", nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	stableswap.RegisterInterfaces(registry)

	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSwap{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgPauseByAlgorithm{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgPauseByPoolIds{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUnpauseByAlgorithm{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUnpauseByPoolIds{})

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var amino = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}
