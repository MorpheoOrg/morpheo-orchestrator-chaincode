# Morpheo Orchestrator Chaincode

This is the Orchestrator of the Morpheo platform with the blockchain. We use the private and permissioned solution called [Hyperledger Fabric](https://hyperledger-fabric.readthedocs.io/en/release/).  
Morpheo chaincode corresponds to the set of smart contracts, which are used to orchestrate operations on the [Morpheo platform](http://morpheo.co/). 
It is the translation of [Morpheo Orchestrator](https://github.com/MorpheoOrg/morpheo-orchestrator) with a blockchain solution.

**Licence:** CECILL 2.1 (compatible with GNU GPL)


## How to interact with the orchestrator

Use the [Morpheo-Fabric-Bootstrap](https://github.com/MorpheoOrg/morpheo-fabric-bootstrap) to create a network to interact with the Orchestrator.  
Once the network is up, the chaincode is installed and instantiated, you can go inside the docker cli to interact with the Orchestrator. Below some interaction examples, do not forget to set the correct environment variable:  
```
peer chaincode query -n mycc -c '{"Args":["queryObject", "algo_1"]}' -C $CHANNEL_NAME
peer chaincode query -n mycc -c '{"Args":["queryObjects", "algo"]}' -C $CHANNEL_NAME
peer chaincode invoke -o orderer.morpheo.co:7050 --tls true --cafile $ORDERER_CA -n mycc -c '{"Args":["registerProblem", "dda81bfc-b5f4-5ba2-b81a-b464248f02d2", "2", "0pa81bfc-b5f4-5ba2-b81a-b464248f02a1, 0kk81bfc-b5f4-5ba2-b81a-b464248f02e3"]}' -C $CHANNEL_NAME
peer chaincode invoke -o orderer.morpheo.co:7050 --tls true --cafile $ORDERER_CA -n mycc -c '{"Args":["registerItem", "algo", "0pa81baa-b5f4-5ba2-b81a-b464248f02d2", "problem_1", "mytopalgo"]}' -C $CHANNEL_NAME
peer chaincode invoke -o orderer.morpheo.co:7050 --tls true --cafile $ORDERER_CA -n mycc -c '{"Args":["registerItem", "data", "9pa81bfc-b5f8-5ba2-b81a-b464248f02d2", "problem_1", "psg"]}' -C $CHANNEL_NAME
peer chaincode query -n mycc -c '{"Args":["queryProblemItems", "data", "problem_1"]}' -C $CHANNEL_NAME
peer chaincode query -n mycc -c '{"Args":["queryStatusLearnuplet", "todo"]}' -C $CHANNEL_NAME
// replace algo_0 with correct key
peer chaincode query -n mycc -c '{"Args":["queryAlgoLearnuplet", "algo_0"]}' -C $CHANNEL_NAME
// replace learnuplet_0 with correct key
peer chaincode invoke -o orderer.morpheo.co:7050 --tls true --cafile $ORDERER_CA -n mycc -c '{"Args":["setUpletWorker", "learnuplet_0", "Arbeiter_12"]}' -C $CHANNEL_NAME   
// replace learnuplet_0 with correct key
peer chaincode invoke -o orderer.morpheo.co:7050 --tls true --cafile $ORDERER_CA -n mycc -c '{"Args":["reportLearn", "learnuplet_0", "done", "0.82", "{\"data_3\": 0.78, \"data_4\": 0.88}", "{\"data_2\": 0.80}"]}' -C $CHANNEL_NAME
```

## Chaincode-docker-devmode  

You can use the `chaincode-docker-devmode` to more easily develop the chaincode, [as detailed here](./chaincode-docker-devmode/README.md)

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
    StorageAddress string `json:"storageAddress"`
    Problem        string `json:"problem"`
    Name           string `json:"name"`
}
```
**Keys**: `data_<uuid>` and `algo_<uuid>`.   
Associated composite keys: `data~problem~key` and `algo~problem~key`.  


#### Problem  

A problem derives from the Problem structure:
```
type Problem struct {
    ObjectType       string   `json:"docType"`
    StorageAddress   string   `json:"storageAddress"`
    SizeTrainDataset int      `json:"sizeTrainDataset"`
    TestData         []string `json:"testData"`
}
```  
**Keys**: `problem_<uuid>`.

#### Learnuplet

A learnuplet derives from the Learnuplet structure:  
```
type Learnuplet struct {
    ObjectType        string             `json:"docType"`
    Problem           map[string]string  `json:"problem"`      // {problemKey: problemStorageAddress}
    Algo              map[string]string  `json:"algo"`         // {algoKey: algoStorageAddress}
    ModelStartAddress string             `json:"modelStartAddress"`
    ModelEndAddress   string             `json:"modelEndAddress"`
    TrainData         map[string]string  `json:"trainData"`    // {data1Key: data1StorageAddress, ...}
    TestData          map[string]string  `json:"testData"`     // {data1Key: data1StorageAddress, ...}
    Worker            string             `json:"worker"`
    Status            string             `json:"status"`
    Rank              int                `json:"rank"`
    Perf              float64            `json:"perf"`
    TrainPerf         map[string]float64 `json:"trainPerf"`
    TestPerf          map[string]float64 `json:"testPerf"`
}
```
**Keys**: `learnuplet_<uuid>`.   
Associated composite key: `learnuplet~algo~key`.  


### Smart Contracts 

#### + `queryObject`: to query a given object

Args:  
- `objectKey`, such as `data_8fa81bfc-b5f4-4ba2-b81a-b464248f02d3`, `learnuplet_ca3a5a53-9684-429f-9896-4f7c94f9def0`  

```
peer chaincode query -n mycc -c '{"Args":["queryObject", "learnuplet_ca3a5a53-9684-429f-9896-4f7c94f9def0"]}' -C $CHANNEL_NAME
```

#### + `queryObjects`: to query objects of a given type

Args:  
- `objectType`, such as `data`, `learnuplet`  

```
peer chaincode query -n mycc -c '{"Args":["queryObjects", "learnuplet"]}' -C $CHANNEL_NAME
```

#### + `queryProblemItems`: to query data or algos related to a problem

Args:  
- `itemType`: `data` or `algo`  
- `problemKey`, such as `problem_8fa81bfc-b5f4-4ba2-b81a-b464248f02d3`  

```
peer chaincode query -n mycc -c '{"Args":["queryProblemItems", "data", "problem_8fa81bfc-b5f4-4ba2-b81a-b464248f02d3"]}' -C $CHANNEL_NAME
```

#### + `registerItem`: to register an algo or data  

TODO: modify function to register several data at a time  

Args:  
- `itemType`: `data` or `algo`  
- `storageAddress`, for now it corresponds to the uuid on Storage, such as `0pa81bfc-b5f4-5ba2-b81a-b464248f02d2`    
- `problemKey`, such as `problem_2`  
- `itemName`, such as `mysuperalgo`  

```
peer chaincode invoke -o orderer.morpheo.co:7050 --tls true --cafile $ORDERER_CA -n mycc -c '{"Args":["registerItem", 
"algo", "0pa81bfc-b5f4-5ba2-b81a-b464248f02d2", "problem_1", "topalgo"]}' -C $CHANNEL_NAME
peer chaincode invoke -o orderer.morpheo.co:7050 --tls true --cafile $ORDERER_CA -n mycc -c '{"Args":["registerItem", "data", "9pa81bfc-b5f4-5ba2-b81a-b464248f02d2", "problem_1", "psg"]}' -C $CHANNEL_NAME
```

#### + `registerProblem`: to register a new problem

Args:  
- `storageAddress`: address of the problem workflow on storage  
- `sizeTrainDataset`: number of train data per mini-batch   
- `testData`: list of test data adresses on storage


```
peer chaincode invoke -o orderer.morpheo.co:7050 --tls true --cafile $ORDERER_CA -n mycc -c '{"Args":["registerProblem", "dda81bfc-b5f4-5ba2-b81a-b464248f02d2", "2", "0pa81bfc-b5f4-5ba2-b81a-b464248f02a1, 0pa81bfc-b5f4-5ba2-b81a-b464248f02e3"]}' -C $CHANNEL_NAME
```


#### + `queryStatusLearnuplet`: to query all learnuplets with a given status

Args:  
- `status`: `todo`, `pending`, `failed`, or `done`

```
peer chaincode query -n mycc -c '{"Args":["queryStatusLearnuplet", "todo"]}' -C $CHANNEL_NAME
```

#### + `queryAlgoLearnuplet`: to query all learnuplets associated to a given algo

Args:  
- `algoKey`: algo key of the algo of interest 

```
peer chaincode query -n mycc -c '{"Args":["queryAlgoLearnuplet", "algo_f50844e0-90e7-4fb8-a2aa-3d7e49204584"]}' -C $CHANNEL_NAME
```

#### + `setUpletWorker`: to set the worker and change the status of a learnuplet  

Args:  
- `learnupletKey`, such as `learnuplet_f50844e0-90e7-4fb8-a2aa-3d7e49204584`  
- `worker`: worker identifier... to be defined

```
peer chaincode invoke -o orderer.morpheo.co:7050 --tls true --cafile $ORDERER_CA -n mycc -c '{"Args":["setUpletWorker", "learnuplet_f50844e0-90e7-4fb8-a2aa-3d7e49204584", "Arbeiter_12"]}' -C $CHANNEL_NAME
```

#### + `reportLearn`: to report the output of a learning task

Args:  
- `learnupletKey`, such as `learnuplet_f50844e0-90e7-4fb8-a2aa-3d7e49204584`  
- `status`: `done` or `failed`  
- `perf`: performance of the model (performance on test data), such as `0.99`        
- `trainPerf`: performances on each train data, such as `{\"data_12\": 0.89, \"data_22\": 0.92, \"data_34\": 0.88, \"data_44\": 0.96}`  
- `testPerf`: performances on each test data, such as `{\"data_2\": 0.82, \"data_4\": 0.94, \"data_6\": 0.88}`  
```
peer chaincode invoke -o orderer.morpheo.co:7050 --tls true --cafile $ORDERER_CA -n mycc -c '{"Args":["reportLearn", "learnuplet_f50844e0-90e7-4fb8-a2aa-3d7e49204584", "done", "0.82", "{\"data_3\": 0.78, \"data_4\": 0.88}", "{\"data_2\": 0.80}"]}' -C $CHANNEL_NAME
```



## TODO

- [ ] register several data at a time  


