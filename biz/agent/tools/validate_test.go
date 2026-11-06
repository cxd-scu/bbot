package tools

import (
	"strings"
	"testing"

	"github.com/cxd-scu/bbot/biz/constant"
	"github.com/cxd-scu/bbot/biz/model"
)

func boolPtr(b bool) *bool { return &b }

// 构造一个合法的「首页 Banner」配置（无时间、无频控）。
func validHomeBanner() model.SlotSetting {
	return model.SlotSetting{
		Page:         constant.PageHome,
		ResourceType: constant.SlotTypeBanner,
		CityIDs:      []int64{100000},
		ImageURL:     "https://example.com/banner.png",
		OpenURL:      "sslocal://micro_webview?title=x&url=https://m.example.com",
		EffectSource: []string{constant.SourceTTMicroApp, constant.SourceWXMicroApp},
	}
}

// 构造一个合法的「新房详情页 Banner」配置（含时间）。
func validDetailBanner() model.SlotSetting {
	return model.SlotSetting{
		Page:         constant.PageNewHouse,
		ResourceType: constant.SlotTypeBanner,
		CityIDs:      []int64{100000},
		ImageURL:     "https://example.com/banner.png",
		OpenURL:      "sslocal://micro_webview?title=x",
		EffectSource: []string{constant.SourceTTMicroApp},
		StartTime:    "2026-01-01 00:00:00",
		EndTime:      "2026-01-31 23:59:59",
	}
}

