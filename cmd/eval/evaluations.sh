#!/bin/bash
# Script to automate timing experiments of the implementation.
if (($# < 1)); then
	CURVE="d159"
else
	CURVE="$1"
fi

HEADER="agents mean       sd         variance"
echo "$HEADER" > performance_${CURVE}_setup.dat
echo "$HEADER" > performance_${CURVE}_encrypt.dat
echo "$HEADER" > performance_${CURVE}_gentoken.dat
echo "$HEADER" > performance_${CURVE}_test.dat

experiment() {
	echo "Running ./performance-evaluation --param "param/${CURVE}.param" --datOutput $@"
	RESULTS=$(./performance-evaluation --param "param/${CURVE}.param" --datOutput "$@")
	echo "$RESULTS" | grep 'Setup'   | cut -b12- >> performance_${CURVE}_setup.dat
	echo "$RESULTS" | grep 'Encrypt' | cut -b12- >> performance_${CURVE}_encrypt.dat
	echo "$RESULTS" | grep 'GenToken'| cut -b12- >> performance_${CURVE}_gentoken.dat
	echo "$RESULTS" | grep 'Test'    | cut -b12- >> performance_${CURVE}_test.dat
}

# This should run in at most 8 hours on an Intel Core i5 for the `f` curve.
# 16 minutes
experiment --setup --encrypt --gentoken --test --experiments 1000 --runs 10 --agents 1
# 32 minutes
experiment --setup --gentoken --test --experiments 1000 --runs 10 --agents  2
# 48 minutes
experiment --setup --gentoken --test --experiments 1000 --runs 10 --agents  3
# 32 minutes
experiment --setup --gentoken --test --experiments  500 --runs 10 --agents  4
# 40 minutes
experiment --setup --gentoken --test --experiments  500 --runs 10 --agents  5
# 80 minutes
experiment --setup --gentoken --test --experiments  500 --runs 10 --agents 10
# 24 minutes
experiment --setup --gentoken --test --experiments  100 --runs 10 --agents 15
# 32 minutes
experiment --setup --gentoken --test --experiments  100 --runs 10 --agents 20
# 40 minutes
experiment --setup --gentoken --test --experiments  100 --runs 10 --agents 25
# 80 minutes
experiment --setup --gentoken --test --experiments  100 --runs 10 --agents 50
# 80 minutes
experiment --setup --gentoken --test --experiments   50 --runs 10 --agents 100
