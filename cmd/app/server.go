package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/manucher051299/crud/pkg/customers"
)

type Server struct {
	mux          *http.ServeMux
	customersSvc *customers.Service
}

func NewServer(mux *http.ServeMux, customersSvc *customers.Service) *Server {
	return &Server{mux: mux, customersSvc: customersSvc}
}
func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) Init() {
	log.Println("Init")
	s.mux.HandleFunc("/customers.getById", s.handleGetCustomerById)
	s.mux.HandleFunc("/customers.getAll", s.handleGetAll)
	s.mux.HandleFunc("/customers.getAllActive", s.handleGetAllActive)
	s.mux.HandleFunc("/customers.save", s.handleSave)
	s.mux.HandleFunc("/customers.removeById", s.handleRemoveCustomerById)
	s.mux.HandleFunc("/customers.blockById", s.handleBlockCustomerById)
	s.mux.HandleFunc("/customers.unblockById", s.handleUnblockCustomerById)

}
func (s *Server) handleGetCustomerById(writer http.ResponseWriter, request *http.Request) {
	log.Println("Get by id")
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item, err := s.customersSvc.ById(request.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleGetAll(writer http.ResponseWriter, request *http.Request) {
	log.Println("Get all")
	items, err := s.customersSvc.All(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(items)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleGetAllActive(writer http.ResponseWriter, request *http.Request) {
	log.Println("Get all Active")
	items, err := s.customersSvc.AllActive(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(items)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleSave(writer http.ResponseWriter, request *http.Request) {
	log.Println("Get all Active")
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	phone := request.URL.Query().Get("phone")
	name := request.URL.Query().Get("name")
	err = s.customersSvc.Save(request.Context(), name, phone, id)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleRemoveCustomerById(writer http.ResponseWriter, request *http.Request) {
	log.Println("Get by id")
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = s.customersSvc.RemoveById(request.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleBlockCustomerById(writer http.ResponseWriter, request *http.Request) {
	log.Println("Get by id")
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = s.customersSvc.BlockById(request.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleUnblockCustomerById(writer http.ResponseWriter, request *http.Request) {
	log.Println("Get by id")
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = s.customersSvc.UnblockById(request.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
