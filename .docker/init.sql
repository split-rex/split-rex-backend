--
-- PostgreSQL database dump
--

-- Dumped from database version 14.2
-- Dumped by pg_dump version 14.2

-- Started on 2023-04-24 19:34:08

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

--
-- TOC entry 3 (class 2615 OID 2200)
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres
--


ALTER SCHEMA public OWNER TO postgres;

--
-- TOC entry 3392 (class 0 OID 0)
-- Dependencies: 3
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: postgres
--

COMMENT ON SCHEMA public IS 'standard public schema';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 215 (class 1259 OID 32826)
-- Name: activities; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.activities (
    activity_id text NOT NULL,
    activity_type text NOT NULL,
    user_id text NOT NULL,
    date timestamp with time zone NOT NULL,
    redirect_id text NOT NULL,
    detail_id text NOT NULL
);


ALTER TABLE public.activities OWNER TO postgres;

--
-- TOC entry 220 (class 1259 OID 32861)
-- Name: expenses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.expenses (
    expense_id text NOT NULL,
    user_id text NOT NULL,
    amount numeric NOT NULL,
    date timestamp with time zone NOT NULL
);


ALTER TABLE public.expenses OWNER TO postgres;

--
-- TOC entry 211 (class 1259 OID 32798)
-- Name: friends; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.friends (
    id text NOT NULL,
    friend_id text,
    req_received text,
    req_sent text
);


ALTER TABLE public.friends OWNER TO postgres;

--
-- TOC entry 216 (class 1259 OID 32833)
-- Name: group_activities; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.group_activities (
    activity_id text NOT NULL,
    group_id text NOT NULL,
    user_id1 text NOT NULL,
    user_id2 text NOT NULL,
    amount numeric NOT NULL,
    date timestamp with time zone NOT NULL
);


ALTER TABLE public.group_activities OWNER TO postgres;

--
-- TOC entry 210 (class 1259 OID 32791)
-- Name: groups; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.groups (
    group_id text NOT NULL,
    name text NOT NULL,
    member_id text NOT NULL,
    start_date timestamp with time zone NOT NULL,
    end_date timestamp with time zone NOT NULL
);


ALTER TABLE public.groups OWNER TO postgres;

--
-- TOC entry 214 (class 1259 OID 32819)
-- Name: items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.items (
    item_id text NOT NULL,
    name text NOT NULL,
    quantity bigint NOT NULL,
    price numeric NOT NULL,
    consumer text NOT NULL
);


ALTER TABLE public.items OWNER TO postgres;

--
-- TOC entry 217 (class 1259 OID 32840)
-- Name: payment_activities; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payment_activities (
    payment_activity_id text NOT NULL,
    name text NOT NULL,
    status text NOT NULL,
    amount numeric NOT NULL,
    group_name text NOT NULL
);


ALTER TABLE public.payment_activities OWNER TO postgres;

--
-- TOC entry 213 (class 1259 OID 32812)
-- Name: payments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payments (
    payment_id text NOT NULL,
    group_id text NOT NULL,
    user_id1 text NOT NULL,
    user_id2 text NOT NULL,
    total_unpaid numeric NOT NULL,
    total_paid numeric NOT NULL,
    status text NOT NULL
);


ALTER TABLE public.payments OWNER TO postgres;

--
-- TOC entry 219 (class 1259 OID 32854)
-- Name: reminder_activities; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.reminder_activities (
    reminder_activity_id text NOT NULL,
    name text NOT NULL,
    group_name text NOT NULL
);


ALTER TABLE public.reminder_activities OWNER TO postgres;

--
-- TOC entry 218 (class 1259 OID 32847)
-- Name: transaction_activities; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transaction_activities (
    transaction_activity_id text NOT NULL,
    name text NOT NULL,
    group_name text NOT NULL
);


ALTER TABLE public.transaction_activities OWNER TO postgres;

--
-- TOC entry 212 (class 1259 OID 32805)
-- Name: transactions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transactions (
    transaction_id text NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    group_id text NOT NULL,
    date timestamp with time zone NOT NULL,
    subtotal numeric NOT NULL,
    tax numeric NOT NULL,
    service numeric NOT NULL,
    total numeric NOT NULL,
    bill_owner text NOT NULL,
    items text NOT NULL
);


