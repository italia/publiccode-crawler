## File and folders description

###Generic file and folders

.circleci: CircleCI folder

.git: git folder

cmd, crawler, httpclient, ipa, jekyll, metrics, vendor, version: golang packages and crawler source.

./data: the tree for publiccode.yml files:

- ./data/<host>/<org>/<repo>/<index>\_publiccode.yml.
- ./data/<host>/<org>/<repo>/<index>\_metadata_publiccode.yml.
- ./data/<host>/<org>/<repo>/gitClone/<git_clone_files>

./docker: the tree for docker config files

./elasticsearch: the tree for elasticsearch mappings and example files

vitality-ranges.yml: ranges for vitality index (Ref: vitality-ranges.md)

whitelistGeneric.yml and whitelistPA.yml: whitelist for Public Administrations and generic organizations (Ref: whitelist.md)

###Configuration files

- Docker and crawler building configurations: .env
- Crawler specific configurations: config.toml
- Crawler domains basic auth: domains.yml
- Crawler vitality-ranges: vitality-ranges.yml
- Crawler whitelist: whitelistGeneric.yml and whitelistPA.yml
