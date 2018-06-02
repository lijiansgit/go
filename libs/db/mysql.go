// TODO
package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	// DefaultCharset 默认字符集
	DefaultCharset = "utf8"
)

type MySQL struct {
	Addr     string
	Username string
	Password string
	DBName   string
	Charset  string
	db       *sql.DB
}

func NewMySQL(addr, username, password, dbName string) *MySQL {
	return &MySQL{
		Addr:     addr,
		Username: username,
		Password: password,
		DBName:   dbName,
		Charset:  DefaultCharset,
	}
}

func (m *MySQL) initDB() (db *sql.DB, err error) {
	dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=%s`,
		m.Username, m.Password, m.Addr, m.DBName, m.Charset)
	db, err = sql.Open(`mysql`, dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (m *MySQL) Query(query string) (res string, err error) {
	db, err := m.initDB()
	if err != nil {
		return res, err
	}

	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return res, err
	}

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return res, err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return res, err
		}

		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			fmt.Println(columns[i], ": ", value)
		}
		fmt.Println("-----------------------------------")
	}
	if err = rows.Err(); err != nil {
		return res, err
	}

	return res, nil
}
