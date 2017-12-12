/*
Copyright Morpheo Org. 2017

 contact@morpheo.co

 This software is part of the Morpheo project, an open-source machine
 learning platform.
 This software is governed by the CeCILL license, compatible with the
 GNU GPL, under French law and abiding by the rules of distribution of
 free software. You can  use, modify and/ or redistribute the software
 under the terms of the CeCILL license as circulated by CEA, CNRS and
 INRIA at the following URL "http://www.cecill.info".

 As a counterpart to the access to the source code and  rights to copy,
 modify and redistribute granted by the license, users are provided only
 with a limited warranty  and the software's author,  the holder of the
 economic rights,  and the successive licensors  have only  limited
 liability.

 In this respect, the user's attention is drawn to the risks associated
 with loading,  using,  modifying and/or developing or reproducing the
 software by the user in light of its specific status of free software,
 that may mean  that it is complicated to manipulate,  and  that  also
 therefore means  that it is reserved for developers  and  experienced
 professionals having in-depth computer knowledge. Users are therefore
 encouraged to load and test the software's suitability as regards their
 requirements in conditions enabling the security of their systems and/or
 data to be ensured and,  more generally, to use and operate it in the
 same conditions as regards security.

 The fact that you are presently reading this means that you have had
 knowledge of the CeCILL license and that you accept its terms.
*/

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	"github.com/satori/go.uuid"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the Item structure, corresponding to algo and data.  Structure tags are used by encoding/json library
// ObjectType belongs to algo, data
type Item struct {
	ObjectType string `json:"docType"`
	Problem    string `json:"problem"`
}

// Define the Problem structure. ObjectType is problem
type Problem struct {
	ObjectType       string   `json:"docType"`
	StorageAddress   string   `json:"storage_address"`
	SizeTrainDataset int      `json:"size_train_dataset"`
	TestDataset      []string `json:"test_dataset"`
}

// Define Learnuplet structure. ObjectType is learnuplet
type Learnuplet struct {
	ObjectType     string             `json:"docType"`
	Problem        string             `json:"problem"`
	ProblemAddress string             `json:"problem_storage_address"`
	Algo           string             `json:"algo"`
	ModelStart     string             `json:"model_start"`
	ModelEnd       string             `json:"model_end"`
	TrainData      []string           `json:"train_data"`
	TestData       []string           `json:"test_data"`
	Worker         string             `json:"worker"`
	Status         string             `json:"status"`
	Rank           int                `json:"rank"`
	Perf           float64            `json:"perf"`
	TrainPerf      map[string]float64 `json:"train_perf"`
	TestPerf       map[string]float64 `json:"test_perf"`
}

