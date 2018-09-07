# hlf-database-app
Skeleton for building a Hyperledger Fabric database app

# QUICKSTART

* Follow [hlf-service-network QUICKSTART instructions](https://github.com/Blockdaemon/hlf-service-network/blob/master/README.md)
* Build the code
```
make
```
* Install the chaincode (only have to do this once)
```
./app init
```
* Some things you can do
```
./app set hello world
./app get hello
```
```
./app store <sample-file>
./app fetch <sample-file> out
diff <sample-file> out
```

# To clear out old chaincode
```
pushd $GOPATH/src/github.com/Blockdaemon/hlf-service-network
docker-compose down
popd
. tools/rm-cc
```

Based on [this tutorial](https://chainhero.io/2018/03/tutorial-build-blockchain-app-2/)

Requires [a service network](https://github.com/Blockdaemon/hlf-service-network)
