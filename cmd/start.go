package cmd

import (
	"fmt"
	"github.com/EscanBE/escan-request-redirector/constants"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	flagRedirectMessage = "redirect-message"
	flagPort            = "port"
	flagEngine          = "engine"
	flagRedirectTimeout = "redirect-timeout"
)

// startCmd represents the start command, it launches the main business logic of this app
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start request redirect service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetBlockExplorer := strings.TrimSuffix(args[0], "/")
		port, _ := cmd.Flags().GetUint16(flagPort)
		redirectTimeout, _ := cmd.Flags().GetDuration(flagRedirectTimeout)

		engine, _ := cmd.Flags().GetString(flagEngine)
		if engine == "" {
			if strings.Contains(targetBlockExplorer, "mintscan.io") {
				engine = constants.EngineMintscan
			} else if strings.Contains(targetBlockExplorer, "silknodes.io") {
				engine = constants.EngineSilkNodes
			} else {
				fmt.Printf("Failed to detect block explorer engine, please specify it with --%s\n", flagEngine)
				os.Exit(1)
			}
		}

		redirectMessage, _ := cmd.Flags().GetString(flagRedirectMessage)
		if redirectMessage == "" {
			redirectMessage = fmt.Sprintf("Temporary unavailable, redirecting to %s", targetBlockExplorer)
		}

		http.HandleFunc("/tx/", func(w http.ResponseWriter, r *http.Request) {
			hash := strings.TrimPrefix(r.URL.Path, "/tx/")
			if !isSafeUrlPath(hash) {
				http.Error(w, "Bad URL", http.StatusBadRequest)
				return
			}
			w.Write(buildRedirectContent(targetBlockExplorer+"/tx/"+hash, redirectMessage, redirectTimeout))
		})
		handleFuncAddress := func(addr string, w http.ResponseWriter, r *http.Request) {
			if !isSafeUrlPath(addr) {
				http.Error(w, "Bad URL", http.StatusBadRequest)
				return
			}
			var url string
			switch engine {
			case constants.EngineMintscan:
				url = targetBlockExplorer + "/address-estimation/" + addr
			case constants.EngineSilkNodes:
				url = targetBlockExplorer + "/account/" + addr
			default:
				url = targetBlockExplorer + "/address/" + addr
			}
			w.Write(buildRedirectContent(url, redirectMessage, redirectTimeout))
		}
		http.HandleFunc("/address/", func(w http.ResponseWriter, r *http.Request) {
			addr := strings.TrimPrefix(r.URL.Path, "/address/")
			handleFuncAddress(addr, w, r)
		})
		http.HandleFunc("/token/", func(w http.ResponseWriter, r *http.Request) {
			addr := strings.TrimPrefix(r.URL.Path, "/token/")
			handleFuncAddress(addr, w, r)
		})
		handleFuncDirectForward := func(w http.ResponseWriter, r *http.Request) {
			w.Write(buildRedirectContent(targetBlockExplorer+"/"+r.URL.Path, redirectMessage, redirectTimeout))
		}
		http.HandleFunc("/validators", handleFuncDirectForward)
		handleFuncReplaceSimplePath := func(path string, w http.ResponseWriter, r *http.Request) {
			w.Write(buildRedirectContent(targetBlockExplorer+"/"+strings.TrimPrefix(path, "/"), redirectMessage, redirectTimeout))
		}
		http.HandleFunc("/blocks", func(w http.ResponseWriter, r *http.Request) {
			var path string
			switch engine {
			case constants.EngineMintscan:
				path = "/block"
			case constants.EngineSilkNodes:
				path = "/block"
			default:
				path = "/blocks"
			}
			handleFuncReplaceSimplePath(path, w, r)
		})
		http.HandleFunc("/txs", func(w http.ResponseWriter, r *http.Request) {
			var path string
			switch engine {
			case constants.EngineMintscan:
				path = "/tx"
			case constants.EngineSilkNodes:
				path = "/block"
			default:
				path = "/txs"
			}
			handleFuncReplaceSimplePath(path, w, r)
		})
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, targetBlockExplorer, http.StatusTemporaryRedirect)
		})

		fmt.Printf("Starting HTTP server at :%d...\n", port)
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	},
}

func init() {
	startCmd.Flags().StringP(flagRedirectMessage, "m", "", "Message to be displayed when redirecting")
	startCmd.Flags().Uint16P(flagPort, "p", 8080, "Port to listen on")
	startCmd.Flags().String(flagEngine, "", fmt.Sprintf("Block explorer engine to redirect to, available options: %s, %s", constants.EngineMintscan, constants.EngineSilkNodes))
	startCmd.Flags().Duration(flagRedirectTimeout, 3*time.Second, "Redirect timeout in seconds")
	rootCmd.AddCommand(startCmd)
}

func isSafeUrlPath(str string) bool {
	for _, c := range str {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '/' || c == '%') {
			return false
		}
	}
	return true
}

func buildRedirectContent(url, messageContent string, timeout time.Duration) []byte {
	if timeout < 1*time.Second {
		timeout = 1 * time.Second
	}
	return []byte(fmt.Sprintf(`
<html>
	<head>
		<title>Escan Redirection</title>
		<script type="text/javascript">
			setTimeout(function() {
				window.location.href = "%s";
			}, %d);
		</script>
	</head>
	<body style="background-color: black; color: white;">
		<h3>%s</h3>
		<p><i>Redirecting in few seconds...</i></p>
	</body>
</html>`, url, timeout.Milliseconds(), messageContent))
}
