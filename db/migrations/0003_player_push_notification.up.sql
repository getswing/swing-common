-- initial schema: users and todos
CREATE TABLE player_push_notifications (
	id uuid NOT NULL,
	"type" varchar(255) NULL,
	title text NULL,
	description text NULL,
	"module" varchar(255) NULL,
	target_id varchar(255) NULL,
	status varchar(255) NULL,
	player_id uuid NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	icon varchar(255) NULL,
	image_id uuid NULL,
	custom_notification_id uuid NULL,
	product_id uuid NULL,
	sent_by varchar(255) NULL,
	tournament_id uuid NULL,
	custom_avatar text NULL,
	"group" public."enum_player_push_notifications_group" DEFAULT 'PLAYER'::enum_player_push_notifications_group NULL,
	CONSTRAINT "player-push-notifications_pkey" PRIMARY KEY (id)
);

ALTER TABLE public.player_push_notifications ADD CONSTRAINT "player-push-notifications_custom_notification_id_fkey" FOREIGN KEY (custom_notification_id) REFERENCES public.custom_notifications(id);
ALTER TABLE public.player_push_notifications ADD CONSTRAINT "player-push-notifications_image_id_fkey" FOREIGN KEY (image_id) REFERENCES image_files(id);
ALTER TABLE public.player_push_notifications ADD CONSTRAINT "player-push-notifications_player_id_fkey" FOREIGN KEY (player_id) REFERENCES players(id);
