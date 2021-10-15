package contract

// Middleware interface
type Middleware interface {
	Setup()
}

// Route interface
type Route interface {
	Setup()
}

// Command interface
type Command interface {
	Setup()
}

// ConfigPath
type ConfigPath string

// AppBooter app boot interface
type AppBooter interface {
	AppBoot()
}

// CLIBooter cli boot interface
type CLIBooter interface {
	CLIBoot()
}
