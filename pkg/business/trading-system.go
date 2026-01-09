//=============================================================================
/*
Copyright © 2024 Andrea Carboni andrea.carboni71@gmail.com

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
	"github.com/tradalia/portfolio-trader/pkg/db"
	"github.com/tradalia/portfolio-trader/pkg/platform"
	"gorm.io/gorm"
)

//=============================================================================

func GetTradingSystems(tx *gorm.DB, c *auth.Context, filter map[string]any, offset int, limit int, details bool) (*[]db.TradingSystem, error) {
	if ! c.Session.IsAdmin() {
		filter["username"] = c.Session.Username
	}

	return db.GetTradingSystems(tx, filter, offset, limit)
}

//=============================================================================

func GetTradingSystem(tx *gorm.DB, c *auth.Context, id uint) (*db.TradingSystem, error) {
	ts,err := db.GetTradingSystemById(tx, id)
	if err != nil {
		return nil, err
	}

	if ts == nil {
		return nil, req.NewNotFoundError("trading system not found : %v", id)
	}

	if ! c.Session.IsAdmin() {
		if ts.Username != c.Session.Username {
			return nil, req.NewForbiddenError("user not allowed : %v", ts.Username)
		}
	}

	return ts, nil
}

//=============================================================================

func DeleteTradingSystem(tx *gorm.DB, id uint) error {
	err := db.DeleteAllTradesByTradingSystemId(tx, id)
	if err != nil {
		return err
	}

	err = db.DeleteAllDailyReturnsByTradingSystemId(tx, id)
	if err != nil {
		return err
	}

	err = db.DeleteTradingFilter(tx, id)
	if err != nil {
		return err
	}

	return db.DeleteTradingSystem(tx, id)
}

//=============================================================================

func GetTrades(tx *gorm.DB, c *auth.Context, id uint) (*[]db.Trade, error) {
	_, err := getTradingSystemAndCheckAccess(tx, c, id)
	if err != nil {
		return nil, err
	}

	return db.FindTradesByTradingSystemId(tx, id)
}

//=============================================================================

func DeleteTrades(tx *gorm.DB, c *auth.Context, id uint) error {
	c.Log.Info("DeleteTrades: Deleting all trades on trading system", "id", id)

	ts, err := getTradingSystemAndCheckAccess(tx, c, id)
	if err != nil {
		return err
	}

	err = db.DeleteAllTradesByTradingSystemId(tx, id)
	if err != nil {
		return err
	}

	err = db.DeleteAllDailyReturnsByTradingSystemId(tx, id)
	if err != nil {
		return err
	}

	ts.FirstTrade      = nil
	ts.LastTrade       = nil
	ts.LastNetProfit   = 0
	ts.LastNetAvgTrade = 0
	ts.LastNumTrades   = 0

	err = db.UpdateTradingSystem(tx, ts)
	if err != nil {
		return err
	}

	err = platform.DeleteEquityChart(c.Session.Username, id)
	if err != nil {

	}

	c.Log.Info("DeleteTrades: Operation ended", "id", id)
	return nil
}

//=============================================================================
