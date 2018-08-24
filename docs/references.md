## References

### domains.yml

Contains all the basic auth token for the repositories APIs in the form `Basic <token>`

```- host: "gitlab.com"
  basic-auth:
    - "Basic <base64-auth-token>"
- host: "bitbucket.org"
  basic-auth:
    - "Basic <base64-auth-token>"
- host: "github.com"
  basic-auth:
    - "Basic <base64-auth-token>"
```

### whitelist.yml

Contains a list of every organization.

```
 whitelistPA.yml is a list of every organization repository with an iPA.

- id: "Comune di Bagnacavallo" # generic name of the organization.
  codice-iPA: "c_a547" # codice-iPA
  organizations: # list of organization urls.
    - "https://github.com/gith002"
```

### amministrazioni.txt

Reference: http://www.indicepa.gov.it/documentale/n-opendata.php

### publiccode.yml parsing and validation

Reference: https://github.com/publiccodenet/publiccode.yml-parser-go
