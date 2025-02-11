#!/bin/bash

# Check if the correct number of arguments is provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <size_exponent> <csv_file>"
    exit 1
fi

# Assign input arguments to variables
size_exponent=$1
csv_file=$2

# Calculate the size as 1 << size_exponent
size=$((1 << size_exponent))

# Function to run the benchmark and append results to the CSV file
run_benchmark() {
    local tags=$1
    local algorithm=$2
    local output

    # Build the binary with optional build tags
    if [ -z "$tags" ]; then
        go build
    else
        go build -tags="$tags"
    fi

    # Run the benchmark and capture output
    output=$( { date; time ./gnark-bench plonk --algo "$algorithm" --count 1 --curve bls12_381 --profile cpu --size "$size"; } 2>&1 )

    # Append the output to the CSV file
    echo "$output" | tee -a "$csv_file"
}

# Run the benchmark without additional build tags
echo "Running benchmark without amd64_adx..."
run_benchmark "" "prove"
run_benchmark "" "verify"

# Run the benchmark with amd64_adx build tags
echo "Running benchmark with amd64_adx..."
run_benchmark "amd64_adx" "prove"
run_benchmark "amd64_adx" "verify"

echo "Benchmark results appended to $csv_file"
