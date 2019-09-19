package bankx_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/coinexchain/dex/modules/asset"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/bankx/internal/keeper"
	bx "github.com/coinexchain/dex/modules/bankx/internal/types"
	"github.com/coinexchain/dex/testapp"
	"github.com/coinexchain/dex/testutil"
	dex "github.com/coinexchain/dex/types"
)

var (
	myaddr   = testutil.ToAccAddress("myaddr")
	fromAddr = testutil.ToAccAddress("fromaddr")
	toAddr   = testutil.ToAccAddress("toaddr")
	feeAddr  = sdk.AccAddress(crypto.AddressHash([]byte(auth.FeeCollectorName)))
	owner    = testutil.ToAccAddress("owner")
)

func defaultContext() (*keeper.Keeper, sdk.Handler, sdk.Context) {
	app := testapp.NewTestApp()
	ctx := sdk.NewContext(app.Cms, abci.Header{}, false, log.NewNopLogger())
	app.BankxKeeper.SetParams(ctx, bx.DefaultParams())
	app.BankxKeeper.SetSendEnabled(ctx, true)
	handler := bankx.NewHandler(app.BankxKeeper)
	cet, _ := asset.NewToken("cet", "cet", sdk.NewInt(200000000000000), owner,
		false, false, false, false,
		"", "", asset.TestIdentityString)
	_ = app.AssetKeeper.SetToken(ctx, cet)
	return &app.BankxKeeper, handler, ctx
}

func TestHandlerMsgSend(t *testing.T) {
	bkx, handle, ctx := defaultContext()
	err := bkx.AddCoins(ctx, fromAddr, dex.NewCetCoins(100000000))
	require.NoError(t, err)

	msgSend := bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 0}
	handle(ctx, msgSend)

	require.Equal(t, sdk.NewInt(0).String(), bkx.GetCoins(ctx, fromAddr).AmountOf("cet").String())
	require.Equal(t, sdk.NewInt(0).String(), bkx.GetCoins(ctx, toAddr).AmountOf("cet").String())
	require.Equal(t, sdk.NewInt(100000000), bkx.GetCoins(ctx, feeAddr).AmountOf("cet"))

	fee := bkx.GetParams(ctx).LockCoinsFee
	err = bkx.AddCoins(ctx, fromAddr, dex.NewCetCoins(1000000000+fee*2))
	require.NoError(t, err)

	handle(ctx, msgSend)
	require.Equal(t, sdk.NewInt(900000000+fee*2), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000), bkx.GetCoins(ctx, feeAddr).AmountOf("cet"))

	handle(ctx, msgSend)
	require.Equal(t, sdk.NewInt(800000000+fee*2), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(200000000), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000), bkx.GetCoins(ctx, feeAddr).AmountOf("cet"))

	newMsg := bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 1}
	handle(ctx, newMsg)
	require.Equal(t, sdk.NewInt(100000000), bkx.GetLockedCoins(ctx, toAddr)[0].Coin.Amount)
	require.Equal(t, sdk.NewInt(700000000+fee), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(200000000), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000+fee), bkx.GetCoins(ctx, feeAddr).AmountOf("cet"))
	require.Equal(t, int64(1), bkx.GetLockedCoins(ctx, toAddr)[0].UnlockTime)

	newMsg2 := bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 2}
	handle(ctx, newMsg2)
	require.Equal(t, sdk.NewInt(100000000), bkx.GetLockedCoins(ctx, toAddr)[0].Coin.Amount)
	require.Equal(t, sdk.NewInt(600000000), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(200000000), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000+fee*2), bkx.GetCoins(ctx, feeAddr).AmountOf("cet"))
	require.Equal(t, int64(1), bkx.GetLockedCoins(ctx, toAddr)[0].UnlockTime)
	require.Equal(t, sdk.NewInt(100000000), bkx.GetLockedCoins(ctx, toAddr)[1].Coin.Amount)
	require.Equal(t, int64(2), bkx.GetLockedCoins(ctx, toAddr)[1].UnlockTime)
}

func TestHandlerMsgSendFail(t *testing.T) {
	bkx, handle, ctx := defaultContext()
	err := bkx.AddCoins(ctx, fromAddr, dex.NewCetCoins(100000000))
	require.NoError(t, err)

	bkx.SetSendEnabled(ctx, false)
	msgSend := bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 0}
	res := handle(ctx, msgSend)
	require.Equal(t, bank.CodeSendDisabled, res.Code)

	bkx.SetSendEnabled(ctx, true)
	msgSend = bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(200000000), UnlockTime: 0}
	res = handle(ctx, msgSend)
	require.Equal(t, sdk.CodeInsufficientCoins, res.Code)

}

func TestHandlerMsgSendUnlockFirst(t *testing.T) {
	bkx, handle, ctx := defaultContext()
	fee := bkx.GetParams(ctx).LockCoinsFee
	err := bkx.AddCoins(ctx, fromAddr, dex.NewCetCoins(1000000000+fee*2))
	require.NoError(t, err)

	msgSend := bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 2}
	handle(ctx, msgSend)
	require.Equal(t, sdk.NewInt(900000000+fee), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(0), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000+fee), bkx.GetCoins(ctx, feeAddr).AmountOf("cet"))

	msgSend2 := bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 2}
	handle(ctx, msgSend2)
	require.Equal(t, sdk.NewInt(800000000), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(0), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(100000000+fee*2), bkx.GetCoins(ctx, feeAddr).AmountOf("cet"))
}

func TestHandleMsgSetMemoRequiredAccountNotExisted(t *testing.T) {
	_, handle, ctx := defaultContext()

	msg := bankx.NewMsgSetTransferMemoRequired(testutil.ToAccAddress("xxx"), true)
	result := handle(ctx, msg)
	require.Equal(t, sdk.CodespaceRoot, result.Codespace)
	require.Equal(t, sdk.CodeUnknownAddress, result.Code)
}

func TestHandleMsgSetMemoRequiredAccountNotActivated(t *testing.T) {
	_, handle, ctx := defaultContext()

	msg := bankx.NewMsgSetTransferMemoRequired(myaddr, true)
	result := handle(ctx, msg)
	require.Equal(t, sdk.CodespaceRoot, result.Codespace)
	require.Equal(t, sdk.CodeUnknownAddress, result.Code)
}

//TODO
//func TestHandleMsgSetMemoRequiredAccountOK(t *testing.T) {
//	bkx, handle, ctx := defaultContext()
//
//	acc := auth.NewBaseAccountWithAddress(myaddr)
//	bkx.Ak.SetAccount(ctx, &acc)
//
//	msg := bx.NewMsgSetTransferMemoRequired(myaddr, true)
//	result := handle(ctx, msg)
//	require.Equal(t, sdk.CodeOK, result.Code)
//
//	accX, _ := bkx.Axk.GetAccountX(ctx, myaddr)
//	require.Equal(t, true, accX.MemoRequired)
//}

func TestUnlockQueueNotAppend(t *testing.T) {
	bkx, handle, ctx := defaultContext()
	err := bkx.AddCoins(ctx, fromAddr, dex.NewCetCoins(10100000000))
	require.NoError(t, err)

	msgSend := bankx.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: dex.NewCetCoins(100000000), UnlockTime: 10000}
	handle(ctx, msgSend)

	//send 0 to toaddr results toAccount to be created
	//to be consistent with cosmos-sdk
	require.Equal(t, sdk.NewInt(0), bkx.GetCoins(ctx, fromAddr).AmountOf("cet"))
	require.Equal(t, sdk.NewInt(0), bkx.GetCoins(ctx, toAddr).AmountOf("cet"))
}
