package product

import (
	"net/http"
	"store-api/types"
	"store-api/utils"

	"github.com/gorilla/mux"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/product", h.handleCreateProduct).Methods(http.MethodPost)
	router.HandleFunc("/products", h.handleGetProducts).Methods(http.MethodGet)
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {

	ps, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, ps)
}
