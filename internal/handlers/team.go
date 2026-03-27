package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"taskflow-backend/internal/services"
)

type TeamHandler struct {
	teamService *services.TeamService
	db          *gorm.DB
}

func NewTeamHandler(db *gorm.DB) *TeamHandler {
	return &TeamHandler{
		teamService: services.NewTeamService(db),
		db:          db,
	}
}

// CreateTeam 鍒涘缓鍥㈤槦w
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	// TODO: 瀹炵幇鍒涘缓鍥㈤槦閫昏緫
}

func (h *TeamHandler) GetListTeams(c *gin.Context) {
	// TODO: 瀹炵幇鑾峰彇鍥㈤槦鍒楄〃閫昏緫
}

// GetTeams 鑾峰彇鍥㈤槦鍒楄〃
func (h *TeamHandler) GetTeams(c *gin.Context) {
	// TODO: 鑾峰彇鐢ㄦ埛鍙備笌鐨勫洟闃熷垪琛?}

// GetTeam 鑾峰彇鍥㈤槦璇︽儏
func (h *TeamHandler) GetTeam(c *gin.Context) {
	// TODO: 瀹炵幇鑾峰彇鍥㈤槦璇︽儏閫昏緫
}

// UpdateTeam 鏇存柊鍥㈤槦
func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	// TODO: 瀹炵幇鏇存柊鍥㈤槦閫昏緫
}

// DeleteTeam 鍒犻櫎鍥㈤槦
func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	// TODO: 瀹炵幇鍒犻櫎鍥㈤槦閫昏緫
}
