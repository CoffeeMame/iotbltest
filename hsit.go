package main

// Cloudant NoSQLDBからBlockchainへのデータ受け渡しテスト
import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
//	"encoding/json"
//	"strconv"
//	"strings"
//	"regexp"
)

//==============================================================================================================================
//name for the key/value that will store a list of all known data
//KVSの末尾に格納するデータ一覧リストのキー（未実装）
//==============================================================================================================================
const DATA_INDEX_STR = "_dataindex"

//==============================================================================================================================
//	 Structure Definitions
//==============================================================================================================================
//	Chaincode - A blank struct for use with Shim (A HyperLedger included go file used for get/put state
//				and other HyperLedger functions)
//==============================================================================================================================
type SimpleChaincode struct {
}

//==============================================================================================================================
//	 Router Functions
//==============================================================================================================================
// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string,args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "add" {
		// Create new data
		return t.add_data(stub, args)
	} else if function == "del" {
		return t.del_data(stub, args)
	}

	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "dummy_query" {
		fmt.Println("hi there " + function)
		return nil, nil
	} else if function == "read" {
		return t.read_data(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}



// ============================================================================================================================
// read_data:データを参照
// ============================================================================================================================
func (t *SimpleChaincode) read_data(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var id, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting ID of the var to query")
	}

	id = args[0]
	// 指定されたIDのデータの中身を取得する
	valAsbytes, err := stub.GetState(id)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + id + "\"}"
		return nil, errors.New(jsonResp)
	}

	//send it onward
	return valAsbytes, nil

}

//=================================================================================================================================
// add_data:データを追加
//=================================================================================================================================
func (t *SimpleChaincode) add_data(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var id, val string
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	//input sanitation
	fmt.Println("- start add data")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}

	// Variables to define the JSON
	id = args[0]
	val = args[1]

	//check if data already exists
	dataAsBytes, err := stub.GetState(id)
	if err != nil {
		return nil, errors.New("Failed to get data name")
	}

	if dataAsBytes != nil {
		fmt.Println("This data arleady exists: " + id)
		//all stop a data by this ID exists
		return nil, errors.New("This data arleady exists")
	}

	//store data with ID as key
	err = stub.PutState(id, []byte(val))
	if err != nil {
		return nil, err
	}

	return []byte(args[0]), nil
}


//=============================================================================================================================
// del_data:IDを指定して荷物を削除する
//=============================================================================================================================
func (t *SimpleChaincode) del_data(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	// 引数の数をチェック
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	// 引数が空でないことをチェック
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}

	// 要素を削除
	err = stub.DelState(args[0])
	if err != nil {
		return nil, errors.New("Fail to Delete Data")
	}

	// IDのスライスから削除（未実装）

	return nil, nil
}
