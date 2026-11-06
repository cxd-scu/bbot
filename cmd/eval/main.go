package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cxd-scu/bbot/biz/agent/chain"
	"github.com/cxd-scu/bbot/biz/agent/graph"
	"github.com/cxd-scu/bbot/biz/agent/tools"
	"github.com/cxd-scu/bbot/biz/eval"
	"github.com/cxd-scu/bbot/biz/infra/llm"
)

// 运行评测：
//
//	go run ./cmd/eval                 # 跑全部
//	go run ./cmd/eval -target=gateway # 只跑 gateway 路由
//	go run ./cmd/eval -target=slot    # 只跑 slot agent 端到端
func main() {
	ctx := context.Background()

	target := "all"
	if len(os.Args) > 1 {
		target = os.Args[1]
	}
	if len(target) > 0 && target[0] == '-' {
		// 兼容 -target=xxx 形式，简单解析
		target = os.Args[1]
		if len(target) > 8 && target[:8] == "-target=" {
			target = target[8:]
		} else {
			target = "all"
		}
	}

	// 初始化（与主服务一致）
	tools.InitTools()
	llm.InitModel(ctx)
	chain.InitChain(ctx)
	graph.InitGraph(ctx)

	switch target {
	case "gateway":
		runGateway(ctx)
	case "slot":
		runSlot(ctx)
	case "faithfulness":
		runFaithfulness(ctx)
	default:
		runGateway(ctx)
		fmt.Println()
		runSlot(ctx)
		fmt.Println()
		runFaithfulness(ctx)
	}
}

func runGateway(ctx context.Context) {
	report, err := eval.EvalGateway(ctx, eval.DefaultGatewayCases)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gateway eval error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(eval.FormatGatewayReport(report))
}

func runSlot(ctx context.Context) {
	report, err := eval.EvalSlotAgent(ctx, eval.DefaultSlotCases)
	if err != nil {
		fmt.Fprintf(os.Stderr, "slot eval error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(eval.FormatSlotReport(report))
}

func runFaithfulness(ctx context.Context) {
	report, err := eval.EvalFaithfulness(ctx, eval.DefaultFaithfulnessCases)
	if err != nil {
		fmt.Fprintf(os.Stderr, "faithfulness eval error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(eval.FormatFaithfulnessReport(report))
}
