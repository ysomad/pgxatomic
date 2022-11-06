# pgxatomic

pgxatomic is a library of tools that allow you to implement transfer of clean control to transactions to a higher level by adding transaction in a `context.Context` using [pgx](https://github.com/jackc/pgx) driver.

![schema](https://i.imgur.com/RpsfuBb.jpg)

## Example Usage
1. You can use `pgxatomic.Pool` within repository implementation. It's simple wrapper around `pgxpool.Pool` which
is wrapping `Query`, `QueryRow` and `Exec` methods with `pgxatomic` query functions.
```go
type orderRepo struct {
    pool *pgxatomic.Pool // pgxpool.Pool wrapper
}

type order struct {
    ID uuid.UUID
    Cost int
}

func (r *orderRepo) Insert(ctx context.Context, cost int) order {
    rows, _ := r.pool.Query(ctx, "insert into order(cost) values ($1) RETURNING id, cost", cost)
    o, _ := pgx.CollectOneRow(rows, pgx.RowToStructByPos[order])
    return o
}
```

Or you can use `Query`, `QueryRow`, `Exec` functions directly from the library.

2. Run wrapped usecase method calls within txFunc using `pgxatomic.runner.Run` function
```go
conf, _ := pgxpool.ParseConfig("postgres://user:pass@localhost:5432/postgres")
pool, _ := pgxpool.NewWithConfig(context.Background(), conf)

r, _ := pgxatomic.NewRunner(pool, pgx.TxOptions{})

_ = r.Run(context.Background(), func(txCtx context.Context) error {
    _ = orderService.Create(txCtx)
    _ = balanceService.Withdraw(txCtx)
    return nil
})
```

Error handling is omitted on purpose, handle all errors!

## Credits
- [Clean transactions in Golang hexagon](https://www.kaznacheev.me/posts/en/clean-transactions-in-hexagon)
