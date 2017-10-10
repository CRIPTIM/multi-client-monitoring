package crypmonsys

import (
	"github.com/Nik-U/pbc"
	"testing"
)

var (
	testSetupKey = NewSetupKey(NewSystemParameters(pbc.GenerateF(160).NewPairing()))
)

func TestBasis(t *testing.T) {
	rulegenerator, agents := testSetupKey.GenerateKeys(3, 8)

	identifier := "identifier"

	ruletoken, err := rulegenerator.NewToken([]int32{16, -1, 12})
	if err != nil {
		t.Fatal("Error creating token: ", err)
	}

	ciphertextsMatch := make([]*Ciphertext, len(agents))
	ciphertextsMatch[0] = agents[0].NewCiphertext(identifier, 16)
	ciphertextsMatch[1] = agents[1].NewCiphertext(identifier, 42)
	ciphertextsMatch[2] = agents[2].NewCiphertext(identifier, 12)

	alarmsystem := NewAlarmSystem(testSetupKey.sp, ruletoken, identifier)

	if alarmsystem.Test(ciphertextsMatch) {
		t.Log("Alarm was raised, as expected.")
	} else {
		t.Fatal("No alarm was raised, whereas an alarm should have been raised.")
	}

	ciphertextsNoMatch := make([]*Ciphertext, len(agents))
	ciphertextsNoMatch[0] = agents[0].NewCiphertext(identifier, 14)
	ciphertextsNoMatch[1] = agents[1].NewCiphertext(identifier, 42)
	ciphertextsNoMatch[2] = agents[2].NewCiphertext(identifier, 12)

	if alarmsystem.Test(ciphertextsNoMatch) {
		t.Fatal("Alarm was raised whereas it should not have.")
	} else {
		t.Log("No alarm was raised, as expected.")
	}
}

func TestIdentifier(t *testing.T) {
	rulegenerator, agents := testSetupKey.GenerateKeys(3, 8)

	identifier := "identifier"
	otherIdentifier := "some other identifier"

	ruletoken, err := rulegenerator.NewToken([]int32{16, -1, 12})
	if err != nil {
		t.Fatal("Error creating token: ", err)
	}

	ciphertextsNoMatch := make([]*Ciphertext, len(agents))
	ciphertextsNoMatch[0] = agents[0].NewCiphertext(identifier, 16)
	ciphertextsNoMatch[1] = agents[1].NewCiphertext(identifier, 42)
	ciphertextsNoMatch[2] = agents[2].NewCiphertext(otherIdentifier, 12)

	alarmsystem := NewAlarmSystem(testSetupKey.sp, ruletoken, identifier)

	if alarmsystem.Test(ciphertextsNoMatch) {
		t.Fatal("Alarm was raised whereas it should not have.")
	} else {
		t.Log("No alarm was raised, as expected.")
	}
}

func benchmarkEncryption(b *testing.B, messageSpaceBitSize int) {
	_, agents := testSetupKey.GenerateKeys(1, messageSpaceBitSize)
	agent := agents[0]
	// Make sure the message to encrypt is the worst message possible
	message := int32(2 ^ messageSpaceBitSize - 1)

	identifier := "identifier"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		agent.NewCiphertext(identifier, message)
	}
}

func BenchmarkEncryption8Bits(b *testing.B) {
	benchmarkEncryption(b, 8)
}

func BenchmarkEncryption16Bits(b *testing.B) {
	benchmarkEncryption(b, 16)
}

func BenchmarkEncryption32Bits(b *testing.B) {
	benchmarkEncryption(b, 32)
}

func benchmarkTest(b *testing.B, numAgents int) {
	rulegenerator, agents := testSetupKey.GenerateKeys(numAgents, 8)

	// rules will be all zeros
	rules := make([]int32, numAgents)
	ruletoken, err := rulegenerator.NewToken(rules)
	if err != nil {
		b.Fatal("Error creating rule: ", err)
	}

	identifier := "identifier"
	ciphertexts := make([]*Ciphertext, len(agents))
	for i := 0; i < numAgents; i++ {
		// All ciphertexts will be one, so there should never be a match
		ciphertexts[i] = agents[i].NewCiphertext(identifier, 1)
	}
	alarmsystem := NewAlarmSystem(testSetupKey.sp, ruletoken, identifier)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if alarmsystem.Test(ciphertexts) {
			b.Fatal("Alarm was raised whereas it should not have.")
		}
	}

}

func BenchmarkTest3Agents(b *testing.B) {
	benchmarkTest(b, 3)
}

func BenchmarkTest10Agents(b *testing.B) {
	benchmarkTest(b, 10)
}

func BenchmarkTest100Agents(b *testing.B) {
	benchmarkTest(b, 100)
}
