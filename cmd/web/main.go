package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"forum/internal/database"
	deliveryHTTP "forum/internal/delivery/http"
	"forum/internal/repository"
	"forum/internal/service"
)

func main() {
	// 1. RESOLVE CONFIGURATION PARAMETERS (With strict fail-safes)
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./forum.db" // Fallback local project runtime path
	}

	schemaPath := os.Getenv("SCHEMA_PATH")
	if schemaPath == "" {
		schemaPath = "./internal/database/schema.sql" // Target schema definitions file
	}

	log.Printf("[INIT] Initializing ViableForum Production Application Node...")
	log.Printf("[INIT] Database Storage Targets: %s", dbPath)
	log.Printf("[INIT] Active Migration Schema Path: %s", schemaPath)

	// 2. INITIALIZE COHESIVE EMBEDDED DATABASES LAYER
	db, err := database.InitDB(dbPath, schemaPath)
	if err != nil {
		log.Fatalf("[CRITICAL] Storage engine initialization error: %v", err)
	}
	defer func() {
		log.Println("[SHUTDOWN] Severing storage connection pools gracefully...")
		_ = db.Close()
	}()

	// 3. EXECUTE REPOSITORY SEEDING ENGINES (Fulfills initial mock conditions natively)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := database.SeedCategories(ctx, db); err != nil {
		log.Fatalf("[CRITICAL] Failed applying fallback configuration metadata: %v", err)
	}
	if err := database.SeedMockData(ctx, db); err != nil {
		log.Printf("[WARN] Mock data sequence skipped or already executed: %v", err)
	}
	log.Println("[INIT] Database Schema migrations and foundational seeds compiled successfully.")

	// 4. DEPENDENCY INJECTION LAYER (Top-Down Wiring)
	// Pure SoC: Data references bubble up sequentially into logic matrices
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	interactionRepo := repository.NewInteractionRepository(db)

	authService := service.NewAuthService(userRepo, sessionRepo)
	postService := service.NewPostService(postRepo)
	commentService := service.NewCommentService(commentRepo)
	filterService := service.NewFilterService(postRepo, interactionRepo)
	interactionService := service.NewInteractionService(interactionRepo)

	// Bundle decoupled references into delivery handlers memory maps
	container := &deliveryHTTP.HandlerContainer{
		AuthService:        authService,
		PostService:        postService,
		CommentService:     commentService,
		InteractionService: interactionService,
		FilterService:      filterService,
	}

	// 5. ASSEMBLE DELIVERY HTTP ROUTER ROUTING TABLES
	router := deliveryHTTP.NewRouter(container)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 6. HEALTHY GRACEFUL INTERRUPT OS SHUTDOWN LIFECYCLE MANAGEMENT
	// This prevents SQLite data file corruptions by draining traffic when containers spin down!
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("[RUNTIME] Architectural engine active. Web platform running on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[CRITICAL] Unhandled transport runtime crashing fault: %v", err)
		}
	}()

	// Wait explicitly for container orchestration terminate commands (e.g., Docker stop)
	<-shutdownChan
	log.Println("[SHUTDOWN] Intercepted shutdown vector sequence. Initiating systems cleanup...")

	graceCtx, graceCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer graceCancel()

	if err := server.Shutdown(graceCtx); err != nil {
		log.Fatalf("[CRITICAL] Server forced into dirty termination layout sequence: %v", err)
	}

	log.Println("[SHUTDOWN] Forum engine node offline safely. Architecture preserved.")
}
