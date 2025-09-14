CREATE TABLE public.movie_genre (
	movie_id int4 NULL,
	genre_id int4 NULL
);


-- public.movie_genre foreign keys

ALTER TABLE public.movie_genre ADD CONSTRAINT movie_genre_genre_id_fkey FOREIGN KEY (genre_id) REFERENCES public.genres(id);
ALTER TABLE public.movie_genre ADD CONSTRAINT movie_genre_movie_id_fkey FOREIGN KEY (movie_id) REFERENCES public.movies(id);