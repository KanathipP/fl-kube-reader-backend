package handlers

import (
	"github.com/KanathipP/fl-kube-reader-backend/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// Healthz godoc
// @Summary     Health check
// @Description Basic health endpoint to verify service is alive
// @Tags        system
// @Success     200 {string} string "OK"
// @Router      /healthz [get]
func Healthz(ctx *fiber.Ctx) error {
	return response.OK(ctx, "OK")
}

// NotFound godoc
// @Summary     Not Found handler
// @Description Handles unknown routes and returns 404 JSON payload
// @Tags        system
// @Failure     404 {object} response.ErrorResponse "Not Found"
// @Router      /{any} [get]
func NotFound(ctx *fiber.Ctx) error {
	return response.NotFound(ctx, "Not Found")
}
