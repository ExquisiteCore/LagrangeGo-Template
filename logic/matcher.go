package logic

import (
	"regexp"
	"strings"

	"github.com/LagrangeDev/LagrangeGo/message"
)

// Matcher 匹配器接口
type Matcher interface {
	Match(ctx *MessageContext) bool
}

// MessageTypeMatcher 消息类型匹配器
type MessageTypeMatcher struct {
	MessageType string
}

func (m *MessageTypeMatcher) Match(ctx *MessageContext) bool {
	switch m.MessageType {
	case "private":
		_, ok := ctx.GetPrivateMessage()
		return ok
	case "group":
		_, ok := ctx.GetGroupMessage()
		return ok
	case "friend_request":
		_, ok := ctx.GetFriendRequest()
		return ok
	default:
		return false
	}
}

// NewMessageTypeMatcher 创建消息类型匹配器
func NewMessageTypeMatcher(msgType string) *MessageTypeMatcher {
	return &MessageTypeMatcher{MessageType: msgType}
}

// TextMatcher 文本匹配器
type TextMatcher struct {
	Pattern   string
	CaseSensitive bool
}

func (m *TextMatcher) Match(ctx *MessageContext) bool {
	text := ctx.GetMessageText()
	if text == "" {
		return false
	}
	
	if !m.CaseSensitive {
		return strings.Contains(strings.ToLower(text), strings.ToLower(m.Pattern))
	}
	return strings.Contains(text, m.Pattern)
}

// NewTextMatcher 创建文本匹配器
func NewTextMatcher(pattern string, caseSensitive bool) *TextMatcher {
	return &TextMatcher{Pattern: pattern, CaseSensitive: caseSensitive}
}

// RegexMatcher 正则表达式匹配器
type RegexMatcher struct {
	regex *regexp.Regexp
}

func (m *RegexMatcher) Match(ctx *MessageContext) bool {
	text := ctx.GetMessageText()
	if text == "" {
		return false
	}
	return m.regex.MatchString(text)
}

// NewRegexMatcher 创建正则表达式匹配器
func NewRegexMatcher(pattern string) (*RegexMatcher, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &RegexMatcher{regex: regex}, nil
}

// PrefixMatcher 前缀匹配器
type PrefixMatcher struct {
	Prefix string
}

func (m *PrefixMatcher) Match(ctx *MessageContext) bool {
	text := ctx.GetMessageText()
	return strings.HasPrefix(text, m.Prefix)
}

// NewPrefixMatcher 创建前缀匹配器
func NewPrefixMatcher(prefix string) *PrefixMatcher {
	return &PrefixMatcher{Prefix: prefix}
}

// SenderMatcher 发送者匹配器
type SenderMatcher struct {
	UserIDs []uint32
}

func (m *SenderMatcher) Match(ctx *MessageContext) bool {
	var senderID uint32
	
	if privateMsg, ok := ctx.GetPrivateMessage(); ok {
		senderID = privateMsg.Sender.Uin
	} else if groupMsg, ok := ctx.GetGroupMessage(); ok {
		senderID = groupMsg.Sender.Uin
	} else {
		return false
	}
	
	for _, userID := range m.UserIDs {
		if userID == senderID {
			return true
		}
	}
	return false
}

// NewSenderMatcher 创建发送者匹配器
func NewSenderMatcher(userIDs ...uint32) *SenderMatcher {
	return &SenderMatcher{UserIDs: userIDs}
}

// GroupMatcher 群聊匹配器
type GroupMatcher struct {
	GroupIDs []uint32
}

func (m *GroupMatcher) Match(ctx *MessageContext) bool {
	groupMsg, ok := ctx.GetGroupMessage()
	if !ok {
		return false
	}
	
	for _, groupID := range m.GroupIDs {
		if groupID == groupMsg.GroupUin {
			return true
		}
	}
	return false
}

// NewGroupMatcher 创建群聊匹配器
func NewGroupMatcher(groupIDs ...uint32) *GroupMatcher {
	return &GroupMatcher{GroupIDs: groupIDs}
}

// CommandMatcher 命令匹配器
type CommandMatcher struct {
	Commands []string
	Prefix   string
}

func (m *CommandMatcher) Match(ctx *MessageContext) bool {
	text := ctx.GetMessageText()
	if text == "" {
		return false
	}
	
	// 检查前缀
	if m.Prefix != "" && !strings.HasPrefix(text, m.Prefix) {
		return false
	}
	
	// 移除前缀
	if m.Prefix != "" {
		text = strings.TrimPrefix(text, m.Prefix)
	}
	
	// 分割命令和参数
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return false
	}
	
	command := parts[0]
	
	// 检查命令是否匹配
	for _, cmd := range m.Commands {
		if cmd == command {
			// 将命令和参数保存到上下文
			ctx.Set("command", command)
			if len(parts) > 1 {
				ctx.Set("args", parts[1:])
			}
			return true
		}
	}
	
	return false
}

// NewCommandMatcher 创建命令匹配器
func NewCommandMatcher(prefix string, commands ...string) *CommandMatcher {
	return &CommandMatcher{Commands: commands, Prefix: prefix}
}

// AndMatcher 逻辑AND匹配器
type AndMatcher struct {
	Matchers []Matcher
}

func (m *AndMatcher) Match(ctx *MessageContext) bool {
	for _, matcher := range m.Matchers {
		if !matcher.Match(ctx) {
			return false
		}
	}
	return true
}

// NewAndMatcher 创建逻辑AND匹配器
func NewAndMatcher(matchers ...Matcher) *AndMatcher {
	return &AndMatcher{Matchers: matchers}
}

// OrMatcher 逻辑OR匹配器
type OrMatcher struct {
	Matchers []Matcher
}

func (m *OrMatcher) Match(ctx *MessageContext) bool {
	for _, matcher := range m.Matchers {
		if matcher.Match(ctx) {
			return true
		}
	}
	return false
}

// NewOrMatcher 创建逻辑OR匹配器
func NewOrMatcher(matchers ...Matcher) *OrMatcher {
	return &OrMatcher{Matchers: matchers}
}

// NotMatcher 逻辑NOT匹配器
type NotMatcher struct {
	Matcher Matcher
}

func (m *NotMatcher) Match(ctx *MessageContext) bool {
	return !m.Matcher.Match(ctx)
}

// NewNotMatcher 创建逻辑NOT匹配器
func NewNotMatcher(matcher Matcher) *NotMatcher {
	return &NotMatcher{Matcher: matcher}
}

// CustomMatcher 自定义匹配器
type CustomMatcher struct {
	MatchFunc func(*MessageContext) bool
}

func (m *CustomMatcher) Match(ctx *MessageContext) bool {
	return m.MatchFunc(ctx)
}

// NewCustomMatcher 创建自定义匹配器
func NewCustomMatcher(matchFunc func(*MessageContext) bool) *CustomMatcher {
	return &CustomMatcher{MatchFunc: matchFunc}
}

// AtMatcher @消息匹配器
type AtMatcher struct {
	BotID uint32
}

func (m *AtMatcher) Match(ctx *MessageContext) bool {
	groupMsg, ok := ctx.GetGroupMessage()
	if !ok {
		return false
	}
	
	// 检查消息中是否包含@机器人
	for _, element := range groupMsg.Elements {
		if atElement, ok := element.(*message.AtElement); ok {
			if atElement.TargetUin == m.BotID {
				return true
			}
		}
	}
	
	return false
}

// NewAtMatcher 创建@匹配器
func NewAtMatcher(botID uint32) *AtMatcher {
	return &AtMatcher{BotID: botID}
}