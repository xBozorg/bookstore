package repository

import (
	"context"

	"github.com/XBozorg/bookstore/entity/book"
)

func (m MySQLRepo) AddAuthor(ctx context.Context, authorName string) (book.Author, error) {

	result, err := m.db.ExecContext(ctx, "INSERT INTO author (name) VALUES (?)", authorName)

	if err != nil {
		return book.Author{}, err
	}

	authorID, err := result.LastInsertId()
	if err != nil {
		return book.Author{}, err
	}

	return book.Author{ID: uint(authorID), Name: authorName}, nil
}

func (m MySQLRepo) GetAuthor(ctx context.Context, authorID uint) (book.Author, error) {

	result := m.db.QueryRowContext(ctx, "SELECT name From author WHERE id = ?", authorID)

	author := book.Author{}
	err := result.Scan(&author.Name)

	if err != nil {
		return book.Author{}, err
	}

	author.ID = authorID
	return author, nil
}

func (m MySQLRepo) GetAuthors(ctx context.Context) ([]book.Author, error) {

	result, err := m.db.QueryContext(ctx, "SELECT id , name FROM author")

	if err != nil {
		return []book.Author{}, err
	}
	defer result.Close()

	authors := []book.Author{}
	for result.Next() {
		var a book.Author
		err := result.Scan(&a.ID, &a.Name)
		if err != nil {
			return []book.Author{}, err
		}
		authors = append(authors, a)
	}

	return authors, nil
}

func (m MySQLRepo) DeleteAuthor(ctx context.Context, authorID uint) error {

	_, err := m.db.ExecContext(ctx, "DELETE FROM author WHERE id = ?", authorID)

	if err != nil {
		return err
	}

	return nil
}

func (m MySQLRepo) AddPublisher(ctx context.Context, publisherName string) (book.Publisher, error) {

	result, err := m.db.ExecContext(ctx, "INSERT IGNORE INTO publisher (name) VALUES (?)", publisherName)

	if err != nil {
		return book.Publisher{}, err
	}

	publisherID, err := result.LastInsertId()
	if err != nil {
		return book.Publisher{}, err
	}

	return book.Publisher{ID: uint(publisherID), Name: publisherName}, nil
}

func (m MySQLRepo) GetPublisher(ctx context.Context, publisherID uint) (book.Publisher, error) {

	result := m.db.QueryRowContext(ctx, "SELECT name FROM publisher WHERE id = ?", publisherID)

	publisher := book.Publisher{}
	err := result.Scan(&publisher.Name)

	if err != nil {
		return book.Publisher{}, err
	}

	publisher.ID = publisherID
	return publisher, nil
}

func (m MySQLRepo) GetAllPublishers(ctx context.Context) ([]book.Publisher, error) {

	result, err := m.db.QueryContext(ctx, "SELECT id , name FROM publisher")

	if err != nil {
		return []book.Publisher{}, err
	}
	defer result.Close()

	publishers := []book.Publisher{}
	for result.Next() {
		var p book.Publisher
		err := result.Scan(&p.ID, &p.Name)
		if err != nil {
			return []book.Publisher{}, err
		}
		publishers = append(publishers, p)
	}

	return publishers, nil
}

func (m MySQLRepo) DeletePublisher(ctx context.Context, publisherId uint) error {

	_, err := m.db.ExecContext(ctx, "DELETE FROM publisher WHERE id = ?", publisherId)

	if err != nil {
		return err
	}

	return nil
}

func (m MySQLRepo) AddTopic(ctx context.Context, topicName string) (book.Topic, error) {

	result, err := m.db.ExecContext(ctx, "INSERT IGNORE INTO topic (name) VALUES (?)", topicName)

	if err != nil {
		return book.Topic{}, err
	}

	topicID, err := result.LastInsertId()
	if err != nil {
		return book.Topic{}, err
	}

	return book.Topic{ID: uint(topicID), Name: topicName}, nil
}

func (m MySQLRepo) GetTopic(ctx context.Context, topicID uint) (book.Topic, error) {

	result := m.db.QueryRowContext(ctx, "SELECT name FROM topic WHERE id = ?", topicID)

	topic := book.Topic{}
	err := result.Scan(&topic.Name)

	if err != nil {
		return book.Topic{}, err
	}

	topic.ID = topicID
	return topic, nil
}

