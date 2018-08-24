## Crawler

**Crawler flow**

![Crawler Flow](https://www.websequencediagrams.com/cgi-bin/cdraw?lz=dGl0bGUgQ3Jhd2xpbmcKCiMgQWN0b3JzLgpwYXJ0aWNpcGFudCBVc2VyAAQNIkNtZCIgYXMgQQAJDklwYQASBUIAHw9yYXdsZXIALQVDADsORWxhc3RpY1NlYXJjaABOBUQAXA5KZWt5bGwAaAVFCgojU3RhcnQgYwCBKgcuClVzZXIgLT4gQTogLi8AEgVlcgAZBiB3aGl0ZWxpc3QueW1sCmFjdGl2YXRlIEEKCiMgVXBkYXRlIElQQSBkYXRhLgpBIC0-IEI6IHUAEwZpcGEAFAUALQpCCkIAHAdkb3dubG9hZCBuZXcgZmlsZQpkZQAcCwojIEVTIG1hcHAAgRwFAFQFQzogR2VuZXJhdGUAEgsgKGlmIG5vdCBleGlzdHMpAIEaCkMKQyAtPiBEOiBjcmVhdGUgInNvZnR3YXJlIgBSCACBRwpEABwRYWRtaW5pc3RyYXRpb24AKQoAgR0LRACBKQxDCgojIExvYWQgQ29uZmlncwCBLQoAEQVjABAGOiBkb21haW5zAIEYEQAiByZQYXJzZQAeCACCcgUAWQ0APxYAgx8JUEEAOh8AHwsASxIKIyBSZXRyaWV2ZSBhbGwgdGhlIHJlcG9zaXRvcmllAIFJCwAeCQATDCAoUHJvY2Vzc1BBAIJ8DWxvb3AgRm9yIGVhY2ggb3JnYW5pegCCUwUKICAgIACBdwhEZXRlY3QgQVBJIChnaXRodWIsIGJpdGJ1Y2tldCwgZ2l0bGFiKQAmDQCBJQkASQwAgSYNAFUNRXh0cmFjdCBwdWJsaWNjb2RlLnltbCByYXcgVVJMAB4OQWRkIHRvAIFSDmNoYW5uZWwKZW5kAINrEQCBdAcAgjMFAIIeFgAaCACCHhVSAIJcCwCCIhsAgwUJeQCCLg1DaGVjawCBTAsgYXZhaWxhYmlsaXQAHA5WYWxpAIczBXJlbW90ZQCBfA8AgwMNU2F2ZQCHGAUgbWV0YQCHRgUACRYAKRtDbG9uAIQ5CwCBJw9hbGN1bGF0ZSAAiFUFaXR5IGluZGV4AIQMCkQAgH8HdG8gRVMgYXMgAIc8CACCbhVGbHVzaCBFUwCIKAoACggAiAMUACkFKACIJwsAh0MNAIdLEgCJaQdBbGlhAIdSCwBGHACKHAYAOCsAiVgJAIssBgCKTAxFAIl0CwAXB1lNTACLCQpFAIZGBUUAGxBhbQCJPwh6aW9uaQCDWwkAGBEAihUILXJpdXNvAAUjb3Blbi1zb3VyYwCEMAoAPRkAiTgRRQoAi2YMQSAKCg&s=rose)

[//]: # "Comment - Link to web diagram page: https://www.websequencediagrams.com/?lz=dGl0bGUgQ3Jhd2xpbmcKCiMgQWN0b3JzLgpwYXJ0aWNpcGFudCBVc2VyAAQNIkNtZCIgYXMgQQAJDklwYQASBUIAHw9yYXdsZXIALQVDADsORWxhc3RpY1NlYXJjaABOBUQAXA5KZWt5bGwAaAVFCgojU3RhcnQgYwCBKgcuClVzZXIgLT4gQTogLi8AEgVlcgAZBiB3aGl0ZWxpc3QueW1sCmFjdGl2YXRlIEEKCiMgVXBkYXRlIElQQSBkYXRhLgpBIC0-IEI6IHUAEwZpcGEAFAUALQpCCkIAHAdkb3dubG9hZCBuZXcgZmlsZQpkZQAcCwojIEVTIG1hcHAAgRwFAFQFQzogR2VuZXJhdGUAEgsgKGlmIG5vdCBleGlzdHMpAIEaCkMKQyAtPiBEOiBjcmVhdGUgInNvZnR3YXJlIgBSCACBRwpEABwRYWRtaW5pc3RyYXRpb24AKQoAgR0LRACBKQxDCgojIExvYWQgQ29uZmlncwCBLQoAEQVjABAGOiBkb21haW5zAIEYEQAiByZQYXJzZQAeCACCcgUAWQ0APxYAgx8JUEEAOh8AHwsASxIKIyBSZXRyaWV2ZSBhbGwgdGhlIHJlcG9zaXRvcmllAIFJCwAeCQATDCAoUHJvY2Vzc1BBAIJ8DWxvb3AgRm9yIGVhY2ggb3JnYW5pegCCUwUKICAgIACBdwhEZXRlY3QgQVBJIChnaXRodWIsIGJpdGJ1Y2tldCwgZ2l0bGFiKQAmDQCBJQkASQwAgSYNAFUNRXh0cmFjdCBwdWJsaWNjb2RlLnltbCByYXcgVVJMAB4OQWRkIHRvAIFSDmNoYW5uZWwKZW5kAINrEQCBdAcAgjMFAIIeFgAaCACCHhVSAIJcCwCCIhsAgwUJeQCCLg1DaGVjawCBTAsgYXZhaWxhYmlsaXQAHA5WYWxpAIczBXJlbW90ZQCBfA8AgwMNU2F2ZQCHGAUgbWV0YQCHRgUACRYAKRtDbG9uAIQ5CwCBJw9hbGN1bGF0ZSAAiFUFaXR5IGluZGV4AIQMCkQAgH8HdG8gRVMgYXMgAIc8CACCbhVGbHVzaCBFUwCIKAoACggAiAMUACkFKACIJwsAh0MNAIdLEgCJaQdBbGlhAIdSCwBGHACKHAYAOCsAiVgJAIssBgCKTAxFAIl0CwAXB1lNTACLCQpFAIZGBUUAGxBhbQCJPwh6aW9uaQCDWwkAGBEAihUILXJpdXNvAAUjb3Blbi1zb3VyYwCEMAoAPRkAiTgRRQoAi2YMQSAKCg&s=rose"

**Crawler execution**

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

2.  if `./ipa/amministrazioni.txt` is older than two days, update to latest version. (if some error occurs, log the error and use the actual `amministrazioni.txt`)
3.  connect to Elasticsearch node. If ko, log the error and stop the crawler. If ok set the mapping for publiccode and administrations.
4.  Read and parse `domains.yml` ([reference](references.md)) and `whitelist.yml` ([reference](references.md))
5.  initialize a repositories channel that will host all the repository infos.
6.  Initialize a goroutine that will process the main organization url (`crawler.ProcessPA`), detect the right host (`crawler.KnownHost`) and, using the right APIs, generate the right URL to the _publiccode.yml_ raw file (`crawler.ProcessPADomain`) and save it to the repositories channel.
7.  With another goroutine, process every repository in repositories channel (`crawler.ProcessRepositories`) checking the availability. If the file is available: validate (`crawler.validateRemoteFile`), save the publiccode.yml (`crawler.SaveToFile`) and the metadata returned from the APIs (`crawler.SaveToFile`) to a file, clone the repository (`crawler.CloneRepository`) locally, calculate the repo activity/vitalityIndex (`crawler.CalculateRepoActivity`) and save the data on ES (`crawler.SaveToES`).
8.  Flush elastic commits.
9.  Update elastic aliases with the alias `publiccode`
10. Generate Jekyll website data for the the OSS catalog.
