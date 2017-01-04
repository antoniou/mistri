package lambda

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Function struct {
	Path     string
	Name     string
	Target   string
	S3Bucket string
	S3Key    string
}

func (f *Function) Setup() {
	f.compile()
	f.install()
	f.cleanup()
}

func (f *Function) compile() {

	var out, stderr bytes.Buffer
	requirementsFile := fmt.Sprintf("%s/requirements.txt", f.Path)

	if _, err := os.Stat(requirementsFile); os.IsNotExist(err) {
		log.Printf("No requirements file found for %s, skipping", f.Path)
		return
	}
	cmd := exec.Command("pip",
		"install",
		"-r",
		requirementsFile,
		"-t",
		f.Name,
	)

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
	}

	log.Printf(out.String())
}

func (f *Function) install() {
	var zipManager Zipper = LambdaZipper{}
	err := zipManager.Zip(f.Path, f.Target)

	var uploader Uploader = S3Uploader{}
	err = uploader.Upload(f)

	if err != nil {
		log.Fatal(err)
	}
}

func (f *Function) cleanup() {
	log.Printf("[DEBUG] Cleaning up %s", f.Target)
	err := os.Remove(f.Target)
	if err != nil {
		log.Fatalf("Could not remove %s: %s", f.Path, err)
	}

}

func NewFunction(attrs map[string]string) *Function {
	path, ok := attrs["path"]
	if !ok {
		path = attrs["name"]
	}

	return &Function{
		Name:     attrs["name"],
		Path:     path,
		S3Bucket: attrs["s3bucket"],
		Target:   strings.Join([]string{attrs["name"], ".zip"}, ""),
		S3Key:    strings.Join([]string{attrs["s3KeyPrefix"], attrs["name"]}, "/"),
	}
}
