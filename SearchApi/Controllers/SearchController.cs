using Microsoft.AspNetCore.Mvc;
using SearchApi.Models;
using SearchApi.Services;

namespace SearchApi.Controllers;

/// <summary>
/// Search API controller
/// </summary>
[ApiController]
[Route("api/[controller]")]
[Produces("application/json")]
public class SearchController : ControllerBase
{
    private readonly ISearchService _searchService;
    private readonly ILogger<SearchController> _logger;

    public SearchController(ISearchService searchService, ILogger<SearchController> logger)
    {
        _searchService = searchService;
        _logger = logger;
    }

    /// <summary>
    /// Search metadata documents
    /// </summary>
    /// <param name="request">Search request parameters</param>
    /// <returns>Search results</returns>
    [HttpPost("search")]
    [ProducesResponseType(typeof(SearchResponse), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    public async Task<ActionResult<SearchResponse>> Search([FromBody] SearchRequest request)
    {
        ArgumentNullException.ThrowIfNull(request);

        try
        {
            var results = await _searchService.SearchAsync(request).ConfigureAwait(false);
            return Ok(results);
        }
        catch (SolrNet.Exceptions.SolrConnectionException ex)
        {
            _logger.LogError(ex, "Solr connection failed during search");
            return BadRequest(new { error = "Search service unavailable", message = ex.Message });
        }
        catch (ArgumentException ex)
        {
            _logger.LogError(ex, "Invalid search query");
            return BadRequest(new { error = "Invalid query", message = ex.Message });
        }
    }

    /// <summary>
    /// Get document by ID
    /// </summary>
    /// <param name="id">Document ID</param>
    /// <returns>Document details</returns>
    [HttpGet("{id}")]
    [ProducesResponseType(typeof(MetadataDocument), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<ActionResult<MetadataDocument>> GetById(string id)
    {
        ArgumentException.ThrowIfNullOrWhiteSpace(id);

        var document = await _searchService.GetByIdAsync(id).ConfigureAwait(false);
        if (document == null)
        {
            return NotFound();
        }
        return Ok(document);
    }

    /// <summary>
    /// Index a new document
    /// </summary>
    /// <param name="document">Document to index</param>
    /// <returns>Success status</returns>
    [HttpPost("index")]
    [ProducesResponseType(StatusCodes.Status201Created)]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    public async Task<IActionResult> IndexDocument([FromBody] MetadataDocument document)
    {
        ArgumentNullException.ThrowIfNull(document);

        var success = await _searchService.IndexDocumentAsync(document).ConfigureAwait(false);
        if (success)
        {
            return CreatedAtAction(nameof(GetById), new { id = document.Id }, document);
        }
        return BadRequest(new { error = "Failed to index document" });
    }

    /// <summary>
    /// Delete document by ID
    /// </summary>
    /// <param name="id">Document ID</param>
    /// <returns>Success status</returns>
    [HttpDelete("{id}")]
    [ProducesResponseType(StatusCodes.Status204NoContent)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<IActionResult> DeleteDocument(string id)
    {
        ArgumentException.ThrowIfNullOrWhiteSpace(id);

        var success = await _searchService.DeleteDocumentAsync(id).ConfigureAwait(false);
        if (success)
        {
            return NoContent();
        }
        return NotFound();
    }
}
