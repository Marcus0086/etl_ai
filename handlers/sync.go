package handlers

import (
	requestModels "etl/pkg/models"
	"etl/pkg/orchestrator"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

func SyncHandler(event *core.ServeEvent) error {
	app := event.App
	event.Router.POST("/api/sync", func(ctx echo.Context) error {
		_, ok := ctx.Get(apis.ContextAuthRecordKey).(*models.Record)
		if !ok {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		var requestBody requestModels.SyncRequest
		if err := ctx.Bind(&requestBody); err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		connectionId := requestBody.ConnectionId
		connectionRecord, err := app.Dao().FindRecordById("connections", connectionId)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		go orchestrator.StartEtlWorkflow(connectionRecord)
		return ctx.JSON(http.StatusCreated, map[string]string{"status": "started etl workflow for connection:" + connectionId})
	}, apis.RequireRecordAuth("users"), apis.ActivityLogger(app))
	return nil
}
