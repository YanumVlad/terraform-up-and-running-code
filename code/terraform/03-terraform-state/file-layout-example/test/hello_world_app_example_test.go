package test

import (
	"github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"testing"
	"fmt"
	"time"
	"strings"
)

func TestHelloWorldAppExample(t *testing.T)  {
	opts := &terraform.Options{
		// You should update this relative path to point at your
		// hello-world-app example directory!
		TerraformDir: "../examples/hello-world-app/standalone",

		Vars: map[string]interface{}{
			"mysql_config": map[string]interface{}{
				"address": "mock-value-for-test",
				"port": 3306,
			},
		},
	}

	// Clean up everything at the end of the test
	defer terraform.Destroy(t, opts)

	// Deploy the example
	terraform.InitAndApply(t, opts)

	// Get the URL of the ALB
	albDnsName := terraform.OutputRequired(t, opts, "alb_dns_name")
	url := fmt.Sprintf("http://%s", albDnsName)

	expectedStatus := 200
	expectedBody := "Hello, World"

	maxRetries := 10
	timeBetweenRetries := 10 * time.Second

	customValidation := func(statusCode int, response string) bool {
		return statusCode == expectedStatus && strings.Contains(response, expectedBody)
	}


	//http_helper.HttpGetWithRetry(
	http_helper.HttpGetWithRetryWithCustomValidation(
		t, 
		url,
		nil,
		maxRetries,
		timeBetweenRetries,
		customValidation,
	)
}