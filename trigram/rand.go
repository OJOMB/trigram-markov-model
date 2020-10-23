package trigram

import (
	"github.com/mroth/weightedrand"
)

// WeightedRandom picks a single option from the provided at random but weighted
// by the mapped integer.
//
// For example, if you have the two options foo=1 and bar=2 then the chance
// that "bar" will be returned is twice as high as the change that "foo" is
// returned.
func WeightedRandom(options map[string]uint) string {
	choices := make([]weightedrand.Choice, 0, len(options))
	for word, weight := range options {
		choices = append(
			choices,
			weightedrand.Choice{
				Item:   word,
				Weight: weight,
			},
		)
	}

	c := weightedrand.NewChooser(choices...)
	return c.Pick().(string)
}
