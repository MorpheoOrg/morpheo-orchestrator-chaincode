package main

import	(
	"testing"
	"fmt"
	// "bytes"
	"encoding/json"
	"strconv"
	"reflect"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	)

func TestRegisterItem(t *testing.T) {
	// ARRANGE
   	smartContract := new(SmartContract)
   	mockStub := shim.NewMockStub("mockstub", smartContract)
   	txId := "mockTxID"
    itemKey := "8fa81bfc-b5f4-4ba2-b81a-b46424800000"
   	args := []string{"data", itemKey, "problem_1",}

   	// ACT
   	mockStub.MockTransactionStart(txId)
   	response := smartContract.registerItem(mockStub, args)
   	mockStub.MockTransactionEnd(txId)

  	// ASSERT
   	// response
   	if s := response.GetStatus(); s != 200 {
      	t.Errorf("the status is %d, instead of 200", s)
      	t.Errorf("message: %s", response.Message)
   	}
   	// storing in db
   	itemAsBytes, err := mockStub.GetState("data_" + itemKey)
   	if err != nil {
		t.Errorf("Get state did not work")
	}
	item := Item{}
	err = json.Unmarshal(itemAsBytes, &item)
    if item.ObjectType != "data" || item.Problem != "problem_1" {
		t.Errorf("Registration of item fails")
	}
}

func TestQueryItem(t *testing.T) {
	// ARRANGE
    smartContract := new(SmartContract)
    mockStub := shim.NewMockStub("mockstub", smartContract)
    txId := "mockTxID"
    itemKey := "8fa81bfc-b5f4-4ba2-b81a-b46424800400"
    args := []string{"algo", itemKey, "problem_1",}

    // ACT
    mockStub.MockTransactionStart(txId)
    // add data item
    smartContract.registerItem(mockStub, args)
    response := smartContract.queryItem(mockStub, []string{"algo_"+ itemKey,})
    mockStub.MockTransactionEnd(txId)
    // format item
    itemQueried := Item{}
    err := json.Unmarshal(response.GetPayload(), &itemQueried)
    if err != nil {
        t.Errorf("Problem wih json.Unmarshal")
    }

    // ASSERT
    // response
    if s := response.GetStatus(); s != 200 {
    	t.Errorf("the status is %d, instead of 200", s)
    	t.Errorf("message: %s", response.Message)
    }
    // collecting in db
    itemAsBytes, err := mockStub.GetState("algo_" + itemKey)
    if err != nil {
        t.Errorf("Get state did not work")
    }
    item := Item{}
    err = json.Unmarshal(itemAsBytes, &item)
    if item.Problem != itemQueried.Problem || item.ObjectType != itemQueried.ObjectType {
        t.Errorf("Query of item fails")
    }
}

func TestRegisterProblem(t *testing.T) {
	// ARRANGE
   	smartContract := new(SmartContract)
   	mockStub := shim.NewMockStub("mockstub", smartContract)
   	txId := "mockTxID"
    problemId := "dda81bfc-b5f4-5ba2-b81a-b464248f02d2"
   	args := []string{
   					problemId, // suffix of problem problemKey
   					"2", // size of learnuplets
   					"data_0, data_1", // test dataset
   				}

   	// ACT
   	mockStub.MockTransactionStart(txId)
   	// add problem
   	response := smartContract.registerProblem(mockStub, args)
   	mockStub.MockTransactionEnd(txId)

  	// ASSERT
   	// response
   	if s := response.GetStatus(); s != 200 {
      	t.Errorf("the status is %d, instead of 200", s)
      	t.Errorf("message: %s", response.Message)
   	}
   	// storing in db
   	problemAsBytes, err := mockStub.GetState("problem_"+problemId)
   	if err != nil {
		t.Errorf("Get state did not work")
	}
	problem := Problem{}
	err = json.Unmarshal(problemAsBytes, &problem)
	if problem.SizeTrainDataset != 2 || problem.TestDataset[1] != "data_1"{
		t.Errorf("Registration of problem fails")
	}
}

func TestInitLedger(t *testing.T) {
   // ARRANGE
      smartContract := new(SmartContract)
   mockStub := shim.NewMockStub("mockstub", smartContract)
   txId := "mockTxID"

    // ACT
   mockStub.MockTransactionStart(txId)
   response := smartContract.initLedger(mockStub)
   mockStub.MockTransactionEnd(txId)

      // ASSERT
      // response
      if s := response.GetStatus(); s != 200 {
         t.Errorf("the status is %d, instead of 200", s)
         t.Errorf("message: %s", response.Message)
      }
      // storing in db
      // data
      itemAsBytes, err := mockStub.GetState("data_0")
      if err != nil {
      t.Errorf("Get state did not work")
   }
   item := Item{}
   err = json.Unmarshal(itemAsBytes, &item)
   if item.Problem != "problem_0"{
      t.Errorf("Registration/query fails")
   }
}

