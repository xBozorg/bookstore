package dto

import "github.com/XBozorg/bookstore/entity/book"

type AddAuthorRequest struct {
	Name string `json:"name"`
}
type AddAuthorResponse struct {
	Author book.Author `json:"author"`
}

type GetAuthorRequest struct {
	AuthorID uint `json:"authorID"`
}
type GetAuthorResponse struct {
	Author book.Author `json:"author"`
}

type GetAuthorsRequest struct{}
type GetAuthorsResponse struct {
	Authors []book.Author `json:"authors"`
}

type DeleteAuthorRequest struct {
	AuthorID uint `json:"authorID"`
}
type DeleteAuthorResponse struct{}

type AddPublisherRequest struct {
	Name string `json:"name"`
}
type AddPublisherResponse struct {
	Publisher book.Publisher `json:"publisher"`
}

type GetPublisherRequest struct {
	PublisherID uint `json:"publisherID"`
}
type GetPublisherResponse struct {
	Publisher book.Publisher `json:"publisher"`
}

type GetPublishersRequest struct{}
type GetPublishersResponse struct {
	Publishers []book.Publisher `json:"publishers"`
}

type DeletePublisherRequest struct {
	PublisherID uint `json:"publisherID"`
}
type DeletePublisherResponse struct{}

type AddTopicRequest struct {
	Name string `json:"name"`
}
type AddTopicResponse struct {
	Topic book.Topic `json:"topic"`
}

type GetTopicRequest struct {
	TopicID uint `json:"topicID"`
}
type GetTopicResponse struct {
	Topic book.Topic `json:"topic"`
}

type GetTopicsRequest struct{}
type GetTopicsResponse struct {
	Topics []book.Topic `json:"topics"`
}

type DeleteTopicRequest struct {
	TopicID uint `json:"topicID"`
}
type DeleteTopicResponse struct{}

type AddLanguageRequest struct {
	LangCode string `json:"langCode"`
}
type AddLanguageResponse struct {
	Language book.Language `json:"language"`
}

type GetLanguageRequest struct {
	LangID uint `json:"langID"`
}
type GetLanguageResponse struct {
	Language book.Language `json:"language"`
}

type GetLanguagesRequest struct{}
type GetLanguagesResponse struct {
	Languages []book.Language `json:"languages"`
}

type DeleteLanguageRequest struct {
	LangID uint `json:"langID"`
}

type DeleteLanguageResponse struct{}

type AddBookRequest struct {
	Book book.Book `json:"book"`
}
type AddBookResponse struct {
	Book book.Book `json:"book"`
}

type SetBookDiscountRequest struct {
	BookID   uint `json:"bookID"`
	Digital  uint `json:"digital"`
	Physical uint `json:"physical"`
}
type SetBookDiscountResponse struct{}

type GetBookRequest struct {
	BookID uint `json:"bookID"`
}
type GetBookResponse struct {
	Book book.Book `json:"book"`
}

type EditBookRequest struct {
	Book book.Book `json:"book"`
}
type EditBookResponse struct {
	Book book.Book `json:"book"`
}

type GetAllBooksRequest struct{}
type GetAllBooksResponse struct {
	Books []book.Book `json:"books"`
}

type GetAuthorBooksRequest struct {
	AuthorID uint `json:"authorID"`
}
type GetAuthorBooksResponse struct {
	Books []book.Book `json:"books"`
}

type GetTopicBooksRequest struct {
	TopicID uint `json:"topicID"`
}
type GetTopicBooksResponse struct {
	Books []book.Book `json:"books"`
}

type GetPublisherBooksRequest struct {
	PublisherID uint `json:"publisherID"`
}
type GetPublisherBooksResponse struct {
	Books []book.Book `json:"books"`
}

type GetLangBooksRequest struct {
	LangID uint `json:"langID"`
}
type GetLangBooksResponse struct {
	Books []book.Book `json:"books"`
}

type DeleteBookRequest struct {
	BookID uint `json:"bookID"`
}
type DeleteBookResponse struct{}
