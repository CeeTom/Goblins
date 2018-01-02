package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"goblins/game"
	"goblins/game/combat"
	"net/http"
	"os"
)

func main() {
	webDir := flag.String("www", ".", "directory for static html files")
	atkDir := flag.String("atk", ".", "directory from which to read attacks")
	brdDir := flag.String("brd", ".", "directory from which to read breeds")
	host := flag.String("host", "127.0.0.1:8001",
		"host:port to listen for ui connections")
	flag.Parse()

	attacks, err := LoadAttackList(*atkDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't read attacks from %s: %s\n", *atkDir,
			err.Error())
		os.Exit(1)
	}

	breeds, err := LoadBreedList(*brdDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't read breeds from %s: %s\n", *brdDir,
			err.Error())
		os.Exit(1)
	}

	http.Handle("/attacks", serveAttacks(attacks))
	http.Handle("/saveAttack", saveAttack(*atkDir, attacks))
	http.Handle("/breeds", serveBreeds(breeds))
	http.Handle("/saveBreed", saveBreed(*brdDir, breeds))
	http.Handle("/staticGameInfo", http.HandlerFunc(serveStaticGameInfo))
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

type StaticInfo struct {
	Stats, ScalingFuncs, DamageTypes, Statuses, Traits []IdNamePair
}

func popIdNamePair(tgt *IdNamePair, v game.EnumId) {
	tgt.Name = v.Name()
	tgt.Id = v.AsU64()
}

func serveStaticGameInfo(w http.ResponseWriter, r *http.Request) {
	var toWrite StaticInfo
	toWrite.Stats = make([]IdNamePair, len(game.AllStats))
	for i, v := range game.AllStats {
		popIdNamePair(&toWrite.Stats[i], v)
	}
	toWrite.ScalingFuncs = make([]IdNamePair, len(combat.AllScalingFuncs))
	for i, v := range combat.AllScalingFuncs {
		popIdNamePair(&toWrite.ScalingFuncs[i], v)
	}
	toWrite.DamageTypes = make([]IdNamePair, len(combat.AllDamageTypes))
	for i, v := range combat.AllDamageTypes {
		popIdNamePair(&toWrite.DamageTypes[i], v)
	}
	toWrite.Statuses = make([]IdNamePair, len(combat.AllStatuses))
	for i, v := range combat.AllStatuses {
		popIdNamePair(&toWrite.Statuses[i], v)
	}
	toWrite.Traits = make([]IdNamePair, len(game.AllTraits))
	for i, v := range game.AllTraits {
		popIdNamePair(&toWrite.Traits[i], v)
	}
	writeAsJson(w, &toWrite)
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
