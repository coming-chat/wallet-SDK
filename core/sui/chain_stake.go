package sui

import (
	"context"
	"strconv"
)

type Validator struct {
	Address    string
	Name       string
	Desc       string
	ImageUrl   string
	ProjectUrl string
	APY        float64
	StakedSui  string
	Epoch      int64
}

func (c *Chain) GetStakeList() ([]*Validator, error) {
	cli, err := c.client()
	if err != nil {
		return nil, err
	}
	state, err := cli.GetSuiSystemState(context.Background())
	if err != nil {
		return nil, err
	}

	var validators = []*Validator{}
	for _, v := range state.Validators.ActiveValidators {
		meta := v.Metadata

		validator := Validator{
			Address:    meta.SuiAddress.String(),
			Name:       string(meta.Name),
			Desc:       string(meta.Description),
			ImageUrl:   string(meta.ImageUrl),
			ProjectUrl: string(meta.ProjectUrl),
			APY:        v.CalculateAPY(state.Epoch),
			StakedSui:  strconv.FormatUint(v.DelegationStakingPool.SuiBalance, 10),
			Epoch:      int64(state.Epoch),
		}
		validators = append(validators, &validator)
	}

	return validators, nil
}
