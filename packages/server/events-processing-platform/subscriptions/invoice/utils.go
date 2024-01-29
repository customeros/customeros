package invoice

import (
	"bytes"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"html/template"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func setEventSpanTagsAndLogFields(span opentracing.Span, evt eventstore.Event) {
	span.SetTag(tracing.SpanTagComponent, constants.ComponentSubscriptionInvoice)
	span.SetTag(tracing.SpanTagAggregateId, evt.GetAggregateID())
}

func FillInvoiceHtmlTemplate(tmpFile *os.File, invoiceData map[string]interface{}) error {

	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "os.Getwd")
	}

	// Build the full path to the template file
	templatePath := filepath.Join(currentDir, "/subscriptions/invoice/pdf_template/index.html")
	templateContent, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return errors.Wrap(err, "ioutil.ReadFile")
	}

	// Convert the template content to a string
	templateString := string(templateContent)

	// Load HTML template
	tmpl, err := template.New("template").Funcs(template.FuncMap{
		"safeHTML": func(text string) template.HTML {
			return template.HTML(text)
		},
	}).Parse(templateString)
	if err != nil {
		return errors.Wrap(err, "template.ParseFiles")
	}

	// Create a buffer to store the filled template
	var tplBuffer bytes.Buffer
	err = tmpl.Execute(&tplBuffer, invoiceData)
	if err != nil {
		return errors.Wrap(err, "tmpl.Execute")
	}

	// Write the filled template to the temporary HTML file
	_, err = tmpFile.Write(tplBuffer.Bytes())
	if err != nil {
		return errors.Wrap(err, "tmpHTMLFile.Write")
	}

	return nil
}

func ConvertInvoiceHtmlToPdf(pdfConverterUrl string, tmpFile *os.File, invoiceData map[string]interface{}) (*[]byte, error) {
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

	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "os.Getwd")
	}
	resourcesPath := filepath.Join(currentDir, "/subscriptions/invoice/pdf_template")

	// Prepare HTTP request
	url := pdfConverterUrl + "/forms/chromium/convert/html"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add files to the request
	// invoice html file
	invoiceHtmlFile, err := getFileByName(tmpFile.Name())
	if err != nil {
		return nil, errors.Wrap(err, "getFileByName")
	}
	err = addMultipartFile(writer, invoiceHtmlFile, "index.html")
	if err != nil {
		return nil, errors.Wrap(err, "addMultipartFile index.html")
	}

	//provider logo
	if providerLogoUrl, ok := invoiceData["ProviderLogoUrl"].(string); ok && providerLogoUrl != "" {
		err = addProviderLogoAsResourceFile(writer, providerLogoUrl)
		if err != nil {
			return nil, errors.Wrap(err, "addProviderLogoAsResourceFile")
		}
	}

	err = addResourceFile(writer, resourcesPath, "/index.css", "index.css")
	if err != nil {
		return nil, errors.Wrap(err, "addResourceFile index.css")
	}

	err = addResourceFile(writer, resourcesPath, "/style.css", "style.css")
	if err != nil {
		return nil, errors.Wrap(err, "addResourceFile style.css")
	}

	err = addResourceFile(writer, resourcesPath, "/fonts.css", "fonts.css")
	if err != nil {
		return nil, errors.Wrap(err, "addResourceFile fonts.css")
	}

	//images
	err = addResourceFile(writer, resourcesPath, "/customer-os.png", "customer-os.png")
	if err != nil {
		return nil, errors.Wrap(err, "addResourceFile customer-os.png")
	}
	err = addResourceFile(writer, resourcesPath, "/line11681-7w4.svg", "line11681-7w4.svg")
	if err != nil {
		return nil, errors.Wrap(err, "addResourceFile line11681-7w4.svg")
	}
	err = addResourceFile(writer, resourcesPath, "/line21681-3s8.svg", "line21681-3s8.svg")
	if err != nil {
		return nil, errors.Wrap(err, "addResourceFile line21681-3s8.svg")
	}
	err = addResourceFile(writer, resourcesPath, "/line31681-nvh.svg", "line31681-nvh.svg")
	if err != nil {
		return nil, errors.Wrap(err, "addResourceFile line31681-nvh.svg")
	}

	err = addMultipartValue(writer, "8.6", "paperWidth")
	if err != nil {
		return nil, errors.Wrap(err, "addMultipartValue paperWidth")
	}
	err = addMultipartValue(writer, "0", "marginTop")
	if err != nil {
		return nil, errors.Wrap(err, "addMultipartValue marginTop")
	}
	err = addMultipartValue(writer, "0", "marginBottom")
	if err != nil {
		return nil, errors.Wrap(err, "addMultipartValue marginBottom")
	}
	err = addMultipartValue(writer, "0", "marginLeft")
	if err != nil {
		return nil, errors.Wrap(err, "addMultipartValue marginLeft")
	}
	err = addMultipartValue(writer, "0", "marginRight")
	if err != nil {
		return nil, errors.Wrap(err, "addMultipartValue marginRight")
	}

	writer.Close()

	// Create HTTP request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequest")
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Perform HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client.Do")
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Error: Unexpected status code")
	}

	// Read the response body
	pdfBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadAll")
	}

	return &pdfBytes, nil
}

func addProviderLogoAsResourceFile(writer *multipart.Writer, logoURL string) error {
	fileExtension := GetFileExtensionFromUrl(logoURL)
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "downloaded-logo-*"+fileExtension)
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Download the image
	response, err := http.Get(logoURL)
	if err != nil {
		fmt.Println("Error downloading image:", err)
		return err
	}
	defer response.Body.Close()

	// Save the image to the temporary file
	_, err = io.Copy(tmpFile, response.Body)
	if err != nil {
		fmt.Println("Error saving image to file:", err)
		return err
	}

	logoFile, err := getFileByName(tmpFile.Name())

	err = addMultipartFile(writer, logoFile, "provider-logo"+fileExtension)
	if err != nil {
		return errors.Wrap(err, "addMultipartFile "+"provider-logo"+fileExtension)
	}

	return nil
}

func getFileByName(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "os.Open")
	}

	return file, nil
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
	file, err := getFileByName(filepath.Join(basePath, fileName))
	if err != nil {
		return errors.Wrap(err, "getFileByName")
	}
	err = addMultipartFile(writer, file, partName)
	if err != nil {
		return errors.Wrap(err, "addMultipartFile "+fileName)
	}

	return nil
}

func GetFileExtensionFromUrl(url string) string {
	ext := filepath.Ext(url)
	if ext == "" {
		contentType := http.DetectContentType([]byte(url))
		switch {
		case strings.HasPrefix(contentType, "image/jpeg"):
			ext = ".jpg"
		case strings.HasPrefix(contentType, "image/png"):
			ext = ".png"
		case strings.HasPrefix(contentType, "image/svg"):
			ext = ".svg"
		default:
			ext = ".bin"
		}
	}
	return ext
}
