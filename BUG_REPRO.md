# Bug Reproduction

## 包的性质

当前 test_model_fix 保存的是被测模型修复后的结果源码，不是初始含 Bug 源码。要复现原始缺陷，必须检出下面固定的 parent SHA；不要在当前修复结果源码上期待重新出现修复前失败。生成系统使用的可信验证补丁和完整验证日志仅在本地留存，不提交到结果分支。

## 问题现象

用毫秒或微秒 epoch 数字作为时间戳摄取日志时，事件时间被截断到整秒：1735689600123 变成 2025-01-01T00:00:00Z，微秒输入同样丢掉亚秒部分，导致同一秒内的事件排序和窗口统计都不对。请修复 epoch 时间戳解析，让毫秒和微秒输入保留亚秒精度，同时保持秒级与带小数的 epoch 输入行为不变，并保证全量测试通过。

## 含 Bug 版本

- 仓库：zhanglei10281852-gif/gogo-47
- 仓库地址：https://github.com/zhanglei10281852-gif/gogo-47.git
- parent SHA：eb5d79c84ca0b12fe24a7084cf65d78e409ae933

## 复现步骤

```bash
git clone -- https://github.com/zhanglei10281852-gif/gogo-47.git bug-repro
cd bug-repro
git checkout --detach eb5d79c84ca0b12fe24a7084cf65d78e409ae933
go test ./internal/ingest -run "^TestParseEpochPreservesSubsecondUnits$" -count=1 -v
```

## 双架构完整错误信息

### linux/amd64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/ingest -run "^TestParseEpochPreservesSubsecondUnits$" -count=1 -v
=== RUN   TestParseEpochPreservesSubsecondUnits
=== RUN   TestParseEpochPreservesSubsecondUnits/milliseconds
    epoch_regression_test.go:25: parseTimestamp(1735689600123) = 2025-01-01T00:00:00Z, want 2025-01-01T00:00:00.123Z
=== RUN   TestParseEpochPreservesSubsecondUnits/microseconds
    epoch_regression_test.go:25: parseTimestamp(1735689600123456) = 2025-01-01T00:00:00Z, want 2025-01-01T00:00:00.123456Z
--- FAIL: TestParseEpochPreservesSubsecondUnits (0.00s)
    --- FAIL: TestParseEpochPreservesSubsecondUnits/milliseconds (0.00s)
    --- FAIL: TestParseEpochPreservesSubsecondUnits/microseconds (0.00s)
FAIL
FAIL	LogPilot/internal/ingest	0.012s
FAIL

```

stderr：

```text
(empty)
```

### linux/arm64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/ingest -run "^TestParseEpochPreservesSubsecondUnits$" -count=1 -v
=== RUN   TestParseEpochPreservesSubsecondUnits
=== RUN   TestParseEpochPreservesSubsecondUnits/milliseconds
    epoch_regression_test.go:25: parseTimestamp(1735689600123) = 2025-01-01T00:00:00Z, want 2025-01-01T00:00:00.123Z
=== RUN   TestParseEpochPreservesSubsecondUnits/microseconds
    epoch_regression_test.go:25: parseTimestamp(1735689600123456) = 2025-01-01T00:00:00Z, want 2025-01-01T00:00:00.123456Z
--- FAIL: TestParseEpochPreservesSubsecondUnits (0.01s)
    --- FAIL: TestParseEpochPreservesSubsecondUnits/milliseconds (0.01s)
    --- FAIL: TestParseEpochPreservesSubsecondUnits/microseconds (0.00s)
FAIL
FAIL	LogPilot/internal/ingest	0.135s
FAIL

```

stderr：

```text
(empty)
```

## 通过条件

毫秒输入解析为 .123，微秒输入解析为 .123456；秒级与小数 epoch、越界校验不回归；双架构定向、全量、build/vet 通过。
