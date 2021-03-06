#!/bin/bash

PAYLOAD='{"Args":["sendTokensFast", "1", "2", "1", "false"]}'

for (( i = 0; i < 100; ++i ))
do
	1>/dev/null 2>&1 peer chaincode invoke --tls true --cafile \
	/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/zak.codes/orderers/orderer.zak.codes/msp/tlscacerts/tlsca.zak.codes-cert.pem \
	-n chaincode_tokens -c "${PAYLOAD}" -C channel3
done
wait
