package eval

import (
	"context"

	"github.com/cxd-scu/bbot/biz/app"
)

// EvalGateway 运行 gateway 路由评测集，返回评测报告。
// 逐条调用 GatewayNluFull，对比实际路由与期望路由。
func EvalGateway(ctx context.Context, cases []GatewayCase) (*GatewayReport, error) {
	report := &GatewayReport{
		Total: len(cases),
	}

	for _, c := range cases {
		out, err := app.GatewayNluFull(ctx, c.Query)
		if err != nil {
			report.Errors = append(report.Errors, GatewayCaseError{Case: c, Err: err})
			continue
		}

		res := GatewayCaseResult{
			Case:       c,
			GotAgent:   out.TargetAgent,
			Confidence: out.Confidence,
			Reason:     out.Reason,
			Correct:    out.TargetAgent == c.WantAgent,
		}
		report.Details = append(report.Details, res)
		if res.Correct {
			report.Correct++
		}
	}

	if report.Total > 0 {
		report.Accuracy = float64(report.Correct) / float64(report.Total)
	}
	return report, nil
}
