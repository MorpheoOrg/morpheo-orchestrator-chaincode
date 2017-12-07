# Morpheo Orchestrator Chaincode


This is the Orchestrator of the Morpheo platform with the blockchain. We use the private and permissioned solution called [Hyperledger Fabric](https://hyperledger-fabric.readthedocs.io/en/release/).  
Morpheo chaincode corresponds to the set of smart contracts, which are used to orchestrate operations on the [Morpheo platform](http://morpheo.co/). 
It is the translation of [Morpheo Orchestrator](https://github.com/MorpheoOrg/morpheo-orchestrator) with a blockchain solution.

**Licence:** CECILL 2.1 (compatible with GNU GPL)


## How to interact with the orchestrator

Use the [Morpheo-Fabric-Bootstrap](https://github.com/MorpheoOrg/morpheo-fabric-bootstrap) to create a network to interact with the Orchestrator.  
Once the network is up, the chaincode is installed and instantiated, you can go inside the docker cli to interact with the Orchestrator. Below some interaction examples, do not forget to set the correct environment variable:  
```
peer chaincode query -n mycc -c '{"Args":["queryItem", "algo_1"]}' -C $CHANNEL_NAME
peer chaincode query -n mycc -c '{"Args":["queryItems", "algo"]}' -C $CHANNEL_NAME
peer chaincode invoke -o orderer.morpheo.co:7050 --tls true --cafile $ORDERER_CA -n mycc -c '{"Args":["registerProblem", "dda81bfc-b5f4-5ba2-b81a-b464248f02d2", "2", "0pa81bfc-b5f4-5ba2-b81a-b464248f02a1, 0kk81bfc-b5f4-5ba2-b81a-b464248f02e3"]}' -C $CHANNEL_NAME
peer chaincode invoke -o orderer.morpheo.co:7050 --tls true --cafile $ORDERER_CA -n mycc -c '{"Args":["registerItem", "algo", "0pa81baa-b5f4-5ba2-b81a-b464248f02d2", "problem_1"]}' -C $CHANNEL_NAME
peer chaincode invoke -o orderer.morpheo.co:7050 --tls true --cafile $ORDERER_CA -n mycc -c '{"Args":["registerItem", "data", "9pa81bfc-b5f8-5ba2-b81a-b464248f02d2", "problem_1"]}' -C $CHANNEL_NAME
peer chaincode query -n mycc -c '{"Args":["queryProblemItems", "data", "problem_1"]}' -C $CHANNEL_NAME
peer chaincode query -n mycc -c '{"Args":["queryStatusLearnuplet", "todo"]}' -C $CHANNEL_NAME
// replace algo_0 with correct key
peer chaincode query -n mycc -c '{"Args":["queryAlgoLearnuplet", "algo_0"]}' -C $CHANNEL_NAME
// replace learnuplet_0 with correct key
peer chaincode invoke -o orderer.morpheo.co:7050 --tls true --cafile $ORDERER_CA -n mycc -c '{"Args":["setUpletWorker", "learnuplet_0", "Arbeiter_12"]}' -C $CHANNEL_NAME   
// replace learnuplet_0 with correct key
peer chaincode invoke -o orderer.morpheo.co:7050 --tls true --cafile $ORDERER_CA -n mycc -c '{"Args":["reportLearn", "learnuplet_0", "done", "0.82", "data_3 0.78, data_4 0.88", "data_2 0.80"]}' -C $CHANNEL_NAME
```


## Smart contracts documentation  


### Elements of the ledger

The ledger is a key value store. 
To be able to make complex queries, such as querying all algorithms related to a problem, we use [`CompositeKey`](https://godoc.org/github.com/hyperledger/fabric/core/chaincode/shim#ChaincodeStub.CreateCompositeKey).

We call an `ObjectType` a type of element of the ledger (similar to a table in a relational database).  

#### Data and Algo

Data and algo are 2 ObjectTypes, which both derive from an Item structure:  
```
type Item struct {
    ObjectType     string `json:"docType"`
    Problem        string `json:"problem"`
}
```
Keys: `data_<uuid>` and `algo_<uuid>`.   
Associated composite keys: `data~problem~key` and `algo~problem~key`.  


#### Problem  

A problem derives from the Problem structure:
```
type Problem struct {
    ObjectType       string   `json:"docType"`
    StorageAddress   string   `json:"storage_address"`
    SizeTrainDataset int      `json:"size_train_dataset"`
    TestDataset      []string `json:"test_dataset"`
}
```  
Keys: `problem_<uuid>`.

#### Learnuplet

A learnuplet derives from the Learnuplet structure:  
```
type Learnuplet struct {
    ObjectType       string             `json:"docType"`
    Problem          string             `json:"problem"`
    ProblemAddress   string             `json:"problem_storage_address"`
    Algo             string             `json:"algo"`
    ModelStart       string             `json:"model_start"`
    ModelEnd         string             `json:"model_end"`
    TrainData        []string           `json:"train_data"`
    TestData         []string           `json:"test_data"`
    Worker           string             `json:"worker"`
    Status           string             `json:"status"`
    Rank             int                `json:"rank"`
    Perf             float64            `json:"perf"`
    TrainPerf        map[string]float64 `json:"train_perf"`
    TestPerf         map[string]float64 `json:"test_perf"`
}
```
Keys: `learnuplet_<uuid>`.   
Associated composite key: `learnuplet~algo~key`.  


### + `queryItem`: to query a given item

Args:  
- `item_key`, such as `data_8fa81bfc-b5f4-4ba2-b81a-b464248f02d3`, `learnuplet_ca3a5a53-9684-429f-9896-4f7c94f9def0`  

```
peer chaincode query -n mycc -c '{"Args":["queryItem", "learnuplet_ca3a5a53-9684-429f-9896-4f7c94f9def0"]}' -C $CHANNEL_NAME
```

### + `queryItems`: to query items of a given object type

Args:  
- `object_type`, such as `data`, `learnuplet`  

```
peer chaincode query -n mycc -c '{"Args":["queryItems", "learnuplet"]}' -C $CHANNEL_NAME
```

### + `queryProblemItems`: to query data or algos related to a problem

Args:  
- `item_type`: `data` or `algo`  
- `problem_key`, such as `problem_8fa81bfc-b5f4-4ba2-b81a-b464248f02d3`  

```
peer chaincode query -n mycc -c '{"Args":["queryProblemItems", "data", "problem_8fa81bfc-b5f4-4ba2-b81a-b464248f02d3"]}' -C $CHANNEL_NAME
```

### + `registerItem`: to register an algo or data  

TODO: modify function to register several data at a time  

Args:  
- `item_type`: `data` or `algo`  
- `storage_address`, such as `https://storage.morpheo.io/algo/0pa81bfc-b5f4-5ba2-b81a-b464248f02d2`    
- `problem_key`, such as `problem_2`  

```
peer chaincode invoke -o orderer.morpheo.co:7050 --tls true --cafile $ORDERER_CA -n mycc -c '{"Args":["registerItem", "algo", "https://storage.morpheo.io/algo/0pa81bfc-b5f4-5ba2-b81a-b464248f02d2", "problem_1"]}' -C $CHANNEL_NAME
peer chaincode invoke -o orderer.morpheo.co:7050 --tls true --cafile $ORDERER_CA -n mycc -c '{"Args":["registerItem", "data", "https://storage.morpheo.io/data/9pa81bfc-b5f4-5ba2-b81a-b464248f02d2", "problem_1"]}' -C $CHANNEL_NAME
```

### + `registerProblem`: to register a new problem

Args:  
- `problem_storage_address`: address of the problem workflow on storage  
- `size_train_dataset`: number of train data per mini-batch   
- `test_data`: list of test data adresses on storage


```
peer chaincode invoke -o orderer.morpheo.co:7050 --tls true --cafile $ORDERER_CA -n mycc -c '{"Args":["registerProblem", "https://storage.morpheo.io/problem/dda81bfc-b5f4-5ba2-b81a-b464248f02d2", "2", "https://storage.morpheo.io/data/0pa81bfc-b5f4-5ba2-b81a-b464248f02a1, https://storage.morpheo.io/data/0pa81bfc-b5f4-5ba2-b81a-b464248f02e3"]}' -C $CHANNEL_NAME
```


### + `queryStatusLearnuplet`: to query all learnuplets with a given status

Args:  
- `status`: `todo`, `pending`, `failed`, or `done`

```
peer chaincode query -n mycc -c '{"Args":["queryStatusLearnuplet", "todo"]}' -C $CHANNEL_NAME
```

### + `queryAlgoLearnuplet`: to query all learnuplets associated to a given algo

Args:  
- `algo_key`: algo key of the algo of interest 

```
peer chaincode query -n mycc -c '{"Args":["queryAlgoLearnuplet", "algo_f50844e0-90e7-4fb8-a2aa-3d7e49204584"]}' -C $CHANNEL_NAME
```

### + `setUpletWorker`: to set the worker and change the status of a learnuplet  

Args:  
- `learnuplet_key`, such as `learnuplet_f50844e0-90e7-4fb8-a2aa-3d7e49204584`  
- `worker`: worker identifier... to be defined

```
peer chaincode invoke -o orderer.morpheo.co:7050 --tls true --cafile $ORDERER_CA -n mycc -c '{"Args":["setUpletWorker", "learnuplet_f50844e0-90e7-4fb8-a2aa-3d7e49204584", "Arbeiter_12"]}' -C $CHANNEL_NAME
```

### + `reportLearn`: to report the output of a learning task

Args:  
- `learnuplet_key`, such as `learnuplet_f50844e0-90e7-4fb8-a2aa-3d7e49204584`  
- `status`: `done` or `failed`  
- `perf`: performance of the model (performance on test data), such as `0.99`        
- `train_perf`: performances on each train data, such as `data_12 0.89, data_22 0.92, data_34 0.88, data_44 0.96`  
- `test_perf`: performances on each test data, such as `data_2 0.82, data_4 0.94, data_6 0.88`  

```
peer chaincode invoke -o orderer.morpheo.co:7050 --tls true --cafile $ORDERER_CA -n mycc -c '{"Args":["reportLearn", "learnuplet_f50844e0-90e7-4fb8-a2aa-3d7e49204584", "done", "0.82", "data_3 0.78, data_4 0.88", "data_2 0.80"]}' -C $CHANNEL_NAME
```


## TODO

- [ ] register several data at a time  


