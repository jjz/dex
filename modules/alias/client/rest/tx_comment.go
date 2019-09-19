package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"net/http"

	"github.com/coinexchain/dex/modules/alias/internal/types"
	"github.com/coinexchain/dex/modules/authx/client/restutil"
)

type AliasUpdateReq struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	Alias     string       `json:"alias"`
	IsAdd     bool         `json:"is_add"`
	AsDefault bool         `json:"as_default"`
}

var _ restutil.RestReq = &AliasUpdateReq{}

func (req *AliasUpdateReq) GetBaseReq() *rest.BaseReq {
	return &req.BaseReq
}

func (req *AliasUpdateReq) GetMsg(w http.ResponseWriter, sender sdk.AccAddress) sdk.Msg {
	return &types.MsgAliasUpdate{
		Owner:     sender,
		Alias:     req.Alias,
		IsAdd:     req.IsAdd,
		AsDefault: req.AsDefault,
	}
}

func aliasUpdateHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	var req AliasUpdateReq
	builder := restutil.NewRestHandlerBuilder(cdc, cliCtx, &req)
	return builder.Build()
}
