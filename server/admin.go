package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/rangidev/rangi/admin"
	"github.com/rangidev/rangi/blueprint"
)

type getItemsQueryParams struct {
	Limit  int   `schema:"limit,required" validate:"gte=1,lte=200"`
	Offset int64 `schema:"offset,required"`
}

func createAdminRouter(s *Server) http.Handler {
	router := chi.NewRouter()
	router.Use(admin.EnsurePermission)
	router.Get("/", s.GetAdminBase)
	router.Get("/login", s.GetAdminLogin)
	router.Post("/login", s.PostAdminLogin)
	router.Get("/dashboard", s.GetAdminDashboard)
	router.Get("/collections/{collection}", s.GetAdminCollection)
	router.Get("/edit/{collection}/{id}", s.GetAdminEdit) // If id == "new", we will display an empty input form
	router.Get("/settings", s.GetAdminSettings)
	router.Post("/{collection}/items", s.PostAdminItem)
	router.Put("/{collection}/items", s.PutAdminItem)
	router.Get("/{collection}/items", s.GetAdminItems) // For possible query parameters see getItemsQueryParams
	return router
}

func (s *Server) GetAdminBase(w http.ResponseWriter, r *http.Request) {
	// Always redirect to dashboard
	// admin.EnsurePermission will redirect to "admin/login" if not authenticated
	http.Redirect(w, r, admin.DashboardPath, http.StatusFound)
}

func (s *Server) GetAdminStatic(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/admin/static/", s.adminStaticServer).ServeHTTP(w, r)
}

func (s *Server) GetAdminLogin(w http.ResponseWriter, r *http.Request) {
	err := s.adminTemplates.Render(w, nil, admin.TemplateLogin, s.collectionLoader, "")
	if err != nil {
		http.Error(w, fmt.Sprintf("error while rendering login template: %v", err), http.StatusInternalServerError)
		return
	}
}

func (s *Server) PostAdminLogin(w http.ResponseWriter, r *http.Request) {
	// Read login data
	email := r.PostFormValue("email")
	if email == "" {
		http.Error(w, "missing value 'email'", http.StatusBadRequest)
		return
	}
	password := r.PostFormValue("password")
	if password == "" {
		http.Error(w, "missing value 'password'", http.StatusBadRequest)
		return
	}
	// Return dummy cookie for now
	c := &http.Cookie{
		Name:     admin.SessionCookieName,
		Value:    "true",
		HttpOnly: true,
		Path:     "/admin/",
		Expires:  time.Now().Add(admin.SessionCookieDuration),
	}
	http.SetCookie(w, c)
	w.Header().Set("HX-Redirect", admin.DashboardPath)
}

