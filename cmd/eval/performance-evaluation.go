// Copyright 2017 Maarten H. Everts and Tim R. van de Kamp.
// All rights reserved.
// Use of this source code is governed by the MIT license that can be
// found in the LICENSE file.

// Program to measure the running time of the scheme's algorithms.
package main

import (
	"crypmonsys"
	"flag"
	"fmt"
	"github.com/Nik-U/pbc"
	"log"
	"math"
	"os"
	"time"
)

const (
	Setup = iota
	Encrypt
	GenToken
	Test
	numberOfAlgorithms
)

// benchmarkSetup executes experiments times the Setup algorithm of the
// scheme.
func benchmarkSetup(pairingGroup *crypmonsys.SystemParameters, experiments, numberOfAgents int) (duration time.Duration) {
	start := time.Now()
	for i := 0; i < experiments; i++ {
		setupKey := crypmonsys.NewSetupKey(pairingGroup)
		setupKey.GenerateKeys(numberOfAgents, 8)
	}

	return time.Since(start)
}

// benchmarkEncrypt executes experiments times the Encrypt algorithm of
// the scheme.
func benchmarkEncrypt(agent *crypmonsys.Agent, experiments int) (duration time.Duration) {
	start := time.Now()
	for i := 0; i < experiments; i++ {
		agent.NewCiphertext("identifier", 16)
	}

	return time.Since(start)
}

// benchmarkGenToken executes experiments times the GenToken algorithm
// of the scheme.
func benchmarkGenToken(rulegenerator *crypmonsys.RuleGenerator, rule []int32, experiments int) (duration time.Duration) {
	start := time.Now()
	for i := 0; i < experiments; i++ {
		rulegenerator.NewToken(rule)
	}

	return time.Since(start)
}

// benchmarkTest executes experiments times the Test algorithm of the
// scheme.
func benchmarkTest(pairingGroup *crypmonsys.SystemParameters, token *crypmonsys.RuleToken, ciphertexts []*crypmonsys.Ciphertext, experiments int) (duration time.Duration) {
	start := time.Now()
	for i := 0; i < experiments; i++ {
		alarmsystem := crypmonsys.NewAlarmSystem(pairingGroup, token, "identifier")
		alarmsystem.Test(ciphertexts)
	}

	return time.Since(start)
}

