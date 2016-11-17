package main

import (
	"database/sql"
	"fmt"
)

// Up is executed when this migration is applied
func Up_20161117123305(txn *sql.Tx) {
	_, err := txn.Query(`CREATE TABLE cashout_transactions (
        id int(10) unsigned NOT NULL AUTO_INCREMENT,
        guid varchar(40) NOT NULL,
		user_guid varchar(40) NOT NULL,
        transaction_guid varchar(40) NOT NULL,
        bank_account_holder_name varchar(100) NULL DEFAULT NULL,
        bank_account_number varchar(50) NULL DEFAULT NULL,
        bank_name varchar(50) NULL DEFAULT NULL,
        bank_country varchar(50) NULL DEFAULT NULL,
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
func Down_20161117123305(txn *sql.Tx) {
	_, err := txn.Query(`DROP TABLE cashout_transactions;`)

	if err != nil {
		fmt.Print(err)
	}
}
