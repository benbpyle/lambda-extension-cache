package routing

import (
	"cache-layer/clients"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

type Config struct {
	HttpClient      *http.Client
	CacheRepository *clients.CacheRepository
	DbRepository    *clients.DbRepository
}

// Start begins running the sidecar
func Start(port string, config *Config) {
	go startHTTPServer(port, config)
}

// Method that responds back with the cached values
func startHTTPServer(port string, config *Config) {
	r := chi.NewRouter()
	r.Get("/{key}", handleValue(config))

	logrus.Infof("Starting server on %s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("error starting the server")
		os.Exit(0)
	}
}

func handleValue(config *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		v := chi.URLParam(r, "key")
		m, err := config.CacheRepository.ReadCache(r.Context(), v)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if m == nil {
			logrus.Debug("Cache miss, reading from table")
			i, err := config.DbRepository.ReadItem(r.Context(), v)

			if err != nil || i == nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			config.CacheRepository.WriteCache(r.Context(), i)

			b, _ := json.Marshal(&i)
			w.Write(b)
		} else {
			logrus.Debug("Cache hit, returning from Momento")
			b, _ := json.Marshal(&m)
			w.Write(b)
		}
	}
}
