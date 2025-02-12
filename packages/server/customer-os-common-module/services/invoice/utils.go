package invoice

import (
	"bytes"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func FillInvoiceHtmlTemplate(ctx context.Context, tmpFile *os.File, invoiceData map[string]interface{}) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FillInvoiceHtmlTemplate")
	defer span.Finish()

	// Read template from the embedded FS.
	templateContent, err := Templates.ReadFile("pdf_template/index.html")
	if err != nil {
		return errors.Wrap(err, "reading embedded template")
	}

	// Convert the template content to a string.
	templateString := string(templateContent)

	// Load HTML template.
	tmpl, err := template.New("template").Funcs(template.FuncMap{
		"safeHTML": func(text string) template.HTML {
			return template.HTML(text)
		},
	}).Parse(templateString)
	if err != nil {
		return errors.Wrap(err, "parsing template")
	}

	// Create a buffer to store the filled template.
	var tplBuffer bytes.Buffer
	err = tmpl.Execute(&tplBuffer, invoiceData)
	if err != nil {
		return errors.Wrap(err, "executing template")
	}

	// Write the filled template to the temporary HTML file.
	_, err = tmpFile.Write(tplBuffer.Bytes())
	if err != nil {
		return errors.Wrap(err, "writing to temporary file")
	}

	return nil
}

func ConvertInvoiceHtmlToPdf(ctx context.Context, fsc interfaces.FileService, pdfConverterUrl string, tmpFile *os.File, invoiceData map[string]interface{}) (*[]byte, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ConvertInvoiceHtmlToPdf")
	defer span.Finish()
	// This is doing a request like this:
	//curl \
	//--request POST 'http://localhost:11006/forms/chromium/convert/html' \
	//--form 'files=@"index.html"' \
	//--form 'files=@"style.css"' \
	//--form 'files=@"index.css"' \
	//--form 'files=@"fonts.css"' \
	//--form 'files=@"customer-os.png"' \
	//--form 'files=@"provider_logo.png"' \
	//-o my.pdf

	// Construct the URL for PDF conversion.
	url := pdfConverterUrl + "/forms/chromium/convert/html"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 1. Add the invoice HTML file (already created and filled in tmpFile).
	invoiceHtmlFile, err := utils.GetFileByName(tmpFile.Name())
	if err != nil {
		tracing.TraceErr(span, fmt.Errorf("getFileByName: %w", err))
		return nil, fmt.Errorf("getFileByName: %w", err)
	}
	err = addMultipartFile(writer, invoiceHtmlFile, "index.html")
	if err != nil {
		tracing.TraceErr(span, fmt.Errorf("addMultipartFile index.html: %w", err))
		return nil, fmt.Errorf("addMultipartFile index.html: %w", err)
	}

	// 2. Add provider logo if available
	if providerLogoRepositoryFileId, ok := invoiceData["ProviderLogoRepositoryFileId"].(string); ok && providerLogoRepositoryFileId != "" {
		file, metadata, err := downloadProviderLogoAsTempFile(ctx, fsc, invoiceData["Tenant"].(string), providerLogoRepositoryFileId, span)
		if err != nil {
			tracing.TraceErr(span, fmt.Errorf("downloadProviderLogoAsTempFile: %w", err))
			return nil, fmt.Errorf("downloadProviderLogoAsTempFile: %w", err)
		}

		fileExtension := GetFileExtensionFromMetadata(metadata)
		err = addMultipartFile(writer, file, "provider-logo"+fileExtension)
		if err != nil {
			tracing.TraceErr(span, fmt.Errorf("addMultipartFile provider-logo%s: %w", fileExtension, err))
			return nil, fmt.Errorf("addMultipartFile provider-logo%s: %w", fileExtension, err)
		}
	}

	// 3. Add static resource files from the embedded FS.
	resourceFiles := []struct {
		FileName string
		PartName string
	}{
		{"index.css", "index.css"},
		{"style.css", "style.css"},
		{"fonts.css", "fonts.css"},
		{"customer-os.png", "customer-os.png"},
		{"preview-stamp.png", "preview-stamp.png"},
		{"line11681-7w4.svg", "line11681-7w4.svg"},
		{"line21681-3s8.svg", "line21681-3s8.svg"},
		{"line31681-nvh.svg", "line31681-nvh.svg"},
	}
	for _, rf := range resourceFiles {
		err = addEmbeddedResourceFile(writer, rf.FileName, rf.PartName)
		if err != nil {
			tracing.TraceErr(span, fmt.Errorf("addEmbeddedResourceFile %s: %w", rf.FileName, err))
			return nil, fmt.Errorf("addEmbeddedResourceFile %s: %w", rf.FileName, err)
		}
	}

	// 4. Add multipart form fields (e.g., paper width, margins, etc.)
	formFields := []struct {
		FieldName string
		Value     string
	}{
		{"paperWidth", "8.6"},
		{"marginTop", "0"},
		{"marginBottom", "0"},
		{"marginLeft", "0"},
		{"marginRight", "0"},
	}
	for _, ff := range formFields {
		err = addMultipartValue(writer, ff.Value, ff.FieldName)
		if err != nil {
			tracing.TraceErr(span, fmt.Errorf("addMultipartValue %s: %w", ff.FieldName, err))
			return nil, fmt.Errorf("addMultipartValue %s: %w", ff.FieldName, err)
		}
	}

	// Close the multipart writer to flush the buffer.
	err = writer.Close()
	if err != nil {
		tracing.TraceErr(span, fmt.Errorf("writer.Close: %w", err))
		return nil, fmt.Errorf("writer.Close: %w", err)
	}

	// Create the HTTP request.
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		tracing.TraceErr(span, fmt.Errorf("http.NewRequest: %w", err))
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Perform the HTTP request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, fmt.Errorf("client.Do: %w", err))
		return nil, fmt.Errorf("client.Do: %w", err)
	}
	defer resp.Body.Close()

	// Check for a successful status code.
	if resp.StatusCode != http.StatusOK {
		span.LogFields(log.String("status_code", resp.Status))
		tracing.TraceErr(span, fmt.Errorf("unexpected status code %v", resp.StatusCode))
		return nil, fmt.Errorf("unexpected status code %v", resp.StatusCode)
	}

	// Read the response body (the PDF bytes)
	pdfBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, fmt.Errorf("io.ReadAll: %w", err))
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	return &pdfBytes, nil
}

