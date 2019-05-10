# hlf-database-app
Skeleton for building a Hyperledger Fabric database app

# QUICKSTART

1. Follow [hlf-service-network QUICKSTART instructions](https://github.com/Blockdaemon/hlf-service-network/blob/master/README.md#quickstart)
2. Build the code
```
go get
make
```
3. Create/update/join channel, then install/instantiate chaincode (only have to do this once):
```
./app init
```

# Some things you can do
## Start up a webapp
```
./app webapp
```
It should be visible here: [http://localhost:3001](http://localhost:3001/)

## Set `hello` to a value, and retrieve it
```
./app set hello world
./app get hello
```

## Store a file and retrieve it
```
./app store <key> <sample-file>
./app fetch <key> out
diff <sample-file> out
```

## Clear out old chaincode
```
pushd $GOPATH/src/github.com/Blockdaemon/hlf-service-network
docker-compose down
popd
make clean-cc
```

## Talk to real external hosts and not a local docker instance
Set `DISABLE_MATCHERS=_` in `config.env`, rerun make

## Talk to a Blockdaemon CA server instead of using `cryptogen`
```
cp examples/config-ca.env config.env
```
Edit `config.env` to taste, then add CA server admin creds to `ca-client/local.env`:
```
CA_USER=<adminuser>
CA_PASS=<adminpass>
```
Then run make:
```
make
```

# References
Based on the [chainhero.io](https://chainhero.io) [tutorial](https://chainhero.io/2018/03/tutorial-build-blockchain-app-2/).

Requires a [service network](https://github.com/Blockdaemon/hlf-service-network) or a [Blockdaemon](https://blockdaemon.com/) hosted HLF network
