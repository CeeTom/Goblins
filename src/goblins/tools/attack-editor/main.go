package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"goblins/game/combat"
	"goblins/game/dataio"
	"net/http"
	"os"
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
	http.Handle("/", http.FileServer(http.Dir(*webDir)))

	err = http.ListenAndServe(*host, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "server stopped: %s\n", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func serveAttacks(attacksP *[]*combat.Attack) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bytes, err := json.Marshal(*attacksP)
		if err != nil {
			errBytes, err := json.Marshal(err.Error())
			if err != nil {
				panic(err.Error())
			}
			w.WriteHeader(500)
			w.Write([]byte(`{"error": `))
			w.Write(errBytes)
			w.Write([]byte{'}'})
		}
		w.Write(bytes)
	})
}
