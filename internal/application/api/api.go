package api

import (
	"github.com/sinakeshmiri/shortener/internal/ports"
)

// Application implements the APIPort interface
type Application struct {
	db        ports.DbPort
	shortener shortener
	urlschan  chan string
}

// NewApplication creates a new Application
func NewApplication(db ports.DbPort, shortener shortener, urlschan chan string) *Application {
	go shortener.Short()
	return &Application{db: db, shortener: shortener, urlschan: urlschan}
}

func (apia Application) NewURL(url, username string) (string, error) {

	urlID := <-apia.urlschan
	apia.db.AddURL(url, urlID, username)
	return string(urlID), nil
}

func (apia Application) GetURL(id string) (string, error) {
	url, err := apia.db.GetURL(id)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (apia Application) AddMetrics(id string) error {
	err := apia.db.AddHit(id)
	if err != nil {
		return err
	}
	return nil
}

func (apia Application) DeleteURL(id, username string) error {
	err := apia.db.DeleteURL(id, username)
	if err != nil {
		return err
	}
	return nil
}

func (apia Application) GetMetrics(username, id string) (map[string]int, error) {
	metrics, err := apia.db.GetHits(username, id)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}
