#!env bash

_lines="100000"
_baseurl="http://api.tech.com/item"

rm -f output.txt
# This should be a LONG value so 8 bytes
#seq $_lines | parallel --bar -q -j 10 sh -c "echo http://api.tech.com/item/{} \$(cat /dev/random | tr -dc '0-9' | head -c 8) >> output.txt"
seq $_lines | xargs -P 10 -I {} sh -c "echo http://api.tech.com/item/{} \$(cat /dev/random | tr -dc '0-9' | head -c 8) >> output.txt"
