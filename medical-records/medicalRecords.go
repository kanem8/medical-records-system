package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/kanem8/fabric-samples/project2/medical-records/chaincode"
	// import for running with microFab:
	// "github.com/hyperledger/fabric-chaincode-go/shim"
)

func main() {

	medicalChaincode, err := contractapi.NewChaincode(&chaincode.MedicalSmartContract{})
	if err != nil {
		log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
	}

	if err := medicalChaincode.Start(); err != nil {
		log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
	}

	// // code for running with microFab:
	// assetChaincode, err := contractapi.NewChaincode(&chaincode.MedicalSmartContract{})
	// if err != nil {
	// 	log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
	// }

	// //The ccid is assigned to the chaincode on install (using the “peer lifecycle chaincode install <package>” command) for instance
	// ccid := "medical-records:43f23ec4b410d1c74e691379904c2ef02fa62fb90136af2c957ed4e4d12c87c3"

	// server := &shim.ChaincodeServer{
	// 	CCID:    ccid,
	// 	Address: "0.0.0.0:9999",
	// 	CC:      assetChaincode,
	// 	TLSProps: shim.TLSProperties{
	// 		Disabled: true,
	// 	},
	// }
	// err2 := server.Start()
	// if err2 != nil {
	// 	fmt.Printf("Error starting Simple chaincode: %s", err)
	// }

}
