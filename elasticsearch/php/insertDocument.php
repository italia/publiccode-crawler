<?php

require 'vendor/autoload.php';

use Elasticsearch\ClientBuilder;

require "config.inc";
require "generatorRandomDocuments.php";

$hosts = [
  [
    'host' => $host,
    'port' => $post,
    'scheme' => $schema,
    'user' => $user,
    'pass' => $password
  ],
];

$client = ClientBuilder::create()->setHosts($hosts)->build();

create_indexes($client, $indexes);
foreach ($aliases as $a => $a_indexes) {
  create_aliases($client, $a, $a_indexes);
}
$posts = insert_jekyll_posts($client, 'jekyll', 'jekyll_data.json');
$documents = insert_publiccode_documents($client, $index, $type);

// insert_suggestions_agencies($client, 'suggestions', 'suggestion');
// insert_suggestions_software_type($client, 'suggestions', 'suggestion');

function create_indexes($client, $indexes) {
  foreach ($indexes as $index => $filename_mapping) {
    $mapping = file_get_contents($filename_mapping);
    $mapping = json_decode($mapping);

    $params = [
      'index' => $index,
      'body' => $mapping,
    ];

    $response = $client->indices()->create($params);
  }
}

function create_aliases($client, $alias, $indexes) {
  $params = [
    'body' => [
      'actions' => []
    ],
  ];

  foreach ($indexes as $index) {
    $params['body']['actions'][] = [
      'add' => [
        'index' => $index,
        'alias' => $alias,
      ]
    ];
  }

  return $client->indices()->updateAliases($params);
}

function insert_jekyll_posts($client, $index, $filename_data) {
  $data = file_get_contents($filename_data);
  $data = json_decode($data);
  $posts = [];

  foreach ($data->data as $post) {

    $params = [
      'index' => $index,
      'type' => $post->_type,
      'id' => $post->_id,
      'body' => $post->_source,
    ];
    $posts[] = $post->_source;

    $response = $client->index($params);
  }

  return $posts;
}

function insert_publiccode_documents($client, $index, $type) {
  $generator = new generatorRandomDocuments();
  $documents = $generator->generateDocuments();

  foreach ($documents as $key => $document) {
    $params = [
      'index' => $index,
      'type' => $type,
      'id' => $key,
      'body' => $document
    ];
    try {
      $response = $client->index($params);
    }
    catch(Exception $e) {
      print_r($e->getMessage() . "\n");
    }
  }

  return $documents;
}

function insert_suggestions_documents($client, $index, $type, $documents) {

  foreach ($documents as $key => $document) {
    $params = [
      'index' => $index,
      'type' => $type,
      'body' => get_suggestion_from_software_document($document),
    ];

    try {
      $response = $client->index($params);
    }
    catch(Exception $e) {
      print_r($e->getMessage() . "\n");
    }
  }
}

function insert_suggestions_posts($client, $index, $type, $posts) {

  foreach ($posts as $key => $post) {
    $params = [
      'index' => $index,
      'type' => $type,
      'body' => get_suggestion_from_post_document($post),
    ];
    try {
      $response = $client->index($params);
    }
    catch(Exception $e) {
      print_r($e->getMessage() . "\n");
    }
  }
}

function insert_suggestions_agencies($client, $index, $type) {
  $agencies = [];
  // get all agencies
  $params = [
    'index' => 'publiccode',
    'type' => 'software',
    'body' => [
      'aggs' => [
        'agencies' => [
          'terms' => [
            'field' =>'legal-mainCopyrightOwner',
            'size' => 1000
          ]
        ]
      ]
    ],
  ];

  $response = $client->search($params);
  foreach ($response['aggregations']['agencies']['buckets'] as $item) {
    $agencies[] = [
      'it' => [
        'suggest-agencies' => explode(' ', $item['key']),
      ],
      'en' => [
        'suggest-agencies' => explode(' ', $item['key']),
      ],
      'agency' => [
        'title' => $item['key']
      ]
    ]; 
  }

  foreach ($agencies as $agency) {
    $params = [
      'index' => $index,
      'type' => $type,
      'body' => $agency,
    ];
    try {
      $response = $client->index($params);
    }
    catch(Exception $e) {
      print_r($e->getMessage() . "\n");
    }
  }
}

function insert_suggestions_software_type($client, $index, $type) {
  $software_types = [];
  // get all agencies
  $params = [
    'index' => 'publiccode',
    'type' => 'software',
    'body' => [
      'aggs' => [
        'software-types-eng' => [
          'terms' => [
            'field' =>'description.eng.genericName.keyword',
            'size' => 1000
          ]
          ],
          'software-types-ita' => [
            'terms' => [
              'field' =>'description.ita.genericName.keyword',
              'size' => 1000
            ]
          ]
      ]
    ],
  ];

  $response = $client->search($params);

  foreach ($response['aggregations']['software-types-eng']['buckets'] as $item) {
    $software_types[] = [
      'en' => [
        'suggest-software-type' => $item['key'],
      ],
      'software_type' => [
        'title' => $item['key']
      ]
    ];
  }
  foreach ($response['aggregations']['software-types-ita']['buckets'] as $item) {
    $software_types[] = [
      'it' => [
        'suggest-software-type' => $item['key'],
      ],
      'software_type' => [
        'title' => $item['key']
      ]
    ];
  }

  foreach ($software_types as $software_type) {
    $params = [
      'index' => $index,
      'type' => $type,
      'body' => $software_type,
    ];
    try {
      $response = $client->index($params);
    }
    catch(Exception $e) {
      print_r($e->getMessage() . "\n");
    }
  }
}

/**
 * Create suggestion object from a software.
 */
function get_suggestion_from_software_document($document) {

  $suggestion = [
    // to have all information on frontend, include also the original software object.
    'software' => $document,
  ];

  foreach ($document['description'] as $lang => $description) {
    $suggest_all = generate_autocomplete_strings($document['name']);
    if (isset($description['localisedName'])) {
      $suggest_all = $suggest_all + generate_autocomplete_strings($description['localisedName']);
    }
    $suggestion[substr($lang, 0, 2)]['suggest-all'] = $suggest_all;
    
    if (empty($document['it-riuso-codiceIPA'])) {
      $suggestion[substr($lang, 0, 2)]['suggest-open-source'] = $suggest_all;
    }
    else {
      $suggestion[substr($lang, 0, 2)]['suggest-reuse-codeipa'] = $suggest_all;
    }
  }

  return $suggestion;
}

/**
 * Create suggestion object from a post.
 */
function get_suggestion_from_post_document($post) {
  return [
    // to have all information on frontend, include also the original post object.
    'post' => $post,
    // to have suggestion language specific, use $post->lang.
    $post->lang => [
      'suggest-all' => generate_autocomplete_strings($post->title),
      'suggest-platforms' => ($post->type == 'projects') ? generate_autocomplete_strings($post->title) : [],
      'suggest-api' => ($post->type == 'api') ? generate_autocomplete_strings($post->title) : [],
    ],
  ];
}

function generate_autocomplete_strings($text) {
  $strings = [];
  $pattern = "/[\w]+/";
  $matches = [];
  $n = preg_match_all($pattern, $text, $matches, PREG_OFFSET_CAPTURE);

  foreach ($matches[0] as $match) {
    $strings[] = substr($text, $match[1]);
  }

  return $strings;
}