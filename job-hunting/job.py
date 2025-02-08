package main

import (
	"fmt"
	"log"
	"net/smtp"
	"crypto/tls"
	"github.com/xuri/excelize/v2"
	"regexp"
	"strings"
)

// Email Configuration
const (
	SMTPServer    = "smtp.gmail.com"
	SMTPPort      = "465"
	EmailAddress  = "xx@gmail.com"  // Replace with your Gmail address
	EmailPassword = "xx"           // Replace with your Gmail app password (not your main Gmail password)
	GoogleDriveCVLink = "https://drive.google.com/file/d/view?usp=sharing"  // Replace with your CV's Google Drive link
)

// Validate email format using regular expression
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// Read the recipient and company name from an Excel sheet
func readExcelData(filePath string) ([]string, []string, []string, error) {
	// Open the Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to open Excel file: %v", err)
	}

	// Create slices to hold the data
	var receiverEmails, companyNames, hrNames []string

	// Loop through the first 100 rows (1-indexed: row 2 to row 101)
	for row := 57; row <= 101; row++ {
		receiverEmail, err := f.GetCellValue("Sheet1", fmt.Sprintf("D%d", row))
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to read receiver email from row %d: %v", row, err)
		}
		companyName, err := f.GetCellValue("Sheet1", fmt.Sprintf("C%d", row))
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to read company name from row %d: %v", row, err)
		}
		hrName, err := f.GetCellValue("Sheet1", fmt.Sprintf("B%d", row))
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to read HR name from row %d: %v", row, err)
		}

		// Append the values to the slices
		receiverEmails = append(receiverEmails, receiverEmail)
		companyNames = append(companyNames, companyName)
		hrNames = append(hrNames, hrName)
	}

	return hrNames, companyNames, receiverEmails, nil
}

// HTML content template with dynamic subject and body
func generateHTMLBody(hrName string, companyName string) string {
	return fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Email Sample</title>
	</head>
	<body>
		<div>
			<p>Hi %s,<p>
			<P>I hope this email finds you well. <br>I am Parth Bagda, a graduate of IIT Guwahati. 
			I am currently working as a Backend Developer at Practo Technologies, with 9 months of experience in Java,
			Spring Boot, and microservices. I am reaching out to inquire about potential software opportunities at %s. 
			Please let me know if there are any openings that align with my skills and experience.
			<p>Please find my cv here: <a href="%s">parth_bagda_sde</a></p>
			<br> Thank you for your time and consideration. </p>
			<p>Parth Bagda <br> 6355351675 </p>
		</div>
	</body>
	</html>
	`, hrName, companyName, GoogleDriveCVLink)
}

// sendEmail sends an HTML email with a specified subject, sender, and receiver.
func sendEmail(to string, senderEmail string, password string, hrName string, companyName string) (string, error) {
	// Create the email message headers and body
	subject := fmt.Sprintf("Seeking Software Development Opportunities at %s", companyName)
	body := generateHTMLBody(hrName, companyName)

	// Construct the message
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"Importance: high\r\n" +
		"\r\n" + body)

	// Connect to the SMTP server
	conn, err := tls.Dial("tcp", SMTPServer+":"+SMTPPort, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         SMTPServer,
	})
	if err != nil {
		return "", fmt.Errorf("failed to connect to the server: %v", err)
	}

	// Create SMTP client
	client, err := smtp.NewClient(conn, SMTPServer)
	if err != nil {
		return "", fmt.Errorf("failed to create SMTP client: %v", err)
	}

	// Authenticate using sender's email and password
	auth := smtp.PlainAuth("", senderEmail, password, SMTPServer)
	if err := client.Auth(auth); err != nil {
		return "", fmt.Errorf("SMTP authentication failed: %v", err)
	}

	// Set the sender and recipient addresses
	if err := client.Mail(senderEmail); err != nil {
		return "", fmt.Errorf("failed to set sender: %v", err)
	}
	if err := client.Rcpt(to); err != nil {
		return "", fmt.Errorf("failed to set recipient: %v", err)
	}

	// Send the email body
	writer, err := client.Data()
	if err != nil {
		return "", fmt.Errorf("failed to get writer for email data: %v", err)
	}
	_, err = writer.Write(msg)
	if err != nil {
		return "", fmt.Errorf("failed to write email data: %v", err)
	}

	// Close the client connection
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close writer: %v", err)
	}
	client.Quit()

	return "Email sent successfully", nil
}



func main() {
	// Read the recipient, company name, and HR name from the Excel file
	hrNames, companyNames, receiverEmails, err := readExcelData("imRBz7Hk9wXHKGSyTjvtEB.xlsx") // Replace with the correct file path
	if err != nil {
		log.Fatalf("Error reading Excel data: %v", err)
	}

	// Prepare slices to collect invalid emails
	var invalidEmails, invalidEmails550 [][]string

	// Print and send emails
	for i := 0; i < len(receiverEmails); i++ {
		receiverEmail := receiverEmails[i]
		hrName := hrNames[i]
		companyName := companyNames[i]

		// Check if the email is valid
		if !isValidEmail(receiverEmail) {
			invalidEmails = append(invalidEmails, []string{hrName, companyName, receiverEmail, "Invalid email format"})
			continue
		}

		// Send email and capture the response
		_, err := sendEmail(receiverEmail, EmailAddress, EmailPassword, hrName, companyName)
		if err != nil {
			// Check for specific 550 5.4.1 error
			if strings.Contains(err.Error(), "550 5.4.1") {
				invalidEmails550 = append(invalidEmails550, []string{hrName, companyName, receiverEmail, "550 5.4.1 Error"})
			} else {
				log.Printf("Error sending email to %s: %v", receiverEmail, err)
			}
		} else {
			fmt.Printf("Email sent successfully to %s!\n", receiverEmail)
		}
	}

	// Save invalid emails with 550 error to Excel
	if len(invalidEmails550) > 0 {
		if err := writeInvalidEmailsToExcel(invalidEmails550, "invalid_emails_550.xlsx"); err != nil {
			log.Printf("Error writing invalid emails to Excel: %v", err)
		} else {
			fmt.Println("Invalid emails with 550 error saved to 'invalid_emails_550.xlsx'")
		}
	}
}
