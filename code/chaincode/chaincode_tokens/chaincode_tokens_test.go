package main

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.Fail()
	}
}

func checkInitFail(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status == shim.OK {
		fmt.Println("Init should fail but it did not", string(res.Message))
		t.Fail()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.Fail()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was", string(bytes), "instead required", value)
		t.Fail()
	}
}

func checkCreateAccount(t *testing.T, stub *shim.MockStub, accID string, name string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("createAccount"), []byte(accID), []byte(name)})
	if res.Status != shim.OK {
		fmt.Println("createAccount", accID, name, "failed", string(res.Message))
		t.Fail()
	}
	if res.Payload == nil {
		fmt.Println("createAccount", accID, "failed to get response")
		t.Fail()
	}
	if string(res.Payload) != "Account created" {
		fmt.Println("createAccount", accID, "response was not as expected")
		t.Fail()
	}
}

func checkCreateAccountFail(t *testing.T, stub *shim.MockStub, accID string, name string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("createAccount"), []byte(accID), []byte(name)})
	if res.Status == shim.OK {
		fmt.Println("createAccount", accID, name, "should fail but did not", string(res.Payload))
		t.Fail()
	}
	expectedPayload := "This account already exists: " + accID
	if string(res.Message) != expectedPayload {
		fmt.Println("Expected payload:", expectedPayload)
		fmt.Println("Instead got this:", string(res.Payload))
		t.Fail()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.Fail()
	}
}

func checkInvokeFail(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status == shim.OK {
		fmt.Println("Invoke", args, "should fail but did not", string(res.Payload))
		t.Fail()
	}
}

func checkInvokeResponse(t *testing.T, stub *shim.MockStub, args [][]byte, expectedPayload string) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.Fail()
	}
	if string(res.Payload) != expectedPayload {
		fmt.Println("Expected payload:", expectedPayload)
		fmt.Println("Instead got this:", string(res.Payload))
		t.Fail()
	}
}

func checkInvokeResponseFail(t *testing.T, stub *shim.MockStub, args [][]byte, expectedMessage string) {
	res := stub.MockInvoke("1", args)
	if res.Status == shim.OK {
		fmt.Println("Invoke", args, "should fail")
		fmt.Println("Instead got payload:", string(res.Payload))
		t.Fail()
	}
	if res.Message != expectedMessage {
		fmt.Println("Expected message:", expectedMessage)
		fmt.Println("Instead got this:", res.Message)
		t.Fail()
	}
}

func Test_Init(t *testing.T) {
	cc := new(Chaincode)
	stub := shim.NewMockStub("tokens_init_test", cc)

	// It should Init 1 account with 10 000 tokens
	checkInit(t, stub, [][]byte{[]byte("1"), []byte("10000")})
	checkState(t, stub, "1", "{\"RecordType\":\"ACCOUNT\",\"AccountID\":\"1\",\"Name\":\"Init_Account\",\"Tokens\":10000}")

	// It should Init 4 accounts with 10 000 tokens
	stub = shim.NewMockStub("tokens_init_test", cc)
	checkInit(t, stub, [][]byte{[]byte("4"), []byte("10000")})
	for i := 1; i < 5; i++ {
		checkState(t, stub, strconv.Itoa(i), "{\"RecordType\":\"ACCOUNT\",\"AccountID\":\""+strconv.Itoa(i)+"\",\"Name\":\"Init_Account\",\"Tokens\":10000}")
	}

	// It should Init 0 accounts with 0 tokens
	stub = shim.NewMockStub("tokens_init_test", cc)
	checkInit(t, stub, [][]byte{[]byte("0"), []byte("0")})

	// It should not Init negative number of accounts
	stub = shim.NewMockStub("tokens_init_test", cc)
	checkInitFail(t, stub, [][]byte{[]byte("-4"), []byte("10000")})

	// It should not Init an account with negative number of tokens
	stub = shim.NewMockStub("tokens_init_test", cc)
	checkInitFail(t, stub, [][]byte{[]byte("1"), []byte("-10")})

	// It should not Init with less args than 2
	stub = shim.NewMockStub("tokens_init_test", cc)
	checkInitFail(t, stub, [][]byte{[]byte("1")})

	// It should not Init with first empty arg
	stub = shim.NewMockStub("tokens_init_test", cc)
	checkInitFail(t, stub, [][]byte{[]byte(""), []byte("10")})

	// It should not Init with second empty arg
	stub = shim.NewMockStub("tokens_init_test", cc)
	checkInitFail(t, stub, [][]byte{[]byte("1"), []byte("")})
}

