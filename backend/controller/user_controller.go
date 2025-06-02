package controller

import (
	"encoding/csv"
	"fmt"
	"go-gin-auth/dto"
	"go-gin-auth/model"
	"go-gin-auth/service"
	"go-gin-auth/utils"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("PPL-K4-2025")

func Register(c *gin.Context) {
	var userDTO dto.RegisterRequestDTO

	if err := c.ShouldBindJSON(&userDTO); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid input", err.Error(), nil)
		return
	}

	user, err := utils.ConvertDTOToUser(userDTO)

	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Error converting DTO to user", err.Error(), nil)
		return
	}

	if err := service.CreateUser(user); err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Failed to create user", err.Error(), nil)
		return
	}

	// Simpan audit log (INSERT)
	user.Password = "-" // Kosongkan password sebelum menyimpan ke log
	if err := service.LogAudit(utils.GetTableName(user),
		fmt.Sprint(user.ID), "INSERT", strconv.FormatUint(uint64(utils.GetCurrentUserID(c)), 10),
		nil, user, "User registered successfully."); err != nil {
		// utils.Respond(c, http.StatusInternalServerError, "Gagal mencatat log audit register:", err.Error(), nil)
		// return
		log.Println("Gagal mencatat log audit register:", err)
	}

	// Menambahkan aktivitas log setelah user berhasil didaftarkan
	// err = service.LogActivity(user.ID, user.FullName, "Register", "User registered successfully.", c)
	// if err != nil {
	// 	utils.Respond(c, http.StatusInternalServerError, "Failed to log activity", err.Error(), nil)
	// 	return
	// }

	utils.Respond(c, http.StatusCreated, "User registered successfully", nil, user)
}

func Login(c *gin.Context) {

	var input dto.LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid request", err.Error(), nil)
		return
	}

	user, err := service.GetUserByEmail(input.Email)
	if err != nil {
		utils.Respond(c, http.StatusUnauthorized, "Login failed", "User not found", nil)
		return
	}

	// Mengecek jika user tidak aktif
	if !user.Active {
		utils.Respond(c, http.StatusUnauthorized, "Login failed", "User is inactive", nil)
		return
	}

	// Cek password yang dimasukkan dengan password yang tersimpan
	match := service.VerifyPassword(input.Password, user.Password)
	if !match {
		err = service.HandleFailedLogin(&user)
		utils.Respond(c, http.StatusUnauthorized, "Login failed", err.Error(), nil)
		return
	}

	// Periksa apakah akun terkunci
	if user.LockedUntil.After(time.Now()) {
		remainingTime := time.Until(user.LockedUntil).Round(time.Minute)
		utils.Respond(c, http.StatusUnauthorized, "Login failed", "akun terkunci sementara, coba lagi setelah "+remainingTime.String(), nil)
		return
	}

	// generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"role":      user.Role, // ← penting untuk middleware
		"full_name": user.FullName,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Could not generate token", err.Error(), nil)
		return
	}

	// Catat log aktivitas login
	err = service.LogActivity(user.ID, user.FullName, "Login", "User logged in successfully", c)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Failed to log activity", err.Error(), nil)
		return
	}

	// Jika login berhasil, reset failed login attempts
	err = service.ResetFailedLoginAttempts(&user)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusOK, "Login successful", nil, gin.H{
		"token": tokenString,
	})
}

func Logout(c *gin.Context) {
	// // Biasanya di sisi frontend: hapus token dari storage.
	// c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
	// Ambil userID dari context atau token (misalnya dari JWT token atau session)
	var input dto.LogoutRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid request", err.Error(), nil)
		return
	}

	// Panggil service untuk logout
	err := service.LogoutUser(input.UserId)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusOK, "Logged out successfully", nil, nil)
}

