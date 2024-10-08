package database

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/rangidev/rangi/collection"
)

func (db *DB) CreateTables(collections []collection.Collection, collectionLoader *collection.CollectionLoader) error {
	for index := range collections {
		err := db.CreateTable(&collections[index], collectionLoader)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) CreateTable(collection *collection.Collection, collectionLoader *collection.CollectionLoader) error {
	var subStatements []string
	// TODO: Make sure field names aren't used twice between default and blueprint fields
	for _, fieldDef := range collection.Blueprint.Fields {
		sqlType, ok := db.castToSQLType(fieldDef.Type)
		if !ok {
			return fmt.Errorf("could not cast type %s", string(fieldDef.Type))
		}
		if sqlType == SQLTypeReference {
			// Get referenced collection
			refCollection, err := collectionLoader.Get(fieldDef.Reference.Collection)
			if err != nil {
				return fmt.Errorf("could not get referenced collection %s: %v", fieldDef.Reference.Collection, err)
			}
			err = db.CreateReferenceTable(collection, refCollection)
			if err != nil {
				return fmt.Errorf("could not create reference table: %v", err)
			}
			// Create reference field in table
		} else {
			statement := fmt.Sprintf("%s %s", fieldDef.Name, sqlType)
			if fieldDef.Required {
				statement = statement + " NOT NULL"
			}
			subStatements = append(subStatements, statement)
		}
	}
	finalStatement := fmt.Sprintf(statementCreateTable, collection.Blueprint.CollectionName, strings.Join(subStatements, ","))
	_, err := db.db.Exec(finalStatement)
	return err
}

func (db *DB) CreateReferenceTable(collection1 *collection.Collection, collection2 *collection.Collection) error {
	if collection1.Blueprint.CollectionName == "" || collection2.Blueprint.CollectionName == "" {
		return errors.New("empty collection string")
	}
	// Sort collections alphabetically to prevent creating reference table twice
	sorted := []string{collection1.Blueprint.CollectionName, collection2.Blueprint.CollectionName}
	slices.Sort(sorted)
	tableName := sorted[0] + "_" + sorted[1]
	var statement string
	switch db.dbType {
	case DatabaseTypeSqlite3:
		statement = fmt.Sprintf(statementCreateReferenceTableSqlite3, tableName, sorted[0], sorted[1], sorted[0], sorted[0], sorted[1], sorted[1])
	case DatabaseTypePostgres:
		statement = fmt.Sprintf(statementCreateReferenceTablePostgres, tableName, sorted[0], sorted[1], sorted[0], sorted[0], sorted[1], sorted[1])
	default:
		return ErrorUnknownDatabaseType
	}
	_, err := db.db.Exec(statement)
	return err
}
