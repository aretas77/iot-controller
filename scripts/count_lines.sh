#!/bin/sh

# This script will calculate total lines for Go and Vue files.
paho_commits=21
current_dir=$(pwd)
go_lines=0
lines_paho=0
vue_lines=0
js_lines=0
py_lines=0

calculate_paho()
{
    paho_dir="${current_dir}"/../paho.mqtt.golang

    echo "Go paho.mqtt.golang lines: "
    lines_paho="$(git -C $paho_dir diff --numstat HEAD~${paho_commits} | awk '{sum += $1} END {print sum}')"
    echo "   $lines_paho total"

    # Parse numbers
    lines_paho=$(echo "$lines_paho" | sed 's/[^0-9]*//g')
}

calculate_go()
{
    go_files="web cmd device types clients utils"

    echo "Go iotctl lines: "
    go_lines=$(find $go_files  -name '*.go' | xargs wc -l | tail -1)
    echo "$go_lines"

    # Parse numbers
    go_lines=$(echo "$go_lines" | sed 's/[^0-9]*//g')
}

calculate_vue()
{
    echo "Vue lines: "
    vue_lines=$(find web/ui/src -name '*.vue' | xargs wc -l | tail -1)
    echo " $vue_lines"

    echo "Js lines: "
    js_lines=$(find web/ui/src -name '*.js' | xargs wc -l | tail -1)
    echo " $js_lines"

    # Parse numbers
    vue_lines=$(echo "$vue_lines" | sed 's/[^0-9]*//g')
    js_lines=$(echo "$js_lines" | sed 's/[^0-9]*//g')
}

calculate_py()
{
    py_files="../iot-hades/ ../paho.mqtt.golang/interpreter.py"

    echo "Python lines: "
    py_lines=$(find $py_files -name '*.py' | xargs wc -l | tail -1)
    echo " $py_lines"

    # Parse numbers
    py_lines=$(echo "$py_lines" | sed 's/[^0-9]*//g')
}


calculate_go
calculate_paho
calculate_vue
calculate_py

echo "Total lines: $((vue_lines + js_lines + go_lines + py_lines + lines_paho))"
