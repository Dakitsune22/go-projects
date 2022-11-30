package users

import (
	"context"
	"log"
	"time"

	"github.com/Dakitsune22/go-projects/VinosByDak/config"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateUserSql(dbms string, user User) error {
	var query string
	switch dbms {
	case "postgres":
		query = "INSERT INTO users (id, password, nickname, name, fname, lname) VALUES ($1, $2, $3, $4, $5, $6)"
	case "mysql":
		query = "INSERT INTO users (id, password, nickname, name, fname, lname) VALUES (?, ?, ?, ?, ?, ?)"
	}
	_, err := config.DbSql.Exec(query, user.Id, user.Password, user.Nickname, user.Name, user.Fname, user.Lname)
	if err != nil {
		return err
	}
	log.Printf("*** CreateUser *** Usuario %s dado de alta correctamente.\n", user.Id)
	return nil
}

func CreateUserMongoDB(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	coll := config.DbMongo.Collection("users")
	result, err := coll.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	log.Printf("*** CreateUser *** Usuario %s dado de alta correctamente con %v.\n", user.Id, result.InsertedID)
	return nil
}

func GetUserSql(dbms string, userId string) (User, error) {
	var query string
	log.Println("*** GetUserSql *** UserId:", userId)
	log.Println("*** GetUserSql *** dbms:", dbms)
	switch dbms {
	case "postgres":
		query = "SELECT * FROM users WHERE id = $1"
	case "mysql":
		query = "SELECT * FROM users WHERE id = ?"
	}
	row := config.DbSql.QueryRow(query, userId)
	user := User{}
	err := row.Scan(&user.Id, &user.Password, &user.Nickname, &user.Name, &user.Fname, &user.Lname, &user.Bdate)
	if err != nil {
		return User{}, err
	}
	log.Println("*** GetUser *** date:", user.Bdate)
	log.Println("*** GetUser *** date (formatted):", user.Bdate.Format(DateYYYYMMDD))
	log.Printf("*** GetUser *** Usuario %s recuperado correctamente.\n", user.Id)
	return user, nil
}

func GetUserMongoDB(userId string) (User, error) {
	coll := config.DbMongo.Collection("users")
	filter := bson.M{"id": userId}
	doc := coll.FindOne(context.TODO(), filter)
	user := User{}
	err := doc.Decode(&user)
	if err != nil {
		return User{}, err
	}
	log.Printf("*** GetUser *** Usuario %s recuperado correctamente.\n", user.Id)
	return user, nil
}

func GetAllUsersSql(dbms string) ([]User, error) {
	query := "SELECT * FROM users"
	rows, err := config.DbSql.Query(query)
	if err != nil {
		return nil, err
	}
	users := make([]User, 0)
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.Id, &user.Password, &user.Nickname, &user.Name, &user.Fname, &user.Lname, &user.Bdate); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	//log.Println("*** GetAllUsersSql *** users:", users)
	log.Printf("*** GetAllUsers *** Número de usuarios recuperados: %d.\n", len(users))
	if err = rows.Err(); err != nil {
		return nil, err
	}
	log.Printf("*** GetAllUsers *** Número de usuarios recuperados: %d.\n", len(users))
	return users, nil
}

func GetAllUsersMongoDB() ([]User, error) {
	coll := config.DbMongo.Collection("users")
	filter := bson.M{}
	users := make([]User, 0)
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	cursor.All(context.TODO(), &users)
	if err != nil {
		return nil, err
	}
	log.Printf("*** GetAllUsers *** Número de usuarios recuperados: %d.\n", len(users))
	return users, nil
}

func DeleteUserSql(dbms string, userId string) error {
	var query string
	switch dbms {
	case "postgres":
		query = "DELETE FROM users WHERE id = $1"
	case "mysql":
		query = "DELETE FROM users WHERE id = ?"
	}
	_, err := config.DbSql.Exec(query, userId)
	if err != nil {
		return err
	}
	log.Printf("*** DeleteUser *** Usuario %s eliminado correctamente.\n", userId)
	return nil
}

func DeleteUserMongoDB(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	coll := config.DbMongo.Collection("users")
	filter := bson.M{"id": userId}
	result, err := coll.DeleteOne(ctx, filter)
	log.Println("*** UpdateUser *** Número de usuarios eliminados:", result.DeletedCount)
	return err
}

func UpdateUserSql(dbms string, user User) error {
	var query string
	switch dbms {
	case "postgres":
		query = "UPDATE users SET id = $1, password = $2, nickname = $3, name = $4, fname = $5, lname = $6, bdate = $7 WHERE id = $8"
	case "mysql":
		query = "UPDATE users SET id = ?, password = ?, nickname = ?, name = ?, fname = ?, lname = ?, bdate = ? WHERE id = ?"
	}
	_, err := config.DbSql.Exec(query, user.Id, user.Password, user.Nickname, user.Name, user.Fname, user.Lname, user.Bdate, user.Id)
	if err != nil {
		return err
	}
	log.Printf("*** UpdateUser *** Usuario %s actualizado correctamente.\n", user.Id)
	return nil
}

func UpdateUserMongoDB(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	coll := config.DbMongo.Collection("users")
	filter := bson.M{"id": user.Id}
	update := bson.M{"$set": user}
	result, err := coll.UpdateOne(ctx, filter, update)
	log.Println("*** UpdateUser *** Número de usuarios modificados:", result.ModifiedCount)
	return err
}

/*func ChangeUserPasswordSql(dbms string, user User) error {
	var query string
	switch dbms {
	case "postgres":
		query = "UPDATE users SET password = $1 WHERE id = $2"
	case "mysql":
		query = "UPDATE users SET password = ? WHERE id = ?"
	}
	_, err := config.DbSql.Exec(query, user.Password, user.Id)
	if err != nil {
		return err
	}
	log.Printf("*** ChangeUserPassword *** Password de usuario %s actualizado correctamente.\n", user.Id)
	return nil
}*/
