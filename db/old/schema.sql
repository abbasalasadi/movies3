--
-- PostgreSQL database dump
--

\restrict Oqp8baNyTzS4OksCUCPHnSlnh5WESUbg5klTUPotz7nFgWiGFe5UIYqJZIYFNv3

-- Dumped from database version 18.1 (Ubuntu 18.1-1.pgdg24.04+2)
-- Dumped by pg_dump version 18.1 (Ubuntu 18.1-1.pgdg24.04+2)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: Lines; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA "Lines";


--
-- Name: MonitorPackages; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA "MonitorPackages";


--
-- Name: References; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA "References";


--
-- Name: Tables; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA "Tables";


--
-- Name: pg_trgm; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA "Tables";


--
-- Name: EXTENSION pg_trgm; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pg_trgm IS 'text similarity measurement and index searching based on trigrams';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: AwardTitleLine; Type: TABLE; Schema: Lines; Owner: -
--

CREATE TABLE "Lines"."AwardTitleLine" (
    "TitleID" integer NOT NULL,
    "EventID" integer NOT NULL,
    "CastID" bigint NOT NULL,
    "AwardYear" character varying(255) NOT NULL,
    "NominationType" smallint NOT NULL,
    "Description" character varying(255) NOT NULL,
    "Category" character varying(255) NOT NULL
);


--
-- Name: CastTitleLine; Type: TABLE; Schema: Lines; Owner: -
--

CREATE TABLE "Lines"."CastTitleLine" (
    "TitleID" integer NOT NULL,
    "CastID" integer NOT NULL,
    "CastType" smallint NOT NULL,
    "CastRole" character varying(255),
    "Sequence" smallint NOT NULL
);


--
-- Name: CertificateTitleLine; Type: TABLE; Schema: Lines; Owner: -
--

CREATE TABLE "Lines"."CertificateTitleLine" (
    "TitleID" integer NOT NULL,
    "CountryID" smallint NOT NULL,
    "CertificateID" integer NOT NULL
);


--
-- Name: CompanyTitleLine; Type: TABLE; Schema: Lines; Owner: -
--

CREATE TABLE "Lines"."CompanyTitleLine" (
    "TitleID" integer NOT NULL,
    "CompanyID" integer NOT NULL
);


--
-- Name: ConnectionTitleLine; Type: TABLE; Schema: Lines; Owner: -
--

CREATE TABLE "Lines"."ConnectionTitleLine" (
    "TitleID" integer NOT NULL,
    "ConnectionTitleID" integer NOT NULL,
    "ConnectionType" smallint NOT NULL
);


--
-- Name: CountryTitleLine; Type: TABLE; Schema: Lines; Owner: -
--

CREATE TABLE "Lines"."CountryTitleLine" (
    "TitleID" integer NOT NULL,
    "CountryID" smallint NOT NULL
);


--
-- Name: FileTitleLine; Type: TABLE; Schema: Lines; Owner: -
--

CREATE TABLE "Lines"."FileTitleLine" (
    "TitleID" integer NOT NULL,
    "QualityID" smallint NOT NULL,
    "DisplayID" integer NOT NULL,
    "AudioLanguageID" integer NOT NULL,
    "SubtitleLanguageID" integer NOT NULL
);


--
-- Name: GenreTitleLine; Type: TABLE; Schema: Lines; Owner: -
--

CREATE TABLE "Lines"."GenreTitleLine" (
    "TitleID" integer NOT NULL,
    "GenreID" integer NOT NULL
);


--
-- Name: KnownAsTitleLine; Type: TABLE; Schema: Lines; Owner: -
--

CREATE TABLE "Lines"."KnownAsTitleLine" (
    "TitleID" integer NOT NULL,
    "KnownAs" character varying(255) NOT NULL
);


--
-- Name: LanguageTitleLine; Type: TABLE; Schema: Lines; Owner: -
--

CREATE TABLE "Lines"."LanguageTitleLine" (
    "TitleID" integer NOT NULL,
    "LanguageID" smallint NOT NULL
);


--
-- Name: SimilaritiesTitleLine; Type: TABLE; Schema: Lines; Owner: -
--

CREATE TABLE "Lines"."SimilaritiesTitleLine" (
    "TitleID" integer NOT NULL,
    "SimilarTitleID" integer NOT NULL
);


--
-- Name: MonitorPackage1; Type: TABLE; Schema: MonitorPackages; Owner: -
--

CREATE TABLE "MonitorPackages"."MonitorPackage1" (
    "Application" integer NOT NULL,
    "Status" smallint,
    "ItemsCompleted" smallint,
    "ActiveTitleID" integer,
    "ClosedAt" timestamp(6) without time zone,
    "Average" numeric(16,2),
    "RunOrder" boolean,
    "ShowFP" boolean
);


--
-- Name: MonitorPackage2; Type: TABLE; Schema: MonitorPackages; Owner: -
--

CREATE TABLE "MonitorPackages"."MonitorPackage2" (
    "Application" integer NOT NULL,
    "Status" smallint,
    "ItemsCompleted" smallint,
    "ActiveTitleID" integer,
    "ClosedAt" timestamp(6) without time zone,
    "Average" numeric(16,2),
    "RunOrder" boolean,
    "ShowFP" boolean
);


--
-- Name: MonitorPackage3; Type: TABLE; Schema: MonitorPackages; Owner: -
--

CREATE TABLE "MonitorPackages"."MonitorPackage3" (
    "Application" integer NOT NULL,
    "Status" smallint,
    "ItemsCompleted" smallint,
    "ActiveTitleID" integer,
    "ClosedAt" timestamp(6) without time zone,
    "Average" numeric(16,2),
    "RunOrder" boolean,
    "ShowFP" boolean
);


--
-- Name: MonitorPackage4; Type: TABLE; Schema: MonitorPackages; Owner: -
--

CREATE TABLE "MonitorPackages"."MonitorPackage4" (
    "Application" integer NOT NULL,
    "Status" smallint,
    "ItemsCompleted" smallint,
    "ActiveTitleID" integer,
    "ClosedAt" timestamp(6) without time zone,
    "Average" numeric(16,2),
    "RunOrder" boolean,
    "ShowFP" boolean
);


--
-- Name: AwardEventRef; Type: TABLE; Schema: References; Owner: -
--

CREATE TABLE "References"."AwardEventRef" (
    "EventID" integer NOT NULL,
    "EventName" character varying(255) NOT NULL
);


--
-- Name: AwardNominationTypeRef; Type: TABLE; Schema: References; Owner: -
--

CREATE TABLE "References"."AwardNominationTypeRef" (
    "NominationTypeID" smallint NOT NULL,
    "NominationType" character varying(255) NOT NULL
);


--
-- Name: CastTypeRef; Type: TABLE; Schema: References; Owner: -
--

CREATE TABLE "References"."CastTypeRef" (
    "CastTypeID" smallint NOT NULL,
    "CastTypeDescription" character varying(255) NOT NULL
);


--
-- Name: CategoryRef; Type: TABLE; Schema: References; Owner: -
--

CREATE TABLE "References"."CategoryRef" (
    "CategoryID" smallint NOT NULL,
    "CategoryDecription" character varying(255) NOT NULL
);


--
-- Name: CertificateCountryRef; Type: TABLE; Schema: References; Owner: -
--

CREATE TABLE "References"."CertificateCountryRef" (
    "CountryID" integer NOT NULL,
    "CertificateID" integer NOT NULL,
    "Age" integer NOT NULL
);


--
-- Name: CertificateRef; Type: TABLE; Schema: References; Owner: -
--

CREATE TABLE "References"."CertificateRef" (
    "CertificateID" integer NOT NULL,
    "CertificateName" character varying(255) NOT NULL
);


--
-- Name: CertificateRef_CertificateID_seq; Type: SEQUENCE; Schema: References; Owner: -
--

ALTER TABLE "References"."CertificateRef" ALTER COLUMN "CertificateID" ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME "References"."CertificateRef_CertificateID_seq"
    START WITH 1
    INCREMENT BY 1
    MINVALUE 0
    NO MAXVALUE
    CACHE 1
);


