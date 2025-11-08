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
        var sw = Stopwatch.StartNew();

        try
        {
            // Build query options
            var options = new QueryOptions
            {
                Rows = request.Rows,
                Start = request.Start
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
            foreach (var filter in request.Filters)
            {
                options.AddFilterQuery(new SolrQueryByField(filter.Key, filter.Value));
            }

            // Execute search
            var results = await Task.Run(() => _solr.Query(request.Query, options));

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
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error executing search query: {Query}", request.Query);
            throw;
        }
    }

    public async Task<MetadataDocument?> GetByIdAsync(string id)
    {
        try
        {
            var results = await Task.Run(() => _solr.Query($"id:{id}"));
            return results.FirstOrDefault();
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error retrieving document by ID: {Id}", id);
            throw;
        }
    }

    public async Task<bool> IndexDocumentAsync(MetadataDocument document)
    {
        try
        {
            await Task.Run(() => _solr.Add(document));
            await Task.Run(() => _solr.Commit());
            _logger.LogInformation("Document indexed successfully: {Id}", document.Id);
            return true;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error indexing document: {Id}", document.Id);
            return false;
        }
    }

    public async Task<bool> DeleteDocumentAsync(string id)
    {
        try
        {
            await Task.Run(() => _solr.Delete(id));
            await Task.Run(() => _solr.Commit());
            _logger.LogInformation("Document deleted successfully: {Id}", id);
            return true;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error deleting document: {Id}", id);
            return false;
        }
    }
}
