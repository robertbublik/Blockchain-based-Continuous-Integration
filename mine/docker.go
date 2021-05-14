package mine

import (
	"fmt"
	"bufio"
	"errors"
 	"github.com/robertbublik/bci/database"
	"encoding/json"
	"encoding/base64"
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

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

var authConfig = types.AuthConfig{
	Username:      "admin",
	Password:      "123",
	ServerAddress: registryUrl,
	}

func DockerBuildAndPush(tx database.Tx, dockerfilePath string, imageTag string) {
	client, err := client.NewEnvClient()
	if err != nil {
		log.Fatalf("Unable to create docker client: %s", err)
	}
	
	tags := []string{imageTag}
	err = buildImage(client, tags, dockerfilePath)
	if err != nil {
		log.Println(err)
	}

	err = pushImage(client, imageTag)
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

	dockerfileReader, err := os.Open(dockerfile)
	if err != nil {
		return err
	}

	readDockerfile, err := ioutil.ReadAll(dockerfileReader)
	if err != nil {
		return err
	}

	tarHeader := &tar.Header{
		Name: dockerfile,
		Size: int64(len(readDockerfile)),
	}

	err = tw.WriteHeader(tarHeader)
    if err != nil {
		return err
    }

    _, err = tw.Write(readDockerfile)
    if err != nil {
		return err
    }

    dockerFileTarReader := bytes.NewReader(buf.Bytes())

	buildOptions := types.ImageBuildOptions{
        Context:    dockerFileTarReader,
        Dockerfile: dockerfile,
        Remove:     true,
		Tags: 		tags,
	}
	imageBuildResponse, err := client.ImageBuild(
        ctx,
        dockerFileTarReader,
		buildOptions, 
	)	

	if err != nil {
		return err
	}

	defer imageBuildResponse.Body.Close()
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		return err
	}

	return nil
}


func pushImage(client *client.Client, imageTag string) error {
	ctx := context.Background()

	authConfigBytes, _ := json.Marshal(authConfig)
	authConfigEncoded := base64.URLEncoding.EncodeToString(authConfigBytes)

	//tag := dockerRegistryUserID + "/node-hello"
	opts := types.ImagePushOptions{RegistryAuth: authConfigEncoded}
	fmt.Printf("Pushing image %s", imageTag)
	imagePushResponse, err := client.ImagePush(ctx, imageTag, opts)
	if err != nil {
		return err
	}

	defer imagePushResponse.Close()
	_, err = io.Copy(os.Stdout, imagePushResponse)
	if err != nil {
		return err
	}

	err = print(imagePushResponse)
	if err != nil {
		return err
	}

	return nil
}

func print(rd io.Reader) error {
	var lastLine string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
		fmt.Println(scanner.Text())
	}

	errLine := &ErrorLine{}
	json.Unmarshal([]byte(lastLine), errLine)
	if errLine.Error != "" {
		return errors.New(errLine.Error)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}