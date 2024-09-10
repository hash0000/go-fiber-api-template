package tale

import (
	"encoding/json"
	"fmt"
	"go-fiber-api-template/app/common/constants"
	"go-fiber-api-template/app/common/helpers"
	"go-fiber-api-template/app/common/types"
	"go-fiber-api-template/app/modules/user"
	"image/color"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"golang.org/x/net/proxy"
)

func addFinalButtons(pdf *fpdf.Fpdf) {
	pdf.SetMargins(0, 0, 0)

	pdf.ImageOptions(filepath.Join(constants.ImageDir, "feedback_button.png"), 30, 210, 150, 30, false, fpdf.ImageOptions{ReadDpi: true, ImageType: "PNG"}, 0, os.Getenv("FEEDBACK_LINK"))

	pdf.ImageOptions(filepath.Join(constants.ImageDir, "create_tale_button.png"), 30, 210+35, 150, 30, false, fpdf.ImageOptions{ReadDpi: true, ImageType: "PNG"}, 0, os.Getenv("TG_BOT_LINK"))

	pdf.SetMargins(10.0, 10.0, 10.0)
}

func addFinalTrialPage(pdf *fpdf.Fpdf) {
	pdf.AddPage()

	pdf.SetMargins(0, 0, 0)

	pdf.ImageOptions(filepath.Join(constants.ImageDir, "trial_end_page.png"), 0, 0, 210.0, 297.0, false, fpdf.ImageOptions{ReadDpi: true, ImageType: "PNG"}, 0, os.Getenv("TRIAL_PAY_LINK"))

	pdf.SetMargins(10.0, 10.0, 10.0)
}

