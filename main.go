package main

import (
	"log"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primarykey"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Phone     string    `gorm:"not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type SessionStore struct {
	sessions map[string]uint
}

var (
	db       *gorm.DB
	sessions = &SessionStore{sessions: make(map[string]uint)}
)

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("crystal_rose.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	db.AutoMigrate(&User{})
	seedDefaultUser()

	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/login")
	})

	app.Get("/login", handleLoginPage)
	app.Post("/api/login", handleLogin)
	app.Get("/register", handleRegisterPage)
	app.Post("/api/register", handleRegister)
	app.Post("/api/validate", handleValidation)
	app.Get("/dashboard", authMiddleware, handleDashboard)
	app.Post("/logout", handleLogout)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🌹 Crystal Rose Garden running on http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}

func seedDefaultUser() {
	var count int64
	db.Model(&User{}).Count(&count)
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("rose2024"), bcrypt.DefaultCost)
		defaultUser := User{
			Email:    "rose@crystal.garden",
			Phone:    "+1 (214) 214-2024",
			Password: string(hashedPassword),
		}
		db.Create(&defaultUser)
		log.Println("🌹 Default user created: rose@crystal.garden / rose2024")
	}
}

func generateSessionID() string {
	return time.Now().Format("20060102150405.000000000")
}

func handleLoginPage(c *fiber.Ctx) error {
	sessionID := c.Cookies("session_id")
	if sessionID != "" {
		if _, exists := sessions.sessions[sessionID]; exists {
			return c.Redirect("/dashboard")
		}
	}
	return c.Render("login", fiber.Map{})
}

func handleRegisterPage(c *fiber.Ctx) error {
	sessionID := c.Cookies("session_id")
	if sessionID != "" {
		if _, exists := sessions.sessions[sessionID]; exists {
			return c.Redirect("/dashboard")
		}
	}
	return c.Render("register", fiber.Map{})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type ValidationRequest struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

func handleLogin(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(APIResponse{Success: false, Message: "Invalid request format", Field: "general"})
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Email == "" {
		return c.JSON(APIResponse{Success: false, Message: "Please enter your email", Field: "email"})
	}

	if req.Password == "" {
		return c.JSON(APIResponse{Success: false, Message: "Please enter your password", Field: "password"})
	}

	var user User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.JSON(APIResponse{Success: false, Message: "No garden found with this email", Field: "email"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(APIResponse{Success: false, Message: "Incorrect secret phrase", Field: "password"})
	}

	sessionID := generateSessionID()
	sessions.sessions[sessionID] = user.ID

	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HTTPOnly: true,
		MaxAge:   86400,
		SameSite: "Lax",
	})

	return c.JSON(APIResponse{Success: true, Message: "Welcome to the Crystal Garden 🌹"})
}

func handleRegister(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(APIResponse{Success: false, Message: "Invalid request format", Field: "general"})
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Phone = strings.TrimSpace(req.Phone)

	if err := validateEmail(req.Email); err != "" {
		return c.JSON(APIResponse{Success: false, Message: err, Field: "email"})
	}

	if err := validatePhone(req.Phone); err != "" {
		return c.JSON(APIResponse{Success: false, Message: err, Field: "phone"})
	}

	if err := validatePassword(req.Password); err != "" {
		return c.JSON(APIResponse{Success: false, Message: err, Field: "password"})
	}

	var existing User
	if db.Where("email = ?", req.Email).First(&existing).Error == nil {
		return c.JSON(APIResponse{Success: false, Message: "This garden already exists", Field: "email"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(APIResponse{Success: false, Message: "Failed to plant your rose", Field: "general"})
	}

	user := User{
		Email:    req.Email,
		Phone:    req.Phone,
		Password: string(hashedPassword),
	}

	if err := db.Create(&user).Error; err != nil {
		return c.JSON(APIResponse{Success: false, Message: "Failed to create garden", Field: "general"})
	}

	sessionID := generateSessionID()
	sessions.sessions[sessionID] = user.ID

	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HTTPOnly: true,
		MaxAge:   86400,
		SameSite: "Lax",
	})

	return c.JSON(APIResponse{Success: true, Message: "Your crystal rose has bloomed 🌹"})
}