ALTER TABLE public.transactions OWNER TO postgres;

--
-- TOC entry 209 (class 1259 OID 32779)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id text NOT NULL,
    name text NOT NULL,
    email text NOT NULL,
    username text NOT NULL,
    color bigint DEFAULT 1 NOT NULL,
    password bytea NOT NULL,
    groups text,
    payment_info text
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 3381 (class 0 OID 32826)
-- Dependencies: 215
-- Data for Name: activities; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.activities (activity_id, activity_type, user_id, date, redirect_id, detail_id) FROM stdin;
ccabfbe9-9a75-480e-b545-b439a8286f81	TRANSACTION	b607d360-3a6f-4654-8710-7c4e9dfcb5ac	2023-04-24 14:04:12.761236+07	0a17d941-aee8-43ce-a0b5-fdd6a07a335c	a189bf30-5463-4a31-af3d-5e4265fab487
8059d1d5-e2a7-4c0e-b7d6-34ccee4c37b7	TRANSACTION	7aa4bce4-584c-404b-92eb-f92e7576c377	2023-04-24 14:04:12.772074+07	0a17d941-aee8-43ce-a0b5-fdd6a07a335c	a189bf30-5463-4a31-af3d-5e4265fab487
114dc87a-89cc-45c0-9a4d-b41486d799d8	TRANSACTION	376d39fe-cf54-442c-a54b-fea91f1bd482	2023-04-24 14:12:13.059631+07	4607d510-654f-40be-9786-34dff56d5d4e	69be7ceb-ff11-4f60-a80c-174fb8b1568c
7f5ead85-596f-4f44-ba75-feec0527f9b0	TRANSACTION	b607d360-3a6f-4654-8710-7c4e9dfcb5ac	2023-04-24 14:12:13.061445+07	4607d510-654f-40be-9786-34dff56d5d4e	69be7ceb-ff11-4f60-a80c-174fb8b1568c
67816889-ec49-46d5-9647-839113a2c2f1	TRANSACTION	376d39fe-cf54-442c-a54b-fea91f1bd482	2023-04-24 14:16:34.429897+07	4d8f5325-e8bf-40a1-a530-367e1232371e	145ffbbc-3ac6-4d82-a29d-72c73cb14ad8
645a5f6e-1bf4-4487-a046-bcbc9bc56194	TRANSACTION	7aa4bce4-584c-404b-92eb-f92e7576c377	2023-04-24 14:16:34.433196+07	4d8f5325-e8bf-40a1-a530-367e1232371e	145ffbbc-3ac6-4d82-a29d-72c73cb14ad8
f600515d-9216-4a8c-89c5-7da024d9acdb	TRANSACTION	b607d360-3a6f-4654-8710-7c4e9dfcb5ac	2023-04-24 14:21:18.909768+07	2c0d42bd-dbb2-43c9-bf11-e55dd93672de	f221a2fe-8a90-44be-a47b-75903a18ecbf
644a5972-8665-4660-80bb-144de0200e06	TRANSACTION	5e41f06f-3364-489b-9b20-c401400efd06	2023-04-24 14:21:18.910731+07	2c0d42bd-dbb2-43c9-bf11-e55dd93672de	f221a2fe-8a90-44be-a47b-75903a18ecbf
ff7d24cb-6663-4a79-9679-68ae4cff11a4	PAYMENT	b607d360-3a6f-4654-8710-7c4e9dfcb5ac	2023-04-24 14:44:41.409316+07	d954f95a-70ba-4de5-a283-1e20224291cb	3afad79a-b884-4479-bb23-7cc6f4b99115
83ff0956-b45b-4435-ac3b-047709b8c97b	PAYMENT	7aa4bce4-584c-404b-92eb-f92e7576c377	2023-04-24 16:18:29.898755+07	2105941b-9c3a-48dd-afd3-6679994b2a5d	c8140a37-a7bd-49a9-a8c5-750acef4f2fc
80757dfa-8ffa-4bce-80ef-b9c918db1c82	PAYMENT	b607d360-3a6f-4654-8710-7c4e9dfcb5ac	2023-04-24 16:19:26.789945+07	2105941b-9c3a-48dd-afd3-6679994b2a5d	fb609452-894d-4d16-b208-32b9eae36512
9de47fbe-a6c3-45b5-a584-4280b94f407e	PAYMENT	376d39fe-cf54-442c-a54b-fea91f1bd482	2023-04-24 14:45:45.664584+07	d954f95a-70ba-4de5-a283-1e20224291cb	d3a6bb82-9dd8-4422-8f1f-29354e1f2c88
fe03dfa1-e1be-4eb5-ae02-edbc75784494	TRANSACTION	8a8ac694-b85d-43ca-b7c9-38f79852eb7d	2023-04-24 18:30:17.733071+07	995a8bb0-4fae-48b2-8c44-67083c26f43d	6afe524f-d3ec-4239-bdfa-afa8b7b89646
3e4e2a75-36dc-4087-adc3-cb7b8b6ea42e	TRANSACTION	175a21a4-c01c-4411-9ae6-f49fcab81bd2	2023-04-24 18:30:17.734297+07	995a8bb0-4fae-48b2-8c44-67083c26f43d	6afe524f-d3ec-4239-bdfa-afa8b7b89646
08a51037-eefa-42f6-85b7-35520b47b19b	TRANSACTION	eb4ff75a-158f-489c-ba4b-44fba94ad7a4	2023-04-24 18:30:17.735286+07	995a8bb0-4fae-48b2-8c44-67083c26f43d	6afe524f-d3ec-4239-bdfa-afa8b7b89646
\.


