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
	GetPublishers(ctx context.Context, req dto.GetPublishersRequest) (dto.GetPublishersResponse, error)
	DeletePublisher(ctx context.Context, req dto.DeletePublisherRequest) (dto.DeletePublisherResponse, error)

	AddTopic(ctx context.Context, req dto.AddTopicRequest) (dto.AddTopicResponse, error)
	GetTopic(ctx context.Context, req dto.GetTopicRequest) (dto.GetTopicResponse, error)
	GetTopics(ctx context.Context, req dto.GetTopicsRequest) (dto.GetTopicsResponse, error)
	DeleteTopic(ctx context.Context, req dto.DeleteTopicRequest) (dto.DeleteTopicResponse, error)

	AddLanguage(ctx context.Context, req dto.AddLanguageRequest) (dto.AddLanguageResponse, error)
	GetLanguage(ctx context.Context, req dto.GetLanguageRequest) (dto.GetLanguageResponse, error)
	GetLanguages(ctx context.Context, req dto.GetLanguagesRequest) (dto.GetLanguagesResponse, error)
	DeleteLanguage(ctx context.Context, req dto.DeleteLanguageRequest) (dto.DeleteLanguageResponse, error)

	AddBook(ctx context.Context, req dto.AddBookRequest) (dto.AddBookResponse, error)
	SetBookDiscount(ctx context.Context, req dto.SetBookDiscountRequest) (dto.SetBookDiscountResponse, error)
	GetBook(ctx context.Context, req dto.GetBookRequest) (dto.GetBookResponse, error)
	EditBook(ctx context.Context, req dto.EditBookRequest) (dto.EditBookResponse, error)
	GetAllBooks(ctx context.Context, req dto.GetAllBooksRequest) (dto.GetAllBooksResponse, error)
	GetAuthorBooks(ctx context.Context, req dto.GetAuthorBooksRequest) (dto.GetAuthorBooksResponse, error)
	GetPublisherBooks(ctx context.Context, req dto.GetPublisherBooksRequest) (dto.GetPublisherBooksResponse, error)
	GetTopicBooks(ctx context.Context, req dto.GetTopicBooksRequest) (dto.GetTopicBooksResponse, error)
	GetLangBooks(ctx context.Context, req dto.GetLangBooksRequest) (dto.GetLangBooksResponse, error)
	DeleteBook(ctx context.Context, req dto.DeleteBookRequest) (dto.DeleteBookResponse, error)

	GetUserDigitalBooks(ctx context.Context, req dto.GetUserDigitalBooksRequest) (dto.GetUserDigitalBooksResponse, error)
	DownloadBook(ctx context.Context, req dto.DownloadBookRequest) (dto.DownloadBookResponse, error)
}

type UseCaseRepo struct {
	repo Repository
}

func New(r Repository) UseCaseRepo {
	return UseCaseRepo{repo: r}
}

func (u UseCaseRepo) AddAuthor(ctx context.Context, req dto.AddAuthorRequest) (dto.AddAuthorResponse, error) {

	author, err := u.repo.AddAuthor(ctx, req.Name)
	if err != nil {
		return dto.AddAuthorResponse{}, err
	}

	return dto.AddAuthorResponse{Author: author}, nil
}

func (u UseCaseRepo) GetAuthor(ctx context.Context, req dto.GetAuthorRequest) (dto.GetAuthorResponse, error) {

	author, err := u.repo.GetAuthor(ctx, req.AuthorID)
	if err != nil {
		return dto.GetAuthorResponse{}, err
	}

	return dto.GetAuthorResponse{Author: author}, nil
}

func (u UseCaseRepo) GetAuthors(ctx context.Context, req dto.GetAuthorsRequest) (dto.GetAuthorsResponse, error) {

	authors, err := u.repo.GetAuthors(ctx)
	if err != nil {
		return dto.GetAuthorsResponse{}, err
	}

	return dto.GetAuthorsResponse{Authors: authors}, nil
}

func (u UseCaseRepo) DeleteAuthor(ctx context.Context, req dto.DeleteAuthorRequest) (dto.DeleteAuthorResponse, error) {

	err := u.repo.DeleteAuthor(ctx, req.AuthorID)
	if err != nil {
		return dto.DeleteAuthorResponse{}, err
	}

	return dto.DeleteAuthorResponse{}, nil
}

