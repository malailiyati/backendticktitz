CREATE TABLE public.orders (
	id serial4 NOT NULL,
	users_id int4 NULL,
	schedule_id int4 NULL,
	payment_id int4 NULL,
	totalprice int4 NULL,
	fullname varchar(255) NOT NULL,
	email varchar(255) NOT NULL,
	phone varchar(100) NOT NULL,
	ispaid bool DEFAULT false NULL,
	qr_code text NOT NULL,
	created_at timestamp NULL,
	updated_at timestamp NULL,
	CONSTRAINT orders_pkey PRIMARY KEY (id)
);


-- public.orders foreign keys

ALTER TABLE public.orders ADD CONSTRAINT orders_payment_id_fkey FOREIGN KEY (payment_id) REFERENCES public.payments(id);
ALTER TABLE public.orders ADD CONSTRAINT orders_schedule_id_fkey FOREIGN KEY (schedule_id) REFERENCES public.schedule(id);
ALTER TABLE public.orders ADD CONSTRAINT orders_users_id_fkey FOREIGN KEY (users_id) REFERENCES public.users(id);