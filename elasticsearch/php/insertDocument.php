<?php

require 'vendor/autoload.php';

use Elasticsearch\ClientBuilder;

require "config.inc";
require "generatorRandomDocuments.php";

$generator = new generatorRandomDocuments();
$documents = $generator->generateDocuments();

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

foreach ($documents as $key => $document) {
  $params = [
    'index' => $index,
    'type' => $type,
    'id' => $key,
    'body' => $document
  ];
  try {
    $client->index($params);
  }
  catch(Exception $e) {
    print_r($e->getMessage());
  }
}
