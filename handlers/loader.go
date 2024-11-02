package handlers

import (
	requestModels "formdata/pkg/models"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

func LoaderHandler(event *core.ServeEvent) error {
	app := event.App
	event.Router.POST("/api/loader", func(c echo.Context) error {
		authRecord, ok := c.Get(apis.ContextAuthRecordKey).(*models.Record)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		var requestBody requestModels.RequestBody
		if err := c.Bind(&requestBody); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		loadersCollection, err := app.Dao().FindCollectionByNameOrId("loaders")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		record := models.NewRecord(loadersCollection)
		record.Set("name", requestBody.Name)
		record.Set("type", requestBody.Type)
		record.Set("config", requestBody.Config)
		record.Set("user_id", authRecord.Id)
		if err := app.Dao().SaveRecord(record); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusCreated, map[string]string{"id": record.Id})
	}, apis.RequireRecordAuth("users"), apis.ActivityLogger(app))
	return nil
}
