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

package performance

import (
	"time"

	"github.com/tradalia/core/datatype"
	"github.com/tradalia/portfolio-trader/pkg/core"
	"github.com/tradalia/portfolio-trader/pkg/core/stats"
	"github.com/tradalia/portfolio-trader/pkg/db"
)

//=============================================================================

func GetPerformanceAnalysis(ts *db.TradingSystem, trades *[]db.Trade, returns *[]db.DailyReturn) *AnalysisResponse {
	res := AnalysisResponse{}
	res.TradingSystem = ts
	res.Trades        = trades

	allEq  , allMaxGrossDD  , allMaxNetDD   := calcEquities(ts, trades, db.TradeTypeAll)
	longEq , longMaxGrossDD , longMaxNetDD  := calcEquities(ts, trades, db.TradeTypeLong)
	shortEq, shortMaxGrossDD, shortMaxNetDD := calcEquities(ts, trades, db.TradeTypeShort)

	res.AllEquities   = allEq
	res.LongEquities  = longEq
	res.ShortEquities = shortEq

	res.Gross.MaxDrawdown.Total = allMaxGrossDD
	res.Gross.MaxDrawdown.Long  = longMaxGrossDD
	res.Gross.MaxDrawdown.Short = shortMaxGrossDD
	res.Net  .MaxDrawdown.Total = allMaxNetDD
	res.Net  .MaxDrawdown.Long  = longMaxNetDD
	res.Net  .MaxDrawdown.Short = shortMaxNetDD

	res.Gross.Profit.Total = calcProfit(allEq  .GrossEquity)
	res.Gross.Profit.Long  = calcProfit(longEq .GrossEquity)
	res.Gross.Profit.Short = calcProfit(shortEq.GrossEquity)
	res.Net  .Profit.Total = calcProfit(allEq  .NetEquity)
	res.Net  .Profit.Long  = calcProfit(longEq .NetEquity)
	res.Net  .Profit.Short = calcProfit(shortEq.NetEquity)

	res.Gross.AverageTrade.Total = calcAvgTrade(res.Gross.Profit.Total, allEq.Trades)
	res.Gross.AverageTrade.Long  = calcAvgTrade(res.Gross.Profit.Long , longEq.Trades)
	res.Gross.AverageTrade.Short = calcAvgTrade(res.Gross.Profit.Short, shortEq.Trades)
	res.Net  .AverageTrade.Total = calcAvgTrade(res.Net  .Profit.Total, allEq.Trades)
	res.Net  .AverageTrade.Long  = calcAvgTrade(res.Net  .Profit.Long , longEq.Trades)
	res.Net  .AverageTrade.Short = calcAvgTrade(res.Net  .Profit.Short, shortEq.Trades)

	calcAggregates   (&res)
	updateGeneralInfo(&res)
	calcDistributions(&res, returns)
	calcRolling      (&res)

	return &res
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func calcEquities(ts *db.TradingSystem, trades *[]db.Trade, tradeType string) (*Equities, float64, float64) {
	timeSlice, grossProfits := core.BuildGrossProfits(trades, tradeType)
	netProfits              := core.BuildNetProfits(grossProfits, ts.CostPerOperation)

	grossEquity := core.BuildEquity(grossProfits)
	netEquity   := core.BuildEquity(netProfits)

	grossDD, maxGrossDD := core.BuildDrawDown(grossEquity)
	netDD  , maxNetDD   := core.BuildDrawDown(netEquity)

	return &Equities{
		Time         : timeSlice,
		GrossEquity  : grossEquity,
		NetEquity    : netEquity,
		GrossDrawdown: grossDD,
		NetDrawdown  : netDD,
		Trades       : len(*timeSlice),
	}, maxGrossDD, maxNetDD
}

//=============================================================================

func calcProfit(equity *[]float64) float64 {
	if equity == nil || len(*equity)==0 {
		return 0
	}

	return (*equity)[len(*equity) -1]
}

//=============================================================================

func calcAvgTrade(value float64, count int) float64 {
	if count == 0 {
		return 0
	}

	return core.Trunc2d(value / float64(count))
}

//=============================================================================
//=== Timezone shifting
//=============================================================================

func calcAggregates(res *AnalysisResponse) {
	calcYearAggregates(res)
}

//=============================================================================

func calcYearAggregates(res *AnalysisResponse) {
	cost := float64(res.TradingSystem.CostPerOperation)
	list := []*AnnualAggregate{}

	var currYear *AnnualAggregate

	for _, tr := range *res.Trades {
		if currYear == nil {
			//--- Beginning of a new year

			currYear = NewAggregate(&tr, cost)
			list = append(list, currYear)
		} else {
			if currYear.Year == tr.ExitDate.Year() {
				//--- Continue on the current year
				currYear.addTrade(&tr, cost)
			} else {
				//--- Continue on the new year

				currYear.consolidate()
				currYear = NewAggregate(&tr, cost)
				list = append(list, currYear)
			}
		}
	}

	if currYear != nil {
		currYear.consolidate()
	}

	res.Aggregates.Annual = &list
}

//=============================================================================
//=== General information
//=============================================================================

func updateGeneralInfo(res *AnalysisResponse) {
	calcFromToDates(res)
}

//=============================================================================

func calcFromToDates(res *AnalysisResponse) {
	numTrades := len(*res.Trades)
	if  numTrades> 0 {
		firstTrade := (*res.Trades)[0].ExitDate
		lastTrade  := (*res.Trades)[numTrades-1].ExitDate

		res.General.FromDate = NewIntDate(firstTrade)
		res.General.ToDate   = NewIntDate(lastTrade)
	}
}

//=============================================================================

func NewIntDate(t *time.Time) datatype.IntDate {
	return datatype.IntDate(t.Year()*10000 + int(t.Month())*100 + t.Day())
}

//=============================================================================
//=== Metrics
//=============================================================================

func calcDistributions(res *AnalysisResponse, returns *[]db.DailyReturn) {
	dist := &res.Distributions
	list := core.ToNonZeroDailyReturnSlice(returns)
	dist.Daily = calcDistribution(list)

	if dist.Daily != nil {
		dist.AnnualSharpeRatio = core.Trunc2d(dist.Daily.SharpeRatio * 16)
		dist.AnnualStandardDev = core.Trunc2d(dist.Daily.StandardDev * 16)
	}

	//--- All (gross + net)

	_, allGross := core.BuildGrossProfits(res.Trades, db.TradeTypeAll)
	allNet      := core.BuildNetProfits(allGross, res.TradingSystem.CostPerOperation)

	dist.TradesAllGross = calcDistribution(*allGross)
	dist.TradesAllNet   = calcDistribution(*allNet)

	//--- Long (gross + net)

	_, longGross := core.BuildGrossProfits(res.Trades, db.TradeTypeLong)
	longNet      := core.BuildNetProfits(longGross, res.TradingSystem.CostPerOperation)

	dist.TradesLongGross = calcDistribution(*longGross)
	dist.TradesLongNet   = calcDistribution(*longNet)

	//--- Short (gross + net)

	_, shortGross := core.BuildGrossProfits(res.Trades, db.TradeTypeShort)
	shortNet      := core.BuildNetProfits(shortGross, res.TradingSystem.CostPerOperation)

	dist.TradesShortGross = calcDistribution(*shortGross)
	dist.TradesShortNet   = calcDistribution(*shortNet)
}

//=============================================================================

func calcDistribution(data []float64) *Distribution {
	if data == nil || len(data) == 0 {
		return nil
	}

	mean    := stats.Mean(data)
	median  := stats.Median(data)
	stdDev  := stats.StdDev(data, mean)
	sharpeR := stats.SharpeRatio(mean, stdDev)
	skewness:= stats.Skewness(mean, median, stdDev)
	percen  := stats.NewPercentile(data)

	perc01 := percen.Get( 1) - mean
	perc30 := percen.Get(30) - mean
	perc70 := percen.Get(70) - mean
	perc99 := percen.Get(99) - mean

	lowerPercRatio := perc01 / perc30
	upperPercRatio := perc99 / perc70

	return &Distribution{
		Mean        : core.Trunc2d(mean),
		Median      : core.Trunc2d(median),
		StandardDev : core.Trunc2d(stdDev),
		SharpeRatio : core.Trunc2d(sharpeR),
		LowerTail   : core.Trunc2d(lowerPercRatio / 4.43),
		UpperTail   : core.Trunc2d(upperPercRatio / 4.43),
		Skewness    : core.Trunc2d(skewness),
		Histogram   : stats.NewHistogram(data),
	}
}

//=============================================================================

func calcRolling(res *AnalysisResponse) {
	costPerOper := res.TradingSystem.CostPerOperation

	for _, tr := range *res.Trades {
		year := tr.EntryDate.Year()
		dow  := int(tr.EntryDate.Weekday())
		mon  := int(tr.EntryDate.Month()) -1

		dowRI := &res.Rolling.Daily  [dow]
		monRI := &res.Rolling.Monthly[mon]

		updateRollingInfo(&tr, dowRI, costPerOper)
		updateRollingInfo(&tr, monRI, costPerOper)

		res.Rolling.DayYoY   = updateYoY(res.Rolling.DayYoY,   year, &tr, dow, costPerOper,  7)
		res.Rolling.MonthYoY = updateYoY(res.Rolling.MonthYoY, year, &tr, mon, costPerOper, 12)
	}
}

//=============================================================================

func updateRollingInfo(tr *db.Trade, ri *RollingInfo, costPerOper float64) {
	ri.Trades.Total++
	ri.GrossReturns.Total += tr.GrossProfit
	ri.NetReturns.Total   += tr.GrossProfit - 2 * costPerOper

	if tr.TradeType == db.TradeTypeLong {
		ri.Trades      .Long++
		ri.GrossReturns.Long += tr.GrossProfit
		ri.NetReturns  .Long += tr.GrossProfit - 2 * costPerOper
	} else {
		ri.Trades      .Short++
		ri.GrossReturns.Short += tr.GrossProfit
		ri.NetReturns  .Short += tr.GrossProfit - 2 * costPerOper
	}
}

//=============================================================================

func updateYoY(list []*YoYRolling, year int, tr *db.Trade, slot int, costPerOper float64, slots int) []*YoYRolling{
	var yoy *YoYRolling

	if list == nil || list[len(list)-1].Year != year {
		yoy = &YoYRolling{
			Year: year,
		}

		list = append(list, yoy)

		for i:=0;i<slots;i++ {
			yoy.Data = append(yoy.Data, &RollingInfo{})
		}
	} else {
		yoy = list[len(list)-1]
	}

	updateRollingInfo(tr, yoy.Data[slot], costPerOper)

	return list
}

//=============================================================================
