package music

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

func SeedMusic(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch {
		case r.Method == "POST":
			CreateAllMusic(w, r)
			return
		case r.Method == "PUT":
			RefreshAllMusic(w, r)
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(fmt.Sprintf("Method not supported %v", r.Method)))
			return
		}
	}
}

func CreateAllMusic(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Posted to seed"))
}

func RefreshAllMusic(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Posted to seed"))
}
