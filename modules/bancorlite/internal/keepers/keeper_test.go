package keepers

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var owner = sdk.AccAddress("user")

func TestBancorInfo_UpdateStockInPool(t *testing.T) {
	type fields struct {
		Owner       sdk.AccAddress
		Token       string
		MaxSupply   sdk.Int
		MaxPrice    sdk.Dec
		Price       sdk.Dec
		StockInPool sdk.Int
		MoneyInPool sdk.Int
	}
	type args struct {
		stockInPool sdk.Int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "positive",
			fields: fields{
				Owner:       owner,
				Token:       "bch",
				MaxSupply:   sdk.NewInt(100),
				MaxPrice:    sdk.NewDec(10),
				StockInPool: sdk.NewInt(10),
			},
			args: args{
				stockInPool: sdk.NewInt(20),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bi := &BancorInfo{
				Owner:       tt.fields.Owner,
				Token:       tt.fields.Token,
				MaxSupply:   tt.fields.MaxSupply,
				MaxPrice:    tt.fields.MaxPrice,
				Price:       tt.fields.Price,
				StockInPool: tt.fields.StockInPool,
				MoneyInPool: tt.fields.MoneyInPool,
			}
			if got := bi.UpdateStockInPool(tt.args.stockInPool); got != tt.want {
				t.Errorf("BancorInfo.UpdateStockInPool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBancorInfo_IsConsistent(t *testing.T) {
	type fields struct {
		Owner       sdk.AccAddress
		Token       string
		MaxSupply   sdk.Int
		MaxPrice    sdk.Dec
		Price       sdk.Dec
		StockInPool sdk.Int
		MoneyInPool sdk.Int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "positive",
			fields: fields{
				Owner:       owner,
				Token:       "bch",
				MaxSupply:   sdk.NewInt(100),
				MaxPrice:    sdk.NewDec(10),
				Price:       sdk.NewDec(1),
				StockInPool: sdk.NewInt(90),
				MoneyInPool: sdk.NewInt(5),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bi := &BancorInfo{
				Owner:       tt.fields.Owner,
				Token:       tt.fields.Token,
				MaxSupply:   tt.fields.MaxSupply,
				MaxPrice:    tt.fields.MaxPrice,
				Price:       tt.fields.Price,
				StockInPool: tt.fields.StockInPool,
				MoneyInPool: tt.fields.MoneyInPool,
			}
			if got := bi.IsConsistent(); got != tt.want {
				t.Errorf("BancorInfo.IsConsistent() = %v, want %v", got, tt.want)
			}
		})
	}
}