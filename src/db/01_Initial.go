package db

import (
	"gopkg.in/go-pg/migrations.v4"
	"fmt"
	"log"
)

func init() {
	err := migrations.Register(func(db migrations.DB) error {
		fmt.Println("Initilizing User, Relation table...")
		_, err := db.Exec(`
		 DO $$
		 BEGIN
		 IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'usertype') THEN
		 	CREATE TYPE usertype AS ENUM ('user', 'admin');
		 END IF;
		 END
		 $$;
		 CREATE TABLE IF NOT EXISTS users (
		   id serial,
		   name varchar(30),
		   type usertype
		 );

		 DO $$
		 BEGIN
		 IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'relationtype') THEN
		    CREATE TYPE relationtype AS ENUM ('like', 'dislike', 'matched');
		 END IF;
		 END
		 $$;
		 CREATE TABLE IF NOT EXISTS relationships (
		   id serial,
		   uid integer,
		   oid integer,
		   type relationtype
		 );
		 `)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("Drop User, Relation table...")
		_, err := db.Exec(`
		DROP TABLE users;
		DROP TABLE relationships;
		DROP TYPE relationtype;
		DROP TYPE usertype;
		`)
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
}
