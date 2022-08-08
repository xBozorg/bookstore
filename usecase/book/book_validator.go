package book

import (
	"context"

	"github.com/XBozorg/bookstore/dto"
)

type (
	ValidateAddAuthor    func(ctx context.Context, req dto.AddAuthorRequest) error
	ValidateGetAuthor    func(ctx context.Context, req dto.GetAuthorRequest) error
	ValidateDeleteAuthor func(ctx context.Context, req dto.DeleteAuthorRequest) error

	ValidateAddPublisher    func(ctx context.Context, req dto.AddPublisherRequest) error
	ValidateGetPublisher    func(ctx context.Context, req dto.GetPublisherRequest) error
	ValidateGetPublishers   func(ctx context.Context, req dto.GetPublishersRequest) error
	ValidateDeletePublisher func(ctx context.Context, req dto.DeletePublisherRequest) error

	ValidateAddTopic    func(ctx context.Context, req dto.AddTopicRequest) error
	ValidateGetTopic    func(ctx context.Context, req dto.GetTopicRequest) error
	ValidateGetTopics   func(ctx context.Context, req dto.GetTopicsRequest) error
	ValidateDeleteTopic func(ctx context.Context, req dto.DeleteTopicRequest) error

	ValidateAddLanguage    func(ctx context.Context, req dto.AddLanguageRequest) error
	ValidateGetLanguage    func(ctx context.Context, req dto.GetLanguageRequest) error
	ValidateGetLanguages   func(ctx context.Context, req dto.GetLanguagesRequest) error
	ValidateDeleteLanguage func(ctx context.Context, req dto.DeleteLanguageRequest) error

	ValidateAddBook           func(ctx context.Context, req dto.AddBookRequest) error
	ValidateGetBook           func(ctx context.Context, req dto.GetBookRequest) error
	ValidateEditBook          func(ctx context.Context, req dto.EditBookRequest) error
	ValidateSetBookDiscount   func(ctx context.Context, req dto.SetBookDiscountRequest) error
	ValidateGetAllBooks       func(ctx context.Context, req dto.GetAllBooksRequest) error
	ValidateGetAuthorBooks    func(ctx context.Context, req dto.GetAuthorBooksRequest) error
	ValidateGetPublisherBooks func(ctx context.Context, req dto.GetPublisherBooksRequest) error
	ValidateGetTopicBooks     func(ctx context.Context, req dto.GetTopicBooksRequest) error
	ValidateGetLangBooks      func(ctx context.Context, req dto.GetLangBooksRequest) error
	ValidateDeleteBook        func(ctx context.Context, req dto.DeleteBookRequest) error
)
