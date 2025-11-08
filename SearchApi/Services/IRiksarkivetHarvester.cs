using SearchApi.Models;

namespace SearchApi.Services;

/// <summary>
/// Service for harvesting metadata from Riksarkivet APIs
/// </summary>
public interface IRiksarkivetHarvester
{
    /// <summary>
    /// Harvest documents from Riksarkivet Search API
    /// </summary>
    Task<List<MetadataDocument>> HarvestFromSearchApiAsync(string query = "*", int maxRecords = 100);

    /// <summary>
    /// Harvest documents from OAI-PMH endpoint
    /// </summary>
    Task<List<MetadataDocument>> HarvestFromOaiPmhAsync(string set = "", int maxRecords = 100);

    /// <summary>
    /// Fetch ALTO XML transcription for a document
    /// </summary>
    Task<string?> FetchAltoXmlAsync(string documentId);

    /// <summary>
    /// Fetch IIIF manifest for a document
    /// </summary>
    Task<string?> FetchIiifManifestAsync(string documentId);
}
