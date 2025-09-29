###Gator - RSS aggregator

- go 1.15
- postgres@17

Installation 
```
go install github.com/sidis405/gator
```

```yaml
version: "2"
sql:
  - schema: "sql/schema"
    queries: "sql/queries"
    engine: "postgresql"
    gen:
      go:
        out: "internal/database"
```

.gatorconfig.json
```json
{
  "db_url": "postgres://root:@localhost:5432/gator?sslmode=disable"
}
```
