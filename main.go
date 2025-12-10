/*
Purpose:
- OSM Notes to CSV File

Description:
- This program requests notes from the OSM database and saves them to a CSV file.

Releases:
- v1.0.0 - 2017/03/01 : initial release
- v1.0.1 - 2017/03/01 : license modified
- v1.0.2 - 2017/03/09 : layout modified
- v1.1.0 - 2017/11/03 : link to note added
- v1.2.0 - 2022/06/17 : modules, recompiled
- v1.3.0 - 2025/12/10 : program revised, go v1.25.5

Author:
- Klaus Tockloth

Copyright and license:
- Copyright (c) 2017-2025 Klaus Tockloth
- MIT license
*/

// main package
package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

// OSMNotes represents OSM notes structure
type OSMNotes struct {
	Type     string `json:"type"`
	Features []struct {
		Type     string `json:"type"`
		Geometry struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
		Properties struct {
			ID          int    `json:"id"`
			URL         string `json:"url"`
			CommentURL  string `json:"comment_url"`
			CloseURL    string `json:"close_url"`
			DateCreated string `json:"date_created"`
			Status      string `json:"status"`
			Comments    []struct {
				Date    string `json:"date"`
				UID     int    `json:"uid"`
				User    string `json:"user"`
				UserURL string `json:"user_url"`
				Action  string `json:"action"`
				Text    string `json:"text"`
				HTML    string `json:"html"`
			} `json:"comments"`
		} `json:"properties"`
	} `json:"features"`
}

// general program info
var (
	progName    = os.Args[0]
	progVersion = "v1.3.0"
	progDate    = "2025/12/10"
	progPurpose = "OSM Notes to CSV File"
	progInfo    = "Requests notes from the OSM database and saves them to a CSV file."
)

// command line options
var (
	bbox   = flag.String("bbox", "", "bounding box (left,bottom,right,top) (required)")
	limit  = flag.Int("limit", 999, "maximum number of notes")
	closed = flag.Bool("closed", false, "include closed notes")
)

/*
main starts this program.
*/
func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	flag.Usage = printProgUsage
	flag.Parse()

	if *bbox == "" {
		fmt.Fprintf(os.Stderr, "\nERROR: Option -bbox required.\n")
		printProgUsage()
	}

	if len(flag.Args()) < 1 {
		fmt.Fprintf(os.Stderr, "\nERROR: Output filename required.\n")
		printProgUsage()
	}

	csvfile := flag.Arg(0)

	osmBaseURI := "https://api.openstreetmap.org/api/0.6/notes.json"
	u, err := url.Parse(osmBaseURI)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	q.Set("bbox", *bbox)
	q.Set("limit", fmt.Sprintf("%d", *limit))
	if *closed {
		q.Set("closed", "1")
	} else {
		q.Set("closed", "0")
	}
	u.RawQuery = q.Encode()

	fmt.Printf("\nRequesting OSM notes ...\n")
	fmt.Printf("  URI  : %s\n", u.String())

	var netClient = &http.Client{
		Timeout: time.Second * 180,
	}

	// HTTP Request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Fatalf("error creating request: %v", err)
	}
	req.Header.Set("User-Agent", fmt.Sprintf("%s/%s (%s)", progName, progVersion, "freizeitkarte@googlemail.com"))

	resp, err := netClient.Do(req)
	if err != nil {
		log.Fatalf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("unexpected http status <%v>", resp.Status)
	}

	fmt.Printf("\nWriting CSV file ...\n")
	fmt.Printf("  FILE : %s\n", csvfile)

	var notes OSMNotes
	if err := json.NewDecoder(resp.Body).Decode(&notes); err != nil {
		log.Fatalf("error decoding JSON: %v", err)
	}

	if notes.Type != "FeatureCollection" {
		log.Fatalf("notes type <%v> unexpected", notes.Type)
	}

	file, err := os.OpenFile(csvfile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("error opening file <%s>: %v", csvfile, err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	header := []string{"Note", "Longitude", "Latitude", "Timestamp", "User", "Action", "Link", "Text"}
	if err := w.Write(header); err != nil {
		log.Fatalf("error writing header: %v", err)
	}

	numRecords := 0
	numNotes := 0

	for _, feature := range notes.Features {
		numNotes++
		lon := fmt.Sprintf("%.7f", feature.Geometry.Coordinates[0])
		lat := fmt.Sprintf("%.7f", feature.Geometry.Coordinates[1])
		noteID := fmt.Sprintf("%d", feature.Properties.ID)
		link := fmt.Sprintf("https://www.openstreetmap.org/note/%s#map=18/%s/%s&layers=N", noteID, lat, lon)

		for _, comment := range feature.Properties.Comments {
			user := comment.User
			if user == "" {
				user = "anonymous"
			}

			record := []string{
				noteID,
				lon,
				lat,
				comment.Date,
				user,
				comment.Action,
				link,
				comment.Text,
			}

			if err := w.Write(record); err != nil {
				log.Fatalf("error writing record: %v", err)
			}
			numRecords++
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatalf("error flushing CSV: %v", err)
	}

	fmt.Printf("  DONE : %d notes, %d records\n\n", numNotes, numRecords)
}

/*
printProgUsage prints program usage.
*/
func printProgUsage() {

	printProgInfo()

	fmt.Printf("\nUsage:\n")
	fmt.Printf("  %s <-bbox=lon,lat,lon,lat> [-limit=n] [-closed] <filename>\n", progName)

	fmt.Printf("\nExample:\n")
	fmt.Printf("  %s -bbox=7.47,51.84,7.78,52.06 osmnotes.csv\n", progName)

	fmt.Printf("\nOptions:\n")
	flag.PrintDefaults()

	fmt.Printf("\nArgument:\n")
	fmt.Printf("  filename string\n")
	fmt.Printf("        name of csv output file (required)\n")

	fmt.Printf("\nRemarks:\n")
	fmt.Printf("  A proxy server can be configured via environment variables:\n")
	fmt.Printf("  temporary: env HTTP_PROXY=http://proxy.server:port %s ...\n", progName)
	fmt.Printf("  permanent: export HTTP_PROXY=http://user:password@proxy.server:port\n")

	fmt.Printf("\nDisclaimer:\n")
	fmt.Printf("  The software is provided 'as is', without warranty of any kind, express or implied, including\n" +
		"  but not limited to the warranties of merchantability, fitness for a particular purpose and\n" +
		"  noninfringement. In no event shall the authors or copyright holders be liable for any claim,\n" +
		"  damages or other liability, whether in an action of contract, tort or otherwise, arising from,\n" +
		"  out of or in connection with the software or the use or other dealings in the software.\n\n")

	os.Exit(1)
}

/*
printProgInfo prints program info.
*/
func printProgInfo() {
	fmt.Printf("\nProgram:\n")
	fmt.Printf("  Name    : %s\n", progName)
	fmt.Printf("  Release : %s - %s\n", progVersion, progDate)
	fmt.Printf("  Purpose : %s\n", progPurpose)
	fmt.Printf("  Info    : %s\n", progInfo)
}
