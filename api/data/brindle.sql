--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5 (Ubuntu 14.5-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.5 (Ubuntu 14.5-0ubuntu0.22.04.1)

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
-- Name: configurations; Type: TABLE; Schema: public; Owner: mothman
--

CREATE TABLE public.configurations (
    id integer NOT NULL,
    username character varying(255) DEFAULT 'no-user'::character varying,
    key character varying(255) DEFAULT 'no-key'::character varying,
    logpath character varying(255) DEFAULT 'no-path'::character varying,
    hosts text[] DEFAULT '{"nohosts check config"}'::text[],
    command character varying(255) DEFAULT 'no-cmd'::character varying,
    timeout bigint DEFAULT 60 NOT NULL,
    port bigint DEFAULT 22 NOT NULL,
    fatal boolean DEFAULT false NOT NULL,
    ordered boolean DEFAULT false NOT NULL,
    name character varying(255) DEFAULT 'empty_config'::character varying NOT NULL
);


ALTER TABLE public.configurations OWNER TO mothman;

--
-- Name: configurations_id_seq; Type: SEQUENCE; Schema: public; Owner: mothman
--

CREATE SEQUENCE public.configurations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.configurations_id_seq OWNER TO mothman;

--
-- Name: configurations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mothman
--

ALTER SEQUENCE public.configurations_id_seq OWNED BY public.configurations.id;


--
-- Name: replies; Type: TABLE; Schema: public; Owner: mothman
--

CREATE TABLE public.replies (
    id integer NOT NULL,
    command_sent timestamp without time zone DEFAULT now() NOT NULL,
    reply_received timestamp without time zone DEFAULT now(),
    reply jsonb DEFAULT '{}'::jsonb,
    config jsonb DEFAULT '{}'::jsonb,
    good boolean DEFAULT false NOT NULL,
    host character varying(255) DEFAULT 'nonE'::character varying NOT NULL
);


ALTER TABLE public.replies OWNER TO mothman;

--
-- Name: replies_id_seq; Type: SEQUENCE; Schema: public; Owner: mothman
--

CREATE SEQUENCE public.replies_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.replies_id_seq OWNER TO mothman;

--
-- Name: replies_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mothman
--

ALTER SEQUENCE public.replies_id_seq OWNED BY public.replies.id;


--
-- Name: targets; Type: TABLE; Schema: public; Owner: mothman
--

CREATE TABLE public.targets (
    id integer NOT NULL,
    port bigint DEFAULT 22 NOT NULL,
    protocol bigint DEFAULT 22 NOT NULL,
    address character varying(255) DEFAULT 'localhost'::character varying NOT NULL,
    "user" character varying(255) DEFAULT 'mrbyte'::character varying,
    key character varying(255) DEFAULT '~/.ssh/id_rsa'::character varying,
    password character varying(255) DEFAULT 'f7219fd24ffb33e607ebda9fcfe22193d2983ef6'::character varying,
    token character varying(255) DEFAULT 'no-tk'::character varying,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.targets OWNER TO mothman;

--
-- Name: targets_id_seq; Type: SEQUENCE; Schema: public; Owner: mothman
--

CREATE SEQUENCE public.targets_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.targets_id_seq OWNER TO mothman;

--
-- Name: targets_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mothman
--

ALTER SEQUENCE public.targets_id_seq OWNED BY public.targets.id;


--
-- Name: tokens; Type: TABLE; Schema: public; Owner: mothman
--

CREATE TABLE public.tokens (
    id integer NOT NULL,
    user_id integer,
    email character varying(255) NOT NULL,
    token character varying(255) NOT NULL,
    token_hash bytea NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    death timestamp with time zone NOT NULL
);


ALTER TABLE public.tokens OWNER TO mothman;

--
-- Name: tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: mothman
--

CREATE SEQUENCE public.tokens_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.tokens_id_seq OWNER TO mothman;

--
-- Name: tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mothman
--

ALTER SEQUENCE public.tokens_id_seq OWNED BY public.tokens.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: mothman
--

CREATE TABLE public.users (
    id integer NOT NULL,
    email character varying(255) NOT NULL,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    password character varying(60) NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    user_active integer DEFAULT 0 NOT NULL,
    settings json
);


ALTER TABLE public.users OWNER TO mothman;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: mothman
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO mothman;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mothman
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: configurations id; Type: DEFAULT; Schema: public; Owner: mothman
--

ALTER TABLE ONLY public.configurations ALTER COLUMN id SET DEFAULT nextval('public.configurations_id_seq'::regclass);


--
-- Name: replies id; Type: DEFAULT; Schema: public; Owner: mothman
--

ALTER TABLE ONLY public.replies ALTER COLUMN id SET DEFAULT nextval('public.replies_id_seq'::regclass);


--
-- Name: targets id; Type: DEFAULT; Schema: public; Owner: mothman
--

ALTER TABLE ONLY public.targets ALTER COLUMN id SET DEFAULT nextval('public.targets_id_seq'::regclass);


--
-- Name: tokens id; Type: DEFAULT; Schema: public; Owner: mothman
--

ALTER TABLE ONLY public.tokens ALTER COLUMN id SET DEFAULT nextval('public.tokens_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: mothman
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: configurations configurations_pkey; Type: CONSTRAINT; Schema: public; Owner: mothman
--

ALTER TABLE ONLY public.configurations
    ADD CONSTRAINT configurations_pkey PRIMARY KEY (id);


--
-- Name: replies replies_pkey; Type: CONSTRAINT; Schema: public; Owner: mothman
--

ALTER TABLE ONLY public.replies
    ADD CONSTRAINT replies_pkey PRIMARY KEY (id);


--
-- Name: targets targets_pkey; Type: CONSTRAINT; Schema: public; Owner: mothman
--

ALTER TABLE ONLY public.targets
    ADD CONSTRAINT targets_pkey PRIMARY KEY (id);


--
-- Name: tokens tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: mothman
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: mothman
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: tokens tokens_relation_1; Type: FK CONSTRAINT; Schema: public; Owner: mothman
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_relation_1 FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

