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

	// Moderation: Add a table (example)
	// client.AddTable(workflowID, documentID, pageID, []int{10,10,100,100}, []string{"col1","col2"}, "unverified", "", []map[string]interface{}{{"row":0,"col":0,"header":"col1","text":"cell","bbox":[]int{10,10,20,20},"verification_status":"unverified","verification_message":""}})
}
