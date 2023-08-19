package api

// Shortner
type Shortner interface {
	Short(url string) ([]byte, error)
}