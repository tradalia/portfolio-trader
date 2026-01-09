//=============================================================================
/*
Copyright © 2025 Andrea Carboni andrea.carboni71@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
//=============================================================================

package quality

import (
	"math"

	"github.com/tradalia/core/datatype"
	"github.com/tradalia/portfolio-trader/pkg/core"
	"github.com/tradalia/portfolio-trader/pkg/db"
	"github.com/tradalia/portfolio-trader/pkg/platform"
	"golang.org/x/exp/stats"
)

//=============================================================================

const (
	DirectionAll  = 3
	VolatilityAll = 4
)

//=============================================================================

func GetQualityAnalysis(ts *db.TradingSystem, trades *[]db.Trade, man *platform.DataProductAnalysisResponse) (*AnalysisResponse, error) {
	res := NewAnalysisResponse()
	res.TradingSystem = ts

	risk, err := core.CalcRisk(trades)
	if err != nil {
		return nil, err
	}

	marketMap := buildMarketMap(man.DailyResults)

	//--- Calc all standard metrics

	for dir := platform.DirectionStrongBear; dir <= platform.DirectionStrongBull; dir++ {
		for vol := platform.VolatilityQuiet; vol <= platform.VolatilityVeryVolatile; vol++ {
			calcQualityCell(res, trades, dir, vol, risk, ts.CostPerOperation, marketMap)
		}
	}

	//--- Calc summary by direction

	for dir := platform.DirectionStrongBear; dir <= platform.DirectionStrongBull; dir++ {
		calcQualityCell(res, trades, dir, VolatilityAll, risk, ts.CostPerOperation, marketMap)
	}

	//--- Calc summary by volatility

	for vol := platform.VolatilityQuiet; vol <= platform.VolatilityVeryVolatile; vol++ {
		calcQualityCell(res, trades, DirectionAll, vol, risk, ts.CostPerOperation, marketMap)
	}

	//--- Calc overall

	calcQualityCell(res, trades, DirectionAll, VolatilityAll, risk, ts.CostPerOperation, marketMap)

	return res, nil
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func buildMarketMap(list []*platform.DailyResult) map[datatype.IntDate]*platform.DailyResult {
	res := make(map[datatype.IntDate]*platform.DailyResult)

	for _,dr := range list {
		res[dr.Date] = dr
	}

	return res
}

//=============================================================================

func calcQualityCell(res *AnalysisResponse, trades *[]db.Trade, dir int, vol int, risk float64, costPerOperation float64, marketMap map[datatype.IntDate]*platform.DailyResult) {
	res.QualityAllGross  [dir+2][vol] = calcQualityMetrics(trades, db.TradeTypeAll,   dir, vol, risk, 0, marketMap)
	res.QualityLongGross [dir+2][vol] = calcQualityMetrics(trades, db.TradeTypeLong,  dir, vol, risk, 0, marketMap)
	res.QualityShortGross[dir+2][vol] = calcQualityMetrics(trades, db.TradeTypeShort, dir, vol, risk, 0, marketMap)

	res.QualityAllNet    [dir+2][vol] = calcQualityMetrics(trades, db.TradeTypeAll,   dir, vol, risk, costPerOperation, marketMap)
	res.QualityLongNet   [dir+2][vol] = calcQualityMetrics(trades, db.TradeTypeLong,  dir, vol, risk, costPerOperation, marketMap)
	res.QualityShortNet  [dir+2][vol] = calcQualityMetrics(trades, db.TradeTypeShort, dir, vol, risk, costPerOperation, marketMap)
}

//=============================================================================

func calcQualityMetrics(trades *[]db.Trade, tradeType string, direction int, volatility int, risk float64, costPerOper float64, marketMap map[datatype.IntDate]*platform.DailyResult) *Metrics {

	//--- Step 1: Collect relevant trades

	var list []float64

	for _, t := range *trades {
		if t.TradeType == tradeType || tradeType == db.TradeTypeAll {
			tradeDir, tradeVol := mapTrade(&t, marketMap)

			if direction == DirectionAll || direction == tradeDir {
				if volatility == VolatilityAll || volatility == tradeVol {
					returns := t.GrossProfit - 2 * costPerOper
					list = append(list, returns / risk)
				}
			}
		}
	}

	//--- Step 2: Calc metrics

	cell := &Metrics{
		Trades: len(list),
	}

	if len(list) > 0 {
		cell.TradesPerc = core.Trunc2d(100 * float64(len(list)) / float64(len(*trades)))
		calcMetrics(list, cell)
	}

	return cell
}

//=============================================================================

func mapTrade(trade *db.Trade, marketMap map[datatype.IntDate]*platform.DailyResult) (int, int) {
	date := datatype.ToIntDate(trade.EntryDate)
	dr, ok := marketMap[date]
	if !ok {
		return 10, 10
	}

	return dr.Direction, dr.Volatility
}

//=============================================================================

func calcMetrics(list []float64, cell *Metrics) {
	mean,stdd := stats.MeanAndStdDev(list)
	listLen   := float64(len(list))
	capLen    := math.Min(listLen, 100)

	if stdd > 0.0 {
		cell.Sqn    = core.Trunc2d(mean / stdd * math.Sqrt(listLen))
		cell.Sqn100 = core.Trunc2d(mean / stdd * math.Sqrt(capLen))
	}

	equity  := core.BuildEquity(&list)
	_,maxDD := core.BuildDrawDown(equity)

	cell.MaxDrawdown = core.Trunc2d(maxDD)
}

//=============================================================================
