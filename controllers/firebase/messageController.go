package firebase

import (
	"AuthApi/initializers"
	"firebase.google.com/go/v4/messaging"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Message(c *gin.Context) {
	var messages struct {
		Title    string `json:"title"`
		Body     string `json:"body"`
		ImageURL string `json:"image_url"`
		Token    string `json:"token"`
	}

	if err := c.BindJSON(&messages); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	client, err := initializers.FB.Messaging(initializers.Ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title:    messages.Title,
			Body:     messages.Body,
			ImageURL: messages.ImageURL,
		},
		Token: messages.Token,
	}

	_, err = client.Send(initializers.Ctx, message)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Succes": "Succes Send Notification",
	})
}

func MessageArray(c *gin.Context) {
	type _messages struct {
		Title    string `json:"title"`
		Body     string `json:"body"`
		ImageURL string `json:"image_url"`
		Token    string `json:"token"`
	}

	var messages struct {
		MessageData []_messages `json:"message_data"`
	}

	if err := c.BindJSON(&messages); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	client, err := initializers.FB.Messaging(initializers.Ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	for _, item := range messages.MessageData {
		message := &messaging.Message{
			Notification: &messaging.Notification{
				Title:    item.Title,
				Body:     item.Body,
				ImageURL: item.ImageURL,
			},
			Android: &messaging.AndroidConfig{
				Notification: &messaging.AndroidNotification{
					Sound: "default",
				},
			},
			Token: item.Token,
			APNS: &messaging.APNSConfig{
				Payload: &messaging.APNSPayload{
					Aps: &messaging.Aps{
						Sound: "default",
					},
				},
			},
		}

		_, err = client.Send(initializers.Ctx, message)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"Succes": "Succes Send Notification",
	})
}
