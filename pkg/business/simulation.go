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

package business

import (
	"github.com/tradalia/core/auth"
	"github.com/tradalia/core/req"
	"github.com/tradalia/portfolio-trader/pkg/business/simulation"
	"github.com/tradalia/portfolio-trader/pkg/core"
	"github.com/tradalia/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func StartSimulation(tx *gorm.DB, c *auth.Context, tsId uint, rq *simulation.Request) error {
	//--- Get trading system

	ts, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return err
	}

	c.Log.Info("StartSimulation: Starting", "id", tsId, "name", ts.Name, "runs", rq.Runs)

	fromTime := calcBackPeriod(rq.DaysBack)

	trades, err := db.FindTradesByTsIdFromTime(tx, ts.Id, fromTime, nil)
	if err != nil {
		return err
	}

	if len(*trades) == 0 {
		return req.NewUnprocessableEntityError("no trades found for given time")
	}

	risk,err := core.CalcRisk(trades)
	if err != nil {
		return err
	}

	simulation.Start(rq, ts, trades, risk)
	c.Log.Info("StartSimulation: Ending", "id", tsId, "name", ts.Name, "runs", rq.Runs)
	return nil
}

//=============================================================================

func StopSimulation(c *auth.Context, tsId uint) bool {
	return simulation.Stop(tsId)
}

//=============================================================================

func GetSimulationResult(c *auth.Context, tsId uint) *simulation.Result {
	return simulation.GetResult(tsId)
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================
