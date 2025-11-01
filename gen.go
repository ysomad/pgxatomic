package pgxatomic

//go:generate mockgen -package pgxatomic -destination mocks_test.go github.com/jackc/pgx/v5 Rows,Row,Tx
//go:generate mockgen -package pgxatomic -source runner.go -destination runner_mocks_test.go txStarter
