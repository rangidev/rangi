package blueprint

var (
	defaultCollections = []string{
		"articles",
		"authors",
	}
)

type Collection struct {
	Blueprint *Blueprint
}

type CollectionLoader struct {
	blueprintsPath string
}

func NewCollectionLoader(blueprintsPath string) *CollectionLoader {
	return &CollectionLoader{blueprintsPath: blueprintsPath}
}

// GetAll
// returns all configured collections
func (cl *CollectionLoader) GetAll() ([]Collection, error) {
	// TODO: Cache collections
	var collections []Collection
	for _, defaultCollection := range defaultCollections {
		blueprint, err := LoadBlueprint(defaultCollection, cl.blueprintsPath)
		if err != nil {
			return nil, err
		}
		collections = append(collections, Collection{Blueprint: blueprint})
	}
	// TODO: Load blueprints from disk
	return collections, nil
}

// Get
// returns the named collection
func (cl *CollectionLoader) Get(name string) (*Collection, error) {
	// TODO: Cache collections
	return getCollection(name, cl.blueprintsPath)
}

func getCollection(name string, blueprintsPath string) (*Collection, error) {
	blueprint, err := LoadBlueprint(name, blueprintsPath)
	if err != nil {
		return nil, err
	}
	return &Collection{Blueprint: blueprint}, nil
}
