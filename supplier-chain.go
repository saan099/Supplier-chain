package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//order for a product
type order struct {
	Order_id      string  `json:"order_id"`
	Product_name  string  `json:"product_name"`
	Quantity      int     `json:"quantity"`
	Total_payment float64 `json:"total_payment"`
	Delivery_date string  `json:"delivery_date"`
	Status        string  `json:"status"`
}

//account for buyer
type buyer struct {
	BuyerId       string  `json:"buyerId"`
	BuyerName     string  `json:"buyerName"`
	BuyerBalance  float64 `json:"buyerBalance"`
	GoodsRecieved []order `json:"goodsRecieved"`
	PendingPay    float64 `json:"pendingPay"`
}

//account for supplier
type supplier struct {
	SupplierId      string  `json:"supplierId"`
	SupplierName    string  `json:"supplierName"`
	SupplierBalance float64 `json:"supplierBalance"`
	GoodsDelivered  []order `json:"goodsDelivered"`
	Loans           []loans `json:"loans"`
}

//account for bank
type bank struct {
	BankId         string  `json:"bankId"`
	BankName       string  `json:"bankName"`
	BankBalance    float64 `json:"bankBalance"`
	Loans          []loans `json:"loans"`
	LoanedAmount   float64 `json:"loanedAmount"`
	AmountRecieved float64 `json:"amountRecieved"`
}

// loan details
type loans struct {
	LoanId     string  `json:"loanId"`
	Orders     []order `json:"orders"`
	LoanAmount float64 `json:"loanAmount"`
	Status     string  `json:"status"`
	Interest   float64 `json:"interest"`
}

var buyerNameKey string = "buyerName"
var supplierNameKey string = "supplierName"
var bankNameKey string = "bankName"
var orderIndex string = "OrderIndex"

type SupplierChaincode struct {
}

//launching chaincode
func main() {

	err := shim.Start(new(SupplierChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

//initializing chaincode
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

	var list []string
	jsonAsbytes, _ := json.Marshal(list)
	err = stub.PutState(orderIndex, jsonAsbytes)
	if err != nil {
		return nil, errors.New("didnt commit node's names")
	}

	return nil, nil
}

//function to invoke all functionality of chaoncode
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
	} else if function == "recieveGoods" {
		return t.RecieveGoods(stub, args)
	} else if function == "generateInvoice" {
		return t.invoiceGeneration(stub, args)
	} else if function == "loanAmount" {
		return t.LoanAmount(stub, args)
	} else if function == "payToBank" {
		return t.PayToBank(stub, args)
	} else if function == "sendProfitsToSupplier" {
		return t.SendProfitsToSupplier(stub, args)
	}

	return nil, errors.New("No function invoked")

}

