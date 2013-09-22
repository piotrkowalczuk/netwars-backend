--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

--
-- Name: do_search(text, integer, integer); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION do_search(search_text text, zakres integer, user_id integer) RETURNS integer
    LANGUAGE plpgsql
    AS $$ 
    DECLARE 
       ret_fsearch_id INTEGER;
       post_number INTEGER :=0;
       topic_number INTEGER:=0;
    BEGIN
    
    insert into f_search (fsearch_range, fsearch_text, user_id, search_start, search_end, search_status) 
		values   (zakres,        search_text,  user_id, now(),        null,       1) returning fsearch_id into ret_fsearch_id;

    BEGIN 
	if zakres =1 or zakres=3 then
		insert into f_result (fsearch_id, fresult_type, fresult_id, fsearch_tsrank)
		SELECT ret_fsearch_id, 1, post_id, ts_rank_cd(to_tsvector('english',post_body), query) 
			FROM forum_post, to_tsquery('english',search_text) query
			WHERE query @@ to_tsvector('english',post_body) LIMIT 1000;
		GET DIAGNOSTICS post_number = ROW_COUNT;
	end if;

	if zakres =2 or zakres=3 then
		insert into f_result (fsearch_id, fresult_type, fresult_id, fsearch_tsrank)
		SELECT ret_fsearch_id, 2, topic_id, ts_rank_cd(to_tsvector('english',topic_name), query) 
			FROM forum_topic, to_tsquery('english',search_text) query
			WHERE query @@ to_tsvector('english',topic_name) LIMIT 1000;
		GET DIAGNOSTICS topic_number = ROW_COUNT;
	end if;

	update f_search set 
		search_end = now(),
		post_count = post_number,
		topic_count = topic_number,
		search_status = 2 where fsearch_id = ret_fsearch_id;
    EXCEPTION 
	WHEN OTHERS THEN 
		update f_search set 
		search_status = -1 where fsearch_id = ret_fsearch_id;	
    END;
    RETURN ret_fsearch_id;
END;
$$;


ALTER FUNCTION public.do_search(search_text text, zakres integer, user_id integer) OWNER TO postgres;

