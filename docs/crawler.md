## Crawler

**Crawler flow**

1.  Start crawler

```Usage:
  crawler [command]

Available Commands:
  clients     List existing supported clients (github, gitlab, bitbucket).
  crawl whitelist.yml       Crawl publiccode.yml file from domains in whitelist.yml file.
  domains     List all the Domains.
  list        List all the PA in the whitelist file.
  one http://repo.url        Crawl publiccode.yml from one single repository url.
  version     Print the version number of the crawler.
```

with command:
`crawler crawl whitelistGeneric.yml whitelistPA.yml`

2.  -> if "./ipa/amministrazioni.txt" is older than two days, update to latest version. (if some error occurs, log the error and use the actual "amministrazioni.txt")
3.  -> connect to Elasticsearch node. If ko, log the error and stop the crawler. If ok set the mapping for publiccode and administrations.
4.  -> read and parse domains.yml (ref: domains.yml) and whitelist.yml (ref: whitelist.yml)
5.  -> initialize a repositories channel that will host all the repository infos.
6.  -> Initialize a goroutine that will process the main organization url (crawler.ProcessPA), detect the right host (crawler.KnownHost) and, using the right APIs, generate the right URL to the publiccode.yml raw file (crawler.ProcessPADomain) and save it to the repositories channel.
7.  -> With another goroutine, process every repository in repositories channel (crawler.ProcessRepositories) checking the availability. If the file is available: validate (crawler.validateRemoteFile), save the publiccode.yml (crawler.SaveToFile) and the metadata returned from the APIs (crawler.SaveToFile) to a file, clone the repository (crawler.CloneRepository) locally, calculate the repo activity/vitalityIndex (crawler.CalculateRepoActivity) and save the data on ES (crawler.SaveToES).
8.  -> Flush elastic commits.
9.  -> Update elastic aliases with the alias `publiccode`
10. -> Generate Jekyll website data for the the OSS catalog.
