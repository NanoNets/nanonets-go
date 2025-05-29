# Nanonets Go SDK

A Go SDK for the Nanonets API, supporting document processing and moderation workflows.

## Installation

```bash
go get nanonets
```

## Usage

```go
package main

import (
    "fmt"
    "nanonets/nanonets"
)

func main() {
    client := nanonets.NewClient("YOUR_API_KEY")
    workflowID := "your_workflow_id"
    documentID := "your_document_id"
    pageID := "your_page_id"

    // List documents
    docs, _ := client.ListDocuments(workflowID, 1, 10)
    fmt.Println("Documents:", docs)

    // Get document
    doc, _ := client.GetDocument(workflowID, documentID)
    fmt.Println("Document:", doc)

    // Get page data
    page, _ := client.GetPageData(workflowID, documentID, pageID)
    fmt.Println("Page:", page)

    // Moderation: Update a field value (example)
    // client.UpdateFieldValue(workflowID, documentID, pageID, "field_data_id", "new_value")
}
```

## Available Methods

- Document Processing:
  - `UploadDocument`
  - `GetDocument`
  - `ListDocuments`
  - `DeleteDocument`
  - `GetDocumentFields`
  - `GetDocumentTables`
  - `GetDocumentOriginalFile`
  - `GetPageData`
- Moderation:
  - `UpdateFieldValue`
  - `AddFieldValue`
  - `DeleteFieldValue`
  - `AddTable`
  - `DeleteTable`
  - `UpdateTableCell`
  - `AddTableCell`
  - `DeleteTableCell`
  - `VerifyField`
  - `VerifyTableCell`
  - `VerifyTable`
  - `VerifyDocument`

For detailed request/response formats, see the [Go SDK docs](../docs/sdk/go-sdk/) and the [Postman collection](../postman/nanonets-document-processing.json). 