--
-- TOC entry 3386 (class 0 OID 32861)
-- Dependencies: 220
-- Data for Name: expenses; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.expenses (expense_id, user_id, amount, date) FROM stdin;
bafedd65-e977-482f-a304-0425849a98e9	175a21a4-c01c-4411-9ae6-f49fcab81bd2	10000	2023-04-24 18:32:12.840214+07
04020528-83f7-4114-b274-145915acc7f3	eb4ff75a-158f-489c-ba4b-44fba94ad7a4	30000	2023-04-24 18:32:12.842317+07
7325d9e3-4cd9-4c73-8990-c19c223039a3	8a8ac694-b85d-43ca-b7c9-38f79852eb7d	10000	2023-04-24 18:32:12.997851+07
\.


--
-- TOC entry 3377 (class 0 OID 32798)
-- Dependencies: 211
-- Data for Name: friends; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.friends (id, friend_id, req_received, req_sent) FROM stdin;
8a8ac694-b85d-43ca-b7c9-38f79852eb7d	["7aa4bce4-584c-404b-92eb-f92e7576c377"]	[]	\N
eb4ff75a-158f-489c-ba4b-44fba94ad7a4	["7aa4bce4-584c-404b-92eb-f92e7576c377"]	[]	\N
2884b0d7-f18e-4663-89d5-4933370e4779	["7aa4bce4-584c-404b-92eb-f92e7576c377"]	[]	\N
5e41f06f-3364-489b-9b20-c401400efd06	["7aa4bce4-584c-404b-92eb-f92e7576c377"]	[]	\N
175a21a4-c01c-4411-9ae6-f49fcab81bd2	["7aa4bce4-584c-404b-92eb-f92e7576c377"]	[]	\N
7aa4bce4-584c-404b-92eb-f92e7576c377	["b607d360-3a6f-4654-8710-7c4e9dfcb5ac","376d39fe-cf54-442c-a54b-fea91f1bd482","8a8ac694-b85d-43ca-b7c9-38f79852eb7d","eb4ff75a-158f-489c-ba4b-44fba94ad7a4","2884b0d7-f18e-4663-89d5-4933370e4779","5e41f06f-3364-489b-9b20-c401400efd06","175a21a4-c01c-4411-9ae6-f49fcab81bd2"]	\N	[]
376d39fe-cf54-442c-a54b-fea91f1bd482	["7aa4bce4-584c-404b-92eb-f92e7576c377", "b607d360-3a6f-4654-8710-7c4e9dfcb5ac"]	[]	\N
b607d360-3a6f-4654-8710-7c4e9dfcb5ac	["7aa4bce4-584c-404b-92eb-f92e7576c377", "376d39fe-cf54-442c-a54b-fea91f1bd482"]	[]	\N
\.


