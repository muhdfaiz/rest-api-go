package main

import (
	"database/sql"
	"fmt"
)

// Up_20160831042931 is executed when this migration is applied
func Up_20160831042931(txn *sql.Tx) {
	_, err := txn.Query(`CREATE TABLE users (
        id int(10) unsigned NOT NULL AUTO_INCREMENT,
        guid varchar(255) NOT NULL,
        facebook_id varchar(100) DEFAULT NULL,
        name varchar(100) NOT NULL,
        email varchar(255) NOT NULL,
        phone_no varchar(20) NOT NULL,
        profile_picture varchar(255) NULL DEFAULT NULL,
        referral_code varchar(20) NULL DEFAULT NULL,
        bank_country varchar(50) NULL DEFAULT NULL,
        bank_name varchar(50) NULL DEFAULT NULL,
        bank_account_name varchar(50) NULL DEFAULT NULL,
        bank_account_number varchar(50) NULL DEFAULT NULL,
        register_by varchar(20) NOT NULL,
        verified int(1) unsigned NOT NULL DEFAULT 0,
        blacklist int(1) DEFAULT 0,
        created_at timestamp NULL DEFAULT NULL,
        updated_at timestamp NULL DEFAULT NULL,
        deleted_at timestamp NULL DEFAULT NULL,
        PRIMARY KEY (id),
        UNIQUE (guid),
        UNIQUE (facebook_id),
        UNIQUE (phone_no),
        UNIQUE (referral_code),
        KEY idx_users_deleted_at (deleted_at)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `)

	if err != nil {
		fmt.Print(err)
	}
}

// Down_20160831042931 is executed when this migration is rolled back
func Down_20160831042931(txn *sql.Tx) {
	_, err := txn.Query(`DROP TABLE users;`)

	if err != nil {
		fmt.Print(err)
	}
}
