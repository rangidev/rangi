package blueprint

import (
	"github.com/google/uuid"
)

type Item map[string]interface{}

func NewItem(col *Collection) (Item, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return Item{
		KeyUUID:       uuid.String(),
		KeyCollection: col.Blueprint.CollectionName,
	}, nil
}
