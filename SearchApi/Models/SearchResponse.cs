namespace SearchApi.Models;

/// <summary>
/// Search response with results and metadata
/// </summary>
public class SearchResponse
{
    /// <summary>
    /// Total number of matching documents
    /// </summary>
    public long TotalResults { get; set; }

    /// <summary>
    /// Search results
    /// </summary>
    public IReadOnlyList<MetadataDocument> Results { get; set; } = Array.Empty<MetadataDocument>();

    /// <summary>
    /// Time taken for the search in milliseconds
    /// </summary>
    public double QueryTime { get; set; }

    /// <summary>
    /// Current page offset
    /// </summary>
    public int Start { get; set; }

    /// <summary>
    /// Number of results per page
    /// </summary>
    public int Rows { get; set; }
}