--
-- Name: forum_post_audit(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION forum_post_audit() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
    BEGIN
        --
        -- Create a row in emp_audit to reflect the operation performed on emp,
        -- make use of the special variable TG_OP to work out the operation.
        --
        IF (TG_OP = 'DELETE') THEN
            INSERT INTO forum_post_audit SELECT nextval('forum_post_audit_post_audit_id_seq'::regclass),'D', now(), OLD.*;
            RETURN OLD;
        ELSIF (TG_OP = 'UPDATE') THEN
            INSERT INTO forum_post_audit SELECT nextval('forum_post_audit_post_audit_id_seq'::regclass),'U', now(), NEW.*;
            RETURN NEW;
        ELSIF (TG_OP = 'INSERT') THEN
            INSERT INTO forum_post_audit SELECT nextval('forum_post_audit_post_audit_id_seq'::regclass),'I', now(), NEW.*;
            RETURN NEW;
        END IF;
        RETURN NULL; -- result is ignored since this is an AFTER trigger
    END;
$$;


ALTER FUNCTION public.forum_post_audit() OWNER TO postgres;

--
-- Name: forum_topic_audit(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION forum_topic_audit() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
    BEGIN
        --
        -- Create a row in emp_audit to reflect the operation performed on emp,
        -- make use of the special variable TG_OP to work out the operation.
        --
        IF (TG_OP = 'DELETE') THEN
            INSERT INTO forum_topic_audit SELECT nextval('forum_topic_audit_topic_audit_id_seq'::regclass),'D', now(), OLD.*;
            RETURN OLD;
        ELSIF (TG_OP = 'UPDATE') THEN
            INSERT INTO forum_topic_audit SELECT nextval('forum_topic_audit_topic_audit_id_seq'::regclass),'U', now(), NEW.*;
            RETURN NEW;
        ELSIF (TG_OP = 'INSERT') THEN
            INSERT INTO forum_topic_audit SELECT nextval('forum_topic_audit_topic_audit_id_seq'::regclass),'I', now(), NEW.*;
            RETURN NEW;
        END IF;
        RETURN NULL; -- result is ignored since this is an AFTER trigger
    END;
$$;


ALTER FUNCTION public.forum_topic_audit() OWNER TO postgres;

--
-- Name: update_sequences(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION update_sequences() RETURNS integer
    LANGUAGE plpgsql
    AS $$
BEGIN
    PERFORM setval('users_user_id_seq',
	(select coalesce(max(user_id),0)+1 from users));
    PERFORM setval('forum_forum_id_seq',
        (select coalesce(max(forum_id),0)+1 from forum));
    PERFORM setval('forum_topic_topic_id_seq',
        (select coalesce(max(topic_id),0)+1 from forum_topic));
    PERFORM setval('forum_post_post_id_seq',
        (select coalesce(max(post_id),0)+1 from forum_post));
        
    RETURN 1;
END;
$$;


ALTER FUNCTION public.update_sequences() OWNER TO postgres;

--
-- Name: users_audit(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION users_audit() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
    BEGIN
        --
        -- Create a row in emp_audit to reflect the operation performed on emp,
        -- make use of the special variable TG_OP to work out the operation.
        --
        IF (TG_OP = 'DELETE') THEN
            INSERT INTO users_audit SELECT nextval('users_audit_user_audit_id_seq'::regclass),'D', now(), OLD.*;
            RETURN OLD;
        ELSIF (TG_OP = 'UPDATE') THEN
            INSERT INTO users_audit SELECT nextval('users_audit_user_audit_id_seq'::regclass),'U', now(), NEW.*;
            RETURN NEW;
        ELSIF (TG_OP = 'INSERT') THEN
            INSERT INTO users_audit SELECT nextval('users_audit_user_audit_id_seq'::regclass),'I', now(), NEW.*;
            RETURN NEW;
        END IF;
        RETURN NULL; -- result is ignored since this is an AFTER trigger
    END;
$$;


ALTER FUNCTION public.users_audit() OWNER TO postgres;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: captcha; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE captcha (
    captcha_id integer NOT NULL,
    session_id text,
    captcha_string text,
    captcha_created timestamp without time zone,
    captcha_used integer
);


ALTER TABLE public.captcha OWNER TO postgres;

--
-- Name: captcha_captcha_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE captcha_captcha_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.captcha_captcha_id_seq OWNER TO postgres;

--
-- Name: captcha_captcha_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE captcha_captcha_id_seq OWNED BY captcha.captcha_id;


--
-- Name: errorlog; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE errorlog (
    errorlog_id integer NOT NULL,
    session_id text,
    script_path text,
    user_id integer,
    globals text,
    ip text,
    errorlog_time timestamp without time zone
);


ALTER TABLE public.errorlog OWNER TO postgres;

--
-- Name: errorlog_errorlog_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE errorlog_errorlog_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.errorlog_errorlog_id_seq OWNER TO postgres;

--
-- Name: errorlog_errorlog_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE errorlog_errorlog_id_seq OWNED BY errorlog.errorlog_id;


--
-- Name: errorlog_pos; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE errorlog_pos (
    errorlog_pos_id integer NOT NULL,
    errorlog_id bigint,
    log_type integer,
    log_message text,
    sql_query text,
    debug_backtrace text
);


ALTER TABLE public.errorlog_pos OWNER TO postgres;

--
-- Name: errorlog_pos_errorlog_pos_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE errorlog_pos_errorlog_pos_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.errorlog_pos_errorlog_pos_id_seq OWNER TO postgres;

--
-- Name: errorlog_pos_errorlog_pos_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE errorlog_pos_errorlog_pos_id_seq OWNED BY errorlog_pos.errorlog_pos_id;


--
-- Name: user_custom_style; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE user_custom_style (
    userid integer NOT NULL,
    style text,
    template text,
    exported character varying(1) DEFAULT 'N'::character varying,
    enabled character varying(1) DEFAULT 'N'::character varying
);


ALTER TABLE public.user_custom_style OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE users (
    user_id integer NOT NULL,
    user_name text,
    user_pass text,
    pass_type smallint,
    pass_salt text,
    last_login timestamp without time zone,
    bad_logins smallint,
    email text,
    ntcnick text,
    nickhistory text,
    user_status smallint,
    change_date timestamp without time zone,
    change_user_id integer,
    change_ip text,
    email_used integer,
    referrer integer,
    gg character varying(20),
    extrainfo text,
    created timestamp without time zone DEFAULT now(),
    trial character varying(1),
    showemail character varying(1) DEFAULT 'T'::character varying,
    refer_count integer DEFAULT 0,
    suspended character varying(1) DEFAULT 'N'::character varying
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: exported_custom_styles; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW exported_custom_styles AS
    SELECT u.user_id, u.user_name, ucs.style, ucs.template FROM (user_custom_style ucs LEFT JOIN users u ON ((ucs.userid = u.user_id))) WHERE ((ucs.exported)::text = 'T'::text);


ALTER TABLE public.exported_custom_styles OWNER TO postgres;

--
-- Name: f_result; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE f_result (
    fsearch_id integer NOT NULL,
    fresult_type integer NOT NULL,
    fresult_id integer NOT NULL,
    fsearch_tsrank real
);


ALTER TABLE public.f_result OWNER TO postgres;

--
-- Name: f_search; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE f_search (
    fsearch_id integer NOT NULL,
    fsearch_range smallint,
    fsearch_text text,
    user_id integer,
    search_status smallint,
    post_count integer,
    topic_count integer,
    search_start timestamp without time zone,
    search_end timestamp without time zone
);


ALTER TABLE public.f_search OWNER TO postgres;

--
-- Name: f_search_fsearch_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE f_search_fsearch_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.f_search_fsearch_id_seq OWNER TO postgres;

--
-- Name: f_search_fsearch_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE f_search_fsearch_id_seq OWNED BY f_search.fsearch_id;


--
-- Name: forum; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE forum (
    forum_id integer NOT NULL,
    forum_name text,
    forum_desc text,
    forum_order smallint,
    forum_type smallint,
    forum_topics integer,
    show_topics integer
);


ALTER TABLE public.forum OWNER TO postgres;

--
-- Name: forum_forum_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE forum_forum_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.forum_forum_id_seq OWNER TO postgres;

--
-- Name: forum_forum_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE forum_forum_id_seq OWNED BY forum.forum_id;


--
-- Name: forum_post; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE forum_post (
    post_id integer NOT NULL,
    topic_id bigint,
    user_id integer,
    user_name text,
    post_date timestamp without time zone,
    post_body text,
    mod_counter integer,
    mod_user_id integer,
    mod_user_name text,
    mod_date timestamp without time zone,
    ip_address text
);


ALTER TABLE public.forum_post OWNER TO postgres;

--
-- Name: forum_post_audit; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE forum_post_audit (
    post_audit_id integer NOT NULL,
    operation character(1) NOT NULL,
    stamp timestamp without time zone NOT NULL,
    post_id bigint,
    topic_id bigint,
    user_id integer,
    user_name text,
    post_date timestamp without time zone,
    post_body text,
    mod_counter integer,
    mod_user_id integer,
    mod_user_name text,
    mod_date timestamp without time zone,
    ip_address text
);


ALTER TABLE public.forum_post_audit OWNER TO postgres;

--
-- Name: forum_post_audit_post_audit_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE forum_post_audit_post_audit_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.forum_post_audit_post_audit_id_seq OWNER TO postgres;

--
-- Name: forum_post_audit_post_audit_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE forum_post_audit_post_audit_id_seq OWNED BY forum_post_audit.post_audit_id;


--
-- Name: forum_post_post_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE forum_post_post_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.forum_post_post_id_seq OWNER TO postgres;

--
-- Name: forum_post_post_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE forum_post_post_id_seq OWNED BY forum_post.post_id;


--
-- Name: forum_topic; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE forum_topic (
    forum_id integer NOT NULL,
    topic_id integer NOT NULL,
    topic_name text,
    first_poster integer,
    first_poster_name text,
    last_poster integer,
    last_poster_name text,
    last_post_id integer,
    last_post_date timestamp without time zone,
    topic_posts integer,
    topic_views integer,
    topic_closed smallint,
    topic_pined smallint,
    topic_visible_from timestamp without time zone,
    topic_visible_to timestamp without time zone,
    topic_deleted integer,
    change_date timestamp without time zone,
    change_user_id integer,
    change_ip text
);


ALTER TABLE public.forum_topic OWNER TO postgres;

--
-- Name: forum_topic_audit; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE forum_topic_audit (
    topic_audit_id integer NOT NULL,
    operation character(1) NOT NULL,
    stamp timestamp without time zone NOT NULL,
    forum_id integer NOT NULL,
    topic_id integer NOT NULL,
    topic_name text,
    first_poster integer,
    first_poster_name text,
    last_poster integer,
    last_poster_name text,
    last_post_id integer,
    last_post_date timestamp without time zone,
    topic_posts integer,
    topic_views integer,
    topic_closed smallint,
    topic_pined smallint,
    topic_visible_from timestamp without time zone,
    topic_visible_to timestamp without time zone,
    topic_deleted integer,
    change_date timestamp without time zone,
    change_user_id integer,
    change_ip text
);


ALTER TABLE public.forum_topic_audit OWNER TO postgres;

--
-- Name: forum_topic_audit_topic_audit_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE forum_topic_audit_topic_audit_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.forum_topic_audit_topic_audit_id_seq OWNER TO postgres;

--
-- Name: forum_topic_audit_topic_audit_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE forum_topic_audit_topic_audit_id_seq OWNED BY forum_topic_audit.topic_audit_id;


--
-- Name: forum_topic_topic_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE forum_topic_topic_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.forum_topic_topic_id_seq OWNER TO postgres;

--
-- Name: forum_topic_topic_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE forum_topic_topic_id_seq OWNED BY forum_topic.topic_id;


--
-- Name: suspicious_logins; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE suspicious_logins (
    user_id integer,
    login_time timestamp without time zone DEFAULT now()
);


ALTER TABLE public.suspicious_logins OWNER TO postgres;

--
-- Name: user_activationkeys; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE user_activationkeys (
    ua_id integer NOT NULL,
    user_id bigint,
    type_id integer,
    activation_chain text,
    create_time timestamp without time zone,
    already_used smallint
);


ALTER TABLE public.user_activationkeys OWNER TO postgres;

--
-- Name: user_activationkeys_ua_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE user_activationkeys_ua_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_activationkeys_ua_id_seq OWNER TO postgres;

--
-- Name: user_activationkeys_ua_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE user_activationkeys_ua_id_seq OWNED BY user_activationkeys.ua_id;


--
-- Name: user_ban; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE user_ban (
    user_id bigint NOT NULL,
    ban_desc text,
    ban_data_do timestamp without time zone,
    banned_by bigint,
    ban_post_id integer
);


ALTER TABLE public.user_ban OWNER TO postgres;

--
-- Name: user_ban_history; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE user_ban_history (
    history_id integer NOT NULL,
    user_id bigint,
    banned_by bigint,
    data timestamp without time zone,
    ban_desc text,
    action_type smallint,
    ban_data_do timestamp without time zone,
    ban_post_id integer
);


ALTER TABLE public.user_ban_history OWNER TO postgres;

--
-- Name: user_ban_history_history_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE user_ban_history_history_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_ban_history_history_id_seq OWNER TO postgres;

--
-- Name: user_ban_history_history_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE user_ban_history_history_id_seq OWNED BY user_ban_history.history_id;


--
-- Name: user_cookie; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE user_cookie (
    user_id integer NOT NULL,
    cookie_id text
);


ALTER TABLE public.user_cookie OWNER TO postgres;

--
-- Name: user_forum; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE user_forum (
    user_id integer NOT NULL,
    forum_id integer NOT NULL,
    show_topics integer,
    forum_order integer
);


ALTER TABLE public.user_forum OWNER TO postgres;

--
-- Name: user_ignore; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE user_ignore (
    user_id integer NOT NULL,
    ignored_id integer NOT NULL
);


ALTER TABLE public.user_ignore OWNER TO postgres;

--
-- Name: user_logins_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE user_logins_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_logins_id_seq OWNER TO postgres;

--
-- Name: user_logins; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE user_logins (
    loginid integer DEFAULT nextval('user_logins_id_seq'::regclass) NOT NULL,
    addr character varying(50),
    clientip character varying(50),
    xforwardedfor character varying(50),
    user_id integer,
    "time" timestamp without time zone DEFAULT now()
);


ALTER TABLE public.user_logins OWNER TO postgres;

--
-- Name: user_stream; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE user_stream (
    user_id integer NOT NULL,
    streamtype integer,
    handle character varying(255),
    wol character varying(1) DEFAULT 'Y'::character varying,
    hots character varying(1) DEFAULT 'Y'::character varying,
    lol character varying(1) DEFAULT 'Y'::character varying,
    bw character varying(1) DEFAULT 'Y'::character varying,
    other character varying(1) DEFAULT 'Y'::character varying
);


ALTER TABLE public.user_stream OWNER TO postgres;

--
-- Name: user_topic; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE user_topic (
    user_id integer NOT NULL,
    topic_id integer NOT NULL,
    observed smallint,
    post_seen integer
);


ALTER TABLE public.user_topic OWNER TO postgres;

--
-- Name: users_audit; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE users_audit (
    user_audit_id integer NOT NULL,
    operation character(1) NOT NULL,
    stamp timestamp without time zone NOT NULL,
    user_id integer NOT NULL,
    user_name text,
    user_pass text,
    pass_type smallint,
    pass_salt text,
    last_login timestamp without time zone,
    bad_logins smallint,
    email text,
    ntcnick text,
    nickhistory text,
    user_status smallint,
    change_date timestamp without time zone,
    change_user_id integer,
    change_ip text,
    email_used integer,
    referrer integer,
    gg character varying(20),
    extrainfo text,
    created timestamp without time zone,
    trial character varying(1),
    showemail character varying(1),
    refer_count integer DEFAULT 0,
    suspended character varying(1)
);


ALTER TABLE public.users_audit OWNER TO postgres;

--
-- Name: users_audit_user_audit_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE users_audit_user_audit_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_audit_user_audit_id_seq OWNER TO postgres;

--
-- Name: users_audit_user_audit_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE users_audit_user_audit_id_seq OWNED BY users_audit.user_audit_id;


--
-- Name: users_bak; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE users_bak (
    user_id integer,
    user_name text,
    user_pass text,
    pass_type smallint,
    pass_salt text,
    last_login timestamp without time zone,
    bad_logins smallint,
    email text,
    ntcnick text,
    nickhistory text,
    user_status smallint,
    change_date timestamp without time zone,
    change_user_id integer,
    change_ip text,
    email_used integer,
    referrer integer,
    gg character varying(20),
    extrainfo text,
    created timestamp without time zone,
    trial character varying(1),
    showemail character varying(1),
    refer_count integer
);


ALTER TABLE public.users_bak OWNER TO postgres;

--
-- Name: users_online; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE users_online (
    username character varying(255) NOT NULL,
    decay timestamp without time zone,
    userid integer
);


ALTER TABLE public.users_online OWNER TO postgres;

--
-- Name: users_user_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE users_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_user_id_seq OWNER TO postgres;

--
-- Name: users_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE users_user_id_seq OWNED BY users.user_id;


--
-- Name: captcha_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY captcha ALTER COLUMN captcha_id SET DEFAULT nextval('captcha_captcha_id_seq'::regclass);


--
-- Name: errorlog_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY errorlog ALTER COLUMN errorlog_id SET DEFAULT nextval('errorlog_errorlog_id_seq'::regclass);


--
-- Name: errorlog_pos_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY errorlog_pos ALTER COLUMN errorlog_pos_id SET DEFAULT nextval('errorlog_pos_errorlog_pos_id_seq'::regclass);


--
-- Name: fsearch_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY f_search ALTER COLUMN fsearch_id SET DEFAULT nextval('f_search_fsearch_id_seq'::regclass);


--
-- Name: forum_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY forum ALTER COLUMN forum_id SET DEFAULT nextval('forum_forum_id_seq'::regclass);


--
-- Name: post_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY forum_post ALTER COLUMN post_id SET DEFAULT nextval('forum_post_post_id_seq'::regclass);


--
-- Name: post_audit_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY forum_post_audit ALTER COLUMN post_audit_id SET DEFAULT nextval('forum_post_audit_post_audit_id_seq'::regclass);


--
-- Name: topic_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY forum_topic ALTER COLUMN topic_id SET DEFAULT nextval('forum_topic_topic_id_seq'::regclass);


--
-- Name: topic_audit_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY forum_topic_audit ALTER COLUMN topic_audit_id SET DEFAULT nextval('forum_topic_audit_topic_audit_id_seq'::regclass);


--
-- Name: ua_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY user_activationkeys ALTER COLUMN ua_id SET DEFAULT nextval('user_activationkeys_ua_id_seq'::regclass);


--
-- Name: history_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY user_ban_history ALTER COLUMN history_id SET DEFAULT nextval('user_ban_history_history_id_seq'::regclass);


--
-- Name: user_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users ALTER COLUMN user_id SET DEFAULT nextval('users_user_id_seq'::regclass);


--
-- Name: user_audit_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users_audit ALTER COLUMN user_audit_id SET DEFAULT nextval('users_audit_user_audit_id_seq'::regclass);


--
-- Name: captcha_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY captcha
    ADD CONSTRAINT captcha_pkey PRIMARY KEY (captcha_id);


--
-- Name: errorlog_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY errorlog
    ADD CONSTRAINT errorlog_pkey PRIMARY KEY (errorlog_id);


--
-- Name: errorlog_pos_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY errorlog_pos
    ADD CONSTRAINT errorlog_pos_pkey PRIMARY KEY (errorlog_pos_id);


--
-- Name: f_result_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY f_result
    ADD CONSTRAINT f_result_pkey PRIMARY KEY (fsearch_id, fresult_type, fresult_id);


--
-- Name: f_search_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY f_search
    ADD CONSTRAINT f_search_pkey PRIMARY KEY (fsearch_id);


--
-- Name: form_post_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY forum_post
    ADD CONSTRAINT form_post_pkey PRIMARY KEY (post_id);


--
-- Name: forum_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY forum
    ADD CONSTRAINT forum_pkey PRIMARY KEY (forum_id);


--
-- Name: forum_post_audit_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY forum_post_audit
    ADD CONSTRAINT forum_post_audit_pkey PRIMARY KEY (post_audit_id);


--
-- Name: forum_topic_audit_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY forum_topic_audit
    ADD CONSTRAINT forum_topic_audit_pkey PRIMARY KEY (topic_audit_id);


--
-- Name: forum_topic_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY forum_topic
    ADD CONSTRAINT forum_topic_pkey PRIMARY KEY (topic_id);


--
-- Name: user_activationkeys_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY user_activationkeys
    ADD CONSTRAINT user_activationkeys_pkey PRIMARY KEY (ua_id);


--
-- Name: user_ban_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY user_ban_history
    ADD CONSTRAINT user_ban_history_pkey PRIMARY KEY (history_id);


--
-- Name: user_ban_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY user_ban
    ADD CONSTRAINT user_ban_pkey PRIMARY KEY (user_id);


--
-- Name: user_cookie_cookie_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY user_cookie
    ADD CONSTRAINT user_cookie_cookie_id_key UNIQUE (cookie_id);


--
-- Name: user_cookie_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY user_cookie
    ADD CONSTRAINT user_cookie_pkey PRIMARY KEY (user_id);


--
-- Name: user_custom_style_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY user_custom_style
    ADD CONSTRAINT user_custom_style_pkey PRIMARY KEY (userid);


--
-- Name: user_forum_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY user_forum
    ADD CONSTRAINT user_forum_pkey PRIMARY KEY (user_id, forum_id);


--
-- Name: user_ignore_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY user_ignore
    ADD CONSTRAINT user_ignore_pkey PRIMARY KEY (user_id, ignored_id);


--
-- Name: user_logins_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY user_logins
    ADD CONSTRAINT user_logins_pkey PRIMARY KEY (loginid);


--
-- Name: user_stream_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY user_stream
    ADD CONSTRAINT user_stream_pkey PRIMARY KEY (user_id);


--
-- Name: user_topic_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY user_topic
    ADD CONSTRAINT user_topic_pkey PRIMARY KEY (user_id, topic_id);


--
-- Name: users_audit_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY users_audit
    ADD CONSTRAINT users_audit_pkey PRIMARY KEY (user_audit_id);


--
-- Name: users_online_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY users_online
    ADD CONSTRAINT users_online_pkey PRIMARY KEY (username);


--
-- Name: users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);


--
-- Name: byforum; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX byforum ON forum_topic USING btree (forum_id, topic_pined, last_post_id);


--
-- Name: bypined; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX bypined ON forum_topic USING btree (topic_pined);


--
-- Name: decay_index; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX decay_index ON users_online USING btree (decay);


--
-- Name: exported_custom_style; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX exported_custom_style ON user_custom_style USING btree (exported);


--
-- Name: forum_post_audit_post_id_post_audit_id_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX forum_post_audit_post_id_post_audit_id_idx ON forum_post_audit USING btree (post_id, post_audit_id);


--
-- Name: forum_post_topic_id_post_id_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX forum_post_topic_id_post_id_idx ON forum_post USING btree (topic_id, post_id);


--
-- Name: forum_post_user_id_post_id_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX forum_post_user_id_post_id_idx ON forum_post USING btree (user_id, post_id);


--
-- Name: forum_topic_forum_id_topic_pined_last_post_id_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX forum_topic_forum_id_topic_pined_last_post_id_idx ON forum_topic USING btree (forum_id, topic_pined DESC, last_post_id DESC);


--
-- Name: forum_topic_topic_id_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX forum_topic_topic_id_idx ON forum_topic USING btree (topic_id);


--
-- Name: fp_ft; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX fp_ft ON forum_post USING gin (to_tsvector('english'::regconfig, post_body));


--
-- Name: regcode_index; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX regcode_index ON users USING btree (md5((user_name || 'pozdro&pocwicz'::text)));


--
-- Name: user_activationkeys_activation_chain_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX user_activationkeys_activation_chain_idx ON user_activationkeys USING btree (activation_chain);


--
-- Name: user_logins_addr_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX user_logins_addr_idx ON user_logins USING btree (addr);


--
-- Name: user_logins_clientip_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX user_logins_clientip_idx ON user_logins USING btree (clientip);


--
-- Name: user_logins_user_id_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX user_logins_user_id_idx ON user_logins USING btree (user_id);


--
-- Name: user_logins_xforwardedfor_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX user_logins_xforwardedfor_idx ON user_logins USING btree (xforwardedfor);


--
-- Name: user_topic_user_id_topic_id_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX user_topic_user_id_topic_id_idx ON user_topic USING btree (user_id, topic_id);


--
-- Name: users_user_name_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX users_user_name_idx ON users USING btree (lower(user_name));


--
-- Name: forum_post_audit; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER forum_post_audit AFTER INSERT OR DELETE OR UPDATE ON forum_post FOR EACH ROW EXECUTE PROCEDURE forum_post_audit();


--
-- Name: forum_topic_audit; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER forum_topic_audit AFTER INSERT OR DELETE OR UPDATE ON forum_topic FOR EACH ROW EXECUTE PROCEDURE forum_topic_audit();


--
-- Name: users_audit; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER users_audit AFTER INSERT OR DELETE OR UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE users_audit();


--
-- Name: user_logins_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY user_logins
    ADD CONSTRAINT user_logins_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(user_id);


--
-- Name: users_referrer_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_referrer_fkey FOREIGN KEY (referrer) REFERENCES users(user_id);


--
-- PostgreSQL database dump complete
--

