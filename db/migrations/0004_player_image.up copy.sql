-- initial schema: users and todos
CREATE TABLE image_files (
	id uuid NOT NULL,
	file_name varchar(255) NULL,
	"path" varchar(255) NULL,
	"type" varchar(255) NULL,
	resolutions json NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	CONSTRAINT image_files_pkey PRIMARY KEY (id)
);