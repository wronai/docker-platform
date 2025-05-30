package utils

// Metadata represents file metadata
type Metadata map[string]interface{}

// ExtractMetadata extracts metadata from a file
func ExtractMetadata(filePath string) (Metadata, error) {
    // Implementation for metadata extraction
    return Metadata{}, nil
}

// CleanMetadata removes sensitive information from metadata
func CleanMetadata(meta Metadata) Metadata {
    // Implementation for metadata cleaning
    return meta
}
