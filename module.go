// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, NASD Inc. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package swap

import (
	"context"
	"encoding/json"
	"fmt"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/event"
	"cosmossdk.io/core/header"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	modulev1 "swap.noble.xyz/api/module/v1"
	stableswapv1 "swap.noble.xyz/api/stableswap/v1"
	swapv1 "swap.noble.xyz/api/v1"
	"swap.noble.xyz/keeper"
	"swap.noble.xyz/types"
	"swap.noble.xyz/types/stableswap"
)

// ConsensusVersion defines the current Noble Swap module consensus version.
const ConsensusVersion = 1

var (
	_ module.AppModuleBasic      = AppModule{}
	_ appmodule.AppModule        = AppModule{}
	_ module.HasConsensusVersion = AppModule{}
	_ module.HasGenesis          = AppModule{}
	_ module.HasGenesisBasics    = AppModuleBasic{}
	_ module.HasServices         = AppModule{}
)

//

type AppModuleBasic struct{}

func NewAppModuleBasic() AppModuleBasic {
	return AppModuleBasic{}
}

func (AppModuleBasic) Name() string { return types.ModuleName }

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}

	if err := stableswap.RegisterQueryHandlerClient(context.Background(), mux, stableswap.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var genesis types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genesis); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}

	return genesis.Validate()
}

//

type AppModule struct {
	AppModuleBasic

	keeper *keeper.Keeper
}

func NewAppModule(keeper *keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(),
		keeper:         keeper,
	}
}

func (AppModule) IsOnePerModuleType() {}

func (AppModule) IsAppModule() {}

func (AppModule) ConsensusVersion() uint64 { return ConsensusVersion }

func (m AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, bz json.RawMessage) {
	var genesis types.GenesisState
	cdc.MustUnmarshalJSON(bz, &genesis)

	InitGenesis(ctx, m.keeper, genesis)
}

func (m AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genesis := ExportGenesis(ctx, m.keeper)
	return cdc.MustMarshalJSON(genesis)
}

func (m AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServer(m.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(m.keeper))

	stableswap.RegisterMsgServer(cfg.MsgServer(), keeper.NewStableSwapMsgServer(m.keeper))
	stableswap.RegisterQueryServer(cfg.QueryServer(), keeper.NewStableSwapQueryServer(m.keeper))
}

func (m AppModule) BeginBlock(ctx context.Context) error {
	return m.keeper.BeginBlocker(ctx)
}

//

