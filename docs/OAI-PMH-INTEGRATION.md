# OAI-PMH Integration Guide

This project integrates with Riksarkivet's OAI-PMH service to harvest archive metadata.

## Riksarkivet OAI-PMH Service

**Base URL**: `https://oai-pmh.riksarkivet.se/OAI`

**Official Documentation**: https://github.com/Riksarkivet/dataplattform/wiki/OAI-PMH

## Supported Metadata Formats

| Prefix | Description |
|--------|-------------|
| `oai_ape_ead` | EAD XML for Archives Portal Europe |
| `oai_ra_ead` | EAD XML adapted for Riksarkivet |

## OAI-PMH Verbs

### 1. Identify
Get service information:
```bash
curl "https://oai-pmh.riksarkivet.se/OAI?verb=Identify"
```

### 2. ListMetadataFormats
List available metadata formats:
```bash
curl "https://oai-pmh.riksarkivet.se/OAI?verb=ListMetadataFormats"
```

### 3. ListAllAuth (Non-standard)
Enumerate all accessible collections/datasets:
```bash
curl "https://oai-pmh.riksarkivet.se/OAI?verb=ListAllAuth"
```

### 4. ListIdentifiers
Get identifiers for records in a dataset:
```bash
curl "https://oai-pmh.riksarkivet.se/OAI?verb=ListIdentifiers&metadataPrefix=oai_ra_ead&set=SE/ULA/10012"
```

With date filtering:
```bash
curl "https://oai-pmh.riksarkivet.se/OAI?verb=ListIdentifiers&metadataPrefix=oai_ra_ead&from=2024-01-01&until=2024-12-31"
```

### 5. ListRecords
Harvest full metadata records:
```bash
curl "https://oai-pmh.riksarkivet.se/OAI?verb=ListRecords&metadataPrefix=oai_ra_ead"
```

### 6. GetRecord
Fetch a specific record by identifier:
```bash
# Example identifier: SE/ULA/10012/A+I
# Note: Spaces must be URL-encoded as "+"
curl "https://oai-pmh.riksarkivet.se/OAI?verb=GetRecord&metadataPrefix=oai_ra_ead&identifier=SE/ULA/10012/A+I"
```

## Identifier Structure

Reference codes follow a hierarchical pattern:

- **Archive**: `SE/ULA/10012`
- **Series**: `SE/ULA/10012/A 1`
- **Volume**: `SE/ULA/10012/A 1/1`

**Important**: Spaces in identifiers must be URL-encoded with `+` characters.

## Integration with This Search API

### Option 1: Use Existing Riksarkivet Tools

Riksarkivet provides tools for working with their OAI-PMH data:
- See: https://github.com/Riksarkivet/dataplattform

### Option 2: Custom Integration

To harvest and index data into this search API:

1. **Harvest OAI-PMH records** using any OAI-PMH client
2. **Parse the EAD XML** to extract metadata
3. **POST to our API** at `/api/search/index`

Example workflow:
```bash
# 1. Harvest records from OAI-PMH
curl "https://oai-pmh.riksarkivet.se/OAI?verb=ListRecords&metadataPrefix=oai_ra_ead" > records.xml

# 2. Parse XML and transform to our JSON format
# (Use xmlstarlet, Python, or other XML tools)

# 3. Index into our search API
curl -X POST "http://localhost:8080/api/search/index" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "SE/ULA/10012/A+I",
    "referenceCode": "SE/ULA/10012/A I",
    "title": "Archive Title",
    "institution": "Riksarkivet",
    "documentType": "ead_archive",
    ...
  }'
```

### Option 3: Python Harvester Example

```python
import requests
from xml.etree import ElementTree as ET

# Harvest from OAI-PMH
response = requests.get(
    "https://oai-pmh.riksarkivet.se/OAI",
    params={
        "verb": "ListRecords",
        "metadataPrefix": "oai_ra_ead"
    }
)

# Parse XML
root = ET.fromstring(response.content)
ns = {
    'oai': 'http://www.openarchives.org/OAI/2.0/',
    'ead': 'urn:isbn:1-931666-22-9'
}

# Extract records
for record in root.findall('.//oai:record', ns):
    identifier = record.find('.//oai:identifier', ns).text
    # Parse EAD metadata...

    # Index to our API
    requests.post(
        "http://localhost:8080/api/search/index",
        json={
            "id": identifier,
            # ... other fields
        }
    )
```

## License

All Riksarkivet data is licensed under CC0 1.0 Universal.

## Resources

- **Riksarkivet Dataplattform**: https://github.com/Riksarkivet/dataplattform
- **OAI-PMH Wiki**: https://github.com/Riksarkivet/dataplattform/wiki/OAI-PMH
- **OAI-PMH Protocol**: https://www.openarchives.org/pmh/
- **EAD Standard**: https://www.loc.gov/ead/
