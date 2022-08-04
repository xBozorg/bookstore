package book

import (
	"context"

	"github.com/XBozorg/bookstore/dto"
)

type (
	ValidateAddAuthor    func(req dto.AddAuthorRequest) error
	ValidateGetAuthor    func(ctx context.Context, req dto.GetAuthorRequest) error
	ValidateGetAuthors   func(ctx context.Context, req dto.GetAuthorsRequest) error
	ValidateDeleteAuthor func(ctx context.Context, req dto.DeleteAuthorRequest) error

	ValidateAddPublisher      func(ctx context.Context, req dto.AddPublisherRequest) error
	ValidateGetPublisher      func(ctx context.Context, req dto.GetPublisherRequest) error
	ValidateGetBookPublishers func(ctx context.Context, req dto.GetBookPublishersRequest) error
	ValidateAllPublishers     func(ctx context.Context, req dto.GetAllPublishersRequest) error
	ValidateDeletePublisher   func(ctx context.Context, req dto.DeletePublisherRequest) error

	ValidateAddTopic      func(ctx context.Context, req dto.AddTopicRequest) error
	ValidateGetTopic      func(ctx context.Context, req dto.GetTopicRequest) error
	ValidateGetBookTopics func(ctx context.Context, req dto.GetBookTopicsRequest) error
	ValidateGetAllTopics  func(ctx context.Context, req dto.GetAllTopicsRequest) error
	ValidateDeleteTopic   func(ctx context.Context, req dto.DeleteTopicRequest) error

	ValidateAddLanguage      func(ctx context.Context, req dto.AddLanguageRequest) error
	ValidateGetLanguage      func(ctx context.Context, req dto.GetLanguageRequest) error
	ValidateGetAllLanguages  func(ctx context.Context, req dto.GetAllLanguagesRequest) error
	ValidateGetBookLanguages func(ctx context.Context, req dto.GetBookLanguagesRequest) error
	ValidateDeleteLanguage   func(ctx context.Context, req dto.DeleteLanguageRequest) error

	ValidateAddCover    func(ctx context.Context, req dto.AddCoverRequest) error
	ValidateGetCover    func(ctx context.Context, req dto.GetCoverRequest) error
	ValidateDeleteCover func(ctx context.Context, req dto.DeleteCoverRequest) error

	ValidateAddBook           func(ctx context.Context, req dto.AddBookRequest) error
	ValidateGetBook           func(ctx context.Context, req dto.GetBookRequest) error
	ValidateEditBook          func(ctx context.Context, req dto.EditBookRequest) error
	ValidateGetAllBooks       func(ctx context.Context, req dto.GetAllBooksRequest) error
	ValidateGetAuthorBooks    func(ctx context.Context, req dto.GetAuthorBooksRequest) error
	ValidateGetPublisherBooks func(ctx context.Context, req dto.GetPublisherBooksRequest) error
	ValidateGetTopicBooks     func(ctx context.Context, req dto.GetTopicBooksRequest) error
	ValidateGetLangBooks      func(ctx context.Context, req dto.GetLangBooksRequest) error
	ValidateDeleteBook        func(ctx context.Context, req dto.DeleteBookRequest) error
)
