package bankx

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/dex/testutil"
)

// MsgSetMemoRequired tests

func TestSetMemoRequiredRoute(t *testing.T) {
	addr := sdk.AccAddress([]byte("addr"))
	msg := NewMsgSetTransferMemoRequired(addr, true)
	require.Equal(t, msg.Route(), "bankx")
	require.Equal(t, msg.Type(), "set_memo_required")
}

func TestSetMemoRequiredValidation(t *testing.T) {
	validAddr := sdk.AccAddress([]byte("addr"))
	var emptyAddr sdk.AccAddress

	// nolint
	testutil.ValidateBasic(t, []testutil.TestCase{
		{true, NewMsgSetTransferMemoRequired(validAddr, true)},
		{true, NewMsgSetTransferMemoRequired(validAddr, false)},
		{false, NewMsgSetTransferMemoRequired(emptyAddr, true)},
		{false, NewMsgSetTransferMemoRequired(emptyAddr, false)},
	})
}

func TestSetMemoRequiredGetSignBytes(t *testing.T) {
	addr := sdk.AccAddress([]byte("addr"))
	msg := NewMsgSetTransferMemoRequired(addr, true)
	sign := msg.GetSignBytes()

	expected := `{"type":"cet-chain/MsgSetMemoRequired","value":{"address":"cosmos1v9jxguspv4h2u","required":true}}`
	require.Equal(t, expected, string(sign))
}

func TestSetMemoRequiredGetSigners(t *testing.T) {
	addr := sdk.AccAddress([]byte("addr"))
	msg := NewMsgSetTransferMemoRequired(addr, true)
	signers := msg.GetSigners()
	require.Equal(t, 1, len(signers))
	require.Equal(t, addr, signers[0])
}