func (m MySQLRepo) GetAllTopics(ctx context.Context) ([]book.Topic, error) {

	result, err := m.db.QueryContext(ctx, "SELECT id , name FROM topic")

	if err != nil {
		return []book.Topic{}, err
	}
	defer result.Close()

	topics := []book.Topic{}
	for result.Next() {
		var t book.Topic
		err := result.Scan(&t.ID, &t.Name)
		if err != nil {
			return []book.Topic{}, err
		}
		topics = append(topics, t)
	}

	return topics, nil
}

func (m MySQLRepo) DeleteTopic(ctx context.Context, topicID uint) error {

	_, err := m.db.ExecContext(ctx, "DELETE FROM topic WHERE id = ?", topicID)

	if err != nil {
		return err
	}

	return nil
}

func (m MySQLRepo) AddLanguage(ctx context.Context, langCode string) (book.Language, error) {

	result, err := m.db.ExecContext(ctx, "INSERT IGNORE INTO language (code) VALUES (?)", langCode)

	if err != nil {
		return book.Language{}, err
	}

	langID, err := result.LastInsertId()
	if err != nil {
		return book.Language{}, err
	}

	return book.Language{ID: uint(langID), Code: langCode}, nil
}

func (m MySQLRepo) GetLanguage(ctx context.Context, langID uint) (book.Language, error) {

	result := m.db.QueryRowContext(ctx, "SELECT code FROM language WHERE id = ?", langID)

	lang := book.Language{}
	err := result.Scan(&lang.Code)

	if err != nil {
		return book.Language{}, err
	}

	lang.ID = langID
	return lang, nil
}

func (m MySQLRepo) GetAllLanguages(ctx context.Context) ([]book.Language, error) {

	result, err := m.db.QueryContext(ctx, "SELECT id , code FROM language")

	if err != nil {
		return []book.Language{}, err
	}
	defer result.Close()

	langs := []book.Language{}
	for result.Next() {
		var l book.Language
		err := result.Scan(&l.ID, &l.Code)
		if err != nil {
			return []book.Language{}, err
		}

		langs = append(langs, l)
	}

	return langs, nil
}

func (m MySQLRepo) DeleteLanguage(ctx context.Context, langID uint) error {

	_, err := m.db.ExecContext(ctx, "DELETE FROM language WHERE id = ?", langID)

	if err != nil {
		return err
	}

	return nil
}

