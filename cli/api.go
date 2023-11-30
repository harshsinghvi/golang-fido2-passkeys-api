package main

import (
	"bytes"
	"encoding/json"
	"harshsinghvi/golang-fido2-passkeys-api/lib/crypto"
	"harshsinghvi/golang-fido2-passkeys-api/utils"
	"io"
	"log"
	"net/http"
)

func login(passkeyId string, serverUrl string) {
	url := getServerURL(serverUrl)

	publicKey, err := crypto.LoadPublicKeyFromFile(PUBLIC_KEY_PATH)
	e(err)
	publicKeyStr, err := crypto.PublicKeyToString(publicKey)
	e(err)

	client := http.Client{}
	req, err := http.NewRequest("GET", url+"/api/login/request-challenge", nil)
	req.Header.Add("Public-Key", publicKeyStr)
	e(err)
	resp, err := client.Do(req)
	e(err)

	// passkeyID := getPasskeyId(passkeyId)
	// resp, err := http.Get(url + "/api/login/request-challenge"+passkeyID)
	// e(err)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// if the status code is not 200, we should log the status code and the
		// status string, then exit with a fatal error
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			// print the response
			data, err := io.ReadAll(resp.Body)
			e(err)

			log.Fatalf("BAD Request status code error: %d %s \n %s", resp.StatusCode, resp.Status, string(data))
		}
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	log.Println("INFO: Challenge Created, Verify to login.")
	data := parseJson(resp)
	challengeID, _ := data["data"].(map[string]interface{})["ChallengeID"].(string)
	challengeStr, _ := data["data"].(map[string]interface{})["ChallengeString"].(string)

	verifyChallenge(url, "passkeyID - not to be used", challengeID, challengeStr)
}

func userReg(name string, email string, serverUrl string) {
	url := getServerURL(serverUrl)

	// privateKey, err := crypto.LoadPrivateKeyFromFile(PRIVATE_KEY_PATH)
	// e(err, "Error in reading Private Key file, pelase generate new files by cli gen")
	publicKey, err := crypto.LoadPublicKeyFromFile(PUBLIC_KEY_PATH)
	e(err, "Error in reading Public Key file, pelase generate new files by cli gen")
	// privateKeyStr, err := crypto.PrivateKeyToString(privateKey) // Convert private key to string
	// e(err)
	publicKeyStr, err := crypto.PublicKeyToString(publicKey) // Convert public key to string
	e(err)

	// jsonValue, err := json.Marshal(map[string]string{"Name": name, "Email": email, "PublicKey": publicKeyStr, "PrivateKey": privateKeyStr})
	jsonValue, err := json.Marshal(map[string]string{"Name": name, "Email": email, "PublicKey": publicKeyStr})
	e(err)
	resp, err := http.Post(url+"/api/registration/user", "application/json", bytes.NewBuffer(jsonValue))
	e(err)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// if the status code is not 200, we should log the status code and the
		// status string, then exit with a fatal error
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			// print the response
			data, err := io.ReadAll(resp.Body)
			e(err)

			log.Fatalf("BAD Request status code error: %d %s \n %s", resp.StatusCode, resp.Status, string(data))
		}
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	log.Println("INFO: User Created.")

	data := parseJson(resp)
	passkeyID, _ := data["data"].(map[string]interface{})["PasskeyID"].(string)
	challengeID, _ := data["data"].(map[string]interface{})["ChallengeID"].(string)
	challengeStr, _ := data["data"].(map[string]interface{})["ChallengeString"].(string)
	verifyChallenge(url, passkeyID, challengeID, challengeStr)
}

func verifyChallenge(url string, passkeyID string, challengeID string, challengeStr string) {
	challenge := decrypt(challengeStr)
	challengeSolution, ok := utils.SolveChallengeString(challenge)
	challengeSignature := sign(challengeSolution)
	if !ok {
		log.Fatal("Something went wrong cant solve challenge.")
	}

	log.Println("Challenge: ", challenge)
	log.Println("Solution: ", challengeSolution)

	jsonValue, err := json.Marshal(map[string]string{"ChallengeID": challengeID, "ChallengeSignature": challengeSignature})
	e(err)
	resp, err := http.Post(url+"/api/login/verify-challenge", "application/json", bytes.NewBuffer(jsonValue))
	e(err)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// if the status code is not 200, we should log the status code and the
		// status string, then exit with a fatal error
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			// print the response
			data, err := io.ReadAll(resp.Body)
			e(err)

			log.Fatalf("BAD Request status code error: %d %s \n %s", resp.StatusCode, resp.Status, string(data))
		}
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	data := parseJson(resp)
	token, _ := data["data"].(map[string]interface{})["Token"].(string)
	config := Config{
		ServerUrl:   url,
		PasskeyID:   passkeyID,
		AccessToken: token,
	}
	writeConfigToFile(config, CONFIG_PATH)
	// Write config file
	log.Println("INFO: Challenge Verification Successful, Passkey Verified Access token stored.")
}

func parseJson(resp *http.Response) map[string]interface{} {
	var data map[string]interface{}
	resBody, err := io.ReadAll(resp.Body)
	e(err)
	err = json.Unmarshal(resBody, &data)
	e(err)
	return data
}
