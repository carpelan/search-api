using SolrNet.Attributes;

namespace SearchApi.Models;

/// <summary>
/// Represents a metadata document in Solr
/// </summary>
public class MetadataDocument
{
    [SolrUniqueKey("id")]
    public string Id { get; set; } = string.Empty;

    [SolrField("title")]
    public string Title { get; set; } = string.Empty;

    [SolrField("description")]
    public string Description { get; set; } = string.Empty;

    [SolrField("author")]
    public string Author { get; set; } = string.Empty;

    [SolrField("created_date")]
    public DateTime CreatedDate { get; set; }

    [SolrField("modified_date")]
    public DateTime ModifiedDate { get; set; }

    [SolrField("tags")]
    public List<string> Tags { get; set; } = new();

    [SolrField("content_type")]
    public string ContentType { get; set; } = string.Empty;

    [SolrField("file_size")]
    public long FileSize { get; set; }

    [SolrField("full_text")]
    public string FullText { get; set; } = string.Empty;
}
