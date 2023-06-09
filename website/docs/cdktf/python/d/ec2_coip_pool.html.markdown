---
subcategory: "Outposts (EC2)"
layout: "aws"
page_title: "AWS: aws_ec2_coip_pool"
description: |-
    Provides details about a specific EC2 Customer-Owned IP Pool
---

# Data Source: aws_ec2_coip_pool

Provides details about a specific EC2 Customer-Owned IP Pool.

This data source can prove useful when a module accepts a coip pool id as
an input variable and needs to, for example, determine the CIDR block of that
COIP Pool.

## Example Usage

The following example returns a specific coip pool ID

```terraform
variable "coip_pool_id" {}

data "aws_ec2_coip_pool" "selected" {
  id = var.coip_pool_id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available
COIP Pools in the current region. The given filters must match exactly one
COIP Pool whose data will be exported as attributes.

* `local_gateway_route_table_id` - (Optional) Local Gateway Route Table Id assigned to desired COIP Pool

* `pool_id` - (Optional) ID of the specific COIP Pool to retrieve.

* `tags` - (Optional) Mapping of tags, each pair of which must exactly match
  a pair on the desired COIP Pool.

More complex filters can be expressed using one or more `filter` sub-blocks,
which take the following arguments:

* `name` - (Required) Name of the field to filter by, as defined by
  [the underlying AWS API](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeCoipPools.html).

* `values` - (Required) Set of values that are accepted for the given field.
  A COIP Pool will be selected if any one of the given values matches.

## Attributes Reference

All of the argument attributes except `filter` blocks are also exported as
result attributes. This data source will complete the data by populating
any fields that are not included in the configuration with the data for
the selected COIP Pool.

In addition, the following attributes are exported:

* `arn` - ARN of the COIP pool
* `pool_cidrs` - Set of CIDR blocks in pool

## Timeouts

[Configuration options](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts):

- `read` - (Default `20m`)

<!-- cache-key: cdktf-0.17.0-pre.15 input-4018a02fb3e3ba399170a12145e567c04dd4dcf3d3de820650bd0da4b7852ab7 -->