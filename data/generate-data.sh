#!env bash

# Read the top 10 from the file `cat ../data/data.txt  | sort -k 2 -n -r | head -n 10`

_lines="${1:-1000}"
_baseurl="http://api.tech.com/item"

rm -f data.txt
seq $_lines | xargs -P 10 -I {} sh -c "echo http://api.tech.com/item/{} \$(od -N 4 -t uL -An /dev/urandom | tr -d \" \") >> data.txt"
