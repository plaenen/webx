package main

import (
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx"
	"github.com/plaenen/webx/cmd/showcase/internal/pages"
	"github.com/plaenen/webx/cmd/showcase/internal/static"
	"github.com/plaenen/webx/session/memory"
	"github.com/plaenen/webx/ui"
	"github.com/plaenen/webx/ui/fileupload"
	"github.com/plaenen/webx/ui/form"
	"github.com/plaenen/webx/ui/markdown"
	"github.com/plaenen/webx/ui/moneyinput"
	"github.com/plaenen/webx/ui/validator"
	"github.com/plaenen/webx/validators"
	"github.com/spf13/cobra"
	"github.com/starfederation/datastar-go/datastar"
)

func main() {
	root := &cobra.Command{
		Use:   "showcase",
		Short: "WebX component showcase",
	}

	root.AddCommand(serveCmd())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func serveCmd() *cobra.Command {
	var (
		port int
		pro  bool
	)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the showcase HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return serve(port, pro)
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 3000, "port to listen on (0 for random)")
	cmd.Flags().BoolVar(&pro, "pro", false, "use Datastar Pro (requires BYOL files in byol/datastar/)")

	return cmd
}

var emailValidator = validator.Handler(func(value string) validator.Result {
	res := validators.Email(value, false)
	return validator.Result{Valid: res.Valid, Error: res.Error}
})

var emailMXValidator = validator.Handler(func(value string) validator.Result {
	res := validators.Email(value, true)
	return validator.Result{Valid: res.Valid, Error: res.Error}
})

var (
	decimalParser         = moneyinput.DecimalHandler()
	moneyParser           = moneyinput.MoneyHandler()
	moneyRestrictedParser = moneyinput.MoneyHandler("USD", "EUR")
	markdownPreview       = markdown.PreviewHandler()
	fileStore             = fileupload.NewStore()
	fileUploader          = fileupload.UploadHandler(fileStore)
	fileRestrictedUpload  = fileupload.UploadHandler(fileStore,
		fileupload.WithAllowedTypes("image/"),
		fileupload.WithMaxFiles(3),
	)
	fileRemover = fileupload.RemoveHandler(fileStore)
)

type loginFormSignals struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var loginFormHandler = form.Handler(
	func(formID string, r *http.Request) []form.FieldError {
		var signals loginFormSignals
		if err := form.ReadSignals(formID, r, &signals); err != nil {
			return []form.FieldError{{Field: "error", Message: "Failed to read form data"}}
		}

		var errs []form.FieldError
		if signals.Email == "" {
			errs = append(errs, form.FieldError{Field: "email_error", Message: "Email is required"})
		} else {
			res := validators.Email(signals.Email, false)
			if !res.Valid {
				errs = append(errs, form.FieldError{Field: "email_error", Message: res.Error})
			}
		}
		if signals.Password == "" {
			errs = append(errs, form.FieldError{Field: "password_error", Message: "Password is required"})
		} else if len(signals.Password) < 8 {
			errs = append(errs, form.FieldError{Field: "password_error", Message: "Password must be at least 8 characters"})
		}
		return errs
	},
	func(formID string, sse *datastar.ServerSentEventGenerator) {
		sanitizedID := strings.ReplaceAll(formID, "-", "_")
		sse.MarshalAndPatchSignals(map[string]any{
			sanitizedID: map[string]any{
				"success": "Login successful!",
			},
		})
	},
)

type contactFormSignals struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

var contactFormHandler = form.Handler(
	func(formID string, r *http.Request) []form.FieldError {
		var signals contactFormSignals
		if err := form.ReadSignals(formID, r, &signals); err != nil {
			return []form.FieldError{{Field: "error", Message: "Failed to read form data"}}
		}

		var errs []form.FieldError
		if signals.Name == "" {
			errs = append(errs, form.FieldError{Field: "name_error", Message: "Name is required"})
		}
		if signals.Email == "" {
			errs = append(errs, form.FieldError{Field: "email_error", Message: "Email is required"})
		} else {
			res := validators.Email(signals.Email, false)
			if !res.Valid {
				errs = append(errs, form.FieldError{Field: "email_error", Message: res.Error})
			}
		}
		if signals.Message == "" {
			errs = append(errs, form.FieldError{Field: "message_error", Message: "Message is required"})
		}
		return errs
	},
	func(formID string, sse *datastar.ServerSentEventGenerator) {
		sanitizedID := strings.ReplaceAll(formID, "-", "_")
		sse.MarshalAndPatchSignals(map[string]any{
			sanitizedID: map[string]any{
				"success": "Message sent successfully!",
			},
		})
	},
)

