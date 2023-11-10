// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iot_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/iot"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfiot "github.com/hashicorp/terraform-provider-aws/internal/service/iot"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccIoTPolicy_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var v iot.GetPolicyOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_iot_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, iot.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPolicyExists(ctx, resourceName, &v),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "iot", fmt.Sprintf("policy/%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "default_version_id", "1"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "policy"),
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

func TestAccIoTPolicy_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var v iot.GetPolicyOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_iot_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, iot.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(ctx, resourceName, &v),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfiot.ResourcePolicy(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccIoTPolicy_update(t *testing.T) {
	ctx := acctest.Context(t)
	var v iot.GetPolicyOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_iot_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, iot.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyConfig_resourceName(rName, fmt.Sprintf("arn:aws:iot:*:*:topic/%s", sdkacctest.RandomWithPrefix(acctest.ResourcePrefix))),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "default_version_id", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPolicyConfig_resourceName(rName, fmt.Sprintf("arn:aws:iot:*:*:topic/%s", sdkacctest.RandomWithPrefix(acctest.ResourcePrefix))),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "default_version_id", "2"),
				),
			},
		},
	})
}

func testAccCheckPolicyDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).IoTConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_iot_policy" {
				continue
			}

			_, err := tfiot.FindPolicyByName(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("IoT Policy %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckPolicyExists(ctx context.Context, n string, v *iot.GetPolicyOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).IoTConn(ctx)

		output, err := tfiot.FindPolicyByName(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccPolicyConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_iot_policy" "test" {
  name = %[1]q

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "iot:*"
      ],
      "Resource": [
        "*"
      ]
    }
  ]
}
EOF
}
`, rName)
}

func testAccPolicyConfig_resourceName(rName, resourceName string) string {
	return fmt.Sprintf(`
resource "aws_iot_policy" "test" {
  name = %[1]q

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "iot:*"
      ],
      "Resource": [
        %[2]q
      ]
    }
  ]
}
EOF
}
`, rName, resourceName)
}
