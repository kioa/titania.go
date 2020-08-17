package tester

import (
	"sort"
	"sync"
)

func Execute(directories, languages []string, async bool) []*ShowUnit {

	outcomes := []*ShowUnit{}

	wg := new(sync.WaitGroup)

	for _, dirname := range directories {

		if async {

			wg.Add(1)
			go execute(dirname, languages, &outcomes, wg)

		} else {

			execute(dirname, languages, &outcomes, nil)

		}
	}

	wg.Wait()

	sort.Slice(outcomes, func(i, j int) bool {
		return outcomes[i].Name < outcomes[j].Name
	})

	return outcomes
}

func execute(
	dirname string, languages []string,
	outcomes *[]*ShowUnit, wg *sync.WaitGroup) {

	if wg != nil {
		defer wg.Done()
	}

	testUnit := NewTestUnit(dirname, languages)

	// 実行するテストがない
	if testUnit == nil {
		return
	}

	fruits := testUnit.Exec(wg != nil)

	outcome := new(ShowUnit)
	outcome.Name = dirname
	outcome.Fruits = fruits

	*outcomes = append(*outcomes, outcome)
}
