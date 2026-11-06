package constant

// 资源位配置领域常量。
//
// 注意：这些枚举值同时以字符串形式硬编码在 biz/model/slot.go 的
// jsonschema:"enum=..." struct tag 中（Go 的 struct tag 无法引用常量），
// 修改此处时务必同步 model 里的 tag。

// 页面
const (
	PageHome       = "home"        // 首页
	PageNewHouse   = "new_house"   // 新房详情页
	PageSecondHand = "second_hand" // 二手房详情页
)

// 资源位类型
const (
	SlotTypeBanner = "banner" // Banner
	SlotTypeWidget = "widget" // 挂件
	SlotTypePopup  = "popup"  // 弹窗
)

// 生效小程序渠道（effect_source）
const (
	SourceTTMicroApp = "tt_microapp" // 字节小程序
	SourceWXMicroApp = "wx_microapp" // 微信小程序
)
