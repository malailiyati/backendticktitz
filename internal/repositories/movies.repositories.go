package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malailiyati/backend/internal/models"
	"github.com/malailiyati/backend/internal/utils"
	"github.com/redis/go-redis/v9"
)

type MovieRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewMovieRepository(db *pgxpool.Pool, rdb *redis.Client) *MovieRepository {
	return &MovieRepository{db: db, rdb: rdb}
}

func (r *MovieRepository) GetUpcomingMovies(ctx context.Context) ([]models.MovieSimpleResponse, error) {
	// cache-aside pattern
	// cek redis terlebih dahulu
	redisKey := "lala:movie-upcoming"
	cmd := r.rdb.Get(ctx, redisKey)
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			log.Printf("Key %s does not exist\n", redisKey)
		} else {
			log.Println("Redis Error. \nCause: ", cmd.Err().Error())
		}
	} else {
		// cache hit
		var cachedMovie []models.MovieSimpleResponse
		cmdByte, err := cmd.Bytes()
		if err != nil {
			log.Println("Internal Server Error.\nCause: ", err.Error())
		} else {
			if err := json.Unmarshal(cmdByte, &cachedMovie); err != nil {
				log.Println("Internal Server Error.\nCause: ", err.Error())
			}
			if len(cachedMovie) > 0 {
				return cachedMovie, nil
			}
		}
	}

	const q = `
		SELECT m.id, m.title, m.poster,
		       string_agg(g.name, ',') AS genres
		FROM movies m
		LEFT JOIN movie_genre mg ON mg.movie_id = m.id
		LEFT JOIN genres g ON g.id = mg.genre_id
		WHERE m.releaseDate > CURRENT_DATE 
		  AND m.deleted_at IS NULL
		GROUP BY m.id, m.title, m.poster
		ORDER BY m.releaseDate ASC, m.id;
	`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.MovieSimpleResponse
	for rows.Next() {
		var m models.MovieSimpleResponse
		var genres sql.NullString

		if err := rows.Scan(&m.ID, &m.Title, &m.Poster, &genres); err != nil {
			return nil, err
		}
		if genres.Valid {
			m.Genres = genres.String
		}
		movies = append(movies, m)
	}
	// renew cache
	bt, err := json.Marshal(movies)
	if err != nil {
		log.Println("Internal Server Error.\nCause: ", err.Error())
	} else {
		if err := r.rdb.Set(ctx, redisKey, string(bt), 5*time.Minute).Err(); err != nil {
			log.Println("Redis Error.\nCause: ", err.Error())
		}
	}

	return movies, nil
}

func (r *MovieRepository) GetPopularMovies(ctx context.Context, limit int) ([]models.MovieSimpleResponse, error) {
	// cache-aside pattern
	// cek redis terlebih dahulu
	redisKey := "lala:movie-popular"
	cmd := r.rdb.Get(ctx, redisKey)
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			log.Printf("Key %s does not exist\n", redisKey)
		} else {
			log.Println("Redis Error. \nCause: ", cmd.Err().Error())
		}
	} else {
		// cache hit
		var cachedMovie []models.MovieSimpleResponse
		cmdByte, err := cmd.Bytes()
		if err != nil {
			log.Println("Internal Server Error.\nCause: ", err.Error())
		} else {
			if err := json.Unmarshal(cmdByte, &cachedMovie); err != nil {
				log.Println("Internal Server Error.\nCause: ", err.Error())
			}
			if len(cachedMovie) > 0 {
				return cachedMovie, nil
			}
		}
	}

	const q = `
		SELECT m.id, m.title, m.poster,
		       string_agg(g.name, ',') AS genres
		FROM movies m
		LEFT JOIN movie_genre mg ON mg.movie_id = m.id
		LEFT JOIN genres g ON g.id = mg.genre_id
		WHERE m.deleted_at IS NULL
		GROUP BY m.id, m.title, m.poster
		ORDER BY m.popularity DESC, m.id
		LIMIT $1;
	`

	rows, err := r.db.Query(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.MovieSimpleResponse
	for rows.Next() {
		var m models.MovieSimpleResponse
		var genres sql.NullString

		if err := rows.Scan(&m.ID, &m.Title, &m.Poster, &genres); err != nil {
			return nil, err
		}
		if genres.Valid {
			m.Genres = genres.String
		}
		movies = append(movies, m)
	}

	// renew cache
	bt, err := json.Marshal(movies)
	if err != nil {
		log.Println("Internal Server Error.\nCause: ", err.Error())
	} else {
		if err := r.rdb.Set(ctx, redisKey, string(bt), 5*time.Minute).Err(); err != nil {
			log.Println("Redis Error.\nCause: ", err.Error())
		}
	}

	return movies, nil
}

