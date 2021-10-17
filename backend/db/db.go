package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	isValidUserSQL = `
SELECT COUNT(id) as cnt 
FROM users 
WHERE username = ? and password = ?`

	getUserSQL = `
SELECT id, username, isadmin FROM users WHERE username = ? LIMIT 1
`
	getUserByIdSQL = `
SELECT id, username, isadmin FROM users WHERE id = ? LIMIT 1
`
	newRefreshToken = `
INSERT INTO refresh_tokens (userid, tokents) VALUES (?,?) ON CONFLICT(userid) DO UPDATE SET tokents = ?;
`
	getUserByRefreshToken = `
SELECT id, username, isadmin 
FROM users u 
INNER JOIN refresh_tokens r ON r.userid = u.id 
WHERE r.tokents=?;
`
	updatePassword = `
UPDATE users SET password = ?
 WHERE username = ?;
`
	schemaSQL = `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    username VARCHAR(32) UNIQUE,
    password VARCHAR(32),
    isadmin BOOLEAN
);
INSERT OR IGNORE INTO users (username, password, isadmin)
VALUES ('root', 'password', true);
CREATE TABLE IF NOT EXISTS refresh_tokens (
    userid INTEGER UNIQUE,
    tokents BIGINT
);
`
	createUser = `
INSERT INTO users (username, password, isadmin) VALUES(?, ?, ?);
`
	removeUser = `
DELETE FROM users WHERE id= ?;
`
	usersSQL = `
SELECT id, username, isadmin FROM users;
`
	updateUser = `
UPDATE users SET password = ?, isadmin = ?
 WHERE id = ?;
`
)

type User struct {
	Id       uint64
	Username string
	Password string
	IsAdmin  bool
}

type DB struct {
	sql *sql.DB
}

func NewDB(dbFile string) (*DB, error) {
	sqlDB, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	if _, err = sqlDB.Exec(schemaSQL); err != nil {
		return nil, err
	}

	db := DB{
		sql: sqlDB,
	}
	return &db, nil
}

func (db *DB) Close() error {
	defer func() {
		err := db.sql.Close()
		if err != nil {
			return
		}
	}()

	return nil
}

func (db *DB) IsValidUser(username string, password string) bool {
	stmt, err := db.sql.Prepare(isValidUserSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var cnt int
	err = stmt.QueryRow(username, password).Scan(&cnt)
	if err != nil {
		log.Fatal(err)
	}
	return cnt == 1
}

func (db *DB) GetUser(username string) (usr User, err error) {
	stmt, err := db.sql.Prepare(getUserSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&usr.Id, &usr.Username, &usr.IsAdmin)
	if err != nil {
		log.Println(err)
	}
	return
}

func (db *DB) GetUserById(id uint64) (usr User, err error) {
	stmt, err := db.sql.Prepare(getUserByIdSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&usr.Id, &usr.Username, &usr.IsAdmin)
	if err != nil {
		log.Println(err)
	}
	return
}

func (db *DB) NewRefreshToken(uid uint64, tokents int64) error {
	tx, err := db.sql.Begin()
	if err != nil {
		return err
	}
	stmt, err := db.sql.Prepare(newRefreshToken)
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.Stmt(stmt).Exec(uid, tokents, tokents)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}

func (db *DB) GetUserByRT(tokents int64) (usr User, err error) {
	stmt, err := db.sql.Prepare(getUserByRefreshToken)
	if err != nil {
		log.Fatal(err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	err = stmt.QueryRow(tokents).Scan(&usr.Id, &usr.Username, &usr.IsAdmin)
	if err != nil {
		log.Println(err)
	}
	return
}

func (db *DB) SetPassword(username string, password string) (usr User, err error) {
	tx, err := db.sql.Begin()
	if err != nil {
		return
	}
	stmt, err := db.sql.Prepare(updatePassword)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = tx.Stmt(stmt).Exec(password, username)
	if err != nil {
		err1 := tx.Rollback()
		if err1 != nil {
			return
		}
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	usr, err = db.GetUser(username)

	return
}

func (db *DB) CreateUser(username string, password string, isadmin bool) (usr User, err error) {
	tx, err := db.sql.Begin()
	if err != nil {
		return
	}
	stmt, err := db.sql.Prepare(createUser)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = tx.Stmt(stmt).Exec(username, password, isadmin)
	if err != nil {
		err1 := tx.Rollback()
		if err1 != nil {
			return
		}
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	usr, err = db.GetUser(username)

	return
}

func (db *DB) UpdateUser(id uint64, password string, isadmin bool) (usr User, err error) {
	tx, err := db.sql.Begin()
	if err != nil {
		return
	}
	stmt, err := db.sql.Prepare(updateUser)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = tx.Stmt(stmt).Exec(password, isadmin, id)
	if err != nil {
		err1 := tx.Rollback()
		if err1 != nil {
			return
		}
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	usr, err = db.GetUserById(id)

	return
}

func (db *DB) GetUsers() (users []User, err error) {
	rows, err := db.sql.Query(usersSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		p := User{}
		err := rows.Scan(&p.Id, &p.Username, &p.IsAdmin)
		if err != nil {
			fmt.Println(err)
			continue
		}
		users = append(users, p)
	}
	return
}

func (db *DB) RemoveUser(id uint64) (err error) {
	tx, err := db.sql.Begin()
	if err != nil {
		return
	}
	stmt, err := db.sql.Prepare(removeUser)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = tx.Stmt(stmt).Exec(id)
	if err != nil {
		err1 := tx.Rollback()
		if err1 != nil {
			return
		}
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}

	return
}
