package db

import (
	"fmt"
	"gopkg.in/pg.v4"
	"config"
	"time"
)

type DataSource interface {
	Users() ([]User, error)
	CreateUser(string) (*User, error)

	UserRelationShips(int64) ([]RelationShip, error)
	CreateOrUpdateRelationShip(int64, int64, uint) error
}

type APIDataSource struct {
	conn *pg.DB
}

type User struct {
	Id     int64  `json:"id"`
	Name   string  `json:"name"`
	Type   string  `json:"type"`
}

func (u User) String() string {
	return fmt.Sprintf("User<%d %s %v>", u.Id, u.Name, u.Type)
}

type RelationShip struct {
	Id       int64
	Uid      int64 // FK to User
	Oid 	 int64 // FK to User
	Type 	 string
}

func (s RelationShip) String() string {
	return fmt.Sprintf("Story<%d %s %s>", s.Id, s.Uid, s.Oid)
}



func NewAPIDataSource(cfg *config.RestfulAPIConfig) *APIDataSource {
	opts := &pg.Options{
		User: cfg.DBUser,
		Password: cfg.DBPassword,
		Database: "api",
		Addr: fmt.Sprintf("%s:%d", cfg.DBHost, cfg.DBPort),
	}
	db := pg.Connect(opts).WithTimeout(30 * time.Second)
	return &APIDataSource{ db }
}

func (ds *APIDataSource) Users() ([]User, error) {
	var users []User
	err := ds.conn.Model(&users).Select()
	return users, err
}

func (ds *APIDataSource) CreateUser(name string) (*User, error) {
	user := &User{
		Name: name,
		Type: "user", // the specification is not clean, just using hard code
	}
	err := ds.conn.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ds *APIDataSource) UserRelationShips(uid int64) ([]RelationShip, error) {
	var rs []RelationShip
	err := ds.conn.Model(rs).Column("Oid", "Type").Where("Uid = ?", uid).Select()
	return rs, err
}

func (ds *APIDataSource) CreateOrUpdateRelationShip(uid, oid int64, tp uint) error {
	return nil
}

func (ds *APIDataSource) Close() {
	ds.conn.Close()
}
