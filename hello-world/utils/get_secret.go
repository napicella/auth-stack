package utils

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func GetSecret() (string, error) {
	s, e := session.NewSession()
	if e != nil {
		return "", e
	}
	svc := secretsmanager.New(s)
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String("jwt-signing-key"),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return "", err
	}
	jsonMap := make(map[string]string)
	e = json.Unmarshal([]byte(*result.SecretString), &jsonMap)
	if e != nil {
		return "", e
	}
	return jsonMap[*result.Name], nil
}
