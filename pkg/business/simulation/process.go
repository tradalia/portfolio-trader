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

package simulation

import (
	"encoding/base64"
	"log/slog"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-analyze/charts"
	"github.com/tradalia/core/datatype"
	"github.com/tradalia/portfolio-trader/pkg/core"
	"github.com/tradalia/portfolio-trader/pkg/db"
	"golang.org/x/exp/stats"
)

//=============================================================================
//===
//=== Simulation process
//===
//=============================================================================

type Process struct {
	ts       *db.TradingSystem
	trades   *[]db.Trade
	req      *Request
	risk     float64
	result   *Result
	stopping bool
}

//=============================================================================

const SimStatusIdle     = "idle"
const SimStatusWaiting  = "waiting"
const SimStatusRunning  = "running"
const SimStatusComplete = "complete"

//=============================================================================

func NewProcess(ts *db.TradingSystem, trades *[]db.Trade, req *Request, risk float64) *Process {
	return &Process{
		ts    : ts,
		trades: trades,
		req   : req,
		risk  : risk,
		result: &Result{
			Status: SimStatusWaiting,
		},
	}
}

//=============================================================================

func (p *Process) Start() {
	slog.Info("SimulationProcess: Starting","id", p.ts.Id)

	p.result = NewResult(p.GetFirstTradeDate(), p.GetLastTradeDate(), p.req.Runs, p.req.InitialCapital, p.req.RuinPercentage, p.risk)
	p.result.Status    = SimStatusRunning
	p.result.StartTime = time.Now()

	rMultGrossAll   := core.CalcRMultiple(p.trades, db.TradeTypeAll,   p.risk, 0)
	rMultGrossLong  := core.CalcRMultiple(p.trades, db.TradeTypeLong,  p.risk, 0)
	rMultGrossShort := core.CalcRMultiple(p.trades, db.TradeTypeShort, p.risk, 0)
	rMultNetAll     := core.CalcRMultiple(p.trades, db.TradeTypeAll,   p.risk, p.ts.CostPerOperation)
	rMultNetLong    := core.CalcRMultiple(p.trades, db.TradeTypeLong,  p.risk, p.ts.CostPerOperation)
	rMultNetShort   := core.CalcRMultiple(p.trades, db.TradeTypeShort, p.risk, p.ts.CostPerOperation)

	p.result.GrossAll = run(rMultGrossAll, p.req)
	p.result.Step++
	if !p.stopping {
		p.result.GrossLong = run(rMultGrossLong, p.req)
		p.result.Step++
		if !p.stopping {
			p.result.GrossShort = run(rMultGrossShort, p.req)
			p.result.Step++
			if !p.stopping {
				p.result.NetAll = run(rMultNetAll, p.req)
				p.result.Step++
				if !p.stopping {
					p.result.NetLong = run(rMultNetLong, p.req)
					p.result.Step++
					if !p.stopping {
						p.result.NetShort = run(rMultNetShort, p.req)
						p.result.Step++
					}
				}
			}
		}
	}

	p.result.Status  = SimStatusComplete
	p.result.EndTime = time.Now()
	slog.Info("SimulationProcess: Ended", "id", p.ts.Id)
}

//=============================================================================

func (p *Process) Stop() {
	p.stopping = true
}

//=============================================================================

func (p *Process) GetResult() *Result {
	return p.result
}

//=============================================================================

func (p *Process) GetFirstTradeDate() datatype.IntDate {
	t := *p.trades
	return datatype.ToIntDate(t[0].ExitDate)
}

//=============================================================================