func TestValidateSlotSetting(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*model.SlotSetting)
		wantErr []string // 期望错误信息中包含的子串，为空表示期望通过
	}{
		{
			name:    "首页 Banner 合法",
			mutate:  func(s *model.SlotSetting) {},
			wantErr: nil,
		},
		{
			name:    "新房详情页 Banner 合法",
			mutate:  func(s *model.SlotSetting) { *s = validDetailBanner() },
			wantErr: nil,
		},
		{
			name: "首页挂件合法（含时间）",
			mutate: func(s *model.SlotSetting) {
				s.ResourceType = constant.SlotTypeWidget
				s.StartTime = "2026-01-01 00:00:00"
				s.EndTime = "2026-01-31 23:59:59"
			},
			wantErr: nil,
		},
		{
			name: "首页弹窗合法（含频控）",
			mutate: func(s *model.SlotSetting) {
				s.ResourceType = constant.SlotTypePopup
				s.StartTime = "2026-01-01 00:00:00"
				s.EndTime = "2026-01-31 23:59:59"
				s.IsLimit = boolPtr(true)
			},
			wantErr: nil,
		},
		{
			name:    "page 缺失",
			mutate:  func(s *model.SlotSetting) { s.Page = "" },
			wantErr: []string{"page 不能为空"},
		},
		{
			name:    "page 非法",
			mutate:  func(s *model.SlotSetting) { s.Page = "foo" },
			wantErr: []string{"page 取值非法"},
		},
		{
			name:    "resource_type 缺失",
			mutate:  func(s *model.SlotSetting) { s.ResourceType = "" },
			wantErr: []string{"resource_type 不能为空"},
		},
		{
			name:    "resource_type 非法",
			mutate:  func(s *model.SlotSetting) { s.ResourceType = "foo" },
			wantErr: []string{"resource_type 取值非法"},
		},
		{
			name: "新房详情页弹窗（不支持）",
			mutate: func(s *model.SlotSetting) {
				s.Page = constant.PageNewHouse
				s.ResourceType = constant.SlotTypePopup
				s.StartTime = "2026-01-01 00:00:00"
				s.EndTime = "2026-01-31 23:59:59"
			},
			wantErr: []string{"新房详情页不支持弹窗"},
		},
		{
			name: "二手房详情页弹窗（不支持）",
			mutate: func(s *model.SlotSetting) {
				s.Page = constant.PageSecondHand
				s.ResourceType = constant.SlotTypePopup
				s.StartTime = "2026-01-01 00:00:00"
				s.EndTime = "2026-01-31 23:59:59"
			},
			wantErr: []string{"二手房详情页不支持弹窗"},
		},
		{
			name:    "city_ids 缺失",
			mutate:  func(s *model.SlotSetting) { s.CityIDs = nil },
			wantErr: []string{"city_ids 不能为空"},
		},
		{
			name:    "city_ids 含非法值",
			mutate:  func(s *model.SlotSetting) { s.CityIDs = []int64{100000, 0} },
			wantErr: []string{"city_ids 含非法城市 ID"},
		},
		{
			name:    "image_url 缺失",
			mutate:  func(s *model.SlotSetting) { s.ImageURL = "" },
			wantErr: []string{"image_url 不能为空"},
		},
		{
			name:    "image_url 非 http",
			mutate:  func(s *model.SlotSetting) { s.ImageURL = "ftp://example.com/a.png" },
			wantErr: []string{"image_url 不是合法的 http(s) 链接"},
		},
		{
			name:    "open_url 缺失",
			mutate:  func(s *model.SlotSetting) { s.OpenURL = "" },
			wantErr: []string{"open_url 不能为空"},
		},
		{
			name:    "open_url 无 scheme",
			mutate:  func(s *model.SlotSetting) { s.OpenURL = "not-a-url" },
			wantErr: []string{"open_url 不是合法的跳转链接"},
		},
		{
			name:    "effect_source 缺失",
			mutate:  func(s *model.SlotSetting) { s.EffectSource = nil },
			wantErr: []string{"effect_source 不能为空"},
		},
		{
			name:    "effect_source 含非法值",
			mutate:  func(s *model.SlotSetting) { s.EffectSource = []string{"tt_microapp", "foo"} },
			wantErr: []string{"effect_source 含非法渠道"},
		},
		{
			name:    "effect_source 重复",
			mutate:  func(s *model.SlotSetting) { s.EffectSource = []string{"tt_microapp", "tt_microapp"} },
			wantErr: []string{"effect_source 渠道重复"},
		},
		{
			name: "首页 Banner 带时间（禁止）",
			mutate: func(s *model.SlotSetting) {
				s.StartTime = "2026-01-01 00:00:00"
				s.EndTime = "2026-01-31 23:59:59"
			},
			wantErr: []string{"首页 Banner 不支持生效时间范围"},
		},
		{
			name:    "首页 Banner 带频控（禁止）",
			mutate:  func(s *model.SlotSetting) { s.IsLimit = boolPtr(true) },
			wantErr: []string{"首页 Banner 不支持用户频控"},
		},
		{
			name: "详情页缺 start_time",
			mutate: func(s *model.SlotSetting) {
				*s = validDetailBanner()
				s.StartTime = ""
			},
			wantErr: []string{"start_time 不能为空"},
		},
		{
			name: "详情页缺时间范围",
			mutate: func(s *model.SlotSetting) {
				*s = validDetailBanner()
				s.StartTime = ""
				s.EndTime = ""
			},
			wantErr: []string{"生效时间范围不完整"},
		},
		{
			name: "start_time 格式错误",
			mutate: func(s *model.SlotSetting) {
				*s = validDetailBanner()
				s.StartTime = "2026/01/01"
			},
			wantErr: []string{"start_time 格式错误"},
		},
		{
			name: "start_time 晚于 end_time",
			mutate: func(s *model.SlotSetting) {
				*s = validDetailBanner()
				s.StartTime = "2026-02-01 00:00:00"
				s.EndTime = "2026-01-01 00:00:00"
			},
			wantErr: []string{"start_time 不能晚于 end_time"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := validHomeBanner()
			tt.mutate(&s)

			resp := ValidateSlotSetting(s)

			if len(tt.wantErr) == 0 {
				if !resp.Valid {
					t.Fatalf("期望通过，实际失败: %v", resp.Errors)
				}
				return
			}

			if resp.Valid {
				t.Fatalf("期望失败，实际通过")
			}
			for _, want := range tt.wantErr {
				if !containsAny(resp.Errors, want) {
					t.Fatalf("期望错误包含 %q，实际: %v", want, resp.Errors)
				}
			}
		})
	}
}

func containsAny(errs []string, substr string) bool {
	for _, e := range errs {
		if strings.Contains(e, substr) {
			return true
		}
	}
	return false
}
