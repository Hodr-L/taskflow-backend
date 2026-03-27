package services

import (
	"errors"
	"regexp"

	"taskflow-backend/internal/models"

	"gorm.io/gorm"
)

var (
	ErrUserExists         = errors.New("鐢ㄦ埛宸插瓨鍦?)
	ErrUserNotFound       = errors.New("鐢ㄦ埛涓嶅瓨鍦?)
	ErrInvalidCredentials = errors.New("鏃犳晥鐨勫嚟鎹?)
	ErrInvalidInput       = errors.New("鏃犳晥鐨勮緭鍏?)
	ErrUserInactive       = errors.New("鐢ㄦ埛璐︽埛鏈縺娲?)
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

// Register 娉ㄥ唽鏂扮敤鎴?func (s *AuthService) Register(username, email, password string) (*models.User, error) {
	// 楠岃瘉杈撳叆
	if err := validateUserInput(username, email, password); err != nil {
		return nil, ErrInvalidInput
	}

	// 妫€鏌ョ敤鎴锋槸鍚﹀凡瀛樺湪
	var count int64
	s.db.Model(&models.User{}).Where("username = ? OR email = ?", username, email).Count(&count)
	if count > 0 {
		return nil, ErrUserExists
	}

	// 鍒涘缓鐢ㄦ埛
	user := &models.User{
		Username: username,
		Email:    email,
		Role:     "user",
		Status:   "active",
	}

	// 璁剧疆瀵嗙爜
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	// 淇濆瓨鍒版暟鎹簱
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Login 鐢ㄦ埛鐧诲綍
func (s *AuthService) Login(email, password string) (*models.User, error) {
	// 鏌ユ壘鐢ㄦ埛
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 妫€鏌ョ敤鎴风姸鎬?	if !user.IsActive() {
		return nil, ErrUserInactive
	}

	// 楠岃瘉瀵嗙爜
	if !user.CheckPassword(password) {
		return nil, ErrInvalidCredentials
	}

	return &user, nil
}

// validateUserInput 楠岃瘉鐢ㄦ埛杈撳叆
func validateUserInput(username, email, password string) error {
	// 楠岃瘉鐢ㄦ埛鍚?	if len(username) < 3 || len(username) > 50 {
		return errors.New("鐢ㄦ埛鍚嶉暱搴﹀繀椤诲湪3-50涓瓧绗︿箣闂?)
	}

	// 楠岃瘉閭鏍煎紡
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("閭鏍煎紡涓嶆纭?)
	}

	// 楠岃瘉瀵嗙爜
	if len(password) < 6 {
		return errors.New("瀵嗙爜闀垮害鑷冲皯6涓瓧绗?)
	}

	return nil
}
