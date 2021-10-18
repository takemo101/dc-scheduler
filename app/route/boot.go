package route

import (
	"github.com/takemo101/dc-scheduler/core"
	"github.com/takemo101/dc-scheduler/core/contract"
	"go.uber.org/fx"
)

// Module ルーティングモジュールDI
var Module = fx.Options(
	fx.Provide(NewApiRoute),
	fx.Provide(NewAdminRoute),
	fx.Provide(NewRoute),
)

// Routes ルーティングコレクション
type Routes []contract.Route

// NewRoute ルーティング 初期化処理
func NewRoute(
	config core.Config,
	api ApiRoute,
	admin AdminRoute,
) Routes {
	return Routes{
		api,
		admin,
	}
}

// Setup ルーティングコレクションの初期化処理実行
func (routes Routes) Setup() {
	for _, route := range routes {
		route.Setup()
	}
}
