# Branching Strategy

当前 fork 采用单主线模式：

- `custom/main`：唯一开发、构建、发布主线
- `upstream/main`：官方同步来源
- `origin/main`：历史镜像分支，保留兼容，不再作为开发或发布主线
- `custom/*`：短期功能、修复、文档分支，从 `custom/main` 切出

## Current Rules

1. 不直接在 `origin/main` 上做自定义提交
2. 日常开发默认基于 `custom/main`
3. `custom/*` 分支合入 `custom/main` 后尽快删除
4. 上游同步先从 `upstream/main` 合入 `custom/main`
5. 构建与自动发布默认绑定 `custom/main`

## Branch Roles

### `custom/main`

当前真实主线。

用途：
- 自定义功能集成
- 自动构建二进制
- 自动发布 GitHub Release
- 服务器部署基线

### `upstream/main`

官方主线，只用来同步上游更新。

用途：
- 查看上游新增功能或修复
- 定期同步到 `custom/main`

### `origin/main`

fork 历史遗留镜像。当前不承担发布职责。

建议：
- 保留以兼容旧链接或旧认知
- 不再直接开发

### `custom/*`

短周期工作分支。

推荐命名：
- `custom/feature-xxx`
- `custom/fix-xxx`
- `custom/docs-xxx`
- `custom/sync-upstream-yyyymmdd`

## Recommended Workflow

### 日常开发

```bash
git fetch origin
git checkout custom/main
git pull --ff-only origin custom/main

git checkout -b custom/feature-xxx
# 开发...
git push -u origin custom/feature-xxx
```

### 合回主线

```bash
git checkout custom/main
git pull --ff-only origin custom/main
git merge --no-ff custom/feature-xxx
git push origin custom/main
```

### 同步上游

```bash
git fetch upstream
git checkout custom/main
git merge upstream/main
git push origin custom/main
```

如果改动较大，先走临时同步分支：

```bash
git fetch upstream
git checkout -b custom/sync-upstream-20260402 origin/custom/main
git merge upstream/main
```

确认通过后再合回 `custom/main`。

## Branch Cleanup Guidance

建议长期只保留：

- `custom/main`
- `upstream/main`
- 少量仍在活跃使用的 `custom/*`

历史功能分支在内容已经进入主线后，应优先清理，避免误把旧专题分支当成当前基线。

## Current Repository Interpretation

按当前仓库状态理解：

- `custom/main`：当前唯一有效主线
- `custom/request-logs-schema`：历史专题分支，已被主线覆盖，不应继续作为基线
- `custom/ai-session-docs`：历史文档分支，内容已整理回主线后可删除

