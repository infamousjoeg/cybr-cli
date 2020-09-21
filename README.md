# pas-api-go
@CyberArk Privileged Access Security (PAS) REST API Client Library

## Required Environment Variables

| Variable Name | Description |
| --- | --- |
| `PAS_HOSTNAME` | Base URL for PAS REST API Web Service |
| `PAS_USERNAME` | Username to authn to PAS REST API |
| `PAS_PASSWORD` | Password associated with `PAS_USERNAME` |
| `PAS_AUTH_TYPE` | Authentication method to use (cyberark or ldap) |

## Testing

1. Download and install [summon](https://cyberark.github.io/summon).
2. [OPTIONAL] Download and install a [summon provider](https://cyberark.github.io/summon/#providers).
   1. I use the `keyring` provider with [conceal](https://github.com/infamousjoeg/conceal).
3. Modify the values in [secrets.yml]() for your environment.
   1. If you did not complete the optional Step #2, you will use literal strings for `PAS_USERNAME` and `PAS_PASSWORD` similar to the values of `PAS_HOSTNAME` and `PAS_AUTH_TYPE`.
4. Run [main.go]() with the command: `summon go run main.go`.

### Successful Output

```shell
$ summon go run main.go

Verify JSON:
{"ApplicationName":"PasswordVault","AuthenticationMethods":[{"Enabled":false,"Id":"windows"},{"Enabled":false,"Id":"pki"},{"Enabled":true,"Id":"cyberark"},{"Enabled":false,"Id":"oraclesso"},{"Enabled":false,"Id":"rsa"},{"Enabled":true,"Id":"radius"},{"Enabled":true,"Id":"ldap"},{"Enabled":true,"Id":"saml"}],"ServerId":"e33e4f16-b637-23e9-8329-ccd02f0102323","ServerName":"Vault"}

Session Token:
YWYyMmYxYTQtNzM5OC00ZWRhLTkeODYtOTkxZTY2NWUzYzNlOzQwNjAwaTI3ODQ4NUg5NkY7MDAwMDAwMDI4NUE4OTcxMkYwNTI1RjY2RjA3QzI5NDRGrEYyOTVCNjVCNjMxODcyMUQ2M0VCMUFEN0VCRjn5MzA0NTIzNDUyMDAwMDAwMDA7
```