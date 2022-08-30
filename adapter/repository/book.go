package repository

import (
	"context"

	"github.com/XBozorg/bookstore/entity/book"
	"github.com/XBozorg/bookstore/entity/order"
)

func (storage Storage) DoesAuthorExist(ctx context.Context, authorID uint) (bool, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM author WHERE id = ?)",
	)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, authorID)

	var doesExist bool
	if err = result.Scan(&doesExist); err != nil {
		return false, err
	}

	return doesExist, nil
}

func (storage Storage) DoesPublisherExist(ctx context.Context, publisherID uint) (bool, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM publisher WHERE id = ?)",
	)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, publisherID)

	var doesExist bool
	if err = result.Scan(&doesExist); err != nil {
		return false, err
	}

	return doesExist, nil
}

func (storage Storage) DoesTopicExist(ctx context.Context, topicID uint) (bool, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM topic WHERE id = ?)",
	)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, topicID)

	var doesExist bool
	if err = result.Scan(&doesExist); err != nil {
		return false, err
	}

	return doesExist, nil
}

func (storage Storage) DoesLanguageExist(ctx context.Context, langID uint) (bool, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM language WHERE id = ?)",
	)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, langID)

	var doesExist bool
	if err = result.Scan(&doesExist); err != nil {
		return false, err
	}

	return doesExist, nil
}

func (storage Storage) DoesBookExist(ctx context.Context, bookID uint) (bool, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM book WHERE id = ?)",
	)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, bookID)

	var doesExist bool
	if err = result.Scan(&doesExist); err != nil {
		return false, err
	}

	return doesExist, nil
}

func (storage Storage) AddAuthor(ctx context.Context, authorName string) (book.Author, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"INSERT INTO author (name) VALUES (?)",
	)
	if err != nil {
		return book.Author{}, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, authorName)
	if err != nil {
		return book.Author{}, err
	}

	authorID, err := result.LastInsertId()
	if err != nil {
		return book.Author{}, err
	}

	return book.Author{
		ID:   uint(authorID),
		Name: authorName,
	}, nil
}

func (storage Storage) GetAuthor(ctx context.Context, authorID uint) (book.Author, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT name From author WHERE id = ?",
	)
	if err != nil {
		return book.Author{}, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, authorID)

	author := book.Author{}
	if err = result.Scan(&author.Name); err != nil {
		return book.Author{}, err
	}

	author.ID = authorID
	return author, nil
}

func (storage Storage) GetAuthors(ctx context.Context) ([]book.Author, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT id , name FROM author",
	)
	if err != nil {
		return []book.Author{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx)
	if err != nil {
		return []book.Author{}, err
	}
	defer result.Close()

	authors := []book.Author{}
	for result.Next() {
		var a book.Author

		if err = result.Scan(
			&a.ID,
			&a.Name,
		); err != nil {
			return []book.Author{}, err
		}
		authors = append(authors, a)
	}

	return authors, nil
}

func (storage Storage) DeleteAuthor(ctx context.Context, authorID uint) error {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"DELETE FROM author WHERE id = ?",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, authorID); err != nil {
		return err
	}

	return nil
}

func (storage Storage) AddPublisher(ctx context.Context, publisherName string) (book.Publisher, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"INSERT INTO publisher (name) VALUES (?)",
	)
	if err != nil {
		return book.Publisher{}, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, publisherName)
	if err != nil {
		return book.Publisher{}, err
	}

	publisherID, err := result.LastInsertId()
	if err != nil {
		return book.Publisher{}, err
	}

	return book.Publisher{
		ID:   uint(publisherID),
		Name: publisherName,
	}, nil
}

func (storage Storage) GetPublisher(ctx context.Context, publisherID uint) (book.Publisher, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT name FROM publisher WHERE id = ?",
	)
	if err != nil {
		return book.Publisher{}, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, publisherID)

	publisher := book.Publisher{}

	if err = result.Scan(&publisher.Name); err != nil {
		return book.Publisher{}, err
	}

	publisher.ID = publisherID
	return publisher, nil
}

func (storage Storage) GetPublishers(ctx context.Context) ([]book.Publisher, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT id , name FROM publisher",
	)
	if err != nil {
		return []book.Publisher{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx)
	if err != nil {
		return []book.Publisher{}, err
	}
	defer result.Close()

	publishers := []book.Publisher{}
	for result.Next() {
		var p book.Publisher

		if err = result.Scan(
			&p.ID,
			&p.Name,
		); err != nil {
			return []book.Publisher{}, err
		}
		publishers = append(publishers, p)
	}

	return publishers, nil
}

