package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/consensys/gnark/frontend"
)

func writeResults(backend string, w *csv.Writer, took time.Duration, ccs frontend.CompiledConstraintSystem, proofSize int) {
	if err := w.Write(benchData{}.headers()); err != nil {
		fmt.Println("error: ", err.Error())
		os.Exit(-1)
	}

	// check memory usage, max ram requested from OS
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	internal, secret, public := ccs.GetNbVariables()
	bData := benchData{
		Backend:             backend,
		Curve:               curveID.String(),
		Algorithm:           *fAlgo,
		NbCoefficients:      ccs.GetNbCoefficients(),
		NbConstraints:       ccs.GetNbConstraints(),
		NbInternalVariables: internal,
		NbSecretVariables:   secret,
		NbPublicVariables:   public,
		RunTime:             took.Milliseconds(),
		MaxRAM:              (m.Sys / 1024 / 1024),
		Throughput:          int(float64(ccs.GetNbConstraints()) / took.Seconds()),
		Count:               *fCount,
		ProofSize:           proofSize,
	}

	if err := w.Write(bData.values()); err != nil {
		panic(err)
	}
	w.Flush()
}
