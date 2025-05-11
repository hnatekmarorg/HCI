data "authentik_flow" "default_provider_authorization_implicit_consent" {
  slug = "default-provider-authorization-implicit-consent"
}

data "authentik_flow" "invalidation_flow" {
  slug = "default-invalidation-flow"
}

data "authentik_certificate_key_pair" "generated" {
  name = "authentik Self-signed Certificate"
}

resource "authentik_provider_oauth2" "bao" {
  signing_key = data.authentik_certificate_key_pair.generated.id
  authorization_flow = data.authentik_flow.default_provider_authorization_implicit_consent.id
  client_id          = "openbao"
  invalidation_flow  = data.authentik_flow.invalidation_flow.id
  name               = "openbao"
  allowed_redirect_uris = [
    {
      url = "https://openbao.hnatekmar.xyz/ui/vault/auth/oidc/oidc/callback",
      matching_mode = "strict"
    },
    {
      url = "https://openbao.hnatekmar.xyz/oidc/callback",
      matching_mode = "strict"
    },
    {
      url = "http://localhost:8250/oidc/callback",
      matching_mode = "strict"
    },
  ]
}

resource "authentik_application" "openbao" {
  name = "openbao"
  slug = "openbao"
  protocol_provider = authentik_provider_oauth2.bao.id
}

resource "vault_jwt_auth_backend" "authentik" {
  path                = "oidc"
  type                = "oidc"
  oidc_discovery_url = "https://authentik.hnatekmar.xyz/application/o/openbao/"
  oidc_client_id = authentik_provider_oauth2.bao.client_id
  oidc_client_secret = authentik_provider_oauth2.bao.client_secret
  default_role = "default"
}

resource "vault_jwt_auth_backend_role" "default" {
  backend = vault_jwt_auth_backend.authentik.path
  allowed_redirect_uris = [
      "https://openbao.hnatekmar.xyz/ui/vault/auth/oidc/oidc/callback",
      "https://openbao.hnatekmar.xyz/oidc/callback",
      "http://localhost:8250/oidc/callback",
  ]
  bound_audiences = [
    authentik_provider_oauth2.bao.client_id
  ]
  role_name  = "default"
  user_claim = "sub"
  token_policies = [
  ]
}