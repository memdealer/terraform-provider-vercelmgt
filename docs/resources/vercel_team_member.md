---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "vercel_team_member Resource - vercel-mgt"
subcategory: ""
description: |-
  Retrieves information about member user in given team
---

# vercel_team_member (Resource)

Retrieves information about member user in given team

## Example Usage

```terraform
resource "vercel_team_member" "qwe" {
    email = "test@gmail.com"
    role  = "OWNER"
}
```

<!-- schema generated by tfplugindocs -->

## Schema

### Required

- `email` (String) Email of the team member
- `role` (String) Role of the team member

### Read-Only

- `uid` (String) Unique identification of the team member on _Vercel_

## Import

Import is supported using the following syntax:

```shell
  terraform import 'vercel_team_member.qwe iliketurtles@gmail.com
```
