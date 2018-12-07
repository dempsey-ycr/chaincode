#/bin/bash

ARCH=$(uname -s)
CURRENT_DIR=${PWD}
DEPLOY_DIR=${CURRENT_DIR}/fabric-deploy2
CONFIG_DIR="fabric1.3_multimachine"

CHANNEL="mychannel"


function printHosts() {
  echo "print every service computer hosts:"
  printf "\n"
  #result=cat /etc/hosts | grep "kafka"
  #echo -e '\e[5;43m'$result'\e[m\n'
  echo "fabric1: 192.168.87.114                                                                   orderer1.example.com"
  echo "fabric2: 192.168.87.118  zookeeper1  kafka1    peer0.org1.example.com  couchdb0.org1.com  orderer2.example.com    ca.org1.example.com"
  echo "fabric3: 192.168.87.128  zookeeper2  kafka2    peer1.org1.example.com  couchdb1.org1.com  orderer3.example.com"
  echo "fabric4: 192.168.87.122  zookeeper3  kafka3    peer0.org2.example.com  couchdb0.org2.com                          ca.org2.example.com"
  echo "fabric5: 192.168.87.125              kafka4    peer1.org2.example.com  couchdb1.org2.com"
  printf "\n"
  echo "Info: The wire surface example will show the hyperledger fabric multi-machine environment: [4zookeeper + 3kafka] and [2Orgs + 4Peers + 3orderer]"
}

# Create a working directory --fabric-deploy
# Create distributed nodes   --orderer... peer...
function createDirs() {
  # find -path "./fabric1.3_multimachine/.git" -prune -o -type f -print|wc -l
  f=`find $CONFIG_DIR -type f|wc -l`
  if [ -e $DEPLOY_DIR -a $f -gt 0 ]
  then
      echo "Info: The working directory already exists"
  else
  mkdir -p $DEPLOY_DIR/{orderer1.example.com,orderer2.example.com,orderer3.example.com,peer0.org1.example.com,peer1.org1.example.com,peer0.org2.example.com,peer1.org2.example.com}
  fi
}

# Add some dependency files. 
# @1. orderer.yaml core.yaml configtx.yaml crypto-config.yaml
# @2. binary files: peer orderer configtxgen ...
function srcFiles() {
  cd ${DEPLOY_DIR}
  if [ -e $CONFIG_DIR ]
  then
      echo "Info: git pull fabric1.3_multimachine..."
      cd ${DEPLOY_DIR}
      git pull
      cd ..

  else
      echo "Info: fabric1.3_multimachine is not exit. git clone fabric1.3_multimachine..."
      git clone https://github.com/dempsey-ycr/fabric1.3_multimachine
      binaryFabric13
  fi
}

# binaryFabric13
# To download the binary files required by fabric, which include orderer, peer, configtxgen, configtxlator, cryptogen, and discover, most of them can also be obtained in another way: 
# cd $GOPATH/src/github.com/hyperledger/fabric
# make peer  
# make orderer
# make ...
function binaryFabric13() {
    if [ $ARCH = "Linux" ]
    then
        echo "Linux ..."
        if [ ! -e "hyperledger-fabric-linux-amd64-1.3.0.tar.gz" ]
        then
            #wget https://nexus.hyperledger.org/content/repositories/releases/org/hyperledger/fabric/hyperledger-fabric/linux-amd64-1.3.0/hyperledger-fabric-linux-amd64-1.3.0.tar.gz
            #tar -vzxf hyperledger-fabric-linux-amd64-1.3.0.tar.gz
            echo "download ...."
        fi
        tar -vzxf hyperledger-fabric-linux-amd64-1.3.0.tar.gz

    elif test $ARCH = "Darwin"
    then
        echo "Darwin ..."
        if [ ! -e "hyperledger-fabric-darwin-amd64-1.3.0.tar.gz" ]
        then
            wget https://nexus.hyperledger.org/content/repositories/releases/org/hyperledger/fabric/hyperledger-fabric/darwin-amd64-1.3.0/hyperledger-fabric-darwin-amd64-1.3.0.tar.gz
            tar -vzxf hyperledger-fabric-darwin-amd64-1.3.0.tar.gz
        fi

    else
        echo "error: Only Linux and Darwin systems are supported"
        exit -1
    fi
}



# 
#
# create certs 
function createCerts() {
    cd $DEPLOY_DIR
    if [ ! -e "certs" ];  then
        mkdir certs
    fi

    ./bin/cryptogen generate --config=crypto-config.yaml --output ./certs
    if [ $? -ne 0 ]; then
        exit -1
    fi
}

# 
#
# create genesisblock
function createGenesisblock() {
    cd $DEPLOY_DIR
    ./bin/configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./genesisblock
    if [ $? -ne 0 ]; then
        echo "create genesisblock failed..."
        exit -1
    fi
    echo "create genesisblock successed..."
}


