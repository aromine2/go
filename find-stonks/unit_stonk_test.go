package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAdjustCurrentPE(t *testing.T) {
	var tests = []struct {
		companyOverviewObject map[string]string
		expected              float64
	}{
		{map[string]string{"PERatio": "2"}, 2 * .8},
		{map[string]string{"PERatio": "123"}, 123 * .8},
		{map[string]string{"PERatio": "300"}, 300 * .8},
	}

	for n, p := range tests {
		testName := fmt.Sprintf("%d", n+1)
		t.Run(testName, func(t *testing.T) {
			actualResult := AdjustCurrentPE(p.companyOverviewObject)
			if !reflect.DeepEqual(actualResult, p.expected) {
				t.Errorf("Expected %v but got %v", p.expected, actualResult)
			}
		})
	}
}

func TestFiveYearMetricProjection(t *testing.T) {
	var tests = []struct {
		metric           string
		percentChange    float64
		financialsObject []map[string]string
		expected         float64
	}{
		{"totalRevenue", .69, []map[string]string{0: {"totalRevenue": "6969696969.00"}}, 78006.1727301504},
		{"totalRevenue", 1.23, []map[string]string{0: {"totalRevenue": "9280208"}}, 385.3722053646542},
		{"totalRevenue", -.3, []map[string]string{0: {"totalRevenue": "9999280208"}}, 2072.9223749651924},
		{"commonStockSharesOutstanding", .05, []map[string]string{0: {"commonStockSharesOutstanding": "8675309"}}, 11.338281906243306},
		{"commonStockSharesOutstanding", .75, []map[string]string{0: {"commonStockSharesOutstanding": "875309"}}, 17.720490458957403},
		{"commonStockSharesOutstanding", -.3, []map[string]string{0: {"commonStockSharesOutstanding": "10000000"}}, 1.3501251069999993},
	}

	for n, p := range tests {
		testName := fmt.Sprintf("%d", n+1)
		t.Run(testName, func(t *testing.T) {
			actualResult := FiveYearMetricProjection(p.metric, p.percentChange, p.financialsObject)
			if !reflect.DeepEqual(actualResult, p.expected) {
				t.Errorf("Expected %v but got %v", p.expected, actualResult)
			}
		})
	}

}

// TODO: Mock out FiveYearMetricProjection
func TestCalculateAverageYearlyMetricGrowth(t *testing.T) {
	testFinancialResponse := FinancialsResponse{Symbol: "AJR", AnnualReports: []map[string]string{
		0: {"totalRevenue": "199999999", "commonStockSharesOutstanding": "199999990"},
		1: {"totalRevenue": "155555555", "commonStockSharesOutstanding": "177777770"},
		2: {"totalRevenue": "111111111", "commonStockSharesOutstanding": "144444440"},
		3: {"totalRevenue": "111111111", "commonStockSharesOutstanding": "111111110"},
	},
		QuarterlyReports: []map[string]string{}}

	var tests = []struct {
		metric           string
		financialsObject FinancialsResponse
		expected         float64
	}{
		{"totalRevenue", testFinancialResponse, 509.6265244984146},
		{"commonStockSharesOutstanding", testFinancialResponse, 587.3865729130903},
	}

	for n, p := range tests {
		testName := fmt.Sprintf("%d", n+1)
		t.Run(testName, func(t *testing.T) {
			actualResult := CalculateAverageYearlyMetricGrowth(p.metric, p.financialsObject)
			if !reflect.DeepEqual(actualResult, p.expected) {
				t.Errorf("Expected %v but got %v", p.expected, actualResult)
			}
		})
	}
}

func TestCalculateAverageMargin(t *testing.T) {
	testFinancialResponse := FinancialsResponse{Symbol: "AJR", AnnualReports: []map[string]string{
		0: {"totalRevenue": "199999999", "netIncome": "19999999"},
		1: {"totalRevenue": "155555555", "netIncome": "17777777"},
		2: {"totalRevenue": "111111111", "netIncome": "14444444"},
		3: {"totalRevenue": "111111111", "netIncome": "11111111"},
	},
		QuarterlyReports: []map[string]string{}}

	var tests = []struct {
		incomeStatement FinancialsResponse
		expected         float64
	}{
		{testFinancialResponse, 0.11476190044129249},
	}

	for n, p := range tests {
		testName := fmt.Sprintf("%d", n+1)
		t.Run(testName, func(t *testing.T) {
			actualResult := CalculateAverageMargin(p.incomeStatement)
			if !reflect.DeepEqual(actualResult, p.expected) {
				t.Errorf("Expected %v but got %v", p.expected, actualResult)
			}
		})
	}
}
