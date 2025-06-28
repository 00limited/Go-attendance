package postgre

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DBOms             *gorm.DB
	DBMaster          *gorm.DB
	DBWms             *gorm.DB
	DBFin             *gorm.DB
	DBLog             *gorm.DB
	DBSfa             *gorm.DB
	DBManage          *gorm.DB
	DBMainOms         *gorm.DB
	DBMainMaster      *gorm.DB
	DBMainWms         *gorm.DB
	DBMainFin         *gorm.DB
	DBMainLog         *gorm.DB
	DBMainSfa         *gorm.DB
	DBMainManage      *gorm.DB
	DBPgBouncerOms    *gorm.DB
	DBPgBouncerMaster *gorm.DB
	DBPgBouncerWms    *gorm.DB
	DBPgBouncerFin    *gorm.DB
	DBPgBouncerLog    *gorm.DB
	DBPgBouncerSfa    *gorm.DB
	DBPgBouncerManage *gorm.DB
}

type postgreHelper struct {
	Username        string
	Password        string
	Host            string
	Port            int
	Name            string
	ApplicationName string
	TimeLocation    *time.Location
}

// Interface ...
type Interface interface {
	Connect() (*gorm.DB, error)
	ConnectionSQL() (*sql.DB, error)
	CloseConnection(db *gorm.DB) error
}

// NewPostgre ...
func NewPostgre(username string, password string, host string, port int, name string, applicationName string, timeLocation *time.Location) Interface {
	return &postgreHelper{
		Username:        username,
		Password:        password,
		Host:            host,
		Port:            port,
		Name:            name,
		ApplicationName: applicationName,
		TimeLocation:    timeLocation,
	}
}

// Connect ...
func (t *postgreHelper) Connect() (*gorm.DB, error) {
	// newLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	// 	logger.Config{
	// 		SlowThreshold:             time.Second,   // Slow SQL threshold
	// 		LogLevel:                  logger.Silent, // Log level
	// 		IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
	// 		Colorful:                  false,         // Disable color
	// 	},
	// )
	var dsn = fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable application_name=%v TimeZone=%v", t.Host, t.Username, t.Password, t.Name, t.Port, t.ApplicationName, t.TimeLocation.String())
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,  // data source name, refer https://github.com/jackc/pgx
		PreferSimpleProtocol: true, // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
	}), &gorm.Config{
		// Logger:               logger.Default.LogMode(logger.Silent),
		DisableAutomaticPing: false,
		NowFunc: func() time.Time {
			layout := "2006-01-02 15:04:05"
			t1 := time.Now()
			t2, _ := time.Parse(layout, t1.Format(layout))
			now, _ := time.ParseInLocation(layout, t2.Format(layout), t.TimeLocation)
			now, _ = time.Parse(layout, now.Format(layout))
			return now
		},
	})
	if err != nil {
		panic(err)
	}

	// err = db.Use(
	// 	dbresolver.Register(dbresolver.Config{}).
	// 		SetConnMaxIdleTime(time.Hour).
	// 		SetConnMaxLifetime(24 * time.Hour).
	// 		SetMaxIdleConns(100).
	// 		SetMaxOpenConns(200),
	// )
	// if err != nil {
	// 	panic(err)
	// }

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	// Atur parameter connection pool
	sqlDB.SetMaxOpenConns(500)                // Maksimum 300 koneksi terbuka (aktif + idle)
	sqlDB.SetMaxIdleConns(10)                 // Maksimum 10 koneksi idle
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Koneksi hidup maksimal 5 menit
	sqlDB.SetConnMaxIdleTime(1 * time.Minute) // Koneksi idle ditutup setelah 1 menit

	db.Debug()
	fmt.Println("DB '" + t.Name + "' Connected!")
	return db, nil
}

// ConnectionSQL ...
func (t *postgreHelper) ConnectionSQL() (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%v user=%s dbname=%s password=%s sslmode=disable", t.Host, t.Port, t.Username, t.Name, t.Password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sql, err := db.DB()
	if err != nil {
		return nil, err
	}
	return sql, nil
}

// CloseConnection ...
func (t *postgreHelper) CloseConnection(db *gorm.DB) error {
	sqlConn, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlConn.Close()
	return nil
}
