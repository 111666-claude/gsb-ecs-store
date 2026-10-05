# ecs-store

按组件索引查询实体：先挂组件，再反复查询。只用 Go 标准库。

```
go test ./...
go run ./cmd/query --sample=order
go run ./cmd/query --sample=disabled
go run ./cmd/query --sample=revive
go run ./cmd/query --sample=work
```

## 口径

- **查询**：返回同时拥有全部给定组件的实体；组件列表为空时返回全部实体。
- **顺序**：结果按实体 ID 升序排列，与挂组件的先后无关。
- **禁用覆盖层**：被 `Disable` 的实体不出现在任何查询结果里；禁用状态只能由显式的启用操作解除，后续挂组件不解除。
- **幂等**：同一实体重复挂同一个组件不产生第二条记录。
- **索引**：查询必须由组件倒排表驱动，从最小的那个集合开始筛。

## 不变量

- 查询结果里同一个实体只出现一次。
- 结果集合与挂组件的顺序无关。
- `Work` 记查询期间检查过的候选实体数，不得随实体总数放大。

## 输出契约

每个查询打印一行：

```
{"query": [...], "hits": [...]}
```

`--sample=work` 只打印 `{"hits": N, "queries": N, "work": N}`。
