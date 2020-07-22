#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

for i in bsondump mongostat mongofiles mongoexport mongoimport mongorestore mongodump mongotop; do
        echo "Building ${i}..."
        go build -o "bin/$i"  "$i/main/$i.go"
done
