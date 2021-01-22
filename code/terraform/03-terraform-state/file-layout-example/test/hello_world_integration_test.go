package test

import (
	"github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/test-structure"
	"testing"
	"fmt"
	"time"
	"strings"
)

const dbDirStage = "../stage/data-stores/mysql"
const appDirStage = "../stage/services/hello-world-app"


func TestHelloWorldAppStage(t *testing.T)  {

	t.Parallel()

	// Deploy the MySQL DB
	dbOpts := createDbOpts(t, dbDirStage)
	defer terraform.Destroy(t, dbOpts)
	terraform.InitAndApply(t, dbOpts)

	// Deploy the hello-world-app
	helloOpts := createHelloOpts(dbOpts, appDirStage)
	defer terraform.Destroy(t, helloOpts)
	terraform.InitAndApply(t, helloOpts)

	// Validate the hello-world-app works
	validateHelloApp(t, helloOpts)
}

func validateHelloApp(t *testing.T, helloOpts *terraform.Options) {
	albDnsName := terraform.OutputRequired(t, helloOpts, "alb_dns_name")
	url := fmt.Sprintf("http://%s", albDnsName)

	maxRetries := 10
	timeBetweenRetries := 10 * time.Second

	http_helper.HttpGetWithRetryWithCustomValidation(
		t,
		url,
		nil,
		maxRetries,
		timeBetweenRetries,
		func(status int, body string) bool {
			return status == 200 &&
				strings.Contains(body, "Hello, world")
		},
	)
}

func createDbOpts(t *testing.T, terraformDir string) *terraform.Options {
	uniqueId := strings.ToLower(random.UniqueId())

	bucketForTesting := "ulny-s3"
	bucketRegionForTesting := "us-east-2"
	dbStateKey := fmt.Sprintf("%s/%s/terraform.tfstate", t.Name(), uniqueId)

	return &terraform.Options{
		TerraformDir: terraformDir,

		Vars: map[string]interface{}{
			"db_name": fmt.Sprintf("test%s", uniqueId),
			"db_password": "password",
			"cluster_name": fmt.Sprintf("test%s", uniqueId),

		},

		BackendConfig: map[string]interface{}{
			"bucket":         bucketForTesting,
			"region":         bucketRegionForTesting,
			"key":            dbStateKey,
			"encrypt":        true,
		},
	}
}

func createHelloOpts(dbOpts *terraform.Options, terraformDir string) *terraform.Options {

	return &terraform.Options{
		TerraformDir: terraformDir,

		Vars: map[string]interface{}{
			"db_remote_state_bucket": dbOpts.BackendConfig["bucket"],
			"db_remote_state_key": dbOpts.BackendConfig["key"],
			"environment": dbOpts.Vars["db_name"],
		},
	}
}

func TestHelloWorldAppStageWithStages(t *testing.T) {
	t.Parallel()

	// Store the function in a short variable name solely to make the
	// code examples fit better in the book.
	stage := test_structure.RunTestStage

	// Deploy the MySQL DB
	defer stage(t, "teardown_db", func() { teardownDb(t, dbDirStage) })
	stage(t, "deploy_db", func() { deployDb(t, dbDirStage) })

	// Deploy the hello-world-app
	defer stage(t, "teardown_app", func() { teardownApp(t, appDirStage) })
	stage(t, "deploy_app", func() { deployApp(t, dbDirStage, appDirStage) })

	// Validate the hello-world-app works
	stage(t, "validate_app", func() { validateApp(t, appDirStage) })

	// Redeploy the hello-world-app
	stage(t, "redeploy_app", func() { redeployApp(t, appDirStage) })
	
}

func redeployApp(t *testing.T, helloAppDir string) {
	helloOpts := test_structure.LoadTerraformOptions(t, helloAppDir)

	albDnsName := terraform.OutputRequired(t, helloOpts, "alb_dns_name")
	url := fmt.Sprintf("http://%s", albDnsName)

	// Start checking every 1s that the app is responding with a 200 OK
	stopChecking := make(chan bool, 1)
	waitGroup, _ := http_helper.ContinuouslyCheckUrl(
		t,
		url,
		stopChecking,
		1*time.Second,
	)

	// Update the server text and redeploy
	newServerText := "Hello, world, v2!"
	helloOpts.Vars["server_text"] = newServerText
	terraform.Apply(t, helloOpts)

	// Make sure the new version deployed
	maxRetries := 10
	timeBetweenRetries := 10 * time.Second
	http_helper.HttpGetWithRetryWithCustomValidation(
		t,
		url,
		nil,
		maxRetries,
		timeBetweenRetries,
		func(status int, body string) bool {
			return status == 200 &&
				strings.Contains(body, newServerText)
		},
	)

	// Stop checking
	stopChecking <- true
	waitGroup.Wait()
}

func deployDb(t *testing.T, dbDir string) {
	dbOpts := createDbOpts(t, dbDir)

	// Save data to disk so that other test stages executed at a later
	// time can read the data back in
	test_structure.SaveTerraformOptions(t, dbDir, dbOpts)

	terraform.InitAndApply(t, dbOpts)
}

func teardownDb(t *testing.T, dbDir string) {
	dbOpts := test_structure.LoadTerraformOptions(t, dbDir)
	defer terraform.Destroy(t, dbOpts)
}

func deployApp(t *testing.T, dbDir string, helloAppDir string) {
	dbOpts := test_structure.LoadTerraformOptions(t, dbDir)
	helloOpts := createHelloOpts(dbOpts, helloAppDir)

	// Save data to disk so that other test stages executed at a later
	// time can read the data back in
	test_structure.SaveTerraformOptions(t, helloAppDir, helloOpts)

	terraform.InitAndApply(t, helloOpts)
}

func teardownApp(t *testing.T, helloAppDir string) {
	helloOpts := test_structure.LoadTerraformOptions(t, helloAppDir)
	defer terraform.Destroy(t, helloOpts)
}

func validateApp(t *testing.T, helloAppDir string)  {
	helloOpts := test_structure.LoadTerraformOptions(t, helloAppDir)
	validateHelloApp(t, helloOpts)
}