package graph

import "context"

func InitGraph(ctx context.Context) {
	err := initGatewayChain(ctx)
	if err != nil {
		panic(err)
	}

	err = initSlotAgentRunner(ctx)
	if err != nil {
		panic(err)
	}

	err = initOncallAgent(ctx)
	if err != nil {
		panic(err)
	}
}
