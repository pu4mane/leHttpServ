package serv

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

type srv struct {
	mu    *sync.Mutex
	stats map[uint]uint
}

func NewSrv() srv {
	return srv{
		mu:    &sync.Mutex{},
		stats: make(map[uint]uint),
	}
}

func (s *srv) Vote(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r := struct {
		CandidateId uint   `json:"candidate_id"`
		Passport    string `json:"passport"`
	}{}

	raw, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(raw, &r); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(r.Passport) == 0 || r.CandidateId == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	s.stats[r.CandidateId]++
	s.mu.Unlock()

	res.WriteHeader(http.StatusOK)
}

func (s *srv) Stats(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	s.mu.Lock()
	stats := s.stats
	s.mu.Unlock()

	raw, err := json.Marshal(stats)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Write(raw)

}
