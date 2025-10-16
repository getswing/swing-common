-- initial schema: users and todos
CREATE TABLE public.players (
	id uuid NOT NULL,
	first_name varchar(255) NULL,
	last_name varchar(255) NULL,
	dob timestamptz NULL,
	referral_code varchar(255) NULL,
	"password" varchar(255) NULL,
	email varchar(255) NULL,
	phone_country_code varchar(255) NULL,
	phone_number varchar(255) NULL,
	nationality varchar(255) NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	experienced varchar(255) NULL,
	image_file_id uuid NULL,
	is_beta_tester bool NULL,
	is_onboarding bool NULL,
	email_verification_code varchar(255) NULL,
	gender varchar(255) DEFAULT 'NO_SPECIFIED'::character varying NULL,
	email_temporary varchar(255) NULL,
	delete_reason varchar(255) NULL,
	is_firstbuy bool NULL,
	username varchar(255) NULL,
	is_default_username bool DEFAULT true NULL,
	is_firstbuy_golf_course bool DEFAULT true NULL,
	onesignal_id uuid NULL,
	subscription_id text NULL,
	partner_subscriptions text NULL,
	CONSTRAINT players_pkey PRIMARY KEY (id)
);

ALTER TABLE players ADD CONSTRAINT players_image_file_id_fkey FOREIGN KEY (image_file_id) REFERENCES image_files(id);