func TestGetProblemItems(t *testing.T) {
	// ARRANGE
   	smartContract := new(SmartContract)
   	mockStub := shim.NewMockStub("mockstub", smartContract)
   	txId := "mockTxID"
   	// prepare variables
   	pbl := "problem_0"
   	itTyp := "data"
    id1 := "8fa81bfc-b5f4-4ba2-b81a-b46424800000"
    id3 := "8fa81bfc-b5f4-4ba2-b81a-b46424800002"
   	args_1 := []string{"data",
					id1,
   					"problem_0",}
   	args_2 := []string{"data",
					"8fa81bfc-b5f4-4ba2-b81a-b46424800001",
   					"problem_1",}
   	args_3 := []string{"data",
					"8fa81bfc-b5f4-4ba2-b81a-b46424800002",
   					"problem_0",}
   	// ACT
   	mockStub.MockTransactionStart(txId)
   // add data items
   	smartContract.registerItem(mockStub, args_1)
   	smartContract.registerItem(mockStub, args_2)
   	smartContract.registerItem(mockStub, args_3)
   	// call the function to be tested
   	res := getProblemItems(mockStub, pbl, itTyp)
   	mockStub.MockTransactionEnd(txId)

  	// ASSERT
   	// response
   	eq := reflect.DeepEqual(res, []string{id1, id3})
   	if !eq {
		t.Errorf("getProblemItems did not work")
	}
}


func TestQueryProblemItems(t *testing.T) {
	// ARRANGE
   	smartContract := new(SmartContract)
   	mockStub := shim.NewMockStub("mockstub", smartContract)
   	txId := "mockTxID"
   	// prepqre variables
    dataKey := "8fa81bfc-b5f4-4ba2-b81a-b46424800001"
   	args := []string{"data", "problem_1"}
   	args_1 := []string{"data",
   					"8fa81bfc-b5f4-4ba2-b81a-b46424800000",
   					"problem_0",}
   	args_2 := []string{"data",
                    dataKey,
   					"problem_1",}

   	// ACT
   	mockStub.MockTransactionStart(txId)
   	// add data items
   	smartContract.registerItem(mockStub, args_1)
   	smartContract.registerItem(mockStub, args_2)
   	// call the function to be tested
   	response := smartContract.queryProblemItems(mockStub, args)
   	mockStub.MockTransactionEnd(txId)
    mapQueried := map[string]string{}
    err := json.Unmarshal(response.GetPayload(), &mapQueried)
    if err != nil {
        t.Errorf("Unmarshal did not work")
    }
    // itemQueriedAsBytes := mapQueried[dataKey]

    // manually collecting in db for comparison
    itemAsBytes, err := mockStub.GetState("data_" + dataKey)
    if err != nil {
        t.Errorf("Get state did not work")
    }
    itemDb := Item{}
    err = json.Unmarshal(itemAsBytes, &itemDb)
    
  	// ASSERT
   	if s := response.GetStatus(); s != 200 {
      	t.Errorf("the status is %d, instead of 200", s)
      	t.Errorf("message: %s", response.Message)
   	}
    fmt.Println("######################")
    fmt.Println(itemDb)
    fmt.Println(mapQueried)
    fmt.Println("######################")
 //   	if data_key != itemQueried{
	// 	t.Errorf("QueryProblemItems did not work")
	// }

}

