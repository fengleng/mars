package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gososy/sorpc/log"
	"sync"
)

var dbConn *sql.DB
var once sync.Once

type QueryResponse struct {
	Err  error                    `json:"errmsg"`
	Data []map[string]interface{} `json:"data"`
}
type ExecResponse struct {
	Err          error `json:"errcode"`
	InsertId     int64 `json:"insert_id"`
	RowsAffected int64 `json:"rows_affected"`
}

func GetDbConnection(User, Pass, Addr, DbName string, MaxIdleConns, MaxOpenConns int) (*sql.DB, error) { //addr like 127.0.0.1:123456
	var err error
	once.Do(func() {
		dsn := fmt.Sprint(User, ":", Pass, "@tcp(", Addr, ")/", DbName)
		dbConn, err = sql.Open("mysql", dsn)
		if err != nil {
			return
		}
		err = dbConn.Ping()
		if err != nil {
			return
		}
		dbConn.SetMaxIdleConns(MaxIdleConns)
		dbConn.SetMaxOpenConns(MaxOpenConns)
	})
	return dbConn, err
}
func FreeConnection(db *sql.DB) {
	db.Close()
}
func ExecMysql(db *sql.DB, resp *ExecResponse, sqlStmt string, args ...interface{}) {
	res, err := dbConn.Exec(sqlStmt, args...)
	if err != nil {
		log.Error(err)
		resp.Err = err
		return
	}
	id, err := res.LastInsertId()
	if err == nil {
		resp.InsertId = id
	}
	affected, err := res.RowsAffected()
	if err == nil {
		resp.RowsAffected = affected
	}
}
func QueryMysql(db *sql.DB, resp *QueryResponse, sqlStmt string, args ...interface{}) {
	rows, err := dbConn.Query(sqlStmt, args...)
	if err != nil {
		log.Error(err)
		resp.Err = err
		return
	}
	defer rows.Close()
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.Error(err)
		resp.Err = err
		return
	}
	resp.Data = make([]map[string]interface{}, 0)
	values := make([]sql.NullString, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Error(err)
			resp.Err = err
			return
		}
		data := make(map[string]interface{})
		for i, v := range values {
			if v.Valid {
				data[columns[i]] = v.String
			} else {
				data[columns[i]] = nil
			}
		}
		resp.Data = append(resp.Data, data)
	}
}
func EscapeString(sql string) string {
	dest := make([]byte, 0, 2*len(sql))
	var escape byte
	for i := 0; i < len(sql); i++ {
		c := sql[i]
		escape = 0
		switch c {
		case 0: /* Must be escaped for 'mysql' */
			escape = '0'
		case '\n': /* Must be escaped for logs */
			escape = 'n'
		case '\r':
			escape = 'r'
		case '\\':
			escape = '\\'
		case '\'':
			escape = '\''
		case '"': /* Better safe than sorry */
			escape = '"'
		case '\032': /* This gives problems on Win32 */
			escape = 'Z'
		}
		if escape != 0 {
			dest = append(dest, '\\', escape)
		} else {
			dest = append(dest, c)
		}
	}
	return string(dest)
}
