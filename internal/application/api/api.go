package api

import (
	"github.com/sinakeshmiri/shortner/internal/ports"
)

// Application implements the APIPort interface
type Application struct {
	db    ports.DbPort
	shortner Shortner
}

// NewApplication creates a new Application
func NewApplication(db ports.DbPort, shortner Shortner) *Application {
	return &Application{db: db, shortner: shortner}
}

func (apia Application) NewURL(url,username string) (string, error) {
	urlID,err:=apia.shortner.Short(url)
	if err != nil {
		return "", err
	}
	apia.db.AddURL(url,string(urlID),username)
	return string(urlID), nil
}

func (apia Application) GetURL(id string) (string, error) {
	url,err:=apia.db.GetURL(id)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (apia Application) AddMetrics(id string) error {
	err:=apia.db.AddMetrics(id)
	if err != nil {
		return err
	}
	return nil
}

func (apia Application) DeleteURL(id,username string) error {
	err:=apia.db.DeleteURL(id,username)
	if err != nil {
		return err
	}
	return nil
}

func (apia Application) GetMetrics(username string) (map[string]int, error) {
	metrics,err:=apia.db.GetMetrics(username)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

