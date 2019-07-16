package authx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Params    Params    `json:"params"`
	AccountXs AccountXs `json:"accountxs"`
}

func NewGenesisState(params Params, accountXs AccountXs) GenesisState {
	return GenesisState{
		Params:    params,
		AccountXs: accountXs,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultParams(), AccountXs{})
}

// InitGenesis - Init store state from genesis data
func InitGenesis(ctx sdk.Context, keeper AccountXKeeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)

	for _, accx := range data.AccountXs {
		accountx := NewAccountX(accx.Address, accx.MemoRequired, accx.LockedCoins, accx.FrozenCoins)
		keeper.SetAccountX(ctx, *accountx)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, keeper AccountXKeeper) GenesisState {
	var accountXs AccountXs
	keeper.IterateAccounts(ctx, func(accountX AccountX) (stop bool) {
		accountXs = append(accountXs, accountX)
		return false
	})

	return NewGenesisState(keeper.GetParams(ctx), accountXs)
}

// ValidateGenesis performs basic validation of asset genesis data returning an
// error for any failed validation criteria.
func (data GenesisState) ValidateGenesis() error {
	limit := data.Params.MinGasPriceLimit
	if limit.IsNegative() {
		return ErrInvalidMinGasPriceLimit(limit)
	}

	//TODO: validate genesis state in

	return nil
}
