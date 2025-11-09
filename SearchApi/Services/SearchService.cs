using SearchApi.Models;
using SolrNet;
using SolrNet.Commands.Parameters;
using System.Diagnostics;

namespace SearchApi.Services;

/// <summary>
/// Implementation of search service using Solr
/// </summary>
public class SearchService : ISearchService
{
    private readonly ISolrOperations<MetadataDocument> _solr;
    private readonly ILogger<SearchService> _logger;

    public SearchService(ISolrOperations<MetadataDocument> solr, ILogger<SearchService> logger)
    {
        _solr = solr;
        _logger = logger;
    }

    public async Task<SearchResponse> SearchAsync(SearchRequest request)
    {
        ArgumentNullException.ThrowIfNull(request);

        var sw = Stopwatch.StartNew();

        try
        {
            // Build query options
            var options = new QueryOptions
            {
                Rows = request.Rows,
                StartOrCursor = new StartOrCursor.Start(request.Start)
            };

            // Add sorting
            if (!string.IsNullOrEmpty(request.SortField))
            {
                options.OrderBy = new[]
                {
                    new SortOrder(request.SortField,
                        request.SortOrder.Equals("asc", StringComparison.OrdinalIgnoreCase)
                            ? Order.ASC
                            : Order.DESC)
                };
            }

            // Add filters
            if (request.Filters.Count > 0)
            {
                options.FilterQueries = new List<ISolrQuery>();
                foreach (var filter in request.Filters)
                {
                    options.FilterQueries.Add(new SolrQueryByField(filter.Key, filter.Value));
                }
            }

            // Execute search
            var results = await Task.Run(() => _solr.Query(request.Query, options)).ConfigureAwait(false);

            sw.Stop();

            _logger.LogInformation("Search completed in {ElapsedMs}ms, found {TotalResults} results",
                sw.ElapsedMilliseconds, results.NumFound);

            return new SearchResponse
            {
                TotalResults = results.NumFound,
                Results = results.ToList(),
                QueryTime = sw.Elapsed.TotalMilliseconds,
                Start = request.Start,
                Rows = request.Rows
            };
        }
        catch (SolrNet.Exceptions.SolrConnectionException ex)
        {
            _logger.LogError(ex, "Solr connection error executing query: {Query}", request.Query);
            throw;
        }
        catch (ArgumentException ex)
        {
            _logger.LogError(ex, "Invalid query syntax: {Query}", request.Query);
            throw;
        }
    }

    public async Task<MetadataDocument?> GetByIdAsync(string id)
    {
        ArgumentException.ThrowIfNullOrWhiteSpace(id);

        try
        {
            var results = await Task.Run(() => _solr.Query($"id:{id}")).ConfigureAwait(false);
            return results.FirstOrDefault();
        }
        catch (SolrNet.Exceptions.SolrConnectionException ex)
        {
            _logger.LogError(ex, "Solr connection error retrieving document by ID: {Id}", id);
            throw;
        }
    }

    public async Task<bool> IndexDocumentAsync(MetadataDocument document)
    {
        ArgumentNullException.ThrowIfNull(document);

        try
        {
            await Task.Run(() => _solr.Add(document)).ConfigureAwait(false);
            await Task.Run(() => _solr.Commit()).ConfigureAwait(false);
            _logger.LogInformation("Document indexed successfully: {Id}", document.Id);
            return true;
        }
        catch (SolrNet.Exceptions.SolrConnectionException ex)
        {
            _logger.LogError(ex, "Solr connection error indexing document: {Id}", document.Id);
            return false;
        }
        catch (InvalidOperationException ex)
        {
            _logger.LogError(ex, "Invalid operation while indexing document: {Id}", document.Id);
            return false;
        }
    }

    public async Task<bool> DeleteDocumentAsync(string id)
    {
        ArgumentException.ThrowIfNullOrWhiteSpace(id);

        try
        {
            await Task.Run(() => _solr.Delete(id)).ConfigureAwait(false);
            await Task.Run(() => _solr.Commit()).ConfigureAwait(false);
            _logger.LogInformation("Document deleted successfully: {Id}", id);
            return true;
        }
        catch (SolrNet.Exceptions.SolrConnectionException ex)
        {
            _logger.LogError(ex, "Solr connection error deleting document: {Id}", id);
            return false;
        }
    }
}
