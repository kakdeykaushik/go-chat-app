package db

// import (
// 	"chat-app/pkg/domain"
// 	"database/sql"
// 	"log"
// )

// type sqlite struct {
// }

// var db *sql.DB

// func GetDatabaseConnection() (*sql.DB, error) {
// 	db, err := sql.Open("sqlite3", "db.sqlite3")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return db, nil
// }

// func CreateDatabaseConnection() {
// 	log.Println("creating a database connection")
// 	conn, err := GetDatabaseConnection()
// 	if err != nil {
// 		log.Println("data base creation failed", err)
// 	}
// 	db = conn
// }

// func CreateDatabase() {
// 	// Create the "todo" table if it doesn't exist
// 	createTableSQL := `
// 		CREATE TABLE IF NOT EXISTS room (
// 		roomId TEXT PRIMARY KEY,
// 		member TEXT,
// 		);`

// 	_, e := db.Exec(createTableSQL)
// 	if e != nil {
// 		log.Fatal("creation of database failed")
// 	}

// }

// func newSqliteStorage() domain.Storage {
// 	panic("not implemented")
// 	// return &sqlite{}
// }

// func (s *sqlite) Get(K any) (any, error) {
// 	panic("not implemented")
// }

// func (s *sqlite) List() ([]*any, error) {
// 	panic("not implemented")
// }

// func (s *sqlite) Save(K any, V any) error {
// 	// stmt, e := db.Prepare("Insert Into todo Values (?,?,?)")
// 	// if e != nil {
// 	// 	log.Fatal("Task creation failed")
// 	// }
// 	// _, err := stmt.Exec(taskId, title, completed)
// 	// if err != nil {
// 	// 	log.Println("Task creation failed with error " + err.Error())
// 	// }
// 	// return nil
// }

// func (s *sqlite) Delete(K any) error {
// 	panic("not implemented")
// }
