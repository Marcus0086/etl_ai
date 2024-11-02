package handlers

import (
	requestModels "formdata/pkg/models"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

func SourceHandler(event *core.ServeEvent) error {
	app := event.App
	event.Router.POST("/api/source", func(ctx echo.Context) error {
		authRecord, ok := ctx.Get(apis.ContextAuthRecordKey).(*models.Record)
		if !ok {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		var requestBody requestModels.RequestBody
		if err := ctx.Bind(&requestBody); err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		sourcesCollection, err := app.Dao().FindCollectionByNameOrId("sources")
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		record := models.NewRecord(sourcesCollection)
		record.Set("name", requestBody.Name)
		record.Set("type", requestBody.Type)
		record.Set("config", requestBody.Config)
		record.Set("user_id", authRecord.Id)
		if err := app.Dao().SaveRecord(record); err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return ctx.JSON(http.StatusCreated, map[string]string{"id": record.Id})
	}, apis.RequireRecordAuth("users"), apis.ActivityLogger(app))
	return nil
}
