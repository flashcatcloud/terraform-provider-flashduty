package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"flashduty": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("FLASHDUTY_APP_KEY"); v == "" {
		t.Skip("FLASHDUTY_APP_KEY not set, skipping acceptance test")
	}
}

func testAccGetEnv(t *testing.T, key string) string {
	t.Helper()
	v := os.Getenv(key)
	if v == "" {
		t.Skipf("Environment variable %s not set, skipping test", key)
	}
	return v
}
