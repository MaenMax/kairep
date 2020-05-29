package main

import (
	"bufio"
	"encoding/hex"
	"errors"
	fernet "fernet-go"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gocql/gocql"
	"golang.org/x/exp/utf8string"
)

var (
	subscription *string = flag.String("endpoint", "", "Endpoint to be removed")

	crypto_key *string = flag.String("crypto_key", "", "Fernet crypto key")

	cass_addr *string = flag.String("cass_addr", "", "Cassandra address")

	cass_user *string = flag.String("cass_user", "", "Cassandra username")

	cass_password *string = flag.String("cass_pass", "", "Cassandra password")

	ask_user *bool = flag.Bool("y", false, "Ask user for delete confirmation of endpoints one by one.")

	err error

	cass_session *gocql.Session
	uaid         string
	chid         string
	pub_key      string
)

const (
	TTL_DURATION time.Duration = 1000 * 1000 * 1000 * 60 * 60 * 24 * 365 * 100
)

func main() {

	flag.Parse()
	fmt.Println("Crypto key is: ", *crypto_key)
	err = start_cassandra()
	if err != nil {
		str := fmt.Sprintf("Cannot start cassandra: '%v'", err)
		panic(str)
	}

	subscriptions := readFile(*subscription)
	fmt.Println("Deleting the list of provided subscriptions...")
	crypto_key_utf8_encoded := utf8string.NewString(*crypto_key)
	key := fernet.MustDecodeKeys(crypto_key_utf8_encoded.String())

	for _, sub := range subscriptions {
		if len(sub) == 183 {

			uaid, chid, pub_key, err = decryptEndpoint(Repad(sub), 2, key)
			if err != nil {
				fmt.Println("Failed to decrypt VAPID subscription data: '%s'.", err)
				return
			}

		} else {

			uaid, chid, pub_key, err = decryptEndpoint(Repad(sub), 1, key)
			if err != nil {
				fmt.Println("Failed to decrypt V1 subscription data: ", err)
				return
			}

		}

		//var text string
		if !*ask_user {
			reader := bufio.NewReader(os.Stdin)
			fmt.Println(fmt.Sprintf("Do you want to delete subscription: %s ?", uaid))
			_, _ = reader.ReadString('\n')
		}

		dropped := DropUser(uaid, cass_session)
		if dropped == false {

			fmt.Println("Error: Cannot drop user: %s", uaid)

			return

		}
		fmt.Println(fmt.Sprintf("uaid: %s, is deleted.", uaid))

	}

	fmt.Println("All subscriptions are successfully deleted.")
}

func decryptEndpoint(encrypted_endpoint string, version int, crypto_key []*fernet.Key) (uaid string, chid string, pub_key string, err error) {

	decoded_endpoint := fernet.VerifyAndDecrypt([]byte(encrypted_endpoint), TTL_DURATION, crypto_key)

	if decoded_endpoint == nil {

		decryption_err := errors.New("[Error]: Invalid Endpoint.")

		return "", "", "", decryption_err
	}

	hex_encoded_endpoint := hex.EncodeToString([]byte(decoded_endpoint))

	if version == 1 {

		uaid, chid = extractSubscription(hex_encoded_endpoint)
		return

	} else {

		uaid, chid, pub_key = extractVapidSubscription(hex_encoded_endpoint)
		return

	}

}

func extractVapidSubscription(subscription string) (uaid string, chid string, pub_key string) {

	uaid = subscription[0:32]
	unformatted_chid := subscription[32:64]
	chid = format_id(unformatted_chid)
	pub_key = subscription[64:]
	return

}
func format_id(chid string) string {

	/*

		This method will format extracted channelID from CassandraDB to be UUAID standard.

		Example:

		Input: 282c8b9184044db89f942c5972d0e55a

		Output: 282c8b91-8404-4db8-9f94-2c5972d0e55a


	*/

	var formated_chid string
	var first_part string
	for pos, char := range chid {

		first_part = first_part + string(char)

		if pos == 7 || pos == 11 || pos == 15 || pos == 19 {

			first_part = first_part + "-"
			formated_chid = formated_chid + first_part
			first_part = ""

		}

		if pos == 31 {

			formated_chid = formated_chid + first_part

		}

	}

	return formated_chid

}

func start_cassandra() (err error) {

	cluster := gocql.NewCluster(*cass_addr)

	cluster.Keyspace = "autopush"
	cluster.Consistency = gocql.Quorum

	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: *cass_user,
		Password: *cass_password,
	}
	cass_session, err = cluster.CreateSession()
	if err != nil {

		fmt.Println("start_cassandra: Failed due to error: %v", err)
		return err
	}
	//	defer cass_session.Close()
	fmt.Println("Connected to Cassandra")
	return nil
}

func DropUser(uaid string, session *gocql.Session) (dropped bool) {

	var err error

	dropped = false

	if err = session.Query("DELETE FROM autopush.router WHERE uaid = ? ",
		uaid).Consistency(gocql.One).Exec(); err != nil {
		fmt.Println(err)
		return false

	}

	dropped = true
	return

}
func extractSubscription(subscription string) (uaid string, chid string) {
	uaid = subscription[0:32]
	unformatted_chid := subscription[32:64]
	chid = format_id(unformatted_chid)
	return

}

func readFile(path string) []string {
	var str_arr []string
	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
		str_arr = append(str_arr, scanner.Text())
	}
	return str_arr
}
func Repad(subscription string) (paded_str string) {
	/*
		Adds padding to strings for base64 decoding.
		base64 encoding requires 'padding' to 4 octet boundries. 'padding' is a '='. So a string like 'abcde' is 5 octets, and would need to be padded to 'abcde==='.
	*/

	padings := len(subscription) % 4

	if padings != 0 {

		for i := 0; i < padings; i++ {

			subscription = subscription + "="

		}

	}

	paded_str = subscription

	return

}
