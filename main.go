package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

type InputData struct {
	VaultAddr  string `json:"vaultAddr"`
	SecretName string `json:"secretName"`
	Region     string `json:"region"`
}

type OutputData struct {
	Message string `json:"message"`
}

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, inputEvent json.RawMessage) (*OutputData, error) {

	log.Println("Vault initialization started")

	inputData := InputData{
		VaultAddr:  "https://vault.platform.corp:8200",
		SecretName: "vault-secrets",
		Region:     "eu-central-1",
	}

	if err := json.Unmarshal(inputEvent, &inputData); err != nil {
		return nil, err
	}

	// Call initialization function
	err := initializeVault(inputData.VaultAddr, inputData.SecretName, inputData.Region)
	if err != nil {
		log.Fatalf("Vault initialization failed: %v", err)
		return nil, err
	}
	log.Println("Vault initialization completed successfully")

	return &OutputData{Message: "Vault initialization completed successfully"}, nil
}
