package db

import (
	"gopkg.in/go-pg/migrations.v4"
	"flag"
	"fmt"
)

func MigrateRun(db *APIDataSource) error {
	oldVersion, newVersion, err := migrations.Run(db.conn, flag.Args()...)
	if err != nil {
		return err
	}
	if newVersion != oldVersion {
		fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("version is %d\n", oldVersion)
	}
	return nil
}

