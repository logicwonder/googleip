package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"

	scribble "github.com/nanobox-io/golang-scribble"
	"gopkg.in/resty.v1"
)

// a fish
type Goog struct {
	SyncToken    string   `json:"syncToken"`
	CreationTime string   `json:"creationTime"`
	Prefixes     []Prefix `json:"prefixes"`
}

type Prefix struct {
	IPv4Prefix string `json:"ipv4Prefix,omitempty"`
	IPv6Prefix string `json:"ipv6Prefix,omitempty"`
}

var client = resty.New().SetTimeout(10 * time.Second)

func main() {
	refresh := flag.Bool("refresh", false, "Refresh IP whitelist")
	dump := flag.String("dump", "", "Dump the allow list. Possible values: ipv4/ipv6")
	flag.Parse()

	fmt.Printf(*dump)
	// create a new scribble database, providing a destination for the database to live
	db, err := scribble.New("./scribbledb", nil)
	if err != nil {
		panic(fmt.Errorf("failed create/open scribbledb: %w", err))
	}

	currentGoog := Goog{}

	if *dump == "ipv4" {

		db.Read("goog", "goog_current", &currentGoog)
		currentIPV4List, _ := getIPWhitelist(currentGoog)
		for _, ip := range currentIPV4List {
			fmt.Println(ip)
		}
		return
	} else if *dump == "ipv6" {
		db.Read("goog", "goog_current", &currentGoog)
		_, currentIPV6List := getIPWhitelist(currentGoog)
		for _, ip := range currentIPV6List {
			fmt.Println(ip)
		}
		return
	} else if *refresh {

		resp, err := client.R().
			Get("https://www.gstatic.com/ipranges/goog.json")
		if err != nil {
			panic(fmt.Errorf("failed to get goog: %w", err))
		}

		if resp.StatusCode() != http.StatusOK {
			panic(fmt.Errorf("Failed to get goog: " + string(resp.Body())))
		}

		newGoog := Goog{}
		err = json.Unmarshal(resp.Body(), &newGoog)
		if err != nil {
			panic(fmt.Errorf("Failed to parse goog response: " + string(resp.Body())))
		}

		//Read current goog
		db.Read("goog", "goog_current", &currentGoog)

		if currentGoog.SyncToken == newGoog.SyncToken {
			fmt.Println("No changes in IP whitelist.")
			return
		} else {
			//Save the old goog
			fmt.Printf("\nArchiving existing IP whitelist. Created On: %s", currentGoog.CreationTime)
			err = db.Write("goog", "goog_"+currentGoog.SyncToken, currentGoog)
			if err != nil {
				fmt.Printf("\nFailed to archive existing IP whitelist: %s", err.Error())
			}
			fmt.Printf("\nArchived to goog_" + currentGoog.SyncToken + ".json")

			// Update the current goog
			db.Write("goog", "goog_current", newGoog)
			if err != nil {
				fmt.Printf("\nFailed to update IP whitelist: %s", err.Error())
			}
			fmt.Printf("Updated IP whitelist. Created On: %s", newGoog.CreationTime)

			currentIPV4List, currentIPV6List := getIPWhitelist(currentGoog)
			newIPV4List, newIPV6List := getIPWhitelist(newGoog)

			additionIPV4List := difference(newIPV4List, currentIPV4List)
			if len(additionIPV4List) > 0 {
				fmt.Printf("\n IPV4 Addition: %v", additionIPV4List)
			} else {
				fmt.Println("\nNo new IPV4")
			}

			deletedIPV4List := difference(currentIPV4List, newIPV4List)
			if len(deletedIPV4List) > 0 {
				fmt.Printf("\n IPV4 Deletion: %v", deletedIPV4List)
			} else {
				fmt.Println("\nNo deletion of IPV4")
			}

			additionIPV6List := difference(newIPV6List, currentIPV6List)
			if len(additionIPV6List) > 0 {
				fmt.Printf("\n IPV6 Addition: %v", additionIPV6List)
			} else {
				fmt.Println("\nNo new IPV6")
			}

			deletedIPV6List := difference(currentIPV6List, newIPV6List)
			if len(deletedIPV6List) > 0 {
				fmt.Printf("\n IPV6 Deletion: %v", deletedIPV6List)
			} else {
				fmt.Println("\nNo deletion of IPV6")
			}
		}
	} else {
		flag.PrintDefaults()
	}
}

//Returns the IPV4 and IPv6 list as string arrays
func getIPWhitelist(goog Goog) ([]string, []string) {
	var ipv4List, ipv6List []string
	for _, p := range goog.Prefixes {
		if p.IPv4Prefix != "" {
			ipv4List = append(ipv4List, p.IPv4Prefix)
		} else if p.IPv6Prefix != "" {
			ipv6List = append(ipv6List, p.IPv6Prefix)
		}

	}
	return ipv4List, ipv6List
}

func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
