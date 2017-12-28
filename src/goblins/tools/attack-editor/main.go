package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"goblins/game"
	"goblins/game/combat"
	"goblins/game/dataio"
	"net/http"
	"os"
	"path"
)

func main() {
	webDir := flag.String("www", ".", "directory for static html files")
	tgtDir := flag.String("dir", ".", "directory from which to read attacks")
	host := flag.String("host", "127.0.0.1:8001",
		"host:port to listen for ui connections")
	flag.Parse()

	attacks, err := dataio.ReadAllAttacks(*tgtDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't read attacks from %s: %s\n", *tgtDir,
			err.Error())
		os.Exit(1)
	}
	attacksP := &attacks

	http.Handle("/attacks", serveAttacks(attacksP))
	http.Handle("/saveAttack", saveAttack(*tgtDir, attacksP))
	http.Handle("/stats", http.HandlerFunc(serveStats))
	http.Handle("/scalingFuncs", http.HandlerFunc(serveScalingFuncs))
	http.Handle("/damageTypes", http.HandlerFunc(serveDamageTypes))
	http.Handle("/statuses", http.HandlerFunc(serveStatuses))
	http.Handle("/", http.FileServer(http.Dir(*webDir)))

	err = http.ListenAndServe(*host, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "server stopped: %s\n", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

type IdNamePair struct {
	Id   uint64
	Name string
}

func popIdNamePair(tgt *IdNamePair, v game.EnumId) {
	tgt.Name = v.Name()
	tgt.Id = v.AsU64()
}

func saveAttack(tgtDir string, attacksP *[]*combat.Attack) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(405)
			return
		}
		var attack combat.Attack
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&attack)
		if err != nil {
			w.WriteHeader(400)
			writeAsJson(w, err.Error())
			return
		}
		fname := fmt.Sprintf("%d.atk", attack.Id)
		fpath := path.Join(tgtDir, fname)
		fw, err := os.Create(fpath)
		if err != nil {
			w.WriteHeader(500)
			writeAsJson(w, err.Error())
			return
		}
		defer fw.Close()
		err = dataio.WriteAttack(fw, &attack)
		if err != nil {
			w.WriteHeader(500)
			writeAsJson(w, err.Error())
			return
		}
		idx := int(attack.Id)
		if idx >= len(*attacksP) {
			newAttacks := make([]*combat.Attack, idx+1)
			copy(newAttacks, *attacksP)
			*attacksP = newAttacks
		}
		(*attacksP)[idx] = &attack
		writeAsJson(w, true)
	})
}

func serveAttacks(attacksP *[]*combat.Attack) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeAsJson(w, *attacksP)
	})
}

func serveStats(w http.ResponseWriter, r *http.Request) {
	arr := make([]IdNamePair, len(combat.AllStats))
	for i, v := range combat.AllStats {
		popIdNamePair(&arr[i], v)
	}
	writeAsJson(w, arr)
}

func serveScalingFuncs(w http.ResponseWriter, r *http.Request) {
	arr := make([]IdNamePair, len(combat.AllScalingFuncs))
	for i, v := range combat.AllScalingFuncs {
		popIdNamePair(&arr[i], v)
	}
	writeAsJson(w, arr)
}

func serveDamageTypes(w http.ResponseWriter, r *http.Request) {
	arr := make([]IdNamePair, len(combat.AllDamageTypes))
	for i, v := range combat.AllDamageTypes {
		popIdNamePair(&arr[i], v)
	}
	writeAsJson(w, arr)
}

func serveStatuses(w http.ResponseWriter, r *http.Request) {
	arr := make([]IdNamePair, len(combat.AllStatuses))
	for i, v := range combat.AllStatuses {
		popIdNamePair(&arr[i], v)
	}
	writeAsJson(w, arr)
}

func writeAsJson(w http.ResponseWriter, v interface{}) {
	bytes, err := json.Marshal(v)
	if err != nil {
		errBytes, err := json.Marshal(err.Error())
		if err != nil {
			panic(err.Error())
		}
		w.WriteHeader(500)
		w.Write([]byte(`{"error": `))
		w.Write(errBytes)
		w.Write([]byte{'}'})
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.Write(bytes)
}