func (AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: swapv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Swap",
					Use:       "swap [amount] [routes] [min]",
					Short:     "Execute a amount swap across specified routes",
					Long:      "Swaps a specified `amount` along the provided `routes`, with a `min` value that sets the minimum acceptable output to protect against slippage.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "amount"},
						{ProtoField: "routes"},
						{ProtoField: "min"},
					},
				},
				{
					RpcMethod:      "WithdrawRewards",
					Use:            "withdraw-rewards",
					Short:          "Collect rewards generated from liquidity provided to pools",
					Long:           "Collects the rewards accrued from swap fees in the pools where liquidity has been provided.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},
				{
					RpcMethod: "WithdrawProtocolFees",
					Use:       "withdraw-protocol-fees [to]",
					Short:     "Collects protocol fees and transfers them to the specified address",
					Long:      "This command collects accumulated protocol fees and transfers them to the specified address. Only the authority can execute this command.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "to"},
					},
				},
				{
					RpcMethod: "PauseByAlgorithm",
					Use:       "pause-by-algorithm [algorithm]",
					Short:     "Pause all the pools operations with the provided algorithm",
					Long:      "Pause all the pools whose operations are associated with the specified algorithm.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "algorithm"},
					},
				},
				{
					RpcMethod: "PauseByPoolIds",
					Use:       "pause-by-pool-ids [pool_ids]",
					Short:     "Pause all the pools identified by the provided pool IDs",
					Long:      "Pause all the pools identified by the provided pool IDs.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "pool_ids"},
					},
				},
				{
					RpcMethod: "UnpauseByAlgorithm",
					Use:       "unpause-by-algorithm [algorithm]",
					Short:     "Unpause all the pools operations with the provided algorithm",
					Long:      "Unpause all the pools whose operations are associated with the specified algorithm.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "algorithm"},
					},
				},
				{
					RpcMethod: "UnpauseByPoolIds",
					Use:       "unpause-by-pool-ids",
					Short:     "Unpause all the pools identified by the provided pool IDs.",
					Long:      "Unpause all the pools identified by the provided pool IDs.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "pool_ids"},
					},
				},
			},
			SubCommands: map[string]*autocliv1.ServiceCommandDescriptor{
				"stableswap": {
					Service: stableswap.Msg_serviceDesc.ServiceName,
					RpcCommandOptions: []*autocliv1.RpcCommandOptions{
						{
							RpcMethod: "CreatePool",
							Use:       "create-pool [pair] [protocol_fee_percentage] [rewards_fee] [initial_a] [future_a] [future_a_time] [rate_multipliers]",
							Short:     "Create a new stable swap pool",
							Long:      "Creates a stable swap pool with specified parameters, including the `pair` coin, `fees`, amplification factors (`initial_a` and `future_a`), and `rate_multipliers` for dynamic rates.",
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{
								{ProtoField: "pair"},
								{ProtoField: "protocol_fee_percentage"},
								{ProtoField: "rewards_fee"},
								{ProtoField: "initial_a"},
								{ProtoField: "future_a"},
								{ProtoField: "future_a_time"},
								{ProtoField: "rate_multipliers", Varargs: true},
							},
						},
						{
							RpcMethod: "UpdatePool",
							Use:       "update-pool [pool_id] [protocol_fee_percentage] [rewards_fee] [initial_a] [future_a] [future_a_time] [rate_multipliers]",
							Short:     "Updates a stable swap pool",
							Long:      "Update a stable swap pool with specified parameters, including the `fees`, amplification factors (`future_a_time` and `future_a`), and `rate_multipliers` for dynamic rates.",
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{
								{ProtoField: "pool_id"},
								{ProtoField: "protocol_fee_percentage"},
								{ProtoField: "rewards_fee"},
								{ProtoField: "initial_a"},
								{ProtoField: "future_a"},
								{ProtoField: "future_a_time"},
								{ProtoField: "rate_multipliers", Varargs: true},
							},
						},
						{
							RpcMethod: "AddLiquidity",
							Use:       "add-liquidity [pool_id] [slippage_percentage] [amount]",
							Short:     "Add liquidity to a specified `StableSwap` pool",
							Long:      "Adds a specified amount of liquidity to the pool identified by `pool_id`.",
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{
								{ProtoField: "pool_id"},
								{ProtoField: "slippage_percentage"},
								{ProtoField: "amount", Varargs: true},
							},
						},
						{
							RpcMethod: "RemoveLiquidity",
							Use:       "remove-liquidity [pool_id] [percentage]",
							Short:     "Remove a percentage of liquidity from a specified pool",
							Long:      "Removes a specified percentage of liquidity from the pool identified by `pool_id`.",
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{
								{ProtoField: "pool_id"},
								{ProtoField: "percentage"},
							},
						},
					},
				},
			},
			EnhanceCustomCommand: true,
		},
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: swapv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "SimulateSwap",
					Use:       "simulate [amount] [routes] [min]",
					Short:     "Simulate a token swap transaction",
					Long:      "Simulate the expected output and associated fees for a token swap, without executing the transaction.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "amount"},
						{ProtoField: "routes"},
						{ProtoField: "min"},
					},
				},
				{
					RpcMethod: "Paused",
					Use:       "paused",
					Short:     "Check if the module is paused",
					Long:      "Queries the module to determine if operations are currently paused, indicating a temporary or emergency halt of activity.",
				},
				{
					RpcMethod:      "Rates",
					Use:            "rates (algorithm)",
					Short:          "Query current swap rates",
					Long:           "Retrieves the swap rates for available algorithms. Optionally specify an `algorithm` to filter results.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "algorithm", Optional: true}},
				},
				{
					RpcMethod: "Rate",
					Use:       "rate [denom] (algorithm)",
					Short:     "Query swap rate for a specific denom",
					Long:      "Returns the swap rate for a specified `denom`, with an optional `algorithm` parameter to filter results.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "denom"},
						{ProtoField: "algorithm", Optional: true},
					},
				},
				{
					RpcMethod: "Pools",
					Use:       "pools",
					Short:     "List all available liquidity pools",
					Long:      "Provides a list of all active liquidity pools within the module.",
				},
				{
					RpcMethod:      "Pool",
					Use:            "pool",
					Short:          "Query details of a specific pool",
					Long:           "Fetches detailed information about a specific liquidity pool, including current liquidity, participants, and fees.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "pool_id"}},
				},
			},
			SubCommands: map[string]*autocliv1.ServiceCommandDescriptor{
				"stableswap": {
					Service: stableswapv1.Query_ServiceDesc.ServiceName,
					RpcCommandOptions: []*autocliv1.RpcCommandOptions{
						{
							RpcMethod:      "PositionsByProvider",
							Use:            "positions-by-provider [provider]",
							Short:          "List positions by a specific provider",
							Long:           "Retrieves all active positions attributed to a specified `provider` within the module.",
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "provider"}},
						},
						{
							RpcMethod:      "UnbondingPositionsByProvider",
							Use:            "unbonding-positions-by-provider [provider]",
							Short:          "List unbonding positions by a specific provider",
							Long:           "Retrieves the unbonding positions attributed to a specified `provider` within the module.",
							PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "provider"}},
						},
					},
				},
			},
		},
	}
}

//

func init() {
	appmodule.Register(&modulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Config *modulev1.Module

	Cdc          codec.Codec
	StoreService store.KVStoreService

	EventService  event.Service
	HeaderService header.Service
	Logger        log.Logger

	AddressCodec  address.Codec
	AccountKeeper types.AccountKeeper
	BankKeeper    types.BankKeeper
}

type ModuleOutputs struct {
	depinject.Out

	Keeper *keeper.Keeper
	Module appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	if in.Config.Authority == "" {
		panic("authority for x/swap module must be set")
	}

	if in.Config.BaseDenom == "" {
		panic("base_denom for x/swap module must be set")
	}

	if in.Config.MaxAddLiquiditySlippagePercentage <= 0 {
		panic("max_add_liquidity_slippage_percentage for x/swap module must be set")
	}

	if in.Config.Stableswap == nil {
		panic("stableswap config for x/swap/stableswap module must be set")
	}

	if in.Config.Stableswap.UnbondingBlockDelta <= 0 {
		panic("unbonding_block_delta for x/swap/stableswap module must be set")
	}

	authority := authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	k := keeper.NewKeeper(
		in.Cdc,
		in.StoreService,
		in.EventService,
		in.HeaderService,
		in.Logger,
		authority.String(),
		in.Config.BaseDenom,
		in.Config.BaseMinimumDeposit,
		in.Config.MaxAddLiquiditySlippagePercentage,
		in.Config.Stableswap,
		in.AddressCodec,
		in.AccountKeeper,
		in.BankKeeper)
	m := NewAppModule(k)

	return ModuleOutputs{Keeper: k, Module: m}
}