--
-- TOC entry 3382 (class 0 OID 32833)
-- Dependencies: 216
-- Data for Name: group_activities; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.group_activities (activity_id, group_id, user_id1, user_id2, amount, date) FROM stdin;
7430ab92-203c-4d32-8d9f-197d99de4cbc	d954f95a-70ba-4de5-a283-1e20224291cb	b607d360-3a6f-4654-8710-7c4e9dfcb5ac	376d39fe-cf54-442c-a54b-fea91f1bd482	10000	2023-04-24 14:45:45.664994+07
9d51bcf5-42bb-4120-81c8-b56607e6452f	2105941b-9c3a-48dd-afd3-6679994b2a5d	b607d360-3a6f-4654-8710-7c4e9dfcb5ac	7aa4bce4-584c-404b-92eb-f92e7576c377	5000	2023-04-24 16:19:26.790757+07
\.


--
-- TOC entry 3376 (class 0 OID 32791)
-- Dependencies: 210
-- Data for Name: groups; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.groups (group_id, name, member_id, start_date, end_date) FROM stdin;
d954f95a-70ba-4de5-a283-1e20224291cb	Group 1	["7aa4bce4-584c-404b-92eb-f92e7576c377","b607d360-3a6f-4654-8710-7c4e9dfcb5ac","376d39fe-cf54-442c-a54b-fea91f1bd482","7aa4bce4-584c-404b-92eb-f92e7576c377"]	2023-04-24 17:19:20.968831+07	2023-04-30 19:19:20.968831+07
e9c4d8cb-edf2-4c7e-b05e-7328810ad345	Group 2	["7aa4bce4-584c-404b-92eb-f92e7576c377","b607d360-3a6f-4654-8710-7c4e9dfcb5ac","376d39fe-cf54-442c-a54b-fea91f1bd482","8a8ac694-b85d-43ca-b7c9-38f79852eb7d","eb4ff75a-158f-489c-ba4b-44fba94ad7a4","2884b0d7-f18e-4663-89d5-4933370e4779","5e41f06f-3364-489b-9b20-c401400efd06","175a21a4-c01c-4411-9ae6-f49fcab81bd2","7aa4bce4-584c-404b-92eb-f92e7576c377"]	2023-04-24 17:19:20.968831+07	2023-04-30 19:19:20.968831+07
27ec8231-71db-4169-8680-9508c61ae559	Group 3	["7aa4bce4-584c-404b-92eb-f92e7576c377","8a8ac694-b85d-43ca-b7c9-38f79852eb7d","eb4ff75a-158f-489c-ba4b-44fba94ad7a4","175a21a4-c01c-4411-9ae6-f49fcab81bd2","7aa4bce4-584c-404b-92eb-f92e7576c377"]	2023-04-24 17:19:20.968831+07	2023-04-30 19:19:20.968831+07
d1cea18e-4d5c-4ee2-84d4-0b8186b65513	Group 4	["7aa4bce4-584c-404b-92eb-f92e7576c377","8a8ac694-b85d-43ca-b7c9-38f79852eb7d","eb4ff75a-158f-489c-ba4b-44fba94ad7a4","b607d360-3a6f-4654-8710-7c4e9dfcb5ac","7aa4bce4-584c-404b-92eb-f92e7576c377"]	2023-04-24 17:19:20.968831+07	2023-04-30 19:19:20.968831+07
2105941b-9c3a-48dd-afd3-6679994b2a5d	Group 5	["7aa4bce4-584c-404b-92eb-f92e7576c377","b607d360-3a6f-4654-8710-7c4e9dfcb5ac","5e41f06f-3364-489b-9b20-c401400efd06","7aa4bce4-584c-404b-92eb-f92e7576c377"]	2023-04-24 17:19:20.968831+07	2023-04-30 19:19:20.968831+07
\.


