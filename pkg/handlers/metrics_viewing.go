package handlers

import (
	"github.com/KanathipP/fl-kube-reader-backend/pkg/response"
	"github.com/KanathipP/fl-kube-reader-backend/pkg/service"
	"github.com/gofiber/fiber/v2"
)

type MetricsViewingHandler struct {
	metricsViewingService *service.MetricsViewingService
}

func NewMetricsViewingHandler(metricsViewingService *service.MetricsViewingService) *MetricsViewingHandler {
	return &MetricsViewingHandler{
		metricsViewingService: metricsViewingService,
	}
}

// GetPodUsage godoc
// @Summary     Get pod resource usage
// @Description Returns CPU (milli) + memory (bytes) usage for a specific pod. If container is provided, returns usage only for that container.
// @Tags        metrics
// @Accept      json
// @Produce     json
// @Param       pod        query     string  true   "Pod name"
// @Param       container  query     string  false  "Container name (optional)"
// @Success     200        {object}  map[string]interface{}  "Pod usage"
// @Failure     400        {object}  response.ErrorResponse   "Bad request"
// @Failure     500        {object}  response.ErrorResponse   "Server error"
// @Router      /metrics-viewing/pods/usage [get]
func (h *MetricsViewingHandler) GetPodUsage(ctx *fiber.Ctx) error {
	pod := ctx.Query("pod", "")
	container := ctx.Query("container", "")

	if pod == "" {
		return response.BadRequest(ctx, "param 'pod' is required")
	}

	usage, err := h.metricsViewingService.GetPodUsage(ctx.Context(), pod, container)
	if err != nil {
		return response.InternalServerError(ctx, "Failed to get pod usage: "+err.Error())
	}

	return response.OK(ctx, usage)
}

// GetNodesUsage godoc
// @Summary     Get cluster node resource usage
// @Description Returns CPU & memory usage and allocatable resources for all nodes in the cluster.
// @Tags        metrics
// @Accept      json
// @Produce     json
// @Success     200  {array}   map[string]interface{}  "Nodes usage list"
// @Failure     500  {object}  response.ErrorResponse   "Server error"
// @Router      /metrics-viewing/nodes/usage [get]
func (h *MetricsViewingHandler) GetNodesUsage(ctx *fiber.Ctx) error {
	usages, err := h.metricsViewingService.GetNodesUsage(ctx.Context())
	if err != nil {
		return response.InternalServerError(ctx, "Failed to get nodes usage: "+err.Error())
	}

	return response.OK(ctx, usages)
}
