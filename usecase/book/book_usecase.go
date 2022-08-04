package book

import (
	"context"

	"github.com/XBozorg/bookstore/dto"
)

type UseCase interface {
	AddAuthor(ctx context.Context, req dto.AddAuthorRequest) (dto.AddAuthorResponse, error)
	GetAuthor(ctx context.Context, req dto.GetAuthorRequest) (dto.GetAuthorResponse, error)
	GetAuthors(ctx context.Context, req dto.GetAuthorsRequest) (dto.GetAuthorsResponse, error)
	DeleteAuthor(ctx context.Context, req dto.DeleteAuthorRequest) (dto.DeleteAuthorResponse, error)

	AddPublisher(ctx context.Context, req dto.AddPublisherRequest) (dto.AddPublisherResponse, error)
	GetPublisher(ctx context.Context, req dto.GetPublisherRequest) (dto.GetPublisherResponse, error)
	GetBookPublishers(ctx context.Context, req dto.GetBookPublishersRequest) (dto.GetBookPublishersResponse, error)
	GetAllPublishers(ctx context.Context, req dto.GetAllPublishersRequest) (dto.GetAllPublishersResponse, error)
	DeletePublisher(ctx context.Context, req dto.DeletePublisherRequest) (dto.DeletePublisherResponse, error)

	AddTopic(ctx context.Context, req dto.AddTopicRequest) (dto.AddTopicResponse, error)
	GetTopic(ctx context.Context, req dto.GetTopicRequest) (dto.GetTopicResponse, error)
	GetBookTopics(ctx context.Context, req dto.GetBookTopicsRequest) (dto.GetAllTopicsResponse, error)
	GetAllTopics(ctx context.Context, req dto.GetAllTopicsRequest) (dto.GetAllTopicsResponse, error)
	DeleteTopic(ctx context.Context, req dto.DeleteTopicRequest) (dto.DeleteTopicResponse, error)

	AddLanguage(ctx context.Context, req dto.AddLanguageRequest) (dto.AddLanguageResponse, error)
	GetLanguage(ctx context.Context, req dto.GetLanguageRequest) (dto.GetLanguageResponse, error)
	GetAllLanguages(ctx context.Context, req dto.GetAllLanguagesRequest) (dto.GetAllLanguagesResponse, error)
	GetBookLanguages(ctx context.Context, req dto.GetBookLanguagesRequest) (dto.GetBookLanguagesResponse, error)
	DeleteLanguage(ctx context.Context, req dto.DeleteLanguageRequest) (dto.DeleteLanguageResponse, error)

	AddCover(ctx context.Context, req dto.AddCoverRequest) (dto.AddCoverResponse, error)
	GetCover(ctx context.Context, req dto.GetCoverRequest) (dto.GetCoverResponse, error)
	DeleteCover(ctx context.Context, req dto.DeleteCoverRequest) (dto.DeleteCoverResponse, error)

	AddBook(ctx context.Context, req dto.AddBookRequest) (dto.AddBookResponse, error)
	GetBook(ctx context.Context, req dto.GetBookRequest) (dto.GetBookResponse, error)
	EditBook(ctx context.Context, req dto.EditBookRequest) (dto.EditBookResponse, error)
	GetAllBooks(ctx context.Context, req dto.GetAllBooksRequest) (dto.GetAllBooksResponse, error)
	GetAuthorBooks(ctx context.Context, req dto.GetAuthorBooksRequest) (dto.GetAuthorBooksResponse, error)
	GetPublisherBooks(ctx context.Context, req dto.GetPublisherBooksRequest) (dto.GetPublisherBooksResponse, error)
	GetTopicBooks(ctx context.Context, req dto.GetTopicBooksRequest) (dto.GetTopicBooksResponse, error)
	GetLangBooks(ctx context.Context, req dto.GetLangBooksRequest) (dto.GetLangBooksResponse, error)
	DeleteBook(ctx context.Context, req dto.DeleteBookRequest) (dto.DeleteBookResponse, error)
}
