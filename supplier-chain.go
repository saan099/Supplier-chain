package main

import (
  "errors"
  "encoding/json"
  "strconv"
  "fmt"
  "github.com/hyperledger/fabric/core/chaincode/shim"
)

type invoice struct{
  Order_id string `json:"order_id"`
  Product_name string `json:"product_name"`
  Quantity int `json:"quantity"`
  Total_payment int `json:"total_payment"`
  Delivery_date string `json:"delivery_date"`
  Interest float64 `json:"interest"`
}

type buyer struct {
  BuyerId string `json:"buyerId"`
  BuyerName string `json:"buyerName"`
  BuyerBalance float64 `json:"buyerBalance"`
  GoodsRecieved string `json:"goodsRecieved"` 
}

type supplier struct {
  SupplierId string `json:"supplierId"`
  SupplierName string `json:"supplierName"`
  SupplierBalance float64 `json:"supplierBalance"`
  GoodsDelivered string `json:"goodsDelivered"`
}
type bank struct {
  BankId string `json:"bankId"`
  BankName string `json:"bankName"`
  BankBalance float64 `json:"bankBalance"`
  LoanedAmount float64 `json:"loanedAmount"`  
}

var buyerNameKey string="buyerName"
var supplierNameKey string="supplierName"
var bankNameKey string="bankName"

type SupplierChaincode struct{
}

func main () {

  err := shim.Start(new(SupplierChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
func (t *SupplierChaincode) Init (stub shim.ChaincodeStubInterface,function string, args []string) ([]byte,error) {
  var err error
  if len(args)!=3 {
    return nil,errors.New("wrong number of arguments")
  }
  if len(args[0])<1||len(args[1])<1||len(args[2])<1 {
    return nil,errors.New("No value given to one field")
  }
	err=stub.PutState(buyerNameKey,[]byte(args[0]))

	err=stub.PutState(supplierNameKey,[]byte(args[1]))

	err=stub.PutState(bankNameKey,[]byte(args[2]))

  if err!=nil {
    return nil,errors.New("didnt commit node's names")
  }

  return nil, nil
}

func (t *SupplierChaincode) Invoke (stub shim.ChaincodeStubInterface, function string, args []string) ([]byte,error) {

  if function=="init" {
    return t.Init(stub,function,args)
  } else if function=="makeOrderDetailsinInvoice" {
    return t.MakeOrderinInvoice(stub,args)
  } else if function=="supplyDetailsinInvoice" {
    return t.SupplyDetailsinInvoice(stub, args)
  } else if function=="bankDetailsinInvoice" {
    return t.BankDetailsinInvoice(stub,args)
  }

  return nil,errors.New("No function invoked")

}

func (t *SupplierChaincode) MakeOrderinInvoice (stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
  var err error
  if len(args)!=3 {
    return nil, errors.New("number of arguments are wrong")
  }
  str:=`{"order_id": "`+args[0]+`", "product_name": "`+args[1]+`", "quantity": `+args[2]+`, "total_payment":`+strconv.Itoa(0)+`,"delivery_date":"`+`null`+`","interest":`+strconv.Itoa(0)+`}`
	err=stub.PutState(args[0],[]byte(str))
  if err!=nil {
    return nil,errors.New("error created in order committed")
  }

  return nil,nil
}
func (t *SupplierChaincode) SupplyDetailsinInvoice (stub shim.ChaincodeStubInterface,args []string) ([]byte,error) {
  var err error
  if len(args)!=3 {
    return nil, errors.New("number of arguemnts are wrong")
  }
  valAsbytes, err:=stub.GetState(args[0])
  if err!=nil {
    return nil, errors.New("state not found")
  }
  inv:=invoice{}
  err=json.Unmarshal(valAsbytes,&inv)
  if err!=nil {
    return nil, err
  }

  inv.Total_payment, _=strconv.Atoi(args[1])
  inv.Delivery_date=args[2]
  str:=`{"order_id": "`+inv.Order_id+`", "product_name": "`+inv.Product_name+`", "quantity": `+strconv.Itoa(inv.Quantity)+`, "total_payment"`+strconv.Itoa(inv.Total_payment)+`,"delivery_date":"`+inv.Delivery_date+`","interest":`+strconv.FormatFloat(inv.Interest, 'E', -1, 32)+`}`
  err=stub.PutState(inv.Order_id,[]byte(str))
  if err!=nil {
    return nil, errors.New("state not comitted")
  }
  return nil,nil

}
func (t *SupplierChaincode) BankDetailsinInvoice (stub shim.ChaincodeStubInterface,args []string) ([]byte,error) {

  var err error
  if len(args)!=2 {
    return nil,errors.New("number of arguments are wrong")
  }

  valAsbytes,err := stub.GetState(args[0])
  if err!=nil {
    return nil, err
  }
  inv:=invoice{}
  err=json.Unmarshal(valAsbytes,&inv)
  if err!=nil {
    return nil, err
  }
  inv.Interest,err =strconv.ParseFloat(args[1], 32)
  if err!=nil {
    return nil, err
  }
  str:=`{"order_id": "`+inv.Order_id+`", "product_name": "`+inv.Product_name+`", "quantity": `+strconv.Itoa(inv.Quantity)+`, "total_payment":`+strconv.Itoa(inv.Total_payment)+`,"delivery_date":"`+inv.Delivery_date+`","interest":`+strconv.FormatFloat(inv.Interest, 'E', -1, 32)+`}`
  err=inv.PutState(args[0],[]byte(str))
  if err!=nil {
    return nil, err
  }
  return nil, nil
}



func (t *SupplierChaincode) Query (stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

  if function=="readOrder" {
    return t.ReadOrder(stub,args)
  }

  return nil,errors.New("quey didnt meet any function")

}

func (t *SupplierChaincode) ReadOrder (stub shim.ChaincodeStubInterface,args []string) ([]byte, error) {

  if len(args)!=1 {
    return nil, errors.New("Wrong numer of arguments")
  }

  valAsbytes,err := stub.GetState(args[0])
  if err!=nil {
    return nil, err
  }
   return valAsbytes, nil

}