--
-- Name: ConnectionTypeRef; Type: TABLE; Schema: References; Owner: -
--

CREATE TABLE "References"."ConnectionTypeRef" (
    "ConnectionTypeID" smallint NOT NULL,
    "ConnectionTypeDescription" character varying(255) NOT NULL
);


--
-- Name: CountryRef; Type: TABLE; Schema: References; Owner: -
--

CREATE TABLE "References"."CountryRef" (
    "CountryID" smallint NOT NULL,
    "CountryName" character varying(255) NOT NULL,
    "CountryCode" character varying(255) NOT NULL
);


--
-- Name: CountryRef_CountryID_seq; Type: SEQUENCE; Schema: References; Owner: -
--

ALTER TABLE "References"."CountryRef" ALTER COLUMN "CountryID" ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME "References"."CountryRef_CountryID_seq"
    START WITH 0
    INCREMENT BY 1
    MINVALUE 0
    NO MAXVALUE
    CACHE 1
);


--
-- Name: DisplayRef; Type: TABLE; Schema: References; Owner: -
--

CREATE TABLE "References"."DisplayRef" (
    "DisplayID" smallint NOT NULL,
    "DisplayType" character varying(255)
);


--
-- Name: GenreRef; Type: TABLE; Schema: References; Owner: -
--

CREATE TABLE "References"."GenreRef" (
    "GenreID" smallint NOT NULL,
    "GenreName" character varying(255) NOT NULL
);


--
-- Name: GenreRef_GenreID_seq; Type: SEQUENCE; Schema: References; Owner: -
--

ALTER TABLE "References"."GenreRef" ALTER COLUMN "GenreID" ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME "References"."GenreRef_GenreID_seq"
    START WITH 1
    INCREMENT BY 1
    MINVALUE 0
    NO MAXVALUE
    CACHE 1
);


--
-- Name: LanguageRef; Type: TABLE; Schema: References; Owner: -
--

CREATE TABLE "References"."LanguageRef" (
    "LanguageID" smallint NOT NULL,
    "LanguageName" character varying(255) NOT NULL,
    "LanguageCode" character varying(255) NOT NULL
);


--
-- Name: LanguageRef_LanguageID_seq; Type: SEQUENCE; Schema: References; Owner: -
--

ALTER TABLE "References"."LanguageRef" ALTER COLUMN "LanguageID" ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME "References"."LanguageRef_LanguageID_seq"
    START WITH 0
    INCREMENT BY 1
    MINVALUE 0
    NO MAXVALUE
    CACHE 1
);


--
-- Name: ParentGuideRef; Type: TABLE; Schema: References; Owner: -
--

CREATE TABLE "References"."ParentGuideRef" (
    "ParentGuideID" smallint NOT NULL,
    "ParentGuideDescription" character varying(255) NOT NULL
);


--
-- Name: QualityRef; Type: TABLE; Schema: References; Owner: -
--

CREATE TABLE "References"."QualityRef" (
    "QualityID" smallint NOT NULL,
    "QualityName" character varying(255)
);


--
-- Name: RecordRef; Type: TABLE; Schema: References; Owner: -
--

CREATE TABLE "References"."RecordRef" (
    "RecordID" smallint NOT NULL,
    "RecordType" character varying(255) NOT NULL
);


--
-- Name: TitleTypeRef; Type: TABLE; Schema: References; Owner: -
--

CREATE TABLE "References"."TitleTypeRef" (
    "TypeID" smallint NOT NULL,
    "TypeName" character varying(255) NOT NULL
);


--
-- Name: CastTable; Type: TABLE; Schema: Tables; Owner: -
--

CREATE TABLE "Tables"."CastTable" (
    "CastID" bigint NOT NULL,
    "CastName" character varying(255) NOT NULL,
    "CastImageURL" character varying(255),
    "IsDirector" boolean,
    "IsWriter" boolean,
    "IsCharacter" boolean,
    "CastDescription" character varying(255)
);


--
-- Name: CompanyTable; Type: TABLE; Schema: Tables; Owner: -
--

CREATE TABLE "Tables"."CompanyTable" (
    "CompanyID" integer NOT NULL,
    "CompanyName" character varying(255) NOT NULL
);


--
-- Name: NotDownloaded; Type: TABLE; Schema: Tables; Owner: -
--

CREATE TABLE "Tables"."NotDownloaded" (
    "TitleID" integer NOT NULL
);


--
-- Name: RequestedTitles; Type: TABLE; Schema: Tables; Owner: -
--

CREATE TABLE "Tables"."RequestedTitles" (
    "TitleID" integer NOT NULL
);


--
-- Name: TitleTable; Type: TABLE; Schema: Tables; Owner: -
--

CREATE TABLE "Tables"."TitleTable" (
    "TitleID" integer NOT NULL,
    "TitleType" smallint NOT NULL,
    "TitleName" character varying(255) NOT NULL,
    "TitleYear" smallint NOT NULL,
    "TitleYearTxt" character varying(255),
    "FolderName" character varying(255) NOT NULL,
    "FolderPath" character varying(255),
    "PosterURL" character varying(255),
    "OriginalTitle" character varying(255),
    "TitleLength" smallint,
    "DateReleased" timestamp(6) without time zone DEFAULT now(),
    "MetacriticRating" smallint,
    "Revenue" bigint,
    "IMDbRating" numeric(4,1),
    "IMDbVotes" integer,
    "Popularity" bigint,
    "ParentID" integer,
    "ParentName" character varying(255),
    "ParentYear" character varying(255),
    "EpisodeSeason" character varying(255),
    "EpisodeNumber" integer,
    "PreviousTitleID" integer,
    "NextTitleID" integer,
    "TotalSeasons" smallint,
    "TotalEpisodes" smallint,
    "TitleSummary" text,
    "TitleStoryLine" text,
    "TitleCertificate" smallint,
    "TitleCategory" smallint,
    "Nudity" smallint,
    "Violence" smallint,
    "Profanity" smallint,
    "AlcoholDrugSmoking" smallint,
    "Frightening" smallint,
    "TitleCountry" smallint,
    "Available" boolean DEFAULT false,
    "DateAdded" timestamp(6) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "DateUpdated" timestamp(6) without time zone DEFAULT now(),
    "Viewed" bigint DEFAULT 0,
    "Played" bigint DEFAULT 0,
    "Liked" bigint DEFAULT 0,
    "UnLiked" bigint DEFAULT 0,
    "PosterDownloaded" boolean DEFAULT false,
    "TitleLanguage" smallint
);


--
-- Name: SelectAvailable; Type: MATERIALIZED VIEW; Schema: Tables; Owner: -
--

CREATE MATERIALIZED VIEW "Tables"."SelectAvailable" AS
 SELECT DISTINCT "TitleTable"."TitleID",
    "TitleTable"."TitleName",
    "TitleTable"."TitleYear",
    "TitleTable"."FolderName",
    "TitleTable"."TitleType",
    lower(("TitleTable"."OriginalTitle")::text) AS "OriginalTitle",
    "TitleTable"."TitleLength",
    "TitleTable"."IMDbRating",
    "TitleTable"."Popularity",
    "TitleTable"."TitleYearTxt",
    "TitleTable"."TitleCountry" AS "Nationality",
    "TitleTable"."DateAdded",
    "TitleTable"."Viewed",
    "TitleTable"."Played",
    "TitleTable"."Liked",
    lower(("CastTable"."CastName")::text) AS "CastName",
    "GenreTitleLine"."GenreID",
    "LanguageRef"."LanguageCode",
    "LanguageRef"."LanguageName",
    "CountryRef"."CountryName",
    "CountryRef"."CountryCode",
    "AwardEventRef"."EventName"
   FROM ((((((((("Tables"."TitleTable"
     JOIN "Lines"."CastTitleLine" ON (("TitleTable"."TitleID" = "CastTitleLine"."TitleID")))
     JOIN "Tables"."CastTable" ON (("CastTitleLine"."CastID" = "CastTable"."CastID")))
     JOIN "Lines"."GenreTitleLine" ON (("TitleTable"."TitleID" = "GenreTitleLine"."TitleID")))
     JOIN "Lines"."LanguageTitleLine" ON (("TitleTable"."TitleID" = "LanguageTitleLine"."TitleID")))
     JOIN "References"."LanguageRef" ON (("LanguageTitleLine"."LanguageID" = "LanguageRef"."LanguageID")))
     JOIN "Lines"."CountryTitleLine" ON (("TitleTable"."TitleID" = "CountryTitleLine"."TitleID")))
     JOIN "References"."CountryRef" ON ((("CountryTitleLine"."CountryID" = "CountryRef"."CountryID") AND ("TitleTable"."TitleCountry" = "CountryRef"."CountryID"))))
     LEFT JOIN "Lines"."AwardTitleLine" ON ((("CastTable"."CastID" = "AwardTitleLine"."CastID") AND ("TitleTable"."TitleID" = "AwardTitleLine"."TitleID"))))
     LEFT JOIN "References"."AwardEventRef" ON (("AwardTitleLine"."EventID" = "AwardEventRef"."EventID")))
  WHERE (("TitleTable"."Available" = true) AND ("TitleTable"."Popularity" IS NOT NULL))
  WITH NO DATA;


