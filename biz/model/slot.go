package model

// SlotSetting 表示一条小程序资源位配置，是 slot agent 各工具之间传递的核心契约对象。
//
// 三种资源位类型（banner=Banner、widget=挂件、popup=弹窗）× 三个页面
// （首页、新房详情页、二手房详情页），理论 9 种组合，但新房/二手房详情页
// 不支持弹窗，实际可配置资源位共 7 种。
//
// 契约设计原则：本结构只承载「语义字段」，即模型能理解、用户能自然语言提供的
// 信息（目标页面/类型、城市、图片、跳转 URL、生效时间、频控等）。
// resource_id、page_code、banner_name/resource_name、resource_type 数字码等
// 「后端标识」由 Go 侧确定性推导，模型不接触、也不得臆造。
//
// 哪些字段必填/禁止由 validate_setting 在 Go 侧确定性校验，这里只描述字段语义，
// 让模型能稳定按语义填参，而不是自行推断业务规则。
type SlotSetting struct {
	// —— 资源位定位（必填，用于唯一确定目标资源位）——
	Page         string `json:"page" jsonschema_description:"目标页面：home=首页, new_house=新房详情页, second_hand=二手房详情页" jsonschema:"enum=home,enum=new_house,enum=second_hand"`
	ResourceType string `json:"resource_type" jsonschema_description:"资源位类型：banner=Banner, widget=挂件, popup=弹窗（新房/二手房详情页不支持弹窗）" jsonschema:"enum=banner,enum=widget,enum=popup"`

	// —— 城市与渠道 ——
	CityIDs      []int64  `json:"city_ids,omitempty" jsonschema_description:"开放城市 ID 列表；全部城市生效需加入所有 city_id"`
	EffectSource []string `json:"effect_source,omitempty" jsonschema_description:"生效的小程序渠道，可多选，缺省通常表示全渠道" jsonschema:"enum=tt_microapp,enum=wx_microapp"`

	// —— 素材 ——
	ImageURL string `json:"image_url,omitempty" jsonschema_description:"图片链接"`
	OpenURL  string `json:"open_url,omitempty" jsonschema_description:"跳转 URL，小程序内跳转通常为 sslocal:// 协议"`

	// —— 生效时间（首页 Banner 无需）——
	StartTime string `json:"start_time,omitempty" jsonschema_description:"生效开始时间，格式 2006-01-02 15:04:05"`
	EndTime   string `json:"end_time,omitempty" jsonschema_description:"生效结束时间，格式 2006-01-02 15:04:05"`

	// —— 用户频控（首页 Banner 无需，通常仅首页弹窗配置）——
	IsLimit *bool `json:"is_limit,omitempty" jsonschema_description:"是否需要对用户做频控；true=需要, false=不需要"`
}

// QuerySettingReq 查询配置中心的筛选条件，所有字段均可选。
type QuerySettingReq struct {
	Page         string `json:"page,omitempty" jsonschema_description:"按页面筛选：home=首页, new_house=新房详情页, second_hand=二手房详情页" jsonschema:"enum=home,enum=new_house,enum=second_hand"`
	ResourceType string `json:"resource_type,omitempty" jsonschema_description:"按资源位类型筛选：banner=Banner, widget=挂件, popup=弹窗" jsonschema:"enum=banner,enum=widget,enum=popup"`
	CityID       *int64 `json:"city_id,omitempty" jsonschema_description:"按城市 ID 筛选"`
	EffectSource string `json:"effect_source,omitempty" jsonschema_description:"按生效渠道筛选" jsonschema:"enum=tt_microapp,enum=wx_microapp"`
}

// QuerySettingResp 查询配置中心的返回结果。
type QuerySettingResp struct {
	Total    int            `json:"total" jsonschema_description:"匹配到的配置条数"`
	Settings []*SlotSetting `json:"settings" jsonschema_description:"匹配到的资源配置列表"`
}

// ValidateSettingResp 配置校验的返回结果。
type ValidateSettingResp struct {
	Valid  bool     `json:"valid" jsonschema_description:"配置是否通过校验"`
	Errors []string `json:"errors,omitempty" jsonschema_description:"校验不通过时的错误字段及原因列表，Valid 为 true 时为空"`
}

// SubmitSettingResp 配置提交的返回结果。
type SubmitSettingResp struct {
	Success  bool   `json:"success" jsonschema_description:"是否提交成功"`
	TicketID string `json:"ticket_id,omitempty" jsonschema_description:"配置变更工单编号，提交成功后返回"`
	Message  string `json:"message,omitempty" jsonschema_description:"提交结果说明，失败时为具体错误原因"`
}

// ClarifyReq need_clarify 工具的输入：需要向用户澄清的问题列表。
type ClarifyReq struct {
	Questions []string `json:"questions" jsonschema_description:"需要向用户澄清的问题列表，只列继续执行所必需的信息，避免重复询问已明确的内容"`
}
