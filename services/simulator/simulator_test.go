package simulator

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestSimulationResultPrintCorrectMap(t *testing.T) {
	res := SimulationResult{
		ResultMap: planetMap{
			"London": {
				name:        "London",
				isDestroyed: true,
				directions:  []mapDirection{{directionType: "west", directionValue: "Bolton"}},
			},
			"Bolton": {
				name:        "Bolton",
				isDestroyed: false,
				directions:  []mapDirection{{directionType: "east", directionValue: "London"}},
			},
		},
	}
	builder := strings.Builder{}
	_ = res.PrintResultMap(&builder)
	require.Equal(t, "Bolton\n", builder.String())
}

func TestCityPrintsOutCorrectBattleMessage(t *testing.T) {
	c := city{
		name:   "Boston",
		aliens: []int64{11, 2, 5},
	}
	require.Equal(t, "Aliens: 👾11, 👾2, 👾5 have met in the city of Boston. ⚔ Battle destroyed the city.", c.battleMessage())
}
