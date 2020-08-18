package tester

import (
	"sync"
)

func Execute(directories, languages []string, async bool) []*ShowUnit {
	testUnits := MakeTestUnits(directories, languages)

	if len(testUnits) == 0 {
		return nil
	}

	outcomes := make([]*ShowUnit, len(testUnits))

	if async {

		wg := new(sync.WaitGroup)

		for i, testUnit := range testUnits {
			wg.Add(1)

			go func(i int, testUnit *TestUnit) {
				defer wg.Done()
				view := testUnit.InitView(true)
				outcome := testUnit.Exec(view)
				outcomes[i] = outcome
			}(i, testUnit)

		}

		wg.Wait()

	} else {

		for i, testUnit := range testUnits {
			view := testUnit.InitView(false)
			outcome := testUnit.Exec(view)
			outcomes[i] = outcome
		}

	}

	return outcomes
}