package routes

import (
	_ "github.com/KanathipP/fl-kube-reader-backend/pkg/docs" // swagger docs

	"github.com/KanathipP/fl-kube-reader-backend/pkg/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(
	app *fiber.App,
	podViewingHandler *handlers.PodViewingHandler,
	metricsViewingHandler *handlers.MetricsViewingHandler,
	logReadingHandler *handlers.LogReadingHandler,
) {
	app.Get("/healthz", handlers.Healthz)

	api := app.Group("/api")

	api.Get("/swagger/*", swagger.HandlerDefault)

	v1 := api.Group("/v1")

	podViewing := v1.Group("/pod-viewing")
	podViewing.Get("/", podViewingHandler.GetPods)

	logReading := v1.Group("/log-reading")
	logReading.Get("/pod", logReadingHandler.GetPodLog)

	metricsViewing := v1.Group("/metrics-viewing")
	metricsViewing.Get("/pods/usage", metricsViewingHandler.GetPodUsage)
	metricsViewing.Get("/nodes/usage", metricsViewingHandler.GetNodesUsage)

	app.Use(handlers.NotFound)
}
