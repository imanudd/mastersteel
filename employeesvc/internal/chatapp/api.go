package chatapp

import (
	"finalproject/config"
	"net/http"

	"github.com/labstack/echo/v4"
)

type handler struct {
	cfg     *config.Config
	service IService
}

func RegisterAPI(r echo.Group, cfg *config.Config, service IService) {
	handler := handler{cfg: cfg, service: IService(service)}

	// r.GET("/get-message", handler.sendMessage)
	r.GET("/testAPI", handler.testAPI)
}

// func (h handler) sendMessage(c echo.Context) error {
// 	// ctx := c.Request().Context()
// 	// resp, _ := h.service.SendMessage(ctx)

// 	// fmt.pr

// 	// if err != nil {
// 	// 	return err
// 	// }
// 	return nil

// }

func (h handler) testAPI(c echo.Context) error {
	type response struct {
		Message string      `json:"message"`
		Code    string      `json:"code"`
		Data    interface{} `json:"data,omitempty"`
	}

	ctx := c.Request().Context()

	resp, err := h.service.SendMessage(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	res := response{
		Message: "success",
		Code:    "200",
		Data:    resp,
	}

	return c.JSON(http.StatusOK, res)
}
