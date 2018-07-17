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
$administrations = insert_administrations_administration($client, 'administrations', 'administrations.csv');

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

function insert_administrations_administration($client, $index, $filename_data) {
  $handle = fopen($filename_data, "r");
  // the first row with headers.
  fgetcsv($handle, 1024, ",");
  while (($data = fgetcsv($handle, 1024, ",")) !== FALSE) {

    $params = [
      'index' => $index,
      'type' => 'administration',
      'body' => [
        'it-riuso-codiceIPA' => $data[0],
        'it-riuso-codiceIPA-label' => $data[1],
      ],
    ];
    $administrations[] = $params['body'];

    $response = $client->index($params);
  }

  return $administrations;
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