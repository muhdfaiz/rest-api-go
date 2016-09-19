package main

import (
	"database/sql"
	"fmt"
)

// Up is executed when this migration is applied
func Up_20160901214114(txn *sql.Tx) {
	_, err := txn.Query(`CREATE TABLE devices (
        id int(10) unsigned NOT NULL AUTO_INCREMENT,
        guid varchar(255) NOT NULL,
        user_guid varchar(255) DEFAULT NULL,
        uuid varchar(255) NOT NULL,
        os varchar(100) NULL DEFAULT NULL,
        model varchar(255) DEFAULT NULL,
        push_token varchar(20) NOT NULL,
        app_version varchar(20) NOT NULL,
        token_expired int(1) unsigned NOT NULL DEFAULT 0,
        created_at timestamp NULL DEFAULT NULL,
        updated_at timestamp NULL DEFAULT NULL,
        deleted_at timestamp NULL DEFAULT NULL,
        PRIMARY KEY (id),
        UNIQUE (uuid),
        UNIQUE (push_token),
        UNIQUE (guid),
        KEY idx_devices_user_guid (user_guid),
        KEY idx_devices_deleted_at (deleted_at)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `)

	if err != nil {
		fmt.Print(err)
	}
}

// Down_20160901214114 is executed when this migration is rolled back
func Down_20160901214114(txn *sql.Tx) {
	_, err := txn.Query(`DROP TABLE devices;`)

	if err != nil {
		fmt.Print(err)
	}
}
