# Medical Records Application

## Link to video demo
https://drive.google.com/file/d/1zEaAOsFUuc_LSPfxpfoIPRF9D4enMdNV/view?usp=share_link

## Prerequisites

The best thing to do to run this code is first follow the guide to download the [prerequisite software](https://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html) and the [Fabric binaries & Fabric samples](https://hyperledger-fabric.readthedocs.io/en/latest/install.html). Once this is done, you should have
the hyperledger fabric-samples repo cloned locally. Clone this repo somewhere in your local machine. Because the rest client relies on the test-network directory of the fabric-samples repo (and uses hardcoded relative paths to access files in there), please copy the project2 folder to the top level of the fabric-samples directory. 

  
## Usage

### First build the chaincode
- cd into the medical-records directory of the project2 repo. 
- Run the following commands to place external dependencies into "vendor" folder: `go mod tidy`, `go mod vendor`
- Setup fabric test network and deploy the medical-records chaincode by [following this instructions](https://hyperledger-fabric.readthedocs.io/en/release-2.4/test_network.html). I have summarised the necessary commands below:

  - `cd fabric-samples/test-network`
  - `./network.sh up` (if you ran this before, first run `./network.sh down` to tear down previous containers). This starts a network with two orgs and an ordering org, each with a single peer.
  - Create a channel (default name is mychannel): `./network.sh createChannel`
  - Install and start the chaincode on the channel: `./network.sh deployCC -ccn basic -ccp ../project2/medical-records -ccl go`.

### Now build & run the rest client
The rest-api-go client which allows us to invoke functions of the smart contract via curl requests is used. This code is tweaked from the original fabric-samples version to fit this application and to be able to help display the JSON data in a more readable fashion. 
- cd into project2/rest-api-go directory
- Download required dependencies using `go mod download`
- Run `go run main.go` to run the REST server

## Sending Requests

Invoke endpoint accepts POST requests with chaincode function and arguments. Query endpoint accepts GET requests with chaincode function and arguments.

Sample request to invoke the "InitLedger" function. This will populate the ledger with two patients (each with two medical records). Response will contain transaction ID for a successful invoke.

``` sh
curl --request POST \
  --url http://localhost:3000/invoke \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data = \
  --data channelid=mychannel \
  --data chaincodeid=basic \
  --data function=InitLedger
```
Sample chaincode query for getting all patients and their medical records details.

``` sh
curl --request GET \
  --url 'http://localhost:3000/query?channelid=mychannel&chaincodeid=basic&function=GetAllMedicalRecords' | jq
  ```

Get a patients medical records by their social security number:
``` sh
curl --request GET \
  --url "http://localhost:3000/query?channelid=mychannel&chaincodeid=basic&function=GetPatientRecordBySSN&args=123-45-6789" | jq
  ```

To add a new medical record for patient with social security number 123-45-6789:
``` sh
curl --request POST \
  --url http://localhost:3000/invoke \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data = \
  --data channelid=mychannel \
  --data chaincodeid=basic \
  --data function=AddPatientRecord \
  --data args=123-45-6789 \
  --data args=2022-03-20 \
  --data args="Dr. Strange" \
  --data args="Crumlin Hospital" \
  --data args="Strep Throat" \
  --data args=Penicillin
```

To register a new patient:
``` sh
curl --request POST \
  --url http://localhost:3000/invoke \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data channelid=mychannel \
  --data chaincodeid=basic \
  --data function=RegisterNewPatient \
  --data args=987-65-4321 \
  --data args="John Doe" \
  --data args=29 \
  --data args=Male
  ```