func (m MySQLRepo) AddBook(ctx context.Context, b book.Book) (book.Book, error) {

	result, err := m.db.ExecContext(ctx, `INSERT IGNORE INTO book 

	(title , isbn , pages , description , year , date , digital_price , 
	physical_price , physical_stock , pdf , epub , djvu , azw , txt , docx , lang_id , cover_front , cover_back , publisher)

		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,

		b.Title, b.ISBN, b.Pages, b.Description, b.Year, b.CreationDate,
		b.Digital.Price, b.Physical.Price, b.Physical.Stock,
		b.Digital.PDF, b.Digital.EPUB, b.Digital.DJVU, b.Digital.AZW, b.Digital.TXT, b.Digital.DOCX, b.Language.ID, b.CoverFront, b.CoverBack, b.Publisher.ID)

	if err != nil {
		return book.Book{}, err
	}

	bookID, err := result.LastInsertId()
	if err != nil {
		return book.Book{}, err
	}

	// Add book authors to book_author table
	for _, author := range b.Authors {
		err = m.AddBookAuthor(ctx, uint(bookID), author.ID)
		if err != nil {
			return book.Book{}, err
		}
	}

	// Add book topics to book_topic table
	for _, topic := range b.Topics {
		err = m.AddBookTopic(ctx, uint(bookID), topic.ID)
		if err != nil {
			return book.Book{}, err
		}
	}

	return b, nil
}

func (m MySQLRepo) AddBookAuthor(ctx context.Context, bookID, authorID uint) error {

	_, err := m.db.ExecContext(ctx, "INSERT IGNORE INTO book_author (book_id , author_id) VALUES (? , ?)", bookID, authorID)

	if err != nil {
		return err
	}

	return nil
}

func (m MySQLRepo) AddBookTopic(ctx context.Context, bookID, topicID uint) error {

	_, err := m.db.ExecContext(ctx, "INSERT IGNORE INTO book_topic (book_id , topic_id) VALUES (? , ?)", bookID, topicID)

	if err != nil {
		return err
	}

	return nil
}

func (m MySQLRepo) SetBookDiscount(ctx context.Context, bookID, digital, physical uint) error {

	_, err := m.db.ExecContext(ctx, "UPDATE book SET digital_discount = ? , physical_discount = ? WHERE id = ?", digital, physical, bookID)

	if err != nil {
		return err
	}

	return nil
}
func (m MySQLRepo) GetBook(ctx context.Context, bookID uint) (book.Book, error) {

	bookResult := m.db.QueryRowContext(ctx, `SELECT title , isbn , pages , description , year , date ,
	digital_price , digital_discount , physical_price , physical_discount , physical_stock ,
	lang_id , cover_front , cover_back FROM book WHERE id = ?`, bookID)

	b := book.Book{}
	b.ID = bookID

	err := bookResult.Scan(&b.Title, &b.ISBN, &b.Pages, &b.Description, &b.Year, &b.CreationDate, &b.Digital.Price,
		&b.Digital.Discount, &b.Physical.Price, &b.Physical.Discount, &b.Physical.Stock, &b.Language.ID, &b.CoverFront, &b.CoverBack)
	if err != nil {
		return book.Book{}, err
	}

	authors, err := m.GetBookAuthors(ctx, bookID)
	if err != nil {
		return book.Book{}, err
	}
	b.Authors = authors

	topics, err := m.GetBookTopics(ctx, bookID)
	if err != nil {
		return book.Book{}, err
	}
	b.Topics = topics

	return b, nil
}

func (m MySQLRepo) GetBookAuthors(ctx context.Context, bookID uint) ([]book.Author, error) {

	result, err := m.db.QueryContext(ctx, "SELECT * FROM author WHERE id IN ( SELECT author_id FROM book_author WHERE book_id = ? )", bookID)

	if err != nil {
		return []book.Author{}, err
	}
	defer result.Close()

	authors := []book.Author{}

	for result.Next() {
		var author book.Author
		err := result.Scan(&author.ID, &author.Name)
		if err != nil {
			return []book.Author{}, err
		}

		authors = append(authors, author)
	}

	return authors, nil
}

func (m MySQLRepo) GetBookTopics(ctx context.Context, bookID uint) ([]book.Topic, error) {

	result, err := m.db.QueryContext(ctx, "SELECT * FROM topic WHERE id IN ( SELECT topic_id FROM book_topic WHERE book_id = ? )", bookID)
	if err != nil {
		return []book.Topic{}, err
	}
	defer result.Close()

	topics := []book.Topic{}

	for result.Next() {
		var topic book.Topic
		err := result.Scan(&topic.ID, &topic.Name)
		if err != nil {
			return []book.Topic{}, err
		}

		topics = append(topics, topic)
	}

	return topics, nil
}

func (m MySQLRepo) EditBook(ctx context.Context, b book.Book) (book.Book, error) {

	_, err := m.db.ExecContext(ctx, `UPDATE book SET
		title=? , isbn=? , pages=? , description=? , year=? , digital_price=? , 
		physical_price=? , physical_stock=? , pdf=? , epub=? , djvu=? , azw=? , txt=? , docx=? , lang_id=? , cover_front=? , cover_back=? , publisher=?
		WHERE id=? `, b.Title, b.ISBN, b.Pages, b.Description, b.Year, b.Digital.Price,
		b.Physical.Price, b.Physical.Stock, b.Digital.PDF, b.Digital.EPUB, b.Digital.DJVU, b.Digital.AZW, b.Digital.TXT, b.Digital.DOCX, b.Language.ID, b.CoverFront, b.CoverBack, b.Publisher.ID, b.ID)

	if err != nil {
		return book.Book{}, err
	}

	return b, nil
}

func (m MySQLRepo) GetAllBooksFull(ctx context.Context) ([]book.Book, error) {

	result, err := m.db.QueryContext(ctx, `SELECT * FROM book`)
	if err != nil {
		return []book.Book{}, err
	}
	defer result.Close()

	books := []book.Book{}

	for result.Next() {
		var b book.Book
		err := result.Scan(&b.ID, &b.Title, &b.ISBN, &b.Pages, &b.Description, &b.Year, &b.CreationDate,
			&b.Digital.Price, &b.Digital.Discount, &b.Physical.Price, &b.Physical.Discount,
			&b.Physical.Stock, &b.Digital.PDF, &b.Digital.EPUB, &b.Digital.DJVU, &b.Digital.AZW, &b.Digital.TXT, &b.Digital.DOCX,
			&b.Language.ID, &b.CoverFront, &b.CoverBack)
		if err != nil {
			return []book.Book{}, err
		}

		books = append(books, b)
	}
	return books, nil
}

func (m MySQLRepo) GetAllBooks(ctx context.Context) ([]book.Book, error) {

	result, err := m.db.QueryContext(ctx, `SELECT id , title , digital_price , digital_discount ,
	physical_price , physical_discount , physical_stock , cover_front FROM book`)
	if err != nil {
		return []book.Book{}, err
	}
	defer result.Close()

	books := []book.Book{}

	for result.Next() {
		var b book.Book
		err := result.Scan(&b.ID, &b.Title, &b.Digital.Price, &b.Digital.Discount, &b.Physical.Price, &b.Physical.Discount,
			&b.Physical.Stock, &b.CoverFront)
		if err != nil {
			return []book.Book{}, err
		}

		books = append(books, b)
	}
	return books, nil
}

func (m MySQLRepo) GetAuthorBooks(ctx context.Context, authorID uint) ([]book.Book, error) {

	result, err := m.db.QueryContext(ctx, `SELECT id , title , digital_price , digital_discount ,
	physical_price , physical_discount , physical_stock , cover_front FROM book
	WHERE id IN (SELECT book_id FROM book_author WHERE author_id = ?)`, authorID)

	if err != nil {
		return []book.Book{}, err
	}
	defer result.Close()

	books := []book.Book{}
	for result.Next() {
		var b book.Book
		err := result.Scan(&b.ID, &b.Title, &b.Digital.Price, &b.Digital.Discount, &b.Physical.Price, &b.Physical.Discount,
			&b.Physical.Stock, &b.CoverFront)
		if err != nil {
			return []book.Book{}, err
		}
		books = append(books, b)
	}

	return books, nil
}

func (m MySQLRepo) GetTopicBooks(ctx context.Context, topicID uint) ([]book.Book, error) {

	result, err := m.db.QueryContext(ctx, `SELECT id , title , digital_price , digital_discount ,
	physical_price , physical_discount , physical_stock , cover_front FROM book
	WHERE id IN (SELECT book_id FROM book_topic WHERE topic_id = ?)`, topicID)

	if err != nil {
		return []book.Book{}, err
	}
	defer result.Close()

	books := []book.Book{}
	for result.Next() {
		var b book.Book
		err := result.Scan(&b.ID, &b.Title, &b.Digital.Price, &b.Digital.Discount, &b.Physical.Price, &b.Physical.Discount,
			&b.Physical.Stock, &b.CoverFront)
		if err != nil {
			return []book.Book{}, err
		}
		books = append(books, b)
	}

	return books, nil
}

func (m MySQLRepo) GetPublisherBooks(ctx context.Context, publisherID uint) ([]book.Book, error) {

	result, err := m.db.QueryContext(ctx, `SELECT id , title , digital_price , digital_discount ,
	physical_price , physical_discount , physical_stock , cover_front FROM book WHERE publisher = ?`, publisherID)

	if err != nil {
		return []book.Book{}, err
	}
	defer result.Close()

	books := []book.Book{}
	for result.Next() {
		var b book.Book
		err := result.Scan(&b.ID, &b.Title, &b.Digital.Price, &b.Digital.Discount, &b.Physical.Price, &b.Physical.Discount,
			&b.Physical.Stock, &b.CoverFront)
		if err != nil {
			return []book.Book{}, err
		}

		books = append(books, b)
	}

	return books, nil
}

func (m MySQLRepo) GetLangBooks(ctx context.Context, langID uint) ([]book.Book, error) {

	result, err := m.db.QueryContext(ctx, `SELECT id , title , digital_price , digital_discount ,
	physical_price , physical_discount , physical_stock , cover_front FROM book WHERE lang_id = ?`, langID)

	if err != nil {
		return []book.Book{}, err
	}
	defer result.Close()

	books := []book.Book{}
	for result.Next() {
		var b book.Book
		err := result.Scan(&b.ID, &b.Title, &b.Digital.Price, &b.Digital.Discount, &b.Physical.Price, &b.Physical.Discount,
			&b.Physical.Stock, &b.CoverFront)
		if err != nil {
			return []book.Book{}, err
		}

		books = append(books, b)
	}

	return books, nil
}

func (m MySQLRepo) DeleteBook(ctx context.Context, bookID uint) error {

	_, err := m.db.ExecContext(ctx, "DELETE FROM book WHERE id = ?", bookID)

	if err != nil {
		return err
	}

	return nil
}
