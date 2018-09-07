# hlf-database-app
Skeleton for building a Hyperledger Fabric database app

# QUICKSTART

* Follow [hlf-service-network QUICKSTART instructions](https://github.com/Blockdaemon/hlf-service-network/blob/master/README.md#quickstart)
* Build the code
```
go get
make
```

* Install the chaincode (only have to do this once)
```
./app init
```

## Some things you can do
### Start up a webapp
```
./app webapp
```
It should be visible here: [http://localhost:3001](http://localhost:3001/)

### Set `hello` to a value, and retrieve it
```
./app set hello world
./app get hello
```

### Store a file and retrieve it
```
./app store <key> <sample-file>
./app fetch <key> out
diff <sample-file> out
```

# To clear out old chaincode
```
pushd $GOPATH/src/github.com/Blockdaemon/hlf-service-network
docker-compose down
popd
. tools/rm-cc
```

# References
Based on [this tutorial](https://chainhero.io/2018/03/tutorial-build-blockchain-app-2/)

Requires [a service network](https://github.com/Blockdaemon/hlf-service-network)
