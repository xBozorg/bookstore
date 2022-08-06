package book

import (
	"context"

	"github.com/XBozorg/bookstore/entity/book"
)

type Repository interface {
	AddAuthor(ctx context.Context, authorName string) (book.Author, error)
	GetAuthor(ctx context.Context, authorID uint) (book.Author, error)
	GetAuthors(ctx context.Context) ([]book.Author, error)
	DeleteAuthor(ctx context.Context, authorID uint) error

	AddPublisher(ctx context.Context, publisherName string) (book.Publisher, error)
	GetPublisher(ctx context.Context, publisherID uint) (book.Publisher, error)
	GetBookPublishers(ctx context.Context, bookID uint) ([]book.Publisher, error)
	GetAllPublishers(ctx context.Context) ([]book.Publisher, error)
	DeletePublisher(ctx context.Context, publisherId uint) error

	AddTopic(ctx context.Context, topicName string) (book.Topic, error)
	GetTopic(ctx context.Context, topicID uint) (book.Topic, error)
	GetAllTopics(ctx context.Context) ([]book.Topic, error)
	DeleteTopic(ctx context.Context, topicID uint) error

	AddLanguage(ctx context.Context, langCode string) (book.Language, error)
	GetLanguage(ctx context.Context, langID uint) (book.Language, error)
	GetAllLanguages(ctx context.Context) ([]book.Language, error)
	GetBookLanguages(ctx context.Context, bookID uint) ([]book.Language, error)
	DeleteLanguage(ctx context.Context, langID uint) error

	AddBook(ctx context.Context, b book.Book) (book.Book, error)
	AddBookAuthor(ctx context.Context, bookID, authorID uint) error
	AddBookTopic(ctx context.Context, bookID, topicID uint) error
	GetBook(ctx context.Context, bookID uint) (book.Book, error)
	GetBookAuthors(ctx context.Context, bookID uint) ([]book.Author, error)
	GetBookTopics(ctx context.Context, bookID uint) ([]book.Topic, error)
	EditBook(ctx context.Context, b book.Book) (book.Book, error)
	GetAllBooks(ctx context.Context) ([]book.Book, error)
	GetAuthorBooks(ctx context.Context, authorID uint) ([]book.Book, error)
	GetPublisherBooks(ctx context.Context, publisherID uint) ([]book.Book, error)
	GetTopicBooks(ctx context.Context, topicID uint) ([]book.Book, error)
	GetLangBooks(ctx context.Context, langID uint) ([]book.Book, error)
	DeleteBook(ctx context.Context, bookID uint) error
}

type ValidatorRepo interface {
	DoesAuthorExist(ctx context.Context, authorID uint) (bool, error)
	DoesPublisherExist(ctx context.Context, publisherID uint) (bool, error)
	DoesTopicExist(ctx context.Context, topicID uint) (bool, error)
	DoesLanguageExist(ctx context.Context, langID uint) (bool, error)
	DoesBookExist(ctx context.Context, bookID uint) (bool, error)
}
