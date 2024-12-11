package types

func DefaultGenesisState() *GenesisState { return &GenesisState{} }

func (gen *GenesisState) Validate() error { return nil }