--
-- Name: SelectNotAvailable; Type: MATERIALIZED VIEW; Schema: Tables; Owner: -
--

CREATE MATERIALIZED VIEW "Tables"."SelectNotAvailable" AS
 SELECT DISTINCT "TitleTable"."TitleID",
    "TitleTable"."TitleName",
    "TitleTable"."TitleYear",
    "TitleTable"."FolderName",
    "TitleTable"."TitleLength",
    "TitleTable"."IMDbRating",
    "TitleTable"."Popularity",
    "TitleTable"."TitleYearTxt",
    "TitleTable"."TitleCountry" AS "Nationality",
    "TitleTable"."DateAdded",
    "TitleTable"."Viewed",
    "TitleTable"."Played",
    "TitleTable"."Liked",
    lower(("CastTable"."CastName")::text) AS "CastName",
    "GenreTitleLine"."GenreID",
    "LanguageRef"."LanguageCode",
    "LanguageRef"."LanguageName",
    "CountryRef"."CountryName",
    "CountryRef"."CountryCode",
    "AwardEventRef"."EventName"
   FROM ((((((((("Tables"."TitleTable"
     JOIN "Lines"."CastTitleLine" ON (("TitleTable"."TitleID" = "CastTitleLine"."TitleID")))
     JOIN "Tables"."CastTable" ON (("CastTitleLine"."CastID" = "CastTable"."CastID")))
     JOIN "Lines"."GenreTitleLine" ON (("TitleTable"."TitleID" = "GenreTitleLine"."TitleID")))
     JOIN "Lines"."LanguageTitleLine" ON (("TitleTable"."TitleID" = "LanguageTitleLine"."TitleID")))
     JOIN "References"."LanguageRef" ON (("LanguageTitleLine"."LanguageID" = "LanguageRef"."LanguageID")))
     JOIN "Lines"."CountryTitleLine" ON (("TitleTable"."TitleID" = "CountryTitleLine"."TitleID")))
     JOIN "References"."CountryRef" ON ((("CountryTitleLine"."CountryID" = "CountryRef"."CountryID") AND ("TitleTable"."TitleCountry" = "CountryRef"."CountryID"))))
     LEFT JOIN "Lines"."AwardTitleLine" ON ((("CastTable"."CastID" = "AwardTitleLine"."CastID") AND ("TitleTable"."TitleID" = "AwardTitleLine"."TitleID"))))
     LEFT JOIN "References"."AwardEventRef" ON (("AwardTitleLine"."EventID" = "AwardEventRef"."EventID")))
  WHERE (("TitleTable"."Available" = false) AND ("TitleTable"."Popularity" IS NOT NULL))
  WITH NO DATA;


--
-- Name: ToBeUpdated; Type: TABLE; Schema: Tables; Owner: -
--

CREATE TABLE "Tables"."ToBeUpdated" (
    "TitleID" integer NOT NULL
);


--
-- Name: TotalSearch; Type: VIEW; Schema: Tables; Owner: -
--

CREATE VIEW "Tables"."TotalSearch" AS
 SELECT "TitleTable"."TitleID",
    "TitleTable"."TitleName",
    "TitleTable"."TitleYear",
    "TitleTable"."FolderName",
    lower(("TitleTable"."OriginalTitle")::text) AS "OriginalTitle",
    "TitleTable"."TitleLength",
    "TitleTable"."IMDbRating",
    "TitleTable"."Popularity",
    "TitleTable"."TitleYearTxt",
    "TitleTable"."TitleCountry" AS "Nationality",
    "TitleTable"."DateAdded",
    "TitleTable"."Viewed",
    "TitleTable"."Played",
    "TitleTable"."Liked",
    lower(("CastTable"."CastName")::text) AS "CastName",
    "GenreTitleLine"."GenreID",
    "LanguageRef"."LanguageCode",
    "LanguageRef"."LanguageName",
    "CountryRef"."CountryName",
    "CountryRef"."CountryCode",
    "AwardEventRef"."EventName"
   FROM ((((((((("Tables"."TitleTable"
     JOIN "Lines"."CastTitleLine" ON (("TitleTable"."TitleID" = "CastTitleLine"."TitleID")))
     JOIN "Tables"."CastTable" ON (("CastTitleLine"."CastID" = "CastTable"."CastID")))
     JOIN "Lines"."GenreTitleLine" ON (("TitleTable"."TitleID" = "GenreTitleLine"."TitleID")))
     JOIN "Lines"."LanguageTitleLine" ON (("TitleTable"."TitleID" = "LanguageTitleLine"."TitleID")))
     JOIN "References"."LanguageRef" ON (("LanguageTitleLine"."LanguageID" = "LanguageRef"."LanguageID")))
     JOIN "Lines"."CountryTitleLine" ON (("TitleTable"."TitleID" = "CountryTitleLine"."TitleID")))
     JOIN "References"."CountryRef" ON ((("CountryTitleLine"."CountryID" = "CountryRef"."CountryID") AND ("TitleTable"."TitleCountry" = "CountryRef"."CountryID"))))
     LEFT JOIN "Lines"."AwardTitleLine" ON ((("CastTable"."CastID" = "AwardTitleLine"."CastID") AND ("TitleTable"."TitleID" = "AwardTitleLine"."TitleID"))))
     LEFT JOIN "References"."AwardEventRef" ON (("AwardTitleLine"."EventID" = "AwardEventRef"."EventID")))
  WHERE (("TitleTable"."Available" = true) AND ("TitleTable"."Popularity" IS NOT NULL));


--
-- Name: CompanyTable; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public."CompanyTable" (
    "CompanyID" integer,
    "CompanyName" character varying(255)
);


--
-- Name: AwardTitleLine AwardTitleLine_pkey; Type: CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."AwardTitleLine"
    ADD CONSTRAINT "AwardTitleLine_pkey" PRIMARY KEY ("TitleID", "EventID", "CastID", "Description", "Category");


--
-- Name: CastTitleLine CastTitleLine_pkey; Type: CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."CastTitleLine"
    ADD CONSTRAINT "CastTitleLine_pkey" PRIMARY KEY ("TitleID", "CastID", "CastType");


--
-- Name: CertificateTitleLine CertificateTitleLine_pkey; Type: CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."CertificateTitleLine"
    ADD CONSTRAINT "CertificateTitleLine_pkey" PRIMARY KEY ("TitleID", "CountryID", "CertificateID");


--
-- Name: CompanyTitleLine CompanyTitleLine_pkey; Type: CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."CompanyTitleLine"
    ADD CONSTRAINT "CompanyTitleLine_pkey" PRIMARY KEY ("TitleID", "CompanyID");


--
-- Name: ConnectionTitleLine ConnectionTitleLine_pkey; Type: CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."ConnectionTitleLine"
    ADD CONSTRAINT "ConnectionTitleLine_pkey" PRIMARY KEY ("TitleID", "ConnectionTitleID", "ConnectionType");