--
-- TOC entry 3380 (class 0 OID 32819)
-- Dependencies: 214
-- Data for Name: items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.items (item_id, name, quantity, price, consumer) FROM stdin;
566ab9b1-7ae9-48a3-85ed-b6aea33909c4	Item 1-1-1	1	10000	["7aa4bce4-584c-404b-92eb-f92e7576c377"]
ca574541-281b-40a8-bf6a-030370b3ad16	Item 1-1-2	1	10000	["7aa4bce4-584c-404b-92eb-f92e7576c377","b607d360-3a6f-4654-8710-7c4e9dfcb5ac"]
ecab99c1-c03f-4ae0-ae4a-ddc6854aec26	Item 1-2-1	1	30000	["376d39fe-cf54-442c-a54b-fea91f1bd482"]
01f06a38-1396-425a-846f-491a2ead3ef0	Item 1-2-2	1	20000	["b607d360-3a6f-4654-8710-7c4e9dfcb5ac"]
304f3f3f-d088-49c5-ab7e-b25a1dcb118f	Item 1-3-1	1	10000	["376d39fe-cf54-442c-a54b-fea91f1bd482"]
e4870b45-d1cb-42c9-9414-74649c474b30	Item 1-3-2	1	50000	["7aa4bce4-584c-404b-92eb-f92e7576c377"]
7b8dbd82-b8be-4b4d-9103-053bcf5e273d	Item 5-1-1	1	35000	["5e41f06f-3364-489b-9b20-c401400efd06"]
05081452-826a-4121-a6b2-edb67e1afdff	Item 5-1-2	1	15000	["b607d360-3a6f-4654-8710-7c4e9dfcb5ac"]
7c4a328c-17f7-4f93-bacf-be3d65c07eb3	Item 3-1-1	1	30000	["eb4ff75a-158f-489c-ba4b-44fba94ad7a4"]
403dd6d7-ef14-462f-83c5-74fffb918975	Item 3-1-2	1	20000	["8a8ac694-b85d-43ca-b7c9-38f79852eb7d","175a21a4-c01c-4411-9ae6-f49fcab81bd2"]
\.


--
-- TOC entry 3383 (class 0 OID 32840)
-- Dependencies: 217
-- Data for Name: payment_activities; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.payment_activities (payment_activity_id, name, status, amount, group_name) FROM stdin;
3afad79a-b884-4479-bb23-7cc6f4b99115	nando	UNCONFIRMED	10000	Group 1
c8140a37-a7bd-49a9-a8c5-750acef4f2fc	ubay	UNCONFIRMED	5000	Group 5
fb609452-894d-4d16-b208-32b9eae36512	patrick	CONFIRMED	5000	Group 5
d3a6bb82-9dd8-4422-8f1f-29354e1f2c88	ubay	CONFIRMED	10000	Group 1
\.


