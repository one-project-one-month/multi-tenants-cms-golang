package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"image/png"
	"os"
	"strings"
)

/*
	AUTHOR:
		Swan Htet Aung Phyo
*/
/*
	PURPOSE:
		How does the mfa work under the hood. For the system, I suggest to use the MFA as another security layer
*/
/*
NOTE:
	TOTP is generated based on the time-window (like 30 second). Not the exact time.
	Server generates the secrete key which is typically base32 length string.
	TOTP function is mathematical form
		TOTP = Truncate(HMAC_SHA1(server secret_key, current_timestamp)
		-- timestamp = floor(unix_time/30) - notice: 30 is time window
	HAMC is the hashing function which will produce the hash of the time window get from totp function
	Calculation = time_stamp = floor(171972005 / 30 ) = 5732400
		HMAC(server secret,5732400) and truncate last digit after 6 digit
	The client like google authenticator will do the same action
	In real app, we have to store server secret with AES 128 encryption and
		Rotate the key regularly

*/

func display(key *otp.Key, data []byte) {
	color.Green("Issuer: %s", key.Issuer())
	color.Yellow("Account Name: %s", key.AccountName())
	color.Green("Secret: %s", key.Secret())
	color.Blue("Writing png file ")
	go func() {
		err := os.WriteFile("qr-code.png", data, 0600)
		if err != nil {
			panic(err.Error())
			return
		}
	}()

	color.Blue("PNG File written")
	color.Yellow("Please add totp to OTP authentication (google authenticator)")
}

func promptCode() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter code: ")
	code, _ := reader.ReadString('\n')
	code = strings.Replace(code, "\n", "", -1)
	return code
}

func main() {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Content Management system",
		AccountName: "content@gmail.com",
	})

	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		panic(err.Error())
	}
	err = png.Encode(&buf, img)
	if err != nil {
		color.Red(err.Error())
		return
	}
	display(key, buf.Bytes())
	code := promptCode()
	if totp.Validate(code, key.Secret()) {
		color.Green("Success")
		removeFile("qr-code.png")
		os.Exit(0)
	}
	color.Red("Failure")
	removeFile("qr-code.png")
}

func removeFile(path string) {
	err := os.Remove(path)
	if err != nil {
		color.Red(err.Error())
	}
	color.Green("Success removal")
}
