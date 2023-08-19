package ports

// APIPort is the technology neutral
// port for driving adapters
type APIPort interface {
	NewURL(url,username string) (string, error)
	GetURL(id string) (string, error)
	AddMetrics(id string) error
	DeleteURL(id,username string) error
	GetMetrics(username string) (map[string]int, error)
}