func (r *MovieRepository) GetMoviesByFilter(ctx context.Context, title, genre string, limit, offset int) ([]models.MovieSimpleResponse, error) {
	var redisKey string
	useCache := offset == 0 // hanya cache page 1

	if useCache {
		redisKey = fmt.Sprintf("movies:filter:title=%s:genre=%s:limit=%d", title, genre, limit)
		if val, err := r.rdb.Get(ctx, redisKey).Result(); err == nil {
			var cached []models.MovieSimpleResponse
			if jsonErr := json.Unmarshal([]byte(val), &cached); jsonErr == nil {
				return cached, nil
			}
		} else if err != redis.Nil {
			log.Println("Redis Error:", err.Error())
		}
	}

	const q = `
		SELECT m.id, m.title, m.poster,
		       string_agg(g.name, ',') AS genres
		FROM movies m
		LEFT JOIN movie_genre mg ON mg.movie_id = m.id
		LEFT JOIN genres g ON g.id = mg.genre_id
		WHERE m.deleted_at IS NULL
		  AND ($1 = '' OR m.title ILIKE '%' || $1 || '%')
		  AND ($2 = '' OR m.id IN (
		      SELECT mg2.movie_id
		      FROM movie_genre mg2
		      JOIN genres g2 ON g2.id = mg2.genre_id
		      WHERE g2.name ILIKE '%' || $2 || '%'
		  ))
		GROUP BY m.id, m.title, m.poster
		ORDER BY m.created_at DESC, m.id
		LIMIT $3 OFFSET $4;
	`

	rows, err := r.db.Query(ctx, q, title, genre, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.MovieSimpleResponse
	for rows.Next() {
		var m models.MovieSimpleResponse
		var genres sql.NullString

		if err := rows.Scan(&m.ID, &m.Title, &m.Poster, &genres); err != nil {
			return nil, err
		}
		if genres.Valid {
			m.Genres = genres.String
		}
		movies = append(movies, m)
	}

	if useCache && len(movies) > 0 {
		bt, _ := json.Marshal(movies)
		_ = r.rdb.Set(ctx, redisKey, bt, 10*time.Minute).Err()
	}

	return movies, nil
}

func (r *MovieRepository) GetMovieDetailByID(ctx context.Context, id int) (*models.MovieDetail, error) {
	sql := `
   SELECT m.id, m.title, m.poster, m.background_poster, 
		       m.releaseDate, m.duration,  
		       m.synopsis, d.name AS director,
		       COALESCE(string_agg(DISTINCT g.name, ', '), '') AS genres,
		       COALESCE(string_agg(DISTINCT c.name, ', '), '') AS casts
		FROM movies m
		LEFT JOIN director d ON d.id = m.director_id
		LEFT JOIN movie_genre mg ON mg.movie_id = m.id
		LEFT JOIN genres g ON g.id = mg.genre_id
		LEFT JOIN movie_cast mc ON mc.movie_id = m.id
		LEFT JOIN casts c ON c.id = mc.cast_id
		WHERE m.id = $1
		GROUP BY m.id, m.title, m.poster, m.background_poster, 
		         m.releaseDate, m.duration, m.synopsis, d.name
    `

	row := r.db.QueryRow(ctx, sql, id)

	var md models.MovieDetail
	var iv pgtype.Interval

	err := row.Scan(
		&md.ID,
		&md.Title,
		&md.Poster,
		&md.BackgroundPoster,
		&md.ReleaseDate,
		&iv,
		&md.Synopsis,
		&md.Director,
		&md.Genres,
		&md.Casts,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	md.Duration = utils.FormatIntervalToText(iv)

	return &md, nil
}