func (s *Server) GetAdminDashboard(w http.ResponseWriter, r *http.Request) {
	err := s.adminTemplates.Render(w, nil, admin.TemplateDashboard, s.collectionLoader, "")
	if err != nil {
		http.Error(w, fmt.Sprintf("error while rendering dashboard template: %v", err), http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetAdminCollection(w http.ResponseWriter, r *http.Request) {
	// Get collection
	collectionData, err := s.getCollection(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Get items
	items, err := s.config.DatabaseInstance.GetItems(collectionData, s.config.AdminItemsLimit, 0)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not get items: %v", err), http.StatusInternalServerError)
		return
	}
	templateData := admin.TemplateData{
		"collection": collectionData.Blueprint.CollectionName,
		"items":      items,
		"limit":      s.config.AdminItemsLimit,
	}
	err = s.adminTemplates.Render(w, templateData, admin.TemplateCollection, s.collectionLoader, "")
	if err != nil {
		http.Error(w, fmt.Sprintf("error while rendering collection template: %v", err), http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetAdminEdit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing 'id' path parameter", http.StatusBadRequest)
		return
	}
	// Get collection
	collectionData, err := s.getCollection(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	item := make(blueprint.Item)
	if id != "new" {
		// Get item
		item, err = s.config.DatabaseInstance.GetItem(collectionData, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not get item: %v", err), http.StatusInternalServerError)
			return
		}
	}
	templateData := admin.TemplateData{
		"collection": collectionData.Blueprint.CollectionName,
		"blueprint":  collectionData.Blueprint,
		"item":       item,
	}
	err = s.adminTemplates.Render(w, templateData, admin.TemplateEdit, s.collectionLoader, "")
	if err != nil {
		http.Error(w, fmt.Sprintf("error while rendering collection template: %v", err), http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetAdminSettings(w http.ResponseWriter, r *http.Request) {
	err := s.adminTemplates.Render(w, nil, admin.TemplateSettings, s.collectionLoader, "")
	if err != nil {
		http.Error(w, fmt.Sprintf("error while rendering settings template: %v", err), http.StatusInternalServerError)
		return
	}
}

func (s *Server) PostAdminItem(w http.ResponseWriter, r *http.Request) {
	// Get collection
	collectionData, err := s.getCollection(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Create item
	item, err := blueprint.NewItem(collectionData)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not create new item: %v", err), http.StatusInternalServerError)
		return
	}
	for _, field := range collectionData.Blueprint.Fields {
		if field.Name == blueprint.KeyUUID || field.Name == blueprint.KeyCollection {
			// Field has been set by item.New
			// Do not overwrite field with user input
			continue
		}
		if field.Type == blueprint.TypeReference {
			// TODO: support this
			continue
		}
		value := r.PostFormValue(field.Name)
		// TODO: Parse value into type that is expected from the blueprint field
		item[field.Name] = value
	}
	if len(item) == 0 {
		http.Error(w, "item is empty", http.StatusBadRequest)
		return
	}
	err = s.config.DatabaseInstance.CreateItem(collectionData, item)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not set item in database: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("HX-Redirect", admin.CollectionPath(collectionData.Blueprint.CollectionName))
}

func (s *Server) PutAdminItem(w http.ResponseWriter, r *http.Request) {
	// Get collection
	collectionData, err := s.getCollection(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Create item
	item := blueprint.Item{}
	for _, field := range collectionData.Blueprint.Fields {
		if field.Name == blueprint.KeyUUID || field.Name == blueprint.KeyCollection {
			// Do not update field in database
			continue
		}
		if field.Type == blueprint.TypeReference {
			// TODO: support this
			continue
		}
		value := r.PostFormValue(field.Name)
		// TODO: Parse value into type that is expected from the blueprint field
		item[field.Name] = value
	}
	if len(item) == 0 {
		http.Error(w, "item is empty", http.StatusBadRequest)
		return
	}
	if _, ok := item[blueprint.KeyID]; !ok {
		http.Error(w, "no id in item", http.StatusBadRequest)
		return
	}
	err = s.config.DatabaseInstance.UpdateItem(collectionData, item)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not set item in database: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("HX-Redirect", admin.CollectionPath(collectionData.Blueprint.CollectionName))
}

func (s *Server) GetAdminItems(w http.ResponseWriter, r *http.Request) {
	// Get collection
	collectionData, err := s.getCollection(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Read query parameters
	var queryParams getItemsQueryParams
	err = s.schemaDecoder.Decode(&queryParams, r.URL.Query())
	if err != nil {
		http.Error(w, fmt.Sprintf("could not decode query parameters: %v", err), http.StatusBadRequest)
		return
	}
	// Check if valid parameters
	err = s.config.Validate.Struct(&queryParams)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid query parameters: %v", err), http.StatusBadRequest)
		return
	}
	// Get items
	items, err := s.config.DatabaseInstance.GetItems(collectionData, queryParams.Limit, queryParams.Offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not get items: %v", err), http.StatusInternalServerError)
		return
	}
	templateData := admin.TemplateData{
		"collection": collectionData.Blueprint.CollectionName,
		"items":      items,
		"limit":      queryParams.Limit,
		"offset":     queryParams.Offset,
	}
	err = s.adminTemplates.Render(w, templateData, admin.TemplateCollection, s.collectionLoader, "list")
	if err != nil {
		http.Error(w, fmt.Sprintf("error while rendering collection template: %v", err), http.StatusInternalServerError)
		return
	}
}

func (s *Server) getCollection(r *http.Request) (*blueprint.Collection, error) {
	collection := chi.URLParam(r, "collection")
	if collection == "" {
		return nil, errors.New("missing 'collection' path parameter")
	}
	// Get collection data
	collectionData, err := s.collectionLoader.Get(collection)
	if err != nil {
		return nil, fmt.Errorf("could not get collection data: %v", err)
	}
	return collectionData, nil
}
