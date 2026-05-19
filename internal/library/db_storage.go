package library

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const BOOK_COLUMNS = "id, title, author, number_of_pages, is_read, created_at, read_at"

type DatabaseStorage struct {
	pool *pgxpool.Pool
}

func NewDatabaseStorage(pool *pgxpool.Pool) *DatabaseStorage {
	return &DatabaseStorage{
		pool: pool,
	}
}

func scanBook(row pgx.Row) (Book, error) {
	var b Book
	var createdAt time.Time
	var readAt *time.Time

	err := row.Scan(
		&b.ID,
		&b.Title,
		&b.Author,
		&b.NumberOfPages,
		&b.IsRead,
		&createdAt,
		&readAt,
	)
	if err != nil {
		return Book{}, err
	}

	b.CreatedAt = PrettyTime(createdAt)
	if readAt != nil {
		prettyReadAt := PrettyTime(*readAt)
		b.ReadAt = &prettyReadAt
	}

	return b, nil
}

func convertError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrBookNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrBookAlreadyExists
	}
	return err
}

func (s *DatabaseStorage) AddBook(ctx context.Context, params BookParams) (Book, error) {
	query := fmt.Sprintf(`
	INSERT INTO books (title, author, number_of_pages)
	VALUES ($1, $2, $3)
	RETURNING %s`, BOOK_COLUMNS)

	row := s.pool.QueryRow(ctx, query, params.Title, params.Author, params.Pages)

	b, err := scanBook(row)

	if err != nil {
		return Book{}, convertError(fmt.Errorf("failed to add book: %w", err))
	}

	return b, nil
}

func (s *DatabaseStorage) GetAllBooks(ctx context.Context, p GetBooksParams) ([]Book, error) {
	query := fmt.Sprintf(`
	SELECT %s FROM books
	WHERE 1=1`, BOOK_COLUMNS)

	var args []any
	argID := 1

	if p.Author != "" {
		args = append(args, "%"+strings.TrimSpace(p.Author)+"%")
		query += fmt.Sprintf(" AND author ILIKE $%d", argID)
		argID++
	}
	if p.IsRead != nil {
		args = append(args, *p.IsRead)
		query += fmt.Sprintf(" AND is_read = $%d", argID)
		argID++
	}

	query += " ORDER BY "
	switch p.SortType {
	case "author":
		query += "author ASC"
	case "title":
		query += "title ASC"
	case "pages":
		query += "number_of_pages ASC"
	default:
		query += "created_at DESC"
	}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	res := make([]Book, 0)

	for rows.Next() {
		b, err := scanBook(rows)
		if err != nil {
			return nil, convertError(fmt.Errorf("failed to get all books: %w", err))
		}

		res = append(res, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *DatabaseStorage) GetBook(ctx context.Context, id uuid.UUID) (Book, error) {
	query := fmt.Sprintf(`
	SELECT %s FROM books
	WHERE id = $1`, BOOK_COLUMNS)

	row := s.pool.QueryRow(ctx, query, id)

	b, err := scanBook(row)

	if err != nil {
		return Book{}, convertError(fmt.Errorf("failed to get book by id: %w", err))
	}

	return b, nil
}

func (s *DatabaseStorage) DeleteBook(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE FROM books
	WHERE id = $1`

	tag, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrBookNotFound
	}
	return nil
}

func (s *DatabaseStorage) MarkBook(ctx context.Context, id uuid.UUID, read bool) (Book, error) {
	query := fmt.Sprintf(`
	UPDATE books
	SET is_read = $2,
	read_at =
	CASE 
		WHEN $2 = TRUE AND is_read = FALSE THEN NOW()
		WHEN $2 = FALSE THEN NULL
		ELSE read_at
	END
	WHERE id = $1
	RETURNING %s`, BOOK_COLUMNS)

	b, err := scanBook(s.pool.QueryRow(ctx, query, id, read))

	if err != nil {
		return Book{}, convertError(fmt.Errorf("failed to mark book read: %w", err))
	}

	return b, nil
}

func (s *DatabaseStorage) PatchBook(ctx context.Context, id uuid.UUID, p UpdateBookParams) (Book, error) {
	clauses := []string{}
	args := []any{}
	argID := 1

	addCol := func(column string, value any) {
		clauses = append(clauses, fmt.Sprintf("%s=$%d", column, argID))
		args = append(args, value)
		argID++
	}

	if p.Title != nil {
		addCol("title", *p.Title)
	}
	if p.Author != nil {
		addCol("author", *p.Author)
	}
	if p.Pages != nil {
		addCol("number_of_pages", *p.Pages)
	}
	if p.IsRead != nil {
		addCol("is_read", *p.IsRead)
		readLogic := fmt.Sprintf(`read_at = 
		CASE
			WHEN $%d=TRUE AND is_read = FALSE THEN NOW()
			WHEN $%d=FALSE THEN NULL
			ELSE read_at
		END`, argID-1, argID-1)

		clauses = append(clauses, readLogic)
	}

	if len(clauses) == 0 {
		return s.GetBook(ctx, id)
	}

	args = append(args, id)
	query := fmt.Sprintf(`UPDATE books SET %s WHERE id=$%d RETURNING %s`,
		strings.Join(clauses, ", "), argID, BOOK_COLUMNS)

	row := s.pool.QueryRow(ctx, query, args...)

	b, err := scanBook(row)

	if err != nil {
		return Book{}, convertError(fmt.Errorf("failed to patch a book: %w", err))
	}

	return b, nil
}
