package vault

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/agentio/v/pkg/pretty"
	"github.com/agentio/v/pkg/vault"
	"github.com/spf13/cobra"
)

var jwksUrl string

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "vault CLUSTER",
		Args: cobra.ExactArgs(1),
		RunE: action,
	}
	cmd.Flags().StringVar(&jwksUrl, "jwks_url", "http://localhost:4646/.well-known/jwks.json", "jwks_url")
	return cmd
}

func action(cmd *cobra.Command, args []string) error {
	var err error
	keys, err := vault.ReadKeys("")
	if err != nil {
		return err
	}
	cluster := args[0]
	var method *AuthMethod
	if method, err = getAuthMethod(keys, cluster); err != nil || method.Accessor == "" {
		if err = createAuthMethod(keys, cluster); err != nil {
			return err
		}
		if method, err = getAuthMethod(keys, cluster); err != nil {
			return err
		}
	}
	if err = configureAuthMethod(keys, cluster, jwksUrl); err != nil {
		return err
	}
	if err = createACLPolicy(keys, cluster, method.Accessor); err != nil {
		return err
	}
	if err = createWorkloadRole(keys, cluster); err != nil {
		return err
	}
	return nil
}

type AuthMethod struct {
	Accessor string `json:"accessor"`
}

func getAuthMethod(keys *vault.VaultKeys, cluster string) (*AuthMethod, error) {
	request, err := http.NewRequest("GET", vault.URL("/v1/sys/auth/jwt-"+cluster), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+keys.RootToken)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s\n", string(pretty.JSON(body)))
	var am AuthMethod
	err = json.Unmarshal(body, &am)
	if err == nil {
		return &am, nil
	}
	return nil, err
}

/*
// Create a new auth method.
POST /v1/sys/auth/jwt-CLUSTER

	{
	  "type": "jwt",
	  "description": "",
	  "config": {
	    "options": null,
	    "default_lease_ttl": "",
	    "max_lease_ttl": "",
	    "force_no_cache": false
	  },
	  "local": false,
	  "seal_wrap": false,
	  "external_entropy_access": false,
	  "options": null
	}
*/
func createAuthMethod(keys *vault.VaultKeys, cluster string) error {
	type AuthMethod struct {
		Type string `json:"type"`
	}
	b, err := json.Marshal(&AuthMethod{Type: "jwt"})
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", vault.URL("/v1/sys/auth/jwt-"+cluster), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+keys.RootToken)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(pretty.JSON(body)))
	return nil
}

/*
// Configure the new auth method.
PUT /v1/auth/jwt-CLUSTER/config

	{
	  "default_role": "CLUSTER-workloads",
	  "jwks_url": "http://localhost:4646/.well-known/jwks.json",
	  "jwt_supported_algs": [
	    "EdDSA",
	    "RS256"
	  ]
	}
*/
func configureAuthMethod(keys *vault.VaultKeys, cluster, jwksUrl string) error {
	type AuthMethodConfig struct {
		DefaultRole      string   `json:"default_role"`
		JwksUrl          string   `json:"jwks_url"`
		JwtSupportedAlgs []string `json:"jwt_supported_algs"`
	}
	b, err := json.Marshal(&AuthMethodConfig{
		DefaultRole:      cluster + "-workloads",
		JwksUrl:          jwksUrl, // TODO: parameterize
		JwtSupportedAlgs: []string{"EdDSA", "RS256"},
	})
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", vault.URL("/v1/auth/jwt-"+cluster+"/config"), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+keys.RootToken)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(pretty.JSON(body)))
	return nil
}

