# pgxatomic

pgxatomic is a library of tools that allow you to implement transfer of clean control to transactions to a higher level by adding transaction in a `context.Context` using [pgx](https://github.com/jackc/pgx) driver.

![schema](https://i.imgur.com/RpsfuBb.jpg)

## Problem
It's common practice to use repository pattern these days so problem of atomic calls of different repositoty methods arises.

For example there it is TWO entities `Order` and `UserBalance` and they has separate repositories. You want to create an order and withdraw amount of money from the user's account balance, of course it has to be atomic, there it is first solutions that come to mind:
- run single query in a transaction or CTE within repository method of Order or UserBalance - impairs code readability by hiding data interaction with two different entities from other layers 
- create separate repository `OrderUserBalance` and run query in a transaction or CTE within - increases amount of code to write and quickly turn into noodles from repositories with 1-2 methods and long ugly names

And there it is also a third solution which is considered in this library:
- ***`Order` and `UserBalance` has their own repositories with simple CRUD queries shares the same context with transaction or without, depends on the caller - cleanest implementation which allows to not worry about transaction in business logic or repository layers***

## Example Usage
1. Repository method has to call wrapped query functions from the package. For example `pgxatomic.Query`
```go
type orderRepo struct {
    pool *pgxpool.Pool
}

type order struct {
    ID uuid.UUID
    Cost int
}

func (r *orderRepo) query(ctx, sql string, args ...any) (pgx.Rows, error) {
    return pgxatomic.Query(ctx, r.pool. sql, args...)
}

func (r *orderRepo) CreateOrder(ctx context.Context, cost int) order {
    rows, _ := r.query(ctx, "insert into order(cost) values ($1)", cost)
    o, _ := pgx.CollectOneRow(rows, pgx.RowToStructByPos[order])
    return o
}
```

2. Wrap usecase method calls within txFunc using `pgxatomic.runner` function
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

## TODO
1. Write tests
2. Write code-generator for repository implementation (generate repo with query methods wrapping pgxatomic.Query etc)

## Credits
- [Clean transactions in Golang hexagon](https://www.kaznacheev.me/posts/en/clean-transactions-in-hexagon)
