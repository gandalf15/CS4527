#!/bin/bash

ENDING='", "marcel", "1", "2"]}'

for (( i = 0; i < 100; ++i )); do
    for (( j = 0; j < 10; ++j )); do
        PAYLOAD='{"Args":["createDataEntryAd", "throughput", "test throughput data", "???", "Celsius", "'
        PAYLOAD=$PAYLOAD$i$j
        PAYLOAD=$PAYLOAD$ENDING
        2>/dev/null 1>&2 peer chaincode invoke --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/zak.codes/orderers/orderer.zak.codes/msp/tlscacerts/tlsca.zak.codes-cert.pem -n chaincode_ad -c "${PAYLOAD}" -C channel2 &
    done
    wait
done
