terraform {
  required_providers {
    vercel = {
      source  = "vercel/vercel"
    }
  }
}

provider "vercel" {
  api_token = var.team_token
  team = var.team_slug
}
resource "vercel_team_member" "qwe" {
    email = "test@gmail.com"
    role  = "OWNER"
}
