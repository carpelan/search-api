using SearchApi.Models;

namespace SearchApi.Services;

/// <summary>
/// Service interface for search operations
/// </summary>
public interface ISearchService
{
    /// <summary>
    /// Search metadata documents
    /// </summary>
    Task<SearchResponse> SearchAsync(SearchRequest request);

    /// <summary>
    /// Get document by ID
    /// </summary>
    Task<MetadataDocument?> GetByIdAsync(string id);

    /// <summary>
    /// Index a new document
    /// </summary>
    Task<bool> IndexDocumentAsync(MetadataDocument document);

    /// <summary>
    /// Delete document by ID
    /// </summary>
    Task<bool> DeleteDocumentAsync(string id);
}
