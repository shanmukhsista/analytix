package lib

import (
	"github.com/ansel1/merry/v2"
	"github.com/gin-gonic/gin"
)

//
// APP error definition
//
type ApiError struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Details map[string]interface{}
}



const DefaultGenericError = `Internal Server Error!`

// Middleware is a handler container for the middleware.
type ErrorMiddleware struct {
	// Debug controls if a call stack should be printed with every error.
	// Defaults to false.
	Debug bool

	// GenericError is a string that is shown on error code 500 or
	// non-merry errors.
	GenericError string

	// LogFunc is the function that gets called each time an error is occurred.
	LogFunc func(err string,code int, vals map[string]interface{})
}

// New returns new middleware container with default options.
// If parameter is true then debug mode is assumed.
func New(debug bool) *ErrorMiddleware {
	return &ErrorMiddleware{Debug:debug,GenericError:DefaultGenericError}
}

// Handler returns the middleware func.
func (m *ErrorMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		c.Next()

		// Skip if no errors
		if c.Errors.Last() == nil {
			return
		}

		// Get last error, clear all errors
		err := c.Errors.Last().Err
		c.Errors = c.Errors[:0]

		// Form the output dict
		// Only takes stuff that has string as a key.
		out := ApiError{Message: err.Error(), Details: map[string]interface{}{}}
		for key, val := range merry.Values(err) {
			if key == "message" || key == "http status code" {
				continue
			}
			if key, ok := key.(string); ok {
				out.Details[key] = val
			}
		}

		// Add the error's stack if Debug is enabled
		if m.Debug {
			out.Details[`stack`] = merry.Stacktrace(err)
		}

		errCode := merry.HTTPCode(err)
		// Log the error
		if m.LogFunc != nil {
			m.LogFunc(err.Error(),errCode,out.Details)
		}

		// Hide error 500
		if merry.HTTPCode(err) == 500 {
			out.Message = m.GenericError
			out.Details = nil
			return
		}

		c.JSON(merry.HTTPCode(err), out)
	}
}