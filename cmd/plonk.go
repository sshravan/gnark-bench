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

	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/consensys/gnark/test/unsafekzg"
	"github.com/pkg/profile"
	"github.com/spf13/cobra"
)

// plonkCmd represents the plonk command
var plonkCmd = &cobra.Command{
	Use:   "plonk",
	Short: "runs benchmarks and profiles using PlonK proof system",
	Run:   runPlonk,
}

func plonkProofSize(proof plonk.Proof) int {
	var buf bytes.Buffer

	// Write the proof to the buffer using WriteRawTo
	bytesWritten, err := proof.WriteTo(&buf)
	if err != nil {
		fmt.Println("Error writing proof:", err)
		return 0
	}

	return int(bytesWritten)
}

func runPlonk(cmd *cobra.Command, args []string) {

	backend := "plonk"
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
		var ccs constraint.ConstraintSystem
		for i := 0; i < *fCount; i++ {
			ccs, err = frontend.Compile(curveID.ScalarField(), scs.NewBuilder, c.Circuit(*fCircuitSize), frontend.WithCapacity(*fCircuitSize))
		}
		stopProfile()
		assertNoError(err)
		writeResults(backend, w, took, ccs, -1)
		return
	}

	ccs, err := frontend.Compile(curveID.ScalarField(), scs.NewBuilder, c.Circuit(*fCircuitSize), frontend.WithCapacity(*fCircuitSize))
	assertNoError(err)

	// create srs
	srs, srsLagrange, err := unsafekzg.NewSRS(ccs)
	assertNoError(err)

	if *fAlgo == "setup" {
		startProfile()
		var err error
		for i := 0; i < *fCount; i++ {
			_, _, err = plonk.Setup(ccs, srs, srsLagrange)
		}
		stopProfile()
		assertNoError(err)
		writeResults(backend, w, took, ccs, -1)
		return
	}

	witness := c.Witness(*fCircuitSize, curveID)
	pk, vk, err := plonk.Setup(ccs, srs, srsLagrange)
	assertNoError(err)

	if *fAlgo == "prove" {
		var proof plonk.Proof
		startProfile()
		for i := 0; i < *fCount; i++ {
			proof, err = plonk.Prove(ccs, pk, witness)
		}
		stopProfile()
		assertNoError(err)
		proofSize := plonkProofSize(proof)
		writeResults(backend, w, took, ccs, proofSize)
		return
	}

	if *fAlgo != "verify" {
		panic("algo at this stage should be verify")
	}

	proof, err := plonk.Prove(ccs, pk, witness)
	assertNoError(err)

	publicWitness, err := witness.Public()
	assertNoError(err)

	startProfile()
	for i := 0; i < *fCount; i++ {
		err = plonk.Verify(proof, vk, publicWitness)
	}
	stopProfile()
	assertNoError(err)
	proofSize := plonkProofSize(proof)
	writeResults(backend, w, took, ccs, proofSize)

}

func init() {
	rootCmd.AddCommand(plonkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// plonkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// plonkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
