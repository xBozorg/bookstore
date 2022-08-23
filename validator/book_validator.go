package validator

import (
	"context"
	"errors"

	"github.com/XBozorg/bookstore/adapter/repository"
	"github.com/XBozorg/bookstore/dto"
	"github.com/XBozorg/bookstore/usecase/book"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func doesAuthorExist(ctx context.Context, repo book.ValidatorRepo) validation.RuleFunc {
	return func(value interface{}) error {
		authorID := value.(uint)

		ok, err := repo.DoesAuthorExist(ctx, authorID)
		if err != nil {
			return err
		}

		if !ok {
			return errors.New("author does not exist")
		}
		return nil
	}
}

func doesBookExist(ctx context.Context, repo book.ValidatorRepo) validation.RuleFunc {
	return func(value interface{}) error {
		bookID := value.(uint)

		ok, err := repo.DoesBookExist(ctx, bookID)
		if err != nil {
			return err
		}

		if !ok {
			return errors.New("book does not exist")
		}
		return nil
	}
}

func doesPublisherExist(ctx context.Context, repo book.ValidatorRepo) validation.RuleFunc {
	return func(value interface{}) error {
		publisherID := value.(uint)

		ok, err := repo.DoesPublisherExist(ctx, publisherID)
		if err != nil {
			return err
		}

		if !ok {
			return errors.New("publisher does not exist")
		}
		return nil
	}
}

func doesTopicExist(ctx context.Context, repo book.ValidatorRepo) validation.RuleFunc {
	return func(value interface{}) error {
		topicID := value.(uint)

		ok, err := repo.DoesTopicExist(ctx, topicID)
		if err != nil {
			return err
		}

		if !ok {
			return errors.New("topic does not exist")
		}
		return nil
	}
}

func doesLangExist(ctx context.Context, repo book.ValidatorRepo) validation.RuleFunc {
	return func(value interface{}) error {
		langID := value.(uint)

		ok, err := repo.DoesLanguageExist(ctx, langID)
		if err != nil {
			return err
		}

		if !ok {
			return errors.New("language does not exist")
		}
		return nil
	}
}

func doesUserAccessBook(ctx context.Context, repo book.ValidatorRepo, userID string) validation.RuleFunc {
	return func(value interface{}) error {
		bookID := value.(uint)

		ok, err := repo.DoesUserAccessBook(ctx, userID, bookID)
		if err != nil {
			return err
		}

		if !ok {
			return errors.New("access denied")
		}

		return nil
	}
}

func ValidateAddAuthor(storage repository.Storage) book.ValidateAddAuthor {
	return func(ctx context.Context, req dto.AddAuthorRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.Name, validation.Required, is.ASCII, validation.Length(4, 100)),
		)
	}
}

func ValidateGetAuthor(storage repository.Storage) book.ValidateGetAuthor {
	return func(ctx context.Context, req dto.GetAuthorRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.AuthorID, validation.Required, validation.By(doesAuthorExist(ctx, storage))),
		)
	}
}

func ValidateDeleteAuthor(storage repository.Storage) book.ValidateDeleteAuthor {
	return func(ctx context.Context, req dto.DeleteAuthorRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.AuthorID, validation.Required, validation.By(doesAuthorExist(ctx, storage))),
		)
	}
}

func ValidateAddPublisher(storage repository.Storage) book.ValidateAddPublisher {
	return func(ctx context.Context, req dto.AddPublisherRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.Name, validation.Required, is.ASCII, validation.Length(1, 100)),
		)
	}
}

func ValidateGetPublisher(storage repository.Storage) book.ValidateGetPublisher {
	return func(ctx context.Context, req dto.GetPublisherRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.PublisherID, validation.Required, validation.By(doesPublisherExist(ctx, storage))),
		)
	}
}

func ValidateDeletePublisher(storage repository.Storage) book.ValidateDeletePublisher {
	return func(ctx context.Context, req dto.DeletePublisherRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.PublisherID, validation.Required, validation.By(doesPublisherExist(ctx, storage))),
		)
	}
}

func ValidateAddTopic(storage repository.Storage) book.ValidateAddTopic {
	return func(ctx context.Context, req dto.AddTopicRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.Name, validation.Required, is.Alpha, validation.Length(2, 30)),
		)
	}
}

func ValidateGetTopic(storage repository.Storage) book.ValidateGetTopic {
	return func(ctx context.Context, req dto.GetTopicRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.TopicID, validation.Required, validation.By(doesTopicExist(ctx, storage))),
		)
	}
}