--
-- Name: CountryTitleLine CountryTitleLine_pkey; Type: CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."CountryTitleLine"
    ADD CONSTRAINT "CountryTitleLine_pkey" PRIMARY KEY ("TitleID", "CountryID");


--
-- Name: FileTitleLine FileTitleLine_pkey; Type: CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."FileTitleLine"
    ADD CONSTRAINT "FileTitleLine_pkey" PRIMARY KEY ("TitleID", "QualityID", "DisplayID", "AudioLanguageID", "SubtitleLanguageID");


--
-- Name: GenreTitleLine GenreTitleLine_pkey; Type: CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."GenreTitleLine"
    ADD CONSTRAINT "GenreTitleLine_pkey" PRIMARY KEY ("TitleID", "GenreID");


--
-- Name: KnownAsTitleLine KnownAsTitleLine_pkey; Type: CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."KnownAsTitleLine"
    ADD CONSTRAINT "KnownAsTitleLine_pkey" PRIMARY KEY ("TitleID", "KnownAs");


--
-- Name: LanguageTitleLine LanguageTitleLine_pkey; Type: CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."LanguageTitleLine"
    ADD CONSTRAINT "LanguageTitleLine_pkey" PRIMARY KEY ("TitleID", "LanguageID");


--
-- Name: SimilaritiesTitleLine SimilaritiesTitleLine_pkey; Type: CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."SimilaritiesTitleLine"
    ADD CONSTRAINT "SimilaritiesTitleLine_pkey" PRIMARY KEY ("TitleID", "SimilarTitleID");


--
-- Name: MonitorPackage2 AppMonitor1_copy1_pkey; Type: CONSTRAINT; Schema: MonitorPackages; Owner: -
--

ALTER TABLE ONLY "MonitorPackages"."MonitorPackage2"
    ADD CONSTRAINT "AppMonitor1_copy1_pkey" PRIMARY KEY ("Application");


--
-- Name: MonitorPackage3 AppMonitor1_copy1_pkey1; Type: CONSTRAINT; Schema: MonitorPackages; Owner: -
--

ALTER TABLE ONLY "MonitorPackages"."MonitorPackage3"
    ADD CONSTRAINT "AppMonitor1_copy1_pkey1" PRIMARY KEY ("Application");


--
-- Name: MonitorPackage4 AppMonitor2_copy1_pkey; Type: CONSTRAINT; Schema: MonitorPackages; Owner: -
--

ALTER TABLE ONLY "MonitorPackages"."MonitorPackage4"
    ADD CONSTRAINT "AppMonitor2_copy1_pkey" PRIMARY KEY ("Application");


--
-- Name: MonitorPackage1 AppMonitor_pkey; Type: CONSTRAINT; Schema: MonitorPackages; Owner: -
--

ALTER TABLE ONLY "MonitorPackages"."MonitorPackage1"
    ADD CONSTRAINT "AppMonitor_pkey" PRIMARY KEY ("Application");


--
-- Name: AwardEventRef AwardEventRef_pkey; Type: CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."AwardEventRef"
    ADD CONSTRAINT "AwardEventRef_pkey" PRIMARY KEY ("EventID");


--
-- Name: AwardNominationTypeRef AwardNominationTypeRef_pkey; Type: CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."AwardNominationTypeRef"
    ADD CONSTRAINT "AwardNominationTypeRef_pkey" PRIMARY KEY ("NominationTypeID");


--
-- Name: CastTypeRef CastTypeRef_pkey; Type: CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."CastTypeRef"
    ADD CONSTRAINT "CastTypeRef_pkey" PRIMARY KEY ("CastTypeID");


--
-- Name: CategoryRef CategoryRef_pkey; Type: CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."CategoryRef"
    ADD CONSTRAINT "CategoryRef_pkey" PRIMARY KEY ("CategoryID");


--
-- Name: CertificateCountryRef CertificateCountryRef_pkey; Type: CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."CertificateCountryRef"
    ADD CONSTRAINT "CertificateCountryRef_pkey" PRIMARY KEY ("CountryID", "CertificateID");


--
-- Name: CertificateRef CertificateRef_CertificateName_key; Type: CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."CertificateRef"
    ADD CONSTRAINT "CertificateRef_CertificateName_key" UNIQUE ("CertificateName");


--
-- Name: CertificateRef CertificateRef_pkey; Type: CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."CertificateRef"
    ADD CONSTRAINT "CertificateRef_pkey" PRIMARY KEY ("CertificateID");


--
-- Name: ConnectionTypeRef ConnectionTypeRef_ConnectionTypeDescription_key; Type: CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."ConnectionTypeRef"
    ADD CONSTRAINT "ConnectionTypeRef_ConnectionTypeDescription_key" UNIQUE ("ConnectionTypeDescription");


--
-- Name: ConnectionTypeRef ConnectionTypeRef_pkey; Type: CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."ConnectionTypeRef"
    ADD CONSTRAINT "ConnectionTypeRef_pkey" PRIMARY KEY ("ConnectionTypeID");


--
-- Name: CountryRef CountryRef_pkey; Type: CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."CountryRef"
    ADD CONSTRAINT "CountryRef_pkey" PRIMARY KEY ("CountryID");


--
-- Name: DisplayRef DisplayRef_pkey; Type: CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."DisplayRef"
    ADD CONSTRAINT "DisplayRef_pkey" PRIMARY KEY ("DisplayID");


--
-- Name: GenreRef GenreRef_GenreName_key; Type: CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."GenreRef"
    ADD CONSTRAINT "GenreRef_GenreName_key" UNIQUE ("GenreName");


--
-- Name: GenreRef GenreRef_pkey; Type: CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."GenreRef"
    ADD CONSTRAINT "GenreRef_pkey" PRIMARY KEY ("GenreID");


--
-- Name: LanguageRef LanguageRef_LanguageCode_key; Type: CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."LanguageRef"
    ADD CONSTRAINT "LanguageRef_LanguageCode_key" UNIQUE ("LanguageCode");


--
-- Name: LanguageRef LanguageRef_LanguageName_key; Type: CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."LanguageRef"
    ADD CONSTRAINT "LanguageRef_LanguageName_key" UNIQUE ("LanguageName");


--
-- Name: LanguageRef LanguageRef_pkey; Type: CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."LanguageRef"
    ADD CONSTRAINT "LanguageRef_pkey" PRIMARY KEY ("LanguageID");


--
-- Name: ParentGuideRef ParentGuideRef_pkey; Type: CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."ParentGuideRef"
    ADD CONSTRAINT "ParentGuideRef_pkey" PRIMARY KEY ("ParentGuideID");


--
-- Name: QualityRef QualityRef_pkey; Type: CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."QualityRef"
    ADD CONSTRAINT "QualityRef_pkey" PRIMARY KEY ("QualityID");


--
-- Name: RecordRef RecordRef_RecordType_key; Type: CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."RecordRef"
    ADD CONSTRAINT "RecordRef_RecordType_key" UNIQUE ("RecordType");


--
-- Name: RecordRef RecordRef_pkey; Type: CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."RecordRef"
    ADD CONSTRAINT "RecordRef_pkey" PRIMARY KEY ("RecordID");


--
-- Name: TitleTypeRef TitleTypeRef_TypeName_key; Type: CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."TitleTypeRef"
    ADD CONSTRAINT "TitleTypeRef_TypeName_key" UNIQUE ("TypeName");


--
-- Name: TitleTypeRef TitleTypeRef_pkey; Type: CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."TitleTypeRef"
    ADD CONSTRAINT "TitleTypeRef_pkey" PRIMARY KEY ("TypeID");


--
-- Name: CastTable CastTable_pkey; Type: CONSTRAINT; Schema: Tables; Owner: -
--

ALTER TABLE ONLY "Tables"."CastTable"
    ADD CONSTRAINT "CastTable_pkey" PRIMARY KEY ("CastID");


--
-- Name: CompanyTable CompanyTable_pkey; Type: CONSTRAINT; Schema: Tables; Owner: -
--

ALTER TABLE ONLY "Tables"."CompanyTable"
    ADD CONSTRAINT "CompanyTable_pkey" PRIMARY KEY ("CompanyID");


