package controller

import (
	"html/template"
	"net/http"
	"strings"

	"forum/domain/entity"
	"forum/usecase"
)

type AuthController struct {
	authService *usecase.AuthService
	postService *usecase.PostService
	templates   *template.Template
}

func NewAuthController(authService *usecase.AuthService, postService *usecase.PostService, templates *template.Template) *AuthController {
	return &AuthController{
		authService: authService,
		postService: postService,
		templates:   templates,
	}
}

func (c *AuthController) HandleSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.RawQuery != "" {
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		c.renderTemplate(w, "register.html", nil)
		return
	}

	if r.Method != http.MethodPost {
		c.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      "Method Not Allowed",
		})
		return
	}
	// if the method is POST that means the user is creating an account
	name := r.PostFormValue("username")
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	_, err := c.authService.Signup(name, email, password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.renderTemplate(w, "register.html", map[string]interface{}{
			"registerError": err.Error(),
			"username":      name,
			"email":         email,
		})
		return
	}

	// Redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (c *AuthController) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// If there are query parameters, redirect to clean URL
		if r.URL.RawQuery != "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		c.renderTemplate(w, "login.html", nil)
		return
	}

	if r.Method != http.MethodPost {
		c.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      "Method Not Allowed",
		})
		return
	}

	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	token, user, err := c.authService.Login(email, password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		c.renderTemplate(w, "login.html", map[string]interface{}{
			"loginError": err.Error(),
			"email":      email,
		})
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
	})

	_ = user

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *AuthController) HandleMainPage(w http.ResponseWriter, r *http.Request) {
	c.ShowMainPage(w, r)
}

func (c *AuthController) HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" && r.Method == http.MethodGet {
		c.ShowMainPage(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/static/") {
		switch r.URL.Path {
		case "/static/css/layout.css", "/static/css/login.css", "/static/css/posts.css", "/static/css/register.css", "/static/css/error.css":
			http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP(w, r)

		case "/static/", "/static/css/", "/static/images/":
			c.ShowErrorPage(w, ErrorMessage{
				StatusCode: http.StatusForbidden,
				Error:      "StatusForbidden",
			})

		default:
			c.ShowErrorPage(w, ErrorMessage{
				StatusCode: http.StatusNotFound,
				Error:      "Page Not Found",
			})
		}
	} else {
		c.ShowErrorPage(w, ErrorMessage{
			StatusCode: http.StatusNotFound,
			Error:      "Page Not Found",
		})
	}
}

func (c *AuthController) HandleLogout(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*entity.User)
	if ok {
		c.authService.Logout(user.ID)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
