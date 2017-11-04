/*
Purpose:
- OSM-Notes to CSV-File

Description:
- This program requests notes from the OSM database and stores them into a CSV file.

Releases:
- 1.0.0 - 2017/03/01 : initial release
- 1.0.1 - 2017/03/01 : license modified
- 1.0.2 - 2017/03/09 : layout modified
- 1.1.0 - 2017/11/03 : link to note added

Author:
- Klaus Tockloth

Copyright and license:
- Copyright (c) 2017 Klaus Tockloth
- MIT license

Permission is hereby granted, free of charge, to any person obtaining a copy of this software
and associated documentation files (the Software), to deal in the Software without restriction,
including without limitation the rights to use, copy, modify, merge, publish, distribute,
sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or
substantial portions of the Software.

The software is provided 'as is', without warranty of any kind, express or implied, including
but not limited to the warranties of merchantability, fitness for a particular purpose and
noninfringement. In no event shall the authors or copyright holders be liable for any claim,
damages or other liability, whether in an action of contract, tort or otherwise, arising from,
out of or in connection with the software or the use or other dealings in the software.

Contact:
- freizeitkarte@googlemail.com

Remarks:
- API description: http://wiki.openstreetmap.org/wiki/API_v0.6#Map_Notes_API

Links:
- https://github.com/Klaus-Tockloth/osmnotes2csv
*/

package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

// OSMNotes (generated with https://mholt.github.io/json-to-go/)
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
	progVersion = "1.1.0"
	progDate    = "2017/11/03"
	progOwner   = "Copyright (c) 2017 Klaus Tockloth"
	progLicense = "MIT license"
	progPurpose = "OSM-Notes -> CSV-File"
	progInfo    = "Requests notes from the OSM database and stores them into a CSV file."
	progContact = "freizeitkarte@googlemail.com"
	progLink    = "https://github.com/Klaus-Tockloth/osmnotes2csv"
)

// debugging
var debug = false

// command line options
var (
	bbox   = flag.String("bbox", "", "bounding box (left,bottom,right,top) (required)")
	limit  = flag.Int("limit", 999, "maximum number of notes")
	closed = flag.Bool("closed", false, "include closed notes")
)

