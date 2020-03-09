#!/bin/sh

# This script is used to simulate IoT server using simple commands.

server=172.18.0.3
port=1883

# $1 - MAC
send_ack() {
    mac=$1
    ackFile="ack.json"
    
    if [ -z ${mac} ]; then
        echo "empty mac" && exit 1
    fi

    ackTopic="control/${mac}/global/ack"
    
cat > ${ackFile} << EOL
{
    "mac": "${mac}",
    "network": "global",
}
EOL

    mosquitto_pub -u mock -P test -h $server -p $port -t $ackTopic -f ack.json
}

case "$1" in
    send_ack)
        send_ack "$2"
        ;;
    *)
        echo "$1"
        exit 1
esac
