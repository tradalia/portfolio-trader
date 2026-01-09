//=============================================================================
/*
Copyright © 2023 Andrea Carboni andrea.carboni71@gmail.com

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

package service

import (
	"github.com/tradalia/core/auth"
	"github.com/tradalia/portfolio-trader/pkg/business"
	"github.com/tradalia/portfolio-trader/pkg/business/filter"
	"github.com/tradalia/portfolio-trader/pkg/business/performance"
	"github.com/tradalia/portfolio-trader/pkg/business/quality"
	"github.com/tradalia/portfolio-trader/pkg/business/simulation"
	"github.com/tradalia/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func getTradingSystems(c *auth.Context) {
	filter := map[string]any{}
	offset, limit, err := c.GetPagingParams()

	if err == nil {
		details, err := c.GetParamAsBool("details", false)

		if err == nil {
			err = db.RunInTransaction(func(tx *gorm.DB) error {
				list, err := business.GetTradingSystems(tx, c, filter, offset, limit, details)

				if err != nil {
					return err
				}

				return c.ReturnList(list, offset, limit, len(*list))
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func getTradingSystem(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		err = db.RunInTransaction(func(tx *gorm.DB) error {
			ts, err2 := business.GetTradingSystem(tx, c, tsId)

			if err2 != nil {
				return err2
			}

			return c.ReturnObject(&ts)
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func getTrades(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		err = db.RunInTransaction(func(tx *gorm.DB) error {
			list, err := business.GetTrades(tx, c, tsId)

			if err != nil {
				return err
			}

			return c.ReturnObject(&list)
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func deleteTrades(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		err = db.RunInTransaction(func(tx *gorm.DB) error {
			err = business.DeleteTrades(tx, c, tsId)

			if err != nil {
				return err
			}

			return c.ReturnObject("")
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func getTradingFilters(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		err = db.RunInTransaction(func(tx *gorm.DB) error {
			filters, err := business.GetTradingFilters(tx, c, tsId)

			if err != nil {
				return err
			}

			return c.ReturnObject(filters)
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func setTradingFilters(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		filters := filter.TradingFilter{}
		err = c.BindParamsFromBody(&filters)

		if err == nil {
			err = db.RunInTransaction(func(tx *gorm.DB) error {
				err = business.SetTradingFilters(tx, c, tsId, &filters)

				if err != nil {
					return err
				}

				return c.ReturnObject("")
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func runFilterAnalysis(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		req := filter.AnalysisRequest{}
		err = c.BindParamsFromBody(&req)

		if err == nil {
			err = db.RunInTransaction(func(tx *gorm.DB) error {
				rep, err := business.RunFilterAnalysis(tx, c, tsId, &req)

				if err != nil {
					return err
				}

				return c.ReturnObject(rep)
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func runPerformanceAnalysis(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		req := performance.AnalysisRequest{}
		err = c.BindParamsFromBody(&req)

		if err == nil {
			err = db.RunInTransaction(func(tx *gorm.DB) error {
				res, err := business.RunPerformanceAnalysis(tx, c, tsId, &req)

				if err != nil {
					return err
				}

				return c.ReturnObject(res)
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func runQualityAnalysis(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		req := quality.AnalysisRequest{}
		err = c.BindParamsFromBody(&req)

		if err == nil {
			err = db.RunInTransaction(func(tx *gorm.DB) error {
				res, err := business.RunQualityAnalysis(tx, c, tsId, &req)

				if err != nil {
					return err
				}

				return c.ReturnObject(res)
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================
//===
//=== Filter optimization
//===
//=============================================================================

func startFilterOptimization(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		req := filter.OptimizationRequest{}
		err = c.BindParamsFromBody(&req)

		if err == nil {
			err = db.RunInTransaction(func(tx *gorm.DB) error {
				err := business.StartFilterOptimization(tx, c, tsId, &req)

				if err != nil {
					return err
				}

				return c.ReturnObject(NewStatusOkResponse())
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func stopFilterOptimization(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		err = business.StopFilterOptimization(c, tsId)

		if err != nil {
			c.ReturnError(err)
		} else {
			_ = c.ReturnObject(NewStatusOkResponse())
		}

		return
	}

	c.ReturnError(err)
}

//=============================================================================

func getFilterOptimizationInfo(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		res, err := business.GetFilterOptimizationInfo(c, tsId)

		if err != nil {
			c.ReturnError(err)
		} else {
			_ = c.ReturnObject(res)
		}

		return
	}

	c.ReturnError(err)
}

//=============================================================================
//===
//=== Simulation
//===
//=============================================================================

func startSimulation(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		req := simulation.Request{}
		err = c.BindParamsFromBody(&req)

		if err == nil {
			err = db.RunInTransaction(func(tx *gorm.DB) error {
				err2 := business.StartSimulation(tx, c, tsId, &req)

				if err2 != nil {
					return err2
				}

				return c.ReturnObject(NewStatusOkResponse())
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func stopSimulation(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		_ = business.StopSimulation(c, tsId)
		_ = c.ReturnObject(NewStatusOkResponse())
	}

	c.ReturnError(err)
}

//=============================================================================

func getSimulationResult(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		res := business.GetSimulationResult(c, tsId)
		_ = c.ReturnObject(res)
		return
	}

	c.ReturnError(err)
}

//=============================================================================
