package route

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/takemo101/dc-scheduler/app/controller/admin"
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/app/middleware"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/core"
)

// AdminRoute is struct
type AdminRoute struct {
	logger              core.Logger
	app                 core.Application
	path                core.Path
	csrf                middleware.Csrf
	value               support.RequestValue
	render              middleware.ViewRender
	authController      controller.SessionAuthController
	dashboardController controller.DashboardController
}

// Setup is setup route
func (r AdminRoute) Setup() {
	r.logger.Info("setup admin-route")

	app := r.app.App

	// root redirect
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect(r.path.URL("system"))
	})

	// admin http route
	http := app.Group(
		"/system",
		r.csrf.CreateHandler("form:csrf_token"),
		r.render.CreateHandler(r.ViewRenderCreateHandler),
	)
	{
		// auth login route
		auth := http.Group(
			"/auth",
		)
		{
			auth.Get("/login", r.authController.LoginForm)
			auth.Post("/login", r.authController.Login)
		}

		// after login route
		system := http.Group(
			"/",
		)
		{
			// dashboard route
			system.Get("/", r.dashboardController.Dashboard)

			// auth logout route
			system.Post("/logout", r.authController.Logout)
		}
	}
}

// NewAdminRoute create new admin route
func NewAdminRoute(
	logger core.Logger,
	app core.Application,
	path core.Path,
	csrf middleware.Csrf,
	value support.RequestValue,
	render middleware.ViewRender,
	authController controller.SessionAuthController,
	dashboardController controller.DashboardController,
) AdminRoute {
	return AdminRoute{
		logger:              logger,
		app:                 app,
		path:                path,
		csrf:                csrf,
		value:               value,
		render:              render,
		authController:      authController,
		dashboardController: dashboardController,
	}
}

// ViewRenderCreateHandler middleware handler
func (r AdminRoute) ViewRenderCreateHandler(c *fiber.Ctx, vr *helper.ViewRender) {
	// load request list
	adminlte, err := r.app.Config.Load("admin-lte")
	if err == nil {
		for k, v := range map[string]string{
			"adminlte_menus":   "menus",
			"adminlte_plugins": "plugins",
		} {
			if value, ok := adminlte[v]; ok {
				vr.SetData(helper.DataMap{
					k: value,
				})
			}
		}
	}

	// load meta
	meta, err := r.app.Config.Load("admin-meta")
	if err == nil {
		vr.SetData(helper.DataMap{
			"admin_meta": meta,
		})
	}

	// csrf-token
	csrfToken := middleware.GetCSRFToken(c)
	// session errors
	errors, _ := middleware.GetSessionErrors(c)
	// session inputs
	inputs, _ := middleware.GetSessionInputs(c)
	// session messages
	messages, _ := middleware.GetSessionMessages(c)

	vr.SetData(helper.DataMap{
		"csrf_token": csrfToken,
		"errors":     errors,
		"inputs":     inputs,
		"messages":   messages,
	})
	vr.SetJS(helper.DataMap{
		"csrfToken": csrfToken,
		"errors":    errors,
		"inputs":    inputs,
		"messages":  messages,
	})
}
