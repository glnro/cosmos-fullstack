# Cosmos Full Stack Journey

## Legend
**This tutorial describes the end to end journey and resulting state transitions of user interacting with the application.**


##### Interactions are documented at three levels of the stack:

**⚛️ APP** - The instance of the application that the user interacts with

**🚀 SDK** - The instance of the application that the Consensus application is communicating with.

**☄️ COMET** - The validator (single instance) that the State Machine is communicating with.

For the purpose of this demonstration, the application used will be a minimal SDK application [Neutrino](https://github.com/fatal-fruit/neutrino)

See the following SDK and Comet forks to 

<hr>

## Setup

1. Clone the application respository
```shell
git clone https://github.com/fatal-fruit/neutrino.git && cd neutrino
```

2. Start the localnet
```shell
make start-localnet
```

## User Journey

### Bank Send

Once the localnet is running, make the following transaction against the application in another session:
```shell
./build/neutrinod tx bank send  $(./build/neutrinod  keys show alice -a --home ~/.neutrinod-liveness --keyring-backend test)  $(./build/neutrinod  keys show bob -a --home ~/.neutrinod-liveness --keyring-backend test) 1000uneutrino --home ~/.neutrinod-liveness --keyring-backend test
```

After the transaction has been included in the block, stop the node and export state to verify the relevant transitions in the bank module.
There should be a `1000uneutrino` change in the bank balances of `Alice` and `Bob`.
```shell
./build/neutrinod  export --home ~/.neutrinod-liveness 
```

#### Sequence

**⚛️ Client Application**

`/x/bank/client/tx.go`
1.  Receive CLI Bank Send 
2. Format message & broadcast to consensus app

**☄️ Consensus Application**

`rpc/core/mempool.go`
3. Receive broadcasted Tx -> call mempool.CheckTx

`mempool/v0/clist_mempool.go`
4. Call CheckTx on ABCI application

**🚀 State Machine Application**

`baseapp/abci.go`
5. Run CheckTx on `MsgSend`
6. Retrieve cached context and multistore
7. Decode Tx
8. Validate Basic
9. NOOP (Insert into app side mempool)
10. Return Success response to Consensu App

**☄️ Consensus Application**

`mempool/v0/clist_mempool.go`
11. Add Tx to Mempool

`/consensus/state.go`
12. Prepare Block Proposal

`/state/execution.go`
13. Call BeginBlock on App
14. Call DeliverTx on App

**🚀 State Machine Application**

`baseapp/abci.go`
15. Run DeliverTx

`x/bank/keeper/msg_server.go`
16. Validate Basic
17. Run Send

` x/bank/keeper/send.go`
18. Subtract coins from sender
19. Add coins to recipient
20. Create account if none exists
21. Emit Transfer Event