func (p *Process) GetLastTradeDate() datatype.IntDate {
	t := *p.trades
	return datatype.ToIntDate(t[len(*p.trades) -1].ExitDate)
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func run(list []float64, req *Request) *Details {
	size := len(list)
	if size == 0 {
		return &Details{}
	}

	sampleSet, maxDrawdowns := buildSampleSet(list, req.Runs)
	sampleSet = addMeanAndStdDev(sampleSet, size)

	p, err := buildChart(sampleSet, req.Width, req.Height)
	if err != nil {
		panic(err)
	}

	buf, err := p.Bytes()
	if err != nil {
		panic(err)
	}

	return &Details{
		Equities    : base64.StdEncoding.EncodeToString(buf),
		MaxDrawdowns: buildDDDistrib(maxDrawdowns),
	}
}

//=============================================================================

func buildSampleSet(list []float64, runs int) ([][]float64, []float64) {
	size := len(list)

	var sampleSet    [][]float64
	var maxDrawdowns []float64

	for i:=0; i<runs; i++ {
		var sample = make([]float64, size)

		for j:=0; j<size; j++ {
			value := list[rand.Intn(size)]
			sample[j] = value
		}

		equity := core.BuildEquity(&sample)
		sampleSet = append(sampleSet, *equity)

		_, maxDD := core.BuildDrawDown(equity)
		maxDrawdowns = append(maxDrawdowns, maxDD)
	}

	return sampleSet, maxDrawdowns
}

//=============================================================================

func addMeanAndStdDev(sampleSet [][]float64, size int) [][]float64 {
	mean, stdDev := buildMeanAndStdDev(sampleSet, size)
	up,down := buildUpAndLowStdDev(mean, stdDev)

	sampleSet = append(sampleSet, make([]float64, size))
	sampleSet = append(sampleSet, mean)
	sampleSet = append(sampleSet, up)
	sampleSet = append(sampleSet, down)

	return sampleSet
}

//=============================================================================

func buildMeanAndStdDev(sampleSet [][]float64, size int) ([]float64, []float64) {
	var mean   []float64
	var stdDev []float64

	for i:=0; i<size; i++ {
		var serie []float64

		for j:=0; j<len(sampleSet); j++ {
			serie = append(serie, sampleSet[j][i])
		}

		m, s := stats.MeanAndStdDev(serie)

		mean   = append(mean,   m)
		stdDev = append(stdDev, s)
	}

	return mean, stdDev
}

//=============================================================================

func buildUpAndLowStdDev(mean, stdDev []float64) ([]float64, []float64) {
	var upStdDev   []float64
	var downStdDev []float64

	for i, v := range mean {
		upStdDev   = append(upStdDev,   v + stdDev[i])
		downStdDev = append(downStdDev, v - stdDev[i])
	}

	return upStdDev, downStdDev
}

//=============================================================================

func buildDDDistrib(data []float64) *Distribution {
	minv := core.CalcMin(data)
	size := int(math.Trunc(math.Abs(minv))) +1

	var xAxis []string
	for i := 1; i <=size; i++ {
		xAxis = append(xAxis, strconv.Itoa(-size + i) +"R")
	}

	yAxis := make([]float64, size)

	for _, value := range data {
		index := size + int(math.Trunc(value)) -1
		yAxis[index]++
	}

	return &Distribution{
		XAxis: xAxis,
		YAxis: yAxis,
	}
}

//=============================================================================
//=== Chart building
//=============================================================================

func buildChart(sampleSet [][]float64, width, height int) (*charts.Painter, error) {
	xAxis := calcXAxis(len(sampleSet[0]))

	opt := charts.NewLineChartOptionWithData(sampleSet)
	opt.XAxis.Title  = "Trades"
	opt.XAxis.Labels = xAxis
	opt.XAxis.LabelFontStyle.FontSize = 8
	opt.YAxis[0].Title = "Cumulative R multiples"
	opt.YAxis[0].LabelFontStyle.FontSize = 8
	opt.LineStrokeWidth = 1
	opt.Theme = opt.Theme.WithSeriesColors(buildColors(len(sampleSet)))

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width       : width,
		Height      : height,
	})
	err := p.LineChart(opt)

	return p, err
}

//=============================================================================

func calcXAxis(size int) []string {
	var axis = make([]string, size)

	for i:=1; i<=size; i++ {
		axis[i-1] = strconv.Itoa(i)
	}

	return axis
}

//=============================================================================

func buildColors(size int) []charts.Color {
	var list []charts.Color

	for i:=0; i<size-4; i++ {
		list = append(list, charts.Color{ 192, 192, 192, 96 })
	}

	list = append(list, charts.Color{ 128, 128, 128, 255 })
	list = append(list, charts.Color{  16,  16,  16, 255 })
	list = append(list, charts.Color{  80,  80,  80, 255 })
	list = append(list, charts.Color{  80,  80,  80, 255 })

	return list
}

//=============================================================================
