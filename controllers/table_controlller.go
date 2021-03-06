package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"w4s/models"
)

func CreateTable(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input models.TableInput
	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var table models.Table
	if db.Where("name = ?", input.Name).Find(&table).RecordNotFound() {
		var user2 models.User
		if err := db.Preload("Profile").Where("email = ? and deleted = ? ", c.Query("e"), false).First(&user2).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "não encontrado o nickname",
			})
			return
		}

		if len(input.Name) >= 20 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Nome muito grande"})
			return
		}

		table.Name = input.Name
		if len(input.Description) >= 360 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Descricao muito grande"})
			return
		}

		table.Description = input.Description
		table.NumberOfParticipants = 1
		/*
			if len(input.Thumbnail) <= 0 {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "sem imagem"})
				return
			}
		*/
		table.Thumbnail = input.Thumbnail
		table.MaxOfParticipants = input.MaxOfParticipants
		if len(input.Links) >= 255 || len(input.Links) == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "link invalido"})
			return
		}
		if len(input.RpgSystem) <= 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "sem sistema de rpg"})
			return
		}
		table.RpgSystem = input.RpgSystem
		table.Links = input.Links
		table.Privacy = input.Privacy
		if err := db.Create(&table).Error; err != nil { //Return the error by JSON / Retornando o erro por JSON
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		tablePermission, err := userPermissionCreate(c, "1", user2.ProfileID, table.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Interno"})
			return
		}
		userAndPermissonAppend(c, table, tablePermission, user2.Profile)
		c.JSON(http.StatusOK, gin.H{"success": "table created"})
		return
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "table name already exists"})
	return
}
func UserJoinTable(c *gin.Context) {
	//Empty parametrs error checking
	if c.Query("nickname") == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user not inform"})
		return
	}
	if c.Query("table") == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "table not inform"})
		return
	}
	//=========================
	db := c.MustGet("db").(*gorm.DB)
	var userTobeAdd models.Profile
	if err := db.Where("nickname = ? AND deleted = ?", c.Query("nickname"), false).Find(&userTobeAdd).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Nenhum registro encontrado",
		})
		return
	}
	var table models.Table
	if err := db.Where("name = ?", c.Query("table")).Preload("User").Find(&table).Error; err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Nenhum registro encontrado",
		})
		return
	}

	for i := 0; i < len(table.User); i++ {
		if table.User[i].ID == userTobeAdd.ID {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User already is in the table"})
			return
		}
	}
	if table.NumberOfParticipants != table.MaxOfParticipants {
		//.Where("name = ? ", c.Query("table"))
		/*
			p := c.Query("p")
			if p == "" || p == "0" {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "permissao invalida"})
				return
			}
		*/
		tablePermission, err := userPermissionCreate(c, "3", userTobeAdd.ID, table.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internos"})
			return
		}
		userAndPermissonAppend(c, table, tablePermission, userTobeAdd)
		db.Model(&table).Update("number_of_participants", table.NumberOfParticipants+1)
		c.JSON(http.StatusOK, gin.H{"success": "join in the table"})
		return
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "table full"})
	return
}
func FindAllTables(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var tables []models.Table

	if err := db.Preload("User").Preload("Permitions").Find(&tables).Error; err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Nenhum registro encontrado",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": tables,
	})
	return
}

func FindAllUserTables(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var tables []models.Profile
	if err := db.Where("nickname = ?", c.Query("nickname")).Preload("Tables").Find(&tables).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Nenhuma mesa encontrada! ",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": tables,
	})
	return
}

func FindOneTables(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var tables []models.Table
	id := c.Params.ByName("id")

	if err := db.Preload("User").Preload("Permitions").Where("id = ?", id).Find(&tables).Error; err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Nenhum registro encontrado",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": tables,
	})
	return
}

func UpdateTable(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var table models.Table
	id := c.Query("id")

	if err := db.Preload("User").Preload("Permitions", "permission NOT IN ('3')").Where("id = ?", id).First(&table).Error; err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Nenhum registro encontrado",
		})
		return
	}
	var profile models.Profile
	for _, permission := range table.Permitions {
		db.Where("id = ?", permission.Permission).Find(&profile)
		if profile.Nickname != c.Query("nickname") {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "não tem permissão de administrador"})
			return
		}
	}
	if err := c.Bind(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	db.Model(&table).Updates(table)
	c.JSON(http.StatusOK, gin.H{
		"success": table,
	})
	return
}

func DeleteTable(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var table models.Table
	id := c.Query("id")

	err := db.Where("id = ?", id).Preload("User").Preload("Permitions", "permission NOT IN ('2,3')").Find(&table).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var profile models.Profile
	for _, permission := range table.Permitions {
		db.Where("id = ?", permission.Permission).Find(&profile)
		if profile.Nickname != c.Query("nickname") {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "não tem permissão de administrador"})
			return
		}
	}
	db.Unscoped().Delete(&table)
	c.JSON(http.StatusOK, gin.H{
		"success": "deleted",
	})
	return
}

// referencia
// https://medium.com/@cgrant/developing-a-simple-crud-api-with-go-gin-and-gorm-df87d98e6ed1

func userAndPermissonAppend(c *gin.Context, table models.Table, tablePermission models.PermissionTable, user models.Profile) {
	db := c.MustGet("db").(*gorm.DB)
	db.Model(&table).Association("Permitions").Append([]*models.PermissionTable{&tablePermission})
	db.Model(&table).Association("User").Append([]*models.Profile{&user})
}
func userPermissionCreate(c *gin.Context, permission string, userID uint, tableID uint) (models.PermissionTable, error) {
	db := c.MustGet("db").(*gorm.DB)

	tablePermission := models.PermissionTable{
		Permission:      permission,
		ProfileNickname: userID,
		TableId:         tableID,
	}
	if err := db.Create(&tablePermission).Error; err != nil { //Return the error by JSON / Retornando o erro por JSON
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return tablePermission, err
	}
	return tablePermission, nil
}

/*func insertPictures(c *gin.Context, TableId uint) {
	db := c.MustGet("db").(*gorm.db)
	var pictures models.Picture
	if err := c.BindJSON(pictures); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	pictures.TableID = TableId
	split := strings.Split(pictures.PictureFile, " ")
	for i := 0; i < len(split); i++ {
		if err := db.Create(&pictures).Error; err != nil { //Return the error by JSON / Retornando o erro por JSON
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
	}
	return
}
*/
