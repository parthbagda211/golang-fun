package main

import (
	"fmt"
	"log"
	"net/smtp"
	"crypto/tls"
	"strings"
	"github.com/xuri/excelize/v2"
)

// Email Configuration
const (
	SMTPServer    = "smtp.gmail.com"
	SMTPPort      = "465"
	EmailAddress  = "parthbagda94@gmail.com"  // Replace with your Gmail address
	EmailPassword = "xxif rrpr xfvy hrlm"           // Replace with your Gmail app password (not your main Gmail password)
	GoogleDriveCVLink = "https://drive.google.com/file/d/1FdhlxPh1aHu4ctOu8e7jlpYKLtKzn-vu/view"  // Replace with your CV's Google Drive link
	LinkedinLink = "https://www.linkedin.com/in/parth-bagda-231a6a205"
)

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
	for row := 2; row <= 14; row++ {
		receiverEmail, err := f.GetCellValue("Sheet1", fmt.Sprintf("C%d", row))
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to read receiver email from row %d: %v", row, err)
		}
		companyName, err := f.GetCellValue("Sheet1", fmt.Sprintf("B%d", row))
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to read company name from row %d: %v", row, err)
		}
		hrName, err := f.GetCellValue("Sheet1", fmt.Sprintf("A%d", row))
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
			Spring Boot, and microservices. I am reaching out to regarding potential software opportunities at %s. 
			Please let me know if there are any openings that align with my skills and experience. <br> I am available for immediate joining and would be excited to contribute my skills and experience to your team. 
			Please find my CV attached for your reference.
			<p>Please find my cv here: <a href="%s">parth_bagda_sde</a></p>
			<br> Thank you for your time and consideration. </p>
			<p>Parth Bagda <br> 6355351675 <br>  <a href="%s">LinkedIn</a></p> </p> 
		</div>
	</body>
	</html>
	`, hrName, companyName, GoogleDriveCVLink,LinkedinLink)
}

// sendEmail sends an HTML email with a specified subject, sender, and receiver.
func sendEmail(to string, senderEmail string, password string, hrName string, companyName string) error {
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
		return fmt.Errorf("failed to connect to the server: %v", err)
	}

	// Create SMTP client
	client, err := smtp.NewClient(conn, SMTPServer)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}

	// Authenticate using sender's email and password
	auth := smtp.PlainAuth("", senderEmail, password, SMTPServer)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %v", err)
	}

	// Set the sender and recipient addresses
	if err := client.Mail(senderEmail); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	// Send the email body
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get writer for email data: %v", err)
	}
	_, err = writer.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write email data: %v", err)
	}

	// Close the client connection
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}
	client.Quit()

	return nil
}

func createRejectedEmailsSheet(rejectedEmails []string) error {
	// Create a new Excel file
	f := excelize.NewFile()

	// Set headers for the sheet
	f.SetCellValue("Sheet1", "A1", "Rejected Email Addresses")

	// Write rejected emails to the sheet
	for i, email := range rejectedEmails {
		cell := fmt.Sprintf("A%d", i+2)
		f.SetCellValue("Sheet1", cell, email)
	}

	// Save the file
	if err := f.SaveAs("rejected_emails.xlsx"); err != nil {
		return fmt.Errorf("failed to save Excel file: %v", err)
	}

	return nil
}

func main() {
	// Read the recipient, company name, and HR name from the Excel file
	hrNames, companyNames, receiverEmails, err := readExcelData("jobs_1738682524976.xlsx") // Replace with the correct file path
	if err != nil {
		log.Fatalf("Error reading Excel data: %v", err)
	}

	var rejectedEmails []string

	// Print the first few records to verify
	for i := 0; i < len(receiverEmails); i++ {
		fmt.Printf("Row %d: %s | %s | %s\n", i+2, hrNames[i], companyNames[i], receiverEmails[i])
	}

	
	for i := 0; i < len(receiverEmails); i++ {
		err := sendEmail(receiverEmails[i], EmailAddress, EmailPassword, hrNames[i], companyNames[i])
		if err != nil {
			// Check if error is due to 550 5.4.1
			if strings.Contains(err.Error(), "5.5.2") {
				rejectedEmails = append(rejectedEmails, receiverEmails[i])
				log.Printf("Recipient %s rejected. Added to rejected emails list.", receiverEmails[i])
			} else {
				log.Printf("Error sending email to %s: %v", receiverEmails[i], err)
			}
		} else {
			fmt.Printf("Email sent successfully to %s!\n", receiverEmails[i])
		}
	}

	// If there are rejected emails, create a new Excel file
	if len(rejectedEmails) > 0 {
		if err := createRejectedEmailsSheet(rejectedEmails); err != nil {
			log.Printf("Failed to create rejected emails sheet: %v", err)
		} else {
			fmt.Println("Rejected email addresses saved to rejected_emails.xlsx")
		}
	}
}
