#!/bin/bash
#
# To insert a document to an index in elasticsearch
#

source config.sh

TODAY=$(date '+%Y%m%d')
YESTERDAY="20180416"
INDEX_TODAY="publiccode_$TODAY"
INDEX_YESTERDAY="publiccode_$YESTERDAY"
ALIAS="publiccode"

generate_document() {
  cat <<EOF
{
  "standard version": "0.0.1",
  "url": "https://github.com/publiccodenet/publiccode.yml",
  "upstream-url": ["https://github.com/publiccodenet/publiccode.yml"],

  "license": "AGPL-3.0-or-later",
  "main-copyright-owner": "City of Amsterdam and many contributors",
  "repo-owner": "City of Amsterday",

  "maintainance-type": "commercial",
  "maintainance-until": "2019-01-01",
  "technical-contacts": [
    {
      "name": "Frank Zappa",
      "email": "frank.zappa@example.com",
      "affiliation": "Comune di Reggio Emilia"
    }
  ],

  "name": "Medusa",
  "logo": "img/logo.jpg",
  "version": "1.0",
  "released": "2017-04-15",
  "platforms": ["web", "Linux"],
  "longdesc-en": "Very long description of this software, also split on multiple rows. You should note what the software is and why one should need it.",
  "longdesc-it": "Descrizione molto lunga di questo software, anche diviso su più righe. Si dovrebbe notare che cos'è il software e perché uno dovrebbe averne bisogno.",
  "shortdesc-en": "A really interesting software.",
  "shortdesc-it": "Un software davvero interessante.",
  "videos": [
    "https://youtube.com/xxxxxxxx"
  ],

  "scope": [
    "it",
    "en"
  ],
  "pa-type": [
    "city",
    "it-ag-lavoro"
  ],
  "category": [
    "it-anagrafe"
  ],
  "tags": [
    "city",
    "employee",
    "public"
  ],
  "used-by": [
    "Comune di Firenze",
    "Comune di Roma"
  ],
  "dependencies": [
    "Oracle 11.4",
    "MySQL"
  ],
  "dependencies-hardware": [
    "NFC Reader (chipset xxx)"
  ],
  "it-use-spid": "yes",
  "it-pagopa": "no",
  "suggest-name": {
    "input": ["Medusa"],
    "weight" : 1
  }
}
EOF
}

curl -u "$BASICAUTH" -X PUT "http://elasticsearch:9200/$INDEX_TODAY/software/1" -H 'Content-Type: application/json' -d"$(generate_document)"