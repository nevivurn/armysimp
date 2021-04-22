package armysimp

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	"cloud.google.com/go/translate"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"github.com/nevivurn/armysimp/report"
	"github.com/nevivurn/armysimp/scrape"
	"github.com/nevivurn/armysimp/thecamp"
)

// Variables populated from the environment
var (
	ytAPIKey      = os.Getenv("YOUTUBE_API_KEY")
	storageBucket = os.Getenv("STORAGE_BUCKET")

	thecampUser = os.Getenv("THECAMP_USER")
	thecampPass = os.Getenv("THECAMP_PASS")
)

// Global shared state.
var global struct {
	st *storage.Client
	tl *translate.Client
	yt *youtube.Service

	cl *scrape.Client
}

// Initialize global state.
func init() {
	ctx := context.Background()
	var err error

	if ytAPIKey == "" || storageBucket == "" {
		panic("required envvars not set")
	}

	global.st, err = storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	global.tl, err = translate.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	global.yt, err = youtube.NewService(ctx, option.WithAPIKey(ytAPIKey))
	if err != nil {
		panic(err)
	}

	global.cl = scrape.New(global.tl, global.yt)
}

// Handler is the function entrypoint.
func Handler(w http.ResponseWriter, r *http.Request) {
	if err := run(r.Context()); err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "OK")
}

func run(ctx context.Context) error {
	data, err := global.cl.Scrape(ctx)
	if err != nil {
		return err
	}

	fullReport, err := report.Generate(data)
	if err != nil {
		return err
	}
	favReport, err := report.GenerateFavorites(data)
	if err != nil {
		return err
	}

	var prefix string
	if err := send(ctx, fullReport); err != nil {
		prefix = fmt.Sprintf("AUTO SEND FAILED: %v\n", err)
	} else {
		prefix = "AUTO SEND SUCCEDED\n"
	}
	fullReport = prefix + fullReport
	favReport = prefix + favReport

	bucket := global.st.Bucket(storageBucket)

	fullWriter := bucket.Object("report.txt").NewWriter(ctx)
	fmt.Fprint(fullWriter, fullReport)
	fullErr := fullWriter.Close()

	favWriter := bucket.Object("report-short.txt").NewWriter(ctx)
	fmt.Fprint(favWriter, favReport)
	favErr := favWriter.Close()

	if fullErr != nil {
		return fullErr
	}
	if favErr != nil {
		return favErr
	}

	return nil
}

func send(ctx context.Context, report string) error {
	camp, err := thecamp.New(ctx, thecampUser, thecampPass)
	if err != nil {
		return err
	}

	code, err := camp.Search(ctx, "성용운", "19991105", "20210429")
	if err != nil {
		return err
	}

	camp.Send(ctx, code, "armysimp", report)

	return nil
}