--
-- Name: NotDownloaded NotDownloaded_pkey; Type: CONSTRAINT; Schema: Tables; Owner: -
--

ALTER TABLE ONLY "Tables"."NotDownloaded"
    ADD CONSTRAINT "NotDownloaded_pkey" PRIMARY KEY ("TitleID");


--
-- Name: RequestedTitles RequestedDownload_pkey; Type: CONSTRAINT; Schema: Tables; Owner: -
--

ALTER TABLE ONLY "Tables"."RequestedTitles"
    ADD CONSTRAINT "RequestedDownload_pkey" PRIMARY KEY ("TitleID");


--
-- Name: RequestedTitles RequestedTitles_TitleID_key; Type: CONSTRAINT; Schema: Tables; Owner: -
--

ALTER TABLE ONLY "Tables"."RequestedTitles"
    ADD CONSTRAINT "RequestedTitles_TitleID_key" UNIQUE ("TitleID");


--
-- Name: TitleTable TitleInfo_pkey; Type: CONSTRAINT; Schema: Tables; Owner: -
--

ALTER TABLE ONLY "Tables"."TitleTable"
    ADD CONSTRAINT "TitleInfo_pkey" PRIMARY KEY ("TitleID");


--
-- Name: ToBeUpdated ToBeUpdated_pkey; Type: CONSTRAINT; Schema: Tables; Owner: -
--

ALTER TABLE ONLY "Tables"."ToBeUpdated"
    ADD CONSTRAINT "ToBeUpdated_pkey" PRIMARY KEY ("TitleID");


--
-- Name: AwardTitleLine_TitleID_EventID_CastID_Description_Category_idx; Type: INDEX; Schema: Lines; Owner: -
--

CREATE UNIQUE INDEX "AwardTitleLine_TitleID_EventID_CastID_Description_Category_idx" ON "Lines"."AwardTitleLine" USING btree ("TitleID", "EventID", "CastID", "Description", "Category");

ALTER TABLE "Lines"."AwardTitleLine" CLUSTER ON "AwardTitleLine_TitleID_EventID_CastID_Description_Category_idx";


--
-- Name: CastTitleLine_CastID_idx; Type: INDEX; Schema: Lines; Owner: -
--

CREATE INDEX "CastTitleLine_CastID_idx" ON "Lines"."CastTitleLine" USING btree ("CastID");


--
-- Name: CastTitleLine_TitleID_CastID_CastType_idx; Type: INDEX; Schema: Lines; Owner: -
--

CREATE UNIQUE INDEX "CastTitleLine_TitleID_CastID_CastType_idx" ON "Lines"."CastTitleLine" USING btree ("TitleID", "CastID", "CastType");

ALTER TABLE "Lines"."CastTitleLine" CLUSTER ON "CastTitleLine_TitleID_CastID_CastType_idx";


--
-- Name: CastTitleLine_TitleID_idx; Type: INDEX; Schema: Lines; Owner: -
--

CREATE INDEX "CastTitleLine_TitleID_idx" ON "Lines"."CastTitleLine" USING btree ("TitleID");


--
-- Name: CertificateTitleLine_TitleID_CountryID_CertificateID_idx; Type: INDEX; Schema: Lines; Owner: -
--

CREATE UNIQUE INDEX "CertificateTitleLine_TitleID_CountryID_CertificateID_idx" ON "Lines"."CertificateTitleLine" USING btree ("TitleID", "CountryID", "CertificateID");

ALTER TABLE "Lines"."CertificateTitleLine" CLUSTER ON "CertificateTitleLine_TitleID_CountryID_CertificateID_idx";


--
-- Name: CompanyTitleLine_TitleID_CompanyID_idx; Type: INDEX; Schema: Lines; Owner: -
--

CREATE UNIQUE INDEX "CompanyTitleLine_TitleID_CompanyID_idx" ON "Lines"."CompanyTitleLine" USING btree ("TitleID", "CompanyID");

ALTER TABLE "Lines"."CompanyTitleLine" CLUSTER ON "CompanyTitleLine_TitleID_CompanyID_idx";


--
-- Name: ConnectionTitleLine_TitleID_ConnectionTitleID_ConnectionTyp_idx; Type: INDEX; Schema: Lines; Owner: -
--

CREATE UNIQUE INDEX "ConnectionTitleLine_TitleID_ConnectionTitleID_ConnectionTyp_idx" ON "Lines"."ConnectionTitleLine" USING btree ("TitleID", "ConnectionTitleID", "ConnectionType");

ALTER TABLE "Lines"."ConnectionTitleLine" CLUSTER ON "ConnectionTitleLine_TitleID_ConnectionTitleID_ConnectionTyp_idx";


--
-- Name: CountryTitleLine_TitleID_CountryID_idx; Type: INDEX; Schema: Lines; Owner: -
--

CREATE UNIQUE INDEX "CountryTitleLine_TitleID_CountryID_idx" ON "Lines"."CountryTitleLine" USING btree ("TitleID", "CountryID");

ALTER TABLE "Lines"."CountryTitleLine" CLUSTER ON "CountryTitleLine_TitleID_CountryID_idx";


--
-- Name: FileTitleLine_TitleID_QualityID_DisplayID_AudioLanguageID_S_idx; Type: INDEX; Schema: Lines; Owner: -
--

CREATE UNIQUE INDEX "FileTitleLine_TitleID_QualityID_DisplayID_AudioLanguageID_S_idx" ON "Lines"."FileTitleLine" USING btree ("TitleID", "QualityID", "DisplayID", "AudioLanguageID", "SubtitleLanguageID");

ALTER TABLE "Lines"."FileTitleLine" CLUSTER ON "FileTitleLine_TitleID_QualityID_DisplayID_AudioLanguageID_S_idx";


--
-- Name: GenreTitleLine_TitleID_GenreID_idx; Type: INDEX; Schema: Lines; Owner: -
--

CREATE UNIQUE INDEX "GenreTitleLine_TitleID_GenreID_idx" ON "Lines"."GenreTitleLine" USING btree ("TitleID", "GenreID");

ALTER TABLE "Lines"."GenreTitleLine" CLUSTER ON "GenreTitleLine_TitleID_GenreID_idx";


--
-- Name: KnownAsTitleLine_TitleID_KnownAs_idx; Type: INDEX; Schema: Lines; Owner: -
--

CREATE UNIQUE INDEX "KnownAsTitleLine_TitleID_KnownAs_idx" ON "Lines"."KnownAsTitleLine" USING btree ("TitleID", "KnownAs");

ALTER TABLE "Lines"."KnownAsTitleLine" CLUSTER ON "KnownAsTitleLine_TitleID_KnownAs_idx";


--
-- Name: LanguageTitleLine_TitleID_LanguageID_idx; Type: INDEX; Schema: Lines; Owner: -
--

CREATE UNIQUE INDEX "LanguageTitleLine_TitleID_LanguageID_idx" ON "Lines"."LanguageTitleLine" USING btree ("TitleID", "LanguageID");

ALTER TABLE "Lines"."LanguageTitleLine" CLUSTER ON "LanguageTitleLine_TitleID_LanguageID_idx";


--
-- Name: SimilaritiesTitleLine_TitleID_SimilarTitleID_idx; Type: INDEX; Schema: Lines; Owner: -
--

CREATE UNIQUE INDEX "SimilaritiesTitleLine_TitleID_SimilarTitleID_idx" ON "Lines"."SimilaritiesTitleLine" USING btree ("TitleID", "SimilarTitleID");

ALTER TABLE "Lines"."SimilaritiesTitleLine" CLUSTER ON "SimilaritiesTitleLine_TitleID_SimilarTitleID_idx";


--
-- Name: AwardEventRef_EventID_idx; Type: INDEX; Schema: References; Owner: -
--

CREATE UNIQUE INDEX "AwardEventRef_EventID_idx" ON "References"."AwardEventRef" USING btree ("EventID");


--
-- Name: AwardNominationTypeRef_NominationTypeID_idx; Type: INDEX; Schema: References; Owner: -
--