func TestCreateLearnpulet(t *testing.T) {
	// TODO check what happens when learnuplet is created without the data existing in the db

	// ARRANGE
   	smartContract := new(SmartContract)
   	mockStub := shim.NewMockStub("mockstub", smartContract)
   	txId := "mockTxID"
	// preparing variables
   	trDataAdd := []string{"8fa81bfc-b5f4-4ba2-b81a-b46424800000",
                          "8fa81bfc-b5f4-4ba2-b81a-b46424800001",
	                      "8fa81bfc-b5f4-4ba2-b81a-b46424800003",
	                      "8fa81bfc-b5f4-4ba2-b81a-b46424800005",}
   	sz_batch := 2
   	teDataAdd := []string{"8fa81bfc-b5f4-4ba2-b81a-b46424800002",
	                      "8fa81bfc-b5f4-4ba2-b81a-b46424800004",}
	pbl := "problem_dda81bfc-b5f4-5ba2-b81a-b464248f0000"
	alg := "algo_0"
	mdlStart := ""
	strtRk := 0
	// args_1 := []string{"data",
 //   					"bobobor",
	// 				"https://storage.morpheo.io/data/8fa81bfc-b5f4-4ba2-b81a-b46424800000",
 //   					"problem_0",
 //   					"14991"}
 //   	args_2 := []string{"data",
 //   					"bobobor",
	// 				"https://storage.morpheo.io/data/8fa81bfc-b5f4-4ba2-b81a-b46424800001",
 //   					"problem_0",
 //   					"14992"}

   	// ACT
   	mockStub.MockTransactionStart(txId)
   	// initialize datatop
   	// mockStub.PutState("datatop", []byte(strconv.Itoa(-1)))
   	// add data item
   	// smartContract.registerItem(mockStub, args_1)
   	// smartContract.registerItem(mockStub, args_2)
   	
   	// add data item
   	nbLearnuplets := createLearnuplet(mockStub, trDataAdd, sz_batch, teDataAdd, pbl, alg, mdlStart, strtRk)
   	mockStub.MockTransactionEnd(txId)

  	// ASSERT
   	// check number of created learnuplets
   	if nbLearnuplets != 2 {
   		t.Errorf("Wrong number of created learnuplets")
   	}
   	// checks learnuplet_0 created in db
   	learnupletAsBytes, err := mockStub.GetState("learnuplet_0")
	if err != nil {
		t.Errorf("Get state did not work")
	}
	learnuplet := Learnuplet{}
	err = json.Unmarshal(learnupletAsBytes, &learnuplet)
	m := map[string]string{
	"data_0": "8fa81bfc-b5f4-4ba2-b81a-b46424800000",
	"data_1": "8fa81bfc-b5f4-4ba2-b81a-b46424800001",}
	eq := reflect.DeepEqual(learnuplet.TrainData, m)
   	if !eq {
		t.Errorf("Creation of learnuplet fails")
	}   	
	// checks learnuplet_1 created in db
   	learnupletAsBytes, err = mockStub.GetState("learnuplet_1")
   	if err != nil {
		t.Errorf("Get state did not work")
	}
	learnuplet = Learnuplet{}
	err = json.Unmarshal(learnupletAsBytes, &learnuplet)
	if learnuplet.Problem != pbl {
		t.Errorf("Creation of learnuplet fails")
	}
}

func TestAlgoLearnuplet(t *testing.T) {
	// ARRANGE
   	smartContract := new(SmartContract)
   	mockStub := shim.NewMockStub("mockstub", smartContract)
   	txId := "mockTxID"
	// preparing variables
	alg := Item{ObjectType: "algo", Problem: "problem_0"}
	args_pbl := []string{
   					"dda81bfc-b5f4-5ba2-b81a-b464248f02d2",
   					"2", // size of learnuplets
   					"data_0, data_1", // test dataset
   				}
	args_1 := []string{"data",
					"8fa81bfc-b5f4-4ba2-b81a-b46424800000",
   					"problem_0",}
   	args_2 := []string{"data",
					"8fa81bfc-b5f4-4ba2-b81a-b46424800001",
   					"problem_0",}
   	args_3 := []string{"data",
					"8fa81bfc-b5f4-4ba2-b81a-b46424800002",
   					"problem_0",}
	args_4 := []string{"data",
					"8fa81bfc-b5f4-4ba2-b81a-b46424800003",
   					"problem_0",}

   	// ACT
   	mockStub.MockTransactionStart(txId)
   	// initialize tops
   	mockStub.PutState("problemtop", []byte(strconv.Itoa(-1)))
   	mockStub.PutState("datatop", []byte(strconv.Itoa(-1)))
   	// add problem and data
   	smartContract.registerProblem(mockStub, args_pbl)
   	smartContract.registerItem(mockStub, args_1)
   	smartContract.registerItem(mockStub, args_2)
   	smartContract.registerItem(mockStub, args_3)
   	smartContract.registerItem(mockStub, args_4)
   	// initialize learnuplettop
   	mockStub.PutState("learnuplettop", []byte(strconv.Itoa(-1)))
   	// add data item
   	nbLearnuplets := algoLearnuplet(mockStub, "algo_8fa81bfc-b5f4-4ba2-b81a-b464248f02d1", alg)
   	mockStub.MockTransactionEnd(txId)

  	// ASSERT
   	// check number of created learnuplets
   	if nbLearnuplets != 1 {
   		t.Errorf("Wrong number of created learnuplets")
   	}
   	// checks learnuplet_0 created in db
   	learnupletAsBytes, err := mockStub.GetState("learnuplet_0")
	if err != nil {
		t.Errorf("Get state did not work")
	}
	learnuplet := Learnuplet{}
	err = json.Unmarshal(learnupletAsBytes, &learnuplet)
	m := []string{
        "8fa81bfc-b5f4-4ba2-b81a-b46424800002", 
        "8fa81bfc-b5f4-4ba2-b81a-b46424800003",
    }
	eq := reflect.DeepEqual(learnuplet.TrainData, m)
   	if !eq {
		t.Errorf("Creation of learnuplet fails")
	}
}