package blueprint

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/rangidev/rangi/sql"
)

var (
	//go:embed blueprint
	blueprintFS embed.FS

	defaultBlueprintFields = []BlueprintField{
		{
			Name:        KeyID,
			DisplayName: "ID",
			Type:        TypeID,
			Required:    true,
			Hidden:      true,
		},
		{
			Name:        KeyUUID,
			DisplayName: "UUID",
			Type:        TypeUUID,
			Required:    true,
			Hidden:      true,
		},
		{
			Name:        KeyCollection,
			DisplayName: "Collection",
			Type:        TypeString,
			Required:    true,
			Hidden:      true,
		},
		{
			Name:        KeyUpdatedAt,
			DisplayName: "Updated at",
			Type:        TypeInt,
			Required:    true,
			Hidden:      true,
		},
		{
			Name:        KeyTitle,
			DisplayName: "Title",
			Type:        TypeString,
			Required:    true,
		},
	}
)

// Blueprint
// "CollectionName" and map keys for "Fields" are vetted so that they can be safely used inside SQL statements
type Blueprint struct {
	CollectionName        string           `json:"collection_name"`
	CollectionDisplayName string           `json:"collection_display_name"`
	Fields                []BlueprintField `json:"fields"`
}

type BlueprintField struct {
	Name        string             `json:"name"` // Name will be used for field name in SQL tables
	DisplayName string             `json:"display_name"`
	Type        Type               `json:"type"`
	Required    bool               `json:"required"`
	Hidden      bool               `json:"hidden"`
	Reference   BlueprintReference `json:"reference"`
}

type BlueprintReference struct {
	Collection    string `json:"collection"`
	MaxReferences int    `json:"max_references"` // -1 means infinite references allowed
}

func LoadBlueprint(collectionName string, blueprintsPath string) (*Blueprint, error) {
	filename := collectionName + ".json"
	// Try directory on disk first
	// If found, it will overwrite the embedded blueprint
	data, err := os.ReadFile(filepath.Join(blueprintsPath, filename))
	if err != nil {
		// Try embed FS
		data, err = blueprintFS.ReadFile(path.Join("blueprint", filename))
		if err != nil {
			return nil, fmt.Errorf("could not read blueprint data: %v", err)
		}
	}
	var blueprint Blueprint
	err = json.Unmarshal(data, &blueprint)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json data: %v", err)
	}
	// Add default fields
	blueprint.Fields = append(defaultBlueprintFields, blueprint.Fields...)
	// Vet field names to prevent SQL injections
	if !sql.AllowedFieldAndTableNameRegex.MatchString(blueprint.CollectionName) {
		return nil, fmt.Errorf("invalid collection name %s", blueprint.CollectionName)
	}
	for _, field := range blueprint.Fields {
		if !sql.AllowedFieldAndTableNameRegex.MatchString(field.Name) {
			return nil, fmt.Errorf("invalid field name %s", field.Name)
		}
	}
	return &blueprint, nil
}