CREATE UNIQUE INDEX "AwardNominationTypeRef_NominationTypeID_idx" ON "References"."AwardNominationTypeRef" USING btree ("NominationTypeID");


--
-- Name: CastTypeRef_CastTypeID_idx; Type: INDEX; Schema: References; Owner: -
--

CREATE UNIQUE INDEX "CastTypeRef_CastTypeID_idx" ON "References"."CastTypeRef" USING btree ("CastTypeID");


--
-- Name: CategoryRef_CategoryID_idx; Type: INDEX; Schema: References; Owner: -
--

CREATE UNIQUE INDEX "CategoryRef_CategoryID_idx" ON "References"."CategoryRef" USING btree ("CategoryID");


--
-- Name: CertificateCountryLine_CountryID_CertificateID_idx; Type: INDEX; Schema: References; Owner: -
--

CREATE UNIQUE INDEX "CertificateCountryLine_CountryID_CertificateID_idx" ON "References"."CertificateCountryRef" USING btree ("CountryID", "CertificateID");


--
-- Name: CertificateRef_CertificateID_idx; Type: INDEX; Schema: References; Owner: -
--

CREATE UNIQUE INDEX "CertificateRef_CertificateID_idx" ON "References"."CertificateRef" USING btree ("CertificateID");


--
-- Name: ConnectionTypeRef_ConnectionTypeDescription_idx; Type: INDEX; Schema: References; Owner: -
--

CREATE UNIQUE INDEX "ConnectionTypeRef_ConnectionTypeDescription_idx" ON "References"."ConnectionTypeRef" USING btree ("ConnectionTypeDescription");


--
-- Name: ConnectionTypeRef_ConnectionTypeID_idx; Type: INDEX; Schema: References; Owner: -
--

CREATE UNIQUE INDEX "ConnectionTypeRef_ConnectionTypeID_idx" ON "References"."ConnectionTypeRef" USING btree ("ConnectionTypeID");


--
-- Name: CountryRef_CountryCode_idx; Type: INDEX; Schema: References; Owner: -
--

CREATE UNIQUE INDEX "CountryRef_CountryCode_idx" ON "References"."CountryRef" USING btree ("CountryCode");


--
-- Name: CountryRef_CountryName_idx; Type: INDEX; Schema: References; Owner: -
--

CREATE UNIQUE INDEX "CountryRef_CountryName_idx" ON "References"."CountryRef" USING btree ("CountryName");


--
-- Name: DisplayRef_DisplayID_idx; Type: INDEX; Schema: References; Owner: -
--

CREATE UNIQUE INDEX "DisplayRef_DisplayID_idx" ON "References"."DisplayRef" USING btree ("DisplayID");


--
-- Name: GenreRef_GenreID_idx; Type: INDEX; Schema: References; Owner: -
--

CREATE UNIQUE INDEX "GenreRef_GenreID_idx" ON "References"."GenreRef" USING btree ("GenreID");


--
-- Name: GenreRef_GenreName_idx; Type: INDEX; Schema: References; Owner: -
--

CREATE UNIQUE INDEX "GenreRef_GenreName_idx" ON "References"."GenreRef" USING btree ("GenreName");


--
-- Name: LanguageRef_LanguageCode_idx; Type: INDEX; Schema: References; Owner: -
--

CREATE UNIQUE INDEX "LanguageRef_LanguageCode_idx" ON "References"."LanguageRef" USING btree ("LanguageCode");


--
-- Name: LanguageRef_LanguageID_idx; Type: INDEX; Schema: References; Owner: -
--

CREATE UNIQUE INDEX "LanguageRef_LanguageID_idx" ON "References"."LanguageRef" USING btree ("LanguageID");


--
-- Name: ParentGuideRef_ParentGuideID_idx; Type: INDEX; Schema: References; Owner: -
--

CREATE UNIQUE INDEX "ParentGuideRef_ParentGuideID_idx" ON "References"."ParentGuideRef" USING btree ("ParentGuideID");


--
-- Name: QualityRef_QualityID_idx; Type: INDEX; Schema: References; Owner: -
--

CREATE UNIQUE INDEX "QualityRef_QualityID_idx" ON "References"."QualityRef" USING btree ("QualityID");


--
-- Name: RecordRef_RecordID_idx; Type: INDEX; Schema: References; Owner: -
--

CREATE UNIQUE INDEX "RecordRef_RecordID_idx" ON "References"."RecordRef" USING btree ("RecordID");


--
-- Name: RecordRef_RecordType_idx; Type: INDEX; Schema: References; Owner: -
--

CREATE UNIQUE INDEX "RecordRef_RecordType_idx" ON "References"."RecordRef" USING btree ("RecordType");


--
-- Name: TitleTypeRef_TypeID_idx; Type: INDEX; Schema: References; Owner: -
--

CREATE UNIQUE INDEX "TitleTypeRef_TypeID_idx" ON "References"."TitleTypeRef" USING btree ("TypeID");


--
-- Name: TitleTypeRef_TypeName_idx; Type: INDEX; Schema: References; Owner: -
--

CREATE UNIQUE INDEX "TitleTypeRef_TypeName_idx" ON "References"."TitleTypeRef" USING btree ("TypeName");


--
-- Name: CastTable_CastID_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE UNIQUE INDEX "CastTable_CastID_idx" ON "Tables"."CastTable" USING btree ("CastID");


--
-- Name: CastTable_CastName_Lower_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "CastTable_CastName_Lower_idx" ON "Tables"."CastTable" USING spgist (lower(("CastName")::text));


--
-- Name: CastTable_CastName_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "CastTable_CastName_idx" ON "Tables"."CastTable" USING spgist ("CastName");


--
-- Name: CastTable_STRUCTURED__idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "CastTable_STRUCTURED__idx" ON "Tables"."CastTable" USING btree ("CastID", "CastName", "CastImageURL", "IsDirector", "IsWriter", "IsCharacter", "CastDescription");

ALTER TABLE "Tables"."CastTable" CLUSTER ON "CastTable_STRUCTURED__idx";


--
-- Name: CompanyTable_CompanyID_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE UNIQUE INDEX "CompanyTable_CompanyID_idx" ON "Tables"."CompanyTable" USING btree ("CompanyID");


--
-- Name: NotDownloaded_TitleID_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE UNIQUE INDEX "NotDownloaded_TitleID_idx" ON "Tables"."NotDownloaded" USING btree ("TitleID");


--
-- Name: RequestedDownload_TitleID_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE UNIQUE INDEX "RequestedDownload_TitleID_idx" ON "Tables"."RequestedTitles" USING btree ("TitleID");


--
-- Name: SelectAvailable_CastName_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "SelectAvailable_CastName_idx" ON "Tables"."SelectAvailable" USING spgist ("CastName");


--
-- Name: SelectAvailable_Clustered_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE UNIQUE INDEX "SelectAvailable_Clustered_idx" ON "Tables"."SelectAvailable" USING btree ("TitleName", "FolderName", "TitleType", "OriginalTitle", "IMDbRating", "Popularity", "Nationality", "CastName", "GenreID", "LanguageCode", "LanguageName", "CountryName", "CountryCode", "EventName", "TitleID", "TitleYear", "TitleLength", "TitleYearTxt", "DateAdded", "Viewed", "Liked", "Played");

ALTER TABLE "Tables"."SelectAvailable" CLUSTER ON "SelectAvailable_Clustered_idx";


--
-- Name: SelectAvailable_CountExpression_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "SelectAvailable_CountExpression_idx" ON "Tables"."SelectAvailable" USING btree (lower(("FolderName")::text), "TitleType", "OriginalTitle");


--
-- Name: SelectAvailable_OriginalTitle_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "SelectAvailable_OriginalTitle_idx" ON "Tables"."SelectAvailable" USING spgist ("OriginalTitle");


--
-- Name: SelectAvailable_Popularity_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "SelectAvailable_Popularity_idx" ON "Tables"."SelectAvailable" USING btree ("Popularity");


--
-- Name: SelectAvailable_TitleID_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "SelectAvailable_TitleID_idx" ON "Tables"."SelectAvailable" USING btree ("TitleID");


