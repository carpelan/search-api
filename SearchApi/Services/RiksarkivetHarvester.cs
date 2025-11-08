using SearchApi.Models;
using System.Text.Json;
using System.Xml.Linq;

namespace SearchApi.Services;

/// <summary>
/// Implementation of Riksarkivet metadata harvester
/// Supports multiple Riksarkivet APIs:
/// - Search API: https://data.riksarkivet.se/api/records
/// - OAI-PMH: https://oai-pmh.riksarkivet.se/OAI
/// - IIIF: https://lbiiif.riksarkivet.se
/// - ALTO XML: https://sok.riksarkivet.se/dokument/alto
/// </summary>
public class RiksarkivetHarvester : IRiksarkivetHarvester
{
    private readonly IHttpClientFactory _httpClientFactory;
    private readonly ILogger<RiksarkivetHarvester> _logger;

    private const string SearchApiUrl = "https://data.riksarkivet.se/api/records";
    private const string OaiPmhUrl = "https://oai-pmh.riksarkivet.se/OAI";
    private const string IiifBaseUrl = "https://lbiiif.riksarkivet.se";
    private const string AltoXmlBaseUrl = "https://sok.riksarkivet.se/dokument/alto";

    public RiksarkivetHarvester(IHttpClientFactory httpClientFactory, ILogger<RiksarkivetHarvester> logger)
    {
        _httpClientFactory = httpClientFactory;
        _logger = logger;
    }

    public async Task<List<MetadataDocument>> HarvestFromSearchApiAsync(string query = "*", int maxRecords = 100)
    {
        var documents = new List<MetadataDocument>();

        try
        {
            var client = _httpClientFactory.CreateClient();
            var url = $"{SearchApiUrl}?query={Uri.EscapeDataString(query)}&rows={maxRecords}";

            _logger.LogInformation("Harvesting from Riksarkivet Search API: {Url}", url);

            var response = await client.GetStringAsync(url);
            var jsonDoc = JsonDocument.Parse(response);

            // Parse response based on actual API structure
            // Note: Adjust based on actual API response format
            if (jsonDoc.RootElement.TryGetProperty("response", out var responseElement) &&
                responseElement.TryGetProperty("docs", out var docsElement))
            {
                foreach (var doc in docsElement.EnumerateArray())
                {
                    var metadata = ParseSearchApiDocument(doc);
                    if (metadata != null)
                    {
                        documents.Add(metadata);
                    }
                }
            }

            _logger.LogInformation("Harvested {Count} documents from Search API", documents.Count);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error harvesting from Search API");
        }

        return documents;
    }

    public async Task<List<MetadataDocument>> HarvestFromOaiPmhAsync(string set = "", int maxRecords = 100)
    {
        var documents = new List<MetadataDocument>();

        try
        {
            var client = _httpClientFactory.CreateClient();
            var url = $"{OaiPmhUrl}?verb=ListRecords&metadataPrefix=oai_dc";
            if (!string.IsNullOrEmpty(set))
            {
                url += $"&set={Uri.EscapeDataString(set)}";
            }

            _logger.LogInformation("Harvesting from OAI-PMH: {Url}", url);

            var response = await client.GetStringAsync(url);
            var xml = XDocument.Parse(response);

            // Parse OAI-PMH response
            XNamespace oai = "http://www.openarchives.org/OAI/2.0/";
            XNamespace dc = "http://purl.org/dc/elements/1.1/";

            var records = xml.Descendants(oai + "record");
            int count = 0;

            foreach (var record in records)
            {
                if (count >= maxRecords)
                    break;

                var metadata = ParseOaiPmhRecord(record, oai, dc);
                if (metadata != null)
                {
                    documents.Add(metadata);
                    count++;
                }
            }

            _logger.LogInformation("Harvested {Count} documents from OAI-PMH", documents.Count);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error harvesting from OAI-PMH");
        }

        return documents;
    }

    public async Task<string?> FetchAltoXmlAsync(string documentId)
    {
        try
        {
            var client = _httpClientFactory.CreateClient();
            var url = $"{AltoXmlBaseUrl}/{documentId}";
            return await client.GetStringAsync(url);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error fetching ALTO XML for document {DocumentId}", documentId);
            return null;
        }
    }