func ValidateDeleteTopic(storage repository.Storage) book.ValidateDeleteTopic {
	return func(ctx context.Context, req dto.DeleteTopicRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.TopicID, validation.Required, validation.By(doesTopicExist(ctx, storage))),
		)
	}
}

func ValidateAddLanguage(storage repository.Storage) book.ValidateAddLanguage {
	return func(ctx context.Context, req dto.AddLanguageRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.LangCode, validation.Required, is.Alpha, validation.Length(2, 2)),
		)
	}
}

func ValidateGetLanguage(storage repository.Storage) book.ValidateGetLanguage {
	return func(ctx context.Context, req dto.GetLanguageRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.LangID, validation.Required, validation.By(doesLangExist(ctx, storage))),
		)
	}
}

func ValidateDeleteLanguage(storage repository.Storage) book.ValidateDeleteLanguage {
	return func(ctx context.Context, req dto.DeleteLanguageRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.LangID, validation.Required, validation.By(doesLangExist(ctx, storage))),
		)
	}
}

func ValidateAddBook(storage repository.Storage) book.ValidateAddBook {
	return func(ctx context.Context, req dto.AddBookRequest) error {
		if errBook := validation.ValidateStruct(&req.Book,
			validation.Field(&req.Book.Title, validation.Required, is.ASCII, validation.Length(1, 100)),
			validation.Field(&req.Book.ISBN, validation.Required, is.ISBN13),
			validation.Field(&req.Book.Pages, validation.Required),
			validation.Field(&req.Book.Description, is.ASCII, validation.Length(0, 500)),
			validation.Field(&req.Book.Year, validation.Required, validation.Date("2006")),
			validation.Field(&req.Book.CreationDate, validation.Required, validation.Date("2006-01-02 15:04:05")),
			validation.Field(&req.Book.CoverFront, validation.Required, is.ASCII, validation.Length(10, 150)),
			validation.Field(&req.Book.CoverBack, validation.Required, is.ASCII, validation.Length(10, 150)),
		); errBook != nil {
			return errBook
		}

		if errDigital := validation.ValidateStruct(&req.Book.Digital,
			validation.Field(&req.Book.Digital.Price, validation.Required),
			validation.Field(&req.Book.Digital.Discount, validation.Max(100)),
			validation.Field(&req.Book.Digital.PDF, is.ASCII, validation.Length(10, 150)),
			validation.Field(&req.Book.Digital.EPUB, is.ASCII, validation.Length(10, 150)),
			validation.Field(&req.Book.Digital.DJVU, is.ASCII, validation.Length(10, 150)),
			validation.Field(&req.Book.Digital.AZW, is.ASCII, validation.Length(10, 150)),
			validation.Field(&req.Book.Digital.TXT, is.ASCII, validation.Length(10, 150)),
			validation.Field(&req.Book.Digital.DOCX, is.ASCII, validation.Length(10, 150)),
		); errDigital != nil {
			return errDigital
		}

		if errPhysical := validation.ValidateStruct(&req.Book.Physical,
			validation.Field(&req.Book.Physical.Price, validation.Required),
			validation.Field(&req.Book.Physical.Discount, validation.Max(100)),
			validation.Field(&req.Book.Physical.Stock, validation.Required),
		); errPhysical != nil {
			return errPhysical
		}

		if errLang := validation.ValidateStruct(&req.Book.Language,
			validation.Field(&req.Book.Language.ID, validation.Required),
		); errLang != nil {
			return errLang
		}

		if errPub := validation.ValidateStruct(&req.Book.Publisher,
			validation.Field(&req.Book.Publisher.ID, validation.Required),
		); errPub != nil {
			return errPub
		}

		return nil
	}
}

func ValidateGetBook(storage repository.Storage) book.ValidateGetBook {
	return func(ctx context.Context, req dto.GetBookRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.BookID, validation.Required, validation.By(doesBookExist(ctx, storage))),
		)
	}
}