--
-- TOC entry 3379 (class 0 OID 32812)
-- Dependencies: 213
-- Data for Name: payments; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.payments (payment_id, group_id, user_id1, user_id2, total_unpaid, total_paid, status) FROM stdin;
db7c784c-b178-44b5-a45b-c67983cb8111	2105941b-9c3a-48dd-afd3-6679994b2a5d	7aa4bce4-584c-404b-92eb-f92e7576c377	5e41f06f-3364-489b-9b20-c401400efd06	-35000	0	UNPAID
304aed85-b1f5-43a3-a313-bd9a97f356d9	d954f95a-70ba-4de5-a283-1e20224291cb	376d39fe-cf54-442c-a54b-fea91f1bd482	b607d360-3a6f-4654-8710-7c4e9dfcb5ac	25000	0	UNPAID
d6c1bd9c-3693-4401-99f2-639eaab12b72	d954f95a-70ba-4de5-a283-1e20224291cb	b607d360-3a6f-4654-8710-7c4e9dfcb5ac	376d39fe-cf54-442c-a54b-fea91f1bd482	-25000	0	UNPAID
0a84fee1-1928-4ee9-8e76-dbf75d09f1d2	2105941b-9c3a-48dd-afd3-6679994b2a5d	7aa4bce4-584c-404b-92eb-f92e7576c377	b607d360-3a6f-4654-8710-7c4e9dfcb5ac	-10000	0	UNPAID
34760fdf-8971-4b27-935d-298e950c2850	2105941b-9c3a-48dd-afd3-6679994b2a5d	b607d360-3a6f-4654-8710-7c4e9dfcb5ac	7aa4bce4-584c-404b-92eb-f92e7576c377	10000	0	UNPAID
65daf1f7-5d16-4773-963f-b7d3b97376f8	27ec8231-71db-4169-8680-9508c61ae559	175a21a4-c01c-4411-9ae6-f49fcab81bd2	eb4ff75a-158f-489c-ba4b-44fba94ad7a4	-30000	0	UNPAID
c64746de-d347-4505-8178-791ce772c1c9	27ec8231-71db-4169-8680-9508c61ae559	eb4ff75a-158f-489c-ba4b-44fba94ad7a4	175a21a4-c01c-4411-9ae6-f49fcab81bd2	30000	0	UNPAID
46988e01-ec82-47b8-989a-79a62c454a58	27ec8231-71db-4169-8680-9508c61ae559	175a21a4-c01c-4411-9ae6-f49fcab81bd2	8a8ac694-b85d-43ca-b7c9-38f79852eb7d	-10000	0	UNPAID
d9660893-6e86-44e3-97ef-93bad0c2dbb6	27ec8231-71db-4169-8680-9508c61ae559	8a8ac694-b85d-43ca-b7c9-38f79852eb7d	175a21a4-c01c-4411-9ae6-f49fcab81bd2	10000	0	UNPAID
8ffa2714-b754-4b27-a2c9-9946af6dc50a	d954f95a-70ba-4de5-a283-1e20224291cb	b607d360-3a6f-4654-8710-7c4e9dfcb5ac	7aa4bce4-584c-404b-92eb-f92e7576c377	0	0	PAID
f4b09292-dec5-4225-8137-9f6f6e85fc96	d954f95a-70ba-4de5-a283-1e20224291cb	7aa4bce4-584c-404b-92eb-f92e7576c377	b607d360-3a6f-4654-8710-7c4e9dfcb5ac	0	0	PAID
715c470f-7208-4c46-b835-84ddcf2fc238	d954f95a-70ba-4de5-a283-1e20224291cb	376d39fe-cf54-442c-a54b-fea91f1bd482	7aa4bce4-584c-404b-92eb-f92e7576c377	5000	0	UNPAID
d5abbda9-4f38-43e7-9293-2c6e72db3528	d954f95a-70ba-4de5-a283-1e20224291cb	7aa4bce4-584c-404b-92eb-f92e7576c377	376d39fe-cf54-442c-a54b-fea91f1bd482	-5000	0	UNPAID
c7eea77e-1522-4ce7-9d43-6b65ee221a19	2105941b-9c3a-48dd-afd3-6679994b2a5d	5e41f06f-3364-489b-9b20-c401400efd06	7aa4bce4-584c-404b-92eb-f92e7576c377	35000	0	UNPAID
\.


--
-- TOC entry 3385 (class 0 OID 32854)
-- Dependencies: 219
-- Data for Name: reminder_activities; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.reminder_activities (reminder_activity_id, name, group_name) FROM stdin;
\.


--
-- TOC entry 3384 (class 0 OID 32847)
-- Dependencies: 218
-- Data for Name: transaction_activities; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.transaction_activities (transaction_activity_id, name, group_name) FROM stdin;
a189bf30-5463-4a31-af3d-5e4265fab487	patrick	Group 1
69be7ceb-ff11-4f60-a80c-174fb8b1568c	patrick	Group 1
145ffbbc-3ac6-4d82-a29d-72c73cb14ad8	ubay	Group 1
f221a2fe-8a90-44be-a47b-75903a18ecbf	patrick	Group 5
6afe524f-d3ec-4239-bdfa-afa8b7b89646	azka	Group 3
\.


