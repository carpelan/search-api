#!/bin/bash
set -e

# Script to download and index Riksarkivet metadata
# Data source: https://sok.riksarkivet.se/data-api/nedladdningsbara-datamangder/arkivmetadata/

API_URL="${API_URL:-http://localhost:5000}"
DATA_URL="${DATA_URL:-https://sok.riksarkivet.se/data-api/nedladdningsbara-datamangder/arkivmetadata/}"

echo "🔄 Fetching Riksarkivet metadata from: $DATA_URL"

# Download the metadata
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

echo "📥 Downloading metadata files..."
wget -r -np -nH --cut-dirs=3 -R "index.html*" "$DATA_URL" || true

# Process each JSON/XML file and index it
echo "📊 Processing and indexing files..."

INDEX_COUNT=0
for file in $(find . -type f \( -name "*.json" -o -name "*.xml" \)); do
    echo "Processing: $file"

    # Extract metadata from filename and content
    BASENAME=$(basename "$file")
    ID=$(echo "$BASENAME" | sed 's/\.[^.]*$//')

    # Create a document for indexing
    # This is a simple example - you'd want to parse the actual file content
    cat > "/tmp/document.json" <<EOF
{
    "id": "${ID}",
    "title": "${BASENAME}",
    "description": "Riksarkivet arkivmetadata",
    "author": "Riksarkivet",
    "createdDate": "$(date -Iseconds)",
    "modifiedDate": "$(date -Iseconds)",
    "tags": ["riksarkivet", "metadata", "archive"],
    "contentType": "application/json",
    "fileSize": $(stat -f%z "$file" 2>/dev/null || stat -c%s "$file"),
    "fullText": "$(head -c 1000 "$file" | tr '\n' ' ')"
}
EOF

    # Index the document
    curl -X POST \
        -H "Content-Type: application/json" \
        -d @/tmp/document.json \
        "${API_URL}/api/search/index" || echo "Failed to index $file"

    INDEX_COUNT=$((INDEX_COUNT + 1))

    if [ $((INDEX_COUNT % 10)) -eq 0 ]; then
        echo "Indexed $INDEX_COUNT documents so far..."
    fi
done

echo "✅ Indexing complete! Total documents indexed: $INDEX_COUNT"

# Cleanup
cd /
rm -rf "$TEMP_DIR"

echo "🔍 You can now search the indexed data!"
