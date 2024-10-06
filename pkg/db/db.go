package db

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

func SetupSchema(app *pocketbase.PocketBase) error {
	sourcesCollection := &models.Collection{
		Name: "sources",
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name:     "name",
				Type:     schema.FieldTypeText,
				Required: true,
			},
			&schema.SchemaField{
				Name:     "type",
				Type:     schema.FieldTypeText,
				Required: true,
			},
			&schema.SchemaField{
				Name:     "config",
				Type:     schema.FieldTypeJson,
				Required: true,
			},
		),
	}
	if err := saveCollection(app, sourcesCollection); err != nil {
		log.Fatal(err)
	}

	loadersCollection := &models.Collection{
		Name: "loaders",
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name:     "name",
				Type:     schema.FieldTypeText,
				Required: true,
			},
			&schema.SchemaField{
				Name:     "type",
				Type:     schema.FieldTypeText,
				Required: true,
			},
			&schema.SchemaField{
				Name:     "config",
				Type:     schema.FieldTypeJson,
				Required: true,
			},
		),
	}
	if err := saveCollection(app, loadersCollection); err != nil {
		log.Fatal(err)
	}

	connectionsCollection := &models.Collection{
		Name: "connections",
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name: "source_id",
				Type: schema.FieldTypeRelation,
				Options: &schema.RelationOptions{
					CollectionId: "sources",
				},
				Required: true,
			},
			&schema.SchemaField{
				Name: "loader_id",
				Type: schema.FieldTypeRelation,
				Options: &schema.RelationOptions{
					CollectionId: "loaders",
				},
				Required: true,
			},
			&schema.SchemaField{
				Name:     "sync_type",
				Type:     schema.FieldTypeText,
				Required: true,
			},
			&schema.SchemaField{
				Name:     "schedule",
				Type:     schema.FieldTypeText,
				Required: false,
			},
			&schema.SchemaField{
				Name:     "config",
				Type:     schema.FieldTypeJson,
				Required: false,
			},
		),
	}
	if err := saveCollection(app, connectionsCollection); err != nil {
		return err
	}
	return nil
}

func saveCollection(app *pocketbase.PocketBase, collection *models.Collection) error {
	dao := app.Dao()
	existingCollection, err := dao.FindCollectionByNameOrId(collection.Name)
	if err != nil {
		// Collection doesn't exist, create it
		return dao.SaveCollection(collection)
	}

	// collection exists, update it
	collection.Id = existingCollection.Id
	collection.MarkAsNotNew()
	return dao.SaveCollection(collection)

}
