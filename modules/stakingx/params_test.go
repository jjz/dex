package stakingx

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

func TestDefaultParams(t *testing.T) {
	params := DefaultParams()
	require.Equal(t, "100000000000000", params.MinSelfDelegation.String())
}

func defaultContext() (sdk.Context, params.Keeper) {
	cdc := codec.New()
	skey := sdk.NewKVStoreKey("test")
	tkey := sdk.NewTransientStoreKey("transient_test")

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(skey, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, db)

	_ = cms.LoadLatestVersion()
	ctx := sdk.NewContext(cms, abci.Header{}, false, log.NewNopLogger())
	paramsKeeper := params.NewKeeper(cdc, skey, tkey, params.DefaultCodespace)

	return ctx, paramsKeeper
}

func TestParamGetSet(t *testing.T) {
	ctx, paramsKeeper := defaultContext()
	subspace := paramsKeeper.Subspace(DefaultParamspace)
	sxk := NewKeeper(sdk.NewKVStoreKey("test"), codec.New(), subspace, nil, &staking.Keeper{}, distribution.Keeper{}, auth.AccountKeeper{}, nil, nil, "")

	testParam := Params{
		MinSelfDelegation:          sdk.ZeroInt(),
		MinMandatoryCommissionRate: DefaultMinMandatoryCommissionRate,
	}

	//expect SetParam don't panic
	require.NotPanics(t, func() { sxk.SetParams(ctx, testParam) }, "stakingx keeper SetParam panics")

	//expect GetParam equals defaultParam
	require.Equal(t, testParam, sxk.GetParams(ctx))
}