--
-- TOC entry 3378 (class 0 OID 32805)
-- Dependencies: 212
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.transactions (transaction_id, name, description, group_id, date, subtotal, tax, service, total, bill_owner, items) FROM stdin;
0a17d941-aee8-43ce-a0b5-fdd6a07a335c	Transaction 1-1	lorem ipsum	d954f95a-70ba-4de5-a283-1e20224291cb	2023-04-24 17:19:20.968831+07	20000	0	0	20000	7aa4bce4-584c-404b-92eb-f92e7576c377	["566ab9b1-7ae9-48a3-85ed-b6aea33909c4","ca574541-281b-40a8-bf6a-030370b3ad16"]
4607d510-654f-40be-9786-34dff56d5d4e	Transaction 1-2	lorem ipsum	d954f95a-70ba-4de5-a283-1e20224291cb	2023-04-24 17:19:20.968831+07	50000	0	0	50000	7aa4bce4-584c-404b-92eb-f92e7576c377	["ecab99c1-c03f-4ae0-ae4a-ddc6854aec26","01f06a38-1396-425a-846f-491a2ead3ef0"]
4d8f5325-e8bf-40a1-a530-367e1232371e	Transaction 1-3	lorem ipsum	d954f95a-70ba-4de5-a283-1e20224291cb	2023-04-25 17:19:20.968831+07	50000	0	0	50000	b607d360-3a6f-4654-8710-7c4e9dfcb5ac	["304f3f3f-d088-49c5-ab7e-b25a1dcb118f","e4870b45-d1cb-42c9-9414-74649c474b30"]
2c0d42bd-dbb2-43c9-bf11-e55dd93672de	Transaction 5-1	lorem ipsum	2105941b-9c3a-48dd-afd3-6679994b2a5d	2023-04-26 17:19:20.968831+07	50000	0	0	50000	7aa4bce4-584c-404b-92eb-f92e7576c377	["7b8dbd82-b8be-4b4d-9103-053bcf5e273d","05081452-826a-4121-a6b2-edb67e1afdff"]
995a8bb0-4fae-48b2-8c44-67083c26f43d	Transaction 3-1	lorem ipsum	27ec8231-71db-4169-8680-9508c61ae559	2023-04-26 17:19:20.968831+07	50000	0	0	50000	175a21a4-c01c-4411-9ae6-f49fcab81bd2	["7c4a328c-17f7-4f93-bacf-be3d65c07eb3","403dd6d7-ef14-462f-83c5-74fffb918975"]
\.


--
-- TOC entry 3375 (class 0 OID 32779)
-- Dependencies: 209
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, name, email, username, color, password, groups, payment_info) FROM stdin;
b607d360-3a6f-4654-8710-7c4e9dfcb5ac	ubay	ubay@gmail.com	ubay	1	\\x24326124313024534c5845654c512f6b526e4e44356467463147684a4f67504e3248456d6d417969746935596a4d2f555a39334d7573396478583779	["d954f95a-70ba-4de5-a283-1e20224291cb","e9c4d8cb-edf2-4c7e-b05e-7328810ad345","d1cea18e-4d5c-4ee2-84d4-0b8186b65513","2105941b-9c3a-48dd-afd3-6679994b2a5d"]	\N
5e41f06f-3364-489b-9b20-c401400efd06	pakfitra	pakfitra@gmail.com	pakfitra	1	\\x243261243130244b7767597452673731726c41645a786f722e7776532e5a624f736d78674b596f7a3932474c543634596633304f2f6d534f36755065	["e9c4d8cb-edf2-4c7e-b05e-7328810ad345","2105941b-9c3a-48dd-afd3-6679994b2a5d"]	\N
7aa4bce4-584c-404b-92eb-f92e7576c377	patrick	patrick@gmail.com	patrick	1	\\x2432612431302446644e74436c5655584d4b464f726c66794e446e394f2e66786e535a796868694c5a35576c6737684b726a347566737244782f7947	["d954f95a-70ba-4de5-a283-1e20224291cb","e9c4d8cb-edf2-4c7e-b05e-7328810ad345","27ec8231-71db-4169-8680-9508c61ae559","d1cea18e-4d5c-4ee2-84d4-0b8186b65513","2105941b-9c3a-48dd-afd3-6679994b2a5d"]	\N
376d39fe-cf54-442c-a54b-fea91f1bd482	nando	nando@gmail.com	nando	1	\\x24326124313024396f6d76554d4a41384753637950674542306856502e4f50524c452f4a6168676f77526b523955723867755656487a5a5555624453	["d954f95a-70ba-4de5-a283-1e20224291cb","e9c4d8cb-edf2-4c7e-b05e-7328810ad345"]	\N
2884b0d7-f18e-4663-89d5-4933370e4779	azhar	azhar@gmail.com	azhar	1	\\x2432612431302436797673684e37346c436e4d33794f464e742f51464f626e4f4b614d73354a5a574b694353793736394561436d57793776744c6571	["e9c4d8cb-edf2-4c7e-b05e-7328810ad345"]	\N
175a21a4-c01c-4411-9ae6-f49fcab81bd2	azka	azka@gmail.com	azka	1	\\x2432612431302449764c3349547834695772636b44626831783671514f727a4f64306e477656706c48734c69726e5848796e56323239494a4c324b75	["e9c4d8cb-edf2-4c7e-b05e-7328810ad345","27ec8231-71db-4169-8680-9508c61ae559"]	\N
8a8ac694-b85d-43ca-b7c9-38f79852eb7d	grace	grace@gmail.com	grace	1	\\x24326124313024496975452e61513151584235435a51354a33335a75755154496f675075496c482f646e536e4e6a766e75666b63494570654b51644f	["e9c4d8cb-edf2-4c7e-b05e-7328810ad345","27ec8231-71db-4169-8680-9508c61ae559","d1cea18e-4d5c-4ee2-84d4-0b8186b65513"]	\N
eb4ff75a-158f-489c-ba4b-44fba94ad7a4	samuel	samuel@gmail.com	samuel	1	\\x2432612431302461354b5a6b6771485552324d726d78576f357367324f4164695a783138364b466a732f6a4d694a444d746350332f4c473853587a65	["e9c4d8cb-edf2-4c7e-b05e-7328810ad345","27ec8231-71db-4169-8680-9508c61ae559","d1cea18e-4d5c-4ee2-84d4-0b8186b65513"]	\N
\.