// Create an ACL Policy.
/*
The request is PUT /v1/sys/policies/acl/CLUSTER-workloads and here is the request body:
{
  "policy": "cGF0aCAic2VjcmV0L2RhdGEve3tpZGVudGl0eS5lbnRpdHkuYWxpYXNlcy5hdXRoX2p3dF84MTNjM2JiMC5tZXRhZGF0YS5ub21hZF9uYW1lc3BhY2V9fS97e2lkZW50aXR5LmVudGl0eS5hbGlhc2VzLmF1dGhfand0XzgxM2MzYmIwLm1ldGFkYXRhLm5vbWFkX2pvYl9pZH19LyoiIHsKICBjYXBhYmlsaXRpZXMgPSBbInJlYWQiXQp9CgpwYXRoICJzZWNyZXQvZGF0YS97e2lkZW50aXR5LmVudGl0eS5hbGlhc2VzLmF1dGhfand0XzgxM2MzYmIwLm1ldGFkYXRhLm5vbWFkX25hbWVzcGFjZX19L3t7aWRlbnRpdHkuZW50aXR5LmFsaWFzZXMuYXV0aF9qd3RfODEzYzNiYjAubWV0YWRhdGEubm9tYWRfam9iX2lkfX0iIHsKICBjYXBhYmlsaXRpZXMgPSBbInJlYWQiXQp9CgpwYXRoICJzZWNyZXQvbWV0YWRhdGEve3tpZGVudGl0eS5lbnRpdHkuYWxpYXNlcy5hdXRoX2p3dF84MTNjM2JiMC5tZXRhZGF0YS5ub21hZF9uYW1lc3BhY2V9fS8qIiB7CiAgY2FwYWJpbGl0aWVzID0gWyJsaXN0Il0KfQoKcGF0aCAic2VjcmV0L21ldGFkYXRhLyoiIHsKICBjYXBhYmlsaXRpZXMgPSBbImxpc3QiXQp9Cg=="
}
$ jq < acl.json .policy -r | base64 -d
path "secret/data/{{identity.entity.aliases.auth_jwt_813c3bb0.metadata.nomad_namespace}}/{{identity.entity.aliases.auth_jwt_813c3bb0.metadata.nomad_job_id}}/*" {
  capabilities = ["read"]
}

path "secret/data/{{identity.entity.aliases.auth_jwt_813c3bb0.metadata.nomad_namespace}}/{{identity.entity.aliases.auth_jwt_813c3bb0.metadata.nomad_job_id}}" {
  capabilities = ["read"]
}

path "secret/metadata/{{identity.entity.aliases.auth_jwt_813c3bb0.metadata.nomad_namespace}}/*" {
  capabilities = ["list"]
}

path "secret/metadata/*" {
  capabilities = ["list"]
}
*/
func createACLPolicy(keys *vault.VaultKeys, cluster, accessor string) error {
	policy := `
path "{{identity.entity.aliases.` + accessor + `.metadata.nomad_namespace}}/{{identity.entity.aliases.` + accessor + `.metadata.nomad_job_id}}/*" {
  capabilities = ["read"]
}

path "{{identity.entity.aliases.` + accessor + `.metadata.nomad_namespace}}/{{identity.entity.aliases.` + accessor + `.metadata.nomad_job_id}}" {
  capabilities = ["read"]
}
	`
	encodedPolicy := base64.StdEncoding.EncodeToString([]byte(policy))
	type ACLPolicy struct {
		Policy string `json:"policy"`
	}
	b, err := json.Marshal(&ACLPolicy{
		Policy: encodedPolicy,
	})
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", vault.URL("/v1/sys/policies/acl/"+cluster+"-workloads"), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+keys.RootToken)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(pretty.JSON(body)))
	return nil
}

// Create the workload role.
/*
The request is PUT /v1/auth/jwt-CLUSTER/role/CLUSTER-workloads and here is the request body:
{
  "bound_audiences": "vault.io",
  "claim_mappings": {
    "nomad_job_id": "nomad_job_id",
    "nomad_namespace": "nomad_namespace",
    "nomad_task": "nomad_task"
  },
  "role_type": "jwt",
  "token_period": "30m",
  "token_policies": [
    "CLUSTER-workloads"
  ],
  "token_type": "service",
  "user_claim": "/nomad_job_id",
  "user_claim_json_pointer": true
}
*/
func createWorkloadRole(keys *vault.VaultKeys, cluster string) error {
	type Role struct {
		BoundAudiences       string            `json:"bound_audiences"`
		ClaimMappings        map[string]string `json:"claim_mappings"`
		RoleType             string            `json:"role_type"`
		TokenPeriod          string            `json:"token_period"`
		TokenPolicies        []string          `json:"token_policies"`
		TokenType            string            `json:"token_type"`
		UserClaim            string            `json:"user_claim"`
		UserClaimJsonPointer bool              `json:"user_claim_json_pointer"`
	}
	b, err := json.Marshal(&Role{
		BoundAudiences: "vault.io",
		ClaimMappings: map[string]string{
			"nomad_job_id":    "nomad_job_id",
			"nomad_namespace": "nomad_namespace",
			"nomad_task":      "nomad_task",
		},
		RoleType:    "jwt",
		TokenPeriod: "30m",
		TokenPolicies: []string{
			cluster + "-workloads",
		},
		UserClaim:            "/nomad_job_id",
		UserClaimJsonPointer: true,
	})
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", vault.URL("/v1/auth/jwt-"+cluster+"/role/"+cluster+"-workloads"), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+keys.RootToken)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(pretty.JSON(body)))
	return nil
}
