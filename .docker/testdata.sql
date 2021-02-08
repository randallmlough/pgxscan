DROP TABLE IF EXISTS "public"."test";

CREATE SEQUENCE IF NOT EXISTS test_id_seq;

CREATE TABLE "public"."test" (
    "id" int4 NOT NULL DEFAULT nextval('test_id_seq'::regclass),
    "int" int4,
    "int_8" int2,
    "int_16" int2,
    "int_32" int4,
    "int_64" int8,
    "uint" int4,
    "uint_8" int2,
    "uint_16" int2,
    "uint_32" int4,
    "uint_64" int8,
    "float_32" numeric,
    "float_64" numeric,
    "rune" int4,
    "byte" int2,
    "string" text,
    "bool" bool,
    "time" timestamp,
    "bytes" varchar(45),
    "string_slice" _text,
    "bool_slice" _bool,
    "int_slice" _int4,
    "float_slice" _numeric,
    "json" json,
    "json_b" jsonb,
    "map" jsonb,
    PRIMARY KEY ("id")
);

INSERT INTO "public"."test" ("id", "int", "int_8", "int_16", "int_32", "int_64", "uint", "uint_8", "uint_16", "uint_32", "uint_64", "float_32", "float_64", "rune", "byte", "string", "bool", "time", "bytes", "string_slice", "bool_slice", "int_slice", "float_slice", "json", "json_b", "map") VALUES
('1', '1', '121', '32761', '2147483641', '9223372036854775801', '11', '121', '32761', '2147483641', '9223372036854775801', '1.2100000381469727', '9715.631', '128512', '97', 'Hello world', 't', '2019-01-01 01:01:01', 'first row', '{cats,dogs}', '{t,f,f,t}', '{1,2,3,4,5}', '{1.2100000381469727,2.2100000381469727,3.2100000381469727,4.210000038146973}', '{"str":"I''m json","int":1,"embedded":{"data":false}}', '{"int": 1, "str": "I''m json b", "embedded": {"data": true}}', '{"key": "value"}'),
('2', '2', '122', '32762', '2147483642', '9223372036854775802', '12', '252', '32762', '2147483642', '9223372036854775802', '1.2200000286102295', '9715.632', '128514', '98', 'foo bar', 't', '2019-01-01 01:01:01', 'second row', '{"john doe","jane smith"}', '{f,f,f,t}', '{6,7,8,9,10}', '{5.210000038146973,6.210000038146973,7.210000038146973,8.210000038146973}', '{"str":"","int":0,"embedded":{"data":false}}', '{"int": 2, "str": "Hi", "embedded": {"data": true}}', '{"marco": "polo"}');

DROP TABLE IF EXISTS "public"."users";

CREATE SEQUENCE IF NOT EXISTS users_id_seq;

-- Table Definition
CREATE TABLE "public"."users" (
    "id" int4 NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    "name" varchar,
    "email" varchar,
    PRIMARY KEY ("id")
);

INSERT INTO "public"."users" ("id", "name", "email") VALUES
('1', 'user01', 'user01@email.com'),
('2', 'user02', 'user02@email.com'),
('3', 'user03', 'user03@email.com'),
('10', NULL, 'user03@email.com');


DROP TABLE IF EXISTS "public"."address";

CREATE SEQUENCE IF NOT EXISTS address_id_seq;

CREATE TABLE "public"."address" (
    "id" int4 NOT NULL DEFAULT nextval('address_id_seq'::regclass),
    "user_id" int4 NOT NULL,
    "line_1" varchar,
    "city" varchar,
    CONSTRAINT "fk_user_address" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id"),
    PRIMARY KEY ("id")
);

INSERT INTO "public"."address" ("id", "user_id", "line_1", "city") VALUES
('1', '1', 'line01_user01', 'city01'),
('2', '2', 'line02_user02', 'city02');
