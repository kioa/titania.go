package main

import (
	"os"

	"github.com/cwd-k2/titania.go/tester"
)

const VERSION = "v0.1.0"

func main() {
	// ターゲットのディレクトリと言語，async
	directories, languages, async := OptParse()
	outcomes := tester.Execute(directories, languages, async)

	// 何もテストが実行されなかった場合
	if outcomes == nil {
		println("Uh, OK, there's no test.")
		os.Exit(1)
	}

	tester.Final(outcomes)
	tester.Print(outcomes)
}
