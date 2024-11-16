package cmd

import (
	"fmt"
	"net"
	"os"
	"strings"

	"remy/lib"
	"remy/server"

	"github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
)

type Parameters struct {
	ServerHost string
	ServerPort uint16
}

const defaultServerHost = "0.0.0.0"
const defaultServerPort = 8888

var rootCmd = &cobra.Command{
	Use:   "remy",
	Short: "Make your phone a remote mouse",
	Run: func(cmd *cobra.Command, args []string) {
		parameters, err := getParameters(cmd)
		if err != nil {
			fmt.Println(err)
			return
		}

		printCmdHeader(parameters)
		if err != nil {
			fmt.Println(err)
			return
		}

		runServer(parameters)
	},
}

func getParameters(cmd *cobra.Command) (*Parameters, error) {
	serverHost, err := cmd.Flags().GetString("address")
	if err != nil {
		return nil, err
	}

	serverPort, err := cmd.Flags().GetUint16("port")
	if err != nil {
		return nil, err
	}

	parameters := &Parameters{
		ServerHost: serverHost,
		ServerPort: serverPort,
	}

	return parameters, nil
}

func printCmdHeader(parameters *Parameters) error {
	displayedServerHost := parameters.ServerHost
	if displayedServerHost == defaultServerHost {
		primaryIP, err := getMyPrimaryIP()
		if err != nil {
			return err
		}
		displayedServerHost = primaryIP
	}

	displayedServerAddr := fmt.Sprintf("https://%s:%d/", displayedServerHost, parameters.ServerPort)
	qrcode, err := qrcode.New(displayedServerAddr, qrcode.Low)
	if err != nil {
		return err
	}

	lib.PrintQR(qrcode)
	fmt.Printf("Available at: %s\n", displayedServerAddr)
	fmt.Println("(available on all other interfaces too)")

	return nil
}

func runServer(parameters *Parameters) {
	serverAddr := fmt.Sprintf("%s:%d", parameters.ServerHost, parameters.ServerPort)

	exitSignal := make(chan interface{})
	go func() {
		err := server.Run(serverAddr)
		if err != nil {
			fmt.Println("Error starting server:", err)
		}
		close(exitSignal)
	}()

	<-exitSignal
}

func getMyPrimaryIP() (string, error) {
	const googleDNSAddr string = "8.8.8.8:53"
	googleAddr, err := net.ResolveTCPAddr("tcp", googleDNSAddr)
	if err != nil {
		return "", err
	}

	conn, err := net.DialTCP("tcp", nil, googleAddr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localIP := strings.Split(conn.LocalAddr().String(), ":")[0]
	return localIP, nil
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("address", "a", defaultServerHost, "bind the HTTP server to this address")
	rootCmd.Flags().Uint16P("port", "p", defaultServerPort, "bind the HTTP server to this port")
}
