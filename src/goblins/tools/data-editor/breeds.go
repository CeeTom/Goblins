package main

import (
	"encoding/json"
	"fmt"
	"goblins/game"
	"goblins/game/dataio"
	"net/http"
	"os"
	"path"
	"sync"
)

type BreedList struct {
	sync.RWMutex
	Breeds []*game.Breed
}

func LoadBreedList(brdDir string) (*BreedList, error) {
	ret := new(BreedList)
	breeds, err := dataio.ReadAllBreeds(brdDir)
	if err != nil {
		return nil, err
	}
	ret.Breeds = breeds
	return ret, nil
}

func saveBreed(brdDir string, breeds *BreedList) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(405)
			return
		}
		var breed game.Breed
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&breed)
		if err != nil {
			w.WriteHeader(400)
			writeAsJson(w, err.Error())
			return
		}

		breeds.Lock()
		defer breeds.Unlock()

		fname := fmt.Sprintf("%d%s", breed.Id, dataio.BreedFileExt)
		fpath := path.Join(brdDir, fname)
		fw, err := os.Create(fpath)
		if err != nil {
			w.WriteHeader(500)
			writeAsJson(w, err.Error())
			return
		}
		defer fw.Close()
		err = dataio.WriteBreed(fw, &breed)
		if err != nil {
			w.WriteHeader(500)
			writeAsJson(w, err.Error())
			return
		}
		idx := int(breed.Id)
		if idx >= len(breeds.Breeds) {
			newBreeds := make([]*game.Breed, idx+1)
			copy(newBreeds, breeds.Breeds)
			breeds.Breeds = newBreeds
		}
		breeds.Breeds[idx] = &breed
		writeAsJson(w, true)
	})
}

func serveBreeds(breeds *BreedList) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		breeds.RLock()
		defer breeds.RUnlock()
		writeAsJson(w, breeds.Breeds)
	})
}
