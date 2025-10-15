package examples
package examples

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/protocyber/kelasgo-api/internal/middleware"
)

// Example 1: Generate API Links in Responses
func GetStudentHandler(c *gin.Context) {
	// Get app URL from context
	appURL, _ := c.Get("app_url") // "http://10.201.0.25:8080"
	
	student := getStudentFromDB() // your logic here
	
	response := gin.H{
		"success": true,
		"data": gin.H{
			"student": student,
			"links": gin.H{
				"self":        fmt.Sprintf("%s/v1/students/%d", appURL, student.ID),
				"edit":        fmt.Sprintf("%s/v1/students/%d", appURL, student.ID),
				"enrollments": fmt.Sprintf("%s/v1/students/%d/enrollments", appURL, student.ID),
				"grades":      fmt.Sprintf("%s/v1/students/%d/grades", appURL, student.ID),
			},
		},
	}
	c.JSON(200, response)
}

// Example 2: Email Templates with API Links
func SendWelcomeEmail(c *gin.Context, userEmail string, resetToken string) {
	appURL, _ := c.Get("app_url")
	
	// Generate reset password link
	resetURL := fmt.Sprintf("%s/v1/auth/reset-password?token=%s", appURL, resetToken)
	
	// Generate verification link  
	verifyURL := fmt.Sprintf("%s/v1/auth/verify-email?token=%s", appURL, resetToken)
	
	emailBody := fmt.Sprintf(`
		Welcome to KelasGo!
		
		Please verify your email: %s
		
		If you need to reset your password: %s
	`, verifyURL, resetURL)
	
	// Send email with these links
	sendEmail(userEmail, "Welcome to KelasGo", emailBody)
}

// Example 3: Webhook URLs
func CreateWebhookHandler(c *gin.Context) {
	appURL, _ := c.Get("app_url")
	
	webhook := gin.H{
		"id":  "webhook_123",
		"url": fmt.Sprintf("%s/v1/webhooks/payment-callback", appURL),
	}
	
	// Register this webhook URL with external service (payment gateway, etc.)
	registerWebhookWithPaymentGateway(webhook["url"].(string))
	
	c.JSON(200, gin.H{
		"success": true,
		"data":    webhook,
	})
}

// Example 4: File Download URLs
func GetFileHandler(c *gin.Context) {
	appURL, _ := c.Get("app_url")
	fileID := c.Param("id")
	
	file := getFileFromDB(fileID)
	
	response := gin.H{
		"success": true,
		"data": gin.H{
			"file": file,
			"download_url": fmt.Sprintf("%s/v1/files/%s/download", appURL, fileID),
			"preview_url":  fmt.Sprintf("%s/v1/files/%s/preview", appURL, fileID),
		},
	}
	c.JSON(200, response)
}

// Example 5: Pagination Links (HATEOAS style)
func ListStudentsWithLinks(c *gin.Context) {
	appURL, _ := c.Get("app_url")
	
	// Your existing pagination logic
	students := getStudentsFromDB()
	currentPage := 1
	totalPages := 10
	
	// Build pagination links
	links := gin.H{
		"self":  fmt.Sprintf("%s/v1/students?page=%d", appURL, currentPage),
		"first": fmt.Sprintf("%s/v1/students?page=1", appURL),
		"last":  fmt.Sprintf("%s/v1/students?page=%d", appURL, totalPages),
	}
	
	if currentPage > 1 {
		links["prev"] = fmt.Sprintf("%s/v1/students?page=%d", appURL, currentPage-1)
	}
	if currentPage < totalPages {
		links["next"] = fmt.Sprintf("%s/v1/students?page=%d", appURL, currentPage+1)
	}
	
	c.JSON(200, gin.H{
		"success": true,
		"data":    students,
		"links":   links,
	})
}

// Helper functions (implement these based on your needs)
func getStudentFromDB() interface{} { return nil }
func getFileFromDB(id string) interface{} { return nil }
func getStudentsFromDB() interface{} { return nil }
func sendEmail(to, subject, body string) {}
func registerWebhookWithPaymentGateway(url string) {}
