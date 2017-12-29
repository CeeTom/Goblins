package main

import (
	"encoding/json"
	"fmt"
	"goblins/game/combat"
	"goblins/game/dataio"
	"net/http"
	"os"
	"path"
	"sync"
)

type AttackList struct {
	sync.RWMutex
	Attacks []*combat.Attack
}

func LoadAttackList(atkDir string) (*AttackList, error) {
	ret := new(AttackList)
	attacks, err := dataio.ReadAllAttacks(atkDir)
	if err != nil {
		return nil, err
	}
	ret.Attacks = attacks
	return ret, nil
}

func saveAttack(atkDir string, attacks *AttackList) http.Handler {
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

		attacks.Lock()
		defer attacks.Unlock()

		fname := fmt.Sprintf("%d.atk", attack.Id)
		fpath := path.Join(atkDir, fname)
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
		if idx >= len(attacks.Attacks) {
			newAttacks := make([]*combat.Attack, idx+1)
			copy(newAttacks, attacks.Attacks)
			attacks.Attacks = newAttacks
		}
		attacks.Attacks[idx] = &attack
		writeAsJson(w, true)
	})
}

func serveAttacks(attacks *AttackList) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attacks.RLock()
		defer attacks.RUnlock()
		writeAsJson(w, attacks.Attacks)
	})
}
