package tail

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccess(t *testing.T) {

	file := setup()
	defer cleanup(file)

	// create stream from tmp file
	ch, err := Stream(context.Background(), file.Name())
	if err != nil {
		log.Fatal(err)
	}

	// push data to stream
	if _, err := file.WriteString("Hello"); err != nil {
		log.Fatal(err)
	}

	// recieve data from stream
	actual := <-ch
	assert.Equal(t, "Hello", actual)

	if _, err := file.WriteString("World"); err != nil {
		log.Fatal(err)
	}
	actual = <-ch
	assert.Equal(t, "World", actual)
}

func TestFails_Process_Does_Not_Exist(t *testing.T) {

	// create stream from tmp file
	ch, err := Stream(context.Background(), "")
	assert.Nil(t, ch)
	assert.Error(t, err)
}

func setup() *os.File {
	tempDirPath, err := ioutil.TempDir("", "myTempDir")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Temp dir created:", tempDirPath)

	// Create a file in new temp directory
	tempFile, err := ioutil.TempFile(tempDirPath, "myTempFile.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Temp file created:", tempFile.Name())

	return tempFile
}

func cleanup(tempFile *os.File) error {
	tempFile.Close()
	// Delete the resources we created
	if err := os.Remove(tempFile.Name()); err != nil {
		log.Fatal(err)
	}
	return nil
}
