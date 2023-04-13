package test

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/PureMLHQ/PureML/packages/purebackend/core"
)

type TestApp struct {
	*core.BaseApp
}

// Cleanup resets the test application state and removes the test
// app's dataDir from the filesystem.
//
// After this call, the app instance shouldn't be used anymore.
func (t *TestApp) Cleanup() {
	err := t.ResetBootstrapState()
	if err != nil {
		panic(err)
	}

	if t.DataDir() != "" {
		os.RemoveAll(t.DataDir())
	}
}

func NewTestApp(optTestDataDir ...string) (*TestApp, error) {
	var testDataDir string
	if len(optTestDataDir) == 0 || optTestDataDir[0] == "" {
		// fallback to the default test data directory
		_, currentFile, _, _ := runtime.Caller(0)
		testDataDir = filepath.Join(path.Dir(currentFile), "data")
	} else {
		testDataDir = optTestDataDir[0]
	}

	tempDir, err := TempDirClone(testDataDir)
	if err != nil {
		return nil, err
	}

	app := core.NewBaseApp(&core.BaseAppConfig{
		DataDir:      tempDir,
		IsDebug:      false,
		DatabaseType: "sqlite3",
	})

	// load data dir and db connections
	if err := app.Bootstrap(); err != nil {
		return nil, err
	}
	app.Settings().AdminAuthToken.Secret = "pureml-test-secret"
	app.Settings().APITokenSecretKey = "pureml-test-secret"

	t := &TestApp{
		BaseApp: app,
	}
	return t, nil
}

// TempDirClone creates a new temporary directory copy from the
// provided directory path.
//
// It is the caller's responsibility to call `os.RemoveAll(tempDir)`
// when the directory is no longer needed!
func TempDirClone(dirToClone string) (string, error) {
	tempDir, err := os.MkdirTemp("", "purebackend_test_*")
	if err != nil {
		return "", err
	}

	// copy everything from testDataDir to tempDir
	if err := copyDir(dirToClone, tempDir); err != nil {
		return "", err
	}

	return tempDir, nil
}

// -------------------------------------------------------------------
// Helpers
// -------------------------------------------------------------------

func copyDir(src string, dest string) error {
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return err
	}

	sourceDir, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceDir.Close()

	items, err := sourceDir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, item := range items {
		fullSrcPath := filepath.Join(src, item.Name())
		fullDestPath := filepath.Join(dest, item.Name())

		var copyErr error
		if item.IsDir() {
			copyErr = copyDir(fullSrcPath, fullDestPath)
		} else {
			copyErr = copyFile(fullSrcPath, fullDestPath)
		}

		if copyErr != nil {
			return copyErr
		}
	}

	return nil
}

func copyFile(src string, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}

	return nil
}
