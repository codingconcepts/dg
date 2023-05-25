package generator

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/samber/lo"
)

type weightedItem struct {
	Value  string
	Weight int
}

type weightedItems struct {
	items       []weightedItem
	totalWeight int
}

func makeWeightedItems(items []weightedItem) weightedItems {
	wi := weightedItems{
		items: items,
	}

	wi.totalWeight = lo.SumBy(items, func(wi weightedItem) int {
		return wi.Weight
	})

	return wi
}

func (wi weightedItems) choose() string {
	randomWeight := gofakeit.IntRange(1, wi.totalWeight)
	for _, i := range wi.items {
		randomWeight -= i.Weight
		if randomWeight <= 0 {
			return i.Value
		}
	}

	return ""
}
