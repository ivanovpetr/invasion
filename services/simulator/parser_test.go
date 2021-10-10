package simulator

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

type errorCase struct{ input, expectedErrorText string }

var errTestCases = []errorCase{
	{
		input:             `London`,
		expectedErrorText: "testing:1:7: unexpected EOF, city must contain at least one direction",
	},
	{
		input:             "London\n",
		expectedErrorText: "testing:2:1: unexpected newline, city must contain at least one direction",
	},
	{
		input:             `London, west=Boston`,
		expectedErrorText: "testing:1:8: expected a valid city name, got London,",
	},
	{
		input:             `London west=Boston*`,
		expectedErrorText: "testing:1:20: expected a valid city name as a mapDirection value, got Boston*",
	},
	{
		input:             `London west=`,
		expectedErrorText: "testing:1:13: Unexpected EOF",
	},
	{
		input:             `London east`,
		expectedErrorText: "testing:1:12: Unexpected EOF",
	},
	{
		input:             `London west=Boston west=Chicago`,
		expectedErrorText: "testing:1:24: got mapDirection type duplication west for city London",
	},
	{
		input:             `London west=Boston east=Boston`,
		expectedErrorText: "testing:1:31: got mapDirection value duplication Boston for city London",
	},
	{
		input: `London west=Boston 
London west=Bolton`,
		expectedErrorText: "testing:2:7: got city duplication for London previously declared on line 1",
	},
	{
		input:             `London nowhere=Boston`,
		expectedErrorText: "testing:1:15: got unexpected mapDirection type nowhere, expected one of south,north,west,east",
	},
	{
		input:             `London west 12`,
		expectedErrorText: "testing:1:15: unexpected token 12, expected =",
	},
}

func TestParserErrorHandling(t *testing.T) {
	for _, tc := range errTestCases {
		prsr := newParser(strings.NewReader(tc.input), "testing")
		err := prsr.parse()
		require.EqualError(t, err, tc.expectedErrorText)
	}
}

func TestHandleNonExistentDirectionValue(t *testing.T) {
	input := `London west=Boston east=Bolton
Bolton west=London `
	expectedError := "city London on line 1 has direction west which points to non existent city Boston"
	prsr := newParser(strings.NewReader(input), "testing")
	err := prsr.parse()
	require.Nil(t, err)
	err = prsr.checkDirectionValuesExistence()
	require.EqualError(t, err, expectedError)
}

func TestFailNonExistentCityCheckOnNonParsedParser(t *testing.T) {
	input := `London east=Bolton
Bolton west=London`
	expectedError := "cannot check direction value for unparsed file"
	prsr := newParser(strings.NewReader(input), "testing")
	err := prsr.checkDirectionValuesExistence()
	require.EqualError(t, err, expectedError)
}

func TestFailToBuildSimulationOnNonParserParser(t *testing.T) {
	input := `London east=Bolton
Bolton west=London`
	expectedError := "cannot build simulation for unparsed file"
	prsr := newParser(strings.NewReader(input), "testing")
	_, err := prsr.buildSimulation()
	require.EqualError(t, err, expectedError)
}

func TestSuccessfullyParseValidInputAndBuildSimulation(t *testing.T) {
	input := `London east=Bolton
Bolton west=London `
	expected := &Simulation{
		initialMap: map[string]*city{
			"London": {
				name:        "London",
				isDestroyed: false,
				directions:  []mapDirection{{directionValue: "Bolton", directionType: "east"}},
			},
			"Bolton": {
				name:        "Bolton",
				isDestroyed: false,
				directions:  []mapDirection{{directionValue: "London", directionType: "west"}},
			},
		},
	}
	prsr := newParser(strings.NewReader(input), "testing")
	err := prsr.parse()
	require.Nil(t, err)
	err = prsr.checkDirectionValuesExistence()
	require.Nil(t, err)
	simulation, err := prsr.buildSimulation()
	require.Nil(t, err)
	require.EqualValues(t, expected, simulation)
}
