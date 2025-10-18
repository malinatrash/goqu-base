# Goqu Base Repository

A generic base repository implementation using [goqu](https://github.com/doug-martin/goqu) SQL builder and PostgreSQL with pgx driver.

## Features

- Generic repository pattern with type safety
- Built-in error handling with predefined error constants
- Support for CRUD operations with flexible filtering
- Pagination support
- Prepared statements for better performance
- Consistent error handling across all methods

## Installation

```bash
go get github.com/malinatrash/goqu-base
```

## Dependencies

- [goqu](https://github.com/doug-martin/goqu) - SQL builder
- [pgx](https://github.com/jackc/pgx) - PostgreSQL driver
- [scany](https://github.com/georgysavva/scany) - SQL result scanning

## Usage

### Define Your Model

```go
type User struct {
    ID       int64     `db:"id" goqu:"skipinsert"`
    Name     string    `db:"name"`
    Email    string    `db:"email"`
    Age      int       `db:"age"`
    CreatedAt time.Time `db:"created_at" goqu:"skipinsert"`
}
```

### Initialize Repository

```go
package main

import (
    "context"
    "log"

    "github.com/jackc/pgx/v5/pgxpool"
    gb "github.com/malinatrash/goqu-base"
)

func main() {
    // Initialize database connection
    pool, err := pgxpool.New(context.Background(), "postgres://user:password@localhost/dbname")
    if err != nil {
        log.Fatal(err)
    }
    defer pool.Close()

    // Create repository instance
    userRepo := gb.NewBaseRepo[User](pool, "users")
}
```

## API Reference

### Create

Create a new record in the database.

```go
func (r *BaseRepo[T]) Create(ctx context.Context, model *T) error
```

**Example:**

```go
user := &User{
    Name:  "John Doe",
    Email: "john@example.com",
    Age:   30,
}

err := userRepo.Create(ctx, user)
if err != nil {
    log.Fatal(err)
}
// user.ID and user.CreatedAt will be populated after creation
```

### GetByID

Retrieve a single record by its ID.

```go
func (r *BaseRepo[T]) GetByID(ctx context.Context, id interface{}) (T, error)
```

**Example:**

```go
user, err := userRepo.GetByID(ctx, 1)
if err != nil {
    if errors.Is(err, gb.ErrNotFound) {
        log.Println("User not found")
        return
    }
    log.Fatal(err)
}
fmt.Printf("Found user: %+v\n", user)
```

### GetBy

Retrieve a single record by custom filters.

```go
func (r *BaseRepo[T]) GetBy(ctx context.Context, filters map[string]any) (T, error)
```

**Example:**

```go
// Find user by email
user, err := userRepo.GetBy(ctx, map[string]any{
    "email": "john@example.com",
})
if err != nil {
    if errors.Is(err, gb.ErrNotFound) {
        log.Println("User not found")
        return
    }
    log.Fatal(err)
}

// Find user by multiple conditions
user, err = userRepo.GetBy(ctx, map[string]any{
    "name": "John Doe",
    "age":  30,
})
```

### FindMany

Retrieve multiple records with pagination and filtering.

```go
func (r *BaseRepo[T]) FindMany(ctx context.Context, page, limit int32, filters map[string]any) ([]T, int32, error)
```

**Example:**

```go
// Get all users (first page, 10 per page)
users, totalPages, err := userRepo.FindMany(ctx, 1, 10)

// Get users with filters
users, totalPages, err = userRepo.FindMany(ctx, 1, 10,
    map[string]any{"age": 30},
)

// Get users older than 25
users, totalPages, err = userRepo.FindMany(ctx, 1, 10,
    map[string]any{"age": goqu.Op{"gt": 25}},
)

fmt.Printf("Found %d users, total pages: %d\n", len(users), totalPages)
```

### FindMany

Расширенный поиск с дополнительными опциями для более гибкой работы с данными.

```go
func (r *BaseRepo[T]) FindMany(ctx context.Context, page, limit int32, opts ...func(*FindManyOptions)) ([]T, int32, error)
```

#### Доступные опции

##### Фильтры (Filters)

```go
// Фильтрация по полям
WithFilters(map[string]any{"status": "active"})

// Фильтрация с использованием goqu операторов
WithFilters(map[string]any{
    "age": goqu.Op{"gte": 18},
    "name": goqu.Op{"like": "John%"},
})
```

##### Сортировка (OrderBy)

```go
// Сортировка по одному полю по возрастанию
WithOrderBy("created_at", false) // false = ASC, true = DESC

// Сортировка по нескольким полям
WithOrderBy("created_at", false),
WithOrderBy("name", true),
```

##### Выбор полей (Select)

```go
// Выбрать только нужные поля для оптимизации
WithSelect("id", "name", "email", "created_at")
```

##### Дополнительные условия (Conditions)

```go
// Добавить кастомные WHERE условия
WithCondition(goqu.C("status").Eq("active")),
WithCondition(goqu.C("age").Gte(18)),
```

##### Группировка (GroupBy)

```go
// Группировка по полям
WithGroupBy("status", "category")
```

##### Условия группировки (Having)

```go
// Условия для группировки
WithHaving(goqu.COUNT("*").Gt(5))
```

##### Текстовый поиск (Search)

```go
// Поиск по нескольким полям с частичным совпадением
WithSearch("john", "name", "email", "description")
```

##### Фильтр по датам (DateRange)

```go
from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
to := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

// Фильтр по диапазону дат
WithDateRange("created_at", &from, &to)
```

##### JOIN операции

```go
// LEFT JOIN с другой таблицей
WithJoin("user_profiles", "p",
    goqu.C("users.profile_id").Eq(goqu.C("p.id")),
    JoinLeft)

// INNER JOIN
WithJoin("orders", "o",
    goqu.C("users.id").Eq(goqu.C("o.user_id")),
    JoinInner)
```

##### Блокировка для конкурентного доступа

```go
// Пессимистическая блокировка
WithLock(LockForUpdate)

// Блокировка без блокировки ключей
WithLock(LockForNoKeyUpdate)
```

##### Специальные режимы

```go
// Отключить пагинацию (получить все результаты)
WithNoPagination()

// Получить только количество записей
WithCountOnly()

// Получить уникальные записи
WithDistinct()

// Максимальное количество результатов
WithMaxResults(1000)
```

#### Примеры использования

```go
// Поиск с сортировкой и выбором полей
users, totalPages, err := repo.FindMany(ctx, 1, 10,
    WithFilters(map[string]any{"status": "active"}),
    WithOrderBy("created_at", false),
    WithSelect("id", "name", "email"),
)

// Поиск с фильтром по датам
users, totalPages, err = repo.FindMany(ctx, 1, 10,
    WithFilters(map[string]any{"status": "active"}),
    WithDateRange("created_at", &from, &to),
)

// Поиск с JOIN'ом и фильтром по датам
users, totalPages, err = repo.FindMany(ctx, 1, 10,
    WithFilters(map[string]any{"type": "premium"}),
    WithJoin("user_profiles", "p", 
        goqu.C("users.profile_id").Eq(goqu.C("p.id")), 
        JoinLeft),
    WithDateRange("users.created_at", &from, &to),
    WithSelect("users.*", "p.first_name", "p.last_name"),
)

// Поиск с группировкой и агрегацией
users, totalPages, err = repo.FindMany(ctx, 1, 10,
    WithFilters(map[string]any{"status": "active"}),
    WithGroupBy("status"),
    WithHaving(goqu.COUNT("*").Gt(5)),
)

// Текстовый поиск с отключением пагинации
users, totalPages, err = repo.FindMany(ctx, 1, 10,
    WithFilters(map[string]any{"type": "premium"}),
    WithSearch("enterprise", "name", "description"),
    WithNoPagination(),
    WithMaxResults(100),
)

// Только подсчет количества записей
_, totalPages, err = repo.FindMany(ctx, 1, 10,
    WithFilters(map[string]any{"status": "active"}),
    WithCountOnly(),
)
```

Update a record by its model (uses the ID field).

```go
func (r *BaseRepo[T]) Update(ctx context.Context, model *T) error
```

**Example:**

```go
user.Name = "Jane Doe"
user.Age = 31

err := userRepo.Update(ctx, &user)
if err != nil {
    if errors.Is(err, gb.ErrNotFound) {
        log.Println("User not found")
        return
    }
    log.Fatal(err)
}
```

### UpdateBy

Update records by custom filters.

```go
func (r *BaseRepo[T]) UpdateBy(ctx context.Context, filters map[string]any, updates map[string]any) error
```

**Example:**

```go
// Update all users with age 30 to age 31
err := userRepo.UpdateBy(ctx,
    map[string]any{"age": 30},           // filters
    map[string]any{"age": 31},           // updates
)

// Update user by email
err = userRepo.UpdateBy(ctx,
    map[string]any{"email": "john@example.com"},
    map[string]any{
        "name": "John Smith",
        "age":  32,
    },
)

if err != nil {
    if errors.Is(err, gb.ErrNotFound) {
        log.Println("No records were updated")
        return
    }
    log.Fatal(err)
}
```

### Delete

Delete a record by its ID.

```go
func (r *BaseRepo[T]) Delete(ctx context.Context, id interface{}) error
```

**Example:**

```go
err := userRepo.Delete(ctx, 1)
if err != nil {
    if errors.Is(err, gb.ErrNotFound) {
        log.Println("User not found")
        return
    }
    log.Fatal(err)
}
```

### DeleteBy

Delete records by custom filters.

```go
func (r *BaseRepo[T]) DeleteBy(ctx context.Context, filters map[string]any) error
```

**Example:**

```go
// Delete all users older than 65
err := userRepo.DeleteBy(ctx, map[string]any{
    "age": goqu.Op{"gt": 65},
})

// Delete user by email
err = userRepo.DeleteBy(ctx, map[string]any{
    "email": "john@example.com",
})

if err != nil {
    if errors.Is(err, gb.ErrNotFound) {
        log.Println("No records were deleted")
        return
    }
    log.Fatal(err)
}
```

## Error Handling

The repository provides predefined error constants for consistent error handling:

```go
var (
    ErrNotFound           = fmt.Errorf("not found")
    ErrInvalidPagination  = fmt.Errorf("invalid pagination parameters")
    ErrBuildSQL           = fmt.Errorf("error building SQL")
    ErrQueryFailed        = fmt.Errorf("query failed")
    ErrInsertFailed       = fmt.Errorf("insert failed")
    ErrUpdateFailed       = fmt.Errorf("update failed")
    ErrDeleteFailed       = fmt.Errorf("delete failed")
    ErrCountFailed        = fmt.Errorf("count query failed")
)
```

**Example error handling:**

```go
user, err := userRepo.GetByID(ctx, 1)
if err != nil {
    switch {
    case errors.Is(err, gb.ErrNotFound):
        // Handle not found case
        return nil, fmt.Errorf("user with ID 1 not found")
    case errors.Is(err, gb.ErrQueryFailed):
        // Handle query error
        return nil, fmt.Errorf("database query failed: %w", err)
    default:
        // Handle other errors
        return nil, fmt.Errorf("unexpected error: %w", err)
    }
}
```

## Advanced Filtering

You can use goqu operators for complex filtering:

```go
import "github.com/doug-martin/goqu/v9"

// Users with age between 25 and 35
users, _, err := userRepo.FindMany(ctx, 1, 10, map[string]any{
    "age": goqu.Op{"between": goqu.Range(25, 35)},
})

// Users with names starting with "John"
users, _, err = userRepo.FindMany(ctx, 1, 10, map[string]any{
    "name": goqu.Op{"like": "John%"},
})

// Users with specific IDs
users, _, err = userRepo.FindMany(ctx, 1, 10, map[string]any{
    "id": goqu.Op{"in": []int{1, 2, 3, 4, 5}},
})

// Users created in the last 30 days
thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
users, _, err = userRepo.FindMany(ctx, 1, 10, map[string]any{
    "created_at": goqu.Op{"gte": thirtyDaysAgo},
})
```

## Complete Example

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "log"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
    gb "github.com/your-username/goqu-base"
)

type User struct {
    ID        int64     `db:"id" goqu:"skipinsert"`
    Name      string    `db:"name"`
    Email     string    `db:"email"`
    Age       int       `db:"age"`
    CreatedAt time.Time `db:"created_at" goqu:"skipinsert"`
}

func main() {
    ctx := context.Background()

    // Initialize database connection
    pool, err := pgxpool.New(ctx, "postgres://user:password@localhost/dbname")
    if err != nil {
        log.Fatal(err)
    }
    defer pool.Close()

    // Create repository
    userRepo := gb.NewBaseRepo[User](pool, "users")

    // Create a new user
    user := &User{
        Name:  "John Doe",
        Email: "john@example.com",
        Age:   30,
    }

    if err := userRepo.Create(ctx, user); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created user with ID: %d\n", user.ID)

    // Get user by ID
    foundUser, err := userRepo.GetByID(ctx, user.ID)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found user: %+v\n", foundUser)

    // Update user
    foundUser.Age = 31
    if err := userRepo.Update(ctx, &foundUser); err != nil {
        log.Fatal(err)
    }

    // Find users with pagination
    users, totalPages, err := userRepo.FindMany(ctx, 1, 10, map[string]any{
        "age": 31,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d users, total pages: %d\n", len(users), totalPages)

    // Delete user
    if err := userRepo.Delete(ctx, user.ID); err != nil {
        if errors.Is(err, gb.ErrNotFound) {
            fmt.Println("User already deleted")
        } else {
            log.Fatal(err)
        }
    }
}
```

## License

MIT License
