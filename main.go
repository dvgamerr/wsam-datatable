package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"regexp"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

func azErrors(err error) {
	if err != nil {
		if serr, ok := err.(azblob.StorageError); ok { // This error is a Service-specific
			switch serr.ServiceCode() { // Compare serviceCode to ServiceCodeXxx constants
			case azblob.ServiceCodeContainerAlreadyExists:
				fmt.Println("Received 409. Container already exists")
				return
			}
		}
		log.Fatal(err)
	}
}

func main() {
	accountName, accountKey := "cgsmpstnonprd", "N2YWp/A+BSIAWfVRcpgUb/wgcu+4MAxnou9kdd8lAOP5PJ88BYyhdTHN+hifKKRceG/dTJmvDNHKkWfz+awYVQ=="
	// accountName, accountKey := "cgpmpstauditprd", "6M3P58Ejq4buz/SJUz5v0xJnZAQ/B+A5hSkFP16vW1Ieuzf+a7iHawNJC3JhloBMApyhV4YkOIAic2AmwYr1vA=="
	containerName, filePrefix := "pmplogs", ""
	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// Create a random string for the quick start container

	// From the Azure portal, get your storage account blob service URL endpoint.
	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	containerURL := azblob.NewContainerURL(*URL, p)

	// // // Create the container
	// fmt.Printf("Creating a container named %s\n", containerName)
	ctx := context.Background() // This example uses a never-expiring context
	// _, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	// azErrors(err)

	// // // Create a file to test the upload and download.
	// fmt.Printf("Creating a dummy file to test the upload and download\n")
	// data := []byte("hello world this is a blob\n")
	// fileName := "xasfarhsdfhgsdfgsdfg.txt"
	// err = ioutil.WriteFile(fileName, data, 0700)
	// azErrors(err)

	// // Here's how to upload a blob.
	// blobURL := containerURL.NewBlockBlobURL(fileName)
	// file, err := os.Open(fileName)
	// azErrors(err)

	// // You can use the low-level PutBlob API to upload files. Low-level APIs are simple wrappers for the Azure Storage REST APIs.
	// // Note that PutBlob can upload up to 256MB data in one shot. Details: https://docs.microsoft.com/en-us/rest/api/storageservices/put-blob
	// // Following is commented out intentionally because we will instead use UploadFileToBlockBlob API to upload the blob
	// // _, err = blobURL.PutBlob(ctx, file, azblob.BlobHTTPHeaders{}, azblob.Metadata{}, azblob.BlobAccessConditions{})
	// // azErrors(err)

	// // The high-level API UploadFileToBlockBlob function uploads blocks in parallel for optimal performance, and can handle large files as well.
	// // This function calls PutBlock/PutBlockList for files larger 256 MBs, and calls PutBlob for any file smaller
	// fmt.Printf("Uploading the file with blob name: %s\n", fileName)
	// _, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{
	// 	BlockSize:   4 * 1024 * 1024,
	// 	Parallelism: 16})
	// azErrors(err)

	// marker := "2!128!MDAwMDUyIXVhdC8yMDIxLzA4LzE0LzEzL3BtcF9wZXJmb3JtYW5jZV9hcHBsaWNhdGlvbkxvZy50eHQhMDAwMDI4ITk5OTktMTItMzFUMjM6NTk6NTkuOTk5OTk5OVoh"
	// List the container that we have created above
	blobItems := []azblob.BlobItemInternal{}
	category := map[string]int{}
	fmt.Println(URL)

	var validName = regexp.MustCompile(`^(.+?)\/`)

	for marker := (azblob.Marker{}); marker.NotDone(); {
		listBlob, err := containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{
			Details: azblob.BlobListingDetails{
				UncommittedBlobs: true,
				Deleted:          true,
			},
			Prefix:     filePrefix,
			MaxResults: 0,
		})

		azErrors(err)
		fmt.Printf("\tFetch: %d\n", len(listBlob.Segment.BlobItems))
		for i := range listBlob.Segment.BlobItems {
			blob := &listBlob.Segment.BlobItems[i]

			match := validName.FindStringSubmatch(blob.Name)
			if len(match) < 2 {
				fmt.Println("\t- Name:", blob.Name)
				continue
			}

			category[match[1]]++
			// duplicate := false
			// for _, g := range category {
			// 	if g == match[1] {
			// 		duplicate = true
			// 		break
			// 	}
			// }
			// if !duplicate {
			// 	category = append(category, match[1])
			// }
		}

		blobItems = append(blobItems, listBlob.Segment.BlobItems...)
		marker = listBlob.NextMarker
	}

	fmt.Printf("Blob Total: %d\n\n", len(blobItems))
	for name, total := range category {
		fmt.Printf("Category: %s (%d)\n", name, total)
	}
	// blob := &blobItems[0]
	// fmt.Print("Name:", blob.Name)
	// // Here's how to download the blob
	// downloadResponse, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false)

	// // NOTE: automatically retries are performed if the connection fails
	// bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})

	// // read the body into a buffer
	// downloadedData := bytes.Buffer{}
	// _, err = downloadedData.ReadFrom(bodyStream)
	// azErrors(err)

	// // The downloaded blob data is in downloadData's buffer. :Let's print it
	// fmt.Printf("Downloaded the blob: " + downloadedData.String())

	// // Cleaning up the quick start by deleting the container and the file created locally
	// fmt.Printf("Press enter key to delete the sample files, example container, and exit the application.\n")
	// bufio.NewReader(os.Stdin).ReadBytes('\n')
	// fmt.Printf("Cleaning up.\n")
	// containerURL.Delete(ctx, azblob.ContainerAccessConditions{})
	// file.Close()
	// os.Remove(fileName)

	// println("GO:: Hello world!")
}

// //export add
// func Add(a, b int) int {
// 	return a + b
// }

// //export update
// func Update() {
// 	println("GO::EVENT Update!")

// 	document := js.Global().Get("document")

// 	aStr := document.Call("getElementById", "a").Get("value").String()
// 	bStr := document.Call("getElementById", "b").Get("value").String()
// 	a, _ := strconv.Atoi(aStr)
// 	b, _ := strconv.Atoi(bStr)
// 	result := Add(a, b)
// 	document.Call("getElementById", "result").Set("value", result)
// }
