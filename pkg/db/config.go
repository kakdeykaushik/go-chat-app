package db

type Configuration struct {
	DBName     string
	Collection string
	Uid        string // wierd, but ok
}

func NewConfiguration(dbName, col, uid string) *Configuration {
	return &Configuration{DBName: dbName, Collection: col, Uid: uid}
}
