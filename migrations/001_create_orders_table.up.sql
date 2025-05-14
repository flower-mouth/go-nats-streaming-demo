--
-- PostgreSQL database dump
--

-- Dumped from database version 15.4
-- Dumped by pg_dump version 15.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: delivery; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.delivery (
                                 order_uid character varying(50),
                                 name character varying(50),
                                 phone character varying(50),
                                 zip character varying(50),
                                 city character varying(50),
                                 address character varying(50),
                                 region character varying(50),
                                 email character varying(50)
);


ALTER TABLE public.delivery OWNER TO postgres;

--
-- Name: items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.items (
                              order_uid character varying(50),
                              chrt_id bigint,
                              track_number character varying(50),
                              price bigint,
                              rid character varying(50),
                              name character varying(50),
                              sale smallint,
                              size character varying(50),
                              total_price bigint,
                              nm_id bigint,
                              brand character varying(50),
                              status smallint
);


ALTER TABLE public.items OWNER TO postgres;

--
-- Name: orders; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.orders (
                               order_uid character varying(50) NOT NULL,
                               track_number character varying(50),
                               entry character varying(50),
                               locale character varying(50),
                               internal_signature character varying(50),
                               customer_id character varying(50),
                               delivery_service character varying(50),
                               shardkey character varying(50),
                               sm_id bigint,
                               date_created character varying(50),
                               oof_shred character varying(50)
);


ALTER TABLE public.orders OWNER TO postgres;

--
-- Name: payment; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payment (
                                order_uid character varying(50),
                                transaction character varying(50),
                                request_id character varying(50),
                                currency character varying(50),
                                provider character varying(50),
                                amount bigint,
                                payment_dt bigint,
                                bank character varying(50),
                                delivery_cost bigint,
                                goods_total bigint,
                                custom_fee smallint
);


ALTER TABLE public.payment OWNER TO postgres;

--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (order_uid);


--
-- Name: delivery delivery_order_uid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.delivery
    ADD CONSTRAINT delivery_order_uid_fkey FOREIGN KEY (order_uid) REFERENCES public.orders(order_uid) ON DELETE CASCADE;


--
-- Name: items items_order_uid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_order_uid_fkey FOREIGN KEY (order_uid) REFERENCES public.orders(order_uid) ON DELETE CASCADE;


--
-- Name: payment payment_order_uid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment
    ADD CONSTRAINT payment_order_uid_fkey FOREIGN KEY (order_uid) REFERENCES public.orders(order_uid) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

