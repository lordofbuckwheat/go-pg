package main

import (
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"time"
)

var db *sqlx.DB

func main() {
	var connectionString string
	flag.StringVar(&connectionString, "connection", "", "")
	flag.Parse()
	var err error
	if db, err = sqlx.Connect("mysql", connectionString); err != nil {
		panic(err)
	} else {
		if err := db.Ping(); err != nil {
			panic(err)
		}
		var eg errgroup.Group
		for i := 0; i < 1000; i++ {
			var i = i
			eg.Go(func() error {
				return Transaction(func(m *NewModel) error {
					fmt.Println("transaction started", i)
					/*if _, err := m.DB.Exec("SELECT id FROM users"); err != nil {
						return errors.WithStack(err)
					}
					fmt.Println("query executed", i)*/
					m.OnAfterCommit(func() error {
						fmt.Println("transaction finished", i)
						return nil
					})
					m.OnAfterRollback(func(err error) {
						fmt.Println("transaction failed", i)
					})
					time.Sleep(5 * time.Second)
					return nil
				})
			})
		}
		if err := eg.Wait(); err != nil {
			panic(err)
		}
	}
	fmt.Println("done")
	<-make(chan struct{})
}

type NewModel struct {
	DB                       *sqlx.Tx
	onAfterCommit            []func() error
	onAfterRollback          []func(error)
	onBeforeCommitOrRollback []func()
}

func Transaction(f func(*NewModel) error) error {
	var noPanic bool
	tx, err := db.Beginx()
	if err != nil {
		return errors.WithStack(err)
	}
	var m = &NewModel{DB: tx}
	defer func() {
		if !noPanic {
			_ = tx.Rollback()
		}
	}()
	if err := f(m); err != nil {
		for _, obcor := range m.onBeforeCommitOrRollback {
			obcor()
		}
		_ = tx.Rollback()
		noPanic = true
		for _, oar := range m.onAfterRollback {
			oar(err)
		}
		return err
	}
	for _, obcor := range m.onBeforeCommitOrRollback {
		obcor()
	}
	if err := tx.Commit(); err != nil {
		noPanic = true
		for _, oar := range m.onAfterRollback {
			oar(err)
		}
		return errors.WithStack(err)
	}
	noPanic = true
	for _, oac := range m.onAfterCommit {
		if err := oac(); err != nil {
			return err
		}
	}
	return nil
}

func (m *NewModel) OnAfterCommit(f func() error) {
	m.onAfterCommit = append(m.onAfterCommit, f)
}

func (m *NewModel) OnAfterRollback(f func(error)) {
	m.onAfterRollback = append(m.onAfterRollback, f)
}

func (m *NewModel) OnBeforeCommitOrRollback(f func()) {
	m.onBeforeCommitOrRollback = append(m.onBeforeCommitOrRollback, f)
}