func Test_InvokeFail(t *testing.T) {
	cc := new(Chaincode)
	stub := shim.NewMockStub("invoke_fail_test", cc)
	args := [][]byte{[]byte("NoFunction"), []byte("test")}
	expectedMessage := "Received unknown function invocation"
	checkInvokeResponseFail(t, stub, args, expectedMessage)
}

func Test_createAccount(t *testing.T) {
	cc := new(Chaincode)
	stub := shim.NewMockStub("create_acc_test", cc)

	// Init 1 account with 10 000 tokens
	checkInit(t, stub, [][]byte{[]byte("1"), []byte("10000")})
	checkCreateAccount(t, stub, "2", "marcel")
	checkCreateAccount(t, stub, "0", "marcel")
	checkCreateAccount(t, stub, "-2", "marcel")
	checkCreateAccountFail(t, stub, "2", "marcel")
}

func Test_deleteAccountByID(t *testing.T) {
	cc := new(Chaincode)
	stub := shim.NewMockStub("tokens_init_test", cc)

	// Init 1 account with 10 000 tokens
	checkInit(t, stub, [][]byte{[]byte("1"), []byte("10000")})
	// it should delete account that have 0 tokens
	checkCreateAccount(t, stub, "2", "marcel")
	args := [][]byte{[]byte("deleteAccountByID"), []byte("2")}
	checkInvoke(t, stub, args)

	// it should not be possible to delete an account that have tokens
	args = [][]byte{[]byte("deleteAccountByID"), []byte("1")}
	checkInvokeFail(t, stub, args)
}

func Test_getAccountByID(t *testing.T) {
	cc := new(Chaincode)
	stub := shim.NewMockStub("get_account_test", cc)

	// Init 1 account with 10 tokens
	checkInit(t, stub, [][]byte{[]byte("1"), []byte("10")})

	// It should get account with ID "1" that was Init
	args := [][]byte{[]byte("getAccountByID"), []byte("1")}
	expectedPayload := "{\"RecordType\":\"ACCOUNT\",\"AccountID\":\"1\",\"Name\":\"Init_Account\",\"Tokens\":10}"
	checkInvokeResponse(t, stub, args, expectedPayload)

	// It should fail with empty string arg
	args = [][]byte{[]byte("getAccountByID"), []byte("")}
	checkInvokeFail(t, stub, args)

	// It should fail with more than one arg
	args = [][]byte{[]byte("getAccountByID"), []byte("1"), []byte("a")}
	checkInvokeFail(t, stub, args)

	/*
		// This cannot be tested because of the limitations of MockStub implementation
		// It should fail to get account that is not created
		args = [][]byte{[]byte("getAccountByID"), []byte("2")}
		expectedMessage := ""
		checkInvokeResponseFail(t, stub, args, expectedMessage)
	*/
}

func Test_getAccountHistoryByID(t *testing.T) {
	cc := new(Chaincode)
	stub := shim.NewMockStub("get_history_test", cc)

	// Init 1 account with 10 tokens
	checkInit(t, stub, [][]byte{[]byte("1"), []byte("10")})
	/*
		// This cannot be tested because of the limitations of MockStub implementation
		// It should return history for account ID "1"
		args := [][]byte{[]byte("getAccountHistoryByID"), []byte("1")}
		expectedPayload := ""
		checkInvokeResponse(t, stub, args, expectedPayload)
	*/

	// It should fail with empty string arg
	args := [][]byte{[]byte("getAccountHistoryByID"), []byte("")}
	expectedMessage := "1st argument must be a non-empty string"
	checkInvokeResponseFail(t, stub, args, expectedMessage)

	// It should fail with more than one args
	args = [][]byte{[]byte("getAccountHistoryByID"), []byte("1"), []byte("lol")}
	expectedMessage = "Incorrect number of arguments. Expecting AccountID"
	checkInvokeResponseFail(t, stub, args, expectedMessage)
}

func Test_queryAccountByName(t *testing.T) {
	cc := new(Chaincode)
	stub := shim.NewMockStub("query_acc_test", cc)

	// Init 1 account with 10 tokens
	checkInit(t, stub, [][]byte{[]byte("1"), []byte("10")})

	// It should return one account
	args := [][]byte{[]byte("queryAccountByName"), []byte("Init_Account")}
	expectedPayload := "[{\"RecordType\":\"ACCOUNT\",\"AccountID\":\"1\",\"Name\":\"Init_Account\",\"Tokens\":10}]"
	checkInvokeResponse(t, stub, args, expectedPayload)

	// It should fail with empty string arg
	args = [][]byte{[]byte("queryAccountByName"), []byte("")}
	expectedMessage := "1st argument must be a non-empty string"
	checkInvokeResponseFail(t, stub, args, expectedMessage)

	// It should fail with more than one args
	args = [][]byte{[]byte("queryAccountByName"), []byte("1"), []byte("lol")}
	expectedMessage = "Incorrect number of arguments. Expecting name of account holder"
	checkInvokeResponseFail(t, stub, args, expectedMessage)
}

