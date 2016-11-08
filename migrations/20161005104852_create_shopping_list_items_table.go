package main

import (
	"database/sql"
	"fmt"
)

// Up is executed when this migration is applied
func Up_20161005104852(txn *sql.Tx) {
	_, err := txn.Query(`CREATE TABLE shopping_list_items (
        id int(10) unsigned NOT NULL AUTO_INCREMENT,
        guid varchar(40) NOT NULL,
        user_guid varchar(40) NOT NULL,
		shopping_list_guid varchar(40) NOT NULL,
        name varchar(255) NOT NULL,
        category varchar(255) NOT NULL,
        sub_category varchar(255) NOT NULL,
        quantity int(6) NOT NULL,
        remark text DEFAULT NULL,
        added_from_deal int(1) NOT NULL DEFAULT 0,
        deal_guid varchar(40) DEFAULT NULL,
        cashback_amount decimal(4,2) DEFAULT NULL,
        add_to_cart int(1) DEFAULT 0 NOT NULL,
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
func Down_20161005104852(txn *sql.Tx) {
	_, err := txn.Query(`DROP TABLE shopping_list_items;`)

	if err != nil {
		fmt.Print(err)
	}
}