/*
 * The Init method is called when the Smart Contract orchestrator is instantiated by the blockchain network
 * Note that chaincode upgrade also calls this function to reset
 * or to migrate data, so be careful to avoid a scenario where you
 * inadvertently clobber your ledger's data!
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	s.initLedger(APIstub)
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryItem" {
		return s.queryItem(APIstub, args)
	} else if function == "queryItems" {
		return s.queryItems(APIstub, args)
	} else if function == "queryProblemItems" {
		return s.queryProblemItems(APIstub, args)
	} else if function == "registerItem" {
		return s.registerItem(APIstub, args)
	} else if function == "registerProblem" {
		return s.registerProblem(APIstub, args)
	} else if function == "queryStatusLearnuplet" {
		return s.queryStatusLearnuplet(APIstub, args)
	} else if function == "queryAlgoLearnuplet" {
		return s.queryAlgoLearnuplet(APIstub, args)
	} else if function == "setUpletWorker" {
		return s.setUpletWorker(APIstub, args)
	} else if function == "reportLearn" {
		return s.reportLearn(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

// ============================================
// initLedger
// ============================================

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Populates the database
	// TODO modify this function to use other functions to register the elements
	fmt.Println("- start Populate Ledger")

	// Adding two problems
	problems := []Problem{
		Problem{ObjectType: "problem", StorageAddress: "2pa81bfc-b5f4-4ba2-b81a-b464248f02d1", SizeTrainDataset: 1, TestDataset: []string{"data_0"}},
		Problem{ObjectType: "problem", StorageAddress: "4za81bfc-b5f4-4ba2-b81a-b464248f02d1", SizeTrainDataset: 2, TestDataset: []string{"data_2"}},
	}
	for i := 0; i < len(problems); i++ {
		problemAsBytes, _ := json.Marshal(problems[i])
		problemKey := fmt.Sprintf("problem_%d", i)
		APIstub.PutState(problemKey, problemAsBytes)
		fmt.Println("-- added", problems[i])
	}

	// Adding algorithms
	// TODO merge the 2 subsection paragraphs
	algos := []Item{
		Item{ObjectType: "algo", Problem: "problem_1"},
		Item{ObjectType: "algo", Problem: "problem_1"},
	}

	for i := 0; i < len(algos); i++ {
		algoAsBytes, _ := json.Marshal(algos[i])
		algoKey := fmt.Sprintf("algo_%d", i)
		APIstub.PutState(algoKey, algoAsBytes)
		fmt.Println("-- added", algos[i])
		// composite key
		indexName := "algo~problem~key"
		algoProblemIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{algos[i].ObjectType, algos[i].Problem, algoKey})
		if err != nil {
			return shim.Error(err.Error())
		}
		value := []byte{0x00}
		APIstub.PutState(algoProblemIndexKey, value)
		// end of composite key
	}

	// Adding data
	datas := []Item{
		Item{ObjectType: "data", Problem: "problem_0"},
		Item{ObjectType: "data", Problem: "problem_0"},
		Item{ObjectType: "data", Problem: "problem_1"},
		Item{ObjectType: "data", Problem: "problem_1"},
		Item{ObjectType: "data", Problem: "problem_1"},
	}
	for i := 0; i < len(datas); i++ {
		dataAsBytes, _ := json.Marshal(datas[i])
		dataKey := fmt.Sprintf("data_%d", i)
		APIstub.PutState(dataKey, dataAsBytes)
		fmt.Println("-- added", datas[i])
		// composite key
		indexName := "data~problem~key"
		dataProblemIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{datas[i].ObjectType, datas[i].Problem, dataKey})
		if err != nil {
			return shim.Error(err.Error())
		}
		value := []byte{0x00}
		APIstub.PutState(dataProblemIndexKey, value)
		// end of composite key
	}

	fmt.Println("- end Populate Ledger")

	return shim.Success(nil)
}

// ======================================================================
// registerProblem - register a new problem, store it into chaincode state
// Should be callable by organizations only
// ======================================================================

func (s *SmartContract) registerProblem(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	// -----------------------------------------
	// Register Problem and associated test data
	// -----------------------------------------

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3: storage_address, size_train_dataset, test_dataset (address_data0, address_data1, ...)")
	}
	//       0          		1		 	         2
	// "storage_address", "size_train_dataset", "test_dataset"

	fmt.Println("- start create problem \n")

	// Clean input data
	sizeTrainDataset, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	testDataAddress := strings.Split(args[2], ",")

	// Create Problem Key
	problemKey := "problem_" + args[0]

	// Store test data
	err, testData := registerTestData(APIstub, problemKey, testDataAddress)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Store Problem
	var problem = Problem{ObjectType: "problem", StorageAddress: args[0], SizeTrainDataset: sizeTrainDataset, TestDataset: testData}
	problemAsBytes, err := json.Marshal(problem)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = APIstub.PutState(problemKey, problemAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("- end create problem")
	return shim.Success(nil)
}

func registerTestData(APIstub shim.ChaincodeStubInterface, problemKey string,
	testDataAddress []string) (err error, testData []string) {
	// ------------------
	// Register test data
	// ------------------
	for _, sdata := range testDataAddress {
		// remove leading and trailing space and split address and owner
		sdata = strings.TrimSpace(sdata)
		// create data key
		dataKey := "data_" + sdata
		// store data
		err, _ = storeItem(APIstub, dataKey, "data", problemKey)
		if err != nil {
			return err, testData
		}
		testData = append(testData, sdata)
		fmt.Printf("-- test data %s registered \n", dataKey)
	}

	return err, testData
}

// ================================================================================
// registerItem - register a new item (algo or data), store it into chaincode state
// ================================================================================

func storeItem(APIstub shim.ChaincodeStubInterface, itemKey string, itemType string,
	problem string) (err error, item Item) {
	// Store item in the ledger and create associated composite key
	// No learnuplet creation

	item = Item{ObjectType: itemType, Problem: problem}

	itemAsBytes, err := json.Marshal(item)
	if err != nil {
		return err, item
	}
	// Store item
	err = APIstub.PutState(itemKey, itemAsBytes)
	if err != nil {
		return err, item
	}

	// Create composite key to enable (itemtype + problem + itemKey)-based range queries,
	// e.g. return all items associated with a given problem
	indexName := item.ObjectType + "~problem~key"
	itemProblemIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{item.ObjectType, item.Problem, itemKey})
	if err != nil {
		return err, item
	}
	emptyValue := []byte{0x00}
	err = APIstub.PutState(itemProblemIndexKey, emptyValue)
	if err != nil {
		return err, item
	}

	return err, item
}

func (s *SmartContract) registerItem(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	// ---------------------------------------------------------------
	// store item in the ledger and create associated learnuplet
	// ---------------------------------------------------------------

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3: itemType, storage_address, problem")
	}
	//   0          	1   		 	2
	// "itemType", "storage_address", "problem"

	fmt.Println("- start create " + args[0])

	// Create item key
	itemKey := args[0] + "_" + args[1]
	// Store item in ledger and create composite key
	err, item := storeItem(APIstub, itemKey, args[0], args[2])
	if err != nil {
		return shim.Error(err.Error())
	}
	// Create associated learnuplet
	if args[0] == "algo" {
		fmt.Println("-- create associated learnuplets")
		algoLearnuplet(APIstub, itemKey, item)
	}
	if args[0] == "data" {
		fmt.Println("-- create associated learnuplets")
		data := []string{itemKey}
		dataLearnuplet(APIstub, data, item.Problem)
	}
	// ==== Algo saved and indexed. Return success ====
	fmt.Println("- end create " + item.ObjectType)
	return shim.Success(nil)
}

// ================================================================================
// queryItem - read one algo/problem/data/... given its key from chaincode state
// ================================================================================

func (s *SmartContract) queryItem(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1: item key")
	}

	key := args[0]
	fmt.Println("- start looking for element with key ", key)
	payload, err := APIstub.GetState(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("- end looking for element with key ", key)
	return shim.Success(payload)
}

// ===================================================================
// queryItems - get all algos/problems/datas/... given its object type
// ===================================================================

func (s *SmartContract) queryItems(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1: object type")
	}

	objectType := args[0]
	fmt.Printf("- start looking for elements of type %s\n", objectType)
	resultsIterator, _ := APIstub.GetStateByRange(objectType+"_", objectType+"_z")
	var items []map[string]interface{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		var item map[string]interface{}
		err = json.Unmarshal(queryResponse.GetValue(), &item)
		if err != nil {
			return shim.Error(err.Error())
		}
		item["key"] = queryResponse.GetKey()
		items = append(items, item)
	}
	fmt.Printf("- end looking for elements of type %s\n", objectType)

	payload, err := json.Marshal(items)
	if err != nil {
		return shim.Error(err.Error())
	}

	//return
	return shim.Success(payload)
}

// ==================================================================
// Query all items corresponding to a problem
// ==================================================================

func getProblemItems(APIstub shim.ChaincodeStubInterface, problem string, itemType string) (itemAddresses []string) {

	fmt.Printf("--- looking for %s associated with %s \n", itemType, problem)

	// Query the itemType~problem~key index by problem
	// This will execute a key range query on all keys starting with 'itemType~problem'
	problemAssociatedItemIterator, err := APIstub.GetStateByPartialCompositeKey(itemType+"~problem~key", []string{itemType, problem})
	if err != nil {
		return itemAddresses
	}
	defer problemAssociatedItemIterator.Close()

	// Iterate through result set and for each algo found
	var i int
	for i = 0; problemAssociatedItemIterator.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the item name from the composite key
		responseRange, err := problemAssociatedItemIterator.Next()
		if err != nil {
			return itemAddresses
		}

		// get the itemType, problem, and key from the composite key
		_, compositeKeyParts, err := APIstub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return itemAddresses
		}
		returnedProblem := compositeKeyParts[1]
		returnedKey := compositeKeyParts[2]

		fmt.Printf("--- found %s associated with %s \n", returnedKey, returnedProblem)

		// Put item key in slice
		newAddress := strings.TrimPrefix(returnedKey, itemType+"_")
		itemAddresses = append(itemAddresses, newAddress)
	}

	return itemAddresses
}

func (s *SmartContract) queryProblemItems(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	//   0      1
	// ObjectType ('algo' or 'data'), "problem", e.g. problem_0
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2: object type and problem key")
	}

	itemType := args[0]
	problem := args[1]
	fmt.Printf("- start of query %s related to %s\n", itemType, problem)

	// Get slice with keys of items associated to the problem
	itemKeys := getProblemItems(APIstub, problem, itemType)

	// Iterate through result set
	results := make(map[string]string)
	for _, key := range itemKeys {

		// Get algo given its key
		value, err := APIstub.GetState(key)
		if err != nil {
			return shim.Error(err.Error())
		}
		results[key] = string(value)
	}
	payload, err := json.Marshal(results)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("- end of query %s related to %s\n", itemType, problem)

	return shim.Success(payload)
}

// =============================================================
// query learnuplet ith specific composite keys (status or algo)
// =============================================================

func getCompositeLearnuplet(APIstub shim.ChaincodeStubInterface, keyRequest string,
	keyValue string) ([]byte, error) {
	// To get all learnuplet having a given status (keyRequest: status, keyValue: todo, ...)
	// or being linked with a given algo (keyRequest: algo, keyValue: algoKey)

	compositeKeyIndex := "learnuplet~" + keyRequest + "~key"
	// Query the learnuplet~<compositeKey>~key index by <compositeKey>
	learnupletIterator, err := APIstub.GetStateByPartialCompositeKey(compositeKeyIndex, []string{"learnuplet", keyValue})
	if err != nil {
		return nil, err
	}
	defer learnupletIterator.Close()

	var learnuplets []map[string]interface{}
	// Iterate through result set
	for i := 0; learnupletIterator.HasNext(); i++ {
		responseRange, err := learnupletIterator.Next()
		if err != nil {
			return nil, err
		}

		// get the ObjectType, status, and key from the composite key
		_, compositeKeyParts, err := APIstub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return nil, err
		}
		returnedKey := compositeKeyParts[2]
		value, _ := APIstub.GetState(returnedKey)
		var learnuplet map[string]interface{}
		err = json.Unmarshal(value, &learnuplet)
		if err != nil {
			return nil, err
		}
		learnuplet["key"] = returnedKey
		learnuplets = append(learnuplets, learnuplet)
	}

	payload, err := json.Marshal(learnuplets)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (s *SmartContract) queryStatusLearnuplet(APIstub shim.ChaincodeStubInterface,
	args []string) sc.Response {

	//   0
	// status
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1: asked learnuplet status")
	}

	status := args[0]
	fmt.Println("- start looking for learnuplet with status ", status)

	payload, err := getCompositeLearnuplet(APIstub, "status", status)
	if err != nil {
		return shim.Error("Problem querying learnuplet depending on status " +
			status + " - " + err.Error())
	}
	fmt.Println("- end looking for learnuplet with status ", status)

	return shim.Success(payload)
}

func (s *SmartContract) queryAlgoLearnuplet(APIstub shim.ChaincodeStubInterface,
	args []string) sc.Response {

	//   0
	// algoKey
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1: algo key")
	}

	algo := args[0]
	fmt.Println("- start looking for learnuplet of algo ", algo)

	payload, err := getCompositeLearnuplet(APIstub, "algo", algo)
	if err != nil {
		return shim.Error("Problem querying learnuplet associated with algo " +
			algo + " - " + err.Error())
	}
	fmt.Println("- end looking for learnuplet of algo ", algo)
	return shim.Success(payload)
}

// ====================================================================
// Learnuplet creation - functions called when registering data or algo
// ====================================================================

func getRankAlgoLearnuplet(APIstub shim.ChaincodeStubInterface, algoKey string) (rank int, modelAddress string) {
	fmt.Printf("--- looking for last learnuplet rank of %s \n", algoKey)

	// Query the learnuplet~algo~key index by algo
	algoAssociatedLearnupletIterator, err := APIstub.GetStateByPartialCompositeKey("learnuplet~algo~key", []string{"learnuplet", algoKey})
	if err != nil {
		return -1, ""
	}
	defer algoAssociatedLearnupletIterator.Close()

	// Iterate through result set and for each learnuplet found
	var newRank int
	var perf, newPerf float64
	rank = 0
	modelAddress = ""
	for i := 0; algoAssociatedLearnupletIterator.HasNext(); i++ {
		responseRange, err := algoAssociatedLearnupletIterator.Next()
		if err != nil {
			return -1, ""
		}

		// get the itemType, problem, and key from the composite key
		_, compositeKeyParts, err := APIstub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return -1, ""
		}
		returnedKey := compositeKeyParts[2]
		value, _ := APIstub.GetState(returnedKey)
		retrievedLearnuplet := Learnuplet{}
		err = json.Unmarshal(value, &retrievedLearnuplet)
		if err != nil {
			fmt.Errorf("Problem Unmarshal %s", returnedKey)
			return -1, ""
		}
		if i == 0 {
			perf = retrievedLearnuplet.Perf
		}
		newRank = retrievedLearnuplet.Rank
		newPerf = retrievedLearnuplet.Perf
		// If better perf, update modelAddess
		if retrievedLearnuplet.Status == "done" && newPerf >= perf {
			perf = newPerf
			modelAddress = retrievedLearnuplet.ModelEnd
		}
		// If greater rank, update rank
		if newRank >= rank {
			rank = newRank
			if retrievedLearnuplet.Status != "done" {
				modelAddress = ""
			}
		}

		fmt.Printf("- for algo %s: found last rank %s and associated model %s \n", algoKey, rank, modelAddress)
	}

	return rank, modelAddress
}

func createLearnuplet(
	APIstub shim.ChaincodeStubInterface, trainData []string, szBatch int,
	testData []string, problem string, problemAddress string, algo string,
	modelStart string, startRank int) (nbNewLearnuplet int) {

	nbNewLearnuplet = 0
	var batchData []string
	// create empty maps for performances
	var trainPerf, testPerf map[string]float64
	trainPerf = make(map[string]float64)
	testPerf = make(map[string]float64)
	// For each mini-batch of data, create a learnuplet
	for i, j := 0, 0; i < len(trainData); i, j = i+szBatch, j+1 {
		if i+szBatch >= len(trainData) {
			batchData = trainData[i:]

		} else {
			batchData = trainData[i : i+szBatch]
		}
		j = j + startRank
		// if not first rank, modelStart is empty, will be filled once first rank has been computed
		// ModelEnd will be sent by compute when it sends status and performances
		learnupletModelStart := ""
		if j == startRank {
			learnupletModelStart = modelStart
		}
		newLearnuplet := Learnuplet{
			ObjectType:     "learnuplet",
			Problem:        problem,
			ProblemAddress: problemAddress,
			Algo:           algo,
			ModelStart:     learnupletModelStart,
			ModelEnd:       "",
			TrainData:      batchData,
			TestData:       testData,
			Worker:         "",
			Status:         "todo",
			Rank:           j,
			Perf:           0,
			TrainPerf:      trainPerf,
			TestPerf:       testPerf,
		}
		// Append to ledger
		learnupletKey := "learnuplet_" + uuid.NewV4().String()
		newLearnupletAsBytes, err := json.Marshal(newLearnuplet)
		if err != nil {
			fmt.Errorf("Problem marshaling ", learnupletKey)
		}
		err = APIstub.PutState(learnupletKey, newLearnupletAsBytes)
		if err != nil {
			fmt.Errorf("Problem putting state of ", learnupletKey)
		} else {
			nbNewLearnuplet++
			// Create composite key learnuplet~algo~key
			indexName := "learnuplet~algo~key"
			learnupletAlgoIndexKey, _ := APIstub.CreateCompositeKey(indexName, []string{"learnuplet", algo, learnupletKey})
			value := []byte{0x00}
			APIstub.PutState(learnupletAlgoIndexKey, value)
			// Create composite key learnuplet~status~key
			indexName = "learnuplet~status~key"
			learnupletStatusIndexKey, _ := APIstub.CreateCompositeKey(indexName, []string{"learnuplet", "todo", learnupletKey})
			APIstub.PutState(learnupletStatusIndexKey, value)
			fmt.Printf("-- creation of %s ok \n", learnupletKey)

		}
	}

	return nbNewLearnuplet
}

func algoLearnuplet(APIstub shim.ChaincodeStubInterface, algoKey string, algo Item) int {
	// ---------------------------------------------
	// create learnuplet when new algo is registered
	// ---------------------------------------------

	problem := algo.Problem

	// Find test data
	value, err := APIstub.GetState(problem)
	if err != nil {
		fmt.Errorf("%s not found", problem)
		return 0
	}
	retrievedProblem := Problem{}
	err = json.Unmarshal(value, &retrievedProblem)
	if err != nil {
		fmt.Errorf("Problem Unmarshal %s", problem)
		return 0
	}
	testData := retrievedProblem.TestDataset
	sizeTrainDataset := retrievedProblem.SizeTrainDataset
	problemAddress := retrievedProblem.StorageAddress
	// Find all active data associated to the same problem and remove test data
	trainData := getProblemItems(APIstub, problem, "data")
	for i := 0; i < len(trainData); i++ {
		itraindata := trainData[i]
		for _, itestdata := range testData {
			if itraindata == itestdata {
				trainData = append(trainData[:i], trainData[i+1:]...)
				i--
				continue
			}
		}
	}
	sort.Strings(trainData)
	// Create learnuplets
	nbNewLearnuplet := createLearnuplet(
		APIstub, trainData, sizeTrainDataset, testData, problem, problemAddress,
		algoKey, algoKey, 0)
	return nbNewLearnuplet
}

func dataLearnuplet(APIstub shim.ChaincodeStubInterface, data []string, problem string) int {
	// ---------------------------------------------
	// create learnuplet when new data is registered
	// ---------------------------------------------

	nbNewLearnuplet := 0

	// Find test data
	value, err := APIstub.GetState(problem)
	if err != nil {
		fmt.Errorf("%s not found", problem)
		return 0
	}
	retrievedProblem := Problem{}
	err = json.Unmarshal(value, &retrievedProblem)
	if err != nil {
		fmt.Errorf("Problem Unmarshal %s", problem)
		return 0
	}
	testData := retrievedProblem.TestDataset
	sizeTrainDataset := retrievedProblem.SizeTrainDataset
	problemAddress := retrievedProblem.StorageAddress
	// Find all active algo associated to the same problem
	algoKeys := getProblemItems(APIstub, problem, "algo")
	// For each algo, find the last rank and create learnuplet
	var rank int
	var modelAddress string
	for _, algoKey := range algoKeys {
		rank, modelAddress = getRankAlgoLearnuplet(APIstub, algoKey)
		nbAlgoNewLearnuplet := createLearnuplet(
			APIstub, data, sizeTrainDataset, testData, problem, problemAddress,
			algoKey, modelAddress, rank)
		nbNewLearnuplet = nbNewLearnuplet + nbAlgoNewLearnuplet
	}
	return nbNewLearnuplet
}

// ====================================================
// Push from Compute: update learnuplets and preduplets
// ====================================================

func (s *SmartContract) setUpletWorker(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	// -----------------------------------------------------------------------------
	// Set a worker for a learnuplet
	// For now, this is a simple function, much more checks will be applied later...
	// -----------------------------------------------------------------------------

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2: uplet_key, worker")
	}
	//       0  		1
	// "uplet_key", "worker"
	upletKey := args[0]
	worker := args[1]
	fmt.Printf("- start set worker for %s \n", upletKey)

	value, _ := APIstub.GetState(upletKey)
	retrievedLearnuplet := Learnuplet{}
	if value == nil {
		return shim.Error("No learnuplet with key - " + upletKey)
	}
	err := json.Unmarshal(value, &retrievedLearnuplet)
	if err != nil {
		return shim.Error("Problem Unmarshal uplet - " + err.Error())
	}
	if retrievedLearnuplet.Status == "pending" {
		return shim.Error("Uplet status is already pending...")
	} else {
		retrievedLearnuplet.Status = "pending"
		retrievedLearnuplet.Worker = worker
		learnupletAsBytes, err := json.Marshal(retrievedLearnuplet)
		if err != nil {
			return shim.Error("Problem (re)marshaling uplet - " + err.Error())
		}
		err = APIstub.PutState(upletKey, learnupletAsBytes)
		if err != nil {
			return shim.Error("Problem storing uplet - " + err.Error())
		}
		// Update associated composite key learnuplet~status~key
		indexName := "learnuplet~status~key"
		emptyValue := []byte{0x00}
		oldLearnupletStatusIndexKey, _ := APIstub.CreateCompositeKey(indexName, []string{"learnuplet", "todo", upletKey})
		APIstub.DelState(oldLearnupletStatusIndexKey)
		learnupletStatusIndexKey, _ := APIstub.CreateCompositeKey(indexName, []string{"learnuplet", "pending", upletKey})
		APIstub.PutState(learnupletStatusIndexKey, emptyValue)
	}
	fmt.Printf("- end set worker for %s \n", upletKey)
	return shim.Success(nil)
}

func (s *SmartContract) reportLearn(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	// -----------------------------------------------------------------------------
	// Set output of a learnuplet, updateing the corresponding learnuplet
	// For now, this is a simple function, much more checks will be applied later...
	// -----------------------------------------------------------------------------

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5: uplet_key, status (failed / done), perf, train_perf (train_data_i perf_i, train_data_j perf_j, ...), test_perf (test_data_i perf_j, test_data_j perf_j, ...)")
	}
	//     0  		   1		  2 		3			4
	// "uplet_key", "status", "perf", "train_perf", "test_perf"

	// TODO: Check args validity (especially uplet_key and status validity)

	upletKey := args[0]
	fmt.Printf("- start Report learning phase of %s \n", upletKey)

	// Get learnuplet
	value, _ := APIstub.GetState(upletKey)
	retrievedLearnuplet := Learnuplet{}
	err := json.Unmarshal(value, &retrievedLearnuplet)
	if err != nil {
		return shim.Error(fmt.Sprintf("Error Unmarshal uplet %s - %s", upletKey, err))
	}

	// Update learnuplet status
	retrievedLearnuplet.Status = args[1]

	// Deal with the status "failed" case
	if retrievedLearnuplet.Status == "failed" {
		// Store updated learnuplet
		learnupletAsBytes, err := json.Marshal(retrievedLearnuplet)
		if err != nil {
			return shim.Error(fmt.Sprintf("Error re-Unmarshal uplet %s - %s", upletKey, err))
		}
		err = APIstub.PutState(upletKey, learnupletAsBytes)
		if err != nil {
			return shim.Error("Problem storing learnuplet - " + err.Error())
		}
		// Update associated composite key learnuplet~status~key
		indexName := "learnuplet~status~key"
		emptyValue := []byte{0x00}
		oldLearnupletStatusIndexKey, _ := APIstub.CreateCompositeKey(indexName, []string{"learnuplet", "pending", upletKey})
		APIstub.DelState(oldLearnupletStatusIndexKey)
		learnupletStatusIndexKey, _ := APIstub.CreateCompositeKey(indexName, []string{"learnuplet", "failed", upletKey})
		APIstub.PutState(learnupletStatusIndexKey, emptyValue)

		fmt.Printf("- end Report learning phase of %s \n", upletKey)
		return shim.Success(nil)
	}

	// Process perf data
	var perf float64
	var trainPerf, testPerf map[string]float64
	trainPerf = make(map[string]float64)
	testPerf = make(map[string]float64)
	if args[2] != "" {
		perf, err = strconv.ParseFloat(args[2], 64)
		if err != nil {
			return shim.Error("Problem parsing performance - " + err.Error())
		}
		// TODO check data addresses correspond to train and test data
		err, trainPerf = mapPerf(args[3])
		if err != nil {
			return shim.Error("Problem parsing train perf - " + err.Error())
		}
		err, testPerf = mapPerf(args[4])
		if err != nil {
			return shim.Error("Problem parsing test perf - " + err.Error())
		}
	}

	// Update Learnuplet Perf results
	retrievedLearnuplet.Perf = perf
	retrievedLearnuplet.TrainPerf = trainPerf
	retrievedLearnuplet.TestPerf = testPerf

	// Store updated learnuplet
	learnupletAsBytes, err := json.Marshal(retrievedLearnuplet)
	if err != nil {
		return shim.Error("Problem (re)marshaling learnuplet - " + err.Error())
	}
	err = APIstub.PutState(upletKey, learnupletAsBytes)
	if err != nil {
		return shim.Error("Problem storing learnuplet - " + err.Error())
	}
	// Update associated composite key learnuplet~status~key
	indexName := "learnuplet~status~key"
	emptyValue := []byte{0x00}
	oldLearnupletStatusIndexKey, _ := APIstub.CreateCompositeKey(indexName, []string{"learnuplet", "pending", upletKey})
	APIstub.DelState(oldLearnupletStatusIndexKey)
	learnupletStatusIndexKey, _ := APIstub.CreateCompositeKey(indexName, []string{"learnuplet", "done", upletKey})
	APIstub.PutState(learnupletStatusIndexKey, emptyValue)

	fmt.Printf("- end Report learning phase of %s \n", upletKey)
	return shim.Success(nil)
}

func mapPerf(perf string) (err error, mapPerf map[string]float64) {
	// Convert string with perf ("dataKey_i perf_i, dataKey_j perf_j, ...") to map
	var keyPerf []string
	var p float64
	mapPerf = make(map[string]float64)
	slicePerf := strings.Split(perf, ",")
	for _, sperf := range slicePerf {
		// remove leading and trailing space and split key and perf
		sperf = strings.TrimSpace(sperf)
		keyPerf = strings.Split(sperf, " ")
		// store in map
		p, err = strconv.ParseFloat(keyPerf[1], 64)
		if err != nil {
			return err, mapPerf
		}
		mapPerf[keyPerf[0]] = p
	}
	return err, mapPerf
}

// ==============================================
// MAIN FUNCTION. Only relevant in unit test mode
// ==============================================
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
