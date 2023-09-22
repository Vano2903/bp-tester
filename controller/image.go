package controller

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"github.com/vano2903/bp-tester/model"
)

// CheckIfImageCompiled will check the output of the image build to see if the image was compiled correctly
func (c *Controller) CheckIfImageCompiled(output []byte) bool {
	imageBuildOutput := string(output)
	lines := strings.Split(imageBuildOutput, "\n")
	return strings.Contains(lines[len(lines)-2], "Successfully tagged")
}

func (c *Controller) CreateBuildContext(ctx context.Context, source []byte) (io.ReadCloser, error) {
	buildDir := c.config.Test.BuildDir
	dockerfileDir := c.config.Test.DockerFilesDir

	// todo: for now it's only go, we can add more languages
	dockerfile, err := os.ReadFile(fmt.Sprintf("%s/%s.dockerfile", dockerfileDir, "go"))
	if err != nil {
		return nil, err
	}
	if err := writeToFile(buildDir+"/dockerfile", dockerfile); err != nil {
		return nil, err
	}
	if err := writeToFile(buildDir+"/main.go", source); err != nil {
		return nil, err
	}

	buildContext, err := archive.TarWithOptions(buildDir, &archive.TarOptions{
		NoLchown: true,
	})
	return buildContext, err
}

func (c *Controller) CreateImageName(code string) string {
	return fmt.Sprintf("bp-image-%s", code)
}

func (c *Controller) CreateImage(ctx context.Context, buildContext io.ReadCloser, imageName, code string) (*model.Image, error) {
	image := new(model.Image)
	image.Name = imageName

	//create the image from the dockerfile
	//we are setting some default labels and the flag -rm -f
	//!should set memory and cpu limit
	resp, err := c.cli.ImageBuild(ctx, buildContext, types.ImageBuildOptions{
		// Squash:     true,
		Dockerfile: "dockerfile",
		Tags:       []string{imageName},
		Labels: map[string]string{
			"builder":      c.config.APP.Name,
			"attempt-code": code,
		},
		Remove:      true,
		ForceRemove: true,
	})
	if err != nil {
		return nil, err
	}

	image.BuildOutput, err = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	image.ImageID, err = c.GetImageIDFromName(ctx, imageName)
	if err != nil {
		return nil, err
	}

	return image, err
}

func (c *Controller) GetImageIDFromName(ctx context.Context, name string) (string, error) {
	var out bytes.Buffer
	cmd := exec.CommandContext(ctx, "docker", "images", "-q", name)
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.ReplaceAll(out.String(), "\n", ""), err
}

func (c *Controller) RemoveImage(ctx context.Context, imageID string) error {
	_, err := c.cli.ImageRemove(ctx, imageID, types.ImageRemoveOptions{
		Force:         true,
		PruneChildren: true,
	})
	return err
}
