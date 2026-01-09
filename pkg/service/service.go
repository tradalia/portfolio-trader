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
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/tradalia/core/auth"
	"github.com/tradalia/core/auth/roles"
	"github.com/tradalia/core/req"
	"github.com/tradalia/portfolio-trader/pkg/app"
)

//=============================================================================

func Init(router *gin.Engine, cfg *app.Config, logger *slog.Logger) {

	ctrl := auth.NewOidcController(cfg.Authentication.Authority, req.GetClient("bf"), logger, cfg)

	router.GET   ("/api/portfolio/v1/trading-systems",                         ctrl.Secure(getTradingSystems,         roles.Admin_User_Service))
	router.GET   ("/api/portfolio/v1/trading-systems/:id",                     ctrl.Secure(getTradingSystem,          roles.Admin_User_Service))
	router.GET   ("/api/portfolio/v1/trading-systems/:id/trades",              ctrl.Secure(getTrades,                 roles.Admin_User_Service))
	router.DELETE("/api/portfolio/v1/trading-systems/:id/trades",              ctrl.Secure(deleteTrades,              roles.Admin_User_Service))
	router.GET   ("/api/portfolio/v1/trading-systems/:id/filters",             ctrl.Secure(getTradingFilters,         roles.Admin_User_Service))
	router.POST  ("/api/portfolio/v1/trading-systems/:id/filters",             ctrl.Secure(setTradingFilters,         roles.Admin_User_Service))
	router.POST  ("/api/portfolio/v1/trading-systems/:id/filter-analysis",     ctrl.Secure(runFilterAnalysis,         roles.Admin_User_Service))
	router.POST  ("/api/portfolio/v1/trading-systems/:id/trading",             ctrl.Secure(setTradingSystemTrading,   roles.Admin_User_Service))
	router.POST  ("/api/portfolio/v1/trading-systems/:id/running",             ctrl.Secure(setTradingSystemRunning,   roles.Admin_User_Service))
	router.POST  ("/api/portfolio/v1/trading-systems/:id/activation",          ctrl.Secure(setTradingSystemActivation,roles.Admin_User_Service))
	router.POST  ("/api/portfolio/v1/trading-systems/:id/active",              ctrl.Secure(setTradingSystemActive,    roles.Admin_User_Service))
	router.POST  ("/api/portfolio/v1/trading-systems/:id/performance-analysis",ctrl.Secure(runPerformanceAnalysis,    roles.Admin_User_Service))
	router.POST  ("/api/portfolio/v1/trading-systems/:id/quality-analysis",    ctrl.Secure(runQualityAnalysis,        roles.Admin_User_Service))

	router.GET   ("/api/portfolio/v1/trading-systems/:id/filter-optimization", ctrl.Secure(getFilterOptimizationInfo, roles.Admin_User_Service))
	router.POST  ("/api/portfolio/v1/trading-systems/:id/filter-optimization", ctrl.Secure(startFilterOptimization,   roles.Admin_User_Service))
	router.DELETE("/api/portfolio/v1/trading-systems/:id/filter-optimization", ctrl.Secure(stopFilterOptimization,    roles.Admin_User_Service))

	router.GET   ("/api/portfolio/v1/trading-systems/:id/simulation",          ctrl.Secure(getSimulationResult,       roles.Admin_User_Service))
	router.POST  ("/api/portfolio/v1/trading-systems/:id/simulation",          ctrl.Secure(startSimulation,           roles.Admin_User_Service))
	router.DELETE("/api/portfolio/v1/trading-systems/:id/simulation",          ctrl.Secure(stopSimulation,            roles.Admin_User_Service))

	router.GET   ("/api/inventory/v1/portfolios",                              ctrl.Secure(getPortfolios,             roles.Admin_User_Service))
	router.GET   ("/api/inventory/v1/portfolio/tree",                          ctrl.Secure(getPortfolioTree,          roles.Admin_User_Service))
	router.POST  ("/api/portfolio/v1/portfolio/monitoring",                    ctrl.Secure(getPortfolioMonitoring,    roles.Admin_User_Service))
}

//=============================================================================

func NewStatusOkResponse() any {
	return struct {
		status string
	}{ status: "ok"}
}

//=============================================================================
