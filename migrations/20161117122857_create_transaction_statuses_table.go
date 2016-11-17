package main

import (
	"database/sql"
	"fmt"
)

// Up is executed when this migration is applied
func Up_20161117122857(txn *sql.Tx) {
	_, err := txn.Query(`CREATE TABLE transaction_statuses (
        id int(10) unsigned NOT NULL AUTO_INCREMENT,
        guid varchar(40) NOT NULL,
		slug varchar(40) NOT NULL,
        name varchar(40) NOT NULL,
        created_at timestamp NULL DEFAULT NULL,
        updated_at timestamp NULL DEFAULT NULL,
        deleted_at timestamp NULL DEFAULT NULL,
        PRIMARY KEY (id),
        UNIQUE (guid)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8; 
    `)

	if err != nil {
		fmt.Print(err)
	}
}

// Down is executed when this migration is rolled back
func Down_20161117122857(txn *sql.Tx) {
	_, err := txn.Query(`DROP TABLE transaction_statuses;`)

	if err != nil {
		fmt.Print(err)
	}
}
