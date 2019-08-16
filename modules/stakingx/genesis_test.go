package stakingx

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGenesisState_Validate(t *testing.T) {
	//valid state
	validState := GenesisState{
		Params: DefaultParams(),
	}
	require.Nil(t, validState.ValidateGenesis())

	//invalidMinSelfDelegation
	invalidMinSelfDelegation := GenesisState{
		Params: Params{
			MinSelfDelegation: sdk.ZeroInt(),
		},
	}
	require.NotNil(t, invalidMinSelfDelegation.ValidateGenesis())

	//invalidNonBondedAddresses
	//nonBondedAddresses := make([]sdk.AccAddress, 2)
	//nonBondedAddresses[0] = testutil.ToAccAddress("myaddr")
	//nonBondedAddresses[1] = testutil.ToAccAddress("myaddr")
	//invalidNonBondedAddresses := GenesisState{
	//	Params: Params{
	//		MinSelfDelegation: sdk.NewInt(DefaultMinSelfDelegation),
	//	},
	//}
	//require.NotNil(t, invalidNonBondedAddresses.ValidateGenesis())
}

func TestDefaultGenesisState(t *testing.T) {
	defautGenesisState := DefaultGenesisState()
	require.Equal(t, DefaultParams(), defautGenesisState.Params)
}
