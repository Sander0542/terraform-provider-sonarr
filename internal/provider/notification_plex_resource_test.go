package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationPlexResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccNotificationPlexResourceConfig("resourcePlexTest", "token123"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_plex.test", "auth_token", "token123"),
					resource.TestCheckResourceAttrSet("sonarr_notification_plex.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccNotificationPlexResourceConfig("resourcePlexTest", "token234"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sonarr_notification_plex.test", "auth_token", "token234"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "sonarr_notification_plex.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_token"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationPlexResourceConfig(name, token string) string {
	return fmt.Sprintf(`
	resource "sonarr_notification_plex" "test" {
		on_download                        = false
		on_upgrade                         = false
		on_rename                          = false
		on_series_delete                   = false
		on_episode_file_delete             = false
		on_episode_file_delete_for_upgrade = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		host = "plex.lcl"
		port = 32400
		auth_token = "%s"
	}`, name, token)
}
