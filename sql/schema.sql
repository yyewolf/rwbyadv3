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
-- Name: auth_github_states_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.auth_github_states_type AS ENUM (
    'check_star'
);


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: auth_github_states; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.auth_github_states (
    state character varying(50) NOT NULL,
    player_id character varying(50) NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    type public.auth_github_states_type NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone
);


--
-- Name: cards; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cards (
    id character varying(50) NOT NULL,
    player_id character varying(50) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone
);


--
-- Name: github_stars; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.github_stars (
    player_id character varying(50) NOT NULL,
    github_user_id character varying(50),
    has_starred boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone
);


--
-- Name: jobs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.jobs (
    id character varying(100) NOT NULL,
    jobkey character varying(100) NOT NULL,
    retries integer DEFAULT 0 NOT NULL,
    run_at timestamp with time zone NOT NULL,
    params jsonb NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone
);


--
-- Name: players; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.players (
    id character varying(50) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying(128) NOT NULL
);


--
-- Name: auth_github_states auth_github_states_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.auth_github_states
    ADD CONSTRAINT auth_github_states_pkey PRIMARY KEY (state);


--
-- Name: cards cards_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cards
    ADD CONSTRAINT cards_pkey PRIMARY KEY (id);


--
-- Name: github_stars github_stars_github_user_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.github_stars
    ADD CONSTRAINT github_stars_github_user_id_key UNIQUE (github_user_id);


--
-- Name: github_stars github_stars_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.github_stars
    ADD CONSTRAINT github_stars_pkey PRIMARY KEY (player_id);


--
-- Name: jobs jobs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.jobs
    ADD CONSTRAINT jobs_pkey PRIMARY KEY (id, jobkey);


--
-- Name: players players_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.players
    ADD CONSTRAINT players_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: auth_github_states auth_github_states_player_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.auth_github_states
    ADD CONSTRAINT auth_github_states_player_id_fkey FOREIGN KEY (player_id) REFERENCES public.players(id);


--
-- Name: cards cards_player_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cards
    ADD CONSTRAINT cards_player_id_fkey FOREIGN KEY (player_id) REFERENCES public.players(id);


--
-- Name: players fk_github_star; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.players
    ADD CONSTRAINT fk_github_star FOREIGN KEY (id) REFERENCES public.github_stars(player_id) DEFERRABLE INITIALLY DEFERRED;


--
-- Name: github_stars fk_player_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.github_stars
    ADD CONSTRAINT fk_player_id FOREIGN KEY (player_id) REFERENCES public.players(id) DEFERRABLE INITIALLY DEFERRED;


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20240430114611'),
    ('20240501142003'),
    ('20240502144720');
