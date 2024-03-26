package cookbook

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/ayhonz/racook/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const jsonContentType = "application/json"

type CookBookStorage interface {
	CreateUser(user database.CreateUserParams) (database.User, error)
	CreateRecipe(recipe database.CreateRecipeParams) (database.Recipe, error)
	GetRecipeByID(id int) (database.Recipe, error)
	GetRecipes() ([]database.Recipe, error)
	Authenticate(email string, password string) (int, error)
}

type CookBookServer struct {
	store          CookBookStorage
	sessionManager *scs.SessionManager
	http.Handler
}

func NewCookBookServer(store CookBookStorage, sessionManager *scs.SessionManager) *CookBookServer {
	server := new(CookBookServer)

	server.store = store
	server.sessionManager = sessionManager

	router := chi.NewRouter()

	v1Router := chi.NewRouter()

	v1Router.Use(middleware.Logger)
	v1Router.Use(middleware.Recoverer)
	v1Router.Use(server.sessionManager.LoadAndSave)

	router.Mount("/v1", v1Router)

	v1Router.Get("/healthz", server.healthHandler)

	v1Router.Get("/recipes", server.getRecipesHandler)
	v1Router.Get("/recipes/{id}", server.getRecipeByIDHandler)

	v1Router.Post("/users/register", server.createUserHandler)
	v1Router.Post("/users/login", server.userLoginHandler)

	v1Router.Group(func(r chi.Router) {
		r.Use(server.requireAuthentication)

		r.Post("/recipes", server.createRecipeHandler)
		r.Post("/users/logout", server.userLogoutHandler)
	})

	server.Handler = router

	return server
}

func (c *CookBookServer) userLoginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := new(parameters)
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("error parsing JSON %v", err))
	}

	id, err := c.store.Authenticate(params.Email, params.Password)
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("error authenticating user %v", err))
		return
	}

	err = c.sessionManager.RenewToken(r.Context())
	if err != nil {
		responseWithError(w, 500, fmt.Sprintf("error renewing session token %v", err))
		return
	}

	c.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	responseWithJSON(w, 200, "OK")
}

func (c *CookBookServer) userLogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := c.sessionManager.RenewToken(r.Context())
	if err != nil {
		responseWithError(w, 500, fmt.Sprintf("error renewing session token %v", err))
		return
	}

	c.sessionManager.Remove(r.Context(), "authenticatedUserID")
	responseWithJSON(w, 200, "OK")
}

func (c *CookBookServer) isAuthenticated(r *http.Request) bool {
	return c.sessionManager.Exists(r.Context(), "authenticatedUserID")
}

func (c *CookBookServer) healthHandler(w http.ResponseWriter, r *http.Request) {
	responseWithJSON(w, http.StatusOK, "OK")
}

func (c *CookBookServer) createRecipeHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Categories  []string `json:"categories"`
	}

	params := new(parameters)
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("error parsing JSON %v", err))
	}

	recipe, err := c.store.CreateRecipe(database.CreateRecipeParams{
		Title:       params.Title,
		Description: params.Description,
		UpdatedAt:   time.Now().UTC(),
		CreatedAt:   time.Now().UTC(),
		Categories:  []string{"test"},
		UserID:      1,
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Could't create recipe %v", err))
		return
	}

	responseWithJSON(w, 202, databaseRecipeToRecipe(recipe))
}

func (c *CookBookServer) getRecipeByIDHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	id, err := strconv.Atoi(idParam)

	recipe, err := c.store.GetRecipeByID(id)
	if err != nil {
		responseWithError(w, 404, "Recipe not found")
		return
	}

	responseWithJSON(w, 200, databaseRecipeToRecipe(recipe))
}

func (c *CookBookServer) getRecipesHandler(w http.ResponseWriter, r *http.Request) {
	// limitStr := r.URL.Query().Get("limit")
	// offsetStr := r.URL.Query().Get("offset")
	// var limit int32 = 10
	// var offset int32 = 0
	// if limitStr != "" {
	// 	var err error
	// 	limitInt, err := strconv.ParseInt(limitStr, 10, 32)
	// 	limit = int32(limitInt)
	// 	if err != nil {
	// 		responseWithError(w, 400, fmt.Sprint("Invalid limit"))
	// 		return
	// 	}
	// }
	// if offsetStr != "" {
	// 	var err error
	// 	offsetInt, err := strconv.ParseInt(offsetStr, 10, 32)
	// 	offset = int32(offsetInt)
	// 	if err != nil {
	// 		responseWithError(w, 400, fmt.Sprint("Invalid offset"))
	// 		return
	// 	}
	// }
	//
	recipes, err := c.store.GetRecipes()
	if err != nil {
		log.Printf("error getting recipes %v", err)
		// probably we should return 422 or something
		responseWithError(w, 404, fmt.Sprint("Recipes not found"))
		return
	}

	response := make([]Recipe, len(recipes))

	for i, recipe := range recipes {
		response[i] = databaseRecipeToRecipe(recipe)
	}

	responseWithJSON(w, 200, response)
}

func (c *CookBookServer) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}
	params := new(parameters)
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("error parsing JSON %v", err))
	}
	user, err := c.store.CreateUser(database.CreateUserParams{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Password:  params.Password,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Could't create user %v", err))
		return

	}

	responseWithJSON(w, 201, databaseUserToUser(user))
}
