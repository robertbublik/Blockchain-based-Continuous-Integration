package mine

import (
	"fmt"
 	"github.com/robertbublik/bci/database"
	//"github.com/robertbublik/bci/node"
	"log"
	"io/ioutil"
	"io"
	"context"
	"bytes"
	"os"
	"archive/tar"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

)

func DockerBuild(tx database.Tx, dockerfilePath string, imageName string) {
	fmt.Printf("dockerfilepath: %s\n imagename: %s", dockerfilePath, imageName)
	
	client, err := client.NewEnvClient()
	if err != nil {
		log.Fatalf("Unable to create docker client: %s", err)
	}
	
	// Client, imagename and Dockerfile location
	tags := []string{imageName}
	err = buildImage(client, tags, dockerfilePath)
	if err != nil {
		log.Println(err)
	}
}

func buildImage(client *client.Client, tags []string, dockerfile string)  error {
	ctx := context.Background()

	// Create a buffer 
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	// Create a filereader
	dockerFileReader, err := os.Open(dockerfile)
	if err != nil {
		return err
	}

	// Read the actual Dockerfile 
	readDockerFile, err := ioutil.ReadAll(dockerFileReader)
	if err != nil {
		return err
	}

	// Make a TAR header for the file
	tarHeader := &tar.Header{
		Name: dockerfile,
		Size: int64(len(readDockerFile)),
	}

	// Writes the header described for the TAR file
	err = tw.WriteHeader(tarHeader)
    if err != nil {
		return err
    }

	// Writes the dockerfile data to the TAR file
    _, err = tw.Write(readDockerFile)
    if err != nil {
		return err
    }

    dockerFileTarReader := bytes.NewReader(buf.Bytes())

	// Define the build options to use for the file
	// https://godoc.org/github.com/docker/docker/api/types#ImageBuildOptions
	buildOptions := types.ImageBuildOptions{
        Context:    dockerFileTarReader,
        Dockerfile: dockerfile,
        Remove:     true,
		Tags: 		tags,
	}

	// Build the actual image
	imageBuildResponse, err := client.ImageBuild(
        ctx,
        dockerFileTarReader,
		buildOptions, 
	)	

	if err != nil {
		return err
	}

	// Read the STDOUT from the build process
	defer imageBuildResponse.Body.Close()
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		return err
	}

	return nil
}