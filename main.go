package main

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var protoFilesUsingProtocGenGoFast = map[string]bool{"": true}

const protoc = "protoc"

func main() {
	pwd, wdErr := os.Getwd()
	if wdErr != nil {
		os.Exit(1)
	}

	GOBIN := GetGoBin()
	protoFilesMap := make(map[string][]string)
	walkErr := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		dir := filepath.Dir(path)
		filename := filepath.Base(path)
		if strings.HasSuffix(filename, ".proto") {
			protoFilesMap[dir] = append(protoFilesMap[dir], path)
		}

		return nil
	})
	if walkErr != nil {
		fmt.Println(walkErr)
		os.Exit(1)
	}
	for _, files := range protoFilesMap {
		for _, relProtoFile := range files {
			var args []string
			if protoFilesUsingProtocGenGoFast[relProtoFile] {
				args = []string{"--gofast_out", pwd, "--plugin", "protoc-gen-gofast=" + GOBIN + "/protoc-gen-gofast"}
			} else {
				args = []string{"--go_out", pwd, "--go-grpc_out", pwd, "--plugin", "protoc-gen-go=" + filepath.Join(GOBIN, ToolsName("protoc-gen-go")), "--plugin", "protoc-gen-go-grpc=" + filepath.Join(GOBIN, ToolsName("protoc-gen-go-grpc"))}
			}
			args = append(args, relProtoFile)
			cmd := exec.Command(protoc, args...)
			cmd.Env = append(cmd.Env, os.Environ()...)
			cmd.Env = append(cmd.Env, "GOBIN="+GOBIN)
			output, cmdErr := cmd.CombinedOutput()
			if len(output) > 0 {
				fmt.Println("cmd:", string(output))
			}
			if cmdErr != nil {
				fmt.Println(cmdErr)
				os.Exit(1)
			}
		}
	}

	moduleName, gmnErr := GetModuleName(pwd)
	if gmnErr != nil {
		fmt.Println(gmnErr)
		os.Exit(1)
	}
	modulePath := filepath.Join(strings.Split(moduleName, "/")...)

	pbGoFilesMap := make(map[string][]string)
	walkErr2 := filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		dir := filepath.Dir(path)
		filename := filepath.Base(path)
		if strings.HasSuffix(filename, ".pb.go") {
			pbGoFilesMap[dir] = append(pbGoFilesMap[dir], path)
		}

		return nil
	})
	if walkErr2 != nil {
		fmt.Println(walkErr2)
		os.Exit(1)
	}

	var err error
	for _, srcPbGoFiles := range pbGoFilesMap {
		for _, srcPbGoFile := range srcPbGoFiles {
			var dstPbGoFile string
			dstPbGoFile, err = filepath.Rel(modulePath, srcPbGoFile)
			if err != nil {
				continue
			}
			err = os.Link(srcPbGoFile, dstPbGoFile)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Printf("'%s' does not exist\n", srcPbGoFile)
					continue
				}
				if os.IsPermission(err) {
					continue
				}
				if os.IsExist(err) {
					err = os.Remove(dstPbGoFile)
					if err != nil {
						fmt.Printf("Failed to delete file '%s'\n", dstPbGoFile)
						continue
					}
					err = os.Rename(srcPbGoFile, dstPbGoFile)
					if err != nil {
						fmt.Printf("Can not move '%s' to '%s'\n", srcPbGoFile, dstPbGoFile)
					}
					continue
				}
			}
			err = os.Rename(srcPbGoFile, dstPbGoFile)
			if err != nil {
				fmt.Printf("Can not move '%s' to '%s'\n", srcPbGoFile, dstPbGoFile)
			}
			continue
		}
	}

	if err == nil {
		err = os.RemoveAll(strings.Split(modulePath, "/")[0])
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("生成proto.go成功")
}

func envFile() (string, error) {
	if file := os.Getenv("GOENV"); file != "" {
		if file == "off" {
			return "", fmt.Errorf("GOENV=off")
		}
		return file, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	if dir == "" {
		return "", fmt.Errorf("missing user-config dir")
	}
	return filepath.Join(dir, "go", "env"), nil
}

func GetRuntimeEnv(key string) (string, error) {
	file, err := envFile()
	if err != nil {
		return "", err
	}
	if file == "" {
		return "", fmt.Errorf("missing runtime env file")
	}
	var data []byte
	var runtimeEnv string
	data, readErr := ioutil.ReadFile(file)
	if readErr != nil {
		return "", readErr
	}
	envStrings := strings.Split(string(data), "\n")
	for _, envItem := range envStrings {
		envItem = strings.TrimSuffix(envItem, "\r")
		envKeyValue := strings.Split(envItem, "=")
		if strings.EqualFold(strings.TrimSpace(envKeyValue[0]), key) {
			runtimeEnv = strings.TrimSpace(envKeyValue[1])
		}
	}
	return runtimeEnv, nil
}

func GetGoBin() string {
	GOBIN := os.Getenv("GOBIN")
	if GOBIN == "" {
		var err error
		// The one set by user by running `go env -w GOBIN=/path`
		GOBIN, err = GetRuntimeEnv("GOBIN")
		if err != nil {
			// The default one that Golang uses
			return filepath.Join(build.Default.GOPATH, "bin")
		}
		if GOBIN == "" {
			return filepath.Join(build.Default.GOPATH, "bin")
		}
		return GOBIN
	}
	return GOBIN
}

func GetModuleName(path string) (string, error) {
	gomodPath := filepath.Join(path, "go.mod")
	gomodBytes, err := ioutil.ReadFile(gomodPath)
	if err != nil {
		return "", err
	}
	gomodContent := string(gomodBytes)
	moduleIdx := strings.Index(gomodContent, "module") + 6
	newLineIdx := strings.Index(gomodContent, "\n")

	var moduleName string
	if moduleIdx >= 0 {
		if newLineIdx >= 0 {
			moduleName = strings.TrimSpace(gomodContent[moduleIdx:newLineIdx])
			moduleName = strings.TrimSuffix(moduleName, "\r")
		} else {
			moduleName = strings.TrimSpace(gomodContent[moduleIdx:])
		}
	} else {
		return "", fmt.Errorf("can not get the value of `module` in path `%s`", gomodPath)
	}
	return moduleName, nil
}

func ToolsName(name string) string {
	if runtime.GOOS == "windows" {
		return strings.Join([]string{name, "exe"}, ".")
	}
	return name
}
