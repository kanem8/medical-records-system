package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Define smart contract.
type MedicalSmartContract struct {
	contractapi.Contract
}

// PatientRecord struct
type PatientRecord struct {
	SSN     string   `json:"ssn"`
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Gender  string   `json:"gender"`
	Records []Record `json:"records"`
}

// Record struct
type Record struct {
	Date      string `json:"date"`
	Doctor    string `json:"doctor"`
	Hospital  string `json:"hospital"`
	Diagnosis string `json:"diagnosis"`
	Treatment string `json:"treatment"`
}

// InitLedger adds a base set of patient records to the ledger
// This won't be exposed by application, it will be invoked manually for demo purposes
func (s *MedicalSmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	// Define the patient records to add to the ledger
	patients := []PatientRecord{
		{
			SSN:    "123-45-6789",
			Name:   "Alice",
			Age:    30,
			Gender: "Female",
			Records: []Record{
				{
					Date:      "2022-01-01",
					Doctor:    "Dr. Smith",
					Hospital:  "General Hospital",
					Diagnosis: "Flu",
					Treatment: "Rest and fluids",
				},
				{
					Date:      "2022-02-15",
					Doctor:    "Dr. Johnson",
					Hospital:  "City Hospital",
					Diagnosis: "Broken leg",
					Treatment: "Cast",
				},
			},
		},
		{
			SSN:    "234-56-7890",
			Name:   "Bob",
			Age:    45,
			Gender: "Male",
			Records: []Record{
				{
					Date:      "2021-11-01",
					Doctor:    "Dr. Lee",
					Hospital:  "General Hospital",
					Diagnosis: "Headache",
					Treatment: "Ibuprofen",
				},
				{
					Date:      "2022-03-10",
					Doctor:    "Dr. Johnson",
					Hospital:  "City Hospital",
					Diagnosis: "Appendicitis",
					Treatment: "Surgery",
				},
			},
		},
	}

	// Add the patient records to the ledger
	for _, patient := range patients {
		fmt.Println("Adding patient: ", patient.Name)
		err := s.CreatePatientRecord(ctx, &patient)
		if err != nil {
			return fmt.Errorf("failed to create patient record: %v", err)
		}
	}

	return nil
}

// StorePatientRecord stores a new patient record in the ledger.
// As with InitLedger func, this won't be exposed by application, it is only invoked by InitLedger for demo purposes
func (s *MedicalSmartContract) CreatePatientRecord(ctx contractapi.TransactionContextInterface, patientRecord *PatientRecord) error {

	fmt.Println("Creating patient with records: ", patientRecord.Name)

	// Convert patient record to JSON bytes
	patientRecordBytes, err := json.Marshal(patientRecord)
	if err != nil {
		return fmt.Errorf("failed to marshal patient record: %v", err)
	}

	// Put patient record to the ledger
	err = ctx.GetStub().PutState(patientRecord.SSN, patientRecordBytes)
	if err != nil {
		return fmt.Errorf("failed to put patient record to the ledger: %v", err)
	}

	return nil
}

// StorePatientRecord stores a new patient record in the ledger.
func (s *MedicalSmartContract) RegisterNewPatient(ctx contractapi.TransactionContextInterface, ssn string, name string, age int, gender string) error {

	// Define the patient records to add to the ledger
	patientRecord := PatientRecord{
		SSN:     ssn,
		Name:    name,
		Age:     age,
		Gender:  gender,
		Records: []Record{},
	}

	fmt.Println("Registering new patient: ", patientRecord.Name)

	// Convert patient record to JSON bytes
	patientRecordBytes, err := json.Marshal(patientRecord)
	if err != nil {
		return fmt.Errorf("failed to marshal patient record: %v", err)
	}

	// Put patient record to the ledger
	err = ctx.GetStub().PutState(patientRecord.SSN, patientRecordBytes)
	if err != nil {
		return fmt.Errorf("failed to put patient record to the ledger: %v", err)
	}

	fmt.Println("New patient ", patientRecord.Name, " is registered by their SSN: ", patientRecord.SSN)

	return nil
}

