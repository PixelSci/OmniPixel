## 1. Dependencies

- [x] 1.1 Promote gorm.io/gorm and gorm.io/driver/postgres to direct dependencies in go.mod, remove pgx/v5 direct dependency

## 2. Domain models

- [x] 2.1 Add GORM struct tags to domain.User (table name, primary key, column mappings)
- [x] 2.2 Add GORM struct tags to domain.Conversation (table name, primary key, column mappings)
- [x] 2.3 Add GORM struct tags to domain.Message (table name, primary key, column mappings)

## 3. Database connection

- [x] 3.1 Rewrite bootstrap/db.go to create *gorm.DB via gorm.Open instead of pgxpool.New
- [x] 3.2 Add AutoMigrate call in the db bootstrap with all domain models
- [x] 3.3 Remove or repurpose db/postgres.go type alias

## 4. Repositories

- [x] 4.1 Rewrite UserRepository to use *gorm.DB and GORM query methods
- [x] 4.2 Rewrite ConversationRepository to use *gorm.DB and GORM query methods

## 5. Wiring

- [x] 5.1 Update bootstrap/wire.go to use *gorm.DB instead of *pgxpool.Pool in Providers setup

## 6. Cleanup

- [x] 6.1 Run go mod tidy to clean up unused indirect dependencies
- [x] 6.2 Verify the application compiles and starts correctly
