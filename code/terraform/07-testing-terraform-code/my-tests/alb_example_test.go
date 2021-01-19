package test

import (
	"testing"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/http-helper"
	"fmt"
	"time"
)

func TestGoIsWorking(t *testing.T)  {
	fmt.Println()
	fmt.Println("If you see this text, it's working!")
	fmt.Println()
}

func TestAlbExample(t *testing.T)  {

	opts := &terraform.Options{
		// You should update this relative path to point at your alb
		// example directory!
		TerraformDir: "../examples/alb",
	}

	// Clean up everything at the end of the test
	defer terraform.Destroy(t, opts)


	// Deploy the example
	terraform.InitAndApply(t, opts)

	albDnsName := terraform.OutputRequired(t, opts, "alb_dns_name")
	url := fmt.Sprintf("http://%s", albDnsName)


	// Test that the ALB's default action is working and returns a 404

	expectedStatus := 404
	expectedBody := "404: page not found"

	maxRetries := 10
	timeBetweenRetries := 10 * time.Second

	http_helper.HttpGetWithRetry(
		t,
		url,
		nil,
		expectedStatus,
		expectedBody,
		maxRetries,
		timeBetweenRetries,
	)



}