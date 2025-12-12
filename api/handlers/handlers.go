package handlers

import (
	"errors"

	"finlog-api/api/constants"
	"finlog-api/api/models/responses"

	"github.com/gofiber/fiber/v2"
)

func HttpError(c *fiber.Ctx, err error) error {
	var errResponse *responses.ErrorResponse
	if errors.As(err, &errResponse) {
		c.Status(errResponse.Status)
	}

	if errResponse == nil {
		errResponse = &responses.ErrorResponse{}
	}

	if app.Config[constants.ServerEnv] == constants.EnvDevelopment {
		errResponse.Debug = errResponse.Error()
	}

	c.Append("Access-Control-Allow-Origin", "*")

	_ = c.JSON(errResponse)

	return nil
}

func HttpSuccess(c *fiber.Ctx, message string, data interface{}) (err error) {
	response := responses.Response{}
	response.Status = fiber.StatusOK
	response.Message = message
	response.Data = data

	c.Append("Access-Control-Allow-Origin", "*")

	_ = c.JSON(response)

	return nil
}
