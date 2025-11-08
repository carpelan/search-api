namespace SearchApi.Models;

/// <summary>
/// Search request parameters
/// </summary>
public class SearchRequest
{
    /// <summary>
    /// Search query string
    /// </summary>
    public string Query { get; set; } = "*:*";

    /// <summary>
    /// Number of results to return
    /// </summary>
    public int Rows { get; set; } = 10;

    /// <summary>
    /// Starting offset for pagination
    /// </summary>
    public int Start { get; set; } = 0;

    /// <summary>
    /// Field to sort by
    /// </summary>
    public string? SortField { get; set; }

    /// <summary>
    /// Sort order (asc or desc)
    /// </summary>
    public string SortOrder { get; set; } = "desc";

    /// <summary>
    /// Filters to apply
    /// </summary>
    public Dictionary<string, string> Filters { get; set; } = new();
}
