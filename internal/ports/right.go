package ports

// DbPort is the port for a db adapter
type DbPort interface {
	CloseDbConnection()
	AddURL(url,urlID,username string) (string, error)
	GetURL(id string) (string, error)
	DeleteURL(id,username string) error
	GetMetrics(username string) (map[string]int, error)
	AddMetrics(id string) error
}