func ValidateEditBook(storage repository.Storage) book.ValidateEditBook {
	return func(ctx context.Context, req dto.EditBookRequest) error {
		if errBook := validation.ValidateStruct(&req.Book,
			validation.Field(&req.Book.ID, validation.Required, validation.By(doesBookExist(ctx, storage))),
			validation.Field(&req.Book.Title, validation.Required, is.ASCII, validation.Length(1, 100)),
			validation.Field(&req.Book.ISBN, validation.Required, is.ISBN13),
			validation.Field(&req.Book.Pages, validation.Required),
			validation.Field(&req.Book.Description, is.ASCII, validation.Length(0, 500)),
			validation.Field(&req.Book.Year, validation.Required, validation.Date("2006")),
			validation.Field(&req.Book.CreationDate, validation.Required, validation.Date("2006-01-02 15:04:05")),
			validation.Field(&req.Book.CoverFront, validation.Required, is.ASCII, validation.Length(10, 150)),
			validation.Field(&req.Book.CoverBack, validation.Required, is.ASCII, validation.Length(10, 150)),
		); errBook != nil {
			return errBook
		}

		if errDigital := validation.ValidateStruct(&req.Book.Digital,
			validation.Field(&req.Book.Digital.Price, validation.Required),
			validation.Field(&req.Book.Digital.Discount, validation.Max(100)),
			validation.Field(&req.Book.Digital.PDF, is.ASCII, validation.Length(10, 150)),
			validation.Field(&req.Book.Digital.EPUB, is.ASCII, validation.Length(10, 150)),
			validation.Field(&req.Book.Digital.DJVU, is.ASCII, validation.Length(10, 150)),
			validation.Field(&req.Book.Digital.AZW, is.ASCII, validation.Length(10, 150)),
			validation.Field(&req.Book.Digital.TXT, is.ASCII, validation.Length(10, 150)),
			validation.Field(&req.Book.Digital.DOCX, is.ASCII, validation.Length(10, 150)),
		); errDigital != nil {
			return errDigital
		}

		if errPhysical := validation.ValidateStruct(&req.Book.Physical,
			validation.Field(&req.Book.Physical.Price, validation.Required),
			validation.Field(&req.Book.Physical.Discount, validation.Max(100)),
			validation.Field(&req.Book.Physical.Stock, validation.Required),
		); errPhysical != nil {
			return errPhysical
		}

		if errLang := validation.ValidateStruct(&req.Book.Language,
			validation.Field(&req.Book.Language.ID, validation.Required),
		); errLang != nil {
			return errLang
		}

		if errPub := validation.ValidateStruct(&req.Book.Publisher,
			validation.Field(&req.Book.Publisher.ID, validation.Required),
		); errPub != nil {
			return errPub
		}

		return nil
	}
}

func ValidateSetBookDiscount(storage repository.Storage) book.ValidateSetBookDiscount {
	return func(ctx context.Context, req dto.SetBookDiscountRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.BookID, validation.Required, validation.By(doesBookExist(ctx, storage))),
			validation.Field(&req.Digital, validation.Max(uint(100))),
			validation.Field(&req.Physical, validation.Max(uint(100))),
		)
	}
}

func ValidateGetAuthorBooks(storage repository.Storage) book.ValidateGetAuthorBooks {
	return func(ctx context.Context, req dto.GetAuthorBooksRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.AuthorID, validation.Required, validation.By(doesAuthorExist(ctx, storage))),
		)
	}
}

func ValidateGetPublisherBooks(storage repository.Storage) book.ValidateGetPublisherBooks {
	return func(ctx context.Context, req dto.GetPublisherBooksRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.PublisherID, validation.Required, validation.By(doesPublisherExist(ctx, storage))),
		)
	}
}

func ValidateGetTopicBooks(storage repository.Storage) book.ValidateGetTopicBooks {
	return func(ctx context.Context, req dto.GetTopicBooksRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.TopicID, validation.Required, validation.By(doesTopicExist(ctx, storage))),
		)
	}
}

func ValidateGetLangBooks(storage repository.Storage) book.ValidateGetLangBooks {
	return func(ctx context.Context, req dto.GetLangBooksRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.LangID, validation.Required, validation.By(doesLangExist(ctx, storage))),
		)
	}
}

func ValidateDeleteBook(storage repository.Storage) book.ValidateDeleteBook {
	return func(ctx context.Context, req dto.DeleteBookRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.BookID, validation.Required, validation.By(doesBookExist(ctx, storage))),
		)
	}
}

func ValidateGetUserDigitalBooks(storage repository.Storage) book.ValidateGetUserDigitalBooks {
	return func(ctx context.Context, req dto.GetUserDigitalBooksRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.UserID, is.UUIDv4, validation.By(doesUserExist(ctx, storage))),
		)
	}
}

func ValidateDownloadBook(storage repository.Storage) book.ValidateDownloadBook {
	return func(ctx context.Context, req dto.DownloadBookRequest) error {
		return validation.ValidateStruct(&req,
			validation.Field(&req.BookID, validation.Required,
				validation.By(doesBookExist(ctx, storage)),
				validation.By(doesUserAccessBook(ctx, storage, req.UserID))),

			validation.Field(&req.Path, is.ASCII, validation.Length(10, 150)),
		)
	}
}
