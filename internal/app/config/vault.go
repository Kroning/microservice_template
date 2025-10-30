package config

import (
	"context"
	"os"
	"time"

	vault "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"{{index .App "git"}}/pkg/logger"
)

const (
	pathApp            = "CONFIG_PATH_APP"
	vaultConfigPathKey = "CONFIG_PATH_KEY"
	vaultAddr          = "VAULT_ADDR"
	vaultRoleId        = "VAULT_ROLE_ID"
	vaultSecretId      = "VAULT_SECRET_ID"
	roleIDFile         = "./data/vault-creds/role_id"
	secretIDFile       = "./data/vault-creds/secret_id"
	enableVault        = "VAULT_ENABLE"
)

var (
	vaultClient *vault.Client
)

func loadFromVault(ctx context.Context) (loaded bool, err error) {
	if os.Getenv(enableVault) != "yes" {
		logger.Info(ctx, "vault is not enabled with os env: "+enableVault+"!=yes")
		return false, nil
	}

	if os.Getenv(vaultAddr) == "" {
		logger.Info(ctx, "vault address empty, skip loading from vault")
		return false, nil
	}
	logger.Info(ctx, "loading from vault", zap.String("addr", os.Getenv(vaultAddr)),
		zap.String("path", os.Getenv(pathApp)), zap.String("key", os.Getenv(vaultConfigPathKey)))
	vaultClient, err = initVaultClient()
	if err != nil {
		return false, err
	}

	fullPath := os.Getenv(pathApp) + "/secret/data/" + os.Getenv(vaultConfigPathKey)

	data, err := vaultClient.Logical().Read(fullPath)
	if err != nil {
		return false, errors.Wrap(err, "vaultClient Read with path: "+fullPath)
	}
	if data == nil {
		logger.Warn(ctx, "path not exists", zap.String("path", fullPath))
		return false, nil
	}

	err = viper.MergeConfigMap(data.Data["data"].(map[string]interface{}))
	if err != nil {
		return false, err
	}

	return true, nil
}

func initVaultClient() (*vault.Client, error) {
	if vaultClient != nil {
		return vaultClient, nil
	}
	if os.Getenv(vaultAddr) == "" {
		return nil, errors.New("empty vault addr")
	}
	vaultConfig := vault.DefaultConfig()
	vaultConfig.Timeout = 10 * time.Second
	client, err := vault.NewClient(vaultConfig)
	if err != nil {
		return nil, errors.Wrap(err, "problem while trying to construct new vault client")
	}

	vaultRole, vaultSecret := os.Getenv(vaultRoleId), os.Getenv(vaultSecretId)

	// Reading from file to test purpose
	if vaultRole == "" || vaultSecret == "" {
		roleIdBytes, err := os.ReadFile(roleIDFile)
		if err != nil {
			return nil, err
		}
		secretIdBytes, err := os.ReadFile(secretIDFile)
		if err != nil {
			return nil, err
		}

		vaultRole, vaultSecret = string(roleIdBytes), string(secretIdBytes)
	}

	token, err := authAppRole(client, vaultRole, vaultSecret)
	if err != nil {
		return nil, err
	}
	client.SetToken(token)
	vaultClient = client
	return vaultClient, nil
}

func authAppRole(client *vault.Client, roleId, secretId string) (string, error) {
	requestPath := "auth/approle/login"
	options := map[string]interface{}{
		"role_id":   roleId,
		"secret_id": secretId,
	}
	secret, err := client.Logical().Write(requestPath, options)
	if err != nil {
		return "", err
	}
	return secret.Auth.ClientToken, nil
}
