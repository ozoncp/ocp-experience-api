package repo

import (
	"context"
	"errors"
	"time"

	"database/sql"
	"database/sql/driver"

	"github.com/DATA-DOG/go-sqlmock"

	sq "github.com/Masterminds/squirrel"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ozoncp/ocp-experience-api/internal/models"
)

var _ = Describe("IRepo", func() {
	var (
		rep    IRepo
		dbMock sqlmock.Sqlmock
		mockCtrl *gomock.Controller
		ctx      context.Context
		db       *sql.DB
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		ctx = context.Background()
	})

	AfterEach(func() {
		mockCtrl.Finish()
		defer db.Close()

		if err := dbMock.ExpectationsWereMet(); err != nil {
			Expect(err).ToNot(HaveOccurred())
		}
	})

	Context("Adding items with no errors. Will not return any remains.", func() {
		JustBeforeEach(func() {
			var err error
			db, dbMock, err = sqlmock.New()

			Expect(err).ToNot(HaveOccurred())

			cache := sq.NewStmtCache(db)
			rep = &Repo{
				builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(cache),
			}
		})

		It("Add experience. Expect new ID generated.", func() {
			newExperience := models.NewExperience(1, 1, 1, time.Time{}, time.Time{}, 1)
			expectedNewId := uint64(1)
			returnRows := sqlmock.NewRows([]string{"id"}).AddRow(expectedNewId)

			dbMock.ExpectPrepare(
				"INSERT INTO experiences \\(user_id,type,from,to,level\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5\\) RETURNING id",
			).
				ExpectQuery().
				WithArgs(newExperience.UserId, newExperience.Type, newExperience.From, newExperience.To, newExperience.Level).
				WillReturnRows(returnRows)

			newId, err := rep.Add(ctx, newExperience)

			Expect(err).ToNot(HaveOccurred())
			Expect(newId).To(Equal(expectedNewId))
		})

		It("Add many requests into repository", func() {
			experiences := []models.Experience{
				models.NewExperience(1, 1, 1, time.Time{}, time.Time{}, 1),
				models.NewExperience(2, 2, 2, time.Time{}, time.Time{}, 2),
				models.NewExperience(3, 3, 3, time.Time{}, time.Time{}, 3),
			}

			expectedQueryArgs := make([]driver.Value, 0, len(experiences) * 3)

			for _, req := range experiences {
				expectedQueryArgs = append(expectedQueryArgs, req.UserId, req.Type, req.From, req.To, req.Level)
			}

			expectedIds := []uint64{1, 2, 3}

			dbMock.ExpectPrepare(
				"INSERT INTO experiences \\(user_id,type,from,to,level\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5\\),\\(\\$6,\\$7,\\$8,\\$9,\\$10\\),\\(\\$11,\\$12,\\$13,\\$14,\\$15\\) RETURNING id",
			).
				ExpectQuery().
				WithArgs(expectedQueryArgs...).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).
					AddRow(1).
					AddRow(2).
					AddRow(3))

			newIds, err := rep.AddExperiences(ctx, experiences)

			Expect(err).ToNot(HaveOccurred())
			Expect(newIds).To(Equal(expectedIds))
		})

		It("Fetch experiences from database", func() {
			dbRows := [][]driver.Value{
				{uint64(1), uint64(1), uint64(1), time.Time{}, time.Time{}, uint64(1)},
				{uint64(2), uint64(2), uint64(2), time.Time{}, time.Time{}, uint64(2)},
				{uint64(3), uint64(3), uint64(3), time.Time{}, time.Time{}, uint64(3)},
			}

			expectedExperiences := make([]models.Experience, 0, len(dbRows))
			returnRows := sqlmock.NewRows([]string{"id", "user_id", "type", "from", "to", "level"})

			for _, row := range dbRows {
				expectedExperiences = append(expectedExperiences, models.Experience{
					Id:     row[0].(uint64),
					UserId: row[1].(uint64),
					Type:   row[2].(uint64),
					From:   row[3].(time.Time),
					To:		row[4].(time.Time),
					Level: 	row[5].(uint64),
				})

				returnRows.AddRow(row...)
			}

			offset, limit := uint64(100), uint64(1000)

			dbMock.ExpectPrepare(
				"SELECT id, user_id, type, from, to, level FROM experiences LIMIT 1000 OFFSET 100",
			).
				ExpectQuery().
				WillReturnRows(returnRows)

			actualExperiences, err := rep.List(ctx, limit, offset)

			Expect(err).ToNot(HaveOccurred())
			Expect(actualExperiences).To(Equal(expectedExperiences))
		})

		It("Remove experience that exists", func() {
			id := uint64(100)
			res := sqlmock.NewResult(0, 1)

			dbMock.ExpectPrepare(
				"DELETE FROM experiences WHERE id = \\$1",
			).
				ExpectExec().
				WithArgs(id).
				WillReturnResult(res)

			found, err := rep.Remove(ctx, id)
			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(Equal(true))
		})

		It("Remove experience that does not exist", func() {
			id := uint64(100)
			res := sqlmock.NewResult(0, 0)

			dbMock.ExpectPrepare(
				"DELETE FROM experiences WHERE id = \\$1",
			).
				ExpectExec().
				WithArgs(id).
				WillReturnResult(res)

			found, err := rep.Remove(ctx, id)

			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(Equal(false))
		})

		It("Return experiences that exists", func() {
			id := uint64(1)
			experience := models.NewExperience(id, 1, 1, time.Time{}, time.Time{}, 1)

			returnRows := sqlmock.
				NewRows([]string{"id", "user_id", "type", "from", "to", "level"}).
				AddRow(experience.Id, experience.UserId, experience.Type, experience.From, experience.To, experience.Level)

			dbMock.ExpectPrepare(
				"SELECT id, user_id, type, from, to, level FROM experiences WHERE id = \\$1",
			).
				ExpectQuery().
				WithArgs(id).
				WillReturnRows(returnRows)

			actualExperience, err := rep.Describe(ctx, id)

			Expect(err).ToNot(HaveOccurred())
			Expect(actualExperience).To(Equal(experience))
		})

		It("Return experience that doest not exist", func() {
			id := uint64(1)

			returnRows := sqlmock.
				NewRows([]string{"id", "user_id", "type", "from", "to", "level"})

			dbMock.ExpectPrepare(
				"SELECT id, user_id, type, from, to, level FROM experiences WHERE id = \\$1",
			).
				ExpectQuery().
				WithArgs(id).
				WillReturnRows(returnRows)

			experience, err := rep.Describe(ctx, id)

			Expect(err).To(Equal(NotFound))
			Expect(experience).To(Equal(models.Experience{}))
		})

		It("Pops up error on List", func() {
			offset, limit := uint64(100), uint64(1000)
			expectedError := errors.New("test error")

			dbMock.ExpectPrepare(
				"SELECT id, user_id, type, from, to, level FROM experiences LIMIT 1000 OFFSET 100",
			).
				ExpectQuery().
				WillReturnError(expectedError)

			_, err := rep.List(ctx, limit, offset)
			Expect(err).To(Equal(expectedError))
		})

		It("Pops up error on Add", func() {
			experience := models.NewExperience(1, 1, 1, time.Time{}, time.Time{}, 1)
			expectedError := errors.New("test error")

			dbMock.ExpectPrepare(
				"INSERT INTO experiences \\(user_id,type,from,to,level\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5\\) RETURNING id",
			).
				ExpectQuery().
				WithArgs(experience.UserId, experience.Type, experience.From, experience.To, experience.Level).
				WillReturnError(expectedError)

			newId, err := rep.Add(ctx, experience)

			Expect(err).To(Equal(expectedError))
			Expect(newId).To(Equal(uint64(0)))
		})

		It("Pops up error on AddExperiences", func() {
			newReq := models.NewExperience(1, 1, 1, time.Time{}, time.Time{}, 1)
			expectedError := errors.New("test error")
			dbMock.ExpectPrepare(
				"INSERT INTO experiences \\(user_id,type,from,to,level\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5\\) RETURNING id",
			).
				ExpectQuery().
				WithArgs(newReq.UserId, newReq.Type, newReq.From, newReq.To, newReq.Level).
				WillReturnError(expectedError)

			_, err := rep.AddExperiences(ctx, []models.Experience{newReq})
			Expect(err).To(Equal(expectedError))
		})

		It("Pops up error on Remove()", func() {
			id := uint64(100)
			expectedError := errors.New("test error")

			dbMock.ExpectPrepare(
				"DELETE FROM experiences WHERE id = \\$1",
			).
				ExpectExec().
				WithArgs(id).
				WillReturnError(expectedError)

			found, err := rep.Remove(ctx, id)

			Expect(err).To(Equal(expectedError))
			Expect(found).To(Equal(false))
		})

		It("Update experience that is exists", func() {
			experience := models.NewExperience(1, 1, 1, time.Now(), time.Now(), 1)
			res := sqlmock.NewResult(0, 1)

			dbMock.ExpectPrepare(
				"UPDATE experiences SET user_id = \\$1, type = \\$2, from = \\$3, to = \\$4, level = \\$5 WHERE id = \\$6",
			).
				ExpectExec().
				WithArgs(experience.UserId, experience.Type, experience.From, experience.To, experience.Level, experience.Id).
				WillReturnResult(res)

			err := rep.Update(ctx, experience)
			Expect(err).ToNot(Equal(NotFound))
		})

		It("Update experience that is not exists", func() {
			experience := models.NewExperience(1, 1, 1, time.Now(), time.Now(), 1)
			res := sqlmock.NewResult(0, 0)

			dbMock.ExpectPrepare(
				"UPDATE experiences SET user_id = \\$1, type = \\$2, from = \\$3, to = \\$4, level = \\$5 WHERE id = \\$6",
			).
				ExpectExec().
				WithArgs(experience.UserId, experience.Type, experience.From, experience.To, experience.Level, experience.Id).
				WillReturnResult(res)

			err := rep.Update(ctx, experience)
			Expect(err).To(Equal(NotFound))
		})
	})
})