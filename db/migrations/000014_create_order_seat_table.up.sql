CREATE TABLE public.order_seat (
	order_id int4 NULL,
	seat_id int4 NULL
);


-- public.order_seat foreign keys

ALTER TABLE public.order_seat ADD CONSTRAINT order_seat_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id);
ALTER TABLE public.order_seat ADD CONSTRAINT order_seat_seat_id_fkey FOREIGN KEY (seat_id) REFERENCES public.seats(id);