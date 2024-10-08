package database

import (
	"fmt"
	"time"

	"github.com/rangidev/rangi/blueprint"
	"github.com/rangidev/rangi/collection"
	"github.com/rangidev/rangi/item"
)

func (db *DB) CreateItem(collection *collection.Collection, item item.Item) error {
	// Set updated_at field
	item[blueprint.KeyUpdatedAt] = time.Now().Unix()
	statementStart := fmt.Sprintf("INSERT INTO %s (", collection.Blueprint.CollectionName)
	statementEnd := "VALUES ("
	fieldAdded := false
	for index := range collection.Blueprint.Fields {
		fieldName := collection.Blueprint.Fields[index].Name
		if fieldName == blueprint.KeyID {
			// "id" field is always set by the DB
			continue
		}
		if _, ok := item[fieldName]; !ok {
			// Blueprint field is not available in item
			continue
		}
		if fieldAdded {
			statementStart += ","
			statementEnd += ","
		}
		statementStart += collection.Blueprint.Fields[index].Name
		statementEnd += ":" + collection.Blueprint.Fields[index].Name
		fieldAdded = true
	}
	statement := statementStart + ") " + statementEnd + ");"
	_, err := db.db.NamedExec(statement, item)
	return err
}

func (db *DB) UpdateItem(collection *collection.Collection, item item.Item) error {
	// Set updated_at field
	item[blueprint.KeyUpdatedAt] = time.Now().Unix()
	statementStart := fmt.Sprintf("UPDATE %s SET (", collection.Blueprint.CollectionName)
	statementEnd := "= ("
	fieldAdded := false
	for index := range collection.Blueprint.Fields {
		fieldName := collection.Blueprint.Fields[index].Name
		if fieldName == blueprint.KeyID || fieldName == blueprint.KeyCollection || fieldName == blueprint.KeyUUID {
			// We do not want to update "id", "collection", and "uuid" fields
			continue
		}
		if _, ok := item[fieldName]; !ok {
			// Blueprint field is not available in item
			continue
		}
		if fieldAdded {
			statementStart += ","
			statementEnd += ","
		}
		statementStart += collection.Blueprint.Fields[index].Name
		statementEnd += ":" + collection.Blueprint.Fields[index].Name
		fieldAdded = true
	}
	statement := statementStart + ") " + statementEnd + ") WHERE id = :id;"
	_, err := db.db.NamedExec(statement, item)
	return err
}
