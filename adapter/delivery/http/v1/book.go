package v1

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/XBozorg/bookstore/adapter/repository"
	"github.com/XBozorg/bookstore/dto"
	"github.com/XBozorg/bookstore/usecase/book"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func AddAuthor(storage repository.Storage, validator book.ValidateAddAuthor) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.AddAuthorRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(storage).AddAuthor(c.Request().Context(), req)
		if err != nil {

			if driverErr, ok := err.(*mysql.MySQLError); ok && driverErr.Number == 1062 {
				return echo.NewHTTPError(http.StatusConflict, "author already exists")
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetAuthor(storage repository.Storage, validator book.ValidateGetAuthor) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetAuthorRequest{}

		aid, err := strconv.ParseUint(c.Param("authorID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		req.AuthorID = uint(aid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, "author does not exist")
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(storage).GetAuthor(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetAuthors(storage repository.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetAuthorsRequest{}

		resp, err := book.New(storage).GetAuthors(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func DeleteAuthor(storage repository.Storage, validator book.ValidateDeleteAuthor) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.DeleteAuthorRequest{}

		aid, err := strconv.ParseUint(c.Param("authorID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		req.AuthorID = uint(aid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, "author does not exist")
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(storage).DeleteAuthor(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func AddPublisher(storage repository.Storage, validator book.ValidateAddPublisher) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.AddPublisherRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(storage).AddPublisher(c.Request().Context(), req)
		if err != nil {

			if driverErr, ok := err.(*mysql.MySQLError); ok && driverErr.Number == 1062 {
				return echo.NewHTTPError(http.StatusConflict, "publisher already exists")
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetPublisher(storage repository.Storage, validator book.ValidateGetPublisher) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetPublisherRequest{}

		pid, err := strconv.ParseUint(c.Param("publisherID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		req.PublisherID = uint(pid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, "publisher does not exist")
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(storage).GetPublisher(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetPublishers(storage repository.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetPublishersRequest{}

		resp, err := book.New(storage).GetPublishers(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func DeletePublisher(storage repository.Storage, validator book.ValidateDeletePublisher) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.DeletePublisherRequest{}

		pid, err := strconv.ParseUint(c.Param("publisherID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		req.PublisherID = uint(pid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, "publisher does not exist")
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(storage).DeletePublisher(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func AddTopic(storage repository.Storage, validator book.ValidateAddTopic) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.AddTopicRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(storage).AddTopic(c.Request().Context(), req)
		if err != nil {

			if driverErr, ok := err.(*mysql.MySQLError); ok && driverErr.Number == 1062 {
				return echo.NewHTTPError(http.StatusConflict, "topic already exists")
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetTopic(storage repository.Storage, validator book.ValidateGetTopic) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetTopicRequest{}

		tid, err := strconv.ParseUint(c.Param("topicID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		req.TopicID = uint(tid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, "topic does not exist")
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(storage).GetTopic(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetTopics(storage repository.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetTopicsRequest{}

		resp, err := book.New(storage).GetTopics(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func DeleteTopic(storage repository.Storage, validator book.ValidateDeleteTopic) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.DeleteTopicRequest{}

		tid, err := strconv.ParseUint(c.Param("topicID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		req.TopicID = uint(tid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, "topic does not exist")
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(storage).DeleteTopic(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func AddLanguage(storage repository.Storage, validator book.ValidateAddLanguage) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.AddLanguageRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(storage).AddLanguage(c.Request().Context(), req)
		if err != nil {

			if driverErr, ok := err.(*mysql.MySQLError); ok && driverErr.Number == 1062 {
				return echo.NewHTTPError(http.StatusConflict, "language already exists")
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetLanguage(storage repository.Storage, validator book.ValidateGetLanguage) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetLanguageRequest{}

		lid, err := strconv.ParseUint(c.Param("langID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		req.LangID = uint(lid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, "language does not exist")
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(storage).GetLanguage(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetLanguages(storage repository.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetLanguagesRequest{}

		resp, err := book.New(storage).GetLanguages(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func DeleteLanguage(storage repository.Storage, validator book.ValidateDeleteLanguage) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.DeleteLanguageRequest{}

		lid, err := strconv.ParseUint(c.Param("langID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		req.LangID = uint(lid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, "language does not exist")
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(storage).DeleteLanguage(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func AddBook(storage repository.Storage, validator book.ValidateAddBook) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.AddBookRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(storage).AddBook(c.Request().Context(), req)
		if err != nil {

			if driverErr, ok := err.(*mysql.MySQLError); ok && driverErr.Number == 1062 {
				return echo.NewHTTPError(http.StatusConflict, "book already exists")
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetBook(storage repository.Storage, validator book.ValidateGetBook) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetBookRequest{}

		bid, err := strconv.ParseUint(c.Param("bookID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		req.BookID = uint(bid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, "book does not exist")
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(storage).GetBook(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func SetBookDiscount(storage repository.Storage, validator book.ValidateSetBookDiscount) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.SetBookDiscountRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		bid, err := strconv.ParseUint(c.Param("bookID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		req.BookID = uint(bid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, "book does not exist")
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(storage).SetBookDiscount(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func EditBook(storage repository.Storage, validator book.ValidateEditBook) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.EditBookRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		bid, err := strconv.ParseUint(c.Param("bookID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		req.Book.ID = uint(bid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, "book does not exist")
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(storage).EditBook(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetAllBooks(storage repository.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetAllBooksRequest{}

		resp, err := book.New(storage).GetAllBooks(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetAuthorBooks(storage repository.Storage, validator book.ValidateGetAuthorBooks) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetAuthorBooksRequest{}

		aid, err := strconv.ParseUint(c.Param("authorID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		req.AuthorID = uint(aid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, "author does not exist")
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(storage).GetAuthorBooks(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetPublisherBooks(storage repository.Storage, validator book.ValidateGetPublisherBooks) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetPublisherBooksRequest{}

		pid, err := strconv.ParseUint(c.Param("publisherID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		req.PublisherID = uint(pid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, "publisher does not exist")
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(storage).GetPublisherBooks(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetTopicBooks(storage repository.Storage, validator book.ValidateGetTopicBooks) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetTopicBooksRequest{}

		tid, err := strconv.ParseUint(c.Param("topicID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		req.TopicID = uint(tid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, "topic does not exist")
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(storage).GetTopicBooks(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetLangBooks(storage repository.Storage, validator book.ValidateGetLangBooks) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetLangBooksRequest{}

		lid, err := strconv.ParseUint(c.Param("langID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		req.LangID = uint(lid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, "language does not exist")
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(storage).GetLangBooks(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func DeleteBook(storage repository.Storage, validator book.ValidateDeleteBook) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.DeleteBookRequest{}

		bid, err := strconv.ParseUint(c.Param("bookID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		req.BookID = uint(bid)

		if err := validator(c.Request().Context(), req); err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				return echo.NewHTTPError(http.StatusNotFound, "book does not exist")
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(storage).DeleteBook(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetUserDigitalBooks(storage repository.Storage, validator book.ValidateGetUserDigitalBooks) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetUserDigitalBooksRequest{}

		userCookie, _ := c.Cookie("ID")
		req.UserID = userCookie.Value

		if err := validator(c.Request().Context(), req); err != nil {

			if err.Error() == "user does not exist" {
				return echo.NewHTTPError(http.StatusNotFound, "user does not exist")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(storage).GetUserDigitalBooks(c.Request().Context(), req)
		if err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func DownloadBook(storage repository.Storage, validator book.ValidateDownloadBook) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.DownloadBookRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		userCookie, _ := c.Cookie("ID")
		req.UserID = userCookie.Value

		bid, err := strconv.ParseUint(c.Param("bookID"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		req.BookID = uint(bid)

		if err := validator(c.Request().Context(), req); err != nil {

			if err.Error() == "book does not exist" {
				return echo.NewHTTPError(http.StatusNotFound, "book does not exist")
			}

			if err.Error() == "access denied" {
				return echo.NewHTTPError(http.StatusForbidden, "you don't have access to this book")
			}

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.File(req.Path)
	}
}