func Test_transferTokens(t *testing.T) {
	cc := new(Chaincode)
	stub := shim.NewMockStub("tokens_init_test", cc)

	// Init 1 account with 10 000 tokens
	checkInit(t, stub, [][]byte{[]byte("1"), []byte("10000")})

	// create another acc without tokens
	checkCreateAccount(t, stub, "2", "marcel")
	// It should transfer tokens
	args := [][]byte{[]byte("transferTokens"), []byte("1"), []byte("2"), []byte("100")}
	checkInvoke(t, stub, args)

	// It should not transfer account that does have enough tokens

	args = [][]byte{[]byte("transferTokens"), []byte("2"), []byte("1"), []byte("101")}
	checkInvokeFail(t, stub, args)

	// It should transfer all tokens available in account
	args = [][]byte{[]byte("transferTokens"), []byte("2"), []byte("1"), []byte("100")}
	checkInvoke(t, stub, args)

	/*
		// These tests cannot be executed in Mock environment yet.
		// They create panic even though chaincode is ok and should return error.
		// It is because of the implementation of Mock state as map and it dereference 0 pointer
		// It should not transfer tokens to account that does not exist
		args = [][]byte{[]byte("transferTokens"), []byte("1"), []byte("3"), []byte("100")}
		checkInvokeFail(t, stub, args)

		// It should not transfer tokens from account that does not exist
		args = [][]byte{[]byte("transferTokens"), []byte("3"), []byte("1"), []byte("100")}
		checkInvokeFail(t, stub, args)

	*/
}

/*
// This test cannot be used because of the limitations of MockChaincode
func Test_getHistoryForAccount(t *testing.T) {
	cc := new(Chaincode)
	stub := shim.NewMockStub("tokens_init_test", cc)

	// Init 1 account with 10 000 tokens
	checkInit(t, stub, [][]byte{[]byte("1"), []byte("10000")})

	// create another acc without tokens
	checkCreateAccount(t, stub, "2", "marcel")
	// It should transfer tokens
	args := [][]byte{[]byte("transferTokens"), []byte("1"), []byte("2"), []byte("100")}
	checkInvoke(t, stub, args)

	args := [][]byte{[]byte("getAccHistory"), []byte("1")}
	checkInvoke(t, stub, args)
}
*/

func Test_getRecipientTx(t *testing.T) {
	cc := new(Chaincode)
	stub := shim.NewMockStub("tokens_init_test", cc)

	// Init 1 account with 10 000 tokens
	checkInit(t, stub, [][]byte{[]byte("1"), []byte("10000")})

	// create another acc without tokens
	checkCreateAccount(t, stub, "2", "marcel")
	// It should transfer tokens
	args := [][]byte{[]byte("transferTokens"), []byte("1"), []byte("2"), []byte("100")}
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.Fail()
	}
	// get the TxID and try to get recipient
	args = [][]byte{[]byte("getRecipientTx"), res.Payload}
	checkInvoke(t, stub, args)

	// It should fail with random TxID
	args = [][]byte{[]byte("getRecipientTx"), []byte("-4863asfaebh")}
	checkInvokeFail(t, stub, args)
}

func Test_addTxAsUsed(t *testing.T) {
	cc := new(Chaincode)
	stub := shim.NewMockStub("tokens_init_test", cc)

	// Init 1 account with 10 000 tokens
	checkInit(t, stub, [][]byte{[]byte("1"), []byte("10000")})

	// create another acc without tokens
	checkCreateAccount(t, stub, "2", "marcel")
	// It should transfer tokens
	args := [][]byte{[]byte("transferTokens"), []byte("1"), []byte("2"), []byte("100")}
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.Fail()
	}
	// get the TxID and add it as used TxID for purchase
	args = [][]byte{[]byte("addTxAsUsed"), res.Payload}
	checkInvoke(t, stub, args)

	// Do the same but now it should fail because it is already used TxID
	args = [][]byte{[]byte("addTxAsUsed"), res.Payload}
	checkInvokeFail(t, stub, args)
}