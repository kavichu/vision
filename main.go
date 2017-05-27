// Sample vision-quickstart uses the Google Cloud Vision API to label an image.
package main

import (
    "fmt"
    "log"
    "os"

    // Imports the Google Cloud Vision API client package.
    vision "cloud.google.com/go/vision/apiv1"
    "golang.org/x/net/context"
    "encoding/json"
    "net/http"
)

type VisionLabels map[string]float32

type PwnedImage struct {
    // original_file
    Original_file string
    Labels VisionLabels
}

func handler (w http.ResponseWriter, req *http.Request) {
    decoder := json.NewDecoder(req.Body)
    print(req.Body)
    var img PwnedImage
    err := decoder.Decode(&img)
    if err != nil {
        panic(err)
    }
    defer req.Body.Close()
    log.Println(img.Original_file)

    ctx := context.Background()

    // Creates a client.
    client, err := vision.NewImageAnnotatorClient(ctx)
    if err != nil {
            log.Fatalf("Failed to create client: %v", err)
    }

    // Sets the name of the image file to annotate.
    filename := fmt.Sprintf("%s/%s", os.Getenv("DOWNLOADS_LOCATION"), img.Original_file)
    file, err := os.Open(filename)
    if err != nil {
            log.Fatalf("Failed to read file: %v", err)
    }
    defer file.Close()
    image, err := vision.NewImageFromReader(file)
    if err != nil {
            log.Fatalf("Failed to create image: %v", err)
    }

    labels, err := client.DetectLabels(ctx, image, nil, 10)
    if err != nil {
            log.Fatalf("Failed to detect labels: %v", err)
    }

    fmt.Println("Labels:")
    labelsMap := make(map[string]float32)
    for _, label := range labels {
        fmt.Println(label.Description)
        fmt.Println(label.Score)
        labelsMap[label.Description] = label.Score
    }

    img.Labels = labelsMap
    jsonData, err := json.Marshal(img)
    if err != nil {
        log.Fatalf("Failed to parse json: %v", err)
    }

    w.Header().Set("Content-Type","application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonData)
}

func main() {
    http.HandleFunc("/detect", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}