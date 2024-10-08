package database

import (
	"fmt"

	"github.com/rangidev/rangi/blueprint"
	"github.com/rangidev/rangi/collection"
	"github.com/rangidev/rangi/item"
)

func (db *DB) GetItems(collection *collection.Collection, limit int, offset int64) ([]item.Item, error) {
	// TODO: Support for reference fields via JOIN
	var results []item.Item
	rows, err := db.db.Queryx(fmt.Sprintf("SELECT * FROM %s ORDER BY %s DESC LIMIT $2 OFFSET $3;", collection.Blueprint.CollectionName, blueprint.KeyUpdatedAt), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %v", err)
	}
	for rows.Next() {
		result := item.Item{}
		err = rows.MapScan(result)
		if err != nil {
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		results = append(results, result)
	}
	return results, nil
}

func (db *DB) GetItem(collection *collection.Collection, id string) (item.Item, error) {
	// TODO: Support for reference fields via JOIN
	result := item.Item{}
	row := db.db.QueryRowx(fmt.Sprintf("SELECT * FROM %s WHERE id=$1;", collection.Blueprint.CollectionName), id)
	err := row.MapScan(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
