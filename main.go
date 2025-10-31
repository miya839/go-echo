package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/labstack/gommon/log"
	"net/http"
	"time"
)

// UserSpecificMiddleware は /users/ 配下でのみ実行されるダミーミドルウェア
func serverHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		// ヘッダ設定処理
		c.Response().Before(func() {
			c.Response().Header().Set(echo.HeaderServer, "Go/Echo Custom Server")
			
		})
		c.Logger().Infof("Request processed in %v", time.Since(start))

		err := next(c)
		return err
	}
}

// リクエストを処理するハンドラ関数
func handler(c echo.Context) error {

	// 1. クエリパラメータの取得
	name := c.QueryParam("name")

	// 取得したデータをJSON形式でレスポンスとして返す
	if name == "" {
		return c.String(http.StatusOK, "Hello, echo api server!")
	} else {
		return c.String(http.StatusOK, fmt.Sprintf("Hello, %s", name))
	}
}

func helloPathHandler(c echo.Context) error {
	// 1. パスパラメータの取得
	name := c.Param("name")

	// 抽出したIDが空でないか、または追加でバリデーションを行う
	if name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name is missing")
	}

	// 取得したデータをJSON形式でレスポンスとして返す
	return c.JSON(http.StatusOK, echo.Map{
		"message": fmt.Sprintf("Hello, %s", name),
	})
}

// User構造体: クライアントから受け取るJSONデータに対応するGoの構造体
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// userCreateHandler は新しいユーザーを作成するためのPOSTリクエストを処理します。
func userCreateHandler(c echo.Context) error {
	// 1. 構造体の定義
	user := new(User)

	// 2．リクエストボディのJSON/XML/Formデータを構造体 'user' にバインド
	if err := c.Bind(user); err != nil {
		// バインドエラー（例: 不正なJSON形式）が発生した場合
		return echo.NewHTTPError(http.StatusBadRequest, "リクエストボディの形式が不正です。")
	}

	// ここでデータベースにデータを保存するなどの実際の処理を行います
	log.Printf("Received new user: Name=%s, Email=%s", user.Name, user.Email)

	return c.JSON(http.StatusCreated, user)
}

// userModifyHandler はユーザーを更新するためのPUTリクエストを処理します。
func userModifyHandler(c echo.Context) error {
	// 1. 構造体の定義
	user := new(User)

	// 2．リクエストボディのJSON/XML/Formデータを構造体 'user' にバインド
	if err := c.Bind(user); err != nil {
		// バインドエラー（例: 不正なJSON形式）が発生した場合
		return echo.NewHTTPError(http.StatusBadRequest, "リクエストボディの形式が不正です。")
	}

	// ここでデータベースにデータを保存するなどの実際の処理を行います
	log.Printf("Received new user: Name=%s, Email=%s", user.Name, user.Email)

	return c.JSON(http.StatusNoContent, user)
}

func main() {
	// Echoインスタンスを作成
	e := echo.New()

	// ロガーミドルウェア: すべてのリクエスト/レスポンス情報をログに出力します。
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.DEBUG)

	// リカバリーミドルウェア: パニックが発生した場合にアプリケーションをクラッシュから保護します。
	e.Use(middleware.Recover())
	e.Use(serverHeader)

	// HTTP GETリクエストに対するハンドラを定義
	helloGroup := e.Group("/hello")
	helloGroup.GET("/", handler)
	helloGroup.GET("/:name", helloPathHandler)

	userGroup := e.Group("/users")
	userGroup.Use(serverHeader)
	userGroup.POST("/", userCreateHandler)
	userGroup.PUT("/", userModifyHandler)
	
	// サーバーをポート8080で起動
	e.Logger.Fatal(e.Start(":8080"))
}