--
-- Name: SelectAvailable_TitleName_Lower_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "SelectAvailable_TitleName_Lower_idx" ON "Tables"."SelectAvailable" USING spgist (lower(("TitleName")::text));


--
-- Name: SelectAvailable_TitleType_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "SelectAvailable_TitleType_idx" ON "Tables"."SelectAvailable" USING btree ("TitleType");


--
-- Name: SelectNotAvailable_Clustered_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "SelectNotAvailable_Clustered_idx" ON "Tables"."SelectNotAvailable" USING btree ("TitleID", "TitleName", "TitleYear", "FolderName", "TitleLength", "IMDbRating", "Popularity", "TitleYearTxt", "Nationality", "DateAdded", "Viewed", "Played", "Liked", "CastName", "GenreID", "LanguageCode", "LanguageName", "CountryName", "CountryCode", "EventName");

ALTER TABLE "Tables"."SelectNotAvailable" CLUSTER ON "SelectNotAvailable_Clustered_idx";


--
-- Name: TitleTable_Available_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "TitleTable_Available_idx" ON "Tables"."TitleTable" USING btree ("Available");


--
-- Name: TitleTable_DateAdded_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "TitleTable_DateAdded_idx" ON "Tables"."TitleTable" USING btree ("DateAdded");


--
-- Name: TitleTable_DateReleased_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "TitleTable_DateReleased_idx" ON "Tables"."TitleTable" USING btree ("DateReleased");


--
-- Name: TitleTable_DateUpdated_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "TitleTable_DateUpdated_idx" ON "Tables"."TitleTable" USING btree ("DateUpdated");


--
-- Name: TitleTable_FolderName_Lower_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "TitleTable_FolderName_Lower_idx" ON "Tables"."TitleTable" USING spgist (lower(("FolderName")::text));


--
-- Name: TitleTable_FolderName_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "TitleTable_FolderName_idx" ON "Tables"."TitleTable" USING spgist ("FolderName");


--
-- Name: TitleTable_IMDbRating_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "TitleTable_IMDbRating_idx" ON "Tables"."TitleTable" USING btree ("IMDbRating");


--
-- Name: TitleTable_Nationality_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "TitleTable_Nationality_idx" ON "Tables"."TitleTable" USING btree ("TitleCountry");


--
-- Name: TitleTable_OriginalTitle_Lower_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "TitleTable_OriginalTitle_Lower_idx" ON "Tables"."TitleTable" USING spgist (lower(("OriginalTitle")::text));


--
-- Name: TitleTable_Popularity_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "TitleTable_Popularity_idx" ON "Tables"."TitleTable" USING btree ("Popularity");


--
-- Name: TitleTable_PosterDownloaded_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "TitleTable_PosterDownloaded_idx" ON "Tables"."TitleTable" USING btree ("PosterDownloaded");


--
-- Name: TitleTable_PosterURL_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "TitleTable_PosterURL_idx" ON "Tables"."TitleTable" USING btree ("PosterURL");


--
-- Name: TitleTable_TitleCategory_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "TitleTable_TitleCategory_idx" ON "Tables"."TitleTable" USING btree ("TitleCategory");


--
-- Name: TitleTable_TitleCertificate_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "TitleTable_TitleCertificate_idx" ON "Tables"."TitleTable" USING btree ("TitleCertificate");


--
-- Name: TitleTable_TitleID_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE UNIQUE INDEX "TitleTable_TitleID_idx" ON "Tables"."TitleTable" USING btree ("TitleID");

ALTER TABLE "Tables"."TitleTable" CLUSTER ON "TitleTable_TitleID_idx";


--
-- Name: TitleTable_TitleName_Lower_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "TitleTable_TitleName_Lower_idx" ON "Tables"."TitleTable" USING spgist (lower(("TitleName")::text));


--
-- Name: TitleTable_TitleName_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "TitleTable_TitleName_idx" ON "Tables"."TitleTable" USING spgist ("TitleName");


--
-- Name: TitleTable_TitleType_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "TitleTable_TitleType_idx" ON "Tables"."TitleTable" USING btree ("TitleType");


--
-- Name: TitleTable_TitleYearTxt_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "TitleTable_TitleYearTxt_idx" ON "Tables"."TitleTable" USING spgist ("TitleYearTxt");


--
-- Name: TitleTable_TitleYear_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "TitleTable_TitleYear_idx" ON "Tables"."TitleTable" USING btree ("TitleYear");


--
-- Name: ToBeUpdated_TitleID_idx; Type: INDEX; Schema: Tables; Owner: -
--

CREATE INDEX "ToBeUpdated_TitleID_idx" ON "Tables"."ToBeUpdated" USING btree ("TitleID");