func (u UseCaseRepo) AddPublisher(ctx context.Context, req dto.AddPublisherRequest) (dto.AddPublisherResponse, error) {

	publisher, err := u.repo.AddPublisher(ctx, req.Name)
	if err != nil {
		return dto.AddPublisherResponse{}, err
	}

	return dto.AddPublisherResponse{Publisher: publisher}, nil
}

func (u UseCaseRepo) GetPublisher(ctx context.Context, req dto.GetPublisherRequest) (dto.GetPublisherResponse, error) {

	publisher, err := u.repo.GetPublisher(ctx, req.PublisherID)
	if err != nil {
		return dto.GetPublisherResponse{}, err
	}

	return dto.GetPublisherResponse{Publisher: publisher}, nil
}

func (u UseCaseRepo) GetPublishers(ctx context.Context, req dto.GetPublishersRequest) (dto.GetPublishersResponse, error) {

	publishers, err := u.repo.GetPublishers(ctx)
	if err != nil {
		return dto.GetPublishersResponse{}, err
	}

	return dto.GetPublishersResponse{Publishers: publishers}, nil
}

func (u UseCaseRepo) DeletePublisher(ctx context.Context, req dto.DeletePublisherRequest) (dto.DeletePublisherResponse, error) {

	err := u.repo.DeletePublisher(ctx, req.PublisherID)
	if err != nil {
		return dto.DeletePublisherResponse{}, err
	}

	return dto.DeletePublisherResponse{}, nil
}

func (u UseCaseRepo) AddTopic(ctx context.Context, req dto.AddTopicRequest) (dto.AddTopicResponse, error) {

	topic, err := u.repo.AddTopic(ctx, req.Name)
	if err != nil {
		return dto.AddTopicResponse{}, err
	}

	return dto.AddTopicResponse{Topic: topic}, nil
}

func (u UseCaseRepo) GetTopic(ctx context.Context, req dto.GetTopicRequest) (dto.GetTopicResponse, error) {

	topic, err := u.repo.GetTopic(ctx, req.TopicID)
	if err != nil {
		return dto.GetTopicResponse{}, err
	}

	return dto.GetTopicResponse{Topic: topic}, nil
}

func (u UseCaseRepo) GetTopics(ctx context.Context, req dto.GetTopicsRequest) (dto.GetTopicsResponse, error) {

	topics, err := u.repo.GetTopics(ctx)
	if err != nil {
		return dto.GetTopicsResponse{}, err
	}

	return dto.GetTopicsResponse{Topics: topics}, nil
}

func (u UseCaseRepo) DeleteTopic(ctx context.Context, req dto.DeleteTopicRequest) (dto.DeleteTopicResponse, error) {

	err := u.repo.DeleteTopic(ctx, req.TopicID)
	if err != nil {
		return dto.DeleteTopicResponse{}, err
	}

	return dto.DeleteTopicResponse{}, nil
}

func (u UseCaseRepo) AddLanguage(ctx context.Context, req dto.AddLanguageRequest) (dto.AddLanguageResponse, error) {

	lang, err := u.repo.AddLanguage(ctx, req.LangCode)
	if err != nil {
		return dto.AddLanguageResponse{}, err
	}

	return dto.AddLanguageResponse{Language: lang}, nil
}

func (u UseCaseRepo) GetLanguage(ctx context.Context, req dto.GetLanguageRequest) (dto.GetLanguageResponse, error) {

	lang, err := u.repo.GetLanguage(ctx, req.LangID)
	if err != nil {
		return dto.GetLanguageResponse{}, err
	}

	return dto.GetLanguageResponse{Language: lang}, nil
}

func (u UseCaseRepo) GetLanguages(ctx context.Context, req dto.GetLanguagesRequest) (dto.GetLanguagesResponse, error) {

	langs, err := u.repo.GetLanguages(ctx)
	if err != nil {
		return dto.GetLanguagesResponse{}, err
	}

	return dto.GetLanguagesResponse{Languages: langs}, nil
}

func (u UseCaseRepo) DeleteLanguage(ctx context.Context, req dto.DeleteLanguageRequest) (dto.DeleteLanguageResponse, error) {

	err := u.repo.DeleteLanguage(ctx, req.LangID)
	if err != nil {
		return dto.DeleteLanguageResponse{}, err
	}

	return dto.DeleteLanguageResponse{}, nil
}

