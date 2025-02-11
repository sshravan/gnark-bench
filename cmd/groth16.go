/*
Copyright © 2021 ConsenSys Software Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/pkg/profile"
	"github.com/spf13/cobra"
)

// groth16Cmd represents the groth16 command
var groth16Cmd = &cobra.Command{
	Use:   "groth16",
	Short: "runs benchmarks and profiles using Groth16 proof system",
	Run:   runGroth16,
}

func groth16ProofSize(proof groth16.Proof) int {
	var buf bytes.Buffer

	// Write the proof to the buffer using WriteRawTo
	bytesWritten, err := proof.WriteTo(&buf)
	if err != nil {
		fmt.Println("Error writing proof:", err)
		return 0
	}

	return int(bytesWritten)
}

func runGroth16(cmd *cobra.Command, args []string) {

	backend := "groth16"
	if err := parseFlags(); err != nil {
		fmt.Println("error: ", err.Error())
		cmd.Help()
		os.Exit(-1)
	}

	// write to stdout
	w := csv.NewWriter(os.Stdout)

	var (
		start time.Time
		took  time.Duration
		prof  interface{ Stop() }
	)

	startProfile := func() {
		start = time.Now()
		if p != nil {
			prof = profile.Start(p, profile.ProfilePath("."), profile.NoShutdownHook)
		}
	}

	stopProfile := func() {
		took = time.Since(start)
		if p != nil {
			prof.Stop()
		}
		took /= time.Duration(*fCount)
	}

	if *fAlgo == "compile" {
		startProfile()
		var err error
		var ccs frontend.CompiledConstraintSystem
		for i := 0; i < *fCount; i++ {
			ccs, err = frontend.Compile(curveID.ScalarField(), r1cs.NewBuilder, c.Circuit(*fCircuitSize), frontend.WithCapacity(*fCircuitSize))
		}
		stopProfile()
		assertNoError(err)
		writeResults(backend, w, took, ccs, -1)
		return
	}

	ccs, err := frontend.Compile(curveID.ScalarField(), r1cs.NewBuilder, c.Circuit(*fCircuitSize), frontend.WithCapacity(*fCircuitSize))
	assertNoError(err)

	if *fAlgo == "setup" {
		startProfile()
		var err error
		for i := 0; i < *fCount; i++ {
			_, _, err = groth16.Setup(ccs)
		}
		stopProfile()
		assertNoError(err)
		writeResults(backend, w, took, ccs, -1)
		return
	}

	witness := c.Witness(*fCircuitSize, curveID)

	if *fAlgo == "prove" {
		var proof groth16.Proof
		pk, err := groth16.DummySetup(ccs)
		assertNoError(err)

		startProfile()
		for i := 0; i < *fCount; i++ {
			proof, err = groth16.Prove(ccs, pk, witness)
		}
		stopProfile()
		assertNoError(err)
		proofSize := groth16ProofSize(proof)
		writeResults(backend, w, took, ccs, proofSize)
		return
	}

	if *fAlgo != "verify" {
		panic("algo at this stage should be verify")
	}
	pk, vk, err := groth16.Setup(ccs)
	assertNoError(err)

	proof, err := groth16.Prove(ccs, pk, witness)
	assertNoError(err)

	publicWitness, err := witness.Public()
	assertNoError(err)
	startProfile()
	for i := 0; i < *fCount; i++ {
		err = groth16.Verify(proof, vk, publicWitness)
	}
	stopProfile()
	assertNoError(err)
	proofSize := groth16ProofSize(proof)
	writeResults(backend, w, took, ccs, proofSize)

}

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.AddCommand(groth16Cmd)

}