func GetUsers(c *gin.Context) {
	// // Cek nilai di context
	// userID, userIDExists := c.Get("user_id")
	// if !userIDExists {
	// 	utils.Respond(c, http.StatusUnauthorized, "Unauthorized", "Missing user_id in context", nil)
	// 	return
	// }

	// // Cek nilai user_id
	// fmt.Printf("User ID: %v\n", userID)
	users, err := service.GetAllUsers()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Failed to get users", err.Error(), nil)
		return
	}
	// Format the response to include readable last login time
	var usersResponse []map[string]interface{}
	for _, user := range users {
		userMap := map[string]interface{}{
			"id":         user.ID,
			"email":      user.Email,
			"phone":      user.Phone,
			"full_name":  user.FullName,
			"role":       user.Role,
			"nip":        user.NIP,
			"active":     user.Active,
			"created_at": user.CreatedAt,
		}

		// Format last login time if it exists
		if !user.LastLoginAt.IsZero() {
			userMap["last_login"] = map[string]interface{}{
				"time":      user.LastLoginAt,
				"formatted": user.LastLoginAt.Format("02 Jan 2006, 15:04:05"),
				"relative":  utils.GetRelativeTimeString(user.LastLoginAt),
			}
		} else {
			userMap["last_login"] = nil
		}
		if !user.LastLogoutAt.IsZero() {
			userMap["last_logout"] = map[string]interface{}{
				"time":      user.LastLogoutAt,
				"formatted": user.LastLogoutAt.Format("02 Jan 2006, 15:04:05"),
				"relative":  utils.GetRelativeTimeString(user.LastLogoutAt),
			}
		} else {
			userMap["last_logout"] = nil
		}
		usersResponse = append(usersResponse, userMap)
	}
	utils.Respond(c, http.StatusOK, "Users retrieved successfully", nil, usersResponse)
}

func GetUser(c *gin.Context) {
	var id uint
	fmt.Sscanf(c.Param("id"), "%d", &id)
	user, err := service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func UpdateUserOri(c *gin.Context) {
	var id uint
	fmt.Sscanf(c.Param("id"), "%d", &id)
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	service.UpdateUser(id, user)
	c.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

func DeleteUserOri(c *gin.Context) {
	var id uint
	fmt.Sscanf(c.Param("id"), "%d", &id)
	service.DeleteUser(id)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func DeleteUser(c *gin.Context) {
	var id uint
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &id); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid ID parameter", err.Error(), nil)
		return
	}
	// Dapatkan informasi user yang akan dihapus
	targetUser, err := service.GetUserByID(id)
	if err != nil {
		utils.Respond(c, http.StatusNotFound, "User not found", err.Error(), nil)
		return
	}

	err = service.DeleteUser(targetUser.ID)
	if err != nil {
		utils.Respond(c, http.StatusNotFound, "User not found", err.Error(), nil)
		return
	}
	// Ambil user yang sedang login dari context (misalnya dari middleware)
	currentUserIDFloat, _ := c.Get("user_id")
	currentUserID := uint(currentUserIDFloat.(float64))
	currentFullName, _ := c.Get("full_name")

	fmt.Printf("Logged-in user: %v (ID %v)\n", currentFullName, currentUserID)

	// Catat log
	description := fmt.Sprintf("User %s (ID %d) deleted user %s (ID %d)",
		currentFullName, currentUserID, targetUser.FullName, targetUser.ID)

	_ = service.LogActivity(currentUserID, currentFullName.(string), "DeleteUser", description, c)

	utils.Respond(c, http.StatusOK, "User deleted successfully", nil, nil)
}

func UpdateUser(c *gin.Context) {
	var id uint
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &id); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid ID parameter", err.Error(), nil)
		return
	}

	var input model.User
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid request body", err.Error(), nil)
		return
	}

	// Cek apakah user dengan ID tersebut ada
	existingUser, err := service.GetUserByID(id)
	if err != nil {
		utils.Respond(c, http.StatusNotFound, "User not found", err.Error(), nil)
		return
	}
	// Simpan salinan before
	before := existingUser

	// Update field yang ingin diubah (hindari overwrite ID/Password langsung)
	existingUser.Email = input.Email
	existingUser.FullName = input.FullName
	existingUser.Role = input.Role
	existingUser.Phone = input.Phone

	if err := service.UpdateUser(id, existingUser); err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Failed to update user", err.Error(), nil)
		return
	}

	// Audit trail setelah update
	if err := service.LogAudit(utils.GetTableName(existingUser), fmt.Sprint(id), "UPDATE", strconv.FormatUint(uint64(utils.GetCurrentUserID(c)), 10), before, existingUser, "User updated successfully"); err != nil {
		log.Println("Gagal mencatat audit update user:", err)
	}

	utils.Respond(c, http.StatusOK, "User updated successfully", nil, existingUser)
}

func DeactivateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid user ID", nil, err.Error())
		return
	}

	// Ambil user terlebih dahulu
	user, err := service.GetUserByID(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusNotFound, "User not found", nil, err.Error())
		return
	}

	// Jika sudah nonaktif
	if !user.Active {
		utils.Respond(c, http.StatusBadRequest, "User is already deactivated", nil, nil)
		return
	}

	// Cek jumlah admin aktif lainnya
	activeAdminsCount, err := service.CountActiveAdmins()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Error", "Failed to count active admins", nil)
		return
	}

	// Jika hanya satu admin yang aktif, tolak permintaan
	if activeAdminsCount <= 1 {
		utils.Respond(c, http.StatusForbidden, "Forbidden", "You cannot deactivate your account, at least one admin must be active", nil)
		return
	}

	err = service.DeactivateUser(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusNotFound, "User not found", nil, err.Error())
		return
	}

	// // Tambahkan log aktivitas
	// service.CreateUserLog(uint(id), "Deactivate User")

	utils.Respond(c, http.StatusOK, "User deactivated successfully", nil, nil)
}
func ReactivateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid user ID", nil, err.Error())
		return
	}

	err = service.ReactivateUser(uint(id))
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "Failed to reactivate user", nil, err.Error())
		return
	}

	utils.Respond(c, http.StatusOK, "User reactivated successfully", nil, nil)
}

func SearchUsers(c *gin.Context) {
	filters := map[string]string{
		"full_name": c.Query("full_name"),
		"email":     c.Query("email"),
		"role":      c.Query("role"),
	}

	users, err := service.SearchUsers(filters)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "Search failed", nil, err.Error())
		return
	}

	utils.Respond(c, http.StatusOK, "Users fetched", nil, users)
}
func ResetUserPassword(c *gin.Context) {
	var id uint
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &id); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid ID parameter", err.Error(), nil)
		return
	}

	var body struct {
		NewPassword string `json:"new_password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Respond(c, http.StatusBadRequest, "Invalid request body", err.Error(), nil)
		return
	}

	if body.NewPassword == "" {
		utils.Respond(c, http.StatusBadRequest, "Password is required", "new_password is empty", nil)
		return
	}

	err := service.UpdateUserPassword(id, body.NewPassword)
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "Failed to update password", err.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusOK, "Password updated successfully", nil, nil)
}
func ExportUsersCSV(c *gin.Context) {
	// Ambil semua user dari service
	users, err := service.GetAllUsers()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "error", err.Error(), "Gagal mengambil data pengguna")
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengguna"})
		return
	}

	// Set header untuk download file CSV
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment;filename=users.csv")
	c.Header("Cache-Control", "no-cache")

	// Tulis CSV langsung ke response writer
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// // Header kolom CSV
	// writer.Write([]string{"Nama", "Email", "Role", "Status"})

	// // Loop semua data user dan tulis ke CSV
	// for _, user := range users {
	// 	writer.Write([]string{
	// 		user.FullName,
	// 		user.Email,
	// 		user.Role,
	// 		strconv.FormatBool(user.Active),
	// 	})
	// }
	// Gunakan refleksi untuk ambil nama field dari struct
	// Daftar field yang ingin di-skip
	skipFields := map[string]bool{
		"Password": true,
	}

	// Ambil header dari struct, kecuali yang di-skip
	userType := reflect.TypeOf(users[0])
	var headers []string
	var fieldIndexes []int
	for i := 0; i < userType.NumField(); i++ {
		fieldName := userType.Field(i).Name
		if skipFields[fieldName] {
			continue
		}
		headers = append(headers, fieldName)
		fieldIndexes = append(fieldIndexes, i)
	}
	writer.Write(headers)

	// Tulis isi baris data
	for _, user := range users {
		val := reflect.ValueOf(user)
		var row []string
		for _, idx := range fieldIndexes {
			field := val.Field(idx)
			row = append(row, fmt.Sprint(field.Interface()))
		}
		writer.Write(row)
	}

}
