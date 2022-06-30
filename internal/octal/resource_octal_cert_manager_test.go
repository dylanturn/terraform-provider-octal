package octal

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceScaffolding(t *testing.T) {
	//t.Skip("resource not yet implemented, remove this once you add your own code")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScaffolding,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

const testAccResourceScaffolding = `
resource "octal_cert_manager" "foo" {
  
}
`
