package main

import (
	"database/sql"
	"fmt"
)

// Up is executed when this migration is applied
func Up_20161117122948(txn *sql.Tx) {
	_, err := txn.Query(`CREATE TABLE transactions (
        id int(10) unsigned NOT NULL AUTO_INCREMENT,
        guid varchar(40) NOT NULL,
		user_guid varchar(40) NOT NULL,
        transaction_type_guid varchar(40) NOT NULL,
		transaction_status_guid varchar(40) NOT NULL,
        amount decimal(9,2) NOT NULL,
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
func Down_20161117122948(txn *sql.Tx) {
	_, err := txn.Query(`DROP TABLE transactions;`)

	if err != nil {
		fmt.Print(err)
	}
}