// Add a new medical record to an existing patient
func (s *MedicalSmartContract) AddPatientRecord(ctx contractapi.TransactionContextInterface,
	ssn string,
	date string,
	doctor string,
	hospital string,
	diagnosis string,
	treatment string) error {

	// Create a new record
	record := Record{
		Date:      date,
		Doctor:    doctor,
		Hospital:  hospital,
		Diagnosis: diagnosis,
		Treatment: treatment,
	}

	// retrieve the patient struct by their ssn
	patientRecordBytes, err := ctx.GetStub().GetState(ssn)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if patientRecordBytes == nil {
		return fmt.Errorf("the patient record %s does not exist. Please register the patient first", ssn)
	}

	// Unmarshal patient record from JSON bytes
	var patientRecord PatientRecord
	err = json.Unmarshal(patientRecordBytes, &patientRecord)
	if err != nil {
		return fmt.Errorf("failed to unmarshal patient record: %v", err)
	}

	// Append the new record to the patient's records slice
	patientRecord.Records = append(patientRecord.Records, record)

	// Marshal the updated patient record into JSON bytes
	patientRecordBytes, err = json.Marshal(patientRecord)
	if err != nil {
		return fmt.Errorf("failed to marshal patient record: %v", err)
	}

	// Update the patient record in the ledger
	err = ctx.GetStub().PutState(ssn, patientRecordBytes)
	if err != nil {
		return fmt.Errorf("failed to update patient record: %v", err)
	}

	return nil
}

// GetAllMedicalRecords returns all patient records stored on the ledger.
func (cc *MedicalSmartContract) GetAllMedicalRecords(ctx contractapi.TransactionContextInterface) ([]*PatientRecord, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var patientRecords []*PatientRecord
	for resultsIterator.HasNext() {
		recordBytes, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var patientRecord PatientRecord
		err = json.Unmarshal(recordBytes.Value, &patientRecord)
		if err != nil {
			return nil, err
		}

		patientRecords = append(patientRecords, &patientRecord)
	}

	return patientRecords, nil
}

// GetPatientRecordBySSN retrieves a patient record from the ledger using ssn.
func (s *MedicalSmartContract) GetPatientRecordBySSN(ctx contractapi.TransactionContextInterface, ssn string) (*PatientRecord, error) {
	patientRecordBytes, err := ctx.GetStub().GetState(ssn)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if patientRecordBytes == nil {
		return nil, fmt.Errorf("the patient record %s does not exist", ssn)
	}

	// Unmarshal patient record from JSON bytes
	var patientRecord PatientRecord
	err = json.Unmarshal(patientRecordBytes, &patientRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal patient record: %v", err)
	}

	return &patientRecord, nil
}

// // GetPatientRecordByName retrieves a patient record from the ledger using name.
// func (s *MedicalSmartContract) GetPatientRecordByName(ctx contractapi.TransactionContextInterface, name string) (*PatientRecord, error) {
//     // Get all patient records from the ledger
//     patientRecordIterator, err := ctx.GetStub().GetStateByRange("", "")
//     if err != nil {
//         return nil, fmt.Errorf("failed to get patient records from the ledger: %v", err)
//     }
//     defer patientRecordIterator.Close()

//     // Iterate through all patient records to find the one with matching name
//     for patientRecordIterator.HasNext() {
//         queryResponse, err := patientRecordIterator.Next()
//         if err != nil {
//             return nil, fmt.Errorf("failed to get next patient record: %v", err)
//         }

//         // Unmarshal patient record from JSON bytes
//         var patientRecord PatientRecord
//         err = json.Unmarshal(queryResponse.Value, &patientRecord)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to unmarshal patient record from JSON bytes: %v", err)
// 		}

// 		// Check if the patient record matches the name provided
// 		if patientRecord.Name == name {
// 			return &patientRecord, nil
// 		}

// 		return nil, fmt.Errorf("patient record not found for name: %s", name)
// 	}
// }
