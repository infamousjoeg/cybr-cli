## cybr logon

Logon to PAS REST API

### Synopsis

Authenticate to PAS REST API using the provided authentication type.
	
	Example Usage:
	$ cybr logon -u $USERNAME -a $AUTH_TYPE -b https://pvwa.example.com
	To bypass TLS verification:
	$ cybr logon -u $USERNAME -a $AUTH_TYPE -b https://pvwa.example.com -i

```
cybr logon [flags]
```

### Options

```
  -a, --auth-type string   Authentication method to logon using [cyberark|ldap|radius]
  -b, --base-url string    Base URL to send Logon request to [https://pvwa.example.com]
  -h, --help               help for logon
  -i, --insecure-tls       If detected, TLS will not be verified
      --non-interactive    If detected, will retrieve the password from the PAS_PASSWORD environment variable
  -u, --username string    Username to logon PAS REST API using
```

### Options inherited from parent commands

```
      --verbose   To enable verbose logging
```

### SEE ALSO

* [cybr](cybr.md)	 - cybr is CyberArk's PAS command-line interface utility

###### Auto generated by spf13/cobra on 28-Jan-2021