package bdd_test

import (
	"testing"

	"github.com/octohelm/x/testing/bdd"
	"github.com/octohelm/x/testing/snapshot"
)

func TestFeature(t *testing.T) {
	t.Run("case 1", bdd.ScenarioT(func(b bdd.T) {
		b.Given("initial value with 1", func(t bdd.T) {
			value := 1

			t.When("add 1", func(c bdd.T) {
				value += 1

				c.Then("value should",
					bdd.Equal(2, value),
				)
			})

			t.When("add 1 again", func(b bdd.T) {
				value += 1

				t.Then("value should not be 2",
					bdd.NotEqual(2, value),
				)

				t.Then("value should be 3",
					bdd.Equal(3, value),
				)
			})
		})
	}))

	t.Run("snapshot", bdd.GivenT(func(b bdd.T) {
		b.Then("match",
			bdd.MatchSnapshot(
				func(s *snapshot.Snapshot) {
					s.Add("x.txt", []byte("1231"))
				},
				"test",
			),
		)
	}))
}
