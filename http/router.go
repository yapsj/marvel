package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yapsj/marvel/service"
	"github.com/yapsj/marvel/util"
)

type Server struct {
	service *service.MarvelService
	logger  *log.Logger
	mux     *mux.Router
}

func NewServer(service *service.MarvelService, w io.Writer, mux *mux.Router) *Server {
	return &Server{
		service: service,
		logger:  util.NewInfoLogger(w),
		mux:     mux,
	}
}

func (s *Server) Serve() {

	router := s.mux

	router.HandleFunc("/characters", s.getCharactersHandler).Methods("GET")
	router.HandleFunc("/characters/{characterId}", s.getCharacterHandler).Methods("GET")
	//router.HandleFunc("/characters", s.clearCharactersCache).Methods("DELETE")
}

func (s *Server) getCharacterHandler(w http.ResponseWriter, r *http.Request) {
	service := s.service
	vars := mux.Vars(r)
	id := vars["characterId"]
	idLookup, err := strconv.Atoi(id)
	if err != nil {
		s.handleError(w, fmt.Errorf("%s is not a valid number", id), http.StatusBadRequest)
		return
	}

	character, err := service.GetCharacter(r.Context(), idLookup)
	if err != nil {
		s.handleError(w, err, http.StatusBadRequest)
		return
	}

	output, err := json.Marshal(character)
	if err != nil {
		s.handleError(w, err, http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(output)
}

func (s *Server) handleError(w http.ResponseWriter, err error, code int) {
	//log error message
	logger := s.logger
	logger.Printf("Error: %s\n", err.Error())

	errorResult := make(map[string]interface{})
	errorResult["code"] = code
	errorResult["error"] = err.Error()

	output, _ := json.Marshal(errorResult)
	http.Error(w, string(output), code)
}

/*
func (s *Server) clearCharactersCache(w http.ResponseWriter, r *http.Request) {
	err := s.service.ClearCharactersCache()
	if err != nil {
		s.handleError(w, err, http.StatusBadRequest)
		return
	}

	result := make(map[string]interface{})
	result["code"] = 200

	json.NewEncoder(w).Encode(result)

}*/

func (s *Server) getCharactersHandler(w http.ResponseWriter, r *http.Request) {
	service := s.service
	result, err := service.GetCharacters(r.Context())
	if err != nil {
		s.handleError(w, err, http.StatusBadRequest)
		return
	}

	output, err := json.Marshal(result)
	if err != nil {
		s.handleError(w, err, http.StatusBadRequest)
		return
	}

	w.Write(output)

}
