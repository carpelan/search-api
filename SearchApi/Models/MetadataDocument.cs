using SolrNet.Attributes;

namespace SearchApi.Models;

/// <summary>
/// Represents a Riksarkivet archive document in Solr
/// Based on real Riksarkivet API structure
/// </summary>
public class MetadataDocument
{
    /// <summary>
    /// Unique identifier for the document
    /// </summary>
    [SolrUniqueKey("id")]
    public string Id { get; set; } = string.Empty;

    /// <summary>
    /// Reference code in Riksarkivet format (e.g., SE/RA/310187/1)
    /// </summary>
    [SolrField("reference_code")]
    public string ReferenceCode { get; set; } = string.Empty;

    /// <summary>
    /// Document title or summary
    /// </summary>
    [SolrField("title")]
    public string Title { get; set; } = string.Empty;

    /// <summary>
    /// Detailed description of the archive material
    /// </summary>
    [SolrField("description")]
    public string Description { get; set; } = string.Empty;

    /// <summary>
    /// Institution holding the archive (e.g., "Riksarkivet")
    /// </summary>
    [SolrField("institution")]
    public string Institution { get; set; } = string.Empty;

    /// <summary>
    /// Institution location
    /// </summary>
    [SolrField("institution_location")]
    public string InstitutionLocation { get; set; } = string.Empty;

    /// <summary>
    /// Archive collection name
    /// </summary>
    [SolrField("collection_name")]
    public string CollectionName { get; set; } = string.Empty;

    /// <summary>
    /// Start date of the archive material
    /// </summary>
    [SolrField("date_start")]
    public DateTime? DateStart { get; set; }

    /// <summary>
    /// End date of the archive material
    /// </summary>
    [SolrField("date_end")]
    public DateTime? DateEnd { get; set; }

    /// <summary>
    /// Date display string (e.g., "1750-1800")
    /// </summary>
    [SolrField("date_display")]
    public string DateDisplay { get; set; } = string.Empty;

    /// <summary>
    /// Page number(s) where content was found
    /// </summary>
    [SolrField("page_numbers")]
    public List<int> PageNumbers { get; set; } = new();

    /// <summary>
    /// Total page count in the document
    /// </summary>
    [SolrField("total_pages")]
    public int TotalPages { get; set; }

    /// <summary>
    /// Full text content (from ALTO XML transcriptions)
    /// </summary>
    [SolrField("full_text")]
    public string FullText { get; set; } = string.Empty;

    /// <summary>
    /// Text snippets with search highlights
    /// </summary>
    [SolrField("text_snippets")]
    public List<string> TextSnippets { get; set; } = new();

    /// <summary>
    /// IIIF manifest URL for the document
    /// </summary>
    [SolrField("iiif_manifest_url")]
    public string IiifManifestUrl { get; set; } = string.Empty;

    /// <summary>
    /// IIIF image base URL
    /// </summary>
    [SolrField("iiif_image_url")]
    public string IiifImageUrl { get; set; } = string.Empty;

    /// <summary>
    /// ALTO XML URL for transcription data
    /// </summary>
    [SolrField("alto_xml_url")]
    public string AltoXmlUrl { get; set; } = string.Empty;

    /// <summary>
    /// Interactive viewer URL (Bildvisning)
    /// </summary>
    [SolrField("viewer_url")]
    public string ViewerUrl { get; set; } = string.Empty;

    /// <summary>
    /// Document type (e.g., "manuscript", "church_record", "photograph")
    /// </summary>
    [SolrField("document_type")]
    public string DocumentType { get; set; } = string.Empty;

    /// <summary>
    /// Subject tags or categories
    /// </summary>
    [SolrField("subjects")]
    public List<string> Subjects { get; set; } = new();

    /// <summary>
    /// Language of the document
    /// </summary>
    [SolrField("language")]
    public string Language { get; set; } = "sv"; // Default to Swedish

    /// <summary>
    /// Access conditions (public, restricted, etc.)
    /// </summary>
    [SolrField("access_conditions")]
    public string AccessConditions { get; set; } = string.Empty;

    /// <summary>
    /// Rights/license (e.g., "CC0")
    /// </summary>
    [SolrField("rights")]
    public string Rights { get; set; } = "CC0";

    /// <summary>
    /// Date the record was created in the index
    /// </summary>
    [SolrField("indexed_date")]
    public DateTime IndexedDate { get; set; } = DateTime.UtcNow;

    /// <summary>
    /// Date the record was last modified in the index
    /// </summary>
    [SolrField("modified_date")]
    public DateTime ModifiedDate { get; set; } = DateTime.UtcNow;

    /// <summary>
    /// Source API (e.g., "search_api", "oai_pmh")
    /// </summary>
    [SolrField("source_api")]
    public string SourceApi { get; set; } = string.Empty;
}
