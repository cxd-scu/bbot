package eval

import (
	"testing"

	"github.com/cxd-scu/bbot/biz/constant"
	"github.com/cxd-scu/bbot/biz/model"
)

func boolPtr(b bool) *bool { return &b }

func TestCompareSetting(t *testing.T) {
	want := &model.SlotSetting{
		Page:         constant.PageHome,
		ResourceType: constant.SlotTypeBanner,
		CityIDs:      []int64{100000},
		ImageURL:     "https://x.com/a.png",
		OpenURL:      "sslocal://x",
		EffectSource: []string{constant.SourceTTMicroApp, constant.SourceWXMicroApp},
	}

	// 完全一致（城市顺序不同应视为相等）
	got := &model.SlotSetting{
		Page:         constant.PageHome,
		ResourceType: constant.SlotTypeBanner,
		CityIDs:      []int64{100000},
		ImageURL:     "https://x.com/a.png",
		OpenURL:      "sslocal://x",
		EffectSource: []string{constant.SourceWXMicroApp, constant.SourceTTMicroApp},
	}
	if diffs := compareSetting(got, want); len(diffs) != 0 {
		t.Fatalf("期望无差异，got=%v", diffs)
	}

	// 字段不一致应报差异
	bad := *got
	bad.Page = constant.PageNewHouse
	bad.CityIDs = []int64{200000}
	if diffs := compareSetting(&bad, want); len(diffs) == 0 {
		t.Fatal("期望报差异")
	}
}

func TestCompareSetting_OnlyWantSetFields(t *testing.T) {
	// want 只设 Page + ResourceType，其余不比对
	want := &model.SlotSetting{Page: constant.PageHome, ResourceType: constant.SlotTypeBanner}
	got := &model.SlotSetting{
		Page:         constant.PageHome,
		ResourceType: constant.SlotTypeBanner,
		CityIDs:      []int64{999}, // 不同，但 want 未设，不应报差异
	}
	if diffs := compareSetting(got, want); len(diffs) != 0 {
		t.Fatalf("want 未设的字段不应比对，got=%v", diffs)
	}
}

func TestEqualHelpers(t *testing.T) {
	if !equalInt64s([]int64{1, 2, 3}, []int64{3, 2, 1}) {
		t.Fatal("equalInt64s 应忽略顺序")
	}
	if equalInt64s([]int64{1, 2}, []int64{1, 3}) {
		t.Fatal("equalInt64s 不同元素应不等")
	}
	if !boolPtrEqual(boolPtr(true), boolPtr(true)) {
		t.Fatal("boolPtrEqual 同值应等")
	}
	if boolPtrEqual(boolPtr(true), boolPtr(false)) {
		t.Fatal("boolPtrEqual 不同值应不等")
	}
	if !boolPtrEqual(nil, nil) {
		t.Fatal("boolPtrEqual nil==nil 应等")
	}
}
