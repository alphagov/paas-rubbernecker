package rubbernecker_test

import (
	"github.com/alphagov/paas-rubbernecker/pkg/rubbernecker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filter", func() {
	Describe("blocked filter", func(){
		It("is applied when the queries contain 'sticker:blocked'", func(){
			queries := []string {"bar", "sticker:blocked", "foo"}

			filter := rubbernecker.BlockedFilter{}
			actual := filter.IsApplied(queries)

			Expect(actual).To(BeTrue())
		})
	})

	Describe("scheduled filter", func() {
		It("is applied when the queries contain 'sticker:scheduled'", func(){
			queries := []string {"bar", "sticker:scheduled", "foo"}

			filter := rubbernecker.ScheduledFilter{}
			actual := filter.IsApplied(queries)

			Expect(actual).To(BeTrue())
		})
	})

	Describe("comments to resolve filter", func(){
		It("is applied when the queries contain 'sticker:comments-to-resolve'", func(){
			queries := []string {"bar", "sticker:comments-to-resolve", "foo"}

			filter := rubbernecker.CommentsToResolveFilter{}
			actual := filter.IsApplied(queries)

			Expect(actual).To(BeTrue())
		})
	})

	Describe("small task filter", func(){
		It("is applied when the queries contain 'sticker:'small' task'", func(){
			queries := []string {"bar", "sticker:'small' task", "foo"}

			filter := rubbernecker.SmallTaskFilter{}
			actual := filter.IsApplied(queries)

			Expect(actual).To(BeTrue())
		})
	})
	
	Describe("pairing filter", func(){
		It("is applied when the queries contain 'sticker:pairing'", func(){
			queries := []string {"bar", "sticker:pairing", "foo"}

			filter := rubbernecker.PairingFilter{}
			actual := filter.IsApplied(queries)

			Expect(actual).To(BeTrue())
		})
	})

	Describe("non-tech filter", func(){
		It("is applied when the queries contain 'sticker:non-tech'", func(){
			queries := []string {"bar", "sticker:non-tech", "foo"}

			filter := rubbernecker.NonTechFilter{}
			actual := filter.IsApplied(queries)

			Expect(actual).To(BeTrue())
		})
	})

	Describe("tech filter", func(){
		It("is applied when the queries contain 'not-sticker:non-tech'", func(){
			queries := []string {"bar", "not-sticker:non-tech", "foo"}

			filter := rubbernecker.TechFilter{}
			actual := filter.IsApplied(queries)

			Expect(actual).To(BeTrue())
		})
	})
})
