# pgxatomic

A library for managing PostgreSQL transactions through context propagation using the [pgx](https://github.com/jackc/pgx) driver.

pgxatomic allows you to pass transactions via `context.Context`, providing cleaner separation between business logic and transaction management. Repository and service layers work with contexts while transaction boundaries are controlled at higher levels.

![schema](https://i.imgur.com/RpsfuBb.jpg)

## Installation

```bash
go get github.com/ysomad/pgxatomic
```

## Usage

### Repository layer

Use `pgxatomic.Pool` in your repository. It wraps `pgxpool.Pool` and automatically uses transactions from context when available.

```go
type orderRepo struct {
    pool *pgxatomic.Pool
}

type order struct {
    ID   uuid.UUID
    Cost int
}

func (r *orderRepo) Insert(ctx context.Context, cost int) order {
    rows, _ := r.pool.Query(ctx, "INSERT INTO orders(cost) VALUES ($1) RETURNING id, cost", cost)
    o, _ := pgx.CollectOneRow(rows, pgx.RowToStructByPos[order])
    return o
}
```

Alternatively, use `Query`, `QueryRow`, and `Exec` functions directly without the pool wrapper.

### Transaction management

Use `Runner` to execute operations within a transaction. The transaction propagates automatically through context.

```go
conf, _ := pgxpool.ParseConfig("postgres://user:pass@localhost:5432/postgres")
pool, _ := pgxpool.NewWithConfig(context.Background(), conf)

runner, _ := pgxatomic.NewRunner(pool, pgx.TxOptions{})

_ = runner.Run(context.Background(), func(txCtx context.Context) error {
    _ = orderService.Create(txCtx)
    _ = balanceService.Withdraw(txCtx)
    return nil
})
```

All operations within `Run` use the same transaction. If the function returns an error, the transaction rolls back. Otherwise, it commits.

Note: Error handling is omitted for brevity. Handle errors appropriately in production code.

## References

- [Clean transactions in Golang hexagon](https://www.kaznacheev.me/posts/en/clean-transactions-in-hexagon)
