package helper

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/takemo101/dc-scheduler/core"
)

// JsonStatus json status message type
type JsonStatus string

// JsonStatus json status messages
const (
	Success JsonStatus = "success"
	Fail    JsonStatus = "fail"
	Error   JsonStatus = "error"
)

// JsonError json error type
type JsonError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// JsonResult json result type
type JsonResult struct {
	Status  JsonStatus `json:"status"`
	Message string     `json:"message"`
	Data    fiber.Map  `json:"data"`
}

// JsonErrorResult json error result type
type JsonErrorResult struct {
	Status JsonStatus `json:"status"`
	Error  JsonError  `json:"error"`
}

// JsonErrorMessagesResult json error messages result type
type JsonErrorMessagesResult struct {
	Status   JsonStatus        `json:"status"`
	Error    JsonError         `json:"error"`
	Messages map[string]string `json:"messages"`
}

// ResponseHelper response helper
type ResponseHelper struct {
	logger    core.Logger
	config    core.Config
	path      core.Path
	render    *ViewRender
	errorPath string
}

// NewResponseHelper response utility
func NewResponseHelper(
	logger core.Logger,
	config core.Config,
	path core.Path,
	render *ViewRender,
) *ResponseHelper {
	return &ResponseHelper{
		logger:    logger,
		config:    config,
		path:      path,
		render:    render,
		errorPath: "error/error",
	}
}

// Back response redirect back
func (helper *ResponseHelper) Back(c *fiber.Ctx) error {
	back := string(c.Request().Header.Referer())
	return c.Redirect(back)
}

// Redirect response redirect
func (helper *ResponseHelper) Redirect(c *fiber.Ctx, path string) error {
	return c.Redirect(helper.path.URL(path))
}

// Json response json
func (helper *ResponseHelper) Json(c *fiber.Ctx, data interface{}) error {
	return c.JSON(data)
}

// JsonSuccess response json success with data
func (helper *ResponseHelper) JsonSuccess(c *fiber.Ctx, message string, data ...fiber.Map) error {
	jsonData := fiber.Map{}
	if len(data) > 0 {
		jsonData = data[0]
	}

	return helper.Json(c, JsonResult{
		Status:  Success,
		Message: message,
		Data:    jsonData,
	})
}

// JsonError response json error
func (helper *ResponseHelper) JsonError(c *fiber.Ctx, err error, code ...int) error {
	statusCode := fiber.StatusInternalServerError
	if len(code) > 0 {
		statusCode = code[0]
	}

	helper.logger.Error(err)

	c.Status(statusCode)
	return helper.Json(c, JsonErrorResult{
		Status: Error,
		Error: JsonError{
			Code:    statusCode,
			Message: helper.CreateErrorMessage(err, statusCode),
		},
	})
}

// JsonErrorMessages response json error with error_messages
func (helper *ResponseHelper) JsonErrorMessages(c *fiber.Ctx, err error, messages map[string]string) error {
	code := fiber.StatusUnprocessableEntity
	c.Status(code)
	return helper.Json(c, JsonErrorMessagesResult{
		Status: Fail,
		Error: JsonError{
			Code:    code,
			Message: err.Error(),
		},
		Messages: messages,
	})
}

// JS set template javascript data
func (helper *ResponseHelper) JS(data DataMap) {
	helper.render.SetJS(data)
}

// Data set template view data
func (helper *ResponseHelper) Data(data DataMap) {
	helper.render.SetData(data)
}

// View render template view
func (helper *ResponseHelper) View(name string, data DataMap) error {
	return helper.render.Render(name, data)
}

// Error render template error with Code
func (helper *ResponseHelper) Error(err error, code ...int) error {
	statusCode := fiber.StatusInternalServerError
	if len(code) > 0 {
		statusCode = code[0]
	}

	helper.logger.Error(err)
	return helper.render.Error(helper.errorPath, helper.CreateErrorMessage(err, statusCode), statusCode)
}

// SetErrorViewPath set error view path
func (helper *ResponseHelper) SetErrorViewPath(path string) {
	helper.errorPath = path
}

// CreateErrorMessage create error response messsage
func (helper *ResponseHelper) CreateErrorMessage(err error, code int) string {
	// create error message
	if helper.config.App.Debug {
		return err.Error()
	}
	return utils.StatusMessage(code)
}
