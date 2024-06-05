-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table "public"."filemeta" drop constraint if exists "filemeta_tigertree_id_fkey";
alter table "public"."task" drop constraint if exists "fk_task_token_id_refs_token";
alter table "public"."tigertree" drop constraint if exists "tigertree_tree_type_full_path_host_operation_id_key";
alter table "public"."tigertree" drop constraint if exists "tigertree_callback_id_fkey";
alter table "public"."tigertree" drop constraint if exists "tigertree_callback_id_tree_type_full_path_host_operation_id_ke";
drop index if exists "public"."tigertree_tree_type_full_path_host_operation_id_key";
drop index if exists "public"."tigertree_callback_id_tree_type_full_path_host_operation_id_ke";
alter table "public"."callback" add column if not exists "tigertree_groups" text[] not null default '{Default}'::text[];
alter table "public"."tigertree" add column if not exists "callback_id" integer;

CREATE UNIQUE INDEX tigertree_callback_id_tree_type_full_path_host_operation_id_ke ON public.tigertree USING btree (callback_id, tree_type, full_path, host, operation_id);
alter table "public"."tigertree" add constraint "tigertree_callback_id_fkey" FOREIGN KEY (callback_id) REFERENCES callback(id) ON UPDATE SET NULL ON DELETE SET NULL not valid;
alter table "public"."tigertree" validate constraint "tigertree_callback_id_fkey";
alter table "public"."tigertree" add constraint "tigertree_callback_id_tree_type_full_path_host_operation_id_ke" UNIQUE using index "tigertree_callback_id_tree_type_full_path_host_operation_id_ke";

alter table "public"."task" add constraint "task_token_id_fkey" FOREIGN KEY (token_id) REFERENCES token(id) ON UPDATE SET NULL ON DELETE SET NULL not valid;
alter table "public"."task" validate constraint "task_token_id_fkey";
alter table "public"."filemeta" add constraint "filemeta_tigertree_id_fkey" FOREIGN KEY (tigertree_id) REFERENCES tigertree(id) ON UPDATE SET NULL ON DELETE SET NULL not valid;
alter table "public"."filemeta" validate constraint "filemeta_tigertree_id_fkey";

-- +migrate Down
-- SQL in section 'Down' is executed when this migration is rolled back
alter table "public"."c2profileparametersinstance" drop constraint if exists "c2profileparametersinstance_c2_profile_parameters_id_operation_";

drop index if exists "public"."c2profileparametersinstance_c2_profile_parameters_id_operation_";