--
-- Name: AwardTitleLine AwardTitleLine_CastID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."AwardTitleLine"
    ADD CONSTRAINT "AwardTitleLine_CastID_fkey" FOREIGN KEY ("CastID") REFERENCES "Tables"."CastTable"("CastID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: AwardTitleLine AwardTitleLine_EventID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."AwardTitleLine"
    ADD CONSTRAINT "AwardTitleLine_EventID_fkey" FOREIGN KEY ("EventID") REFERENCES "References"."AwardEventRef"("EventID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: AwardTitleLine AwardTitleLine_NominationType_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."AwardTitleLine"
    ADD CONSTRAINT "AwardTitleLine_NominationType_fkey" FOREIGN KEY ("NominationType") REFERENCES "References"."AwardNominationTypeRef"("NominationTypeID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: AwardTitleLine AwardTitleLine_TitleID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."AwardTitleLine"
    ADD CONSTRAINT "AwardTitleLine_TitleID_fkey" FOREIGN KEY ("TitleID") REFERENCES "Tables"."TitleTable"("TitleID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: CastTitleLine CastTitleLine_CastID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."CastTitleLine"
    ADD CONSTRAINT "CastTitleLine_CastID_fkey" FOREIGN KEY ("CastID") REFERENCES "Tables"."CastTable"("CastID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: CastTitleLine CastTitleLine_CastType_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."CastTitleLine"
    ADD CONSTRAINT "CastTitleLine_CastType_fkey" FOREIGN KEY ("CastType") REFERENCES "References"."CastTypeRef"("CastTypeID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: CastTitleLine CastTitleLine_TitleID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."CastTitleLine"
    ADD CONSTRAINT "CastTitleLine_TitleID_fkey" FOREIGN KEY ("TitleID") REFERENCES "Tables"."TitleTable"("TitleID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: CertificateTitleLine CertificateTitleLine_CertificateID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."CertificateTitleLine"
    ADD CONSTRAINT "CertificateTitleLine_CertificateID_fkey" FOREIGN KEY ("CertificateID") REFERENCES "References"."CertificateRef"("CertificateID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: CertificateTitleLine CertificateTitleLine_CountryID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."CertificateTitleLine"
    ADD CONSTRAINT "CertificateTitleLine_CountryID_fkey" FOREIGN KEY ("CountryID") REFERENCES "References"."CountryRef"("CountryID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: CertificateTitleLine CertificateTitleLine_TitleID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."CertificateTitleLine"
    ADD CONSTRAINT "CertificateTitleLine_TitleID_fkey" FOREIGN KEY ("TitleID") REFERENCES "Tables"."TitleTable"("TitleID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: CompanyTitleLine CompanyTitleLine_CompanyID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."CompanyTitleLine"
    ADD CONSTRAINT "CompanyTitleLine_CompanyID_fkey" FOREIGN KEY ("CompanyID") REFERENCES "Tables"."CompanyTable"("CompanyID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: CompanyTitleLine CompanyTitleLine_TitleID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."CompanyTitleLine"
    ADD CONSTRAINT "CompanyTitleLine_TitleID_fkey" FOREIGN KEY ("TitleID") REFERENCES "Tables"."TitleTable"("TitleID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: ConnectionTitleLine ConnectionTitleLine_ConnectionType_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."ConnectionTitleLine"
    ADD CONSTRAINT "ConnectionTitleLine_ConnectionType_fkey" FOREIGN KEY ("ConnectionType") REFERENCES "References"."ConnectionTypeRef"("ConnectionTypeID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: ConnectionTitleLine ConnectionTitleLine_TitleID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."ConnectionTitleLine"
    ADD CONSTRAINT "ConnectionTitleLine_TitleID_fkey" FOREIGN KEY ("TitleID") REFERENCES "Tables"."TitleTable"("TitleID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: CountryTitleLine CountryTitleLine_CountryID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."CountryTitleLine"
    ADD CONSTRAINT "CountryTitleLine_CountryID_fkey" FOREIGN KEY ("CountryID") REFERENCES "References"."CountryRef"("CountryID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: CountryTitleLine CountryTitleLine_TitleID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."CountryTitleLine"
    ADD CONSTRAINT "CountryTitleLine_TitleID_fkey" FOREIGN KEY ("TitleID") REFERENCES "Tables"."TitleTable"("TitleID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: FileTitleLine FileTitleLine_AudioLanguageID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."FileTitleLine"
    ADD CONSTRAINT "FileTitleLine_AudioLanguageID_fkey" FOREIGN KEY ("AudioLanguageID") REFERENCES "References"."LanguageRef"("LanguageID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: FileTitleLine FileTitleLine_DisplayID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."FileTitleLine"
    ADD CONSTRAINT "FileTitleLine_DisplayID_fkey" FOREIGN KEY ("DisplayID") REFERENCES "References"."DisplayRef"("DisplayID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: FileTitleLine FileTitleLine_QualityID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."FileTitleLine"
    ADD CONSTRAINT "FileTitleLine_QualityID_fkey" FOREIGN KEY ("QualityID") REFERENCES "References"."QualityRef"("QualityID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: FileTitleLine FileTitleLine_SubtitleLanguageID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."FileTitleLine"
    ADD CONSTRAINT "FileTitleLine_SubtitleLanguageID_fkey" FOREIGN KEY ("SubtitleLanguageID") REFERENCES "References"."LanguageRef"("LanguageID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: FileTitleLine FileTitleLine_TitleID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."FileTitleLine"
    ADD CONSTRAINT "FileTitleLine_TitleID_fkey" FOREIGN KEY ("TitleID") REFERENCES "Tables"."TitleTable"("TitleID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: GenreTitleLine GenreTitleLine_GenreID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."GenreTitleLine"
    ADD CONSTRAINT "GenreTitleLine_GenreID_fkey" FOREIGN KEY ("GenreID") REFERENCES "References"."GenreRef"("GenreID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: GenreTitleLine GenreTitleLine_TitleID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."GenreTitleLine"
    ADD CONSTRAINT "GenreTitleLine_TitleID_fkey" FOREIGN KEY ("TitleID") REFERENCES "Tables"."TitleTable"("TitleID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: KnownAsTitleLine KnownAsTitleLine_TitleID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."KnownAsTitleLine"
    ADD CONSTRAINT "KnownAsTitleLine_TitleID_fkey" FOREIGN KEY ("TitleID") REFERENCES "Tables"."TitleTable"("TitleID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: LanguageTitleLine LanguageTitleLine_LanguageID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."LanguageTitleLine"
    ADD CONSTRAINT "LanguageTitleLine_LanguageID_fkey" FOREIGN KEY ("LanguageID") REFERENCES "References"."LanguageRef"("LanguageID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: LanguageTitleLine LanguageTitleLine_TitleID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."LanguageTitleLine"
    ADD CONSTRAINT "LanguageTitleLine_TitleID_fkey" FOREIGN KEY ("TitleID") REFERENCES "Tables"."TitleTable"("TitleID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: SimilaritiesTitleLine SimilaritiesTitleLine_TitleID_fkey; Type: FK CONSTRAINT; Schema: Lines; Owner: -
--

ALTER TABLE ONLY "Lines"."SimilaritiesTitleLine"
    ADD CONSTRAINT "SimilaritiesTitleLine_TitleID_fkey" FOREIGN KEY ("TitleID") REFERENCES "Tables"."TitleTable"("TitleID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: CertificateCountryRef CertificateCountryRef_CertificateID_fkey; Type: FK CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."CertificateCountryRef"
    ADD CONSTRAINT "CertificateCountryRef_CertificateID_fkey" FOREIGN KEY ("CertificateID") REFERENCES "References"."CertificateRef"("CertificateID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: CertificateCountryRef CertificateCountryRef_CountryID_fkey; Type: FK CONSTRAINT; Schema: References; Owner: -
--

ALTER TABLE ONLY "References"."CertificateCountryRef"
    ADD CONSTRAINT "CertificateCountryRef_CountryID_fkey" FOREIGN KEY ("CountryID") REFERENCES "References"."CountryRef"("CountryID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: TitleTable TitleInfo_AlcoholDrugSmoking_fkey; Type: FK CONSTRAINT; Schema: Tables; Owner: -
--

ALTER TABLE ONLY "Tables"."TitleTable"
    ADD CONSTRAINT "TitleInfo_AlcoholDrugSmoking_fkey" FOREIGN KEY ("AlcoholDrugSmoking") REFERENCES "References"."ParentGuideRef"("ParentGuideID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: TitleTable TitleInfo_Frightening_fkey; Type: FK CONSTRAINT; Schema: Tables; Owner: -
--

ALTER TABLE ONLY "Tables"."TitleTable"
    ADD CONSTRAINT "TitleInfo_Frightening_fkey" FOREIGN KEY ("Frightening") REFERENCES "References"."ParentGuideRef"("ParentGuideID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: TitleTable TitleInfo_Nudity_fkey; Type: FK CONSTRAINT; Schema: Tables; Owner: -
--

ALTER TABLE ONLY "Tables"."TitleTable"
    ADD CONSTRAINT "TitleInfo_Nudity_fkey" FOREIGN KEY ("Nudity") REFERENCES "References"."ParentGuideRef"("ParentGuideID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: TitleTable TitleInfo_Profanity_fkey; Type: FK CONSTRAINT; Schema: Tables; Owner: -
--

ALTER TABLE ONLY "Tables"."TitleTable"
    ADD CONSTRAINT "TitleInfo_Profanity_fkey" FOREIGN KEY ("Profanity") REFERENCES "References"."ParentGuideRef"("ParentGuideID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: TitleTable TitleInfo_TitleCategory_fkey; Type: FK CONSTRAINT; Schema: Tables; Owner: -
--

ALTER TABLE ONLY "Tables"."TitleTable"
    ADD CONSTRAINT "TitleInfo_TitleCategory_fkey" FOREIGN KEY ("TitleCategory") REFERENCES "References"."CategoryRef"("CategoryID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: TitleTable TitleInfo_TitleCertificate_fkey; Type: FK CONSTRAINT; Schema: Tables; Owner: -
--

ALTER TABLE ONLY "Tables"."TitleTable"
    ADD CONSTRAINT "TitleInfo_TitleCertificate_fkey" FOREIGN KEY ("TitleCertificate") REFERENCES "References"."CertificateRef"("CertificateID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: TitleTable TitleInfo_TitleType_fkey; Type: FK CONSTRAINT; Schema: Tables; Owner: -
--

ALTER TABLE ONLY "Tables"."TitleTable"
    ADD CONSTRAINT "TitleInfo_TitleType_fkey" FOREIGN KEY ("TitleType") REFERENCES "References"."TitleTypeRef"("TypeID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: TitleTable TitleInfo_Violence_fkey; Type: FK CONSTRAINT; Schema: Tables; Owner: -
--

ALTER TABLE ONLY "Tables"."TitleTable"
    ADD CONSTRAINT "TitleInfo_Violence_fkey" FOREIGN KEY ("Violence") REFERENCES "References"."ParentGuideRef"("ParentGuideID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: TitleTable TitleTable_Nationality_fkey; Type: FK CONSTRAINT; Schema: Tables; Owner: -
--

ALTER TABLE ONLY "Tables"."TitleTable"
    ADD CONSTRAINT "TitleTable_Nationality_fkey" FOREIGN KEY ("TitleCountry") REFERENCES "References"."CountryRef"("CountryID") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict Oqp8baNyTzS4OksCUCPHnSlnh5WESUbg5klTUPotz7nFgWiGFe5UIYqJZIYFNv3

