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
		&model.Superadmin{},
		&model.UsersWithVisibleAnswers{},
		&model.Submission{},
		&model.Answer{},
		&model.SuperadminPermission{},
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
			-- Drop all tables and views
			FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
				EXECUTE 'DROP TABLE IF EXISTS public.' || quote_ident(r.tablename) || ' CASCADE';
			END LOOP;
		END $$;`

	err := db.Exec(statement).Error

	if err != nil {
		return fmt.Errorf("failed to create enum type: %v", err)
	}

	types := []string{
		"gender_enum",
		"questionnaires_status_enum",
		"questionnaires_sequence_enum",
		"questionnaires_visibility_enum",
		"questions_type_enum",
		"submissions_status_enum",
	}

	// Iterate and delete each type
	for _, t := range types {
		statement := fmt.Sprintf(`DROP TYPE IF EXISTS public.%s CASCADE`, t)
		if err := db.Exec(statement).Error; err != nil {
			return fmt.Errorf("failed to delete type %s: %v", t, err)
		}
	}

	return nil
}
