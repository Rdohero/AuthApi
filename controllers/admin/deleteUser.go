package admin

import (
	"AuthApi/initializers"
	"AuthApi/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	result := initializers.DB.Where("id = ?", id).Find(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Terjadi kesalahan dalam mencari pengguna.",
		})
		return
	}

	os.Remove(user.Foto)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "User yang ingin dihapus tidak ditemukan",
		})
		return
	}

	initializers.DB.Where("id = ?", id).Delete(&user)

	c.JSON(http.StatusOK, gin.H{
		"Succes": "User telah terhapus",
	})
}
