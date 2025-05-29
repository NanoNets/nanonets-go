# Nanonets Go

[![Go Report Card](https://goreportcard.com/badge/github.com/NanoNets/nanonets-go)](https://goreportcard.com/report/github.com/NanoNets/nanonets-go)

A Go SDK for the Nanonets API, supporting document processing, workflow management, and moderation workflows.



---

**Nanonets** is an AI-powered Intelligent Document Processing platform that helps you:
- Extract structured data from invoices, receipts, forms, and more documents 
- Supports pdf, images (jpg, png, tiff), excel files, scanned documents and photos
- Automate data entry and document workflows
- Convert unstructured documents into machine-readable formats
- Integrate advanced OCR and table extraction into your apps

---

## Installation

```bash
go get github.com/NanoNets/nanonets-go/nanonets
```

## Authentication

Sign up and get your API key from your [Nanonets dashboard](https://app.nanonets.com/#/keys).

The SDK uses your API key for authentication. You can set it in two ways:

1. **Environment variable:**
   ```bash
   export NANONETS_API_KEY='your_api_key'
   ```
2. **Direct initialization:**
   ```go
   import "github.com/NanoNets/nanonets-go/nanonets"
   client := nanonets.NewClient("your_api_key")
   ```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/NanoNets/nanonets-go/nanonets"
)

func main() {
    client := nanonets.NewClient("YOUR_API_KEY")

    // Create a workflow
    workflow, err := client.Workflows.Create(nanonets.CreateWorkflowRequest{
        Description:  "SDK Example Workflow",
        WorkflowType: "", // Instant learning
    })
    if err != nil {
        fmt.Println("Error creating workflow:", err)
        return
    }
    workflowID := workflow.WorkflowID
    fmt.Println("Created workflow:", workflowID)

    // Configure fields and table headers
    fields := []nanonets.Field{
        {Name: "invoice_number"},
        {Name: "total_amount"},
        {Name: "invoice_date"},
    }
    tableHeaders := []nanonets.TableHeader{
        {Name: "item_description"},
        {Name: "quantity"},
        {Name: "unit_price"},
        {Name: "total"},
    }
    err = client.Workflows.SetFields(workflowID, nanonets.SetFieldsRequest{
        Fields: fields,
        TableHeaders: tableHeaders,
    })
    if err != nil {
        fmt.Println("Error setting fields:", err)
        return
    }
    fmt.Println("Configured fields and table headers.")

    // Upload a document from file
    uploadResult, err := client.Documents.Upload(workflowID, nanonets.UploadDocumentRequest{
        File:     "/path/to/document.pdf",
        Async:    false,
        Metadata: map[string]string{"test": "true"},
    })
    if err != nil {
        fmt.Println("Error uploading document:", err)
        return
    }
    fmt.Println("Upload result:", uploadResult)

    // Upload a document from URL
    uploadResultURL, err := client.Documents.UploadFromURL(workflowID, nanonets.UploadDocumentFromURLRequest{
        URL:      "https://example.com/document.pdf",
        Async:    false,
        Metadata: map[string]string{"test": "true"},
    })
    if err != nil {
        fmt.Println("Error uploading document from URL:", err)
        return
    }
    fmt.Println("Upload from URL result:", uploadResultURL)

    // List documents (paginated)
    documents, err := client.Documents.ListWithPagination(workflowID, 1, 10)
    if err != nil {
        fmt.Println("Error listing documents:", err)
        return
    }
    fmt.Println("Documents:", documents)

    // Get a document
    if len(documents) > 0 {
        docID := documents[0].DocumentID
        doc, err := client.Documents.Get(workflowID, docID)
        if err != nil {
            fmt.Println("Error getting document:", err)
        } else {
            fmt.Println("Document:", doc)
        }
    }
}
```

## Features

- **Workflow Management:** Create, list, get, set fields, update/delete fields, update metadata/settings, get types
- **Document Processing:** Upload (file/URL), list (paginated), get, delete, get fields/tables, get original file
- **Moderation:** Update/add/delete/verify fields, add/delete/update/verify tables and cells

## Error Handling

The SDK provides idiomatic Go error handling. Check errors returned from all SDK methods:

```go
workflow, err := client.Workflows.Create(nanonets.CreateWorkflowRequest{...})
if err != nil {
    // Handle error (authentication, validation, etc.)
    fmt.Println("Error:", err)
}
```

## Best Practices

1. **Resource Management**
   ```go
   // Use defer for cleanup if needed
   func processDocument(client *nanonets.Client, document string) error {
       defer func() {
           // Cleanup if needed
       }()
       return nil
   }
   ```
2. **Batch Processing**
   ```go
   // Process documents in batches
   func processBatch(client *nanonets.Client, documents []string, batchSize int) error {
       for i := 0; i < len(documents); i += batchSize {
           end := i + batchSize
           if end > len(documents) {
               end = len(documents)
           }
           batch := documents[i:end]
           // Process batch
       }
       return nil
   }
   ```
3. **Error Recovery**
   ```go
   import "github.com/cenkalti/backoff"

   // Retry with exponential backoff
   func processWithRetry(client *nanonets.Client, document string) error {
       operation := func() error {
           return processDocument(client, document)
       }
       return backoff.Retry(operation, backoff.NewExponentialBackOff())
   }
   ```

## Documentation & Support

- [Full Go SDK Documentation](../docs/sdk/go-sdk/)
- [API Reference](https://nanonets.com/documentation/)
- [Nanonets Dashboard](https://app.nanonets.com/)
- [Get your API Key](https://app.nanonets.com/#/keys)

---

For detailed request/response formats, see the [Go SDK docs](../docs/sdk/go-sdk/) and the [Postman collection](../postman/nanonets-document-processing.json). 