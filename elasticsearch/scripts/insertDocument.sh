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
  },
  "metadata-repo": {
    "scm": "git",
    "website": "",
    "has_wiki": false,
    "name": "BlowMeAway",
    "links": {
      "watchers": {
        "href": "https://api.bitbucket.org/2.0/repositories/blowmeaway/blowmeaway/watchers"
      },
      "branches": {
        "href": "https://api.bitbucket.org/2.0/repositories/blowmeaway/blowmeaway/refs/branches"
      },
      "tags": {
        "href": "https://api.bitbucket.org/2.0/repositories/blowmeaway/blowmeaway/refs/tags"
      },
      "commits": {
        "href": "https://api.bitbucket.org/2.0/repositories/blowmeaway/blowmeaway/commits"
      },
      "clone": [{
        "href": "https://bitb001@bitbucket.org/blowmeaway/blowmeaway.git",
        "name": "https"
      }, {
        "href": "git@bitbucket.org:blowmeaway/blowmeaway.git",
        "name": "ssh"
      }],
      "self": {
        "href": "https://api.bitbucket.org/2.0/repositories/blowmeaway/blowmeaway"
      },
      "source": {
        "href": "https://api.bitbucket.org/2.0/repositories/blowmeaway/blowmeaway/src"
      },
      "html": {
        "href": "https://bitbucket.org/blowmeaway/blowmeaway"
      },
      "avatar": {
        "href": "https://bitbucket.org/blowmeaway/blowmeaway/avatar/32/"
      },
      "hooks": {
        "href": "https://api.bitbucket.org/2.0/repositories/blowmeaway/blowmeaway/hooks"
      },
      "forks": {
        "href": "https://api.bitbucket.org/2.0/repositories/blowmeaway/blowmeaway/forks"
      },
      "downloads": {
        "href": "https://api.bitbucket.org/2.0/repositories/blowmeaway/blowmeaway/downloads"
      },
      "pullrequests": {
        "href": "https://api.bitbucket.org/2.0/repositories/blowmeaway/blowmeaway/pullrequests"
      }
    },
    "fork_policy": "allow_forks",
    "uuid": "{75b388a7-cef5-41ff-bdcc-5d4cda1f7569}",
    "language": "java",
    "created_on": "2014-11-08T08:26:52.583284+00:00",
    "mainbranch": {
      "type": "branch",
      "name": "master"
    },
    "full_name": "blowmeaway/blowmeaway",
    "has_issues": false,
    "owner": {
      "username": "blowmeaway",
      "display_name": "BlowMeAway",
      "type": "team",
      "uuid": "{35edbd3f-c06c-4361-92cb-2d2979688430}",
      "links": {
        "self": {
          "href": "https://api.bitbucket.org/2.0/teams/blowmeaway"
        },
        "html": {
          "href": "https://bitbucket.org/blowmeaway/"
        },
        "avatar": {
          "href": "https://bitbucket.org/account/blowmeaway/avatar/32/"
        }
      }
    },
    "updated_on": "2016-11-24T12:09:52.635503+00:00",
    "size": 5286793,
    "type": "repository",
    "slug": "blowmeaway",
    "is_private": false,
    "description": "app@night",
    "project": {
      "key": "PROJ",
      "type": "project",
      "uuid": "{2becd404-ef65-4a26-96f0-2b8c84d1d86b}",
      "links": {
        "self": {
          "href": "https://api.bitbucket.org/2.0/teams/blowmeaway/projects/PROJ"
        },
        "html": {
          "href": "https://bitbucket.org/account/user/blowmeaway/projects/PROJ"
        },
        "avatar": {
          "href": "https://bitbucket.org/account/user/blowmeaway/projects/PROJ/avatar/32"
        }
      },
      "name": "Untitled project"
    },
    "parent": {
      "links": {
        "self": {
          "href": ""
        },
        "html": {
          "href": ""
        },
        "avatar": {
          "href": ""
        }
      },
      "type": "",
      "name": "",
      "full_name": "",
      "uuid": ""
    }
  }
}
EOF
}

curl -u "$BASICAUTH" -X PUT "$ELASTICSEARCH_URL/$INDEX_TODAY/software/1" -H 'Content-Type: application/json' -d"$(generate_document)"