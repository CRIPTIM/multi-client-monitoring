// Copyright 2017 Maarten H. Everts and Tim R. van de Kamp.
// All rights reserved.
// Use of this source code is governed by the MIT license that can be
// found in the LICENSE file.

// An example program for running the scheme.
package main

import (
	"crypmonsys"
	"fmt"
	"github.com/Nik-U/pbc"
	"log"
)

func main() {
	fmt.Println("Cryptographic Monitoring System - sample code")
	sp := crypmonsys.NewSystemParameters(pbc.GenerateF(160).NewPairing())
	setupKey := crypmonsys.NewSetupKey(sp)
	rulegenerator, agents := setupKey.GenerateKeys(3, 8)

	identifier := "identifier"

	ruletoken, err := rulegenerator.NewToken([]int32{16, -1, 12})
	if err != nil {
		log.Fatal("Error creating token: ", err)
	}

	ciphertextsMatch := make([]*crypmonsys.Ciphertext, len(agents))
	ciphertextsMatch[0] = agents[0].NewCiphertext(identifier, 16)
	ciphertextsMatch[1] = agents[1].NewCiphertext(identifier, 42)
	ciphertextsMatch[2] = agents[2].NewCiphertext(identifier, 12)

	alarmsystem := crypmonsys.NewAlarmSystem(sp, ruletoken, identifier)

	if alarmsystem.Test(ciphertextsMatch) {
		fmt.Println("Alarm was raised, as expected.")
	} else {
		log.Fatal("No alarm was raised, whereas an alarm should have been raised.")
	}
}