func getTaleFromChatGpt(prompt string) (types.StoryChapter, error) {
	var responseFormat interface{}
	if err := json.Unmarshal([]byte(constants.GptJsonOption), &responseFormat); err != nil {
		slog.Error("Error unmarshaling", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return types.StoryChapter{}, err
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"model": "gpt-4o-2024-08-06",
		"messages": []map[string]string{
			{"role": "system", "content": constants.GptSystemPromptV1},
			{"role": "user", "content": prompt},
		},
		"response_format": responseFormat,
		"max_tokens":      3024,
	})
	if err != nil {
		slog.Error("Error marshalling", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return types.StoryChapter{}, err
	}

	client := resty.New()

	if os.Getenv("HTTP_PROXY_TCP") != "" {
		proxyURL, err := url.Parse(os.Getenv("HTTP_PROXY_TCP"))
		if err != nil {
			slog.Error("Error parsing proxy URL", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
			return types.StoryChapter{}, err
		}

		client.SetProxy(proxyURL.String())
	}

	if os.Getenv("HTTPS_PROXY_SOCKS5") != "" {
		proxyURL, _ := url.Parse(os.Getenv("HTTPS_PROXY_SOCKS5"))
		dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			slog.Error("Error parsing proxy URL", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
			return types.StoryChapter{}, err
		}

		transport := &http.Transport{Dial: dialer.Dial}

		client.SetTransport(transport)
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY")).
		SetBody(requestBody).
		Post(constants.Url.ChatGptRequest)
	if err != nil {
		slog.Error(fmt.Sprintf("Error sending POST request to: %s", constants.Url.ChatGptRequest), slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return types.StoryChapter{}, err
	}

	if resp.StatusCode() != 200 {
		slog.Error(fmt.Sprintf("Unexpected status code: %d, body: %s, from: %s", resp.StatusCode(), resp.String(), constants.Url.ChatGptRequest), slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return types.StoryChapter{}, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var gptResponse types.OpenAIResponseType
	err = json.Unmarshal(resp.Body(), &gptResponse)
	if err != nil {
		slog.Error("Error unmarshalling response body", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return types.StoryChapter{}, err
	}

	if len(gptResponse.Choices) == 0 {
		slog.Error("No choices returned from GPT", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return types.StoryChapter{}, err
	}

	var story types.StoryChapter
	err = json.Unmarshal([]byte(gptResponse.Choices[0].Message.Content), &story)
	if err != nil {
		slog.Error("Error unmarshalling story chapters", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return types.StoryChapter{}, err
	}

	return story, nil
}

func generatePdf(tale types.StoryChapter, storeDir string, sessionID string, trialTale int8) (string, error) {
	pdf := setupPdf()

	var chaptersToOrder []types.Chapter
	var chaptersToWait []types.Chapter
	skipImageIndex := 0

	if trialTale == constants.TrialTale.No {
		chaptersToOrder = tale.Chapters
		chaptersToWait = tale.Chapters
	} else if trialTale == constants.TrialTale.Start {
		chaptersToOrder = append(chaptersToOrder, tale.Chapters[0])
		chaptersToWait = append(chaptersToWait, tale.Chapters[0])
	} else if trialTale == constants.TrialTale.Finish {
		chaptersToOrder = tale.Chapters[1:]
		chaptersToWait = tale.Chapters
		skipImageIndex = 1
	}

	for index, chapter := range chaptersToOrder {
		err := getImageFromFabula(chapter.PicGeneration, index+skipImageIndex, sessionID)
		if err != nil {
			slog.Error("Error creating order for one of image", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
			return "", err
		}
	}

	imagePaths, err := waitForAllImages(sessionID, chaptersToWait)
	if err != nil {
		slog.Error("Error running waitForAllImages", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return "", err
	}

	const (
		pageWidth    = 210.0
		imageHeight  = 120.0
		textStartY   = imageHeight + 13.0
		titleTextGap = 8.0
		marginWidth  = 10.0
	)

	err = addTitlePage(pdf, tale.Title)
	if err != nil {
		slog.Error("Error adding title page", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return "", err
	}
	_, lineHt := pdf.GetFontSize()
	lineHt *= 1.5

	pdf.SetFont("MontserratAlternates", "", 16)

	for i, chapter := range chaptersToWait {
		pdf.AddPage()

		if i < len(imagePaths) {
			addImageToPdf(pdf, imagePaths[i], pageWidth, imageHeight)
		}

		pdf.SetY(textStartY)

		pdf.SetFont("MontserratAlternates", "B", 32)
		pdf.MultiCell(pageWidth-20, 11, chapter.Title, "", "", false)

		pdf.Ln(titleTextGap)

		pdf.SetFont("MontserratAlternates", "", 19)
		pdf.MultiCell(pageWidth-2*marginWidth, 8, chapter.Text, "", "", false)
	}

	if trialTale == constants.TrialTale.No || trialTale == constants.TrialTale.Finish {
		pdf.AddPage()
		pdf.SetY(20)
		pdf.SetFont("MontserratAlternates", "B", 32)
		pdf.MultiCell(pageWidth-20, 11, "Вопросы для обсуждения", "", "", false)
		pdf.Ln(titleTextGap)
		for _, question := range tale.QuestionsAboutTale {
			pdf.SetFont("MontserratAlternates", "", 20)
			pdf.MultiCell(pageWidth-2*marginWidth, 8, fmt.Sprintf("• %s", question), "", "", false)

			pdf.Ln(titleTextGap)
		}

		addFinalButtons(pdf)
	} else {
		addFinalTrialPage(pdf)
	}

	outputPath, err := savePdf(pdf, storeDir, sessionID, trialTale)
	if err != nil {
		slog.Error("Error running savePdf", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return "", err
	}

	return outputPath, nil
}

func setPageBackgroundColor(pdf *fpdf.Fpdf, bgColor color.RGBA) {
	pdf.SetFillColor(int(bgColor.R), int(bgColor.G), int(bgColor.B))
	pdf.Rect(0, 0, 210, 297, "F")
}

func setupPdf() *fpdf.Fpdf {
	pdf := fpdf.New("P", "mm", "A4", "")
	fontRDir := filepath.Join(constants.FontsDir, "MontserratAlternates-Regular.ttf")
	fontBDir := filepath.Join(constants.FontsDir, "MontserratAlternates-Bold.ttf")
	fontIDir := filepath.Join(constants.FontsDir, "MontserratAlternates-Italic.ttf")
	fontPribambasDir := filepath.Join(constants.FontsDir, "Pribambas-Regular.ttf")
	pdf.AddUTF8Font("MontserratAlternates", "", fontRDir)
	pdf.AddUTF8Font("MontserratAlternates", "B", fontBDir)
	pdf.AddUTF8Font("MontserratAlternates", "I", fontIDir)
	pdf.AddUTF8Font("Pribambas", "", fontPribambasDir)
	pdf.SetFont("MontserratAlternates", "", 16)

	pdf.SetHeaderFunc(func() {
		setPageBackgroundColor(pdf, color.RGBA{R: 254, G: 249, B: 243, A: 255})
	})

	return pdf
}

func addTitlePage(pdf *fpdf.Fpdf, title string) error {
	pdf.AddPage()

	headerHeight := 10.0
	titleHeight := 20.0
	spacing := 5.0
	startY := 105.0
	pageWidth := 210.0
	marginWidth := 10.0

	pdf.SetY(startY)
	pdf.SetTextColor(255, 204, 0)
	pdf.SetFont("Pribambas", "", 45)
	pdf.CellFormat(0, headerHeight, os.Getenv("COMPANY_NAME"), "", 1, "C", false, 0, "")

	startY += headerHeight + spacing + spacing

	pdf.SetY(startY)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Pribambas", "", 60)
	pdf.SetLeftMargin(marginWidth)
	pdf.SetRightMargin(marginWidth)
	pdf.MultiCell(pageWidth-2*marginWidth, titleHeight, title, "", "C", false)

	startY = pdf.GetY() + spacing

	return nil
}

func addImageToPdf(pdf *fpdf.Fpdf, imagePath string, pageWidth float64, imageHeight float64) {
	pdf.SetMargins(0, 0, 0)

	pdf.ImageOptions(imagePath, 0, 0, pageWidth, imageHeight, false, fpdf.ImageOptions{ReadDpi: true, ImageType: "JPG"}, 0, "")

	pdf.SetMargins(10.0, 10.0, 10.0)
}

func savePdf(pdf *fpdf.Fpdf, storeDir string, sessionID string, trial int8) (string, error) {
	var fileName string

	if trial == constants.TrialTale.Start {
		fileName = fmt.Sprintf("trial_tale_%s.pdf", sessionID)
	} else {
		fileName = fmt.Sprintf("tale_%s.pdf", sessionID)
	}

	relativePathTale := filepath.Join(storeDir, fileName)

	if err := pdf.OutputFileAndClose(relativePathTale); err != nil {
		slog.Error("Error saving PDF", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return "", err
	}

	return fileName, nil
}

func getImageFromFabula(prompt string, index int, sessionID string) error {
	requestBody, err := json.Marshal(map[string]interface{}{
		"sd_model":            "dreamshaper_xl_light",
		"sampling_steps":      9,
		"guidance_scale":      3,
		"width":               816,
		"height":              512,
		"product_id":          15,
		"prompt":              prompt + ", fantasy style and illustrations for children books, fantasy art, bright colors, soft lighting, detailed background, the style of digital painting, the atmosphere of magic and adventure",
		"nsfw_check":          true,
		"negative_prompt":     "realistic, photorealistic, 3d, worst quality, low quality, jpeg and png artifacts, interlocked fingers, wrong eyes, naked, nipples, nudes, tits, lowres, over-smooth, text, watermark, words, brands",
		"webhook_url":         os.Getenv("WEBHOOK_TXT_TO_IMG_URL"),
		"custom_webhook_data": fmt.Sprintf("{\"session_id\": \"%s\", \"index\": \"%s\"}", sessionID, strconv.Itoa(index)),
	})
	if err != nil {
		slog.Error("Error marshalling request body for image generation", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return err
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Api-Key "+os.Getenv("FABULA_API_KEY")).
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Post(constants.Url.FabulaCreateText2Image)

	if err != nil {
		slog.Error("Error sending POST request for image generation", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		slog.Error(fmt.Sprintf("Unexpected status code: %d, body: %s, to: %s", resp.StatusCode(), resp.String(), constants.Url.FabulaCreateText2Image), slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode(), resp.String())
	}

	return nil
}

func getFileFromUrl(url string) (string, error) {
	client := resty.New()

	resp, err := client.R().Get(url)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while sending request to: %s, Code: %d, body: %s", url, resp.StatusCode(), resp.String()), slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return "", err
	}

	fileData := resp.Body()

	randomId := uuid.New().String()
	filePath := filepath.Join(constants.TempImgToTextJsonDir, fmt.Sprintf("img_to_text_%s.json", randomId))

	err = os.WriteFile(filePath, fileData, 0644)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while saving file from: %s", url), slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return "", fmt.Errorf("error while saving file from: %s", url)
	}

	return filePath, nil
}

func convertJsonToImgToTextJsonType(filePath string) (types.ImgToTextJsonType, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		slog.Error("Error reading file", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return types.ImgToTextJsonType{}, fmt.Errorf("error reading file")
	}

	jsonContent := string(file)

	var result types.ImgToTextJsonType
	err = json.Unmarshal([]byte(jsonContent), &result)
	if err != nil {
		slog.Error("Error unmarshalling story chapters", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return types.ImgToTextJsonType{}, fmt.Errorf("error unmarshalling story chapters")
	}

	err = os.Remove(filePath)
	if err != nil {
		slog.Error(fmt.Sprintf("Error deleting file: %s", filePath), slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
	}

	return result, nil
}

func waitForAllImages(sessionID string, chapters []types.Chapter) ([]string, error) {
	time.Sleep(3 * time.Second)

	numImages := len(chapters)
	imagePaths := make([]string, numImages)

	timeout := time.After(20 * time.Minute)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			slog.Error(fmt.Sprintf("Timeout: not all images were available after 20 minutes. SessionID: %s", sessionID), slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)))
			return nil, fmt.Errorf("timeout: not all images were available after 20 minutes")
		case <-ticker.C:
			allImagesAvailable := true

			for i := 0; i < numImages; i++ {
				imageName := fmt.Sprintf("image_%d.jpg", i)
				imagePath := filepath.Join(constants.StoreDir, sessionID, imageName)

				if _, err := os.Stat(imagePath); os.IsNotExist(err) {
					allImagesAvailable = false
					break
				} else {
					imagePaths[i] = imagePath
				}
			}

			if allImagesAvailable {
				return imagePaths, nil
			}
		}
	}
}

func checkIfUserExistsAndHasToken(userID int64) (int8, error) {
	userEntity, status, err := user.SelectOne(nil, nil, userID)
	if err != nil {
		return 0, err
	}

	if status == constants.NotFound {
		return 1, nil
	}

	if *userEntity.TokenNumber == 0 {
		return 2, nil
	}

	return 0, nil
}

func checkIfUserAllowTrial(userID int64) (int8, error) {
	userEntity, status, err := user.SelectOne(nil, nil, userID)
	if err != nil {
		return 0, err
	}

	if status == constants.NotFound {
		return 1, nil
	}

	if *userEntity.UseTrial == true {
		return 2, nil
	}

	return 0, nil
}

func saveToFile(data interface{}, dir, fileName string) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	filePath := filepath.Join(dir, fileName)

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling data to JSON: %v", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("error writing JSON to file: %v", err)
	}

	return nil
}

func loadTaleFromFile(generationId string) (types.StoryChapter, error) {
	filePath := filepath.Join(constants.StoreDir, generationId, fmt.Sprintf("tale_%s.json", generationId))

	file, err := os.ReadFile(filePath)
	if err != nil {
		return types.StoryChapter{}, fmt.Errorf("error reading file: %v", err)
	}

	var tale types.StoryChapter
	err = json.Unmarshal(file, &tale)
	if err != nil {
		return types.StoryChapter{}, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return tale, nil
}

func sendTaleToTelegramBot(userID int64, fileDownloadUrl string) {
	requestBody, err := json.Marshal(map[string]string{
		"id":  strconv.Itoa(int(userID)),
		"url": fileDownloadUrl,
	})
	if err != nil {
		slog.Error("Error unmarshalling JSON", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Post(constants.Url.UploadTaleTg)

	if err != nil {
		slog.Error(fmt.Sprintf("Error sending POST to: %s", constants.Url.UploadTaleTg), slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
	}

	if resp.StatusCode() != http.StatusOK {
		slog.Error("Error response status code", slog.Int("status", resp.StatusCode()), slog.String("body", resp.String()), slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)))
	}
}

func sendTaleToTelegramAnalytic(fileDownloadUrl string) error {
	requestBodyAnalytic, err := json.Marshal(map[string]string{
		"text": fileDownloadUrl,
	})
	if err != nil {
		slog.Error("Error unmarshalling JSON", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return err
	}

	clientAnalytic := resty.New()
	resp, err := clientAnalytic.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBodyAnalytic).
		Post(constants.Url.UploadTgAnalytic)

	if err != nil {
		slog.Error(fmt.Sprintf("Error sending POST to: %s", constants.Url.UploadTaleTg), slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		slog.Error("Error response status code", slog.Int("status", resp.StatusCode()), slog.String("body", resp.String()), slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)))
		return fmt.Errorf("error sending POST to: %s", constants.Url.UploadTaleTg)
	}

	return nil
}
