package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/hashicorp/vault/api"
)

func initializeVault(vaultAddr string, secretName string, region string) error {

	config := &api.Config{
		Address: vaultAddr,
	}

	config.ConfigureTLS(&api.TLSConfig{
		Insecure: true,
	})

	client, err := api.NewClient(config)
	if err != nil {
		return err
	}

	initRequest := &api.InitRequest{
		RecoveryShares:    1,
		RecoveryThreshold: 1,
	}

	initResponse, err := client.Sys().Init(initRequest)
	if err != nil {
		return err
	}

	rootToken := initResponse.RootToken
	recoveryKey := initResponse.RecoveryKeysB64[0]
	return storeVaultSecrets(rootToken, recoveryKey, secretName, region)
}

func storeVaultSecrets(rootToken, recoveryKey, secretName, region string) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	svc := secretsmanager.New(sess)
	secretString, err := json.Marshal(map[string]string{
		"root_token":   rootToken,
		"recovery_key": recoveryKey,
	})
	if err != nil {
		return err
	}

	// Try to update the secret
	_, err = svc.UpdateSecret(&secretsmanager.UpdateSecretInput{
		SecretId:     aws.String(secretName),
		SecretString: aws.String(string(secretString)),
	})

	// Check if the error is because the secret does not exist
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == secretsmanager.ErrCodeResourceNotFoundException {
			// Create the secret if it does not exist
			_, err = svc.CreateSecret(&secretsmanager.CreateSecretInput{
				Name:         aws.String(secretName),
				SecretString: aws.String(string(secretString)),
			})
			if err != nil {
				return fmt.Errorf("failed to create secret: %w", err)
			}
		} else {
			return fmt.Errorf("failed to update secret: %w", err)
		}
	}

	return nil
}
