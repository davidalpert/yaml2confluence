package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/thanhpk/randstr"
)

const DEFAULT_SECRET_FILENAME = ".secret"

func GetDefaultY2cHomeDir() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(dirname, ".y2c")
}

func GetDefaultSecretPath() string {
	return filepath.Join(GetDefaultY2cHomeDir(), DEFAULT_SECRET_FILENAME)
}

func GetSecret() (string, error) {
	secret, err := os.ReadFile(GetDefaultSecretPath())
	if err != nil {
		return "", err
	}

	return string(secret), nil
}

func GetSecretAndGenerateIfMissing() string {
	secret, err := GetSecret()
	if err == nil {
		return secret
	}

	return GenerateSecret()
}

func GenerateSecret() string {
	secret := []byte(randstr.String(24))

	dir := GetDefaultY2cHomeDir()
	secretPath := GetDefaultSecretPath()
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(secretPath, secret, 0600)
	if err != nil {
		panic(err)
	}

	fmt.Println("Generated secret key: " + secretPath)

	return string(secret)
}

func CreateInstanceDirectory(baseDir string, name string, config string) {
	var dir string
	current, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if baseDir != "" {
		if filepath.IsAbs(baseDir) {
			dir = baseDir
		} else {
			absPath, err := filepath.Abs(filepath.Join(current, baseDir))
			if err != nil {
				panic(err)
			}
			dir = absPath
		}
	} else {
		dir = current
	}

	instanceDir := filepath.Join(dir, name)
	spacesDir := filepath.Join(instanceDir, "spaces")
	templatesDir := filepath.Join(instanceDir, "templates")
	configFile := filepath.Join(instanceDir, "config.yml")

	// create instance directory
	err = os.Mkdir(instanceDir, 0755)
	if err != nil {
		if os.IsExist(err) {
			fmt.Println(instanceDir + " already exists.")
			os.Exit(1)
		} else {
			panic(err)
		}
	}
	fmt.Println("Created directory " + instanceDir)
	// create spaces directory
	err = os.Mkdir(spacesDir, 0755)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created directory " + spacesDir)
	// create templates directory
	err = os.Mkdir(templatesDir, 0755)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created directory " + templatesDir)

	// write config.yml
	err = os.WriteFile(configFile, []byte(config), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created file " + configFile)

}

func ResolveAbsolutePathDir(dir string) string {
	if filepath.IsAbs(dir) {
		return dir
	}

	current, err := fs.Getwd()
	if err != nil {
		panic(err)
	}

	if dir == "" {
		return current
	}

	absPath, err := filepath.Abs(filepath.Join(current, dir))
	if err != nil {
		panic(err)
	}

	return absPath
}

func ResolveAbsolutePathFile(file string) string {
	if file == "" {
		panic("no file provided")
	}

	return ResolveAbsolutePathDir(file)
}
