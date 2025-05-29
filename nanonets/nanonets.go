package nanonets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// Client represents a Nanonets API client
type Client struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

// NewClient creates a new Nanonets API client
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: "https://app.nanonets.com/api/v4",
		Client:  &http.Client{},
	}
}

// Workflows represents the workflows API
type Workflows struct {
	client *Client
}

// Create creates a new workflow
func (w *Workflows) Create(req CreateWorkflowRequest) (*Workflow, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := w.client.Client.Post(
		fmt.Sprintf("%s/workflows", w.client.BaseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var workflow Workflow
	if err := json.NewDecoder(resp.Body).Decode(&workflow); err != nil {
		return nil, err
	}
	return &workflow, nil
}

// Get retrieves a workflow by ID
func (w *Workflows) Get(workflowID string) (*Workflow, error) {
	resp, err := w.client.Client.Get(fmt.Sprintf("%s/workflows/%s", w.client.BaseURL, workflowID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var workflow Workflow
	if err := json.NewDecoder(resp.Body).Decode(&workflow); err != nil {
		return nil, err
	}
	return &workflow, nil
}

// List retrieves all workflows
func (w *Workflows) List() ([]Workflow, error) {
	resp, err := w.client.Client.Get(fmt.Sprintf("%s/workflows", w.client.BaseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var workflows []Workflow
	if err := json.NewDecoder(resp.Body).Decode(&workflows); err != nil {
		return nil, err
	}
	return workflows, nil
}

// SetFields sets fields and table headers for a workflow
func (w *Workflows) SetFields(workflowID string, req SetFieldsRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/workflows/%s/fields", w.client.BaseURL, workflowID), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// UpdateField updates a field in a workflow
func (w *Workflows) UpdateField(workflowID, fieldID string, req UpdateFieldRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/workflows/%s/fields/%s", w.client.BaseURL, workflowID, fieldID), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// DeleteField deletes a field from a workflow
func (w *Workflows) DeleteField(workflowID, fieldID string) error {
	httpReq, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/workflows/%s/fields/%s", w.client.BaseURL, workflowID, fieldID), nil)
	if err != nil {
		return err
	}

	resp, err := w.client.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// UpdateMetadata updates metadata for a workflow
func (w *Workflows) UpdateMetadata(workflowID string, req UpdateMetadataRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/workflows/%s", w.client.BaseURL, workflowID), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// UpdateSettings updates settings for a workflow
func (w *Workflows) UpdateSettings(workflowID string, req UpdateSettingsRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/workflows/%s/settings", w.client.BaseURL, workflowID), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// GetTypes retrieves available workflow types
func (w *Workflows) GetTypes() ([]WorkflowType, error) {
	resp, err := w.client.Client.Get(fmt.Sprintf("%s/workflows/types", w.client.BaseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var types []WorkflowType
	if err := json.NewDecoder(resp.Body).Decode(&types); err != nil {
		return nil, err
	}
	return types, nil
}

// Documents represents the documents API
type Documents struct {
	client *Client
}

// Upload uploads a document to a workflow
func (d *Documents) Upload(workflowID string, req UploadDocumentRequest) (*UploadResult, error) {
	file, err := os.Open(req.File)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(req.File))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	writer.WriteField("async", fmt.Sprintf("%v", req.Async))
	for key, value := range req.Metadata {
		writer.WriteField(key, value)
	}
	writer.Close()

	resp, err := d.client.Client.Post(
		fmt.Sprintf("%s/workflows/%s/documents", d.client.BaseURL, workflowID),
		writer.FormDataContentType(),
		body,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result UploadResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves a document by ID
func (d *Documents) Get(workflowID, documentID string) (*Document, error) {
	resp, err := d.client.Client.Get(fmt.Sprintf("%s/workflows/%s/documents/%s", d.client.BaseURL, workflowID, documentID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var document Document
	if err := json.NewDecoder(resp.Body).Decode(&document); err != nil {
		return nil, err
	}
	return &document, nil
}

// List retrieves all documents for a workflow
func (d *Documents) List(workflowID string) ([]Document, error) {
	resp, err := d.client.Client.Get(fmt.Sprintf("%s/workflows/%s/documents", d.client.BaseURL, workflowID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var documents []Document
	if err := json.NewDecoder(resp.Body).Decode(&documents); err != nil {
		return nil, err
	}
	return documents, nil
}

// Delete deletes a document
func (d *Documents) Delete(workflowID, documentID string) error {
	httpReq, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/workflows/%s/documents/%s", d.client.BaseURL, workflowID, documentID), nil)
	if err != nil {
		return err
	}

	resp, err := d.client.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// GetFields retrieves fields for a document
func (d *Documents) GetFields(workflowID, documentID string) ([]Field, error) {
	resp, err := d.client.Client.Get(fmt.Sprintf("%s/workflows/%s/documents/%s/fields", d.client.BaseURL, workflowID, documentID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var fields []Field
	if err := json.NewDecoder(resp.Body).Decode(&fields); err != nil {
		return nil, err
	}
	return fields, nil
}

// GetTables retrieves tables for a document
func (d *Documents) GetTables(workflowID, documentID string) ([]Table, error) {
	resp, err := d.client.Client.Get(fmt.Sprintf("%s/workflows/%s/documents/%s/tables", d.client.BaseURL, workflowID, documentID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tables []Table
	if err := json.NewDecoder(resp.Body).Decode(&tables); err != nil {
		return nil, err
	}
	return tables, nil
}

// Moderation represents the moderation API
type Moderation struct {
	client *Client
}

// UpdateField updates a field value
func (m *Moderation) UpdateField(workflowID, documentID, pageID, fieldDataID string, req UpdateFieldRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/workflows/%s/documents/%s/pages/%s/fields/%s", m.client.BaseURL, workflowID, documentID, pageID, fieldDataID), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// AddField adds a field value
func (m *Moderation) AddField(workflowID, documentID, pageID string, req AddFieldRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/workflows/%s/documents/%s/pages/%s/fields", m.client.BaseURL, workflowID, documentID, pageID), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// DeleteField deletes a field value
func (m *Moderation) DeleteField(workflowID, documentID, pageID, fieldDataID string) error {
	httpReq, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/workflows/%s/documents/%s/pages/%s/fields/%s", m.client.BaseURL, workflowID, documentID, pageID, fieldDataID), nil)
	if err != nil {
		return err
	}

	resp, err := m.client.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// AddTable adds a table
func (m *Moderation) AddTable(workflowID, documentID, pageID string, req AddTableRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/workflows/%s/documents/%s/pages/%s/tables", m.client.BaseURL, workflowID, documentID, pageID), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// DeleteTable deletes a table
func (m *Moderation) DeleteTable(workflowID, documentID, pageID, tableID string) error {
	httpReq, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/workflows/%s/documents/%s/pages/%s/tables/%s", m.client.BaseURL, workflowID, documentID, pageID, tableID), nil)
	if err != nil {
		return err
	}

	resp, err := m.client.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// UpdateTableCell updates a table cell
func (m *Moderation) UpdateTableCell(workflowID, documentID, pageID, tableID, cellID string, req UpdateTableCellRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/workflows/%s/documents/%s/pages/%s/tables/%s/cells/%s", m.client.BaseURL, workflowID, documentID, pageID, tableID, cellID), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// AddTableCell adds a table cell
func (m *Moderation) AddTableCell(workflowID, documentID, pageID, tableID string, req AddTableCellRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/workflows/%s/documents/%s/pages/%s/tables/%s/cells", m.client.BaseURL, workflowID, documentID, pageID, tableID), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// DeleteTableCell deletes a table cell
func (m *Moderation) DeleteTableCell(workflowID, documentID, pageID, tableID, cellID string) error {
	httpReq, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/workflows/%s/documents/%s/pages/%s/tables/%s/cells/%s", m.client.BaseURL, workflowID, documentID, pageID, tableID, cellID), nil)
	if err != nil {
		return err
	}

	resp, err := m.client.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// VerifyField verifies a field
func (m *Moderation) VerifyField(workflowID, documentID, pageID, fieldDataID string, req VerifyFieldRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/workflows/%s/documents/%s/pages/%s/fields/%s/verify", m.client.BaseURL, workflowID, documentID, pageID, fieldDataID), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// VerifyTableCell verifies a table cell
func (m *Moderation) VerifyTableCell(workflowID, documentID, pageID, tableID, cellID string, req VerifyTableCellRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/workflows/%s/documents/%s/pages/%s/tables/%s/cells/%s/verify", m.client.BaseURL, workflowID, documentID, pageID, tableID, cellID), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// VerifyTable verifies a table
func (m *Moderation) VerifyTable(workflowID, documentID, pageID, tableID string, req VerifyTableRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/workflows/%s/documents/%s/pages/%s/tables/%s/verify", m.client.BaseURL, workflowID, documentID, pageID, tableID), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// VerifyDocument verifies a document
func (m *Moderation) VerifyDocument(workflowID, documentID string, req VerifyDocumentRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/workflows/%s/documents/%s/verify", m.client.BaseURL, workflowID, documentID), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// CreateWorkflowRequest represents a request to create a workflow
type CreateWorkflowRequest struct {
	Description  string `json:"description"`
	WorkflowType string `json:"workflow_type"`
}

// SetFieldsRequest represents a request to set fields and table headers
type SetFieldsRequest struct {
	Fields       []Field       `json:"fields"`
	TableHeaders []TableHeader `json:"table_headers"`
}

// UpdateFieldRequest represents a request to update a field
type UpdateFieldRequest struct {
	Name string `json:"name"`
}

// UpdateMetadataRequest represents a request to update metadata
type UpdateMetadataRequest struct {
	Description string `json:"description"`
}

// UpdateSettingsRequest represents a request to update settings
type UpdateSettingsRequest struct {
	TableCapture bool `json:"table_capture"`
}

// UploadDocumentRequest represents a request to upload a document
type UploadDocumentRequest struct {
	File     string            `json:"file"`
	Async    bool              `json:"async"`
	Metadata map[string]string `json:"metadata"`
}

// Workflow represents a workflow
type Workflow struct {
	WorkflowID string `json:"workflow_id"`
}

// Document represents a document
type Document struct {
	DocumentID string `json:"document_id"`
}

// Field represents a field
type Field struct {
	Name string `json:"name"`
}

// TableHeader represents a table header
type TableHeader struct {
	Name string `json:"name"`
}

// Table represents a table
type Table struct {
	TableID string `json:"table_id"`
}

// UploadResult represents the result of an upload
type UploadResult struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

// AddFieldRequest represents a request to add a field
type AddFieldRequest struct {
	FieldName           string    `json:"field_name"`
	Value               string    `json:"value"`
	Bbox                []float64 `json:"bbox"`
	Confidence          float64   `json:"confidence"`
	VerificationStatus  string    `json:"verification_status"`
	VerificationMessage string    `json:"verification_message"`
}

// AddTableRequest represents a request to add a table
type AddTableRequest struct {
	Bbox                []float64 `json:"bbox"`
	Headers             []string  `json:"headers"`
	VerificationStatus  string    `json:"verification_status"`
	VerificationMessage string    `json:"verification_message"`
	Cells               []Cell    `json:"cells"`
}

// AddTableCellRequest represents a request to add a table cell
type AddTableCellRequest struct {
	Row                 int       `json:"row"`
	Col                 int       `json:"col"`
	Header              string    `json:"header"`
	Text                string    `json:"text"`
	Bbox                []float64 `json:"bbox"`
	VerificationStatus  string    `json:"verification_status"`
	VerificationMessage string    `json:"verification_message"`
}

// UpdateTableCellRequest represents a request to update a table cell
type UpdateTableCellRequest struct {
	Value string `json:"value"`
}

// VerifyFieldRequest represents a request to verify a field
type VerifyFieldRequest struct {
	VerificationStatus  string `json:"verification_status"`
	VerificationMessage string `json:"verification_message"`
}

// VerifyTableCellRequest represents a request to verify a table cell
type VerifyTableCellRequest struct {
	VerificationStatus  string `json:"verification_status"`
	VerificationMessage string `json:"verification_message"`
}

// VerifyTableRequest represents a request to verify a table
type VerifyTableRequest struct {
	VerificationStatus  string `json:"verification_status"`
	VerificationMessage string `json:"verification_message"`
}

// VerifyDocumentRequest represents a request to verify a document
type VerifyDocumentRequest struct {
	VerificationStatus  string `json:"verification_status"`
	VerificationMessage string `json:"verification_message"`
}

// Cell represents a cell in a table
type Cell struct {
	Row                 int       `json:"row"`
	Col                 int       `json:"col"`
	Header              string    `json:"header"`
	Text                string    `json:"text"`
	Bbox                []float64 `json:"bbox"`
	VerificationStatus  string    `json:"verification_status"`
	VerificationMessage string    `json:"verification_message"`
}

// UploadFromURL uploads a document to a workflow from a URL
func (d *Documents) UploadFromURL(workflowID string, req UploadDocumentFromURLRequest) (*UploadResult, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := d.client.Client.Post(
		fmt.Sprintf("%s/workflows/%s/documents/url", d.client.BaseURL, workflowID),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result UploadResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListWithPagination retrieves documents for a workflow with pagination
func (d *Documents) ListWithPagination(workflowID string, page, limit int) ([]Document, error) {
	url := fmt.Sprintf("%s/workflows/%s/documents?page=%d&limit=%d", d.client.BaseURL, workflowID, page, limit)
	resp, err := d.client.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var documents []Document
	if err := json.NewDecoder(resp.Body).Decode(&documents); err != nil {
		return nil, err
	}
	return documents, nil
}

// GetOriginalFile downloads the original document file
func (d *Documents) GetOriginalFile(workflowID, documentID string) ([]byte, error) {
	url := fmt.Sprintf("%s/workflows/%s/documents/%s/original", d.client.BaseURL, workflowID, documentID)
	resp, err := d.client.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// WorkflowType represents a workflow type
// Add this type for GetTypes
// Example: {"id": "invoice", "name": "Invoice", "description": "Extracts data from invoices"}
type WorkflowType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UploadDocumentFromURLRequest represents a request to upload a document from a URL
type UploadDocumentFromURLRequest struct {
	URL      string            `json:"url"`
	Async    bool              `json:"async"`
	Metadata map[string]string `json:"metadata"`
}
