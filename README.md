# hlf-database-app

Skeleton for building a Hyperledger Fabric database app

## Requirements

### MacOS

```shell
brew install jq
```

### debian/ubuntu

```shell
sudo apt install jq
```

## QUICKSTART

1. Follow [hlf-service-network QUICKSTART instructions](https://github.com/Blockdaemon/hlf-service-network/blob/master/README.md#quickstart)
2. Configure environment

   ```shell
   make config.env
   ```

   Edit config.env to taste (see below for Blockdaemon CA-server setup)

3. Build binary and make credentials:

   ```shell
   make
   ```

3. Create/update/join channel, then install/instantiate chaincode (only have to do this once):

   ```shell
   ./app init
   ```

## Some things you can do

### Start up a webapp

```shell
./app webapp
```

It should be visible here: [http://localhost:3001](http://localhost:3001/)

### Set `hello` to a value, and retrieve it

```shell
./app set hello world
./app get hello
```

### Store a file and retrieve it

```shell
./app store <key> <sample-file>
./app fetch <key> out
diff <sample-file> out
```

### Clear out old chaincode

```shell
pushd $GOPATH/src/github.com/Blockdaemon/hlf-service-network
docker-compose down
popd
make clean-cc
```

### Talk to real external hosts and not a local docker instance

Set `DISABLE_MATCHERS=_` in `config.env`, rerun make

### Talk to a Blockdaemon CA server instead of using `cryptogen`

```shell
cp examples/config-ca.env config.env
```

Edit `config.env` to taste, then add CA server admin creds to `ca-client/local.env`:

```shell
CA_USER=<adminuser>
CA_PASS=<adminpass>
```

Make it a bit more secure, then run make:

```shell
chmod og-rw ca-client/local.env
make
```

## Bugs

* You may need to check out the repo in your `GOPATH` or chaincode installation may not work.
* If you change `config.env`, and you are using `ca-client`, you may need to `rm -rf ca-client/crypto-config`
* `.app get/set` may hang waiting for the transaction ack event. For some reason they go missing sometimes. CTRL-C to abort, the transaction usually went through

## References

Based on the [chainhero.io](https://chainhero.io) [tutorial](https://chainhero.io/2018/03/tutorial-build-blockchain-app-2/).

Requires a [service network](https://github.com/Blockdaemon/hlf-service-network) or a [Blockdaemon](https://blockdaemon.com/) hosted HLF network
