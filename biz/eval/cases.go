package eval

import "github.com/cxd-scu/bbot/biz/constant"

// GatewayCase 一条 gateway 路由评测用例。
type GatewayCase struct {
	Name      string // 用例名
	Query     string // 用户输入
	WantAgent string // 期望路由到的 agent（slot / oncall / need_classify）
}

// DefaultGatewayCases 默认的 gateway 路由评测集（黄金用例）。
// 覆盖 slot（资源位配置）、oncall（信息查询）、need_classify（无法归类）及边界表达。
// 后续可迁移到 JSON/YAML 文件由非工程同学维护。
var DefaultGatewayCases = []GatewayCase{
	// slot —— 资源位配置
	{"首页banner配置", "帮我在首页配一个banner", constant.SlotAgent},
	{"首页banner英文", "首页 banner 配置", constant.SlotAgent},
	{"新房挂件查询", "查一下新房详情页的挂件现在配的是什么", constant.SlotAgent},
	{"二手房banner修改", "二手房详情页的banner改成这张图", constant.SlotAgent},
	{"首页弹窗频控", "首页弹窗要加用户频控", constant.SlotAgent},
	{"提交配置工单", "帮我把这个配置提交成变更工单", constant.SlotAgent},
	{"带参数的配置", "首页banner，城市100000，图片链接是xxx", constant.SlotAgent},
	{"模糊配置意图", "配一下", constant.SlotAgent},
	{"查询现有配置", "给我看看首页现在配了什么", constant.SlotAgent},

	// oncall —— 信息查询（RAG）
	{"谁值班", "今天谁值班", constant.OncallAgent},
	{"值班安排", "值班安排是什么", constant.OncallAgent},
	{"告警处理", "收到告警了怎么处理", constant.OncallAgent},
	{"流程咨询", "故障上报的流程是什么", constant.OncallAgent},
	{"手册查询", "操作手册里怎么写交接班", constant.OncallAgent},
	{"oncall关键词", "oncall", constant.OncallAgent},

	// need_classify —— 无法归类
	{"天气闲聊", "今天天气怎么样", constant.NeedClassify},
	{"写诗", "帮我写一首诗", constant.NeedClassify},
	{"打招呼", "你好呀", constant.NeedClassify},
	{"无关业务", "帮我订个机票", constant.NeedClassify},
}
