package v1

import (
	"net/http"
	"strconv"
	"strings"

	repository "github.com/XBozorg/bookstore/adapter/repository/mysql"
	"github.com/XBozorg/bookstore/dto"
	"github.com/XBozorg/bookstore/usecase/book"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func AddAuthor(repo repository.MySQLRepo, validator book.ValidateAddAuthor) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.AddAuthorRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(repo).AddAuthor(c.Request().Context(), req)
		if err != nil {

			if driverErr, ok := err.(*mysql.MySQLError); ok && driverErr.Number == 1062 {
				return echo.NewHTTPError(http.StatusConflict, "author already exists")
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetAuthor(repo repository.MySQLRepo, validator book.ValidateGetAuthor) echo.HandlerFunc {
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

		resp, err := book.New(repo).GetAuthor(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetAuthors(repo repository.MySQLRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetAuthorsRequest{}

		resp, err := book.New(repo).GetAuthors(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func DeleteAuthor(repo repository.MySQLRepo, validator book.ValidateDeleteAuthor) echo.HandlerFunc {
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

		resp, err := book.New(repo).DeleteAuthor(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func AddPublisher(repo repository.MySQLRepo, validator book.ValidateAddPublisher) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.AddPublisherRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(repo).AddPublisher(c.Request().Context(), req)
		if err != nil {

			if driverErr, ok := err.(*mysql.MySQLError); ok && driverErr.Number == 1062 {
				return echo.NewHTTPError(http.StatusConflict, "publisher already exists")
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetPublisher(repo repository.MySQLRepo, validator book.ValidateGetPublisher) echo.HandlerFunc {
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

		resp, err := book.New(repo).GetPublisher(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetPublishers(repo repository.MySQLRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetPublishersRequest{}

		resp, err := book.New(repo).GetPublishers(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func DeletePublisher(repo repository.MySQLRepo, validator book.ValidateDeletePublisher) echo.HandlerFunc {
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

		resp, err := book.New(repo).DeletePublisher(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func AddTopic(repo repository.MySQLRepo, validator book.ValidateAddTopic) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.AddTopicRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(repo).AddTopic(c.Request().Context(), req)
		if err != nil {

			if driverErr, ok := err.(*mysql.MySQLError); ok && driverErr.Number == 1062 {
				return echo.NewHTTPError(http.StatusConflict, "topic already exists")
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetTopic(repo repository.MySQLRepo, validator book.ValidateGetTopic) echo.HandlerFunc {
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

		resp, err := book.New(repo).GetTopic(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetTopics(repo repository.MySQLRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetTopicsRequest{}

		resp, err := book.New(repo).GetTopics(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func DeleteTopic(repo repository.MySQLRepo, validator book.ValidateDeleteTopic) echo.HandlerFunc {
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

		resp, err := book.New(repo).DeleteTopic(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func AddLanguage(repo repository.MySQLRepo, validator book.ValidateAddLanguage) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.AddLanguageRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(repo).AddLanguage(c.Request().Context(), req)
		if err != nil {

			if driverErr, ok := err.(*mysql.MySQLError); ok && driverErr.Number == 1062 {
				return echo.NewHTTPError(http.StatusConflict, "language already exists")
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetLanguage(repo repository.MySQLRepo, validator book.ValidateGetLanguage) echo.HandlerFunc {
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

		resp, err := book.New(repo).GetLanguage(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetLanguages(repo repository.MySQLRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetLanguagesRequest{}

		resp, err := book.New(repo).GetLanguages(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func DeleteLanguage(repo repository.MySQLRepo, validator book.ValidateDeleteLanguage) echo.HandlerFunc {
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

		resp, err := book.New(repo).DeleteLanguage(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func AddBook(repo repository.MySQLRepo, validator book.ValidateAddBook) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.AddBookRequest{}

		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := validator(c.Request().Context(), req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := book.New(repo).AddBook(c.Request().Context(), req)
		if err != nil {

			if driverErr, ok := err.(*mysql.MySQLError); ok && driverErr.Number == 1062 {
				return echo.NewHTTPError(http.StatusConflict, "book already exists")
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetBook(repo repository.MySQLRepo, validator book.ValidateGetBook) echo.HandlerFunc {
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

		resp, err := book.New(repo).GetBook(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func SetBookDiscount(repo repository.MySQLRepo, validator book.ValidateSetBookDiscount) echo.HandlerFunc {
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

		resp, err := book.New(repo).SetBookDiscount(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func EditBook(repo repository.MySQLRepo, validator book.ValidateEditBook) echo.HandlerFunc {
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

		resp, err := book.New(repo).EditBook(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetAllBooks(repo repository.MySQLRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := dto.GetAllBooksRequest{}

		resp, err := book.New(repo).GetAllBooks(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetAuthorBooks(repo repository.MySQLRepo, validator book.ValidateGetAuthorBooks) echo.HandlerFunc {
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

		resp, err := book.New(repo).GetAuthorBooks(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetPublisherBooks(repo repository.MySQLRepo, validator book.ValidateGetPublisherBooks) echo.HandlerFunc {
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

		resp, err := book.New(repo).GetPublisherBooks(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetTopicBooks(repo repository.MySQLRepo, validator book.ValidateGetTopicBooks) echo.HandlerFunc {
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

		resp, err := book.New(repo).GetTopicBooks(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetLangBooks(repo repository.MySQLRepo, validator book.ValidateGetLangBooks) echo.HandlerFunc {
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

		resp, err := book.New(repo).GetLangBooks(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func DeleteBook(repo repository.MySQLRepo, validator book.ValidateDeleteBook) echo.HandlerFunc {
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

		resp, err := book.New(repo).DeleteBook(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func GetUserDigitalBooks(repo repository.MySQLRepo, validator book.ValidateGetUserDigitalBooks) echo.HandlerFunc {
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

		resp, err := book.New(repo).GetUserDigitalBooks(c.Request().Context(), req)
		if err != nil {

			if strings.Contains(err.Error(), "no rows") {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func DownloadBook(repo repository.MySQLRepo, validator book.ValidateDownloadBook) echo.HandlerFunc {
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
