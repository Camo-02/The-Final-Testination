package functionality

import (
	"backend/database"
	"backend/database/entity"
	"backend/utils"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AvailableLevelDTO struct {
	GameId        string `json:"game_id"`
	Title         string `json:"title"`
	GameOrder     int    `json:"game_order"`
	MaxScore      int    `json:"max_score"`
	ScoreAchieved int    `json:"score_achieved"`
	Description   string `json:"description"`
}

type LevelProgress struct {
	Title                  string     `json:"title"`
	Score                  int        `json:"score"`
	StartTime              time.Time  `json:"start_time"`
	EndTime                *time.Time `json:"end_time"`
	TextualHintPointsUsed  int        `json:"textual_hint_points_used"`
	HintSolutionPointsUsed int        `json:"hint_solution_points_used"`
	TimeFreezePointsUsed   int        `json:"time_freeze_points_used"`
}
type ProfileDTO struct {
	ProfileImage string          `json:"profileImage"`
	Username     string          `json:"username"`
	Email        string          `json:"email"`
	Levels       []LevelProgress `json:"levels"`
}

type IconDTO struct {
	Id  string `json:"id"`
	Svg string `json:"svg"`
}

func PlayerCreate(database *database.FinalTestinationDB, username string, password string, email string) (*entity.Player, error) {
	hashedPassword, err := utils.GenerateHash(password)
	if err != nil {
		return nil, err
	}
	player := entity.Player{
		Model: utils.Model{
			ID: uuid.New().String(),
		},
		Username: username,
		Password: hashedPassword,
		Email:    email,
		IconID:   "1",
		Secure:   false,
		SameSite: "None",
	}

	result := database.Orm.Create(&player)

	if result.Error != nil {
		return nil, result.Error
	}

	return &player, nil
}

func PlayerGetByID(database *database.FinalTestinationDB, id string) (*entity.Player, error) {
	var player entity.Player
	result := database.Orm.Where("id = ?", id).First(&player)

	return &player, result.Error
}

// Google does not provide the password, so we need to check the email
func PlayerGetByEmail(database *database.FinalTestinationDB, email string) (*entity.Player, error) {
	var player entity.Player
	result := database.Orm.Where("email = ?", email).First(&player)

	// Checks if the error is "record not found"
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil for the player and nil for the error
		}
		return nil, result.Error // Return other errors
	}

	return &player, nil
}

func PlayerUpdate(database *database.FinalTestinationDB, player *entity.Player) error {
	result := database.Orm.Save(player)
	return result.Error
}

func PlayerGetByUsernameAndPassword(database *database.FinalTestinationDB, username string, password string) (*entity.Player, error) {
	var player entity.Player
	result := database.Orm.Where("username = ?", username).First(&player)

	if result.Error != nil {
		return nil, result.Error
	}

	if !utils.CompareHash(player.Password, password) {
		return nil, errors.New("user not found")
	}

	return &player, nil
}

func PlayerGetByEmailAndPassword(database *database.FinalTestinationDB, email string, password string) (*entity.Player, error) {
	var player entity.Player
	result := database.Orm.Where("email = ?", email).First(&player)

	if result.Error != nil {
		return nil, result.Error
	}

	if !utils.CompareHash(player.Password, password) {
		return nil, errors.New("user not found")
	}

	return &player, nil
}

func GetPlayerLevels(database *database.FinalTestinationDB, id string) ([]AvailableLevelDTO, error) {
	// TODO: get played levels with related info
	var res []AvailableLevelDTO
	var res2 []AvailableLevelDTO

	// TODO: query can be optimized

	// get already completed levels (max_score will be 0)
	result := database.Orm.Table("\"player_games\" AS pg").
		Joins("JOIN games ON pg.game_id = games.id").
		Where("pg.player_id = ? AND pg.end_time IS NOT NULL", id).
		Select("game_id, title, game_order, max_score, pg.score as score_achieved, description").
		Order("game_order").
		Scan(&res)

	if result.Error != nil {
		return nil, result.Error
	}

	nextLevelOrder := 1
	// If the player has already completed at least one level,
	// the order of the next level is the order of the last completed level plus one.
	// Otherwise, the order is 1
	if len(res) > 0 {
		nextLevelOrder = res[len(res)-1].GameOrder + 1
	}

	// get next available level (normal) and those beyond (locked, negative score)
	//TODO: check that failed attempts are registered, then change attempts for the level to be completed
	//TODO: the frontend uses max_score as part of the logic to "display available levels", if the real max_score is sent then the frontend needs to be updated as well
	result2 := database.Orm.Raw("SELECT * FROM (? UNION ?) AS T ORDER BY T.game_order",
		database.Orm.
			Table("\"games\" AS g").
			Where("g.game_order = ?", nextLevelOrder).
			Select("id AS game_id, title, game_order, max_score, 0 as score_achieved, description"),
		database.Orm.
			Table("\"games\" AS g").
			Where("g.game_order > ?", nextLevelOrder).
			Select("id AS game_id, title, game_order, -1 AS max_score, 0 as score_achieved, 'locked' AS description"),
	).Scan(&res2)

	if result2.Error != nil {
		return nil, result2.Error
	}

	finalRes := append(res, res2...)

	return finalRes, nil
}

func Profile(database *database.FinalTestinationDB, player *entity.Player) (*ProfileDTO, error) {
	var levelProgress []LevelProgress
	res := database.Orm.Table("\"player_games\" AS pg").
		Joins("JOIN games ON pg.game_id = games.id ").
		Where("pg.player_id = ?", player.ID).
		Select("title,score, start_time,end_time,textual_hint_points_used,hint_solution_points_used,time_freeze_points_used").
		Order("pg.end_time DESC").
		Scan(&levelProgress)

	if res.Error != nil {
		return nil, res.Error
	}

	var propic string
	res = database.Orm.Table("icons").
		Joins("JOIN players ON icons.id = players.icon_id ").
		Where("players.id = ?", player.ID).
		Select("svg").
		Scan(&propic)

	if res.Error != nil {
		return nil, res.Error
	}

	var profile ProfileDTO
	profile.Levels = levelProgress
	profile.Email = player.Email
	profile.Username = player.Username
	profile.ProfileImage = propic

	return &profile, res.Error
}

func GetAvailableIcons(database *database.FinalTestinationDB) ([]IconDTO, error) {

	var icons []IconDTO

	res := database.Orm.Table("icons").
		Select("id, svg").
		Scan(&icons)

	if res.Error != nil {
		return nil, res.Error
	}

	return icons, nil
}

func ChangeIcon(database *database.FinalTestinationDB, player string, icon string) error {
	res := database.Orm.Model(&entity.Player{}).
		Where("id = ?", player).
		Update("icon_id", icon)
	return res.Error
}
