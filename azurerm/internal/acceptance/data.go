package acceptance

import (
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azuread/azuread"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/provider"
)

type TestData struct {
	// Locations is a set of Azure Regions which should be used for this Test
	Locations Regions

	// RandomString is a random integer which is unique to this test case
	RandomInteger int

	// RandomString is a random 5 character string is unique to this test case
	RandomString string

	// ResourceName is the fully qualified resource name, comprising of the
	// resource type and then the resource label
	// e.g. `azurerm_resource_group.test`
	ResourceName string

	// ResourceType is the Terraform Resource Type - `azurerm_resource_group`
	ResourceType string

	// Environment is a struct containing Details about the Azure Environment
	// that we're running against
	Environment azure.Environment

	// EnvironmentName is the name of the Azure Environment where we're running
	EnvironmentName string

	// resourceLabel is the local used for the resource - generally "test""
	resourceLabel string
}

// BuildTestData generates some test data for the given resource
func BuildTestData(t *testing.T, resourceType string, resourceLabel string) TestData {
	azureProvider := provider.AzureProvider().(*schema.Provider)

	AzureProvider = azureProvider
	SupportedProviders = map[string]terraform.ResourceProvider{
		"azurerm": azureProvider,
		"azuread": azuread.Provider().(*schema.Provider),
	}

	env, err := Environment()
	if err != nil {
		t.Fatalf("Error retrieving Environment: %+v", err)
	}

	testData := TestData{
		RandomInteger:   tf.AccRandTimeInt(),
		RandomString:    acctest.RandString(5),
		ResourceName:    fmt.Sprintf("%s.%s", resourceType, resourceLabel),
		Environment:     *env,
		EnvironmentName: EnvironmentName(),

		ResourceType:  resourceType,
		resourceLabel: resourceLabel,
	}

	if features.UseDynamicTestLocations() {
		testData.Locations = availableLocations()
	} else {
		testData.Locations = Regions{
			Primary:   os.Getenv("ARM_TEST_LOCATION"),
			Secondary: os.Getenv("ARM_TEST_LOCATION_ALT"),
			Ternary:   os.Getenv("ARM_TEST_LOCATION_ALT2"),
		}
	}

	return testData
}

func (td *TestData) RandomIntOfLength(len int) int {
	// len should not be
	//  - greater then 18, longest a int can represent
	//  - less then 8, as that gives us YYMMDDRR
	if 8 > len || len > 18 {
		panic(fmt.Sprintf("Invalid Test: RandomIntOfLength: len is not between 8 or 18 inclusive"))
	}

	r := td.RandomInteger % 100

	// 18 - just return the int
	if len >= 18 {
		return td.RandomInteger
	}

	// 16-17 just strip off the last 1-2 digits
	if len >= 16 {
		return td.RandomInteger / int(math.Pow10(18-len))
	}

	// 8-15 remove the last x digits
	a := td.RandomInteger / int(math.Pow10(18-len+2))

	// multiply by 100 and add last two digits of randomness back in
	return (a * 100) + r
}