func (storage Storage) DeletePublisher(ctx context.Context, publisherId uint) error {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"DELETE FROM publisher WHERE id = ?",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, publisherId); err != nil {
		return err
	}

	return nil
}

func (storage Storage) AddTopic(ctx context.Context, topicName string) (book.Topic, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"INSERT INTO topic (name) VALUES (?)",
	)
	if err != nil {
		return book.Topic{}, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, topicName)
	if err != nil {
		return book.Topic{}, err
	}

	topicID, err := result.LastInsertId()
	if err != nil {
		return book.Topic{}, err
	}

	return book.Topic{
		ID:   uint(topicID),
		Name: topicName,
	}, nil
}

func (storage Storage) GetTopic(ctx context.Context, topicID uint) (book.Topic, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT name FROM topic WHERE id = ?",
	)
	if err != nil {
		return book.Topic{}, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, topicID)

	topic := book.Topic{}
	if err = result.Scan(&topic.Name); err != nil {
		return book.Topic{}, err
	}

	topic.ID = topicID
	return topic, nil
}

func (storage Storage) GetTopics(ctx context.Context) ([]book.Topic, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT id , name FROM topic",
	)
	if err != nil {
		return []book.Topic{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx)
	if err != nil {
		return []book.Topic{}, err
	}
	defer result.Close()

	topics := []book.Topic{}
	for result.Next() {
		var t book.Topic

		if err = result.Scan(
			&t.ID,
			&t.Name,
		); err != nil {
			return []book.Topic{}, err
		}
		topics = append(topics, t)
	}

	return topics, nil
}

func (storage Storage) DeleteTopic(ctx context.Context, topicID uint) error {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"DELETE FROM topic WHERE id = ?",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, topicID); err != nil {
		return err
	}

	return nil
}

func (storage Storage) AddLanguage(ctx context.Context, langCode string) (book.Language, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"INSERT INTO language (code) VALUES (?)",
	)
	if err != nil {
		return book.Language{}, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, langCode)
	if err != nil {
		return book.Language{}, err
	}

	langID, err := result.LastInsertId()
	if err != nil {
		return book.Language{}, err
	}

	return book.Language{
		ID:   uint(langID),
		Code: langCode,
	}, nil
}

func (storage Storage) GetLanguage(ctx context.Context, langID uint) (book.Language, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT code FROM language WHERE id = ?",
	)
	if err != nil {
		return book.Language{}, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx, langID)

	lang := book.Language{}
	if err = result.Scan(&lang.Code); err != nil {
		return book.Language{}, err
	}

	lang.ID = langID
	return lang, nil
}

func (storage Storage) GetLanguages(ctx context.Context) ([]book.Language, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT id , code FROM language",
	)
	if err != nil {
		return []book.Language{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx)
	if err != nil {
		return []book.Language{}, err
	}
	defer result.Close()

	langs := []book.Language{}
	for result.Next() {
		var l book.Language

		if err := result.Scan(
			&l.ID,
			&l.Code,
		); err != nil {
			return []book.Language{}, err
		}

		langs = append(langs, l)
	}

	return langs, nil
}

func (storage Storage) DeleteLanguage(ctx context.Context, langID uint) error {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"DELETE FROM language WHERE id = ?",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, langID); err != nil {
		return err
	}

	return nil
}

