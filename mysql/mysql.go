package mysql

import (
	"database/sql"
	"errors"
	"time"
)

type connector interface {
	Prepare(string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type DB struct {
	conn *sql.DB
	name string
}

func NewDB(name, dsn string, maxLifeTime time.Duration, maxOpenConn, maxIdleConn int) (*DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(maxLifeTime) //最大连接周期，超过时间的连接就close
	db.SetMaxOpenConns(maxOpenConn)    //设置最大连接数
	db.SetMaxIdleConns(maxIdleConn)    //设置闲置连接数
	return &DB{
		name: name,
		conn: db,
	}, nil
}

var NO_DATA_TO_BIND = errors.New("mysql: no data to bind")

func (this *DB) QueryOne(destObject interface{}, sql string, params ...interface{}) error {
	return queryOne(this.conn, destObject, sql, params)
}

func (this *DB) Query(destObject interface{}, sql string, params ...interface{}) error {
	return query(this.conn, destObject, sql, params)
}

func (this *DB) Insert(sql string, params ...interface{}) (int64, error) {
	return insert(this.conn, sql, params)
}

func (this *DB) Execute(sql string, params ...interface{}) (int64, error) {
	return execute(this.conn, sql, params)
}

func (this *DB) Begin() (*TX, error) {
	name := this.name + "-" + token()
	conn, err := this.conn.Begin()
	if err != nil {
		return nil, err
	}
	tx := &TX{name: name, conn: conn}
	return tx, nil
}

type TX struct {
	name string
	conn *sql.Tx
}
func (this *TX) QueryOne(destObject interface{}, sql string, params ...interface{}) error {
	return queryOne(this.conn, destObject, sql, params)
}

func (this *TX) Query(destObject interface{}, sql string, params ...interface{}) error {
	return query(this.conn, destObject, sql, params)
}

func (this *TX) Insert(sql string, params ...interface{}) (int64, error) {
	return insert(this.conn, sql, params)
}

func (this *TX) Execute(sql string, params ...interface{}) (int64, error) {
	return execute(this.conn, sql, params)
}
func (this *TX) Commit() error {
	if err := this.conn.Commit(); err != nil {
		return err
	}
	return nil
}

func (this *TX) Rollback() error {
	if err := this.conn.Rollback(); err != nil {
		return err
	}
	return nil
}
