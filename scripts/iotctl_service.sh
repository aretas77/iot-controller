#!/bin/sh

# This script is used to simulate IoT server using simple commands.

server=172.18.0.3
port=1883
tmp=./tmp

# $1 - MAC
send_ack() {
    mac=$1
    ackFile="ack.json"

    if [ -z "${mac}" ]; then
        echo "empty mac" && exit 1
    fi

    ackTopic="control/global/${mac}/ack"

cat > ${tmp}/${ackFile} << EOL
{ "mac": "${mac}", "network": "global" }
EOL

    mosquitto_pub -u mock -P test -h "$server" -p "$port" -t "$ackTopic" -f "$tmp"/ack.json
}

# send_stats is used to send some statistics for IoT Hades application.
# $1 - MAC
send_stats() {
	mac=$1
	statFile="stat.json"

	if [ -z "${mac}" ]; then
		echo "empty mac" && exit 1
	fi

	statTopic="node/global/${mac}/hades/statistics"

cat > ${tmp}/${statFile} << EOL
{ "mac": "${mac}" }
EOL

	mosquitto_pub -u mock -P test -h "$server" -p "$port" -t "$statTopic" -f "$tmp"/"$statFile"
}

# send_model_request is used to send a request to the Hades service for a request
# of a new tflite model.
# $1 - MAC
send_model_request() {
    mac=$1
    requestFile="req.json"

    if [ -z "${mac}" ]; then
        echo "empty mac" && exit 1
    fi

    reqTopic="hades/global/${mac}/model/request"

cat > ${tmp}/${requestFile} << EOL
{ "mac": ${mac}" }
EOL

	mosquitto_pub -u mock -P test -h "$server" -p "$port" -t "$reqTopic" -f "$tmp"/"$requestFile"
}

case "$1" in
    send_ack)
        send_ack "$2"
        ;;
	send_stats)
		send_stats "$2"
		;;
    send_model_request)
        send_model_request "$2"
        ;;
    *)
        if [ -z "$1" ]; then
            echo "Empty args, usage: [send_ack {MAC}] [send_stats {MAC}] [send_model_request {MAC}]"
        fi
        echo "Received argument: $1"
        exit 1
esac
