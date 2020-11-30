package app

import (
	"log"
	"encoding/json"
	"errors"
	"strconv"
	"net/http"
	"github.com/khiki1995/gosql/pkg/customers"
)

type Server struct {
	mux			 *http.ServeMux
	customersSvc *customers.Service
}

func NewServer(mux *http.ServeMux, customersSvc *customers.Service) *Server {
	return &Server{mux: mux, customersSvc: customersSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) Init() {
	s.mux.HandleFunc("/customers.getById", s.handleGetCustomerByID)
	s.mux.HandleFunc("/customers.getAll", s.handleGetAllCustomer)
	s.mux.HandleFunc("/customers.getAllActive", s.handleGetAllActiveCustomer)
	s.mux.HandleFunc("/customers.save", s.handleSaveCustomer)
	s.mux.HandleFunc("/customers.removeById", s.handleRemoveByIdCustomer)
	s.mux.HandleFunc("/customers.blockById", s.handleBlockByIdCustomer)
	s.mux.HandleFunc("/customers.unblockById", s.handleUnBlockByIdCustomer)
}
func (s *Server) handleUnBlockByIdCustomer(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	item, err := s.customersSvc.UnBlockByID(request.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}
func (s *Server) handleBlockByIdCustomer(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	item, err := s.customersSvc.BlockByID(request.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}
func (s *Server) handleRemoveByIdCustomer(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}

	item, err := s.customersSvc.RemoveByID(request.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	
	data, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}
func (s *Server) handleSaveCustomer(writer http.ResponseWriter, request *http.Request) {
	idParam := request.PostFormValue("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	var customer = &customers.Customer{
		ID: 		 id,
		Name: 	  	 request.PostFormValue("name"),
		Phone: 		 request.PostFormValue("phone"),		 
	}
	item, err := s.customersSvc.Save(request.Context(), customer)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	
	data, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}
func (s *Server) handleGetAllActiveCustomer(writer http.ResponseWriter, request *http.Request) {
	items, err := s.customersSvc.AllActive(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	
	data, err := json.Marshal(items)
	if err != nil {
		log.Println(err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}
func (s *Server) handleGetAllCustomer(writer http.ResponseWriter, request *http.Request) {
	items, err := s.customersSvc.All(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	
	data, err := json.Marshal(items)
	if err != nil {
		log.Println(err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}
func (s *Server) handleGetCustomerByID(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}

	item, err := s.customersSvc.ByID(request.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	
	data, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}