package migrations

import (
	"context"
	"fmt"

	"github.com/dhax/go-base/models"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] add moto core tables...")

		// Register the junction model explicitly
		db.RegisterModel((*models.CombinedGroupGroups)(nil))

		// Create time-related tables
		if _, err := db.NewCreateTable().
			Model((*models.Timespan)(nil)).
			IfNotExists().
			Exec(ctx); err != nil {
			return err
		}

		if _, err := db.NewCreateTable().
			Model((*models.Datespan)(nil)).
			IfNotExists().
			Exec(ctx); err != nil {
			return err
		}

		// Create room table
		if _, err := db.NewCreateTable().
			Model((*models.Room)(nil)).
			IfNotExists().
			Exec(ctx); err != nil {
			return err
		}

		// Create user-related tables
		if _, err := db.NewCreateTable().
			Model((*models.CustomUser)(nil)).
			IfNotExists().
			Exec(ctx); err != nil {
			return err
		}

		if _, err := db.NewCreateTable().
			Model((*models.PedagogicalSpecialist)(nil)).
			IfNotExists().
			ForeignKey(`("custom_user_id") REFERENCES "custom_users" ("id") ON DELETE CASCADE`).
			Exec(ctx); err != nil {
			return err
		}

		// Create group tables
		if _, err := db.NewCreateTable().
			Model((*models.Group)(nil)).
			IfNotExists().
			ForeignKey(`("room_id") REFERENCES "rooms" ("id") ON DELETE SET NULL`).
			ForeignKey(`("representative_id") REFERENCES "pedagogical_specialists" ("id") ON DELETE SET NULL`).
			Exec(ctx); err != nil {
			return err
		}

		if _, err := db.NewCreateTable().
			Model((*models.CombinedGroup)(nil)).
			IfNotExists().
			ForeignKey(`("specific_group_id") REFERENCES "groups" ("id") ON DELETE SET NULL`).
			Exec(ctx); err != nil {
			return err
		}

		// Create student tables
		if _, err := db.NewCreateTable().
			Model((*models.Student)(nil)).
			IfNotExists().
			ForeignKey(`("custom_user_id") REFERENCES "custom_users" ("id") ON DELETE CASCADE`).
			ForeignKey(`("group_id") REFERENCES "groups" ("id") ON DELETE CASCADE`).
			Exec(ctx); err != nil {
			return err
		}

		if _, err := db.NewCreateTable().
			Model((*models.Visit)(nil)).
			IfNotExists().
			ForeignKey(`("student_id") REFERENCES "students" ("id") ON DELETE CASCADE`).
			ForeignKey(`("room_id") REFERENCES "rooms" ("id") ON DELETE CASCADE`).
			ForeignKey(`("timespan_id") REFERENCES "timespans" ("id") ON DELETE CASCADE`).
			Exec(ctx); err != nil {
			return err
		}

		if _, err := db.NewCreateTable().
			Model((*models.Feedback)(nil)).
			IfNotExists().
			ForeignKey(`("student_id") REFERENCES "students" ("id") ON DELETE CASCADE`).
			Exec(ctx); err != nil {
			return err
		}

		// Create M2M junction table for CombinedGroup and Group
		if _, err := db.NewCreateTable().
			Model((*models.CombinedGroupGroups)(nil)).
			IfNotExists().
			ForeignKey(`("combined_group_id") REFERENCES "combined_groups" ("id") ON DELETE CASCADE`).
			ForeignKey(`("group_id") REFERENCES "groups" ("id") ON DELETE CASCADE`).
			Exec(ctx); err != nil {
			return err
		}

		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] drop moto core tables...")

		// Drop tables in reverse order
		if _, err := db.NewDropTable().
			Model((*models.Feedback)(nil)).
			IfExists().
			Cascade().
			Exec(ctx); err != nil {
			return err
		}

		if _, err := db.NewDropTable().
			Model((*models.Visit)(nil)).
			IfExists().
			Cascade().
			Exec(ctx); err != nil {
			return err
		}

		if _, err := db.NewDropTable().
			Model((*models.Student)(nil)).
			IfExists().
			Cascade().
			Exec(ctx); err != nil {
			return err
		}

		if _, err := db.NewDropTable().
			Model((*models.CombinedGroupGroups)(nil)).
			IfExists().
			Cascade().
			Exec(ctx); err != nil {
			return err
		}

		if _, err := db.NewDropTable().
			Model((*models.CombinedGroup)(nil)).
			IfExists().
			Cascade().
			Exec(ctx); err != nil {
			return err
		}

		if _, err := db.NewDropTable().
			Model((*models.Group)(nil)).
			IfExists().
			Cascade().
			Exec(ctx); err != nil {
			return err
		}

		if _, err := db.NewDropTable().
			Model((*models.PedagogicalSpecialist)(nil)).
			IfExists().
			Cascade().
			Exec(ctx); err != nil {
			return err
		}

		if _, err := db.NewDropTable().
			Model((*models.CustomUser)(nil)).
			IfExists().
			Cascade().
			Exec(ctx); err != nil {
			return err
		}

		if _, err := db.NewDropTable().
			Model((*models.Room)(nil)).
			IfExists().
			Cascade().
			Exec(ctx); err != nil {
			return err
		}

		if _, err := db.NewDropTable().
			Model((*models.Datespan)(nil)).
			IfExists().
			Cascade().
			Exec(ctx); err != nil {
			return err
		}

		if _, err := db.NewDropTable().
			Model((*models.Timespan)(nil)).
			IfExists().
			Cascade().
			Exec(ctx); err != nil {
			return err
		}

		return nil
	})
}
