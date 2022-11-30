package config

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var DbSql *sql.DB
var DbMongo *mongo.Database

func OpenSqlDB(dbms string) error {
	var err error
	var db string
	var dbmsLog string
	switch dbms {
	case "postgres":
		dbmsLog = "PostgreSQL"
		db = "vinos"
		dbUser := "postgres"
		dbPassword := "password"
		DbSql, err = sql.Open("postgres", "postgres://"+dbUser+":"+dbPassword+"@localhost/"+db+"?sslmode=disable")
	case "mysql":
		dbmsLog = "MySQL"
		db = "vinos"
		dbUser := "root"
		dbPassword := "mysqlDGVDST99"
		DbSql, err = sql.Open("mysql", dbUser+":"+dbPassword+"@tcp(localhost:3307)/"+db+"?charset=utf8&parseTime=true")
	}
	if err != nil {
		log.Printf("*** %s: %v\n", dbmsLog, err)
		return err
	}
	if err = DbSql.Ping(); err != nil {
		log.Printf("*** %s: %v\n", dbmsLog, err)
		return err
	}
	log.Printf("*** %s: Conexión con base de datos '%s' establecida correctamente.\n", dbmsLog, db)
	return nil
}

func OpenMongoDB() error {
	url := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	if err != nil {
		log.Println("*** MongoDB:", err)
		return err
	}
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Println("*** MongoDB:", err)
		return err
	}
	DbMongo = client.Database("vinos")
	log.Println("*** MongoDB: Conexión con base de datos 'Vinos' establecida correctamente.")
	return nil
}
