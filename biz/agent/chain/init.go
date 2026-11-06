package chain

import "context"

func InitChain(ctx context.Context) {
	err := initShortMemorySummarizer(ctx)
	if err != nil {
		panic(err)
	}

	err = initImageAnalyzer(ctx)
	if err != nil {
		panic(err)
	}
}
