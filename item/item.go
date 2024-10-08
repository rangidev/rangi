package item

import (
	"github.com/google/uuid"

	"github.com/rangidev/rangi/blueprint"
	"github.com/rangidev/rangi/collection"
)

type Item map[string]interface{}

func New(col *collection.Collection) (Item, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return Item{
		blueprint.KeyUUID:       uuid.String(),
		blueprint.KeyCollection: col.Blueprint.CollectionName,
	}, nil
}
