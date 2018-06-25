/**
 * Elasticsearch connection parameters are loaded by /assets/js/config.js
 */

$( document ).ready(function() {
  var client = new elasticsearch.Client({
    // log: 'trace'
    'host': host
  });

  // Get all tags terms from ES
  getAllFilterTerms(client);

  $('#es-automplete-input').on('input', null, client, executeAutoCompleteESQuery);
  $('#es-search-button').on('click', null, client, executeSearchCallback);
  $('#es-tags-list').on('click', null, client, activateTermFilter);
  $('#es-intended-audience-only-for-list').on('click', null, client, activateTermFilter);
  $('#es-term-active').on('click', null, client, deActivateTermFilter);
  $('input[name=sort-by-date]').on('change', null, client, onSortChange);
  
});

function executeAutoCompleteESQuery(event) {
  event.preventDefault();
  client = event.data;

  var language = $('#language').val();
  var field_autocomplete = $('#es-select-type-query').val();

  field_autocomplete = language.substring(0, 2) + '.' + field_autocomplete;

  /**
   * In Elasticsearch are defined the following fields in order to use for suggestions
   *  - (it|en).suggest-all           - index: suggestions
   *  - (it|en).suggest-platforms     - index: suggestions
   *  - (it|en).suggest-software-type - index  suggestions
   *  - (it|en).suggest-api           - index: suggestions
   *  - (it|en).suggest-agencies      - index: suggestions
   *  - (it|en).suggest-reuse-codeipa - index: suggestions
   *  - (it|en).suggest-open-source   - index: suggestions
   */

  var params = {
    'index': 'suggestions',
    'body': {
      'suggest': {
        'search_string': {
          'prefix': event.target.value,
          'completion': {
            'field' :  field_autocomplete,
            'size': 10
          }
        }
      }
    }
  };

  console.log(params.body);

  client.search(params).then(
    function(body){
      console.log(body);
      $('#es-automplete-results').text("");
      var search_string = body.suggest.search_string.pop();
      $.each(search_string.options, function(index, option){
        console.log(option);
        var title = '';
        // post - this can be moved into a function
        if (typeof option._source.post !== 'undefined') {
          title = option._source.post.title;
        }

        // software - this can be moved into a function
        if (typeof option._source.software !== 'undefined') {
          if (typeof option._source.software[language.substring(0, 2)] !== 'undefined' && typeof option._source.software[language.substring(0, 2)].localisedName !== 'undefined') {
            title =  option._source.software[language.substring(0, 2)].localisedName; 
          }
          else {
            title =  option._source.software.name; 
          }
        }

        // software - this can be moved into a function
        if (typeof option._source.agencies !== 'undefined') {
          title = option._source.agencies.title;
        }

        // software - this can be moved into a function
        if (typeof option._source.software_type !== 'undefined') {
          title = option._source.software_type.title;
        }

        $('#es-automplete-results').append("<div>" + title + "</div>" );
      });
    },
    function(error){console.log(error);}
  );
}

function executeSearchCallback(event) {
  event.preventDefault();
  client = event.data;
  executeSearchESQuery(client);
}

/**
 * Activate a term filter for the next search.
 * 
 * @param {*} event 
 */
function activateTermFilter(event) {
  event.preventDefault();
  client = event.data;

  console.log(event.target);
  $(event.target).appendTo('#es-term-active');
}

/**
 * Remove a term from activated filter section.
 *
 * @param {*} event 
 */
function deActivateTermFilter(event) {
  event.preventDefault();
  client = event.data;

  console.log(event.target);
  var term = $(event.target).attr('es-name');
  $(event.target).appendTo('#es-'+term+'-list');
}

/**
 * Build and execute a query toward elasticsearch. Write results on page.
 * 
 * @param {*} client 
 */
