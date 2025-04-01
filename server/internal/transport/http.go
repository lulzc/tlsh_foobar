package transport

import (
	"encoding/json"
	"log"
	"net/http"
	"tlsh_foobar/server/internal/tlsh_foo"
)

type TlshItem struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
}

type Server struct {
	mux *http.ServeMux
}

func NewServer(tlshSvc *tlsh_foo.Service) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tlsh", func(w http.ResponseWriter, r *http.Request) {
		response, err := json.Marshal(tlshSvc.GetAll())
		if err != nil {
			log.Println(err)
		}
		_, err = w.Write(response)
		if err != nil {
			log.Println(err)
		}
	})

	// curl -i -X POST --data '{"name":"derName","hash":"T12345"}'  http://localhost:8080/tlsh
	mux.HandleFunc("POST /tlsh", func(writer http.ResponseWriter, request *http.Request) {
		var t TlshItem
		err := json.NewDecoder(request.Body).Decode(&t)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		tlshSvc.Add(t.Name, t.Hash)

		writer.WriteHeader(http.StatusCreated)
		return
	})

	mux.HandleFunc("GET /search", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query().Get("q")
		if query == "" {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		results := tlshSvc.Search(query)
		b, err := json.Marshal(results)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = writer.Write(b)
		if err != nil {
			log.Println(err)
			return
		}
	})

	mux.HandleFunc("GET /distance", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query().Get("q")
		if query == "" {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		results := tlshSvc.Distance(query)
		b, err := json.Marshal(results)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = writer.Write(b)
		if err != nil {
			log.Println(err)
			return
		}
	})

	return &Server{
		mux: mux,
	}
}

func (s *Server) Serve() error {
	return http.ListenAndServe(":8080", s.mux)

}