func serve(port int, pro bool) error {
	readmeBytes, err := os.ReadFile("README.md")
	if err != nil {
		slog.Warn("could not read README.md", "error", err)
	}
	readme := string(readmeBytes)

	r := chi.NewRouter()

	// Session + CSRF middleware
	store := memory.New()
	r.Use(webx.SessionMiddleware(store))
	r.Use(webx.SecurityHeadersMiddleware())

	// Set dev-mode flag and base path on every request
	const basePath = "/showcase"
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wctx := webx.FromContext(r.Context())
			wctx.DevMode = true
			wctx.DatastarPro = pro
			wctx.ShowDatastarInspector = pro
			wctx.BasePath = basePath
			next.ServeHTTP(w, r.WithContext(wctx.WithContext(r.Context())))
		})
	})

	// Serve static files (css, js) at /assets/
	staticFS, _ := fs.Sub(static.Static, "static")
	r.Handle("/assets/*", cacheControl(http.StripPrefix("/assets/", http.FileServerFS(staticFS))))

	// Serve byol files (datastar pro) at /assets/js/datastar/
	if pro {
		byolFS, _ := fs.Sub(static.Byol, "byol/datastar")
		r.Handle("/assets/js/datastar/*", cacheControl(http.StripPrefix("/assets/js/datastar/", http.FileServerFS(byolFS))))
	}

	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// Pages
	r.Get("/", templ.Handler(pages.Home(readme)).ServeHTTP)
	r.Get("/components/button", templ.Handler(pages.Buttons()).ServeHTTP)
	r.Get("/components/card", templ.Handler(pages.Cards()).ServeHTTP)
	r.Get("/components/drawer", templ.Handler(pages.Drawers()).ServeHTTP)
	r.Get("/components/accordion", templ.Handler(pages.Accordions()).ServeHTTP)
	r.Get("/components/alert", templ.Handler(pages.Alerts()).ServeHTTP)
	r.Get("/components/avatar", templ.Handler(pages.Avatars()).ServeHTTP)
	r.Get("/components/calendar", templ.Handler(pages.Calendars()).ServeHTTP)
	r.Get("/components/chat", templ.Handler(pages.Chats()).ServeHTTP)
	r.Get("/components/badge", templ.Handler(pages.Badges()).ServeHTTP)
	r.Get("/components/carousel", templ.Handler(pages.Carousels()).ServeHTTP)
	r.Get("/components/breadcrumbs", templ.Handler(pages.Breadcrumbs()).ServeHTTP)
	r.Get("/components/calendar-advanced", templ.Handler(pages.CalendarAdvanced()).ServeHTTP)
	r.Get("/components/dock", templ.Handler(pages.Docks()).ServeHTTP)
	r.Get("/components/dropdown", templ.Handler(pages.Dropdowns()).ServeHTTP)
	r.Get("/components/fab", templ.Handler(pages.Fabs()).ServeHTTP)
	r.Get("/components/fieldset", templ.Handler(pages.Fieldsets()).ServeHTTP)
	r.Get("/components/footer", templ.Handler(pages.Footers()).ServeHTTP)
	r.Get("/components/file-input", templ.Handler(pages.FileInputs()).ServeHTTP)
	r.Get("/components/filter", templ.Handler(pages.Filters()).ServeHTTP)
	r.Get("/components/label", templ.Handler(pages.Labels()).ServeHTTP)
	r.Get("/components/hover-gallery", templ.Handler(pages.HoverGalleries()).ServeHTTP)
	r.Get("/components/indicator", templ.Handler(pages.Indicators()).ServeHTTP)
	r.Get("/components/join", templ.Handler(pages.Joins()).ServeHTTP)
	r.Get("/components/kbd", templ.Handler(pages.Kbds()).ServeHTTP)
	r.Get("/components/link", templ.Handler(pages.Links()).ServeHTTP)
	r.Get("/components/list", templ.Handler(pages.Lists()).ServeHTTP)
	r.Get("/components/loading", templ.Handler(pages.Loadings()).ServeHTTP)
	r.Get("/components/menu", templ.Handler(pages.Menus()).ServeHTTP)
	r.Get("/components/modal", templ.Handler(pages.Modals()).ServeHTTP)
	r.Get("/components/radio", templ.Handler(pages.Radios()).ServeHTTP)
	r.Get("/components/range", templ.Handler(pages.RangeInputs()).ServeHTTP)
	r.Get("/components/rating", templ.Handler(pages.Ratings()).ServeHTTP)
	r.Get("/components/progress", templ.Handler(pages.Progresses()).ServeHTTP)
	r.Get("/components/radial-progress", templ.Handler(pages.RadialProgresses()).ServeHTTP)
	r.Get("/components/mockup-code", templ.Handler(pages.MockupCodes()).ServeHTTP)
	r.Get("/components/navbar", templ.Handler(pages.Navbars()).ServeHTTP)
	r.Get("/components/pagination", templ.Handler(pages.Paginations()).ServeHTTP)
	r.Get("/components/stat", templ.Handler(pages.Stats()).ServeHTTP)
	r.Get("/components/status", templ.Handler(pages.Statuses()).ServeHTTP)
	r.Get("/components/steps", templ.Handler(pages.Stepss()).ServeHTTP)
	r.Get("/components/select", templ.Handler(pages.SelectInputs()).ServeHTTP)
	r.Get("/components/separator", templ.Handler(pages.Separators()).ServeHTTP)
	r.Get("/components/skeleton", templ.Handler(pages.Skeletons()).ServeHTTP)
	r.Get("/components/tab", templ.Handler(pages.Tabs()).ServeHTTP)
	r.Get("/components/table", templ.Handler(pages.Tables()).ServeHTTP)
	r.Get("/components/textarea", templ.Handler(pages.Textareas()).ServeHTTP)
	r.Get("/components/text-rotate", templ.Handler(pages.TextRotates()).ServeHTTP)
	r.Get("/components/timeline", templ.Handler(pages.Timelines()).ServeHTTP)
	r.Get("/components/toast", templ.Handler(pages.Toasts()).ServeHTTP)
	r.Get("/components/toggle", templ.Handler(pages.Toggles()).ServeHTTP)
	r.Get("/components/tooltip", templ.Handler(pages.Tooltips()).ServeHTTP)
	r.Get("/components/theme-controller", templ.Handler(pages.ThemeControllers()).ServeHTTP)
	r.Get("/components/validator", templ.Handler(pages.Validators()).ServeHTTP)
	r.Get("/components/markdown", templ.Handler(pages.Markdowns()).ServeHTTP)
	r.Get("/components/money", templ.Handler(pages.Moneys()).ServeHTTP)
	r.Get("/components/money-input", templ.Handler(pages.MoneyInputs()).ServeHTTP)
	r.Get("/components/stack", templ.Handler(pages.Stacks()).ServeHTTP)
	r.Get("/components/form", templ.Handler(pages.Forms()).ServeHTTP)
	r.Get("/components/file-upload", templ.Handler(pages.FileUploads()).ServeHTTP)

	// SSE API endpoints
	r.Route(basePath, func(r chi.Router) {
		ui.RegisterRoutes(r)
		r.Get("/api/validate/email", emailValidator)
		r.Get("/api/validate/email-mx", emailMXValidator)
		r.Get("/api/parse/decimal", decimalParser)
		r.Get("/api/parse/money", moneyParser)
		r.Get("/api/parse/money-restricted", moneyRestrictedParser)
		r.Post("/api/preview/markdown", markdownPreview)
		r.Post("/api/form/login", loginFormHandler)
		r.Post("/api/form/contact", contactFormHandler)
		r.Post("/api/upload/files", fileUploader)
		r.Post("/api/upload/files-restricted", fileRestrictedUpload)
		r.Post("/api/upload/remove", fileRemover)
	})

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("listen on port %d: %w", port, err)
	}

	dsMode := "open-source"
	if pro {
		dsMode = "pro"
	}
	slog.Info("server started", "address", fmt.Sprintf("http://localhost:%d", ln.Addr().(*net.TCPAddr).Port), "datastar", dsMode)

	return http.Serve(ln, r)
}

// cacheControl wraps a handler to set Cache-Control: no-cache so browsers
// revalidate via ETags on each request.
func cacheControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		h.ServeHTTP(w, r)
	})
}
