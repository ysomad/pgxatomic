# pgxatomic

pgxatomic is a library of tools that allow you to implement transfer of clean control to transactions to a higher level by adding transaction in a `context.Context` using [pgx](https://github.com/jackc/pgx) driver.

![schema](https://i.imgur.com/RpsfuBb.jpg)

## Problem
It's common practice to use repository pattern these days so problem of atomic calls of different repositoty methods arises.

For example there is TWO entities `Order` and `UserBalance` and they has separate repositories. You want to create an order and withdraw amount of money from the user's account balance, of course it has to be atomic, there is first solutions that come to mind:
- run single query in a transaction or CTE within repository method of `Order` or `UserBalance` - impairs code readability by hiding data interaction with two different entities from other layers
- create separate repository `OrderUserBalance` and run query in a transaction or CTE within - increases amount of code to write and quickly turn into noodles from repositories with 1-2 methods and long ugly names

And there is also a third solution which is considered in pgxatomic:
- ***`Order` and `UserBalance` has their own repositories with simple CRUD queries shares the same context with transaction or without, depends on the caller - cleanest implementation which allows to not worry about transaction in business logic or repository layers***

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
