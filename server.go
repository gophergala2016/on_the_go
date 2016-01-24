/**
 * This file provided by Facebook is for non-commercial testing and evaluation
 * purposes only. Facebook reserves all rights not expressly granted.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * FACEBOOK BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
 * WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	
)

type map struct {
	
	name string `json:"name"`
	latitude   string `json:"latitude"`
	longitude   string `json:"longitude"`
	Menu   string `json:"Menu"`
}

const dataFile = "./spaza.json"

var mapMutex = new(sync.Mutex)

// Handle map
func handlemap(w http.ResponseWriter, r *http.Request) {
	// Since multiple requests could come in at once, ensure we have a lock
	// around all file operations
	mapMutex.Lock()
	defer mapMutex.Unlock()

	// Stat the file, so we can find its current permissions
	fi, err := os.Stat(dataFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to stat the data file (%s): %s", dataFile, err), http.StatusInternalServerError)
		return
	}

	// Read the map from the file.
	mapData, err := ioutil.ReadFile(dataFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read the data file (%s): %s", dataFile, err), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "POST":
		// Decode the JSON data
		var map []map
		if err := json.Unmarshal(mapData, &map); err != nil {
			http.Error(w, fmt.Sprintf("Unable to Unmarshal map from data file (%s): %s", dataFile, err), http.StatusInternalServerError)
			return
		}

		// Add a new map to the in memory slice of map
		//map = append(map, map{ID: time.Now().UnixNano() / 1000000, Author: r.FormValue("author"), Text: r.FormValue("text")})

		// Marshal the map to indented json.
		mapData, err = json.MarshalIndent(map, "", "    ")
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to marshal map to json: %s", err), http.StatusInternalServerError)
			return
		}

		// Write out the map to the file, preserving permissions
		err := ioutil.WriteFile(dataFile, mapData, fi.Mode())
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to write map to data file (%s): %s", dataFile, err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		io.Copy(w, bytes.NewReader(mapData))

	case "GET":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		// stream the contents of the file to the response
		io.Copy(w, bytes.NewReader(mapData))

	default:
		// Don't know the method, so error
		http.Error(w, fmt.Sprintf("Unsupported method: %s", r.Method), http.StatusMethodNotAllowed)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	http.HandleFunc("/api/map", handlemap)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Println("Server started: http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