func (u UseCaseRepo) AddBook(ctx context.Context, req dto.AddBookRequest) (dto.AddBookResponse, error) {

	book, err := u.repo.AddBook(ctx, req.Book)
	if err != nil {
		return dto.AddBookResponse{}, err
	}

	return dto.AddBookResponse{Book: book}, nil
}

func (u UseCaseRepo) SetBookDiscount(ctx context.Context, req dto.SetBookDiscountRequest) (dto.SetBookDiscountResponse, error) {

	err := u.repo.SetBookDiscount(ctx, req.BookID, req.Digital, req.Physical)
	if err != nil {
		return dto.SetBookDiscountResponse{}, err
	}

	return dto.SetBookDiscountResponse{}, nil
}

func (u UseCaseRepo) GetBook(ctx context.Context, req dto.GetBookRequest) (dto.GetBookResponse, error) {

	book, err := u.repo.GetBook(ctx, req.BookID)
	if err != nil {
		return dto.GetBookResponse{}, err
	}

	return dto.GetBookResponse{Book: book}, nil
}

func (u UseCaseRepo) EditBook(ctx context.Context, req dto.EditBookRequest) (dto.EditBookResponse, error) {

	book, err := u.repo.EditBook(ctx, req.Book)
	if err != nil {
		return dto.EditBookResponse{}, nil
	}

	return dto.EditBookResponse{Book: book}, nil
}

func (u UseCaseRepo) GetAllBooksFull(ctx context.Context, req dto.GetAllBooksRequest) (dto.GetAllBooksResponse, error) {

	books, err := u.repo.GetAllBooksFull(ctx)
	if err != nil {
		return dto.GetAllBooksResponse{}, err
	}

	return dto.GetAllBooksResponse{Books: books}, nil
}

func (u UseCaseRepo) GetAllBooks(ctx context.Context, req dto.GetAllBooksRequest) (dto.GetAllBooksResponse, error) {

	books, err := u.repo.GetAllBooks(ctx)
	if err != nil {
		return dto.GetAllBooksResponse{}, err
	}

	return dto.GetAllBooksResponse{Books: books}, nil
}

func (u UseCaseRepo) GetAuthorBooks(ctx context.Context, req dto.GetAuthorBooksRequest) (dto.GetAuthorBooksResponse, error) {

	books, err := u.repo.GetAuthorBooks(ctx, req.AuthorID)
	if err != nil {
		return dto.GetAuthorBooksResponse{}, err
	}

	return dto.GetAuthorBooksResponse{Books: books}, nil
}

func (u UseCaseRepo) GetTopicBooks(ctx context.Context, req dto.GetTopicBooksRequest) (dto.GetTopicBooksResponse, error) {

	books, err := u.repo.GetTopicBooks(ctx, req.TopicID)
	if err != nil {
		return dto.GetTopicBooksResponse{}, err
	}

	return dto.GetTopicBooksResponse{Books: books}, nil
}

func (u UseCaseRepo) GetPublisherBooks(ctx context.Context, req dto.GetPublisherBooksRequest) (dto.GetPublisherBooksResponse, error) {

	books, err := u.repo.GetPublisherBooks(ctx, req.PublisherID)
	if err != nil {
		return dto.GetPublisherBooksResponse{}, err
	}

	return dto.GetPublisherBooksResponse{Books: books}, nil
}

func (u UseCaseRepo) GetLangBooks(ctx context.Context, req dto.GetLangBooksRequest) (dto.GetLangBooksResponse, error) {

	books, err := u.repo.GetLangBooks(ctx, req.LangID)
	if err != nil {
		return dto.GetLangBooksResponse{}, err
	}
	return dto.GetLangBooksResponse{Books: books}, nil
}

func (u UseCaseRepo) DeleteBook(ctx context.Context, req dto.DeleteBookRequest) (dto.DeleteBookResponse, error) {

	err := u.repo.DeleteBook(ctx, req.BookID)
	if err != nil {
		return dto.DeleteBookResponse{}, err
	}

	return dto.DeleteBookResponse{}, nil
}

func (u UseCaseRepo) GetUserDigitalBooks(ctx context.Context, req dto.GetUserDigitalBooksRequest) (dto.GetUserDigitalBooksResponse, error) {

	books, err := u.repo.GetUserDigitalBooks(ctx, req.UserID)
	if err != nil {
		return dto.GetUserDigitalBooksResponse{}, err
	}

	return dto.GetUserDigitalBooksResponse{Books: books}, nil
}
