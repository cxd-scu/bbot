package guard

import "regexp"

// Result 注入检测结果。
type Result struct {
	Suspicious bool     // 是否疑似 prompt 注入
	Reasons    []string // 命中的规则
}

// 注入检测规则。刻意保持「高精度、低误报」：只匹配明确要求越权/改角色/探询系统提示的输入，
// 避免误伤正常运营请求（如「帮我配置首页banner」）。
var patterns = []struct {
	re     *regexp.Regexp
	reason string
}{
	{regexp.MustCompile(`(?i)忽略\s*(以上|上述|之前|此前|所有|先前)?\s*(的)?\s*(指令|指示|规则|设定|约束|提示)`), "要求忽略指令"},
	{regexp.MustCompile(`(?i)ignore\s+(all|any|the)?\s*(previous|above|prior|earlier)?\s*(instructions?|rules?|prompts?|guidelines?)`), "要求忽略指令（英文）"},
	{regexp.MustCompile(`(?i)(system\s*prompt|系统提示(词)?|初始指令)`), "探询系统提示"},
	{regexp.MustCompile(`(?i)(你现在是|请扮演|扮演一?个|假装你是|现在开始扮演|你的新身份是)`), "要求改变角色"},
	{regexp.MustCompile(`(?i)(忘记你|忘掉你|推翻你|清除你|忽略你)`), "要求清除/忽略设定"},
	{regexp.MustCompile(`(?i)(jailbreak|越狱|do\s*anything\s*now|\bDAN\b)`), "越狱关键词"},
	{regexp.MustCompile(`(?i)(override|bypass|绕过|无视)\s*(the\s+)?(指令|规则|限制|instruction|rule|restriction)`), "要求绕过规则"},
}

// Check 检测用户输入是否疑似 prompt 注入，返回检测结果。
func Check(input string) Result {
	res := Result{}
	for _, p := range patterns {
		if p.re.MatchString(input) {
			res.Suspicious = true
			res.Reasons = append(res.Reasons, p.reason)
		}
	}
	return res
}
