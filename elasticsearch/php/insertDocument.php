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
create_aliases($client, $alias, $indexes);
$posts = insert_jekyll_posts($client, 'jekyll', 'jekyll_data.json');
$documents = insert_publiccode_documents($client, $index, $type);
insert_suggestions_documents($client, 'suggestions', 'suggestion', $documents);
insert_suggestions_posts($client, 'suggestions', 'suggestion', $posts);

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

  foreach ($indexes as $index => $filename_mapping) {
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

function get_suggestion_from_software_document($document) {

  $suggestion = [
    'title' => $document['name'],
    'suggest-all' => explode(' ', $document['name']),
    // 'suggest-platforms' => [],
    'suggest-software-type' => $document['tags'],
    // 'suggest-api' => [],
    // 'suggest-agencies' => explode(' ', $document['legal-main-copyright-owner']),
    'suggest-agencies' => $document['legal-main-copyright-owner'],
  ];

  foreach ($document['description'] as $lang => $description) {
    $suggest_all = $document['name'];
    if (isset($description['localised-name'])) {
      $suggest_all .= ' ' . $description['localised-name'];
      $suggestion['title-' . substr($lang, 0, 2)] = $description['localised-name'];
    }

    $suggestion['suggest-all-' . substr($lang, 0, 2)] = explode(' ', $suggest_all);
  }

  return $suggestion;
}

function get_suggestion_from_post_document($post) {
  return [
    'title' => $post->title,
    // to have suggestion language specific, use $post['lang'].
    'suggest-all' => explode(' ', $post->title),
    // 'suggest-all-' . $post->lang => explode(' ', $suggest_all),
    'suggest-all-' . $post->lang => explode(' ', $post->title),
    // 'suggest-platforms-' . $post->lang => ($post->type == 'projects') ? explode(' ', $suggest_all) : [],
    'suggest-platforms-' . $post->lang => ($post->type == 'projects') ? $post->title : [],
    // 'suggest-api-' . $post->lang => ($post->type == 'api') ? explode(' ', $suggest_all) : [],
    'suggest-api-' . $post->lang => ($post->type == 'api') ? $post->title : [],
  ];
}
