package handlers

import (
	requestModels "formdata/pkg/models"
	"formdata/pkg/orchestrator"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

func ConnectionHandler(event *core.ServeEvent) error {
	app := event.App

	event.Router.POST("/api/connection", func(ctx echo.Context) error {
		authRecord, ok := ctx.Get(apis.ContextAuthRecordKey).(*models.Record)
		if !ok {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		var connectionRequest requestModels.ConnectionBody
		if err := ctx.Bind(&connectionRequest); err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		sourceRecord, err := app.Dao().FindRecordById("sources", connectionRequest.SourceID)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid source_id"})
		}

		loaderRecord, err := app.Dao().FindRecordById("loaders", connectionRequest.LoaderID)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid loader_id"})
		}
		connectionCollection, err := app.Dao().FindCollectionByNameOrId("connections")
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		record := models.NewRecord(connectionCollection)
		record.Set("source_id", sourceRecord.Id)
		record.Set("loader_id", loaderRecord.Id)
		record.Set("sync_type", connectionRequest.SyncType)
		record.Set("schedule", connectionRequest.Schedule)
		record.Set("config", connectionRequest.Config)
		record.Set("user_id", authRecord.Id)

		if err := app.Dao().SaveRecord(record); err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		if err := orchestrator.ConfigureOrchestrator(&app, record); err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return ctx.JSON(http.StatusCreated, map[string]string{"id": record.Id})
	}, apis.RequireRecordAuth("users"), apis.ActivityLogger(app))
	return nil
}
