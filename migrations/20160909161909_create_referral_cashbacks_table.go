package main

import (
	"database/sql"
	"fmt"
)

// Up is executed when this migration is applied
func Up_20160909161909(txn *sql.Tx) {
	_, err := txn.Query(`CREATE TABLE referral_cashbacks (
        id int(10) unsigned NOT NULL AUTO_INCREMENT,
        guid varchar(255) NOT NULL,
        referrer_guid varchar(40) NOT NULL,
        referent_guid varchar(40) NOT NULL,
        cashback_amount decimal(4,2) NOT NULL,
        created_at timestamp NULL DEFAULT NULL,
        updated_at timestamp NULL DEFAULT NULL,
        deleted_at timestamp NULL DEFAULT NULL,
        PRIMARY KEY (id),
        UNIQUE (guid),
        KEY idx_referral_cashbacks_referrer_guid (referrer_guid),
        KEY idx_referral_cashbacks_referent_guid (referent_guid)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `)

	if err != nil {
		fmt.Print(err)
	}
}

// Down is executed when this migration is rolled back
func Down_20160909161909(txn *sql.Tx) {
	_, err := txn.Query(`DROP TABLE referral_cashbacks;`)

	if err != nil {
		fmt.Print(err)
	}
}
