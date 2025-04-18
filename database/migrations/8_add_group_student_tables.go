package migrations

import (
	"context"
	"fmt"

	"github.com/dhax/go-base/models"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] add group, student, visit, and feedback tables...")

		// Register models with junction tables for many-to-many relationships
		db.RegisterModel((*models.GroupSupervisor)(nil))
		db.RegisterModel((*models.CombinedGroupGroup)(nil))
		db.RegisterModel((*models.CombinedGroupSpecialist)(nil))

		// Create groups table
		_, err := db.NewCreateTable().
			Model((*models.Group)(nil)).
			IfNotExists().
			ForeignKey(`("room_id") REFERENCES "rooms" ("id") ON DELETE SET NULL`).
			ForeignKey(`("representative_id") REFERENCES "pedagogical_specialists" ("id") ON DELETE SET NULL`).
			Exec(ctx)
		if err != nil {
			return err
		}

		// Create group_supervisors junction table
		_, err = db.NewCreateTable().
			Model((*models.GroupSupervisor)(nil)).
			IfNotExists().
			ForeignKey(`("group_id") REFERENCES "groups" ("id") ON DELETE CASCADE`).
			ForeignKey(`("specialist_id") REFERENCES "pedagogical_specialists" ("id") ON DELETE CASCADE`).
			Exec(ctx)
		if err != nil {
			return err
		}

		// Create combined_groups table
		_, err = db.NewCreateTable().
			Model((*models.CombinedGroup)(nil)).
			IfNotExists().
			ForeignKey(`("specific_group_id") REFERENCES "groups" ("id") ON DELETE SET NULL`).
			Exec(ctx)
		if err != nil {
			return err
		}

		// Create combined_group_groups junction table
		_, err = db.NewCreateTable().
			Model((*models.CombinedGroupGroup)(nil)).
			IfNotExists().
			ForeignKey(`("combinedgroup_id") REFERENCES "combined_groups" ("id") ON DELETE CASCADE`).
			ForeignKey(`("group_id") REFERENCES "groups" ("id") ON DELETE CASCADE`).
			Exec(ctx)
		if err != nil {
			return err
		}

		// Create combined_group_specialists junction table
		_, err = db.NewCreateTable().
			Model((*models.CombinedGroupSpecialist)(nil)).
			IfNotExists().
			ForeignKey(`("combinedgroup_id") REFERENCES "combined_groups" ("id") ON DELETE CASCADE`).
			ForeignKey(`("specialist_id") REFERENCES "pedagogical_specialists" ("id") ON DELETE CASCADE`).
			Exec(ctx)
		if err != nil {
			return err
		}

		// Create students table
		_, err = db.NewCreateTable().
			Model((*models.Student)(nil)).
			IfNotExists().
			ForeignKey(`("custom_user_id") REFERENCES "custom_users" ("id") ON DELETE CASCADE`).
			ForeignKey(`("group_id") REFERENCES "groups" ("id") ON DELETE CASCADE`).
			Exec(ctx)
		if err != nil {
			return err
		}

		// Create visits table
		_, err = db.NewCreateTable().
			Model((*models.Visit)(nil)).
			IfNotExists().
			ForeignKey(`("student_id") REFERENCES "students" ("id") ON DELETE CASCADE`).
			ForeignKey(`("room_id") REFERENCES "rooms" ("id") ON DELETE CASCADE`).
			ForeignKey(`("timespan_id") REFERENCES "timespans" ("id") ON DELETE CASCADE`).
			Exec(ctx)
		if err != nil {
			return err
		}

		// Create feedback table
		_, err = db.NewCreateTable().
			Model((*models.Feedback)(nil)).
			IfNotExists().
			ForeignKey(`("student_id") REFERENCES "students" ("id") ON DELETE CASCADE`).
			Exec(ctx)
		if err != nil {
			return err
		}

		fmt.Println(" done")
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] drop group and student related tables...")

		// Drop tables in reverse order to handle foreign key constraints
		_, err := db.NewDropTable().
			Model((*models.Feedback)(nil)).
			IfExists().
			Cascade().
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewDropTable().
			Model((*models.Visit)(nil)).
			IfExists().
			Cascade().
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewDropTable().
			Model((*models.Student)(nil)).
			IfExists().
			Cascade().
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewDropTable().
			Model((*models.CombinedGroupSpecialist)(nil)).
			IfExists().
			Cascade().
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewDropTable().
			Model((*models.CombinedGroupGroup)(nil)).
			IfExists().
			Cascade().
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewDropTable().
			Model((*models.CombinedGroup)(nil)).
			IfExists().
			Cascade().
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewDropTable().
			Model((*models.GroupSupervisor)(nil)).
			IfExists().
			Cascade().
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewDropTable().
			Model((*models.Group)(nil)).
			IfExists().
			Cascade().
			Exec(ctx)
		if err != nil {
			return err
		}

		fmt.Println(" done")
		return nil
	})
}
