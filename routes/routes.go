package routes

import (
	"formdata/handlers"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func SetupRoutes(app *pocketbase.PocketBase) error {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		var err error
		err = handlers.SourceHandler(e)
		if err != nil {
			return err
		}
		err = handlers.LoaderHandler(e)
		if err != nil {
			return err
		}
		err = handlers.ConnectionHandler(e)
		if err != nil {
			return err
		}
		err = handlers.SyncHandler(e)
		if err != nil {
			return err
		}
		return nil
	})
	return nil
}
