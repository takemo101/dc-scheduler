package route

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/takemo101/dc-scheduler/app/controller/user"
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/app/middleware"
	"github.com/takemo101/dc-scheduler/app/support"
	"github.com/takemo101/dc-scheduler/core"
)

// AdminRoute 管理者ルート
type UserRoute struct {
	logger              core.Logger
	app                 core.Application
	path                core.Path
	sessionStore        core.SessionStore
	csrf                middleware.Csrf
	value               support.ContextValue
	auth                middleware.SessionUserAuth
	render              middleware.ViewRender
	authController      controller.SessionAuthController
	dashboardController controller.DashboardController
	userController      controller.UserController
	accountController   controller.AccountController
	botController       controller.BotController
	messageController   controller.PostMessageController
	immediateController controller.ImmediatePostController
	scheduleController  controller.SchedulePostController
	regularController   controller.RegularPostController
	timingController    controller.RegularTimingController
}

// NewUserRoute コンストラクタ
func NewUserRoute(
	logger core.Logger,
	app core.Application,
	path core.Path,
	sessionStore core.SessionStore,
	csrf middleware.Csrf,
	value support.ContextValue,
	auth middleware.SessionUserAuth,
	render middleware.ViewRender,
	authController controller.SessionAuthController,
	dashboardController controller.DashboardController,
	userController controller.UserController,
	accountController controller.AccountController,
	botController controller.BotController,
	messageController controller.PostMessageController,
	immediateController controller.ImmediatePostController,
	scheduleController controller.SchedulePostController,
	regularController controller.RegularPostController,
	timingController controller.RegularTimingController,
) UserRoute {
	return UserRoute{
		logger:              logger,
		app:                 app,
		path:                path,
		sessionStore:        sessionStore,
		csrf:                csrf,
		value:               value,
		auth:                auth,
		render:              render,
		authController:      authController,
		dashboardController: dashboardController,
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
func (r UserRoute) Setup() {
	r.logger.Info("setup user-route")

	app := r.app.App

	// root redirect
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect(r.path.URL("user/auth/login"))
	})

	// user http route
	http := app.Group(
		"/user",
		r.csrf.CreateHandler("form:csrf_token"),
		r.render.CreateHandler(r.ViewRenderCreateHandler),
	)
	{
		// auth login route
		auth := http.Group(
			"/auth",
			r.auth.CreateHandler(false, "user"),
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
			auth.Get("/regist", r.userController.RegistFrom)
			auth.Post("/regist", r.userController.Regist)
			auth.Get("/activation/:signature", r.userController.Activation)
		}

		// after login route
		user := http.Group(
			"/",
			r.auth.CreateHandler(true, "user/auth/login"),
		)
		{
			// dashboard route
			user.Get("/", r.dashboardController.Dashboard)

			// account route
			account := user.Group("/account")
			{
				account.Get("/edit", r.accountController.Edit)
				account.Put("/update", r.accountController.Update)
			}

			// bot route
			bot := user.Group("/bot")
			{
				bot.Get("/", r.botController.Index)
				bot.Get("/create", r.botController.Create)
				bot.Post("/store", r.botController.Store)
				bot.Get("/:id/edit", r.botController.Edit)
				bot.Put("/:id/update", r.botController.Update)
				bot.Delete("/:id/delete", r.botController.Delete)

				bot.Get("/:id/immediate/create", r.immediateController.Create)
				bot.Post("/:id/immediate/store", r.immediateController.Store)

				bot.Get("/:id/schedule/create", r.scheduleController.Create)
				bot.Post("/:id/schedule/store", r.scheduleController.Store)

				bot.Get("/:id/regular/create", r.regularController.Create)
				bot.Post("/:id/regular/store", r.regularController.Store)
			}

			// message route
			message := user.Group("/message")
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
			user.Post("/logout", r.authController.Logout)
		}
	}
}

// ViewRenderCreateHandler Viewへのデータを設定する
func (r UserRoute) ViewRenderCreateHandler(c *fiber.Ctx, vr *helper.ViewRender) (err error) {
	// load Context list
	adminlte, err := r.app.Config.Load("user-lte")
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
	meta, err := r.app.Config.Load("user-meta")
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

	// user auth
	var auth helper.DataMap
	userContext, create := support.CreateSessionUserAuthContext(
		c,
		r.sessionStore,
	)

	if create == nil {
		userAuth, _ := userContext.UserAuth()
		auth = helper.DataMap{
			"id":    userAuth.ID().Value(),
			"name":  userAuth.Name().Value(),
			"email": userAuth.Email().Value(),
		}
	}

	data := helper.DataMap{
		"csrf_token": csrfToken,
		"errors":     errors,
		"inputs":     inputs,
		"messages":   messages,
		"auth":       auth,
	}

	vr.SetData(data)
	vr.SetJS(data)

	return err
}
