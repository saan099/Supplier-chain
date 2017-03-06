package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type invoice struct {
	Order_id      string  `json:"order_id"`
	Product_name  string  `json:"product_name"`
	Quantity      int     `json:"quantity"`
	Total_payment int     `json:"total_payment"`
	Delivery_date string  `json:"delivery_date"`
	Interest      float64 `json:"interest"`
}

type buyer struct {
	BuyerId       string    `json:"buyerId"`
	BuyerName     string    `json:"buyerName"`
	BuyerBalance  float64   `json:"buyerBalance"`
	GoodsRecieved []invoice `json:"goodsRecieved"`
}

type supplier struct {
	SupplierId      string    `json:"supplierId"`
	SupplierName    string    `json:"supplierName"`
	SupplierBalance float64   `json:"supplierBalance"`
	GoodsDelivered  []invoice `json:"goodsDelivered"`
}
type bank struct {
	BankId       string  `json:"bankId"`
	BankName     string  `json:"bankName"`
	BankBalance  float64 `json:"bankBalance"`
	LoanedAmount float64 `json:"loanedAmount"`
}

var buyerNameKey string = "buyerName"
var supplierNameKey string = "supplierName"
var bankNameKey string = "bankName"

type SupplierChaincode struct {
}

func main() {

	err := shim.Start(new(SupplierChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
func (t *SupplierChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error
	if len(args) != 3 {
		return nil, errors.New("wrong number of arguments")
	}
	if len(args[0]) < 1 || len(args[1]) < 1 || len(args[2]) < 1 {
		return nil, errors.New("No value given to one field")
	}
	err = stub.PutState(buyerNameKey, []byte(args[0]))

	err = stub.PutState(supplierNameKey, []byte(args[1]))

	err = stub.PutState(bankNameKey, []byte(args[2]))

	if err != nil {
		return nil, errors.New("didnt commit node's names")
	}

	return nil, nil
}

func (t *SupplierChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "init" {
		return t.Init(stub, function, args)
	} else if function == "makeOrderDetailsinInvoice" {
		return t.MakeOrderinInvoice(stub, args)
	} else if function == "supplyDetailsinInvoice" {
		return t.SupplyDetailsinInvoice(stub, args)
	} else if function == "bankDetailsinInvoice" {
		return t.BankDetailsinInvoice(stub, args)
	} else if function == "initializeBuyer" {
		return t.InitializeBuyer(stub, args)
	} else if function == "initializeSupplier" {
		return t.InitializeSupplier(stub, args)
	} else if function == "initializeBank" {
		return t.InitializeBank(stub, args)
	} else if function == "addBalanceinBank" {
		return t.addBalanceinBank(stub, args)
	} else if function == "addBalanceinSupplier" {
		return t.addBalanceinSupplier(stub, args)
	} else if function == "addBalanceinBuyer" {
		return t.addBalanceinBuyer(stub, args)
	} else if function == "deliverGoods" {
		return t.DeliverGoods(stub, args)
	}

	return nil, errors.New("No function invoked")

}

func (t *SupplierChaincode) MakeOrderinInvoice(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 3 {
		return nil, errors.New("number of arguments are wrong")
	}
	str := `{"order_id": "` + args[0] + `", "product_name": "` + args[1] + `", "quantity": ` + args[2] + `, "total_payment":` + strconv.Itoa(0) + `,"delivery_date":"` + `null` + `","interest":` + strconv.FormatFloat(0.0, 'f', -1, 32) + `}`
	err = stub.PutState(args[0], []byte(str))
	if err != nil {
		return nil, errors.New("error created in order committed")
	}

	return nil, nil
}
func (t *SupplierChaincode) SupplyDetailsinInvoice(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 3 {
		return nil, errors.New("number of arguemnts are wrong")
	}
	valAsbytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("state not found")
	}
	inv := invoice{}
	err = json.Unmarshal(valAsbytes, &inv)
	if err != nil {
		return nil, err
	}

	inv.Total_payment, _ = strconv.Atoi(args[1])
	inv.Delivery_date = args[2]
	str := `{"order_id": "` + inv.Order_id + `", "product_name": "` + inv.Product_name + `", "quantity": ` + strconv.Itoa(inv.Quantity) + `, "total_payment":` + strconv.Itoa(inv.Total_payment) + `,"delivery_date":"` + inv.Delivery_date + `","interest":` + strconv.FormatFloat(inv.Interest, 'f', -1, 32) + `}`
	err = stub.PutState(inv.Order_id, []byte(str))
	if err != nil {
		return nil, errors.New("state not comitted")
	}
	return nil, nil

}
func (t *SupplierChaincode) BankDetailsinInvoice(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	if len(args) != 2 {
		return nil, errors.New("number of arguments are wrong")
	}

	valAsbytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, err
	}
	inv := invoice{}
	err = json.Unmarshal(valAsbytes, &inv)
	if err != nil {
		return nil, err
	}

	inv.Interest, _ = strconv.ParseFloat(args[1], 64)

	str := `{"order_id": "` + inv.Order_id + `", "product_name": "` + inv.Product_name + `", "quantity": ` + strconv.Itoa(inv.Quantity) + `, "total_payment":` + strconv.Itoa(inv.Total_payment) + `,"delivery_date":"` + inv.Delivery_date + `","interest":` + strconv.FormatFloat(inv.Interest, 'f', -1, 32) + `}`
	err = stub.PutState(args[0], []byte(str))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SupplierChaincode) InitializeBuyer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 3 {
		return nil, errors.New("wrong number of arguments")
	}

	var acc = buyer{}
	acc.BuyerId = args[0]
	acc.BuyerName = args[1]
	acc.BuyerBalance, _ = strconv.ParseFloat(args[2], 64)
	var goodsrecieved []invoice
	acc.GoodsRecieved = goodsrecieved
	jsonAsbytes, _ := json.Marshal(acc)

	//str := `{"buyerId":"` + args[0] + `","buyerName":"` + args[1] + `","buyerBalance":` + args[2] + `,"goodsRecieved":"` + string(goodsAsbytes[:]) + `"}`

	err = stub.PutState(args[0], jsonAsbytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SupplierChaincode) InitializeSupplier(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	if len(args) != 3 {
		return nil, errors.New("wrong number of arguments")
	}

	var acc = supplier{}
	acc.SupplierId = args[0]
	acc.SupplierName = args[1]
	acc.SupplierBalance, _ = strconv.ParseFloat(args[2], 64)

	var goodsDelivered []invoice
	acc.GoodsDelivered = goodsDelivered
	jsonAsbytes, err := json.Marshal(acc)

	err = stub.PutState(args[0], jsonAsbytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
func (t *SupplierChaincode) InitializeBank(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	if len(args) != 3 {
		return nil, errors.New("wrong number of arguments")
	}
	str := `{"bankId":"` + args[0] + `","bankName":"` + args[1] + `","bankBalance":` + args[2] + `,"loanedAmount":0}`
	err = stub.PutState(args[0], []byte(str))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SupplierChaincode) addBalanceinBank(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 2 {
		return nil, errors.New("wrong number of arguments")
	}
	valAsbytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, err
	}
	acc := bank{}
	err = json.Unmarshal(valAsbytes, &acc)
	if err != nil {
		return nil, err
	}
	addedAmout, _ := strconv.ParseFloat(args[1], 64)
	acc.BankBalance += addedAmout
	str := `{"bankId":"` + acc.BankId + `","bankName":"` + acc.BankName + `","bankBalance":` + strconv.FormatFloat(acc.BankBalance, 'f', -1, 32) + `,"loanedAmount":` + strconv.FormatFloat(acc.LoanedAmount, 'f', -1, 32) + `}`
	err = stub.PutState(args[0], []byte(str))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SupplierChaincode) addBalanceinBuyer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 2 {
		return nil, errors.New("wrong number of arguments")
	}
	valAsbytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, err
	}
	acc := buyer{}
	err = json.Unmarshal(valAsbytes, &acc)
	if err != nil {
		return nil, err
	}
	addedAmount, _ := strconv.ParseFloat(args[1], 64)
	acc.BuyerBalance += addedAmount
	jsonAsbytes, err := json.Marshal(acc)
	//str := `{"buyerId":"` + acc.BuyerId + `","buyerName":"` + acc.BuyerName + `","buyerBalance":` + strconv.FormatFloat(acc.BuyerBalance, 'f', -1, 32) + `,"goodsRecieved":"` + acc.GoodsRecieved + `"}`
	err = stub.PutState(args[0], jsonAsbytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SupplierChaincode) addBalanceinSupplier(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 2 {
		return nil, errors.New("wrong number of arguments")
	}
	valAsbytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, err
	}
	acc := supplier{}
	err = json.Unmarshal(valAsbytes, &acc)
	if err != nil {
		return nil, err
	}
	addedAmout, _ := strconv.ParseFloat(args[1], 64)
	acc.SupplierBalance += addedAmout
	jsonAsbytes, _ := json.Marshal(acc)

	err = stub.PutState(args[0], jsonAsbytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SupplierChaincode) DeliverGoods(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	if len(args) != 2 {
		return nil, errors.New("wrong number of arguments")
	}
	supplierId := args[0]
	orderId := args[1]

	supplierDetailsAsbytes, err := stub.GetState(supplierId)
	if err != nil {
		return nil, err
	}
	orderDetailsAsbytes, err := stub.GetState(orderId)
	if err != nil {
		return nil, err
	}
	order := invoice{}

	err = json.Unmarshal(orderDetailsAsbytes, &order)
	if err != nil {
		return nil, err
	}
	acc := supplier{}

	err = json.Unmarshal(supplierDetailsAsbytes, &acc)
	if err != nil {
		return nil, err
	}
	acc.GoodsDelivered = append(acc.GoodsDelivered, order)

	jsonAsbytes, err := json.Marshal(acc)
	err = stub.PutState(supplierId, jsonAsbytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SupplierChaincode) RecieveGoods(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	if len(args) != 2 {
		return nil, errors.New("wrong number of arguments")
	}
	buyerId := args[0]
	orderId := args[1]

	buyerDetailsAsbytes, err := stub.GetState(buyerId)
	if err != nil {
		return nil, err
	}
	orderDetailsAsbytes, err := stub.GetState(orderId)
	if err != nil {
		return nil, err
	}
	order := invoice{}

	err = json.Unmarshal(orderDetailsAsbytes, &order)
	if err != nil {
		return nil, err
	}
	acc := buyer{}

	err = json.Unmarshal(buyerDetailsAsbytes, &acc)
	if err != nil {
		return nil, err
	}
	acc.GoodsRecieved = append(acc.GoodsRecieved, order)

	jsonAsbytes, err := json.Marshal(acc)
	err = stub.PutState(buyerId, jsonAsbytes)
	if err != nil {
		return nil, err
	}
	return nil, nil

}

//Queries---
func (t *SupplierChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "read" {
		return t.Read(stub, args)
	}

	return nil, errors.New("quey didnt meet any function")

}

func (t *SupplierChaincode) Read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Wrong numer of arguments")
	}

	valAsbytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, err
	}
	return valAsbytes, nil

}