func main() {
	// Program options
	var param = flag.String("param", "", "file containing the curve parameters")
	var datOutput = flag.Bool("datOutput", false, "output dat files")
	var experiments = flag.Int("experiments", 100, "number of experiments to run")
	var numberOfAgents = flag.Int("agents", 5, "number of agents in the system")
	var runs = flag.Int("runs", 5, "number of runs of the experiments")
	var runSetup = flag.Bool("setup", false, "benchmark the Setup algorithm")
	var runEncrypt = flag.Bool("encrypt", false, "benchmark the Encrypt algorithm")
	var runGenToken = flag.Bool("gentoken", false, "benchmark the GenToken algorithm")
	var runTest = flag.Bool("test", false, "benchmark the Test algorithm")

	flag.Parse()

	if *param == "" {
		log.Fatal("You have to specify a parameters file.")
	}

	if !(*runSetup || *runEncrypt || *runGenToken || *runTest) {
		log.Fatal("Please select at least one algorithm to run.")
	}

	if !*datOutput {
		fmt.Println("Performance evaluation of the Cryptographic Monitoring System")
		fmt.Printf("Doing %d runs, with %d experiments each.\n", *runs, *experiments)
		fmt.Printf("Number of agents: %d\n", *numberOfAgents)
	}

	// Setup a basic system
	paramFile, err := os.Open(*param)
	if err != nil {
		log.Fatal(err)
	}
	curve, err := pbc.NewParams(paramFile)
	paramFile.Close()
	if err != nil {
		panic(err)
	}
	sp := crypmonsys.NewSystemParameters(curve.NewPairing())
	setupKey := crypmonsys.NewSetupKey(sp)
	rulegenerator, agents := setupKey.GenerateKeys(*numberOfAgents, 8)

	ciphertextsMatch := make([]*crypmonsys.Ciphertext, *numberOfAgents)
	for agent := 0; agent < *numberOfAgents; agent++ {
		ciphertextsMatch[agent] = agents[agent].NewCiphertext("identifier", 16)
	}

	rule := make([]int32, *numberOfAgents)
	for agent := 0; agent < *numberOfAgents; agent++ {
		rule[agent] = 16
	}
	token, _ := rulegenerator.NewToken(rule)

	// Individual runs
	run := make([][]time.Duration, numberOfAlgorithms)
	for i := range run {
		run[i] = make([]time.Duration, *runs)
	}
	// Sample mean
	sampleMean := make([]float64, numberOfAlgorithms)
	// Sample variance
	sampleVariance := make([]float64, numberOfAlgorithms)

	// Run the experiments
	fmt.Fprint(os.Stderr, "Progress: 0%")
	for i := 0; i < *runs; i++ {
		// Benchmark Setup
		if *runSetup {
			run[Setup][i] = benchmarkSetup(sp, *experiments, *numberOfAgents)
			sampleMean[Setup] += run[Setup][i].Seconds()
		}

		// Benchmark Encrypt
		if *runEncrypt {
			run[Encrypt][i] = benchmarkEncrypt(agents[0], *experiments)
			sampleMean[Encrypt] += run[Encrypt][i].Seconds()
		}

		// Benchmark GenToken
		if *runGenToken {
			run[GenToken][i] = benchmarkGenToken(rulegenerator, rule, *experiments)
			sampleMean[GenToken] += run[GenToken][i].Seconds()
		}

		// Benchmark Test
		if *runTest {
			run[Test][i] = benchmarkTest(sp, token, ciphertextsMatch, *experiments)
			sampleMean[Test] += run[Test][i].Seconds()
		}

		progress := float32(i+1) / float32(*runs)
		fmt.Fprintf(os.Stderr, "\rProgress: %d%%", uint(progress*100))
	}
	fmt.Fprintln(os.Stderr, "")

	// Compute the sample mean: $\bar{X} = \frac{1}{n} \sum_{i=1}^n X_i$
	for i := range sampleMean {
		sampleMean[i] /= float64(*experiments) * float64(*runs)
	}

	// Compute the sample variance: $\frac{1}{n-1} \sum_{i=i}^n (X_i - \bar{X}_i)^2$
	for i := 0; i < *runs; i++ {
		for j := range sampleVariance {
			difference := (run[j][i].Seconds() / float64(*experiments)) - sampleMean[j]
			sampleVariance[j] += (difference * difference)
		}
	}
	for i := range sampleVariance {
		sampleVariance[i] /= (float64(*runs) - 1)
	}

	// Output
	if *datOutput {
		if *runSetup {
			fmt.Printf("%-10s %-5d %10.7f %10.7f %10.7f\n", "Setup", *numberOfAgents, sampleMean[Setup], math.Sqrt(sampleVariance[Setup]), sampleVariance[Setup])
		}
		if *runEncrypt {
			fmt.Printf("%-10s %-5d %10.7f %10.7f %10.7f\n", "Encrypt", *numberOfAgents, sampleMean[Encrypt], math.Sqrt(sampleVariance[Encrypt]), sampleVariance[Encrypt])
		}
		if *runGenToken {
			fmt.Printf("%-10s %-5d %10.7f %10.7f %10.7f\n", "GenToken", *numberOfAgents, sampleMean[GenToken], math.Sqrt(sampleVariance[GenToken]), sampleVariance[GenToken])
		}
		if *runTest {
			fmt.Printf("%-10s %-5d %10.7f %10.7f %10.7f\n", "Test", *numberOfAgents, sampleMean[Test], math.Sqrt(sampleVariance[Test]), sampleVariance[Test])
		}
	} else {
		fmt.Printf("%-10s + %-10s + %-10s + %-10s\n", "Algorithm", "Mean (s)", "SD σ (s)", "Var. (s²)")
		if *runSetup {
			fmt.Printf("%-10s | %10.7f | %10.7f | %10.7f\n", "Setup", sampleMean[Setup], math.Sqrt(sampleVariance[Setup]), sampleVariance[Setup])
		}
		if *runEncrypt {
			fmt.Printf("%-10s | %10.7f | %10.7f | %10.7f\n", "Encrypt", sampleMean[Encrypt], math.Sqrt(sampleVariance[Encrypt]), sampleVariance[Encrypt])
		}
		if *runGenToken {
			fmt.Printf("%-10s | %10.7f | %10.7f | %10.7f\n", "GenToken", sampleMean[GenToken], math.Sqrt(sampleVariance[GenToken]), sampleVariance[GenToken])
		}
		if *runTest {
			fmt.Printf("%-10s | %10.7f | %10.7f | %10.7f\n", "Test", sampleMean[Test], math.Sqrt(sampleVariance[Test]), sampleVariance[Test])
		}
	}
}
