package eval

import (
	"github.com/cxd-scu/bbot/biz/constant"
	"github.com/cxd-scu/bbot/biz/model"
)

// SlotEvalCase 一条 slot agent 端到端评测用例。
type SlotEvalCase struct {
	Name  string // 用例名
	Query string // 用户输入

	// WantClarify 为 true：期望 agent 触发澄清（缺参数/不支持组合），而非编造参数直接生成。
	WantClarify bool

	// WantSetting 非 nil：期望 agent 生成的配置，做字段级比对。
	WantSetting *model.SlotSetting
}

// DefaultSlotCases 默认的 slot agent 端到端评测集。
// 重点验证「缺参数澄清不编造」「信息齐全生成正确配置」「不支持组合拒绝」三类行为。
var DefaultSlotCases = []SlotEvalCase{
	// —— 缺参数 → 应澄清，不编造 ——
	{Name: "缺图片链接应澄清", Query: "帮我配一个首页banner，城市100000", WantClarify: true},
	{Name: "缺跳转链接应澄清", Query: "首页banner，城市100000，图片https://x.com/a.png", WantClarify: true},

	// —— 信息齐全 → 应生成正确配置 ——
	{
		Name:  "首页banner信息齐全",
		Query: "首页banner，城市100000，图片https://x.com/a.png，跳转sslocal://micro_webview?url=x，字节和微信都上",
		WantSetting: &model.SlotSetting{
			Page:         constant.PageHome,
			ResourceType: constant.SlotTypeBanner,
			CityIDs:      []int64{100000},
			ImageURL:     "https://x.com/a.png",
			OpenURL:      "sslocal://micro_webview?url=x",
			EffectSource: []string{constant.SourceTTMicroApp, constant.SourceWXMicroApp},
		},
	},

	// —— 不支持组合 → 应澄清/拒绝，而非生成 ——
	{Name: "新房详情页弹窗不支持", Query: "新房详情页配一个弹窗", WantClarify: true},
	{Name: "二手房详情页弹窗不支持", Query: "二手房详情页加个弹窗", WantClarify: true},
}
