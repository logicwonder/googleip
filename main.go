package main

import (
	"encoding/json"
	"fmt"

	scribble "github.com/nanobox-io/golang-scribble"
)

// a fish
type Fish struct{ Name string }

func main() {

	// create a new scribble database, providing a destination for the database to live
	db, _ := scribble.New("./fish", nil)

	// add some fish to the database
	for _, name := range []string{"onefish", "twofish", "redfish", "bluefish"} {
		db.Write("fish", name, Fish{Name: name})
	}

	// Read one fish from the database
	onefish := Fish{}
	db.Read("fish", "onefish", &onefish)

	fmt.Printf("It's a fish! %#v\n", onefish)

	// Read more fish from the database
	morefish, _ := db.ReadAll("fish")

	// iterate over morefish creating a new fish for each record
	fishies := []Fish{}
	for _, fish := range morefish {
		f := Fish{}
		json.Unmarshal([]byte(fish), &f)
		fishies = append(fishies, f)
	}

	fmt.Printf("It's a lot of fish! %#v\n", fishies)

	// Delete onefish from the database
	//db.Delete("fish", "onefish")

	// Delete all fish from the database
	//db.Delete("fish", "")
}
