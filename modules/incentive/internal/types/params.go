package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/params"
)

var _ params.ParamSet = &Params{}

var (
	KeyIncentiveDefaultRewardPerBlock = []byte("incentiveDefaultRewardPerBlock")
	KeyIncentivePlans                 = []byte("incentivePlans")
)

type Params struct {
	DefaultRewardPerBlock uint16 `json:"default_reward_per_block"`
	Plans                 []Plan `json:"plans"`
}

type Plan struct {
	StartHeight    int64 `json:"start_height"`
	EndHeight      int64 `json:"end_height"`
	RewardPerBlock int64 `json:"reward_per_block"`
	TotalIncentive int64 `json:"total_incentive"`
}

func DefaultParams() Params {
	return Params{
		DefaultRewardPerBlock: 0,
		Plans: []Plan{
			{0, 10512000, 10e8, 105120000e8},
			{10512000, 21024000, 8e8, 84096000e8},
			{21024000, 31536000, 6e8, 63072000e8},
			{31536000, 42048000, 4e8, 42048000e8},
			{42048000, 52560000, 2e8, 21024000e8},
		},
	}
}

// ParamKeyTable type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyIncentiveDefaultRewardPerBlock, Value: &p.DefaultRewardPerBlock},
		{Key: KeyIncentivePlans, Value: &p.Plans},
	}
}

func (p Params) String() string {
	s := fmt.Sprintf(`Incentive Params:
  DefaultRewardPerBlock: %d`,
		p.DefaultRewardPerBlock)

	for _, p := range p.Plans {
		s += fmt.Sprintf("\n  Plan: StartHeight=%d EndHeight=%d RewardPerBlock=%d TotalIncentive=%d",
			p.StartHeight, p.EndHeight, p.RewardPerBlock, p.TotalIncentive)
	}

	return s
}