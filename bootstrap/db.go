package bootstrap

import (
	"io"
	"log"
	"os"
	"strconv"
	"time"
	"webapp_gin/app/models"
	"webapp_gin/global"

	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitializeDB() *gorm.DB {

	switch global.App.Config.Database.Driver {
	case "mysql":
		return initMySqlGorm()
	default:
		return initMySqlGorm()
	}
}

func initMySqlGorm() *gorm.DB {
	dbConfig := global.App.Config.Database

	if dbConfig.Database == "" {
		return nil
	}
	dsn := dbConfig.UserName + ":" + dbConfig.Password + "@tcp(" + dbConfig.Host + ":" + strconv.Itoa(dbConfig.Port) + ")/" +
		dbConfig.Database + "?charset=" + dbConfig.Charset + "&parseTime=True&loc=Local"
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string default length
		DisableDatetimePrecision:  true,  // disable datetime precision，version before MySQL 5.6 unsupported
		DontSupportRenameIndex:    true,  // disable rename index，version before MySQL 5.7 and mariadb unsupported
		DontSupportRenameColumn:   true,  // disable rename column，version before MySQL 8 and MariaDB unsupported
		SkipInitializeWithVersion: false, // disable init with version
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,            // disable foreign key constraint when migrating
		Logger:                                   getGormLogger(), // use customized Logger
	}); err != nil {
		zap.L().Error("mysql connect failed, err:", zap.Any("err", err))
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
		sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
		initMySqlTables(db)
		return db
	}
}

func getGormLogWriter() logger.Writer {
	var writer io.Writer

	if global.App.Config.Database.EnableFileLogWriter {

		writer = &lumberjack.Logger{
			Filename:   global.App.Config.Log.RootDir + "/" + global.App.Config.Database.LogFilename,
			MaxSize:    global.App.Config.Log.MaxSize,
			MaxBackups: global.App.Config.Log.MaxBackups,
			MaxAge:     global.App.Config.Log.MaxAge,
			Compress:   global.App.Config.Log.Compress,
		}
	} else {
		//default writter
		writer = os.Stdout
	}
	return log.New(writer, "\r\n", log.LstdFlags)
}

func getGormLogger() logger.Interface {
	var logMode logger.LogLevel

	switch global.App.Config.Database.LogMode {
	case "silent":
		logMode = logger.Silent
	case "error":
		logMode = logger.Error
	case "warn":
		logMode = logger.Warn
	case "info":
		logMode = logger.Info
	default:
		logMode = logger.Info
	}

	return logger.New(getGormLogWriter(), logger.Config{
		SlowThreshold:             200 * time.Millisecond,                          // slow query threshold
		LogLevel:                  logMode,                                         // log level
		IgnoreRecordNotFoundError: false,                                           // ignore record not found error
		Colorful:                  !global.App.Config.Database.EnableFileLogWriter, // disable colorful printing
	})
}

func initMySqlTables(db *gorm.DB) {
	err := db.AutoMigrate(
		//models.WechatUser{},
		models.SquareSetting{},
		models.RecommendSetting{},
		//models.Pet{},
		models.KnowMe{},
		models.FocusOn{},
		models.View{},
	)
	if err != nil {
		global.App.Log.Error("migrate table failed", zap.Any("err", err))
		os.Exit(0)
	}
}
