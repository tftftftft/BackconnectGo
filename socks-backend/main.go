package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

var db *sql.DB
var store = sessions.NewCookieStore([]byte("secret-key"))

func main() {
	ConnectDB()
	setupRouterAndStartServer()
}

func setupRouterAndStartServer() {
	router := gin.Default()
	router.Use(CORSMiddleware())

	// Apply the middleware to routes that require authentication
	authorized := router.Group("/")
	authorized.Use(AuthRequired())
	{
		authorized.GET("/fetchProxies", handleFetchProxies)
	}

	router.POST("/api/login", handleLogin)

	router.Run(":30000")
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := store.Get(c.Request, "user-session")
		if err != nil {
			log.Println("Error getting session:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			c.Abort()
			return
		}

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func handleLogin(c *gin.Context) {
	var credentials Credentials

	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	isValid, err := areCredentialsValid(credentials)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	if isValid {
		session, _ := store.Get(c.Request, "user-session")
		session.Values["authenticated"] = true
		session.Values["username"] = credentials.Username
		session.Save(c.Request, c.Writer)
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "invalid"})
	}
}

func areCredentialsValid(creds Credentials) (bool, error) {
	query := `SELECT COUNT(*) FROM Users WHERE Username = ? AND Password = ?`
	var count int
	err := db.QueryRow(query, creds.Username, creds.Password).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func FetchDataByUserID(db *sql.DB, userID string) ([]ProxyInfo, error) {
	query := "SELECT ServerIP, ServerListeningPort, ProxyIP, CountryCode, Region, City, Zip, Mobile, Proxy, Hosting FROM Proxies WHERE UserID = ?"
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var proxies []ProxyInfo
	for rows.Next() {
		var proxy ProxyInfo
		if err := rows.Scan(&proxy.ServerIP, &proxy.ServerListeningPort, &proxy.ProxyIP, &proxy.CountryCode, &proxy.Region, &proxy.City, &proxy.Zip, &proxy.Mobile, &proxy.Proxy, &proxy.Hosting); err != nil {
			return nil, err
		}
		proxies = append(proxies, proxy)
	}
	return proxies, nil
}

func getUserIDByUsername(username string) (string, error) {
	query := `SELECT UserID FROM Users WHERE Username = ?`
	var userID string
	err := db.QueryRow(query, username).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func handleFetchProxies(c *gin.Context) {
	log.Println("handleFetchProxies called") // Logging when the function is called

	// Fetch session data
	session, err := store.Get(c.Request, "user-session")
	if err != nil {
		log.Println("Error getting session:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}
	log.Println("Session data fetched") // Logging after fetching session data

	// Fetch username from session
	username, ok := session.Values["username"].(string)
	if !ok {
		log.Println("Username not found in session")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}
	log.Printf("Username from session: %s\n", username) // Logging the username

	// Fetch user ID
	userID, err := getUserIDByUsername(username)
	if err != nil {
		log.Printf("Error fetching user ID for username %s: %s", username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}
	log.Printf("User ID for username %s: %s\n", username, userID) // Logging the user ID

	// Fetch proxies
	proxies, err := FetchDataByUserID(db, userID)
	if err != nil {
		log.Printf("Error fetching proxies for user ID %s: %s", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	log.Printf("Fetched %d proxies for user ID %s and username %s\n", len(proxies), userID, username)
	log.Println(proxies) // Logging the proxies

	c.JSON(http.StatusOK, gin.H{"data": proxies})
}
