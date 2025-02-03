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

package mocks

import (
	"testing"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	codectestutil "github.com/cosmos/cosmos-sdk/codec/testutil"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	swap "swap.noble.xyz"
	modulev1 "swap.noble.xyz/api/module/v1"
	"swap.noble.xyz/keeper"
	"swap.noble.xyz/types"
)

func SwapKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	return SwapKeeperWithKeepers(
		t,
		AccountKeeper{
			Accounts:      make(map[string]sdk.AccountI),
			AccountNumber: 0,
		},
		BankKeeper{
			Restriction: NoOpSendRestrictionFn,
		},
	)
}

func SwapKeeperWithKeepers(t testing.TB, account AccountKeeper, bank BankKeeper) (*keeper.Keeper, sdk.Context) {
	key := storetypes.NewKVStoreKey(types.ModuleName)
	tkey := storetypes.NewTransientStoreKey("transient_authority")
	wrapper := testutil.DefaultContextWithDB(t, key, tkey)

	cfg := MakeTestEncodingConfig("noble")
	types.RegisterInterfaces(cfg.InterfaceRegistry)

	k := keeper.NewKeeper(
		cfg.Codec,
		runtime.NewKVStoreService(key),
		runtime.ProvideEventService(),
		runtime.ProvideHeaderInfoService(&runtime.AppBuilder{}),
		log.NewNopLogger(),
		"authority",
		"uusdn",
		1e6,
		&modulev1.StableSwap{
			UnbondingBlockDelta: 10,
		},
		address.NewBech32Codec("noble"),
		account,
		bank,
	)

	// bank = bank.WithSendCoinsRestriction(fun)
	k.SetBankKeeper(bank)

	swap.InitGenesis(wrapper.Ctx, k, *types.DefaultGenesisState())
	return k, wrapper.Ctx
}

// MakeTestEncodingConfig is a modified testutil.MakeTestEncodingConfig that
// sets a custom Bech32 prefix in the interface registry.
func MakeTestEncodingConfig(prefix string, modules ...module.AppModuleBasic) moduletestutil.TestEncodingConfig {
	aminoCodec := codec.NewLegacyAmino()
	interfaceRegistry := codectestutil.CodecOptions{
		AccAddressPrefix: prefix,
	}.NewInterfaceRegistry()
	codec := codec.NewProtoCodec(interfaceRegistry)

	encCfg := moduletestutil.TestEncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             codec,
		TxConfig:          tx.NewTxConfig(codec, tx.DefaultSignModes),
		Amino:             aminoCodec,
	}

	mb := module.NewBasicManager(modules...)

	std.RegisterLegacyAminoCodec(encCfg.Amino)
	std.RegisterInterfaces(encCfg.InterfaceRegistry)
	mb.RegisterLegacyAminoCodec(encCfg.Amino)
	mb.RegisterInterfaces(encCfg.InterfaceRegistry)

	return encCfg
}
