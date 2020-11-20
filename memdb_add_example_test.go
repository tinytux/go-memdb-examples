package memdb_examples_test

import (
	"fmt"

	"github.com/hashicorp/go-memdb"
)

type TeamType string

const (
	TeamA TeamType = "Team A"
	TeamB TeamType = "Team B"
)

type player struct {
	ID   string
	Name string
	Team TeamType
}

func createSchema() *memdb.DBSchema {
	return &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"player": &memdb.TableSchema{
				Name: "player",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
					"name": &memdb.IndexSchema{
						Name:    "name",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Name"},
					},
					"team": &memdb.IndexSchema{
						Name:    "team",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Team"},
					},
				},
			},
		},
	}
}

func addEntries(db *memdb.MemDB) {
	// Start write transaction
	txn := db.Txn(true)

	// Add players
	players := []*player{
		&player{ID: "P01", Name: "Joe", Team: TeamA},
		&player{ID: "P02", Name: "Sam", Team: TeamB},
	}
	for _, p := range players {
		if err := txn.Insert("player", p); err != nil {
			panic(err)
		}
	}

	// Commit the transaction
	txn.Commit()
}

func printEntries(db *memdb.MemDB) {
	// Create read-only transaction
	txn := db.Txn(false)
	defer txn.Abort()

	// List all the players
	it, err := txn.Get("player", "id")
	if err != nil {
		panic(err)
	}

	fmt.Println("All the players:")
	for obj := it.Next(); obj != nil; obj = it.Next() {
		p := obj.(*player)
		fmt.Printf("  %v: %s\n", p.Team, p.Name)
	}
}

// Create and initialize new memdb, add and query entries.
func Example_addEntries() {
	// Create the DB schema
	schema := createSchema()

	// Create empty database
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}

	addEntries(db)
	printEntries(db)

	// Output:
	// All the players:
	//   Team A: Joe
	//   Team B: Sam
}