# Move the configuration file to the peer and orderer directories
# Replace the domain and OrgxMSP
function moveReplace() {
    cp $CONFIG_DIR/{crypto-config.yaml,configtx.yaml} ./

    createCerts
    createGenesisblock

    #for peer in peer0.org1.example.com peer1.org1.example.com peer0.org2.example.com peer1.org2.example.com;
    for peer in $(ls -d peer*.example.com)
    do
        org=${peer#*.}
        cp ./bin/peer $peer
        cp -rf ./certs/peerOrganizations/${org}/peers/${peer}/* $peer
        cp $CONFIG_DIR/core.yaml $peer
        sed -i "s/peer0.org1.example.com/${peer}/g" $peer/core.yaml

        #if [ $peer = "peer0.org2.example.com" -o $peer = "peer1.org2.example.com" ]; then
        if [ $peer = "org2.example.com" ]; then
            sed -i 's/Org1MSP/Org2MSP/g' $peer/core.yaml
        fi

    done

    #for orderer in orderer1.example.com orderer2.example.com orderer3.example.com; 
    for orderer in $(ls -d orderer*.example.com)
    do
        cp ./bin/orderer $orderer
        cp -rf ./certs/ordererOrganizations/example.com/orderers/${orderer}/* $orderer
        cp $CONFIG_DIR/orderer.yaml $orderer
        cp ./genesisblock $orderer
        sed -i "s/orderer1.example.com/${orderer}/g" $orderer/orderer.yaml
        sed -i "s/Orderer1MSP/OrdererOrg/g" $orderer/orderer.yaml
    done
}

function clearWorkspace() {
    cd $DEPLOY_DIR
    if [ $? -eq 0 ];then
       echo "=========="
       # shopt -s extglob
       # rm -rf !\( *.tar.gz \)
       # shopt -u extglob
       find . -type f,d -not \( -name '*.tar.gz' -or -name 'fabric1.3_multimachine' \) -delete
    fi
}


# 
function createUsers() {
    cd $DEPLOY_DIR

    mkdir Admin@org1.example.com Admin@org2.example.com

    cp -rf  certs/peerOrganizations/org1.example.com/users/Admin@org1.example.com/*  Admin@org1.example.com/
    cp -rf  certs/peerOrganizations/org2.example.com/users/Admin@org2.example.com/*  Admin@org2.example.com/
    cp  peer0.org1.example.com/core.yaml  Admin@org1.example.com/
    cp  peer0.org2.example.com/core.yaml  Admin@org2.example.com/

    echo "#!/bin/bash
PATH=\`pwd\`/../bin:\$PATH

export FABRIC_CFG_PATH=\`pwd\`

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_CERT_FILE=./tls/client.crt
export CORE_PEER_TLS_KEY_FILE=./tls/client.key

export CORE_PEER_MSPCONFIGPATH=./msp
export CORE_PEER_ADDRESS=peer0.org1.example.com:7051
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_TLS_ROOTCERT_FILE=./tls/ca.crt
export CORE_PEER_ID=cli
export CORE_LOGGING_LEVEL=INFO

peer \$*" > Admin@org1.example.com/peer.sh
    chmod +x Admin@org1.example.com/peer.sh

    cp Admin@org1.example.com/peer.sh  Admin@org2.example.com
    sed -i 's/peer0.org1.example.com/peer0.org2.example.com/g' Admin@org2.example.com/peer.sh
    sed -i 's/Org1MSP/Org2MSP/g' Admin@org2.example.com/peer.sh

    cp -rf Admin@org1.example.com User1@org1.example.com
    rm -rf User1@org1.example.com/{msp,tls}
    cp -rf  certs/peerOrganizations/org1.example.com/users/User1@org1.example.com/* User1@org1.example.com/
}


# create channelFile
# set two anchorPeers in channelFile
function channelFile() {
    channelTx=${CHANNEL}.tx
    cd $DEPLOY_DIR

    # create channel
    ./bin/configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ${channelTx} -channelID $CHANNEL
    if [ $? -ne 0 ];then
            echo "create channel file[$channelTx] failed..."
            exit -1
    fi

    # create two anchor peers
    ./bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate Org1MSPanchors.tx -channelID $CHANNEL -asOrg Org1MSP
    ./bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate Org2MSPanchors.tx -channelID $CHANNEL -asOrg Org2MSP
    if [ $? -ne 0 ];then
            echo "set channel anchor node failed..."
            exit -1
    fi

    # move $channelTx to each users
    cp ./$channelTx Admin@org1.example.com
    cp ./$channelTx Admin@org2.example.com

    # The cli accesses orderer, who needs to authenticate the user, so the cli needs to have an orderer authorization certificate
    cp certs/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem  Admin\@org1.example.com/
    cp certs/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem  User1\@org1.example.com/
    cp certs/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem  Admin\@org2.example.com/

    echo "create channelFile successed..."
    return 0
}


#ssh root@192.168.87.114
#ssh root@192.168.87.118
#ssh root@192.168.87.128
#ssh root@192.168.87.122

if [ $1x = "clear"x ]; then
    echo "clear workspace..."
    clearWorkspace
    exit 0
fi

# Print the domain name mapping for each host
printHosts

# 
createDirs
srcFiles
moveReplace
createUsers
channelFile


## 
## g g v    shift G " + Y
##
