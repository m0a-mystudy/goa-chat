package main

import (
	"fmt"
	goaclient "github.com/goadesign/goa/client"
	"github.com/m0a-mystudy/goa-chat/client"
	"github.com/m0a-mystudy/goa-chat/tool/cli"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"time"
)

func main() {
	// Create command line parser
	app := &cobra.Command{
		Use:   "Chat API-cli",
		Short: `CLI client for the Chat API service`,
	}

	// Create client struct
	httpClient := newHTTPClient()
	c := client.New(goaclient.HTTPClientDoer(httpClient))

	// Register global flags
	app.PersistentFlags().StringVarP(&c.Scheme, "scheme", "s", "", "Set the requests scheme")
	app.PersistentFlags().StringVarP(&c.Host, "host", "H", "localhost:8080", "API hostname")
	app.PersistentFlags().DurationVarP(&httpClient.Timeout, "timeout", "t", time.Duration(20)*time.Second, "Set the request timeout")
	app.PersistentFlags().BoolVar(&c.Dump, "dump", false, "Dump HTTP request and response.")

	// Register signer flags
	var user, pass string
	app.PersistentFlags().StringVar(&user, "user", "", "Username used for authentication")
	app.PersistentFlags().StringVar(&pass, "pass", "", "Password used for authentication")

	// Parse flags and setup signers
	app.ParseFlags(os.Args)
	basicAuthSigner := newBasicAuthSigner(user, pass)

	// Initialize API client
	c.SetBasicAuthSigner(basicAuthSigner)
	c.UserAgent = "Chat API-cli/0"

	// Register API commands
	cli.RegisterCommands(app, c)

	// Execute!
	if err := app.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(-1)
	}
}

// newHTTPClient returns the HTTP client used by the API client to make requests to the service.
func newHTTPClient() *http.Client {
	// TBD: Change as needed (e.g. to use a different transport to control redirection policy or
	// disable cert validation or...)
	return http.DefaultClient
}

// newBasicAuthSigner returns the request signer used for authenticating
// against the basic_auth security scheme.
func newBasicAuthSigner(user, pass string) goaclient.Signer {
	return &goaclient.BasicSigner{
		Username: user,
		Password: pass,
	}

}
