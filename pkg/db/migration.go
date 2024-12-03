package db

import (
	"fmt"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	"gorm.io/gorm"
)

func migrate(gormDB *gorm.DB) error {
	err := insertTypesTable(gormDB)
	if err != nil {
		return err
	}

	return gormDB.AutoMigrate(
		&model.User{},
		&model.Questionnaire{},
		&model.Notification{},
		&model.Option{},
		&model.Permission{},
		&model.Role{},
		&model.RolePermission{},
		&model.RoleUser{},
		&model.SuperAdmin{},
		&model.UsersWithVisibleAnswers{},
		&model.Submission{},
		&model.Answer{},
		&model.TwoFACode{},
	)
}

func insertTypesTable(db *gorm.DB) error {
	statements := []string{
		// Create gender_enum
		`DO $$ BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'gender_enum') THEN
				CREATE TYPE gender_enum AS ENUM ('male', 'female');
			END IF;
		END $$;`,

		// Create questionnaires_status_enum
		`DO $$ BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'questionnaires_status_enum') THEN
				CREATE TYPE questionnaires_status_enum AS ENUM ('open', 'closed', 'cancelled');
			END IF;
		END $$;`,

		// Create questionnaires_sequence_enum
		`DO $$ BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'questionnaires_sequence_enum') THEN
				CREATE TYPE questionnaires_sequence_enum AS ENUM ('random', 'sequential');
			END IF;
		END $$;`,

		// Create questionnaires_visibility_enum
		`DO $$ BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'questionnaires_visibility_enum') THEN
				CREATE TYPE questionnaires_visibility_enum AS ENUM ('everybody', 'admin_and_owner', 'nobody');
			END IF;
		END $$;`,

		// Create questions_type_enum
		`DO $$ BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'questions_type_enum') THEN
				CREATE TYPE questions_type_enum AS ENUM ('multichoice', 'descriptive');
			END IF;
		END $$;`,

		// Create submissions_status_enum
		`DO $$ BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'submissions_status_enum') THEN
				CREATE TYPE submissions_status_enum AS ENUM ('answering', 'submitted', 'cancelled', 'closed');
			END IF;
		END $$;`,
	}

	for _, statement := range statements {
		err := db.Exec(statement).Error
		if err != nil {
			return fmt.Errorf("failed to create enum type: %v", err)
		}
	}
	return nil
}

func deleteAllTablesAndTypes(db *gorm.DB) error {
	statement := `DO $$ 
		DECLARE 
			r RECORD;
		BEGIN
			-- Drop all views
			FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
				EXECUTE 'DROP TABLE IF EXISTS public.' || r.tablename || ' CASCADE';
			END LOOP;

			-- Drop all types
			FOR r IN (SELECT typname FROM pg_type WHERE typnamespace = (SELECT oid FROM pg_catalog.pg_namespace WHERE nspname = 'public')) LOOP
				EXECUTE 'DROP TYPE IF EXISTS public.' || r.typname || ' CASCADE';
			END LOOP;
		END $$;`

	err := db.Exec(statement).Error
	if err != nil {
		return fmt.Errorf("failed to create enum type: %v", err)
	}

	return nil
}
