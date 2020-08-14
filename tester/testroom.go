package tester

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cwd-k2/titania.go/client"
)

// Config
// test configs
type Config struct {
	Host                    string   `json:"host"`
	APIKey                  string   `json:"api_key"`
	SourceCodeDirectories   []string `json:"source_code_directories"`
	TestCaseDirectories     []string `json:"test_case_directories"`
	TestCaseInputExtension  string   `json:"test_case_input_extension"`
	TestCaseOutputExtension string   `json:"test_case_output_extension"`
	MaxProcesses            uint     `json:"max_processes"`
}

// TestRoom
// contains paiza.io API client, config, and map of TestUnits
type TestRoom struct {
	Client    *client.Client
	Config    *Config
	TestUnits map[string]*TestUnit
	TestCases map[string]*TestCase
}

// returns map[string]*TestRoom
func MakeTestRooms(directories []string) map[string]*TestRoom {
	testRooms := make(map[string]*TestRoom)
	// この前にディレクトリ直下に titania.json がいるか確認したい
	for _, dirname := range directories {
		// ディレクトリ直下の titania.json を読んで設定を作る

		baseDirectoryPath, err := filepath.Abs(dirname)
		// ここのエラーは公式のドキュメント見てもわからんのだけど何？
		if err != nil {
			fmt.Println(err)
			continue
		}

		configFileName := path.Join(baseDirectoryPath, "titania.json")
		configRawData, err := ioutil.ReadFile(configFileName)
		// File Read 失敗
		if err != nil {
			fmt.Printf("[SKIP] Couldn't read %s.\n", configFileName)
			continue
		}

		// ようやく設定の構造体を作れる
		config := new(Config)

		// JSON パース失敗
		if err := json.Unmarshal(configRawData, config); err != nil {
			fmt.Printf("[SKIP] Couldn't parse %s.\n%s\n", configFileName, err)
			continue
		}

		// paiza.io API クライアント
		client := new(client.Client)
		client.Host = config.Host
		client.APIKey = config.APIKey

		// テストケース
		testCases := MakeTestCases(
			baseDirectoryPath,
			config.TestCaseDirectories,
			config.TestCaseInputExtension,
			config.TestCaseOutputExtension)

		// テストユニット
		testUnits := MakeTestUnits(
			baseDirectoryPath,
			config.SourceCodeDirectories)

		testRooms[dirname] = new(TestRoom)
		testRooms[dirname].Client = client
		testRooms[dirname].Config = config
		testRooms[dirname].TestUnits = testUnits
		testRooms[dirname].TestCases = testCases

	}

	return testRooms
}

// 実際に paiza.io の API を利用して実行結果をもらう
func (testRoom *TestRoom) Exec() {

	for unitName, testUnit := range testRoom.TestUnits {
		fmt.Printf("  [UNIT] %s\n", unitName)
		fmt.Printf("    [LANG] %s\n", strings.ToUpper(testUnit.Language))
		wg := new(sync.WaitGroup)

		for caseName, testCase := range testRoom.TestCases {

			wg.Add(1)

			exec := func(caseName string, testCase *TestCase) {
				defer wg.Done()

				runnersCreateResponse :=
					testRoom.Client.RunnersCreate(
						testUnit.SourceCode,
						testUnit.Language,
						testCase.Input)

				runnersGetDetailsResponse :=
					testRoom.Client.RunnersGetDetails(runnersCreateResponse.ID)

				if runnersGetDetailsResponse.STDOUT == testCase.Output {
					fmt.Printf("    [CASE] %s [OK]\n", caseName)
				} else {
					fmt.Printf("    [CASE] %s [NG]\n", caseName)
				}

			}

			go exec(caseName, testCase)

		}
		wg.Wait()
	}
}
