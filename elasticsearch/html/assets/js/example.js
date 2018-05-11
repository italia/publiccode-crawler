/**
 * Elasticsearch connection parameters are loaded by /assets/js/config.js
 */

$( document ).ready(function() {
  var client = new elasticsearch.Client({
    // log: 'trace'
    'host': host
  });

  $('#es-automplete-input').on('input', null, client, executeAutoCompleteESQuery);
});

function executeAutoCompleteESQuery(event) {
  event.preventDefault();
  client = event.data;
  client.search({
    index: 'publiccode',
    body: {
      suggest: {
        names: {
          prefix: event.target.value,
          completion: {
            field : "suggest-name",
            size: 10
          }
        }
      }
    }
  }).then(
    function(body){
      $('#es-automplete-results').text("");
      var names = body.suggest.names.pop();
      $.each(names.options, function(index, option){
        $('#es-automplete-results').append("<div>" + option._source.name + "</div>" );
      });
    },
    function(error){console.log(error);}
  );
}
