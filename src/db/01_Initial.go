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
		 IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'relationstate') THEN
		    CREATE TYPE relationstate AS ENUM ('liked', 'disliked', 'matched');
		 END IF;
		 END
		 $$;
		 DO $$
		 BEGIN
		 IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'relationtype') THEN
		    CREATE TYPE relationtype AS ENUM ('relationship', 'watch');
		 END IF;
		 END
		 $$;
		 CREATE TABLE IF NOT EXISTS relation_ships (
		   id serial,
		   uid integer,
		   oid integer,
		   state relationstate,
		   type relationtype
		 );
		 `)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("Drop User, Relation table...")
		_, err := db.Exec(`
		DROP TABLE users;
		DROP TABLE relation_ships;
		DROP TYPE relationtype;
		DROP TYPE relationstate;
		DROP TYPE usertype;
		`)
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
}