    public async Task<string?> FetchIiifManifestAsync(string documentId)
    {
        try
        {
            var client = _httpClientFactory.CreateClient();
            var url = $"{IiifBaseUrl}/collection/arkiv/{documentId}";
            return await client.GetStringAsync(url);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error fetching IIIF manifest for document {DocumentId}", documentId);
            return null;
        }
    }

    private MetadataDocument? ParseSearchApiDocument(JsonElement doc)
    {
        try
        {
            var id = doc.TryGetProperty("id", out var idProp) ? idProp.GetString() : null;
            if (string.IsNullOrEmpty(id))
                return null;

            return new MetadataDocument
            {
                Id = id,
                ReferenceCode = GetStringProperty(doc, "reference", "arkivRef") ?? "unknown",
                Title = GetStringProperty(doc, "title", "label") ?? "Untitled",
                Description = GetStringProperty(doc, "description", "content") ?? "",
                Institution = GetStringProperty(doc, "institution") ?? "Riksarkivet",
                InstitutionLocation = "Sweden",
                CollectionName = GetStringProperty(doc, "collection") ?? "Riksarkivet Collection",
                DateDisplay = GetStringProperty(doc, "dateRange", "dateDisplay") ?? "",
                FullText = GetStringProperty(doc, "text", "fullText") ?? "",
                IiifManifestUrl = $"{IiifBaseUrl}/collection/arkiv/{id}",
                IiifImageUrl = IiifBaseUrl,
                AltoXmlUrl = $"{AltoXmlBaseUrl}/{id}",
                ViewerUrl = $"https://sok.riksarkivet.se/bildvisning/{id}",
                DocumentType = GetStringProperty(doc, "type") ?? "archive_document",
                Language = "sv",
                Rights = "CC0",
                SourceApi = "search_api",
                IndexedDate = DateTime.UtcNow,
                ModifiedDate = DateTime.UtcNow
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error parsing Search API document");
            return null;
        }
    }

    private MetadataDocument? ParseOaiPmhRecord(XElement record, XNamespace oai, XNamespace dc)
    {
        try
        {
            var header = record.Element(oai + "header");
            var metadata = record.Element(oai + "metadata");

            if (header == null || metadata == null)
                return null;

            var identifier = header.Element(oai + "identifier")?.Value;
            if (string.IsNullOrEmpty(identifier))
                return null;

            var dcMetadata = metadata.Elements().FirstOrDefault()?.Elements();
            if (dcMetadata == null)
                return null;

            return new MetadataDocument
            {
                Id = identifier,
                ReferenceCode = GetDcElement(dcMetadata, dc, "identifier") ?? identifier,
                Title = GetDcElement(dcMetadata, dc, "title") ?? "Untitled",
                Description = GetDcElement(dcMetadata, dc, "description") ?? "",
                Institution = GetDcElement(dcMetadata, dc, "publisher") ?? "Riksarkivet",
                InstitutionLocation = "Sweden",
                CollectionName = GetDcElement(dcMetadata, dc, "source") ?? "Riksarkivet Collection",
                DateDisplay = GetDcElement(dcMetadata, dc, "date") ?? "",
                Language = GetDcElement(dcMetadata, dc, "language") ?? "sv",
                Rights = GetDcElement(dcMetadata, dc, "rights") ?? "CC0",
                DocumentType = GetDcElement(dcMetadata, dc, "type") ?? "archive_document",
                SourceApi = "oai_pmh",
                IndexedDate = DateTime.UtcNow,
                ModifiedDate = DateTime.UtcNow
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error parsing OAI-PMH record");
            return null;
        }
    }

    private static string? GetStringProperty(JsonElement element, params string[] propertyNames)
    {
        foreach (var name in propertyNames)
        {
            if (element.TryGetProperty(name, out var prop) && prop.ValueKind == JsonValueKind.String)
            {
                return prop.GetString();
            }
        }
        return null;
    }

    private static string? GetDcElement(IEnumerable<XElement> elements, XNamespace dc, string name)
    {
        return elements.FirstOrDefault(e => e.Name == dc + name)?.Value;
    }
}