func (storage Storage) AddBook(ctx context.Context, b book.Book) (book.Book, error) {

	result, err := storage.MySQL.ExecContext(ctx,
		`INSERT INTO book 
		(
			title , isbn , pages , description , year , date , digital_price , 
			physical_price , physical_stock , pdf , epub , djvu , azw , txt ,
			docx , lang_id , cover_front , cover_back , publisher , availability
		)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,

		b.Title,
		b.ISBN,
		b.Pages,
		b.Description,
		b.Year,
		b.CreationDate,
		b.Digital.Price,
		b.Physical.Price,
		b.Physical.Stock,
		b.Digital.PDF,
		b.Digital.EPUB,
		b.Digital.DJVU,
		b.Digital.AZW,
		b.Digital.TXT,
		b.Digital.DOCX,
		b.Language.ID,
		b.CoverFront,
		b.CoverBack,
		b.Publisher.ID,
		b.Availability,
	)

	if err != nil {
		return book.Book{}, err
	}

	bookID, err := result.LastInsertId()
	if err != nil {
		return book.Book{}, err
	}

	// Add book authors to book_author table
	for _, author := range b.Authors {
		if err = storage.AddBookAuthor(ctx, uint(bookID), author.ID); err != nil {
			return book.Book{}, err
		}
	}

	// Add book topics to book_topic table
	for _, topic := range b.Topics {
		if err = storage.AddBookTopic(ctx, uint(bookID), topic.ID); err != nil {
			return book.Book{}, err
		}
	}

	return b, nil
}

func (storage Storage) AddBookAuthor(ctx context.Context, bookID, authorID uint) error {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"INSERT IGNORE INTO book_author (book_id , author_id) VALUES (? , ?)",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx,
		bookID,
		authorID,
	); err != nil {
		return err
	}

	return nil
}

func (storage Storage) AddBookTopic(ctx context.Context, bookID, topicID uint) error {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"INSERT IGNORE INTO book_topic (book_id , topic_id) VALUES (? , ?)",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx,
		bookID,
		topicID,
	); err != nil {
		return err
	}

	return nil
}

func (storage Storage) SetBookDiscount(ctx context.Context, bookID, digital, physical uint) error {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"UPDATE book SET digital_discount = ? , physical_discount = ? WHERE id = ?",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx,
		digital,
		physical,
		bookID,
	); err != nil {
		return err
	}

	return nil
}

func (storage Storage) GetBook(ctx context.Context, bookID uint) (book.Book, error) {

	bookResult := storage.MySQL.QueryRowContext(ctx,

		`SELECT title , isbn , pages , description , year , date , 
		digital_price , digital_discount , physical_price , physical_discount , physical_stock , 
		lang_id , cover_front , cover_back , availability 
		FROM book 
		WHERE id = ?`,

		bookID,
	)

	b := book.Book{}
	b.ID = bookID

	if err := bookResult.Scan(
		&b.Title,
		&b.ISBN,
		&b.Pages,
		&b.Description,
		&b.Year,
		&b.CreationDate,
		&b.Digital.Price,
		&b.Digital.Discount,
		&b.Physical.Price,
		&b.Physical.Discount,
		&b.Physical.Stock,
		&b.Language.ID,
		&b.CoverFront,
		&b.CoverBack,
		&b.Availability,
	); err != nil {
		return book.Book{}, err
	}

	authors, err := storage.GetBookAuthors(ctx, bookID)
	if err != nil {
		return book.Book{}, err
	}
	b.Authors = authors

	topics, err := storage.GetBookTopics(ctx, bookID)
	if err != nil {
		return book.Book{}, err
	}
	b.Topics = topics

	return b, nil
}

func (storage Storage) GetBookAuthors(ctx context.Context, bookID uint) ([]book.Author, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT * FROM author WHERE id IN ( SELECT author_id FROM book_author WHERE book_id = ? )",
	)
	if err != nil {
		return []book.Author{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx, bookID)
	if err != nil {
		return []book.Author{}, err
	}
	defer result.Close()

	authors := []book.Author{}

	for result.Next() {
		var author book.Author

		if err = result.Scan(
			&author.ID,
			&author.Name,
		); err != nil {
			return []book.Author{}, err
		}

		authors = append(authors, author)
	}

	return authors, nil
}

func (storage Storage) GetBookTopics(ctx context.Context, bookID uint) ([]book.Topic, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT * FROM topic WHERE id IN ( SELECT topic_id FROM book_topic WHERE book_id = ? )",
	)
	if err != nil {
		return []book.Topic{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx, bookID)
	if err != nil {
		return []book.Topic{}, err
	}
	defer result.Close()

	topics := []book.Topic{}

	for result.Next() {
		var topic book.Topic

		if err = result.Scan(
			&topic.ID,
			&topic.Name,
		); err != nil {
			return []book.Topic{}, err
		}

		topics = append(topics, topic)
	}

	return topics, nil
}

func (storage Storage) EditBook(ctx context.Context, b book.Book) (book.Book, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		`UPDATE book SET
		title=? , isbn=? , pages=? , description=? , year=? , digital_price=? , 
		physical_price=? , physical_stock=? , pdf=? , epub=? , djvu=? , azw=? , 
		txt=? , docx=? , lang_id=? , cover_front=? , cover_back=? , publisher=? , availability=?
		WHERE id=? `,
	)
	if err != nil {
		return book.Book{}, err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx,
		b.Title,
		b.ISBN,
		b.Pages,
		b.Description,
		b.Year,
		b.Digital.Price,
		b.Physical.Price,
		b.Physical.Stock,
		b.Digital.PDF,
		b.Digital.EPUB,
		b.Digital.DJVU,
		b.Digital.AZW,
		b.Digital.TXT,
		b.Digital.DOCX,
		b.Language.ID,
		b.CoverFront,
		b.CoverBack,
		b.Publisher.ID,
		b.Availability,
		b.ID,
	); err != nil {
		return book.Book{}, err
	}

	return b, nil
}

func (storage Storage) GetAllBooksFull(ctx context.Context) ([]book.Book, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"SELECT * FROM book",
	)
	if err != nil {
		return []book.Book{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx)
	if err != nil {
		return []book.Book{}, err
	}
	defer result.Close()

	books := []book.Book{}

	for result.Next() {
		var b book.Book

		if err = result.Scan(
			&b.ID,
			&b.Title,
			&b.ISBN,
			&b.Pages,
			&b.Description,
			&b.Year,
			&b.CreationDate,
			&b.Digital.Price,
			&b.Digital.Discount,
			&b.Physical.Price,
			&b.Physical.Discount,
			&b.Physical.Stock,
			&b.Digital.PDF,
			&b.Digital.EPUB,
			&b.Digital.DJVU,
			&b.Digital.AZW,
			&b.Digital.TXT,
			&b.Digital.DOCX,
			&b.Language.ID,
			&b.CoverFront,
			&b.CoverBack,
			&b.Availability,
		); err != nil {
			return []book.Book{}, err
		}

		books = append(books, b)
	}

	return books, nil
}

func (storage Storage) GetAllBooks(ctx context.Context) ([]book.Book, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		`SELECT id , title , digital_price , digital_discount , physical_price , 
		physical_discount , physical_stock , cover_front , availability FROM book`,
	)
	if err != nil {
		return []book.Book{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx)
	if err != nil {
		return []book.Book{}, err
	}
	defer result.Close()

	books := []book.Book{}

	for result.Next() {
		var b book.Book

		if err = result.Scan(
			&b.ID,
			&b.Title,
			&b.Digital.Price,
			&b.Digital.Discount,
			&b.Physical.Price,
			&b.Physical.Discount,
			&b.Physical.Stock,
			&b.CoverFront,
			&b.Availability,
		); err != nil {
			return []book.Book{}, err
		}

		books = append(books, b)
	}
	return books, nil
}

func (storage Storage) GetAuthorBooks(ctx context.Context, authorID uint) ([]book.Book, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		`SELECT id , title , digital_price , digital_discount ,physical_price , 
		physical_discount , physical_stock , cover_front , availability FROM book 
		WHERE id IN ( SELECT book_id FROM book_author WHERE author_id = ? )`,
	)
	if err != nil {
		return []book.Book{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx, authorID)
	if err != nil {
		return []book.Book{}, err
	}
	defer result.Close()

	books := []book.Book{}
	for result.Next() {
		var b book.Book

		if err = result.Scan(
			&b.ID,
			&b.Title,
			&b.Digital.Price,
			&b.Digital.Discount,
			&b.Physical.Price,
			&b.Physical.Discount,
			&b.Physical.Stock,
			&b.CoverFront,
			&b.Availability,
		); err != nil {
			return []book.Book{}, err
		}
		books = append(books, b)
	}

	return books, nil
}

func (storage Storage) GetTopicBooks(ctx context.Context, topicID uint) ([]book.Book, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		`SELECT id , title , digital_price , digital_discount , physical_price , 
		physical_discount , physical_stock , cover_front , availability FROM book 
		WHERE id IN ( SELECT book_id FROM book_topic WHERE topic_id = ? )`,
	)
	if err != nil {
		return []book.Book{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx, topicID)
	if err != nil {
		return []book.Book{}, err
	}
	defer result.Close()

	books := []book.Book{}
	for result.Next() {
		var b book.Book

		if err = result.Scan(
			&b.ID,
			&b.Title,
			&b.Digital.Price,
			&b.Digital.Discount,
			&b.Physical.Price,
			&b.Physical.Discount,
			&b.Physical.Stock,
			&b.CoverFront,
			&b.Availability,
		); err != nil {
			return []book.Book{}, err
		}
		books = append(books, b)
	}

	return books, nil
}

func (storage Storage) GetPublisherBooks(ctx context.Context, publisherID uint) ([]book.Book, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		`SELECT id , title , digital_price , digital_discount , physical_price , 
		physical_discount , physical_stock , cover_front , availability FROM book 
		WHERE publisher = ?`,
	)
	if err != nil {
		return []book.Book{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx, publisherID)
	if err != nil {
		return []book.Book{}, err
	}
	defer result.Close()

	books := []book.Book{}
	for result.Next() {
		var b book.Book

		if err = result.Scan(
			&b.ID,
			&b.Title,
			&b.Digital.Price,
			&b.Digital.Discount,
			&b.Physical.Price,
			&b.Physical.Discount,
			&b.Physical.Stock,
			&b.CoverFront,
			&b.Availability,
		); err != nil {
			return []book.Book{}, err
		}

		books = append(books, b)
	}

	return books, nil
}

func (storage Storage) GetLangBooks(ctx context.Context, langID uint) ([]book.Book, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		`SELECT id , title , digital_price , digital_discount , physical_price , 
		physical_discount , physical_stock , cover_front , availability FROM book 
		WHERE lang_id = ?`,
	)
	if err != nil {
		return []book.Book{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx, langID)
	if err != nil {
		return []book.Book{}, err
	}
	defer result.Close()

	books := []book.Book{}
	for result.Next() {
		var b book.Book

		if err = result.Scan(
			&b.ID,
			&b.Title,
			&b.Digital.Price,
			&b.Digital.Discount,
			&b.Physical.Price,
			&b.Physical.Discount,
			&b.Physical.Stock,
			&b.CoverFront,
			&b.Availability,
		); err != nil {
			return []book.Book{}, err
		}

		books = append(books, b)
	}

	return books, nil
}

func (storage Storage) DeleteBook(ctx context.Context, bookID uint) error {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		"DELETE FROM book WHERE id = ?",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, bookID); err != nil {
		return err
	}

	return nil
}

func (storage Storage) GetUserDigitalBooks(ctx context.Context, userID string) ([]book.Book, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		`SELECT id , title , isbn , pages , description , year , 
		pdf , epub , djvu , azw , txt , docx , 
		lang_id , cover_front , publisher FROM book 
		WHERE book.id IN 
		(SELECT book_id FROM item 
			WHERE item.order_id = (SELECT id FROM orders WHERE user_id = ? AND status != ?) 
			AND 
			type = 0
		)`,
	)
	if err != nil {
		return []book.Book{}, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx,
		userID,
		order.StatusCreated,
	)
	if err != nil {
		return []book.Book{}, err
	}

	books := []book.Book{}
	for result.Next() {
		var b book.Book

		if err = result.Scan(
			&b.ID,
			&b.Title,
			&b.ISBN,
			&b.Pages,
			&b.Description,
			&b.Year,
			&b.Digital.PDF,
			&b.Digital.EPUB,
			&b.Digital.DJVU,
			&b.Digital.AZW,
			&b.Digital.TXT,
			&b.Digital.DOCX,
			&b.Language.ID,
			&b.CoverFront,
			&b.Publisher.ID,
		); err != nil {
			return []book.Book{}, err
		}

		books = append(books, b)
	}

	return books, nil
}

func (storage Storage) DoesUserAccessBook(ctx context.Context, userID string, bookID uint) (bool, error) {

	stmt, err := storage.MySQL.PrepareContext(ctx,
		`SELECT 1 FROM item WHERE type = ? 
		AND book_id = ? 
		AND order_id IN ( SELECT id FROM orders WHERE user_id = ? AND status != ? )`,
	)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	result := stmt.QueryRowContext(ctx,
		order.Digital,
		bookID,
		userID,
		order.StatusCreated,
	)

	var access bool
	if err = result.Scan(&access); err != nil {
		return false, err
	}

	return access, nil
}
