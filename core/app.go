package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Application struct {
	App      *fiber.App
	Config   Config
	Path     Path
	Template TemplateEngine
}

// NewApplication create a new application
func NewApplication(
	config Config,
	path Path,
	template TemplateEngine,
) Application {
	kbyte := 1024
	app := fiber.New(fiber.Config{
		Prefork:         config.Server.Prefork,
		StrictRouting:   config.Server.Strict,
		CaseSensitive:   config.Server.Case,
		ETag:            config.Server.Etag,
		BodyLimit:       config.Server.BodyLimit * kbyte * kbyte,
		Concurrency:     config.Server.Concurrency * kbyte,
		ReadTimeout:     config.Server.Timeout.Read * time.Second,
		WriteTimeout:    config.Server.Timeout.Write * time.Second,
		IdleTimeout:     config.Server.Timeout.Idel * time.Second,
		ReadBufferSize:  config.Server.Buffer.Read * kbyte,
		WriteBufferSize: config.Server.Buffer.Write * kbyte,
		Views:           template.Engine,
	})
	return Application{
		App:    app,
		Config: config,
		Path:   path,
	}
}

// Run is start server
func (app *Application) Run() {
	app.Setup()

	app.App.Listen(fmt.Sprintf("%s:%d", app.Config.App.Host, app.Config.App.Port))
}

// Setup is all setup
func (app *Application) Setup() {
	app.setupTimeLocale()
	app.setupStatic()
	app.setupMiddleware()
}

// setupTimeLocal is setup time locale
func (app *Application) setupTimeLocale() {
	time.Local = time.FixedZone(app.Config.Time.Zone, app.Config.Time.Offset)
	time.LoadLocation(app.Config.Time.Location)
}

// setupStatic is setup static path
func (app *Application) setupStatic() {
	app.App.Static(
		app.Config.Static.Prefix,
		app.Path.Current(app.Config.Static.Root),
		fiber.Static{
			Index: app.Config.Static.Index,
		},
	)
}

// setupMiddleware is setup middleware
func (app *Application) setupMiddleware() {
	app.App.Use(logger.New())
	app.App.Use(recover.New())
	app.App.Use(cache.New(cache.Config{
		Expiration:   app.Config.Cache.Expiration * time.Minute,
		CacheControl: app.Config.Cache.Control,
	}))
}

// Environment check env
func (app *Application) Environment(env string) bool {
	return strings.ToLower(app.Config.App.Env) == strings.ToLower(env)
}
