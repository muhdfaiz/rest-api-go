package main

import (
	"database/sql"
	"fmt"
)

// Up is executed when this migration is applied
func Up_20161103114727(txn *sql.Tx) {
	_, err := txn.Query(`CREATE TABLE deal_cashback_transactions (
        id int(10) unsigned NOT NULL AUTO_INCREMENT,
        guid varchar(40) NOT NULL,
        user_guid varchar(40) NOT NULL,
		receipt_id varchar(40) NOT NULL,
		receipt_image varchar(255) NOT NULL,
		verification_date timestamp NULL DEFAULT NULL,
		remark text DEFAULT NULL,
		status int(2) NOT NULL,
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
func Down_20161103114727(txn *sql.Tx) {
	_, err := txn.Query(`DROP TABLE deal_cashback_transactions;`)

	if err != nil {
		fmt.Print(err)
	}
}
