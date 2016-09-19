package main

import (
	"database/sql"
	"fmt"
)

// Up is executed when this migration is applied
func Up_20160905094332(txn *sql.Tx) {
	_, err := txn.Query(`CREATE TABLE sms_histories (
        id int(10) unsigned NOT NULL AUTO_INCREMENT,
        guid varchar(255) NOT NULL,
        user_guid varchar(255) NOT NULL,
        provider varchar(255) NOT NULL,
        sms_id varchar(100) NOT NULL,
        text varchar(255) NOT NULL,
        recipient_no varchar(20) NOT NULL,
        verification_code varchar(255) NOT NULL,
        status int(1) unsigned NOT NULL DEFAULT 0,
        created_at timestamp NULL DEFAULT NULL,
        updated_at timestamp NULL DEFAULT NULL,
        deleted_at timestamp NULL DEFAULT NULL,
        PRIMARY KEY (id),
        UNIQUE (guid),
        KEY idx_devices_user_guid (user_guid),
        KEY idx_devices_deleted_at (deleted_at)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `)

	if err != nil {
		fmt.Print(err)
	}
}

// Down is executed when this migration is rolled back
func Down_20160905094332(txn *sql.Tx) {
	_, err := txn.Query(`DROP TABLE sms_histories;`)

	if err != nil {
		fmt.Print(err)
	}
}