function executeSearchESQuery(client) {
  console.log("EXECUTE QUERY");

  var query = {
    aggs: {},
    
  };
  var filter = [];
  var sort = [];
  var language = $('#language').val();
  /*** execute full text query ***/

  // Add fields corresponding to the current frontend language.
  var must = {
    'multi_match': {
      'query': $('#es-search-input').val(),
      'fields': ['name', 'description.'+language+'.short-description', 'description.'+language+'.short-description', 'title', 'description']
    }
  };

  /*** execute query filtered by tag ***/

  // first, take tags selected
  var tags = [];
  var intended_audience_only_for = [];
  $('#es-term-active .es-term.tags').each(function(index, element){
    tags.push({
      value:$(element).attr('es-value'),
      name: $(element).attr('es-name')
    });
  });
  $('#es-term-active .es-term.intended-audience-only-for').each(function(index, element){
    intended_audience_only_for.push({
      value:$(element).attr('es-value'),
      name: $(element).attr('es-name')
    });
  });

  console.log("TAGS: ");
  console.log(tags);
  console.log("PATYPE: ");
  console.log(intended_audience_only_for);

  if (tags && tags.length) {
    console.log(tags);

    // filter have to be populated with all filters active
    // for AND query filtes use an distinct object, with term key, for each filter
    $.each(tags, function(index, t){
      var value = t.value;
      var name = t.name;
      term = {};
      term[name] = value;
      filter.push(
        {
          'term': term
        }
      );
    });

    // for OR query filtes use only one object with terms key
    // query.query = {
    //   bool: {
    //     filter: [
    //       {
    //         terms: {
    //           tags: tags
    //         }
    //       }
    //     ]
    //   }
    // };

    // bucket query, to include tags terms presents in the current search query results.
    // query.aggs = {
    //   'tags': {
    //     'filter': {
    //       'terms': {'tags': tags}
    //     },
    //     'aggs': {
    //       'tags': {
    //         'terms': {
    //           'field':'tags'
    //         }
    //       }
    //     }
    //   }
    // };
  }
  
  if (intended_audience_only_for && intended_audience_only_for.length) {
    console.log(intended_audience_only_for);

    // filter have to be populated with all filters active
    // for AND query filtes use an distinct object, with term key, for each filter
    $.each(intended_audience_only_for, function(index, t){
      var value = t.value;
      var name = t.name;
      term = {};
      term[name] = value;
      filter.push(
        {
          'term': term
        }
      );
    });

    // for OR query filtes use only one object with terms key
    // query.query = {
    //   bool: {
    //     filter: [
    //       {
    //         terms: {
    //           tags: tags
    //         }
    //       }
    //     ]
    //   }
    // };

    // bucket query
    // query.aggs = {
    //   'tags': {
    //     'filter': {
    //       'terms': {'tags': tags}
    //     },
    //     'aggs': {
    //       'tags': {
    //         'terms': {
    //           'field':'tags'
    //         }
    //       }
    //     }
    //   }
    // };
  }

  query.query = {
    'bool': {
      'filter': filter
    }
  };

  if ($('#es-search-input').val() != '') {
    query.query.bool.must = must;
  }

  // Sort
  if ($('input[name=sort-by-date]:checked').val() !== undefined) {
    sort.push({
      'release-date' : {'order' : $('input[name=sort-by-date]:checked').val() }
    });
  }

  query.sort = sort;
  console.log("EXECUTE THIS QUERY:");
  console.log(query);

  client.search({
    index: 'publiccode',
    body: query
  }).then(
    function(data){
      $('#es-results').text('');
      console.log(data);
      $.each(data.hits.hits, function(index, result){
        var title = (result._type == 'software') ? result._source.name : result._source.title;
        $('#es-results').append("<div class='es-result'>"+title+" ("+result._id+")</div>");
      });
    },
    function(error){
      $('#es-results').text('');      
      console.log(error);
    }
  );
}

/**
 * Gather all filter terms for: tags, pa-type, and write them on page.
 * 
 * @param {*} client 
 */
function getAllFilterTerms(client) {
  client.search({
    index: 'publiccode',
    // type: 'post',
    body: {
      aggs: {
        'tags': {
          terms: {
            field:'tags'
          }
        },
        'intended-audience-only-for': {
          terms: {
            field:'intended-audience-only-for'
          }
        }
      }
    }
  }).then(
    function(data){
      console.log(data);
      buckets = data.aggregations['tags'].buckets;
      $.each(buckets, function(index, bucket){
        $('#es-tags-list').append("<span class='es-term tags' es-value='"+bucket.key+"' es-name='tags'>" + bucket.key + " ("+bucket.doc_count+")</span>" );
      });

      buckets = data.aggregations['intended-audience-only-for'].buckets;
      $.each(buckets, function(index, bucket){
        $('#es-intended-audience-only-for-list').append("<span class='es-term intended-audience-only-for' es-value='"+bucket.key+"' es-name='intended-audience-only-for'>" + bucket.key + " ("+bucket.doc_count+")</span>" );
      });

    },
    function(error){console.log(error);}
  );
}

/**
 * Sort Results
 */

function onSortChange(event) {
  event.preventDefault();
  client = event.data;

  console.log("SORT CHANGE");
  executeSearchESQuery(client);
}
