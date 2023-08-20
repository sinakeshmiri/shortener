package ports

// DbPort is the port for a db adapter
type DbPort interface {
	CloseDbConnection()
	AddURL(url,urlID,username string) error
	GetURL(id string) (string, error)
	DeleteURL(id,username string) error
	GetHits(username string) (map[string]int, error)
	AddHit(id string) error
}