/*
Initialize program.
*/
func init() {

	// init Logger
	log.SetPrefix("\nFATAL ERROR ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
}

/*
Start program.
*/
func main() {

	flag.Usage = printProgUsage
	flag.Parse()

	if *bbox == "" {
		fmt.Printf("\nERROR:\n")
		fmt.Printf("  Option -bbox required.\n")
		printProgUsage()
	}

	if len(flag.Args()) < 1 {
		fmt.Printf("\nERROR:\n")
		fmt.Printf("  Output filename required.\n")
		printProgUsage()
	}

	csvfile := flag.Arg(0)

	includeClosed := 0
	if *closed {
		includeClosed = 1
	}

	osmBaseURI := "http://api.openstreetmap.org/api/0.6/notes.json"
	osmRequestURI := fmt.Sprintf("%s?bbox=%s&limit=%d&closed=%d", osmBaseURI, *bbox, *limit, includeClosed)

	fmt.Printf("\nRequesting OSM notes ...\n")
	fmt.Printf("  URI  : %s\n", osmRequestURI)

	var netClient = &http.Client{
		Timeout: time.Second * 180,
	}

	resp, err := netClient.Get(osmRequestURI)
	if err != nil {
		log.Fatalf("error <%v> at netClient.Get()", err)
	}

	if resp.Status != "200 OK" {
		log.Fatalf("http status <%v> not expected", resp.Status)
	}

	if debug {
		dumpBody := true
		var dump []byte
		dump, err = httputil.DumpResponse(resp, dumpBody)
		if err != nil {
			log.Fatalf("error <%v> at httputil.DumpResponse()", err)
		}
		fmt.Printf("\nResponse dump (body = %v) ...\n%s\n", dumpBody, dump)
	}

	rb, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatalf("error <%v> at ioutil.ReadAll()", err)
	}

	fmt.Printf("\nWriting CSV file ...\n")
	fmt.Printf("  FILE : %s\n", csvfile)

	var notes OSMNotes
	if err = json.Unmarshal(rb, &notes); err != nil {
		log.Fatalf("error <%v> at json.Unmarshal()", err)
	}

	if notes.Type != "FeatureCollection" {
		log.Fatalf("notes type <%v> unexpected", notes.Type)
	}

	if debug {
		for _, feature := range notes.Features {
			fmt.Printf("----------------------------------------\n\n")
			fmt.Printf("feature.Type                   = %v\n", feature.Type)
			fmt.Printf("feature.Geometry.Type          = %v\n", feature.Geometry.Type)
			fmt.Printf("feature.Geometry.Coordinates   = %v\n", feature.Geometry.Coordinates)
			fmt.Printf("feature.Properties.ID          = %v\n", feature.Properties.ID)
			fmt.Printf("feature.Properties.URL         = %v\n", feature.Properties.URL)
			fmt.Printf("feature.Properties.CommentURL  = %v\n", feature.Properties.CommentURL)
			fmt.Printf("feature.Properties.CloseURL    = %v\n", feature.Properties.CloseURL)
			fmt.Printf("feature.Properties.DateCreated = %v\n", feature.Properties.DateCreated)
			fmt.Printf("feature.Properties.Status      = %v\n", feature.Properties.Status)
			fmt.Printf("\n")
			for _, comment := range feature.Properties.Comments {
				fmt.Printf("  comment.Date    = %v\n", comment.Date)
				fmt.Printf("  comment.UID     = %v\n", comment.UID)
				fmt.Printf("  comment.User    = %v\n", comment.User)
				fmt.Printf("  comment.UserURL = %v\n", comment.UserURL)
				fmt.Printf("  comment.Action  = %v\n", comment.Action)
				fmt.Printf("  comment.Text    = %v\n", comment.Text)
				fmt.Printf("  comment.HTML    = %v\n", comment.HTML)
				fmt.Printf("\n")
			}
		}
	}

	// O_WRONLY = open the file write-only
	// O_TRUNC = if possible, truncate file when opened
	// Modus = O_CREATE = create a new file if none exists
	// 0666 = read & write
	file, err := os.OpenFile(csvfile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error <%v> at os.OpenFile(); file = <%v>", err, csvfile)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalf("error <%v> at file.Close(); file = <%v>", err, csvfile)
		}
	}()

	// CSV writer
	w := csv.NewWriter(file)

	// CSV record buffer
	record := make([]string, 8)

	// CSV header
	header := []string{"Note", "Longitude", "Latitude", "Timestamp", "User", "Action", "Link", "Text"}
	if err := w.Write(header); err != nil {
		log.Fatalf("error <%v> at w.Write()", err)
	}

	numRecords := 0
	numNotes := 0
	for _, feature := range notes.Features {
		record[0] = fmt.Sprintf("%v", feature.Properties.ID)
		record[1] = fmt.Sprintf("%v", feature.Geometry.Coordinates[0])
		record[2] = fmt.Sprintf("%v", feature.Geometry.Coordinates[1])
		numNotes++

		for _, comment := range feature.Properties.Comments {
			record[3] = comment.Date
			user := comment.User
			if user == "" {
				user = "anonym"
			}
			record[4] = user
			record[5] = comment.Action
			record[6] = fmt.Sprintf("https://www.openstreetmap.org/note/%v#map=18/%v/%v&layers=N",
				feature.Properties.ID, feature.Geometry.Coordinates[1], feature.Geometry.Coordinates[0])
			record[7] = comment.Text

			// CSV record
			if err := w.Write(record); err != nil {
				log.Fatalf("error <%v> at w.Write()", err)
			}
			numRecords++
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatalf("error <%v> at w.Flush()", err)
	}

	fmt.Printf("  DONE : %d notes, %d records\n\n", numNotes, numRecords)
}

/*
Print program usage.
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
	fmt.Printf("  A proxy server can be configured via the program environment:\n")
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
Print program info.
*/
func printProgInfo() {

	fmt.Printf("\nProgram:\n")
	fmt.Printf("  Name    : %s\n", progName)
	fmt.Printf("  Release : %s - %s\n", progVersion, progDate)
	fmt.Printf("  Purpose : %s\n", progPurpose)
	fmt.Printf("  Info    : %s\n", progInfo)
	fmt.Printf("  Owner   : %s\n", progOwner)
	fmt.Printf("  License : %s\n", progLicense)
	fmt.Printf("  Contact : %s\n", progContact)
	fmt.Printf("  Link    : %s\n", progLink)
}
