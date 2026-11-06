package tools

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/cxd-scu/bbot/biz/constant"
	"github.com/cxd-scu/bbot/biz/model"
)

// timeLayout 生效时间字段的固定格式。
const timeLayout = "2006-01-02 15:04:05"

// ValidateSlotSetting 对一条资源位配置做确定性业务规则校验。
//
// 这是「Agent 决策 + Go 确定性执行」中的确定性兜底：模型负责理解意图、填语义字段，
// 业务规则（必填/禁止、枚举、类型、时间、城市）由本函数硬编码保证，不依赖模型自觉。
//
// 返回所有错误（不 fail-fast），Errors 为空且 Valid=true 表示通过。
//
// 规则总表（页面 × 类型）：
//   - page / resource_type / city_ids / image_url / open_url / effect_source：所有资源位必填。
//   - start_time / end_time：非「首页 Banner」必填；「首页 Banner」禁止出现。
//   - is_limit：「首页 Banner」禁止出现；其余可选（通常仅首页弹窗配置）。
//   - 新房 / 二手房详情页不支持弹窗。
func ValidateSlotSetting(s model.SlotSetting) model.ValidateSettingResp {
	var errs []string

	// 1. 资源位定位
	if s.Page == "" {
		errs = append(errs, "page 不能为空")
	} else if !isValidPage(s.Page) {
		errs = append(errs, fmt.Sprintf("page 取值非法: %q（允许: %s/%s/%s）", s.Page, constant.PageHome, constant.PageNewHouse, constant.PageSecondHand))
	}

	if s.ResourceType == "" {
		errs = append(errs, "resource_type 不能为空")
	} else if !isValidResourceType(s.ResourceType) {
		errs = append(errs, fmt.Sprintf("resource_type 取值非法: %q（允许: %s/%s/%s）", s.ResourceType, constant.SlotTypeBanner, constant.SlotTypeWidget, constant.SlotTypePopup))
	}

	// 2. 非法组合：详情页不支持弹窗
	if isValidPage(s.Page) && s.ResourceType == constant.SlotTypePopup && s.Page != constant.PageHome {
		errs = append(errs, fmt.Sprintf("%s详情页不支持弹窗（page=%s, resource_type=%s）", pageLabel(s.Page), s.Page, s.ResourceType))
	}

	// 3. 开放城市
	if len(s.CityIDs) == 0 {
		errs = append(errs, "city_ids 不能为空")
	} else {
		for _, id := range s.CityIDs {
			if id <= 0 {
				errs = append(errs, fmt.Sprintf("city_ids 含非法城市 ID: %d（应为正整数）", id))
			}
		}
	}

	// 4. 素材
	if strings.TrimSpace(s.ImageURL) == "" {
		errs = append(errs, "image_url 不能为空")
	} else if !isHTTPURL(s.ImageURL) {
		errs = append(errs, fmt.Sprintf("image_url 不是合法的 http(s) 链接: %q", s.ImageURL))
	}

	if strings.TrimSpace(s.OpenURL) == "" {
		errs = append(errs, "open_url 不能为空")
	} else if !hasURLScheme(s.OpenURL) {
		errs = append(errs, fmt.Sprintf("open_url 不是合法的跳转链接: %q", s.OpenURL))
	}

	// 5. 生效渠道
	if len(s.EffectSource) == 0 {
		errs = append(errs, "effect_source 不能为空")
	} else {
		seen := make(map[string]bool, len(s.EffectSource))
		for _, src := range s.EffectSource {
			if !isValidSource(src) {
				errs = append(errs, fmt.Sprintf("effect_source 含非法渠道: %q（允许: %s/%s）", src, constant.SourceTTMicroApp, constant.SourceWXMicroApp))
				continue
			}
			if seen[src] {
				errs = append(errs, fmt.Sprintf("effect_source 渠道重复: %q", src))
			}
			seen[src] = true
		}
	}

	// 6. 生效时间与频控（首页 Banner 特殊）
	validateTimeAndLimit(s, &errs)

	return model.ValidateSettingResp{
		Valid:  len(errs) == 0,
		Errors: errs,
	}
}

// validateTimeAndLimit 校验生效时间与频控字段，规则因资源位而异。
func validateTimeAndLimit(s model.SlotSetting, errs *[]string) {
	isHomeBanner := s.Page == constant.PageHome && s.ResourceType == constant.SlotTypeBanner

	if isHomeBanner {
		// 首页 Banner 无时间、无频控
		if strings.TrimSpace(s.StartTime) != "" || strings.TrimSpace(s.EndTime) != "" {
			*errs = append(*errs, "首页 Banner 不支持生效时间范围（start_time/end_time 应留空）")
		}
		if s.IsLimit != nil {
			*errs = append(*errs, "首页 Banner 不支持用户频控（is_limit 应不填）")
		}
		return
	}

	// 非首页 Banner：生效时间必填
	startEmpty := strings.TrimSpace(s.StartTime) == ""
	endEmpty := strings.TrimSpace(s.EndTime) == ""
	switch {
	case startEmpty && endEmpty:
		*errs = append(*errs, "生效时间范围不完整：start_time 与 end_time 均需提供")
	case startEmpty:
		*errs = append(*errs, "start_time 不能为空")
	case endEmpty:
		*errs = append(*errs, "end_time 不能为空")
	default:
		st, err1 := time.Parse(timeLayout, s.StartTime)
		et, err2 := time.Parse(timeLayout, s.EndTime)
		if err1 != nil {
			*errs = append(*errs, fmt.Sprintf("start_time 格式错误: %q（应为 %s）", s.StartTime, timeLayout))
		}
		if err2 != nil {
			*errs = append(*errs, fmt.Sprintf("end_time 格式错误: %q（应为 %s）", s.EndTime, timeLayout))
		}
		if err1 == nil && err2 == nil && !st.Before(et) {
			*errs = append(*errs, "start_time 不能晚于 end_time")
		}
	}
}

func isValidPage(p string) bool {
	switch p {
	case constant.PageHome, constant.PageNewHouse, constant.PageSecondHand:
		return true
	}
	return false
}

func isValidResourceType(t string) bool {
	switch t {
	case constant.SlotTypeBanner, constant.SlotTypeWidget, constant.SlotTypePopup:
		return true
	}
	return false
}

func isValidSource(src string) bool {
	return src == constant.SourceTTMicroApp || src == constant.SourceWXMicroApp
}

// isHTTPURL 校验是否为合法的 http/https 链接。
func isHTTPURL(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

// hasURLScheme 校验是否为带 scheme 的合法链接（如 sslocal://、http(s)://）。
func hasURLScheme(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}
	return u.Scheme != ""
}

func pageLabel(p string) string {
	switch p {
	case constant.PageHome:
		return "首页"
	case constant.PageNewHouse:
		return "新房"
	case constant.PageSecondHand:
		return "二手房"
	}
	return p
}
