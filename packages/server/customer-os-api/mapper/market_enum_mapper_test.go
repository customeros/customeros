package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"testing"
)

func TestMapMarketFromModel(t *testing.T) {
	testCases := []struct {
		input          *model.Market
		expectedOutput string
	}{
		{input: nil, expectedOutput: ""},
		{input: utils.ToPtr(model.MarketB2b), expectedOutput: "B2B"},
		{input: utils.ToPtr(model.MarketB2c), expectedOutput: "B2C"},
		{input: utils.ToPtr(model.MarketMarketplace), expectedOutput: "Marketplace"},
	}

	for _, testCase := range testCases {
		output := MapMarketFromModel(testCase.input)
		if output != testCase.expectedOutput {
			t.Errorf("Expected output: %s, but got: %s", testCase.expectedOutput, output)
		}
	}
}

func TestMapMarketToModel(t *testing.T) {
	testCases := []struct {
		input          string
		expectedOutput *model.Market
	}{
		{input: "B2B", expectedOutput: utils.ToPtr(model.MarketB2b)},
		{input: "B2C", expectedOutput: utils.ToPtr(model.MarketB2c)},
		{input: "Marketplace", expectedOutput: utils.ToPtr(model.MarketMarketplace)},
		{input: "Invalid", expectedOutput: nil},
	}

	for _, testCase := range testCases {
		output := MapMarketToModel(testCase.input)
		if testCase.expectedOutput == nil && output != nil {
			t.Errorf("Expected output: %v, but got: %v", testCase.expectedOutput, output)
		} else if testCase.expectedOutput != nil && output.String() != testCase.expectedOutput.String() {
			t.Errorf("Expected output: %v, but got: %v", testCase.expectedOutput, output)
		}
	}
}
