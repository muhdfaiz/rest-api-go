package main

import (
	"database/sql"
	"fmt"
)

// Up is executed when this migration is applied
func Up_20161102100625(txn *sql.Tx) {
	_, err := txn.Query(`CREATE TABLE deal_cashbacks (
        id int(10) unsigned NOT NULL AUTO_INCREMENT,
        guid varchar(40) NOT NULL,
        user_guid varchar(40) NOT NULL,
		shopping_list_guid varchar(40) NOT NULL,
		deal_guid varchar(40) NOT NULL,
		deal_cashback_transaction_guid varchar(40) DEFAULT NULL,
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
func Down_20161102100625(txn *sql.Tx) {
	_, err := txn.Query(`DROP TABLE deal_cashbacks;`)

	if err != nil {
		fmt.Print(err)
	}
}
