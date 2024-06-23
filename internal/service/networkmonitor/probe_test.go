// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package networkmonitor_test

import (
	"context"
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfnetworkmonitor "github.com/hashicorp/terraform-provider-aws/internal/service/networkmonitor"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccNetworkMonitorProbe_basic(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_networkmonitor_probe.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProbeDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccProbeConfig_basic(rName, "10.0.0.1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProbeExists(ctx, resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "address_family"),
					resource.TestCheckResourceAttrSet(resourceName, names.AttrARN),
					resource.TestCheckResourceAttr(resourceName, "destination", "10.0.0.1"),
					resource.TestCheckNoResourceAttr(resourceName, "destination_port"),
					resource.TestCheckNoResourceAttr(resourceName, "packet_size"),
					resource.TestCheckResourceAttrSet(resourceName, "probe_id"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, acctest.Ct0),
					resource.TestCheckResourceAttrSet(resourceName, names.AttrVPCID),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkMonitorProbe_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_networkmonitor_probe.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, "networkMonitor"),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProbeDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccProbeConfig_basic(rName, "10.0.0.1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProbeExists(ctx, resourceName),
					acctest.CheckFrameworkResourceDisappears(ctx, acctest.Provider, tfnetworkmonitor.ResourceProbe, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccNetworkMonitorProbe_tags(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_networkmonitor_probe.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProbeDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccProbeConfig_tags1(rName, acctest.CtKey1, acctest.CtValue1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProbeExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, acctest.Ct1),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsKey1, acctest.CtValue1),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccProbeConfig_tags2(rName, acctest.CtKey1, acctest.CtValue1Updated, acctest.CtKey2, acctest.CtValue2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProbeExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, acctest.Ct2),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsKey1, acctest.CtValue1Updated),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsKey2, acctest.CtValue2),
				),
			},
			{
				Config: testAccProbeConfig_tags1(rName, acctest.CtKey2, acctest.CtValue2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProbeExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, acctest.Ct1),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsKey2, acctest.CtValue2),
				),
			},
		},
	})
}

func TestAccNetworkMonitorProbe_update(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_networkmonitor_probe.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProbeDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccProbeConfig_full(rName, "10.0.0.1", 8080, 256),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProbeExists(ctx, resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "address_family"),
					resource.TestCheckResourceAttrSet(resourceName, names.AttrARN),
					resource.TestCheckResourceAttr(resourceName, "destination", "10.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "destination_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "packet_size", "256"),
					resource.TestCheckResourceAttrSet(resourceName, "probe_id"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, acctest.Ct0),
					resource.TestCheckResourceAttrSet(resourceName, names.AttrVPCID),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccProbeConfig_full(rName, "10.0.0.2", 8443, 512),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProbeExists(ctx, resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "address_family"),
					resource.TestCheckResourceAttrSet(resourceName, names.AttrARN),
					resource.TestCheckResourceAttr(resourceName, "destination", "10.0.0.2"),
					resource.TestCheckResourceAttr(resourceName, "destination_port", "8443"),
					resource.TestCheckResourceAttr(resourceName, "packet_size", "512"),
					resource.TestCheckResourceAttrSet(resourceName, "probe_id"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, acctest.Ct0),
					resource.TestCheckResourceAttrSet(resourceName, names.AttrVPCID),
				),
			},
		},
	})
}

func testAccCheckProbeDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).NetworkMonitorClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_networkmonitor_probe" {
				continue
			}

			_, err := tfnetworkmonitor.FindProbeByTwoPartKey(ctx, conn, rs.Primary.Attributes["monitor_name"], rs.Primary.Attributes["probe_id"])

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("CloudWatch Network Monitor Probe %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckProbeExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).NetworkMonitorClient(ctx)

		_, err := tfnetworkmonitor.FindProbeByTwoPartKey(ctx, conn, rs.Primary.Attributes["monitor_name"], rs.Primary.Attributes["probe_id"])

		return err
	}
}

func testAccProbeConfig_base(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 1), fmt.Sprintf(`
resource "aws_networkmonitor_monitor" "test" {
  aggregation_period = 30
  monitor_name       = %[1]q

  tags = {
    Name = %[1]q
  }
}
`, rName))
}

func testAccProbeConfig_basic(rName, destination string) string {
	return acctest.ConfigCompose(testAccProbeConfig_base(rName), fmt.Sprintf(`
resource "aws_networkmonitor_probe" "test" {
  monitor_name = aws_networkmonitor_monitor.test.monitor_name
  destination  = %[2]q
  protocol     = "TCP"
  source_arn   = aws_subnet.test[0].arn
}
`, rName, destination))
}

func testAccProbeConfig_full(rName, destination string, port, packetSize int) string {
	return acctest.ConfigCompose(testAccProbeConfig_base(rName), fmt.Sprintf(`
resource "aws_networkmonitor_probe" "test" {
  monitor_name     = aws_networkmonitor_monitor.test.monitor_name
  destination      = %[2]q
  destination_port = %[3]d
  protocol         = "TCP"
  source_arn       = aws_subnet.test[0].arn
  packet_size      = %[4]d
}
`, rName, destination, port, packetSize))
}

func testAccProbeConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccProbeConfig_base(rName), fmt.Sprintf(`
resource "aws_networkmonitor_probe" "test" {
  monitor_name = aws_networkmonitor_monitor.test.monitor_name
  destination  = "10.0.0.1"
  protocol     = "TCP"
  source_arn   = aws_subnet.test[0].arn

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1))
}

func testAccProbeConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccProbeConfig_base(rName), fmt.Sprintf(`
resource "aws_networkmonitor_probe" "test" {
  monitor_name = aws_networkmonitor_monitor.test.monitor_name
  destination  = "10.0.0.1"
  protocol     = "TCP"
  source_arn   = aws_subnet.test[0].arn

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}
