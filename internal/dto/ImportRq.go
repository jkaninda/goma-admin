package dto

// ImportResult holds the result of an import operation.
type ImportResult struct {
	Created int      `json:"created" description:"Number of newly created items"`
	Updated int      `json:"updated" description:"Number of updated items"`
	Errors  []string `json:"errors,omitempty" description:"List of errors encountered during import"`
}
