package tools

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cxd-scu/bbot/biz/infra/card"
	"github.com/cxd-scu/bbot/biz/model"
)

var (
	GenerateSettingTool tool.InvokableTool
	QuerySettingTool    tool.InvokableTool
	SubmitSettingTool   tool.InvokableTool
	NeedClarifyTool     tool.InvokableTool
	ValidateSettingTool tool.InvokableTool
)

func initSlotTools() error {
	var err error

	GenerateSettingTool, err = utils.InferTool(
		"generate_setting",
		"根据用户的资源位配置需求，生成符合对应资源位规则和 JSON 结构的配置内容。"+
			"适用于用户提出新增或修改资源位配置的请求。"+
			"输入为语义化配置（page + resource_type 唯一确定目标资源位，其余为城市、图片、跳转 URL、渠道、生效时间、频控等适用参数）。"+
			"不同资源位使用不同的配置结构和必填规则，必须根据目标资源位生成对应配置，不能混用不同类型的 JSON 结构。"+
			"首页 Banner 无需生效时间范围；新房和二手房详情页不支持弹窗。"+
			"当必要参数缺失或目标资源位无法唯一确定时，不得自行猜测或编造参数，应先调用 need_clarify 澄清。"+
			"返回生成后的配置内容。此工具只负责生成配置，不负责查询现有配置，也不负责提交工单或使配置生效。",
		generateSetting,
	)
	if err != nil {
		return err
	}

	QuerySettingTool, err = utils.InferTool(
		"query_setting",
		"查询配置中心中已有的资源位配置，用于了解当前配置内容、核对目标资源位以及为配置变更提供依据。"+
			"当用户询问当前资源位配置，或提出新增、修改配置的需求且需要确认现有配置时，应调用此工具。"+
			"输入为筛选条件（页面、资源位类型、城市、渠道均可选，越明确越精准）。"+
			"返回匹配到的配置列表及条数。如果配置中心只返回原始 JSON，应结合返回内容进一步识别城市、渠道、资源位及其配置规则。"+
			"不得将历史示例、模型记忆或推测当作当前线上配置。"+
			"此工具只负责查询，不负责生成新配置、修改配置或提交工单。",
		querySetting,
	)
	if err != nil {
		return err
	}

	SubmitSettingTool, err = utils.InferTool(
		"submit_setting",
		"将已经生成并确认的资源位配置提交至配置中心，创建配置变更工单。"+
			"仅当目标资源位明确、配置内容完整、适用规则校验通过，且满足要求的用户确认条件时调用。"+
			"提交前应确保配置与目标资源位的结构匹配，并检查必填字段、字段类型、城市范围、渠道、生效时间及频控等适用规则。"+
			"不得提交未经确认的猜测性配置，不得擅自补造城市 ID、图片链接、跳转 URL 或其他业务参数。"+
			"返回提交结果（success / ticket_id / message）。此工具的作用是创建变更工单，不代表配置已经审核通过或生效。"+
			"只有在配置中心明确返回成功后，才能告知用户工单已创建，并如实返回工单编号及相关状态；"+
			"如果提交失败，应反馈实际错误，不得声称提交成功。",
		submitSetting,
	)
	if err != nil {
		return err
	}

	NeedClarifyTool, err = utils.InferTool(
		"need_clarify",
		"判断当前资源位请求是否缺少完成查询或配置任务所必需的信息，或存在无法唯一确定的业务含义。"+
			"适用于无法确定目标页面、资源位类型、城市范围、图片链接、跳转 URL、适用渠道、"+
			"生效时间或其他当前资源位必需参数的情况。"+
			"输入为需要向用户澄清的问题列表（questions）。应结合用户当前消息和已有对话上下文判断，避免重复询问用户已经明确提供的信息。"+
			"只询问当前任务继续执行所必需的缺失信息，并提出具体、清晰的问题。"+
			"必须区分不同资源位的规则：例如首页 Banner 无需生效时间，新房和二手房详情页不支持弹窗。"+
			"不能因为某字段在其他资源位中必填，就要求所有资源位都提供该字段。"+
			"如果目标资源位和必需参数均已明确，则无需澄清，应继续执行相应任务。"+
			"此工具会暂停等待用户补充，用户补充后返回补全后的配置信息（SlotSetting）；不负责查询配置、生成配置或提交工单。",
		needClarify,
	)
	if err != nil {
		return err
	}

	ValidateSettingTool, err = utils.InferTool(
		"validate_setting",
		"校验生成的资源位配置是否符合对应资源位的业务规则和配置结构。"+
			"根据目标页面和资源位类型，检查配置的字段、类型、枚举值、城市 ID、渠道、生效时间及频控等适用约束。"+
			"首页、新房详情页、二手房详情页使用不同的配置结构和规则；新房和二手房详情页不支持弹窗。"+
			"对于首页 Banner，不应要求提供生效时间范围或频控参数。"+
			"返回校验结果（valid / errors）。校验失败时返回具体的错误字段和原因，供 Agent 修正配置或向用户说明。"+
			"此工具只负责校验，不负责补造缺失参数、查询配置或提交工单。",
		validateSetting,
	)
	if err != nil {
		return err
	}

	return nil
}

// TODO(实现): 将语义化 SlotSetting 映射为对应资源位的后端 JSON
// （由 page+resource_type 推导 resource_id / page_code / banner_name 或 resource_name / resource_type 数字码等）。
func generateSetting(ctx context.Context, req model.SlotSetting) (model.SlotSetting, error) {
	return req, nil
}

// TODO(实现): 接入配置中心查询，按筛选条件返回已有配置。
func querySetting(ctx context.Context, req model.QuerySettingReq) (model.QuerySettingResp, error) {
	return model.QuerySettingResp{}, nil
}

// TODO(实现): 接入配置中心提交，创建配置变更工单并返回工单号。
func submitSetting(ctx context.Context, req model.SlotSetting) (model.SubmitSettingResp, error) {
	return model.SubmitSettingResp{Success: true}, nil
}

func needClarify(ctx context.Context, req model.ClarifyReq) (*model.SlotSetting, error) {
	wasInterrupted, _, _ := tool.GetInterruptState[*model.AgentContext](ctx)
	if !wasInterrupted {
		card.ClarifyCard(ctx) // TODO: 将 req.Questions 传入卡片渲染
		return nil, tool.Interrupt(ctx, req)
	}

	isResume, _, data := tool.GetResumeContext[*model.SlotSetting](ctx)
	if !isResume {
		card.ClarifyCard(ctx)
		return nil, tool.Interrupt(ctx, req)
	}

	return data, nil
}

// validateSetting 校验配置是否符合目标资源位的业务规则，由确定性校验逻辑 ValidateSlotSetting 实现。
func validateSetting(ctx context.Context, req model.SlotSetting) (model.ValidateSettingResp, error) {
	return ValidateSlotSetting(req), nil
}
