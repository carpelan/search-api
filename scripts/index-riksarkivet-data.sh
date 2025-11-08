#!/bin/bash
set -e

# Riksarkivet Data Indexing Script
# Based on actual Riksarkivet APIs:
# - Search API: https://data.riksarkivet.se/api/records
# - OAI-PMH: https://oai-pmh.riksarkivet.se/OAI
# - IIIF: https://lbiiif.riksarkivet.se
# - ALTO XML: https://sok.riksarkivet.se/dokument/alto

API_URL="${API_URL:-http://localhost:8080}"
SEARCH_QUERY="${SEARCH_QUERY:-*}"
MAX_RECORDS="${MAX_RECORDS:-100}"

echo "🔄 Indexing Riksarkivet data using Search API"
echo "API URL: $API_URL"
echo "Search Query: $SEARCH_QUERY"
echo "Max Records: $MAX_RECORDS"

# Function to index a document
index_document() {
    local doc_json="$1"

    curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "$doc_json" \
        "${API_URL}/api/search/index" > /dev/null

    if [ $? -eq 0 ]; then
        return 0
    else
        echo "Failed to index document"
        return 1
    fi
}

echo ""
echo "📡 Fetching data from Riksarkivet Search API..."

# Fetch from Riksarkivet Search API
RIKSARKIVET_API="https://data.riksarkivet.se/api/records"
RESPONSE=$(curl -s "${RIKSARKIVET_API}?query=${SEARCH_QUERY}&rows=${MAX_RECORDS}" || echo "{}")

# Check if we got valid JSON
if ! echo "$RESPONSE" | jq empty 2>/dev/null; then
    echo "⚠️  Invalid JSON response from Riksarkivet API"
    echo "Attempting OAI-PMH harvest instead..."

    # Fallback to OAI-PMH
    OAI_ENDPOINT="https://oai-pmh.riksarkivet.se/OAI"

    echo "📥 Harvesting metadata via OAI-PMH..."
    OAI_RESPONSE=$(curl -s "${OAI_ENDPOINT}?verb=ListRecords&metadataPrefix=oai_dc" || echo "")

    if [ -n "$OAI_RESPONSE" ]; then
        echo "✅ Received OAI-PMH response"
        echo "💾 Saving to oai-pmh-records.xml for manual processing"
        echo "$OAI_RESPONSE" > oai-pmh-records.xml
        echo ""
        echo "Note: OAI-PMH responses require XML parsing."
        echo "Consider using tools like:"
        echo "  - xmlstarlet: xmlstarlet sel -t -m '//record' -v 'metadata' oai-pmh-records.xml"
        echo "  - Python: import xml.etree.ElementTree as ET"
        echo ""
        exit 0
    fi

    echo "❌ Could not fetch data from Riksarkivet APIs"
    exit 1
fi

# Parse and index records
INDEX_COUNT=0
FAILED_COUNT=0

echo "📊 Processing records..."
echo ""

# Extract records from response (structure depends on actual API)
# This is a template - adjust based on actual API response structure
RECORD_COUNT=$(echo "$RESPONSE" | jq '.response.docs | length' 2>/dev/null || echo "0")

if [ "$RECORD_COUNT" -eq 0 ]; then
    echo "⚠️  No records found in response"
    echo "Response preview:"
    echo "$RESPONSE" | jq '.' | head -20
    exit 0
fi

echo "Found $RECORD_COUNT records"
echo ""

for i in $(seq 0 $((RECORD_COUNT - 1))); do
    # Extract record data (adjust field names based on actual API)
    RECORD=$(echo "$RESPONSE" | jq ".response.docs[$i]")

    # Extract fields (these are examples - adjust to actual API structure)
    ID=$(echo "$RECORD" | jq -r '.id // empty')
    REF_CODE=$(echo "$RECORD" | jq -r '.reference // .arkivRef // empty')
    TITLE=$(echo "$RECORD" | jq -r '.title // .label // empty')
    INSTITUTION=$(echo "$RECORD" | jq -r '.institution // "Riksarkivet"')
    DATE_START=$(echo "$RECORD" | jq -r '.dateStart // .fromDate // empty')
    DATE_END=$(echo "$RECORD" | jq -r '.dateEnd // .toDate // empty')
    FULL_TEXT=$(echo "$RECORD" | jq -r '.text // .content // empty' | head -c 5000)

    # Skip if no ID
    if [ -z "$ID" ] || [ "$ID" = "null" ]; then
        continue
    fi

    # Create document JSON matching our model
    DOC_JSON=$(cat <<EOF
{
    "id": "${ID}",
    "referenceCode": "${REF_CODE:-unknown}",
    "title": "${TITLE:-Untitled}",
    "description": "",
    "institution": "${INSTITUTION}",
    "institutionLocation": "Sweden",
    "collectionName": "Riksarkivet Collection",
    "dateStart": "${DATE_START:-2000-01-01T00:00:00Z}",
    "dateEnd": "${DATE_END:-2000-12-31T23:59:59Z}",
    "dateDisplay": "${DATE_START:-?}-${DATE_END:-?}",
    "pageNumbers": [],
    "totalPages": 0,
    "fullText": $(echo "$FULL_TEXT" | jq -Rs .),
    "textSnippets": [],
    "iiifManifestUrl": "https://lbiiif.riksarkivet.se/collection/arkiv/${ID}",
    "iiifImageUrl": "https://lbiiif.riksarkivet.se",
    "altoXmlUrl": "https://sok.riksarkivet.se/dokument/alto/${ID}",
    "viewerUrl": "https://sok.riksarkivet.se/bildvisning/${ID}",
    "documentType": "archive_document",
    "subjects": ["riksarkivet", "archive"],
    "language": "sv",
    "accessConditions": "public",
    "rights": "CC0",
    "indexedDate": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "modifiedDate": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "sourceApi": "search_api"
}
EOF
)

    # Index the document
    if index_document "$DOC_JSON"; then
        INDEX_COUNT=$((INDEX_COUNT + 1))
        echo "✓ Indexed: $TITLE (${REF_CODE:-$ID})"
    else
        FAILED_COUNT=$((FAILED_COUNT + 1))
        echo "✗ Failed: $TITLE (${REF_CODE:-$ID})"
    fi

    # Progress indicator
    if [ $((INDEX_COUNT % 10)) -eq 0 ] && [ $INDEX_COUNT -gt 0 ]; then
        echo "   Progress: $INDEX_COUNT indexed, $FAILED_COUNT failed"
    fi
done

echo ""
echo "════════════════════════════════════════"
echo "📈 Indexing Summary"
echo "════════════════════════════════════════"
echo "✅ Successfully indexed: $INDEX_COUNT"
echo "❌ Failed: $FAILED_COUNT"
echo "📊 Total processed: $((INDEX_COUNT + FAILED_COUNT))"
echo "════════════════════════════════════════"
echo ""
echo "🔍 You can now search the indexed data at:"
echo "   ${API_URL}/api/search/search"
echo ""
echo "Example search:"
echo "curl -X POST ${API_URL}/api/search/search \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -d '{\"query\": \"*:*\", \"rows\": 10}'"