--
-- TOC entry 3225 (class 2606 OID 32832)
-- Name: activities activities_activity_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.activities
    ADD CONSTRAINT activities_activity_id_key UNIQUE (activity_id);


--
-- TOC entry 3235 (class 2606 OID 32867)
-- Name: expenses expenses_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expenses
    ADD CONSTRAINT expenses_pkey PRIMARY KEY (expense_id);


--
-- TOC entry 3217 (class 2606 OID 32804)
-- Name: friends friends_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.friends
    ADD CONSTRAINT friends_pkey PRIMARY KEY (id);


--
-- TOC entry 3227 (class 2606 OID 32839)
-- Name: group_activities group_activities_activity_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.group_activities
    ADD CONSTRAINT group_activities_activity_id_key UNIQUE (activity_id);


--
-- TOC entry 3215 (class 2606 OID 32797)
-- Name: groups groups_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.groups
    ADD CONSTRAINT groups_pkey PRIMARY KEY (group_id);


--
-- TOC entry 3223 (class 2606 OID 32825)
-- Name: items items_item_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_item_id_key UNIQUE (item_id);


--
-- TOC entry 3229 (class 2606 OID 32846)
-- Name: payment_activities payment_activities_payment_activity_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_activities
    ADD CONSTRAINT payment_activities_payment_activity_id_key UNIQUE (payment_activity_id);


--
-- TOC entry 3221 (class 2606 OID 32818)
-- Name: payments payments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_pkey PRIMARY KEY (payment_id);


--
-- TOC entry 3233 (class 2606 OID 32860)
-- Name: reminder_activities reminder_activities_reminder_activity_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reminder_activities
    ADD CONSTRAINT reminder_activities_reminder_activity_id_key UNIQUE (reminder_activity_id);


--
-- TOC entry 3231 (class 2606 OID 32853)
-- Name: transaction_activities transaction_activities_transaction_activity_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transaction_activities
    ADD CONSTRAINT transaction_activities_transaction_activity_id_key UNIQUE (transaction_activity_id);


--
-- TOC entry 3219 (class 2606 OID 32811)
-- Name: transactions transactions_transaction_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_transaction_id_key UNIQUE (transaction_id);


--
-- TOC entry 3209 (class 2606 OID 32788)
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- TOC entry 3211 (class 2606 OID 32786)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 3213 (class 2606 OID 32790)
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


-- Completed on 2023-04-24 19:34:09

--
-- PostgreSQL database dump complete
--
