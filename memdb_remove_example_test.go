package memdb_examples_test

import (
	"fmt"

	"github.com/hashicorp/go-memdb"
)

// Remove entries for a memdb.
func Example_removeEntries() {
	// Create the DB schema
	schema := createSchema()

	// Create empty database
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}

	addEntries(db)

	fmt.Println("Before:")
	printEntries(db)

	// Create write transaction
	txn := db.Txn(true)
	defer txn.Abort()

	// Entry must exist
	raw, err := txn.First("player", "name", "Sam")
	if err != nil {
		panic(err)
	}

	// Delete entry
	fmt.Printf("Removing player >%s< (>%s<)...\n", raw.(*player).Name, raw.(*player).Team)
	txn.Delete("player", raw)
	txn.Commit()

	fmt.Println("After:")
	printEntries(db)

	// Output:
	// Before:
	// All the players:
	//   Team A: Joe
	//   Team B: Sam
	// Removing player >Sam< (>Team B<)...
	// After:
	// All the players:
	//   Team A: Joe

}
