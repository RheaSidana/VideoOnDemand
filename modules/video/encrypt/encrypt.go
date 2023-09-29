package videoEncryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"vod/model"
)

func encrypt(videoLinks []model.VideoLinks) (
	map[model.VideoLinks]model.VideoLinks,
	error,
) {
	//create a folder to save encrypted files
	encryptionFolder, err := createFolderForEncryption()
	if err != nil {
		return map[model.VideoLinks]model.VideoLinks{}, err
	}

	//create map to save old and new video locations
	encodedToEncrypted := make(map[model.VideoLinks]model.VideoLinks)

	//encrypt videos
	encodedToEncrypted, err = encryptVideos(
		videoLinks,
		encryptionFolder, encodedToEncrypted,
	)
	if err != nil {
		return map[model.VideoLinks]model.VideoLinks{}, err
	}

	//delete encoded folders/files
	encodedFolder := getEncodedFolder(videoLinks[0].Link)
	err = deleteFolder(encodedFolder)
	if err != nil {
		return map[model.VideoLinks]model.VideoLinks{}, err
	}

	return encodedToEncrypted, nil
}

func getEncodedFolder(videoPath string) string {
	pathParts := strings.Split(videoPath, "/")
	len := len(pathParts)
	removePath := pathParts[len-2] + "/" + pathParts[len-1]

	return strings.Replace(videoPath, removePath, "", -1)
}

type output struct {
	InputVideo  model.VideoLinks
	OutputVideo string
	Err         error
}

func (o output) ToString() string {
	return fmt.Sprintf(
		"InputVideo: { %s }, OutputVideo: %s, Err: %v",
		o.InputVideo.ToString(),
		o.OutputVideo,
		o.Err,
	)
}

func encryptVideos(
	videoLinks []model.VideoLinks,
	folderPath string,
	encodedToEncrypted map[model.VideoLinks]model.VideoLinks,
) (
	map[model.VideoLinks]model.VideoLinks,
	error,
) {
	ch := make(chan output, len(videoLinks))

	var wg sync.WaitGroup

	for _, videoLink := range videoLinks {
		wg.Add(1)

		// multiple files
		go encryption(ch, &wg, videoLink, folderPath)
	}

	wg.Wait()
	close(ch)

	for res := range ch {
		fmt.Println("err: ", res.Err)
		if res.Err != nil {
			continue
		}
		encodedToEncrypted = updateEncodedToEncrypted(
			res.InputVideo,
			res.OutputVideo,
			encodedToEncrypted,
		)
	}
	if len(encodedToEncrypted) == 0 {
		return map[model.VideoLinks]model.VideoLinks{},
			errors.New("error occurred")
	}
	return encodedToEncrypted, nil
}

func encryption(
	ch chan output,
	wg *sync.WaitGroup,
	videoLink model.VideoLinks,
	folderPath string,
) {
	defer wg.Done()
	var outputObj output
	inputVideoPath := videoLink.Link
	outputVideoPath, err := getEncryptedOutputVideo(
		folderPath, videoLink.Link,
	)
	if err != nil {
		outputObj.Err = err
		ch <- outputObj
		return
	}
	outputObj.InputVideo = videoLink
	outputObj.OutputVideo = outputVideoPath
	// fmt.Println("\n-----------------\n",outputObj.ToString(), "\n-----------------")

	key := []byte(os.Getenv("SECRET2"))[:32]

	inputVideo, err := os.Open(inputVideoPath)
	if err != nil {
		outputObj.Err = err
		ch <- outputObj
		return
	}
	defer inputVideo.Close()

	outputVideo, err := os.Create(outputVideoPath)
	if err != nil {
		outputObj.Err = err
		ch <- outputObj
		return
	}
	defer outputVideo.Close()

	//cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Error creating AES cipher:", err)
		outputObj.Err = err
		// ch <- outputObj
		return
	}

	//random initialization vector (IV)
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println("Error generating IV:", err)
		outputObj.Err = err
		ch <- outputObj
		return
	}

	// Write the IV to the beginning of the encrypted file
	_, err = outputVideo.Write(iv)
	if err != nil {
		fmt.Println("Error writing IV:", err)
		outputObj.Err = err
		ch <- outputObj
		return
	}

	//cipher block mode
	cipherStream := cipher.NewCFBEncrypter(block, iv)

	//buffer for reading and writing chunks of data
	bufferSize := 1024
	buffer := make([]byte, bufferSize)

	err = encrypting(
		*inputVideo, *outputVideo,
		buffer, cipherStream,
	)
	if err != nil {
		fmt.Println("Error in encrypting", err)
		outputObj.Err = err
		ch <- outputObj
		return
	}

	ch <- outputObj
}

func updateEncodedToEncrypted(old model.VideoLinks,
	newLink string,
	encodedToEncrypted map[model.VideoLinks]model.VideoLinks,
) map[model.VideoLinks]model.VideoLinks {
	encyptedVideoLink := old
	encyptedVideoLink.Link = newLink
	// encyptedVideoLink := model.VideoLinks{
	// 	VideoMetaDataID: old.VideoMetaDataID,
	// 	EncodedLink:     newLink,
	// }
	encodedToEncrypted[old] = encyptedVideoLink

	return encodedToEncrypted
}

func encrypting(
	inputVideo, outputVideo os.File,
	buffer []byte,
	cipherStream cipher.Stream,
) error {
	// Encrypt and write the video data
	for {
		n, err := inputVideo.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading input file:", err)
			// return map[model.VideoLinks]model.VideoLinks{}, err
			return err
		}

		if n == 0 {
			break
		}

		// Encrypt the chunk of data
		cipherStream.XORKeyStream(buffer[:n], buffer[:n])

		// Write the encrypted data to the output file
		_, err = outputVideo.Write(buffer[:n])
		if err != nil {
			fmt.Println("Error writing encrypted data:", err)
			// return map[model.VideoLinks]model.VideoLinks{}, err
			return err
		}
	}

	return nil
}

func getEncryptedOutputVideo(folderPath, videoLink string) (string, error) {
	if !fileExists(videoLink) {
		return "", errors.New("file does not exists")
	}

	resolution, bitrate, vidTitle, extension := extractVideoDetails(videoLink)
	return getOutputVideo(
		resolution, bitrate, vidTitle,
		extension, folderPath,
	), nil
}

func getOutputVideo(
	resolution, bitrate, vidTitle,
	extension, folderPath string,
) string {
	loc := folderPath + "/" + vidTitle + "/" + bitrate
	_ = createFolder(loc)
	file := resolution + "encrypt" + "." + extension

	return loc + "/" + file
}

func extractVideoDetails(inputVideoPath string) (
	string, string, string, string,
) {
	pathVals := strings.Split(inputVideoPath, "/")
	len := len(pathVals)

	resolution := strings.Split(pathVals[len-1], ".")[0]
	extension := strings.Split(pathVals[len-1], ".")[1]
	bitrate := pathVals[len-2]
	vidTitle := pathVals[len-3]

	return resolution,
		bitrate,
		vidTitle,
		extension
}

func createFolderForEncryption() (string, error) {
	saveToLocation := BASE_DIR
	folder := "Encryption"

	folderPath := saveToLocation + folder
	err := createFolder(folderPath)
	if err != nil {
		return "", err
	}

	return folderPath, nil
}
