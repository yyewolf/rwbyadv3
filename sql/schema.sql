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


--
-- Name: loot_boxes_type; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.loot_boxes_type AS ENUM (
    'classic',
    'rare',
    'special',
    'limited'
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
    deleted_at timestamp with time zone,
    xp integer DEFAULT 0 NOT NULL,
    next_level_xp integer NOT NULL,
    card_type character varying(50) NOT NULL,
    individual_value double precision NOT NULL,
    rarity integer NOT NULL,
    level integer NOT NULL,
    buffs integer NOT NULL
);


--
-- Name: cards_stats; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cards_stats (
    card_id character varying(50) NOT NULL,
    health integer NOT NULL,
    armor integer NOT NULL,
    damage integer NOT NULL,
    healing integer NOT NULL,
    speed integer NOT NULL,
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
    deleted_at timestamp with time zone,
    last_run_id bigint DEFAULT 0 NOT NULL,
    recurring boolean DEFAULT false NOT NULL,
    delta_time bigint DEFAULT 0 NOT NULL,
    errored boolean DEFAULT false NOT NULL
);


--
-- Name: loot_boxes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.loot_boxes (
    id character varying(50) NOT NULL,
    player_id character varying(50) NOT NULL,
    type public.loot_boxes_type NOT NULL,
    metadata json,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone
);


--
-- Name: player_card_favorites; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.player_card_favorites (
    player_id character varying(50) NOT NULL,
    card_id character varying(50) NOT NULL,
    "position" integer NOT NULL
);


--
-- Name: player_cards; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.player_cards (
    player_id character varying(50) NOT NULL,
    card_id character varying(50) NOT NULL,
    "position" integer NOT NULL
);


--
-- Name: player_cards_deck; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.player_cards_deck (
    player_id character varying(50) NOT NULL,
    card_id character varying(50) NOT NULL,
    "position" integer NOT NULL
);


--
-- Name: players; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.players (
    id character varying(50) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    selected_card_id character varying(50)
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
-- Name: cards_stats cards_stats_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cards_stats
    ADD CONSTRAINT cards_stats_pkey PRIMARY KEY (card_id);


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
-- Name: loot_boxes loot_boxes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.loot_boxes
    ADD CONSTRAINT loot_boxes_pkey PRIMARY KEY (id);


--
-- Name: player_card_favorites player_card_favorites_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.player_card_favorites
    ADD CONSTRAINT player_card_favorites_pkey PRIMARY KEY (player_id, card_id);


--
-- Name: player_cards_deck player_cards_deck_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.player_cards_deck
    ADD CONSTRAINT player_cards_deck_pkey PRIMARY KEY (player_id, card_id);


--
-- Name: player_cards player_cards_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.player_cards
    ADD CONSTRAINT player_cards_pkey PRIMARY KEY (player_id, card_id);


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
-- Name: cards_stats cards_stats_card_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cards_stats
    ADD CONSTRAINT cards_stats_card_id_fkey FOREIGN KEY (card_id) REFERENCES public.cards(id);


--
-- Name: cards fk_cards_stats; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cards
    ADD CONSTRAINT fk_cards_stats FOREIGN KEY (id) REFERENCES public.cards_stats(card_id) DEFERRABLE INITIALLY DEFERRED;


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
-- Name: loot_boxes loot_boxes_player_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.loot_boxes
    ADD CONSTRAINT loot_boxes_player_id_fkey FOREIGN KEY (player_id) REFERENCES public.players(id);


--
-- Name: player_card_favorites player_card_favorites_card_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.player_card_favorites
    ADD CONSTRAINT player_card_favorites_card_id_fkey FOREIGN KEY (card_id) REFERENCES public.cards(id);


--
-- Name: player_card_favorites player_card_favorites_player_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.player_card_favorites
    ADD CONSTRAINT player_card_favorites_player_id_fkey FOREIGN KEY (player_id) REFERENCES public.players(id);


--
-- Name: player_cards player_cards_card_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.player_cards
    ADD CONSTRAINT player_cards_card_id_fkey FOREIGN KEY (card_id) REFERENCES public.cards(id) ON DELETE CASCADE;


--
-- Name: player_cards_deck player_cards_deck_card_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.player_cards_deck
    ADD CONSTRAINT player_cards_deck_card_id_fkey FOREIGN KEY (card_id) REFERENCES public.cards(id);


--
-- Name: player_cards_deck player_cards_deck_player_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.player_cards_deck
    ADD CONSTRAINT player_cards_deck_player_id_fkey FOREIGN KEY (player_id) REFERENCES public.players(id);


--
-- Name: player_cards player_cards_player_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.player_cards
    ADD CONSTRAINT player_cards_player_id_fkey FOREIGN KEY (player_id) REFERENCES public.players(id);


--
-- Name: players players_selected_card_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.players
    ADD CONSTRAINT players_selected_card_id_fkey FOREIGN KEY (selected_card_id) REFERENCES public.cards(id);


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20240430114611'),
    ('20240501142003'),
    ('20240502144720'),
    ('20240503174603'),
    ('20240505092653'),
    ('20240506144248'),
    ('20240507105741');
