package db

import (
	"fmt"
	"gopkg.in/pg.v4"
	"config"
	"time"
	"log"
	"os"
)

const (
	RELATIONSHIP_TYPE_RS = "relationship"
	RELATIONSHIP_TYPE_WATCH = "watch"
	RELATIONSHIP_LIKE = "liked"
	RELATIONSHIP_DISLIKE = "disliked"
	RELATIONSHIP_MATCHED = "matched"
)

type DataSource interface {
	Users() ([]User, error)
	CreateUser(string) (*User, error)

	UserRelationShips(int64) ([]RelationShip, error)
	UserRelationShip(int64, int64, string)(*RelationShip, error)
	UserRelationShipCount(int64, int64, string)(int, error)
	CreateOrUpdateRelationShip(int64, int64, string, string) (*RelationShip, error)
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
	Id       int64 `json:"-"`
	Uid      int64 `json:"-"`// FK to User
	Oid 	 int64 `json:"user_id"`// FK to User
	State    string `json:"state"`
	Type 	 string `json:"type"`
}

func (s RelationShip) String() string {
	return fmt.Sprintf("RelationShip<%d %d %d>", s.Id, s.Uid, s.Oid)
}



func NewAPIDataSource(cfg *config.RestfulAPIConfig) *APIDataSource {
	opts := &pg.Options{
		User: cfg.DBUser,
		Password: cfg.DBPassword,
		Database: cfg.DBName,
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
	rs := []RelationShip{}
	err := ds.conn.Model(&rs).Column("uid", "oid", "state", "type").Where("uid = ?", uid).Select()
	return rs, err
}

func (ds *APIDataSource) UserRelationShip(uid, oid int64, tp string)(*RelationShip, error) {
	myRS := RelationShip{
		Uid: uid,
		Oid: oid,
		Type: tp,
	}
	err := ds.conn.Model(&myRS).Where("uid = ?", uid).Where("oid = ?", oid).Where("type = ?", tp).Select()
	return &myRS, err
}

func (ds *APIDataSource) UserRelationShipCount(uid, oid int64, tp string)(int, error) {
	myRS := RelationShip{
		Uid: uid,
		Oid: oid,
		Type: tp,
	}
	return ds.conn.Model(&myRS).Where("uid = ?", uid).Where("oid = ?", oid).Where("type = ?", tp).Count()
}


func (ds *APIDataSource) CreateOrUpdateRelationShip(uid, oid int64, st string, tp string) (*RelationShip, error) {
	rs := RelationShip{
		Uid: uid,
		Oid: oid,
		State: st,
		Type: tp,
	}
	created, err := ds.conn.Model(&rs).Where("uid = ?", uid).Where("oid = ?", oid).Where("type = ?", tp).SelectOrCreate()
	if err != nil {
		return &rs, err
	}
	if !created {
		rs.State = st
		err = ds.conn.Update(&rs)
		return &rs, err
	}
	return &rs, nil
}

func (ds *APIDataSource) Close() {
	ds.conn.Close()
}
