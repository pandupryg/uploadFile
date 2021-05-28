package cloudbucket

import (
	"io"
	"net/http"
	"net/url"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
)

var (
	storageClient *storage.Client
)

// HandleFileUploadToBucket uploads file to bucket
func HandleFileUploadToBucket(c *gin.Context) {
	bucket := "test-my-bucket-project" //your bucket name

	var checkerr error

	ctx := appengine.NewContext(c.Request)

	storageClient, checkerr = storage.NewClient(ctx, option.WithCredentialsFile("keys.json"))
	if checkerr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": checkerr.Error(),
			"error":   true,
		})
		return
	}

	f, uploadedFile, checkerr := c.Request.FormFile("file")
	if checkerr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": checkerr.Error(),
			"error":   true,
		})
		return
	}

	defer f.Close()

	sw := storageClient.Bucket(bucket).Object(uploadedFile.Filename).NewWriter(ctx)

	if _, checkerr := io.Copy(sw, f); checkerr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": checkerr.Error(),
			"error":   true,
		})
		return
	}

	if checkerr := sw.Close(); checkerr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": checkerr.Error(),
			"error":   true,
		})
		return
	}

	u, checkerr := url.Parse("/" + bucket + "/" + sw.Attrs().Name)
	if checkerr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": checkerr.Error(),
			"Error":   true,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "file uploaded successfully",
		"pathname": u.EscapedPath(),
	})
}
