package handlers

import (
	"log/slog"
	"strconv"

	"github.com/KanathipP/fl-kube-reader-backend/pkg/response"
	"github.com/KanathipP/fl-kube-reader-backend/pkg/service"
	"github.com/gofiber/fiber/v2"
)

type LogReadingHandler struct {
	logReadingService *service.LogReadingService
}

func NewLogReadingHandler(logReadingService *service.LogReadingService) *LogReadingHandler {
	return &LogReadingHandler{logReadingService: logReadingService}
}

// GetPodLog godoc
// @Summary     Get pod logs
// @Description Returns container logs for a pod. Supports tailing N lines and fetching previous container logs.
// @Tags        logs
// @Accept      json
// @Produce     json
// @Param       pod        query     string  true   "Pod name"
// @Param       container  query     string  false  "Container name (optional)"
// @Param       tail       query     int     false  "Tail last N lines (default: 0 = all)"
// @Param       previous   query     bool    false  "Get logs from the previously terminated container"
// @Success     200        {object}  map[string]interface{}  "Log payload"
// @Failure     400        {object}  response.ErrorResponse   "Bad request"
// @Failure     500        {object}  response.ErrorResponse   "Server error"
// @Router      /logs/pod [get]
func (h *LogReadingHandler) GetPodLog(ctx *fiber.Ctx) error {
	pod := ctx.Query("pod", "")
	if pod == "" {
		return response.BadRequest(ctx, "param 'pod' is required")
	}

	container := ctx.Query("container", "")

	tailStr := ctx.Query("tail", "0")
	tail := 0
	if tailStr != "" {
		if n, err := strconv.Atoi(tailStr); err != nil || n < 0 {
			return response.BadRequest(ctx, "param 'tail' must be a non-negative integer")
		} else {
			tail = n
		}
	}

	prevStr := ctx.Query("previous", "false")
	previous := false
	if prevStr != "" {
		if b, err := strconv.ParseBool(prevStr); err != nil {
			return response.BadRequest(ctx, "param 'previous' must be a boolean")
		} else {
			previous = b
		}
	}

	slog.Info("GetPodLog",
		"pod", pod,
		"container", container,
		"tail", tail,
		"previous", previous,
	)

	logText, err := h.logReadingService.ReadPodLog(ctx.Context(), pod, container, tail, previous)
	if err != nil {
		return response.InternalServerError(ctx, "failed to read pod log: "+err.Error())
	}

	// Return a simple JSON envelope so clients get metadata + text.
	return response.OK(ctx, fiber.Map{
		"pod":       pod,
		"container": container,
		"tail":      tail,
		"previous":  previous,
		"log":       logText,
	})
}
