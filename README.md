# pgxatomic

pgxatomic is a set of functions that can be used to wrap repository method calls by adding transaction in `context.Context` on the higher level using [pgx](https://github.com/jackc/pgx) driver.

![schema](https://i.imgur.com/RpsfuBb.jpg)

## Problem
It's common practice to use repository pattern these days so problem of atomic calls of different repositoty methods arises.
For example there it is TWO entities `Order` and `UserBalance` and they has separate repositories. You want to create an order and withdraw amount of money from the user's account balance, of course it has to be atomic, there it is first solutions that come to mind:
- run single query in a transaction or CTE within repository method of Order or UserBalance - impairs code readability by hiding data interaction with two different entities from other layers of application into one of repositories
- create separate repository `OrderUserBalance` and run query in a transaction or CTE within - increases amount of code to write and quickly turn into noodles from repositories with 1-2 methods and long ugly names

And there it is also a third solution which is considered in this repository:
- ***`Order` and `UserBalance` has their own repositories with simple CRUD queries shares the same context with transaction or without, depends on the caller - cleanest implementation which allows to not worry about transaction in business logic or repository layers***

## Usage
1. Repository method has to call wrapped query functions from the package. For example `atomic.Query`
```go
type repo struct {
    pool *pgxpool.Pool
}

type user struct {
    ID uuid.UUID
    Name string
}

func (r *repo) getUserByID(ctx context.Context, id uuid.UUID) user {
    rows, _ := atomic.Query(ctx, r.pool, "select * from user where id = $1", id)
    u, _ := pgx.CollectOneRow(rows, pgx.RowToStructByPos[user])
    return u
}
```

2. Wrap usecase method calls within txFunc using `atomic.Run` function
```go
_ = atomic.Run(context.Background(), pool, func(txCtx context.Context) error {
    _ = orderService.Create(txCtx)
    _ = balanceService.Withdraw(txCtx)
    return nil
})
```

Or its possible to use `pgxatomic.runner`:
```go
conf, _ := pgxpool.ParseConfig("postgres://user:pass@localhost:5432/postgres")
pool, _ := pgxpool.NewWithConfig(context.Background(), conf)

r := atomic.NewRunner(pool)

_ = r.Run(context.Background(), func(txCtx context.Context) error {
    _ = orderService.Create(txCtx)
    _ = balanceService.Withdraw(txCtx)
    return nil
})
```

Error handling is omitted on purpose, handle all errors!

## TODO
1. Add examples
2. Add clean RunWithOpts function to run BeginTx
3. Write tests
4. Write code-generator for DB implementation
5. Write code-generator for Runner implementation

## Credits
- [Clean transactions in Golang hexagon](https://www.kaznacheev.me/posts/en/clean-transactions-in-hexagon)