// addEmbeddedResourceFile reads the given file from the embedded FS (Templates)
// and adds it as a file part to the multipart writer.
func addEmbeddedResourceFile(writer *multipart.Writer, fileName, partName string) error {
	data, err := Templates.ReadFile("pdf_template/" + fileName)
	if err != nil {
		return errors.Wrapf(err, "failed to read embedded file %s", fileName)
	}
	part, err := writer.CreateFormFile("files", partName)
	if err != nil {
		return errors.Wrapf(err, "failed to create form file for %s", fileName)
	}
	_, err = part.Write(data)
	if err != nil {
		return errors.Wrapf(err, "failed to write data for %s", fileName)
	}
	return nil
}

func downloadProviderLogoAsTempFile(ctx context.Context, fileService interfaces.FileService, tenant, repositoryFileId string, span opentracing.Span) (*os.File, *interfaces.File, error) {
	fileMetadata, err := fileService.GetById(ctx, repositoryFileId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "fileService.GetById"))
		return nil, nil, err
	}
	if fileMetadata == nil {
		tracing.TraceErr(span, errors.Errorf("File with id %v not found", repositoryFileId))
		return nil, nil, errors.Errorf("File with id %v not found", repositoryFileId)
	}
	fileBytes, err := fileService.GetFileBytes(ctx, repositoryFileId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "fileService.GetFileBytes"))
		return nil, nil, err
	}

	fileExtension := GetFileExtensionFromMetadata(fileMetadata)
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "downloaded-logo-*"+fileExtension)
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return nil, nil, err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Save the image to the temporary file
	_, err = io.Copy(tmpFile, bytes.NewReader(*fileBytes))
	if err != nil {
		fmt.Println("Error copying file to temporary file:", err)
		return nil, nil, err
	}

	fileByName, err := utils.GetFileByName(tmpFile.Name())
	if err != nil {
		fmt.Println("Error getting file by name:", err)
		return nil, nil, err
	}
	return fileByName, fileMetadata, nil
}

func addMultipartValue(writer *multipart.Writer, value string, partName string) error {
	part, err := writer.CreateFormField(partName)
	if err != nil {
		return errors.Wrap(err, "writer.CreateFormFile")
	}
	_, err = part.Write([]byte(value))
	if err != nil {
		return errors.Wrap(err, "part.Write")
	}
	return nil
}

func addMultipartFile(writer *multipart.Writer, file *os.File, partName string) error {
	part, err := writer.CreateFormFile("files", partName)
	if err != nil {
		return errors.Wrap(err, "writer.CreateFormFile")
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return errors.Wrap(err, "io.Copy")
	}
	return nil
}

func addResourceFile(writer *multipart.Writer, basePath, fileName, partName string) error {
	file, err := utils.GetFileByName(filepath.Join(basePath, fileName))
	if err != nil {
		return errors.Wrap(err, "getFileByName")
	}
	err = addMultipartFile(writer, file, partName)
	if err != nil {
		return errors.Wrap(err, "addMultipartFile "+fileName)
	}

	return nil
}

func GetFileExtensionFromMetadata(metadata *interfaces.File) string {
	if metadata == nil {
		return ""
	}
	return strings.Split(metadata.MimeType, "/")[1]
}
