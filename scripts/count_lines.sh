#!/bin/sh

# This script will calculate total lines for Go and Vue files.
go_lines=0
vue_lines=0
js_lines=0

calculate_go() {
    echo "Go lines: "
    go_lines=$(find web cmd device types -name '*.go' | xargs wc -l | tail -1)
    echo "$go_lines"

    # Parse numbers
    go_lines=$(echo "$go_lines" | sed 's/[^0-9]*//g')
}

calculate_vue() {
    echo "Vue lines: "
    vue_lines=$(find web/ui/src -name '*.vue' | xargs wc -l | tail -1)
    echo "$vue_lines"

    echo "Js lines: "
    js_lines=$(find web/ui/src -name '*.js' | xargs wc -l | tail -1)
    echo "$js_lines"

    # Parse numbers
    vue_lines=$(echo "$vue_lines" | sed 's/[^0-9]*//g')
    js_lines=$(echo "$js_lines" | sed 's/[^0-9]*//g')
}

calculate_go
calculate_vue

echo "Total lines: $((vue_lines + js_lines + go_lines))"
