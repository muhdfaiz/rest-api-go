package main

import (
	"database/sql"
	"fmt"
)

// Up is executed when this migration is applied
func Up_20160923095545(txn *sql.Tx) {
	_, err := txn.Query(`CREATE TABLE occasions (
        id int(10) unsigned NOT NULL AUTO_INCREMENT,
        guid varchar(40) NOT NULL,
		slug varchar(50) NOT NULL,
        name varchar(255) NOT NULL,
        image varchar(255) NOT NULL,
		active int(1) unsigned NULL DEFAULT 0,
        created_at timestamp NULL DEFAULT NULL,
        updated_at timestamp NULL DEFAULT NULL,
        deleted_at timestamp NULL DEFAULT NULL,
        PRIMARY KEY (id),
        UNIQUE (guid),
		UNIQUE (slug)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8; 
    `)

	if err != nil {
		fmt.Print(err)
	}
}

// Down is executed when this migration is rolled back
func Down_20160923095545(txn *sql.Tx) {
	_, err := txn.Query(`DROP TABLE occasions;`)

	if err != nil {
		fmt.Print(err)
	}
}