//make order for product from buyer
func (t *SupplierChaincode) MakeOrderinInvoice(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 3 {
		return nil, errors.New("number of arguments are wrong")
	}
	var o = order{}
	o.Order_id = args[0]
	o.Product_name = args[1]
	o.Quantity, _ = strconv.Atoi(args[2])
	o.Total_payment, _ = strconv.ParseFloat(`0`, 64)
	o.Delivery_date = "null"
	o.Status = "pending"
	//str := `{"order_id": "` + args[0] + `", "product_name": "` + args[1] + `", "quantity": ` + args[2] + `, "total_payment":` + strconv.ParseFloat(`0`, 64) + `,"delivery_date":"` + `null` + `","status":"` + `pending` + `"}`
	orderAsbytes, err := json.Marshal(o)
	err = stub.PutState(args[0], orderAsbytes)
	if err != nil {
		return nil, errors.New("error created in order committed")
	}
	var list []string
	valAsbytes, err := stub.GetState(orderIndex)
	err = json.Unmarshal(valAsbytes, &list)
	list = append(list, args[0])
	jsonAsbytes, _ := json.Marshal(list)
	err = stub.PutState(orderIndex, jsonAsbytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//enter supply details for the order by supplier
func (t *SupplierChaincode) SupplyDetailsinInvoice(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 3 {
		return nil, errors.New("number of arguemnts are wrong")
	}
	valAsbytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("state not found")
	}
	inv := order{}
	err = json.Unmarshal(valAsbytes, &inv)
	if err != nil {
		return nil, err
	}

	inv.Total_payment, _ = strconv.ParseFloat(args[1], 64)
	inv.Delivery_date = args[2]
	inv.Status = "processing"
	jsonAsbytes, _ := json.Marshal(inv)
	//str := `{"order_id": "` + inv.Order_id + `", "product_name": "` + inv.Product_name + `", "quantity": ` + strconv.Itoa(inv.Quantity) + `, "total_payment":` + strconv.ParseFloat(`0`, 64) + `,"delivery_date":"` + inv.Delivery_date + `","status":"` + `processing` + `"}`
	err = stub.PutState(inv.Order_id, jsonAsbytes)
	if err != nil {
		return nil, errors.New("state not committed")
	}
	return nil, nil

}

//entering banking details in order by bank
func (t *SupplierChaincode) BankDetailsinInvoice(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	/*
		var err error
		if len(args) != 2 {
			return nil, errors.New("number of arguments are wrong")
		}

		valAsbytes, err := stub.GetState(args[0])
		if err != nil {
			return nil, err
		}
		inv := order{}
		err = json.Unmarshal(valAsbytes, &inv)
		if err != nil {
			return nil, err
		}

		inv.Interest, _ = strconv.ParseFloat(args[1], 64)

		str := `{"order_id": "` + inv.Order_id + `", "product_name": "` + inv.Product_name + `", "quantity": ` + strconv.Itoa(inv.Quantity) + `, "total_payment":` + strconv.Itoa(inv.Total_payment) + `,"delivery_date":"` + inv.Delivery_date + `","interest":` + strconv.FormatFloat(inv.Interest, 'f', -1, 32) + `}`
		err = stub.PutState(args[0], []byte(str))
		if err != nil {
			return nil, err
		}*/
	return nil, nil
}

//initialize buyer's account
func (t *SupplierChaincode) InitializeBuyer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 3 {
		return nil, errors.New("wrong number of arguments")
	}

	var acc = buyer{}
	acc.BuyerId = args[0]
	acc.BuyerName = args[1]
	acc.BuyerBalance, _ = strconv.ParseFloat(args[2], 64)
	var goodsrecieved []order
	acc.GoodsRecieved = goodsrecieved
	acc.PendingPay = 0
	jsonAsbytes, _ := json.Marshal(acc)

	//str := `{"buyerId":"` + args[0] + `","buyerName":"` + args[1] + `","buyerBalance":` + args[2] + `,"goodsRecieved":"` + string(goodsAsbytes[:]) + `"}`

	err = stub.PutState(args[0], jsonAsbytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//initialize supplier's account
func (t *SupplierChaincode) InitializeSupplier(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	if len(args) != 3 {
		return nil, errors.New("wrong number of arguments")
	}

	var acc = supplier{}
	acc.SupplierId = args[0]
	acc.SupplierName = args[1]
	acc.SupplierBalance, _ = strconv.ParseFloat(args[2], 64)
	var l []loans
	acc.Loans = l

	var goodsDelivered []order
	acc.GoodsDelivered = goodsDelivered
	jsonAsbytes, err := json.Marshal(acc)

	err = stub.PutState(args[0], jsonAsbytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//initialize banker's account
func (t *SupplierChaincode) InitializeBank(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	if len(args) != 3 {
		return nil, errors.New("wrong number of arguments")
	}
	//str := `{"bankId":"` + args[0] + `","bankName":"` + args[1] + `","bankBalance":` + args[2] + `,"loanedAmount":0}`
	var l []loans
	acc := bank{}
	acc.BankId = args[0]
	acc.BankName = args[1]
	acc.BankBalance, _ = strconv.ParseFloat(args[2], 64)
	acc.Loans = l
	acc.AmountRecieved = 0
	jsonAsbytes, _ := json.Marshal(acc)
	err = stub.PutState(args[0], jsonAsbytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//add balance in banker's account
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
	jsonAsbytes, _ := json.Marshal(acc)
	//str := `{"bankId":"` + acc.BankId + `","bankName":"` + acc.BankName + `","bankBalance":` + strconv.FormatFloat(acc.BankBalance, 'f', -1, 32) + `,"loanedAmount":` + strconv.FormatFloat(acc.LoanedAmount, 'f', -1, 32) + `}`
	err = stub.PutState(args[0], jsonAsbytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//adding balance in buyer's account
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

//adding balance in supplier's account
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

//delivering goods to buyer by supplier
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
	o := order{}

	err = json.Unmarshal(orderDetailsAsbytes, &o)
	if err != nil {
		return nil, err
	}
	o.Status = `delivered`
	orderAsjson, _ := json.Marshal(o)
	err = stub.PutState(orderId, orderAsjson)

	acc := supplier{}

	err = json.Unmarshal(supplierDetailsAsbytes, &acc)
	if err != nil {
		return nil, err
	}
	acc.GoodsDelivered = append(acc.GoodsDelivered, o)

	jsonAsbytes, err := json.Marshal(acc)
	err = stub.PutState(supplierId, jsonAsbytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//recieving goods from supplier by buyer
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
	o := order{}

	err = json.Unmarshal(orderDetailsAsbytes, &o)
	if err != nil {
		return nil, err
	}
	o.Status = `recieved`
	orderAsjson, _ := json.Marshal(o)
	err = stub.PutState(orderId, orderAsjson)

	acc := buyer{}

	err = json.Unmarshal(buyerDetailsAsbytes, &acc)
	if err != nil {
		return nil, err
	}
	acc.GoodsRecieved = append(acc.GoodsRecieved, o)
	acc.PendingPay += o.Total_payment
	jsonAsbytes, err := json.Marshal(acc)
	err = stub.PutState(buyerId, jsonAsbytes)
	if err != nil {
		return nil, err
	}
	return nil, nil

}

//loan request by supplier by providing an invoice
func (t *SupplierChaincode) invoiceGeneration(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	var i int
	if len(args) < 4 {
		return nil, errors.New("wrong number of arguments")
	}
	loan := loans{}

	numberOfOrders, err := strconv.Atoi(args[0])

	for i = 0; i < numberOfOrders; i++ {
		o := order{}
		valAsbytes, er := stub.GetState(args[i+1])
		if er != nil {
			return nil, er
		}
		err = json.Unmarshal(valAsbytes, o)
		loan.Orders = append(loan.Orders, o)
	}
	i = i + 1
	loan.LoanId = args[i]

	i = i + 1
	loan.LoanAmount, _ = strconv.ParseFloat(args[i], 64)
	loan.Status = "pending"
	loan.Interest = 0
	bankAcc := bank{}
	i = i + 1
	valueAsbytes, err := stub.GetState(args[i])
	err = json.Unmarshal(valueAsbytes, &bankAcc)
	bankAcc.Loans = append(bankAcc.Loans, loan)
	valAsbytes, err := json.Marshal(bankAcc)
	err = stub.PutState(args[i], valAsbytes)
	if err != nil {
		return nil, err
	}
	i = i + 1
	supplierAcc := supplier{}
	supplierAsbytes, err := stub.GetState(args[i])
	err = json.Unmarshal(supplierAsbytes, &supplierAcc)
	supplierAcc.Loans = append(supplierAcc.Loans, loan)
	jsonAsbytes, err := json.Marshal(supplierAcc)
	err = stub.PutState(args[i], jsonAsbytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

//accepting loan request from supplier and aloan transaction
func (t *SupplierChaincode) LoanAmount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 4 {
		return nil, errors.New("wrong number of arguments")
	}
	var temp float64
	supplierId := args[0]
	bankId := args[1]
	loanId := args[2]
	bankAcc := bank{}
	bankAsbytes, err := stub.GetState(bankId)
	err = json.Unmarshal(bankAsbytes, &bankAcc)
	if err != nil {
		return nil, err
	}
	for i := range bankAcc.Loans {
		if bankAcc.Loans[i].LoanId == loanId {
			bankAcc.LoanedAmount += bankAcc.Loans[i].LoanAmount
			bankAcc.Loans[i].Status = "accepted"
			bankAcc.Loans[i].Interest, _ = strconv.ParseFloat(args[3], 64)
			temp = bankAcc.Loans[i].LoanAmount
			bankAcc.BankBalance -= temp
			jsonAsbytes, err := json.Marshal(bankAcc)
			err = stub.PutState(bankId, jsonAsbytes)
			if err != nil {
				return nil, err
			}
			supplierAcc := supplier{}
			supplierAsbytes, err := stub.GetState(supplierId)
			err = json.Unmarshal(supplierAsbytes, &supplierAcc)
			for i := range supplierAcc.Loans {
				if supplierAcc.Loans[i].LoanId == loanId {
					supplierAcc.Loans[i].Status = "accepted"
					supplierAcc.Loans[i].Interest, _ = strconv.ParseFloat(args[3], 64)
				}
			}

			supplierAcc.SupplierBalance += temp
			jsonAsbyte, err := json.Marshal(supplierAcc)
			err = stub.PutState(supplierId, jsonAsbyte)
			if err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

//paying the required amount to bank set by supplier including supplier's profit
func (t *SupplierChaincode) PayToBank(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 3 {
		return nil, errors.New("wrong number of arguments")
	}
	buyerId := args[0]
	bankId := args[1]
	amountToPay, _ := strconv.ParseFloat(args[2], 64)

	buyerAsbytes, err := stub.GetState(buyerId)
	bankAsbytes, err := stub.GetState(bankId)

	buyerAcc := buyer{}
	bankAcc := bank{}
	err = json.Unmarshal(buyerAsbytes, &buyerAcc)
	err = json.Unmarshal(bankAsbytes, &bankAcc)
	buyerAcc.PendingPay -= amountToPay
	buyerAcc.BuyerBalance -= amountToPay
	bankAcc.BankBalance += amountToPay
	bankAcc.AmountRecieved += amountToPay
	buyerjson, _ := json.Marshal(buyerAcc)
	bankjson, _ := json.Marshal(bankAcc)

	err = stub.PutState(buyerId, buyerjson)
	if err != nil {
		return nil, err
	}
	err = stub.PutState(bankId, bankjson)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//separating the interested amount from amount recieved from buyer and passing profits to supplier
func (t *SupplierChaincode) SendProfitsToSupplier(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("wrong number of arguments")
	}
	supplierAcc := supplier{}
	supplierAsbytes, err := stub.GetState(args[0])
	err = json.Unmarshal(supplierAsbytes, &supplierAcc)
	if err != nil {
		return nil, err
	}
	bankAcc := bank{}
	bankAsbytes, err := stub.GetState(args[1])
	err = json.Unmarshal(bankAsbytes, &bankAcc)
	temp := bankAcc.LoanedAmount + (bankAcc.LoanedAmount*bankAcc.Loans[0].Interest)/100
	if temp < bankAcc.AmountRecieved {
		bankAcc.BankBalance -= bankAcc.AmountRecieved - temp
		supplierAcc.SupplierBalance += bankAcc.AmountRecieved - temp
		supplierAcc.Loans[0].Status = "completed"
		bankAcc.Loans[0].Status = "completed"
		bankAcc.AmountRecieved -= temp
		jsonBank, _ := json.Marshal(bankAcc)
		jsonSupplier, _ := json.Marshal(supplierAcc)

		err = stub.PutState(args[1], jsonBank)
		err = stub.PutState(args[0], jsonSupplier)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil

}

//Queries---
func (t *SupplierChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "read" {
		return t.Read(stub, args)
	} else if function == "readAllOrders" {
		return t.ReadAllOrders(stub, args)
	}

	return nil, errors.New("quey didnt meet any function")

}

//reading any required state/account
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

//reading all orders made by buyer
func (t *SupplierChaincode) ReadAllOrders(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, errors.New("wrong number of arguments")
	}
	var list []string
	var orders []order
	valAsbytes, err := stub.GetState(orderIndex)
	err = json.Unmarshal(valAsbytes, &list)
	for i := range list {
		orderAsbytes, err := stub.GetState(list[i])
		if err != nil {
			return nil, err
		}
		o := order{}
		json.Unmarshal(orderAsbytes, &o)
		orders = append(orders, o)

	}
	jsonAsbytes, err := json.Marshal(orders)
	if err != nil {
		return nil, err
	}
	return jsonAsbytes, nil

}
