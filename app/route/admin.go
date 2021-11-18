package route

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/takemo101/dc-scheduler/app/controller/admin"
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/app/middleware"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/core"
)

// AdminRoute 管理者ルート
type AdminRoute struct {
	logger              core.Logger
	app                 core.Application
	path                core.Path
	sessionStore        core.SessionStore
	csrf                middleware.Csrf
	value               support.ContextValue
	auth                middleware.SessionAdminAuth
	role                middleware.SessionAdminRole
	render              middleware.ViewRender
	authController      controller.SessionAuthController
	dashboardController controller.DashboardController
	adminController     controller.AdminController
	userController      controller.UserController
	accountController   controller.AccountController
	botController       controller.BotController
	messageController   controller.PostMessageController
	immediateController controller.ImmediatePostController
	scheduleController  controller.SchedulePostController
	regularController   controller.RegularPostController
	timingController    controller.RegularTimingController
}

// NewAdminRoute コンストラクタ
func NewAdminRoute(
	logger core.Logger,
	app core.Application,
	path core.Path,
	sessionStore core.SessionStore,
	csrf middleware.Csrf,
	value support.ContextValue,
	auth middleware.SessionAdminAuth,
	role middleware.SessionAdminRole,
	render middleware.ViewRender,
	authController controller.SessionAuthController,
	dashboardController controller.DashboardController,
	adminController controller.AdminController,
	userController controller.UserController,
	accountController controller.AccountController,
	botController controller.BotController,
	messageController controller.PostMessageController,
	immediateController controller.ImmediatePostController,
	scheduleController controller.SchedulePostController,
	regularController controller.RegularPostController,
	timingController controller.RegularTimingController,
) AdminRoute {
	return AdminRoute{
		logger:              logger,
		app:                 app,
		path:                path,
		sessionStore:        sessionStore,
		csrf:                csrf,
		value:               value,
		auth:                auth,
		role:                role,
		render:              render,
		authController:      authController,
		dashboardController: dashboardController,
		adminController:     adminController,
		userController:      userController,
		accountController:   accountController,
		botController:       botController,
		messageController:   messageController,
		immediateController: immediateController,
		scheduleController:  scheduleController,
		regularController:   regularController,
		timingController:    timingController,
	}
}

// Setup ルートのセットアップ
func (r AdminRoute) Setup() {
	r.logger.Info("setup admin-route")

	app := r.app.App

	// root redirect
	/*
		app.Get("/", func(c *fiber.Ctx) error {
			return c.Redirect(r.path.URL("system"))
		})
	*/

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
			r.auth.CreateHandler(false, "system"),
			// エラーページのビューパスを変更
			func(c *fiber.Ctx) error {
				response := r.value.GetResponseHelper(c)
				response.SetErrorViewPath("error/auth")
				return c.Next()
			},
		)
		{
			auth.Get("/login", r.authController.LoginForm)
			auth.Post("/login", r.authController.Login)
		}

		// after login route
		system := http.Group(
			"/",
			r.auth.CreateHandler(true, "system/auth/login"),
		)
		{
			// dashboard route
			system.Get("/", r.dashboardController.Dashboard)

			// admin route
			admin := system.Group("/admin", r.role.CreateHandler(
				[]string{"system"},
			))
			{
				admin.Get("/", r.adminController.Index)
				admin.Get("/create", r.adminController.Create)
				admin.Post("/store", r.adminController.Store)
				admin.Get("/:id/edit", r.adminController.Edit)
				admin.Put("/:id/update", r.adminController.Update)
				admin.Delete("/:id/delete", r.adminController.Delete)
			}

			// account route
			account := system.Group("/account")
			{
				account.Get("/edit", r.accountController.Edit)
				account.Put("/update", r.accountController.Update)
			}

			// user route
			user := system.Group("/user")
			{
				user.Get("/", r.userController.Index)
				user.Get("/create", r.userController.Create)
				user.Post("/store", r.userController.Store)
				user.Get("/:id/edit", r.userController.Edit)
				user.Put("/:id/update", r.userController.Update)
				user.Delete("/:id/delete", r.userController.Delete)
			}

			// bot route
			bot := system.Group("/bot")
			{
				bot.Get("/", r.botController.Index)
				bot.Get("/:id/edit", r.botController.Edit)
				bot.Put("/:id/update", r.botController.Update)
				bot.Delete("/:id/delete", r.botController.Delete)
			}
			// message route
			message := system.Group("/message")
			{
				message.Get("/immediate", r.immediateController.Index)

				message.Get("/schedule", r.scheduleController.Index)
				message.Get("/schedule/:id/edit", r.scheduleController.Edit)
				message.Put("/schedule/:id/update", r.scheduleController.Update)

				message.Get("/regular", r.regularController.Index)
				message.Get("/regular/:id/edit", r.regularController.Edit)
				message.Put("/regular/:id/update", r.regularController.Update)
				message.Get("/regular/:id/timing/edit", r.timingController.Edit)
				message.Post("/regular/:id/timing/add", r.timingController.Add)
				message.Delete("/regular/:id/timing/remove/:day_of_week/:hour_time", r.timingController.Remove)

				message.Get("/history", r.messageController.History)
				message.Delete("/:id/delete", r.messageController.Delete)
			}

			// auth logout route
			system.Post("/logout", r.authController.Logout)

			// not found
			system.All("*", func(c *fiber.Ctx) error {
				response := r.value.GetResponseHelper(c)
				return response.Error(fiber.ErrNotFound)
			})
		}
	}
}

// ViewRenderCreateHandler Viewへのデータを設定する
func (r AdminRoute) ViewRenderCreateHandler(c *fiber.Ctx, vr *helper.ViewRender) (err error) {
	// load Context list
	adminlte, err := r.app.Config.Load("admin-lte")
	if err == nil {
		for k, v := range map[string]string{
			"adminlte_menus":   "menus",
			"adminlte_plugins": "plugins",
			"adminlte_links":   "links",
		} {
			if value, ok := adminlte[v]; ok {
				vr.SetData(helper.DataMap{
					k: value,
				})
			}
		}
	} else {
		return err
	}

	// load meta
	meta, err := r.app.Config.Load("admin-meta")
	if err == nil {
		vr.SetData(helper.DataMap{
			"site_meta": meta,
		})
	} else {
		return err
	}

	// csrf-token
	csrfToken := middleware.GetCSRFToken(c)

	// session errors
	errors, _ := r.sessionStore.GetSessionErrors(c)

	// session inputs
	inputs, _ := r.sessionStore.GetSessionInputs(c)

	// session messages
	messages, _ := r.sessionStore.GetSessionMessages(c)

	// admin auth
	var auth helper.DataMap
	adminContext, create := support.CreateSessionAdminAuthContext(
		c,
		r.sessionStore,
	)

	if create == nil {
		adminAuth, _ := adminContext.AdminAuth()
		auth = helper.DataMap{
			"id":    adminAuth.ID().Value(),
			"name":  adminAuth.Name().Value(),
			"email": adminAuth.Email().Value(),
			"role":  adminAuth.Role().Value(),
		}
	}

	data := helper.DataMap{
		"csrf_token":   csrfToken,
		"errors":       errors,
		"inputs":       inputs,
		"messages":     messages,
		"auth":         auth,
		"account_link": "system/account/edit",
	}

	vr.SetData(data)
	vr.SetJS(data)

	return err
}
