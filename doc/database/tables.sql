CREATE TYPE gender_enum AS ENUM ('male', 'female');
CREATE TABLE IF NOT EXISTS "users" (
  "id" serial primary key,
  "first_name" varchar null,
  "last_name" varchar null,
  "birthdate" varchar null,
  "city" varchar null,
  "national_code" varchar(10) unique,
  "gender" gender_enum null,
  "email" varchar unique,
  "password" varchar,
  "credit" bigint,
  "created_at" timestamp DEFAULT now()
);

CREATE TYPE questionnaires_status_enum AS ENUM ('open', 'closed', 'cancelled');
CREATE TYPE questionnaires_sequence_enum AS ENUM ('random', 'sequential');
CREATE TYPE questionnaires_visibility_enum AS ENUM ('everybody', 'admin_and_owner', 'nobody');

CREATE TABLE IF NOT EXISTS "questionnaires" (
  "id" serial primary key,
  "owner_id" bigint,
  "status" questionnaires_status_enum DEFAULT 'open',
  "can_submit_from" timestamp,
  "can_submit_until" timestamp,
  "max_minutes_to_response" int,
  "max_minutes_to_change_answer" int,
  "max_minutes_to_giveback_answer" int,
  "random_or_sequential" questionnaires_sequence_enum,
  "can_back_to_previous_question" bool,
  "title" varchar,
  "max_allowed_submissions_count" int,
  "answers_visibile_for" questionnaires_visibility_enum,
  "created_at" timestamp DEFAULT now(),
  CONSTRAINT "FK_questionnaires.owner_id"
    FOREIGN KEY ("owner_id")
      REFERENCES "users"("id")
);

CREATE TABLE IF NOT EXISTS "notifications" (
  "id" serial primary key,
  "user_id" bigint,
  "title" varchar,
  "description" varchar,
  "is_seen" bool DEFAULT false,
  "created_at" timestamp DEFAULT now(),
  CONSTRAINT "FK_notifications.user_id"
    FOREIGN KEY ("user_id")
      REFERENCES "users"("id")
);

CREATE TABLE IF NOT EXISTS "super_admins" (
  "id" serial primary key,
  "user_id" bigint,
  "granted_by" bigint null,
  "created_at" timestamp DEFAULT now(),
  CONSTRAINT "FK_super_admins.user_id"
    FOREIGN KEY ("user_id")
      REFERENCES "users"("id")
);

CREATE TYPE questions_type_enum AS ENUM ('multioption', 'descriptive');
CREATE TABLE IF NOT EXISTS "questions" (
  "id" serial primary key,
  "title" varchar,
  "type" questions_type_enum,
  "questionnaire_id" bigint,
  "order" int,
  "file_path" varchar,
  "depends_on_question_id" bigint null,
  "depends_on_option_id" bigint null,
  "created_at" timestamp DEFAULT now(),
  CONSTRAINT "FK_questions.questionnaire_id"
    FOREIGN KEY ("questionnaire_id")
      REFERENCES "questionnaires"("id")
);

CREATE TABLE IF NOT EXISTS "options" (
  "id" serial primary key,
  "question_id" bigint,
  "order" int,
  "caption" varchar,
  "is_correct" bool null,
  CONSTRAINT "FK_options.question_id"
    FOREIGN KEY ("question_id")
      REFERENCES "questions"("id")
);


CREATE TYPE submissions_status_enum AS ENUM ('answering', 'submitted', 'cancelled', 'closed');
CREATE TABLE IF NOT EXISTS "submissions" (
  "id" serial primary key,
  "questionnaire_id" bigint,
  "user_id" bigint,
  "status" submissions_status_enum DEFAULT 'answering',
  "last_answered_question_id" bigint,
  "submitted_at" timestamp null,
  "spent_minutes" int null,
  CONSTRAINT "FK_submissions.questionnaire_id"
    FOREIGN KEY ("questionnaire_id")
      REFERENCES "questionnaires"("id"),
  CONSTRAINT "FK_submissions.user_id"
    FOREIGN KEY ("user_id")
      REFERENCES "users"("id")
);

CREATE TABLE IF NOT EXISTS "answers" (
  "id" serial primary key,
  "submission_id" bigint,
  "question_id" bigint,
  "user_id" bigint,
  "option_id" bigint null,
  "answer_text" varchar null,
  "created_at" timestamp DEFAULT now(),
  "updated_at" timestamp null,
  CONSTRAINT "FK_answers.question_id"
    FOREIGN KEY ("question_id")
      REFERENCES "questions"("id"),
  CONSTRAINT "FK_answers.option_id"
    FOREIGN KEY ("option_id")
      REFERENCES "options"("id"),
  CONSTRAINT "FK_answers.submission_id"
    FOREIGN KEY ("submission_id")
      REFERENCES "submissions"("id"),
  CONSTRAINT "FK_answers.user_id"
    FOREIGN KEY ("user_id")
      REFERENCES "users"("id")
);

CREATE TABLE IF NOT EXISTS "roles" (
  "id" serial primary key,
  "name" varchar unique
);

CREATE TABLE IF NOT EXISTS "role_user" (
  "id" serial primary key,
  "role_id" bigint,
  "user_id" bigint,
  CONSTRAINT "FK_role_user.role_id"
    FOREIGN KEY ("role_id")
      REFERENCES "roles"("id"),
  CONSTRAINT "FK_role_user.user_id"
    FOREIGN KEY ("user_id")
      REFERENCES "users"("id")
);

CREATE TABLE IF NOT EXISTS "role_permission" (
  "id" serial primary key,
  "questionnaire_id" bigint,
  "role_id" bigint,
  "permission_id" bigint,
  "expire_at" timestamp,
  CONSTRAINT "FK_role_permission.role_id"
    FOREIGN KEY ("role_id")
      REFERENCES "roles"("id"),
  CONSTRAINT "FK_role_permission.permission_id"
    FOREIGN KEY ("permission_id")
      REFERENCES "permissions"("id"),
  CONSTRAINT "FK_role_permission.user_id"
    FOREIGN KEY ("questionnaire_id")
      REFERENCES "questionnaires"("id")
);


CREATE TABLE IF NOT EXISTS "permissions" (
  "id" serial primary key,
  "name" varchar unique,
  "description" varchar(500)
);

CREATE TABLE IF NOT EXISTS "users_with_visible_answers" (
  "id" serial primary key,
  "role_permission_id" bigint,
  "user_id" bigint,
  CONSTRAINT "FK_users_with_visible_answers.role_permission_id"
    FOREIGN KEY ("role_permission_id")
      REFERENCES "role_permission"("id"),
  CONSTRAINT "FK_users_with_visible_answers.user_id"
    FOREIGN KEY ("user_id")
      REFERENCES "users"("id")
);


CREATE INDEX idx_questionnaire_id ON "submissions"("questionnaire_id");
CREATE INDEX idx_email ON "users"("email");