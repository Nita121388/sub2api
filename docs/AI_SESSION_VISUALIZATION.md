# AI Session Visualization

状态：Partially Landed
更新时间：2026-04-02

## Goal

在 `custom/main` 主线上持续演进“AI 会话内容可视化”能力：

- 用户只能查看自己的请求会话
- 请求与响应内容可视化展示
- 保留元数据、正文、错误信息、性能指标
- 兼顾排障体验与主链路稳定性

## What Has Already Landed

基于当前主线状态，以下能力已经落地：

- 请求日志 schema 与 migration
- 请求日志元数据采集
- 请求日志列表接口与详情接口
- 正文持久化
- 用户侧请求日志列表页
- 用户侧详情页
- 会话时间线 / 可视化消息查看器
- 原始请求与响应正文回退显示

当前产品命名已回退为“会话详情”，但底层数据模型仍为 `request_logs`。

## Core Principles

1. 仅查看自己的会话
2. 默认先展示可理解的会话内容，再保留原始正文用于排障
3. 日志写入失败不阻断主请求
4. 大正文、流式响应必须支持截断或聚合策略

## Data Model

### `request_logs`

保存元数据，例如：

- `id`, `request_id`, `user_id`, `api_key_id`
- `model`, `endpoint`, `method`
- `status_code`, `error_code`, `error_message`
- `input_tokens`, `output_tokens`, `total_cost`
- `duration_ms`, `first_token_ms`
- `stream`, `ip_address`, `user_agent`
- `created_at`

### `request_log_payloads`

保存正文，例如：

- `request_body`
- `response_body`
- 正文字节数
- 截断标记

## Current UI Shape

### List Page

- 分页列表
- API Key / 模型 / 状态码 / 时间范围筛选
- 费用、时延、Token、创建时间展示
- 可进入单条详情

### Detail Page

- 概览卡片
- 会话时间线
- 请求 / 响应消息块
- 元数据、路由、性能、错误信息
- 原始请求体 / 响应体

## Remaining Work

以下项仍值得继续推进：

- 用户级保存策略配置
- 正文保留时长与容量上限
- 后台清理任务策略文档化
- 导出能力（JSON / Markdown）
- 管理员全局检索是否开放
- 更明确的敏感字段脱敏策略

## Suggested Next Phases

### Phase A

- 补全文档与配置说明
- 明确保存策略的默认值与权限边界

### Phase B

- 增加用户级日志设置页
- 增加保留策略 API

### Phase C

- 增加清理任务观测能力
- 增加导出与分享前脱敏逻辑

## Branch Note

本主题最初来自历史文档分支 `custom/ai-session-docs`，现已按当前 `custom/main` 主线策略重新整理。

