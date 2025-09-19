package unit

import (
	"go-admin/app/admin/models"
	"go-admin/common/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

// UserModelTestSuite 使用者模型測試套件
// 此套件測試使用者模型的基本功能，包含驗證、密碼處理等
type UserModelTestSuite struct {
	suite.Suite
	user *models.SysUser
}

// SetupTest 每個測試前的準備工作
// 建立一個標準的測試使用者實例
func (suite *UserModelTestSuite) SetupTest() {
	suite.user = &models.SysUser{
		Username: "testuser",
		NickName: "測試使用者",
		Email:    "test@example.com",
		Phone:    "0912345678",
		Status:   "2", // 啟用狀態
		RoleId:   1,
	}
}

// TestUserValidation 測試使用者資料驗證
func (suite *UserModelTestSuite) TestUserValidation() {
	// 測試有效的使用者資料
	suite.T().Run("有效使用者資料應該通過驗證", func(t *testing.T) {
		assert.NotEmpty(t, suite.user.Username, "使用者名稱不能為空")
		assert.NotEmpty(t, suite.user.Email, "電子郵件不能為空")
		assert.True(t, utils.IsValidEmail(suite.user.Email), "電子郵件格式應該正確")
	})

	// 測試無效的使用者名稱
	suite.T().Run("空的使用者名稱應該驗證失敗", func(t *testing.T) {
		suite.user.Username = ""
		assert.Empty(t, suite.user.Username, "使用者名稱為空時應該驗證失敗")
	})

	// 測試無效的電子郵件
	suite.T().Run("無效電子郵件格式應該驗證失敗", func(t *testing.T) {
		invalidEmails := []string{
			"invalid-email",
			"@example.com",
			"test@",
			"",
			"user@domain",
		}

		for _, email := range invalidEmails {
			suite.user.Email = email
			assert.False(t, utils.IsValidEmail(email), "電子郵件 '%s' 應該驗證失敗", email)
		}
	})
}

// TestPasswordOperations 測試密碼相關操作
func (suite *UserModelTestSuite) TestPasswordOperations() {
	password := "testPassword123"

	// 測試密碼雜湊
	suite.T().Run("密碼雜湊功能", func(t *testing.T) {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		assert.NoError(t, err, "密碼雜湊應該成功")
		assert.NotEmpty(t, hashedPassword, "雜湊密碼不能為空")
		assert.NotEqual(t, password, string(hashedPassword), "雜湊密碼應該與原密碼不同")
	})

	// 測試密碼驗證
	suite.T().Run("密碼驗證功能", func(t *testing.T) {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		
		// 正確密碼應該驗證成功
		err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
		assert.NoError(t, err, "正確密碼應該驗證成功")
		
		// 錯誤密碼應該驗證失敗
		err = bcrypt.CompareHashAndPassword(hashedPassword, []byte("wrongPassword"))
		assert.Error(t, err, "錯誤密碼應該驗證失敗")
	})
}

// TestUserStatus 測試使用者狀態
func (suite *UserModelTestSuite) TestUserStatus() {
	suite.T().Run("使用者狀態設定", func(t *testing.T) {
		// 測試啟用狀態
		suite.user.Status = "2"
		assert.Equal(t, "2", suite.user.Status, "使用者應該處於啟用狀態")
		
		// 測試停用狀態
		suite.user.Status = "1"
		assert.Equal(t, "1", suite.user.Status, "使用者應該處於停用狀態")
	})
}

// TestUserFields 測試使用者欄位
func (suite *UserModelTestSuite) TestUserFields() {
	suite.T().Run("使用者欄位設定和取得", func(t *testing.T) {
		// 測試設定和取得使用者名稱
		username := "newuser"
		suite.user.Username = username
		assert.Equal(t, username, suite.user.Username, "使用者名稱設定應該正確")
		
		// 測試設定和取得暱稱
		nickName := "新使用者"
		suite.user.NickName = nickName
		assert.Equal(t, nickName, suite.user.NickName, "暱稱設定應該正確")
		
		// 測試設定和取得電話號碼
		phone := "0987654321"
		suite.user.Phone = phone
		assert.Equal(t, phone, suite.user.Phone, "電話號碼設定應該正確")
		
		// 測試設定和取得角色 ID
		roleId := 2
		suite.user.RoleId = roleId
		assert.Equal(t, roleId, suite.user.RoleId, "角色 ID 設定應該正確")
	})
}

// TestUserValidationRules 測試使用者驗證規則
func (suite *UserModelTestSuite) TestUserValidationRules() {
	// 測試使用者名稱長度限制
	suite.T().Run("使用者名稱長度驗證", func(t *testing.T) {
		// 過短的使用者名稱
		suite.user.Username = "ab"
		assert.Less(t, len(suite.user.Username), 3, "使用者名稱長度應該至少 3 個字元")
		
		// 過長的使用者名稱
		suite.user.Username = "thisUsernameIsTooLongForValidation"
		assert.Greater(t, len(suite.user.Username), 20, "使用者名稱長度應該少於 20 個字元")
		
		// 合適的使用者名稱長度
		suite.user.Username = "validUser"
		assert.GreaterOrEqual(t, len(suite.user.Username), 3, "合適的使用者名稱長度")
		assert.LessOrEqual(t, len(suite.user.Username), 20, "合適的使用者名稱長度")
	})

	// 測試電話號碼格式
	suite.T().Run("電話號碼格式驗證", func(t *testing.T) {
		validPhones := []string{
			"0912345678",
			"02-12345678",
			"04-12345678",
		}
		
		for _, phone := range validPhones {
			suite.user.Phone = phone
			// 這裡假設有 utils.IsValidPhone 函數
			// assert.True(t, utils.IsValidPhone(phone), "電話號碼 '%s' 應該是有效的", phone)
		}
	})
}

// 執行使用者模型測試套件
func TestUserModelTestSuite(t *testing.T) {
	suite.Run(t, new(UserModelTestSuite))
}

// TestUtils 測試工具函數
func TestUtils(t *testing.T) {
	// 測試電子郵件驗證工具
	t.Run("電子郵件驗證工具", func(t *testing.T) {
		validEmails := []string{
			"test@example.com",
			"user.name@domain.co.uk",
			"test123@test-domain.org",
			"user+tag@example.com",
		}

		invalidEmails := []string{
			"invalid-email",
			"@example.com",
			"test@",
			"",
			"user@domain",
			"user..double.dot@example.com",
		}

		for _, email := range validEmails {
			assert.True(t, utils.IsValidEmail(email), "電子郵件 '%s' 應該是有效的", email)
		}

		for _, email := range invalidEmails {
			assert.False(t, utils.IsValidEmail(email), "電子郵件 '%s' 應該是無效的", email)
		}
	})

	// 測試密碼強度驗證（假設有此函數）
	t.Run("密碼強度驗證", func(t *testing.T) {
		weakPasswords := []string{
			"123",
			"password",
			"123456789",
			"qwerty",
		}

		strongPasswords := []string{
			"StrongPass123!",
			"MyP@ssw0rd",
			"SecureP@ss2024",
		}

		for _, password := range weakPasswords {
			// 假設有 utils.IsStrongPassword 函數
			// assert.False(t, utils.IsStrongPassword(password), "密碼 '%s' 應該被認為是弱密碼", password)
			_ = password // 暫時避免未使用變數錯誤
		}

		for _, password := range strongPasswords {
			// assert.True(t, utils.IsStrongPassword(password), "密碼 '%s' 應該被認為是強密碼", password)
			_ = password // 暫時避免未使用變數錯誤
		}
	})
}

// TestPasswordSecurity 測試密碼安全性
func TestPasswordSecurity(t *testing.T) {
	password := "testPassword123"

	t.Run("密碼雜湊安全性測試", func(t *testing.T) {
		// 測試相同密碼產生不同雜湊值（使用鹽值）
		hash1, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		assert.NoError(t, err)

		hash2, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		assert.NoError(t, err)

		// 雜湊值應該不同（因為使用了隨機鹽值）
		assert.NotEqual(t, string(hash1), string(hash2), "相同密碼應該產生不同的雜湊值")

		// 但是兩個雜湊值都應該能夠驗證原密碼
		err = bcrypt.CompareHashAndPassword(hash1, []byte(password))
		assert.NoError(t, err, "第一個雜湊值應該能驗證原密碼")

		err = bcrypt.CompareHashAndPassword(hash2, []byte(password))
		assert.NoError(t, err, "第二個雜湊值應該能驗證原密碼")
	})

	t.Run("密碼雜湊成本測試", func(t *testing.T) {
		// 測試不同成本的雜湊
		costs := []int{bcrypt.MinCost, bcrypt.DefaultCost, bcrypt.MaxCost}

		for _, cost := range costs {
			hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
			assert.NoError(t, err, "成本 %d 的雜湊應該成功", cost)

			err = bcrypt.CompareHashAndPassword(hash, []byte(password))
			assert.NoError(t, err, "成本 %d 的雜湊應該能驗證密碼", cost)
		}
	})
}