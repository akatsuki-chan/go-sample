package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"moul.io/zapgorm"
)

const VERSION = "1.0.0"
const location = "Asia/Tokyo"

func init() {
	// $HOME/config.ymlか同じディレクトリにあるconfig.ymlを設定ファイルとして読み込む
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	viper.SetDefault("debug", true)
	viper.SetDefault("env", "prod")
}

func main() {
	var version bool
	flag.BoolVar(&version, "v", false, "print version.")
	flag.Parse()

	// バージョン表示（バイナリになるため）
	if version {
		fmt.Printf("version: %s\n", VERSION)
		os.Exit(0)
	}

	// 環境変数読み込み。ファイルがなければ、エラー無視。
	_ = godotenv.Load(fmt.Sprintf(".%s.env", viper.GetString("env")))

	logger, _ := zap.NewDevelopment()
	if !viper.GetBool("debug") {
		logger, _ = zap.NewProduction()
	}
	defer logger.Sync()
	defer zap.S().Sync()

	// グローバルなロガーの設定
	logger = logger.
		With(
			zap.String("version", VERSION),
			zap.String("env", viper.GetString("ENV")),
		).
		Named("app")
	zap.ReplaceGlobals(logger)

	// config.ymlのdbの価でデータベース接続。
	db, err := gorm.Open("mysql", viper.Get("db"))
	if err != nil {
		zap.S().Fatal(err)
		panic(err)
	}
	defer db.Close()

	// gormのSQLのログをzapに出力
	db.LogMode(viper.GetBool("debug"))
	db.SetLogger(zapgorm.New(zap.L().Named("gorm")))

	rows, err := db.Table("sample").
		Where("id = ?", 1).
		Rows()
	if err != nil {
		zap.S().Error(err)
	}
	defer rows.Close()

	zap.S().Info("Hoge")
}
