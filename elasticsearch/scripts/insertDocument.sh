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
  "publiccode-yaml-version": "http://w3id.org/publiccode/version/0.1",
  "name": "Medusa",
  "application-suite": "MegaProductivitySuite",
  "url": "https://example.com/italia/medusa.git",
  "landing-url": "https://example.com/italia/medusa",
  "is-based-on": ["https://github.com/italia/otello.git"],
  "software-version": "1.0",
  "release-date": "2017-04-15",
  "logo": "img/logo.svg",
  "monochrome-logo": "img/logo-mono.svg",
  "platforms": [
    "android",
    "ios"
  ],
  "tags": [
    "cms",
    "productivity",
    "it-portale-trasparenza"
  ],
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
