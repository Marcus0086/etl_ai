package routes

import (
	"formdata/handlers"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func SetupRoutes(app *pocketbase.PocketBase) error {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		if err := handlers.SourceHandler(e); err != nil {
			return err
		}
		if err := handlers.LoaderHandler(e); err != nil {
			return err
		}
		if err := handlers.ConnectionHandler(e); err != nil {
			return err
		}
		if err := handlers.SyncHandler(e); err != nil {
			return err
		}
		return nil
	})
	return nil
}
