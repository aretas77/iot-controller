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

	# construct a sample request - a request could also be empty.
cat > ${tmp}/${requestFile} << EOL
{ "mac": ${mac}" }
EOL

	# possible to send a null message - '--null-mesage'
	mosquitto_pub -u mock -P test -h "$server" -p "$port" -t "$reqTopic" -f "$tmp"/"$requestFile"
}

# send_model_file is used to send a file (in our case - a model) to the device.
# A sent model should be received and processed by paho.mqtt library and NOT by device.
# The file will be sent 
# $1 - MAC
# $2 - filename
send_model_file() {
	mac=$1
	file=$2

	if [ -z "$mac" ] || [ -z "$file" ]; then
		printf "MAC or File is empty.\n"
		printf "MAC	 = %s.\n"
		printf "File = %s.\n"
	fi

	sendTopic="node/global/${mac}/hades/model/receive"

	mosquitto_pub -u mock -P test -h "$server" -p "$port" -t "$sendTopic"
}

# print_usage will print some information on how to use this script.
print_usage() {
cat <<EOF
Usage: $0 [options]

send_ack		| {MAC}				this will send ack to IoT Controller.
send_stats		| {MAC}				this will send a mock statistic entry to IoT controller.
send_model_request	| {MAC}				this will send a new request for model to the Hades.
send_model_file		| {MAC} {FILE}			this will send a file as bytes to paho.mqtt library.
EOF
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
	send_model_file)
		send_model_file "$2" "$3"
		;;
    *)
        echo "Received argument: $1"
        [ -z "$1" ] && print_usage
        exit 1
esac