func handleValidation(c *fiber.Ctx) error {
	var req ValidationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(APIResponse{Success: false, Message: "Invalid request"})
	}

	switch req.Field {
	case "email":
		if err := validateEmail(req.Value); err != "" {
			return c.JSON(APIResponse{Success: false, Message: err, Field: "email"})
		}
		var existing User
		if db.Where("email = ?", strings.ToLower(strings.TrimSpace(req.Value))).First(&existing).Error == nil {
			return c.JSON(APIResponse{Success: false, Message: "This garden already exists", Field: "email"})
		}
		return c.JSON(APIResponse{Success: true, Message: "Garden available ✓", Field: "email"})

	case "phone":
		if err := validatePhone(req.Value); err != "" {
			return c.JSON(APIResponse{Success: false, Message: err, Field: "phone"})
		}
		return c.JSON(APIResponse{Success: true, Message: "Valid ✓", Field: "phone"})

	case "password":
		if err := validatePassword(req.Value); err != "" {
			return c.JSON(APIResponse{Success: false, Message: err, Field: "password"})
		}
		return c.JSON(APIResponse{Success: true, Message: "Strong phrase ✓", Field: "password"})

	case "login_email":
		email := strings.ToLower(strings.TrimSpace(req.Value))
		if email == "" {
			return c.JSON(APIResponse{Success: false, Message: "Please enter email", Field: "email"})
		}
		var user User
		if db.Where("email = ?", email).First(&user).Error != nil {
			return c.JSON(APIResponse{Success: false, Message: "Garden not found", Field: "email"})
		}
		return c.JSON(APIResponse{Success: true, Message: "Garden found ✓", Field: "email"})
	}

	return c.JSON(APIResponse{Success: true})
}

func validateEmail(email string) string {
	email = strings.TrimSpace(email)
	if email == "" {
		return "Email is required"
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return "Please enter a valid email"
	}
	return ""
}

func validatePhone(phone string) string {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return "Phone is required"
	}
	digits := regexp.MustCompile(`\d`).FindAllString(phone, -1)
	if len(digits) < 10 {
		return "At least 10 digits needed"
	}
	if len(digits) > 15 {
		return "Too many digits"
	}
	return ""
}

func validatePassword(password string) string {
	if password == "" {
		return "Password is required"
	}
	if len(password) < 6 {
		return "At least 6 characters"
	}
	if len(password) > 72 {
		return "Too long"
	}

	hasLetter := false
	hasNumber := false
	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if unicode.IsDigit(char) {
			hasNumber = true
		}
	}

	if !hasLetter {
		return "Include a letter"
	}
	if !hasNumber {
		return "Include a number"
	}

	return ""
}

func authMiddleware(c *fiber.Ctx) error {
	sessionID := c.Cookies("session_id")
	if sessionID == "" {
		return c.Redirect("/login")
	}

	userID, exists := sessions.sessions[sessionID]
	if !exists {
		return c.Redirect("/login")
	}

	var user User
	if err := db.First(&user, userID).Error; err != nil {
		delete(sessions.sessions, sessionID)
		return c.Redirect("/login")
	}

	c.Locals("user", user)
	return c.Next()
}

func handleDashboard(c *fiber.Ctx) error {
	user := c.Locals("user").(User)
	return c.Render("dashboard", fiber.Map{
		"Email": user.Email,
		"Phone": user.Phone,
	})
}

func handleLogout(c *fiber.Ctx) error {
	sessionID := c.Cookies("session_id")
	if sessionID != "" {
		delete(sessions.sessions, sessionID)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    "",
		HTTPOnly: true,
		MaxAge:   -1,
		SameSite: "Lax",
	})

	return c.Redirect("/login")
}
