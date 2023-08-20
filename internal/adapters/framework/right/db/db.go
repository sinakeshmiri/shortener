package db

import (
	"time"

	"github.com/gocql/gocql"
)


// Adapter implements the DbPort interface
type Adapter struct {
	session *gocql.Session
}

// NewAdapter creates a new Adapter
func NewAdapter(addresses[]string, keyspace string) (*Adapter, error) {
	// connect to the cluster
	///// Set up the ScyllaDB cluster configuration
	cluster := gocql.NewCluster(addresses...) // Replace with your ScyllaDB node(s) addresses
	cluster.Keyspace = keyspace          // Replace with your keyspace name
	cluster.Timeout = 5 * time.Second

	// Create a session to interact with the database
	session, err := cluster.CreateSession()
	if err != nil {
		return nil,err
	}
	return &Adapter{session: session}, nil
}

func (da Adapter) CloseDbConnection() {
	da.session.Close()
}


func(da Adapter)GetURL(id string) (string, error){
	var url string
	query := da.session.Query("SELECT url FROM urls WHERE id = ?", id)
	iter := query.Iter()
	for iter.Scan(&url){}
	if err := iter.Close(); err != nil {
		return "",err
	}
	return url, nil
}

func (da Adapter)AddURL(url,urlID,username string)  error{
	err := da.session.Query("INSERT INTO urls (id, url, username) VALUES (?, ?, ?)", urlID, url, username).Exec()
	if err != nil {
		return err
	}
	return  nil
}
func (da Adapter) DeleteURL(id,username string) error{
	err := da.session.Query("DELETE FROM urls WHERE id = ? AND username = ?", id, username).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (da Adapter)GetHits(username string) (map[string]int, error){
	var id string
	var hits int
	metrics:=make(map[string]int)
	query := da.session.Query("SELECT id,hits FROM metrics WHERE username = ?", username)
	iter := query.Iter()
	for iter.Scan(&id,&hits){
		metrics[id]=hits
	}
	if err := iter.Close(); err != nil {
		return nil,err
	}
	return metrics, nil
}

func (da Adapter) AddHit(id string) error{
	err := da.session.Query("UPDATE metrics SET hits = hits + 1 WHERE id = ?", id).Exec()
	if err != nil {
		return err
	}
	return  nil
}
	
	