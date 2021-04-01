package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"scratch-test/api/config"
	"scratch-test/api/router"
	"scratch-test/api/utils"
	"scratch-test/db"
	"scratch-test/models"
	"sync"
)

// Setup the basic stuffs in here
func setup(urls ...string) (db.IDatabase, error) {
	resolver, err := models.SetupClinic()
	if err != nil {
		return nil, err
	}

	database := db.CreateMemoryDB()
	err = database.Init(resolver)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	for _, u := range urls {

		resp, err := http.Get(u)
		if err != nil {
			return nil, err
		}

		var clinics []models.Clinic
		if err := json.NewDecoder(resp.Body).Decode(&clinics); err != nil {
			return nil, err
		}
		for i := range clinics {
			clinics[i].Init()
		}

		wg.Add(1)
		go func() {
			err = database.Update(func(tx db.ITx) error {
				for _, cl := range clinics {
					tx.Insert(cl)
				}
				return nil
			})

			if err != nil {
				log.Println(err)
				return
			}
			wg.Done()
		}()
	}
	// Perform other task while waiting for the database to be filled up

	// Wait for the database updating to complete if it has yet to
	wg.Wait()

	return database, nil
}

// Run initializes the server and runs it.
func Run() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	remoteData := []string{
		"https://storage.googleapis.com/scratchpay-code-challenge/dental-clinics.json",
		"https://storage.googleapis.com/scratchpay-code-challenge/vet-clinics.json",
	}
	database, err := setup(remoteData...)
	if err != nil {
		log.Println(err)
		return
	}
	config.Database = database

	listen(8080)
}

func listen(port int) {
	fmt.Println("Listening...:", port)
	r := router.NewRouter()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}

// Just for testing purpose
func sampleDB() {
	req := "http://localhost:8080?from=11:00&state=ak&limit=2&all=true"
	u, err := url.ParseRequestURI(req)
	if err != nil {
		log.Println(err)
		return
	}

	filter, limit := utils.FilterURLQuery(u.Query())

	entries, err := config.Database.Read(filter, limit)
	if err != nil {
		log.Println(err)
		return
	}

	bts, err := json.MarshalIndent(entries, "", "\t")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(bts))
}
